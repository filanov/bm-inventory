package bminventory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"text/template"

	"github.com/filanov/bm-inventory/internal/cluster"
	"github.com/filanov/bm-inventory/internal/host"
	"github.com/filanov/bm-inventory/internal/installcfg"
	"github.com/filanov/bm-inventory/models"
	"github.com/filanov/bm-inventory/pkg/filemiddleware"
	"github.com/filanov/bm-inventory/pkg/job"
	logutil "github.com/filanov/bm-inventory/pkg/log"
	"github.com/filanov/bm-inventory/restapi/operations/installer"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const baseHref = "/api/bm-inventory/v1"
const kubeconfigPrefix = "generate-kubeconfig"

const defaultJobNamespace = "default"

const (
	ClusterStatusReady      = "ready"
	ClusterStatusInstalling = "installing"
	ClusterStatusInstalled  = "installed"
	ClusterStatusError      = "error"
)

const (
	ResourceKindHost    = "Host"
	ResourceKindCluster = "Cluster"
)

const (
	bootstrap = "bootstrap"
)

type Config struct {
	ImageBuilder        string `envconfig:"IMAGE_BUILDER" default:"quay.io/oscohen/installer-image-build"`
	ImageBuilderCmd     string `envconfig:"IMAGE_BUILDER_CMD" default:"echo hello"`
	AgentDockerImg      string `envconfig:"AGENT_DOCKER_IMAGE" default:"quay.io/oamizur/agent:latest"`
	KubeconfigGenerator string `envconfig:"KUBECONFIG_GENERATE_IMAGE" default:"quay.io/oscohen/ignition-manifests-and-kubeconfig-generate"`
	InventoryURL        string `envconfig:"INVENTORY_URL" default:"10.35.59.36"`
	InventoryPort       string `envconfig:"INVENTORY_PORT" default:"30485"`
	S3EndpointURL       string `envconfig:"S3_ENDPOINT_URL" default:"http://10.35.59.36:30925"`
	S3Bucket            string `envconfig:"S3_BUCKET" default:"test"`
	AwsAccessKeyID      string `envconfig:"AWS_ACCESS_KEY_ID" default:"accessKey1"`
	AwsSecretAccessKey  string `envconfig:"AWS_SECRET_ACCESS_KEY" default:"verySecretKey1"`
	InstallerImage      string `envconfig:"INSTALLER_IMAGE" default:"quay.io/ocpmetal/assisted-installer:stable"`
}

const ignitionConfigFormat = `{
"ignition": { "version": "3.0.0" },
  "passwd": {
    "users": [
      {
        "groups": [
          "sudo",
          "docker"
        ],
        "name": "core",
        "passwordHash": "$6$MWO4bibU8TIWG0XV$Hiuj40lWW7pHiwJmXA8MehuBhdxSswLgvGxEh8ByEzeX2D1dk87JILVUYS4JQOP45bxHRegAB9Fs/SWfszXa5."
      }
	 {{.userSshKey}}
    ]
  },
"systemd": {
"units": [{
"name": "agent.service",
"enabled": true,
"contents": "[Service]\nType=simple\nEnvironment=HTTPS_PROXY={{.ProxyURL}}\nEnvironment=HTTP_PROXY={{.ProxyURL}}\nEnvironment=http_proxy={{.ProxyURL}}\nEnvironment=https_proxy={{.ProxyURL}}\nExecStartPre=docker run --privileged --rm -v /usr/local/bin:/hostbin {{.AgentDockerImg}} cp /usr/bin/agent /hostbin\nExecStart=/usr/local/bin/agent --host {{.InventoryURL}} --port {{.InventoryPort}} --cluster-id {{.clusterId}}\n\n[Install]\nWantedBy=multi-user.target"
}]
}
}`

type debugCmd struct {
	cmd    string
	stepID string
}

type bareMetalInventory struct {
	Config
	imageBuildCmd []string
	db            *gorm.DB
	debugCmdMap   map[strfmt.UUID]debugCmd
	debugCmdMux   sync.Mutex
	log           logrus.FieldLogger
	job           job.API
	hostApi       host.API
	clusterApi    cluster.API
}

func NewBareMetalInventory(db *gorm.DB, log logrus.FieldLogger, hostApi host.API, clusterApi cluster.API, cfg Config,
	jobApi job.API) *bareMetalInventory {

	b := &bareMetalInventory{
		db:          db,
		log:         log,
		Config:      cfg,
		debugCmdMap: make(map[strfmt.UUID]debugCmd),
		hostApi:     hostApi,
		clusterApi:  clusterApi,
		job:         jobApi,
	}

	if cfg.ImageBuilderCmd != "" {
		b.imageBuildCmd = strings.Split(cfg.ImageBuilderCmd, " ")
	}
	return b
}

func buildHref(base, id string) *string {
	return swag.String(fmt.Sprintf("%s/%ss/%s", baseHref, base, id))
}

// create discovery image generation job, return job name and error
func (b *bareMetalInventory) createImageJob(cluster *models.Cluster, jobName, imgName, ignitionConfig string) *batch.Job {
	return &batch.Job{
		TypeMeta: meta.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      jobName,
			Namespace: "default",
		},
		Spec: batch.JobSpec{
			BackoffLimit: swag.Int32(2),
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Name:      jobName,
					Namespace: "default",
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            "image-creator",
							Image:           b.Config.ImageBuilder,
							Command:         b.imageBuildCmd,
							ImagePullPolicy: "IfNotPresent",
							Env: []core.EnvVar{
								{
									Name:  "S3_ENDPOINT_URL",
									Value: b.S3EndpointURL,
								},
								{
									Name:  "IGNITION_CONFIG",
									Value: ignitionConfig,
								},
								{
									Name:  "IMAGE_NAME",
									Value: imgName,
								},
								{
									Name:  "S3_BUCKET",
									Value: b.S3Bucket,
								},
								{
									Name:  "aws_access_key_id",
									Value: b.AwsAccessKeyID,
								},
								{
									Name:  "aws_secret_access_key",
									Value: b.AwsSecretAccessKey,
								},
							},
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}
}

func (b *bareMetalInventory) formatIgnitionFile(cluster *models.Cluster, params installer.GenerateClusterISOParams) (string, error) {
	var ignitionParams = map[string]string{
		"userSshKey":     b.getUserSshKey(params),
		"AgentDockerImg": b.AgentDockerImg,
		"InventoryURL":   b.InventoryURL,
		"InventoryPort":  b.InventoryPort,
		"clusterId":      cluster.ID.String(),
		"ProxyURL":       params.ImageCreateParams.ProxyURL,
	}
	tmpl, err := template.New("ignitionConfig").Parse(ignitionConfigFormat)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, ignitionParams); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (b *bareMetalInventory) getUserSshKey(params installer.GenerateClusterISOParams) string {
	sshKey := params.ImageCreateParams.SSHPublicKey
	if sshKey == "" {
		return ""
	}
	return fmt.Sprintf(`,{
		"name": "systemUser",
		"passwordHash": "$6$MWO4bibU8TIWG0XV$Hiuj40lWW7pHiwJmXA8MehuBhdxSswLgvGxEh8ByEzeX2D1dk87JILVUYS4JQOP45bxHRegAB9Fs/SWfszXa5.",
		"sshAuthorizedKeys": [
		"%s"],
		"groups": [ "sudo" ]}`, sshKey)
}

func (b *bareMetalInventory) RegisterCluster(ctx context.Context, params installer.RegisterClusterParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	id := strfmt.UUID(uuid.New().String())
	log.Infof("Register cluster: %s with id %s", swag.StringValue(params.NewClusterParams.Name), id)
	cluster := models.Cluster{
		ID:                       &id,
		Href:                     buildHref(ResourceKindCluster, id.String()),
		Kind:                     swag.String(ResourceKindCluster),
		APIVip:                   params.NewClusterParams.APIVip,
		BaseDNSDomain:            params.NewClusterParams.BaseDNSDomain,
		ClusterNetworkCidr:       params.NewClusterParams.ClusterNetworkCidr,
		ClusterNetworkHostPrefix: params.NewClusterParams.ClusterNetworkHostPrefix,
		DNSVip:                   params.NewClusterParams.DNSVip,
		IngressVip:               params.NewClusterParams.IngressVip,
		Name:                     swag.StringValue(params.NewClusterParams.Name),
		OpenshiftVersion:         swag.StringValue(params.NewClusterParams.OpenshiftVersion),
		PullSecret:               params.NewClusterParams.PullSecret,
		ServiceNetworkCidr:       params.NewClusterParams.ServiceNetworkCidr,
		SSHPublicKey:             params.NewClusterParams.SSHPublicKey,
		UpdatedAt:                strfmt.DateTime{},
	}

	err := b.clusterApi.RegisterCluster(ctx, &cluster)
	if err != nil {
		log.Errorf("failed to register cluster %s ", swag.StringValue(params.NewClusterParams.Name))
		return installer.NewRegisterClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	return installer.NewRegisterClusterCreated().WithPayload(&cluster)
}

func (b *bareMetalInventory) DeregisterCluster(ctx context.Context, params installer.DeregisterClusterParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var cluster models.Cluster

	if err := b.db.First(&cluster, "id = ?", params.ClusterID).Error; err != nil {
		return installer.NewDeregisterClusterNotFound().
			WithPayload(generateError(http.StatusNotFound))
	}

	err := b.clusterApi.DeregisterCluster(ctx, &cluster)
	if err != nil {
		log.WithError(err).Errorf("failed to deregister cluster cluster %s", params.ClusterID)
		return installer.NewDeregisterClusterNotFound().
			WithPayload(generateError(http.StatusNotFound))
	}

	return installer.NewDeregisterClusterNoContent()
}

func (b *bareMetalInventory) DownloadClusterISO(ctx context.Context, params installer.DownloadClusterISOParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	if err := b.db.First(&models.Cluster{}, "id = ?", params.ClusterID).Error; err != nil {
		log.WithError(err).Errorf("failed to get cluster %s", params.ClusterID)
		return installer.NewDownloadClusterISONotFound().
			WithPayload(generateError(http.StatusNotFound))
	}
	imgName := getImageName(params.ClusterID, params.ImageID)
	imageURL := fmt.Sprintf("%s/%s/%s", b.S3EndpointURL, b.S3Bucket, imgName)

	log.Info("Image URL: ", imageURL)
	resp, err := http.Get(imageURL)
	if err != nil {
		log.WithError(err).Errorf("Failed to get ISO: %s", imgName)
		return installer.NewDownloadClusterISOInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		log.WithError(fmt.Errorf("%d - %s", resp.StatusCode, string(b))).
			Errorf("Failed to get ISO: %s", imgName)
		if resp.StatusCode == http.StatusNotFound {
			return installer.NewDownloadClusterISONotFound().
				WithPayload(generateError(http.StatusNotFound))
		}
		return installer.NewDownloadClusterISOInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	return filemiddleware.NewResponder(installer.NewDownloadClusterISOOK().WithPayload(resp.Body),
		fmt.Sprintf("%s-cluster-%s-discovery.iso", params.ImageID.String(), params.ClusterID.String()))
}

// GenerateClusterISO and return image ID that can be used to download the ISO
func (b *bareMetalInventory) GenerateClusterISO(ctx context.Context, params installer.GenerateClusterISOParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	log.Infof("prepare image for cluster %s", params.ClusterID)
	var cluster models.Cluster
	if err := b.db.First(&cluster, "id = ?", params.ClusterID).Error; err != nil {
		log.WithError(err).Errorf("failed to get cluster %s", params.ClusterID)
		return installer.NewGenerateClusterISONotFound().
			WithPayload(generateError(http.StatusNotFound))
	}
	// generating a new uuid for each call to prevent races between concurrent requests
	imgId := strfmt.UUID(uuid.New().String())
	imgName := getImageName(params.ClusterID, imgId)
	// max job name is 63 chars
	jobName := fmt.Sprintf("create-image-%s-%s", cluster.ID, imgId)[:63]

	ignitionConfig, formatErr := b.formatIgnitionFile(&cluster, params)
	if formatErr != nil {
		log.WithError(formatErr).Errorf("failed to format ignition config file for cluster %s", cluster.ID)
		return installer.NewGenerateClusterISOInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	if err := b.job.Create(ctx, b.createImageJob(&cluster, jobName, imgName, ignitionConfig)); err != nil {
		log.WithError(err).Error("failed to create image job")
		return installer.NewGenerateClusterISOInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	if err := b.job.Monitor(ctx, jobName, defaultJobNamespace); err != nil {
		log.WithError(err).Error("image creation failed")
		return installer.NewGenerateClusterISOInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	log.Info("Generated cluster <%s> image <%s> with ignition config %s", params.ClusterID, imgId, ignitionConfig)
	return installer.NewGenerateClusterISOCreated().
		WithPayload(&installer.GenerateClusterISOCreatedBody{ImageID: imgId})
}

func getImageName(clusterID, id strfmt.UUID) string {
	return fmt.Sprintf("discovery-image-%s-%s", clusterID.String(), id)
}

func (b *bareMetalInventory) InstallCluster(ctx context.Context, params installer.InstallClusterParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var cluster models.Cluster

	tx := b.db.Begin()
	if tx.Error != nil {
		log.WithError(tx.Error).Errorf("failed to start db transaction")
		return installer.NewInstallClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}
	defer func() {
		if r := recover(); r != nil {
			log.Error("update cluster failed")
			tx.Rollback()
		}
	}()

	if err := tx.Preload("Hosts").First(&cluster, "id = ?", params.ClusterID).Error; err != nil {
		return installer.NewInstallClusterNotFound().
			WithPayload(generateError(http.StatusNotFound))
	}

	if err := b.clusterApi.Install(ctx, &cluster, tx); err != nil {
		log.WithError(err).Errorf("failed to install cluster %s", cluster.ID.String())
		tx.Rollback()
		return installer.NewInstallClusterConflict().WithPayload(generateError(http.StatusConflict))
	}

	// Temporary hack - use debug API for setting the executing install command:
	masterNodesIds, err := b.clusterApi.GetMasterNodesIds(ctx, &cluster, tx)
	if err != nil {
		tx.Rollback()
		log.WithError(err).Errorf("failed to get cluster %s master node id's", cluster.ID)
		return installer.NewInstallClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}
	if err = b.addInstallCommand(ctx, masterNodesIds, log, params, cluster); err != nil {
		tx.Rollback()
		log.WithError(err).Errorf("failed to add install command to cluster <%s>", params.ClusterID)
		return installer.NewInstallClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	// move hosts states to installing
	for i := range cluster.Hosts {
		if _, err = b.hostApi.Install(ctx, cluster.Hosts[i], tx); err != nil {
			log.WithError(err).Errorf("failed to install hosts <%s> in cluster: %s",
				cluster.Hosts[i].ID.String(), cluster.ID.String())
			tx.Rollback()
			return installer.NewInstallClusterConflict().WithPayload(generateError(http.StatusConflict))
		}
	}
	if err := b.generateClusterInstallConfig(ctx, cluster); err != nil {
		tx.Rollback()
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.WithError(err).Errorf("failed to commit cluster %s changes on installation", cluster.ID.String())
		return installer.NewInstallClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}
	if err := b.db.Preload("Hosts").First(&cluster, "id = ?", params.ClusterID).Error; err != nil {
		return installer.NewInstallClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}
	return installer.NewInstallClusterOK().WithPayload(&cluster)
}

func (b *bareMetalInventory) generateClusterInstallConfig(ctx context.Context, cluster models.Cluster) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)

	cfg, err := installcfg.GetInstallConfig(&cluster)
	if err != nil {
		log.WithError(err).Errorf("failed to get install config for cluster %s", cluster.ID)
		return installer.NewInstallClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}
	jobName := fmt.Sprintf("%s-%s-%s", kubeconfigPrefix, cluster.ID.String(), uuid.New().String())[:63]
	if err := b.job.Create(ctx, b.createKubeconfigJob(&cluster, jobName, cfg)); err != nil {
		log.WithError(err).Errorf("Failed to create kubeconfig generation job %s for cluster %s", jobName, cluster.ID)
		return installer.NewInstallClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	if err := b.job.Monitor(ctx, jobName, defaultJobNamespace); err != nil {
		log.WithError(err).Errorf("Generating kubeconfig files %s failed for cluster %s", jobName, cluster.ID)
		return installer.NewInstallClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}
	return nil
}

func (b *bareMetalInventory) addInstallCommand(ctx context.Context, masterNodesIds []*strfmt.UUID,
	log logrus.FieldLogger, params installer.InstallClusterParams, cluster models.Cluster) error {
	// set one of the master nodes as bootstrap
	bootstrapId := masterNodesIds[len(masterNodesIds)-1]
	log.Debugf("Bootstrap ID is %s", bootstrapId)

	const cmdTmpl = `sudo podman run -v /dev:/dev:rw -v /opt:/opt:rw --privileged --pid=host  {{.INSTALLER}} --role {{.ROLE}}  --cluster-id {{.CLUSTER_ID}}  --host {{.HOST}} --port {{.PORT}} --boot-device {{.BOOT_DEVICE}} --host-id {{.HOST_ID}}`

	t, err := template.New("cmd").Parse(cmdTmpl)
	if err != nil {
		return err
	}

	data := map[string]string{
		"HOST":        b.InventoryURL,
		"PORT":        b.InventoryPort,
		"CLUSTER_ID":  string(params.ClusterID),
		"ROLE":        "",
		"INSTALLER":   b.Config.InstallerImage,
		"BOOT_DEVICE": "",
	}
	for i := range cluster.Hosts {
		role := cluster.Hosts[i].Role
		if cluster.Hosts[i].ID == bootstrapId {
			role = bootstrap
		}
		data["ROLE"] = role
		disks, err := b.hostApi.GetHostValidDisks(cluster.Hosts[i])
		if err != nil {
			log.Errorf("Failed to get valid disks on host with id %s", cluster.Hosts[i].ID)
			return err
		}
		data["BOOT_DEVICE"] = fmt.Sprintf("/dev/%s", disks[0].Name)
		data["HOST_ID"] = string(*cluster.Hosts[i].ID)
		buf := &bytes.Buffer{}
		if err := t.Execute(buf, data); err != nil {
			return err
		}
		command := buf.String()
		b.SetDebugStep(ctx, installer.SetDebugStepParams{
			ClusterID: params.ClusterID,
			HostID:    *cluster.Hosts[i].ID,
			Step:      &models.DebugStep{Command: &command},
		})
	}
	return nil
}

func (b *bareMetalInventory) UpdateCluster(ctx context.Context, params installer.UpdateClusterParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var cluster models.Cluster
	log.Info("update cluster ", params.ClusterID)

	tx := b.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("update cluster failed")
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		log.WithError(tx.Error).Error("failed to start transaction")
	}

	if err := tx.First(&cluster, "id = ?", params.ClusterID).Error; err != nil {
		log.WithError(err).Errorf("failed to get cluster: %s", params.ClusterID)
		tx.Rollback()
		return installer.NewUpdateClusterNotFound().WithPayload(generateError(http.StatusNotFound))
	}

	cluster.Name = params.ClusterUpdateParams.Name
	cluster.APIVip = params.ClusterUpdateParams.APIVip
	cluster.BaseDNSDomain = params.ClusterUpdateParams.BaseDNSDomain
	cluster.ClusterNetworkCidr = params.ClusterUpdateParams.ClusterNetworkCidr
	cluster.ClusterNetworkHostPrefix = params.ClusterUpdateParams.ClusterNetworkHostPrefix
	cluster.DNSVip = params.ClusterUpdateParams.DNSVip
	cluster.IngressVip = params.ClusterUpdateParams.IngressVip
	cluster.PullSecret = params.ClusterUpdateParams.PullSecret
	cluster.ServiceNetworkCidr = params.ClusterUpdateParams.ServiceNetworkCidr
	cluster.SSHPublicKey = params.ClusterUpdateParams.SSHPublicKey

	if err := tx.Model(&cluster).Update(cluster).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Errorf("failed to update cluster: %s", params.ClusterID)
		return installer.NewUpdateClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	for i := range params.ClusterUpdateParams.HostsRoles {
		log.Infof("Update host %s to role: %s", params.ClusterUpdateParams.HostsRoles[i].ID,
			params.ClusterUpdateParams.HostsRoles[i].Role)
		var host models.Host
		if err := tx.First(&host, "id = ? and cluster_id = ?",
			params.ClusterUpdateParams.HostsRoles[i].ID, params.ClusterID).Error; err != nil {
			tx.Rollback()
			log.WithError(err).Errorf("failed to find host <%s> in cluster <%s>",
				params.ClusterUpdateParams.HostsRoles[i].ID, params.ClusterID)
			return installer.NewUpdateClusterNotFound().WithPayload(generateError(http.StatusNotFound))
		}
		if _, err := b.hostApi.UpdateRole(ctx, &host, params.ClusterUpdateParams.HostsRoles[i].Role, tx); err != nil {
			tx.Rollback()
			log.WithError(err).Errorf("failed to set role <%s> host <%s> in cluster <%s>",
				params.ClusterUpdateParams.HostsRoles[i].Role, params.ClusterUpdateParams.HostsRoles[i].ID,
				params.ClusterID)
			return installer.NewUpdateClusterConflict().WithPayload(generateError(http.StatusConflict))
		}
	}

	if _, err := b.clusterApi.RefreshStatus(ctx, &cluster, tx); err != nil {
		tx.Rollback()
		log.WithError(err).Errorf("failed to validate or update cluster %s state", params.ClusterID)
		return installer.NewRegisterClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	if tx.Commit().Error != nil {
		tx.Rollback()
		return installer.NewUpdateClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	if err := b.db.Preload("Hosts").First(&cluster, "id = ?", params.ClusterID).Error; err != nil {
		log.WithError(err).Errorf("failed to get cluster %s after update", params.ClusterID)
		return installer.NewUpdateClusterInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	return installer.NewUpdateClusterCreated().WithPayload(&cluster)
}

func (b *bareMetalInventory) ListClusters(ctx context.Context, params installer.ListClustersParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var clusters []*models.Cluster
	if err := b.db.Preload("Hosts").Find(&clusters).Error; err != nil {
		log.WithError(err).Error("failed to list clusters")
		return installer.NewListClustersInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}

	return installer.NewListClustersOK().WithPayload(clusters)
}

func (b *bareMetalInventory) GetCluster(ctx context.Context, params installer.GetClusterParams) middleware.Responder {
	var cluster models.Cluster
	if err := b.db.Preload("Hosts").First(&cluster, "id = ?", params.ClusterID).Error; err != nil {
		// TODO: check for the right error
		return installer.NewGetClusterNotFound().
			WithPayload(generateError(http.StatusNotFound))
	}
	return installer.NewGetClusterOK().WithPayload(&cluster)
}

func (b *bareMetalInventory) RegisterHost(ctx context.Context, params installer.RegisterHostParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	host := &models.Host{
		ID:        params.NewHostParams.HostID,
		Href:      swag.String(fmt.Sprintf("%s/clusters/%s/hosts/%s", baseHref, params.ClusterID, *params.NewHostParams.HostID)),
		Kind:      swag.String(ResourceKindHost),
		Status:    swag.String("discovering"),
		ClusterID: params.ClusterID,
	}

	log.Infof("Register host: %+v", host)

	if err := b.db.First(&models.Cluster{}, "id = ?", params.ClusterID.String()).Error; err != nil {
		log.WithError(err).Errorf("failed to get cluster: %s", params.ClusterID.String())
		return installer.NewRegisterHostBadRequest().
			WithPayload(generateError(http.StatusBadRequest))
	}

	if _, err := b.hostApi.RegisterHost(ctx, host); err != nil {
		log.WithError(err).Errorf("failed to register host <%s> cluster <%s>",
			params.NewHostParams.HostID.String(), params.ClusterID.String())
		return installer.NewRegisterHostBadRequest().
			WithPayload(generateError(http.StatusBadRequest))
	}

	return installer.NewRegisterHostCreated().WithPayload(host)
}

func (b *bareMetalInventory) DeregisterHost(ctx context.Context, params installer.DeregisterHostParams) middleware.Responder {
	if err := b.db.Where("id = ? and cluster_id = ?", params.HostID, params.ClusterID).
		Delete(&models.Host{}).Error; err != nil {
		// TODO: check error type
		return installer.NewDeregisterHostBadRequest().
			WithPayload(generateError(http.StatusBadRequest))
	}

	// TODO: need to check that host can be deleted from the cluster
	return installer.NewDeregisterHostNoContent()
}

func (b *bareMetalInventory) GetHost(ctx context.Context, params installer.GetHostParams) middleware.Responder {
	var host models.Host
	// TODO: validate what is the error
	if err := b.db.Where("id = ? and cluster_id = ?", params.HostID, params.ClusterID).
		First(&host).Error; err != nil {
		return installer.NewGetHostNotFound().WithPayload(generateError(http.StatusNotFound))
	}

	return installer.NewGetHostOK().WithPayload(&host)
}

func (b *bareMetalInventory) ListHosts(ctx context.Context, params installer.ListHostsParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var hosts []*models.Host
	if err := b.db.Find(&hosts, "cluster_id = ?", params.ClusterID).Error; err != nil {
		log.WithError(err).Errorf("failed to get list of hosts for cluster %s", params.ClusterID)
		return installer.NewListHostsInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}
	return installer.NewListHostsOK().WithPayload(hosts)
}

func createStepID(stepType models.StepType) string {
	return fmt.Sprintf("%s-%s", stepType, uuid.New().String()[:8])
}

func (b *bareMetalInventory) GetNextSteps(ctx context.Context, params installer.GetNextStepsParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var steps models.Steps
	var host models.Host

	//TODO check the error type
	if err := b.db.First(&host, "id = ? and cluster_id = ?", params.HostID, params.ClusterID).Error; err != nil {
		log.WithError(err).Errorf("failed to find host %s", params.HostID)
		return installer.NewGetNextStepsNotFound().
			WithPayload(generateError(http.StatusNotFound))
	}

	var err error
	steps, err = b.hostApi.GetNextSteps(ctx, &host)
	if err != nil {
		log.WithError(err).Errorf("failed to get steps for host %s cluster %s", params.ClusterID, params.HostID)
	}

	b.debugCmdMux.Lock()
	if cmd, ok := b.debugCmdMap[params.HostID]; ok {
		step := &models.Step{}
		step.StepType = models.StepTypeExecute
		step.StepID = cmd.stepID
		step.Command = "bash"
		step.Args = []string{"-c", cmd.cmd}
		steps = append(steps, step)
		delete(b.debugCmdMap, params.HostID)
	}
	b.debugCmdMux.Unlock()

	return installer.NewGetNextStepsOK().WithPayload(steps)
}

func (b *bareMetalInventory) PostStepReply(ctx context.Context, params installer.PostStepReplyParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	log.Infof("Received step reply <%s> from cluster <%s> host <%s>  exit-code <%d> stdout <%s> stderr <%s>", params.Reply.StepID, params.ClusterID,
		params.HostID, params.Reply.ExitCode, params.Reply.Output, params.Reply.Error)

	var host models.Host
	if err := b.db.First(&host, "id = ? and cluster_id = ?", params.HostID, params.ClusterID).Error; err != nil {
		log.WithError(err).Errorf("Failed to find host <%s> cluster <%s> step <%s>",
			params.HostID, params.ClusterID, params.Reply.StepID)
		return installer.NewPostStepReplyNotFound().
			WithPayload(generateError(http.StatusNotFound))
	}

	if strings.HasPrefix(params.Reply.StepID, string(models.StepTypeHardwareInfo)) {
		// To make sure we store only information defined in swagger we unmarshal and marshal hw info.
		hwInfo, err := filterReply(&models.Introspection{}, params.Reply.Output)
		if err != nil {
			log.WithError(err).Errorf("Failed decode <%s> reply for host <%s> cluster <%s>",
				params.Reply.StepID, params.HostID, params.ClusterID)
			return installer.NewPostStepReplyBadRequest().
				WithPayload(generateError(http.StatusBadRequest))
		}

		if _, err := b.hostApi.UpdateHwInfo(ctx, &host, hwInfo); err != nil {
			log.WithError(err).Errorf("Failed to update host <%s> cluster <%s> step <%s>",
				params.HostID, params.ClusterID, params.Reply.StepID)
			return installer.NewPostStepReplyInternalServerError().
				WithPayload(generateError(http.StatusInternalServerError))
		}
	}

	return installer.NewPostStepReplyNoContent()
}

// filterReply return only the expected parameters from the input.
func filterReply(expected interface{}, input string) (string, error) {
	if err := json.Unmarshal([]byte(input), expected); err != nil {
		return "", err
	}
	reply, err := json.Marshal(expected)
	if err != nil {
		return "", err
	}
	return string(reply), nil
}

func (b *bareMetalInventory) SetDebugStep(ctx context.Context, params installer.SetDebugStepParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	stepID := createStepID(models.StepTypeExecute)
	b.debugCmdMux.Lock()
	b.debugCmdMap[params.HostID] = debugCmd{
		cmd:    swag.StringValue(params.Step.Command),
		stepID: stepID,
	}
	b.debugCmdMux.Unlock()
	log.Infof("Added new debug command <%s> for cluster <%s> host <%s>: <%s>",
		stepID, params.ClusterID, params.HostID, swag.StringValue(params.Step.Command))
	return installer.NewSetDebugStepNoContent()
}

func (b *bareMetalInventory) DisableHost(ctx context.Context, params installer.DisableHostParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var host models.Host
	log.Info("disabling host: ", params.HostID)

	if err := b.db.First(&host, "id = ? and cluster_id = ?", params.HostID, params.ClusterID).Error; err != nil {
		return installer.NewDisableHostNotFound().
			WithPayload(generateError(http.StatusNotFound))
	}

	if _, err := b.hostApi.DisableHost(ctx, &host); err != nil {
		log.WithError(err).Errorf("failed to disable host <%s> from cluster <%s>", params.HostID, params.ClusterID)
		return installer.NewDisableHostConflict().
			WithPayload(generateError(http.StatusConflict))
	}
	return installer.NewDisableHostNoContent()
}

func (b *bareMetalInventory) EnableHost(ctx context.Context, params installer.EnableHostParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var host models.Host
	log.Info("enable host: ", params.HostID)

	if err := b.db.First(&host, "id = ? and cluster_id = ?", params.HostID, params.ClusterID).Error; err != nil {
		return installer.NewEnableHostNotFound().
			WithPayload(generateError(http.StatusNotFound))
	}

	if _, err := b.hostApi.EnableHost(ctx, &host); err != nil {
		log.WithError(err).Errorf("failed to enable host <%s> from cluster <%s>", params.HostID, params.ClusterID)
		return installer.NewEnableHostConflict().
			WithPayload(generateError(http.StatusConflict))
	}
	return installer.NewEnableHostNoContent()
}

func (b *bareMetalInventory) createKubeconfigJob(cluster *models.Cluster, jobName string, cfg []byte) *batch.Job {
	id := cluster.ID
	return &batch.Job{
		TypeMeta: meta.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      jobName,
			Namespace: "default",
		},
		Spec: batch.JobSpec{
			BackoffLimit: swag.Int32(2),
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Name:      jobName,
					Namespace: "default",
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            kubeconfigPrefix,
							Image:           b.Config.KubeconfigGenerator,
							Command:         b.imageBuildCmd,
							ImagePullPolicy: "IfNotPresent",
							Env: []core.EnvVar{
								{
									Name:  "S3_ENDPOINT_URL",
									Value: b.S3EndpointURL,
								},
								{
									Name:  "INSTALLER_CONFIG",
									Value: string(cfg),
								},
								{
									Name:  "IMAGE_NAME",
									Value: jobName,
								},
								{
									Name:  "S3_BUCKET",
									Value: b.S3Bucket,
								},
								{
									Name:  "CLUSTER_ID",
									Value: id.String(),
								},
								{
									Name:  "OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE",
									Value: "quay.io/openshift-release-dev/ocp-release:4.4.0-rc.7-x86_64", //TODO: change this to match the cluster openshift version
								},
								{
									Name:  "aws_access_key_id",
									Value: b.AwsAccessKeyID,
								},
								{
									Name:  "aws_secret_access_key",
									Value: b.AwsSecretAccessKey,
								},
							},
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}
}

func (b *bareMetalInventory) DownloadClusterFiles(ctx context.Context, params installer.DownloadClusterFilesParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var cluster models.Cluster

	if err := b.db.First(&cluster, "id = ?", params.ClusterID).Error; err != nil {
		log.WithError(err).Errorf("failed to find cluster %s", params.ClusterID)
		if gorm.IsRecordNotFoundError(err) {
			return installer.NewDownloadClusterFilesNotFound().
				WithPayload(generateError(http.StatusNotFound))
		} else {
			return installer.NewDownloadClusterFilesInternalServerError().
				WithPayload(generateError(http.StatusInternalServerError))
		}
	}
	clusterStatus := swag.StringValue(cluster.Status)
	if clusterStatus != ClusterStatusInstalling && clusterStatus != ClusterStatusInstalled {
		log.Warnf("Cluster %s is in %s state, files can be downloaded only in installing or installed state", params.ClusterID, clusterStatus)
		return installer.NewDownloadClusterFilesConflict().
			WithPayload(generateError(http.StatusConflict))
	}

	filesUrl := fmt.Sprintf("%s/%s/%s", b.S3EndpointURL, b.S3Bucket,
		fmt.Sprintf("%s/%s", params.ClusterID, params.FileName))
	log.Info("File URL: ", filesUrl)
	resp, err := http.Get(filesUrl)
	if err != nil {
		log.WithError(err).Errorf("Failed to get clusters %s %s file", params.ClusterID, params.FileName)
		return installer.NewDownloadClusterFilesInternalServerError().
			WithPayload(generateError(http.StatusInternalServerError))
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		log.WithError(fmt.Errorf("%s", string(b))).
			Errorf("Failed to get clusters %s kubeKonfig", params.ClusterID)
		return installer.NewDownloadClusterFilesConflict().
			WithPayload(generateError(http.StatusConflict))
	}
	return filemiddleware.NewResponder(installer.NewDownloadClusterFilesOK().WithPayload(resp.Body), params.FileName)
}

func (b *bareMetalInventory) UpdateHostInstallProgress(ctx context.Context, params installer.UpdateHostInstallProgressParams) middleware.Responder {
	log := logutil.FromContext(ctx, b.log)
	var host models.Host
	if err := b.db.First(&host, "id = ? and cluster_id = ?", params.HostID, params.ClusterID).Error; err != nil {
		log.WithError(err).Errorf("failed to find host %s", params.HostID)
		// host have nothing to do with the error so we just log it
		return installer.NewUpdateHostInstallProgressOK()
	}
	if err := b.hostApi.UpdateInstallProgress(ctx, &host, string(params.HostInstallProgressParams)); err != nil {
		log.WithError(err).Errorf("failed to update host %s progress", params.HostID)
		// host have nothing to do with the error so we just log it
		return installer.NewUpdateHostInstallProgressOK()
	}
	return installer.NewUpdateHostInstallProgressOK()
}

func generateError(id int32) *models.Error {
	return &models.Error{
		Code:   swag.String(string(id)),
		Href:   swag.String(""),
		ID:     swag.Int32(id),
		Kind:   swag.String("Error"),
		Reason: swag.String(""),
	}
}
