package job

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/filanov/bm-inventory/internal/common"
	"github.com/filanov/bm-inventory/internal/events"
	"github.com/filanov/bm-inventory/pkg/generator"
	logutil "github.com/filanov/bm-inventory/pkg/log"
	"github.com/sirupsen/logrus"
)

type LocalJob interface {
	Execute(pythonCommand string, pythonFilePath string, envVars []string, log logrus.FieldLogger) error
	generator.ISOInstallConfigGenerator
}

type localJob struct {
	Config
	log logrus.FieldLogger
}

func NewLocalJob(log logrus.FieldLogger, cfg Config) *localJob {
	return &localJob{
		Config: cfg,
		log:    log,
	}
}

func (j *localJob) Execute(pythonCommand string, pythonFilePath string, envVars []string, log logrus.FieldLogger) error {
	cmd := exec.Command(pythonCommand, pythonFilePath)
	cmd.Env = envVars
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Infoln(pythonFilePath + " executed with error: " + err.Error())
		log.Infoln("envVars: " + strings.Join(envVars, ","))
		return err
	}
	log.Infoln(cmd.Stdout)
	return nil
}

// creates install config
func (j *localJob) GenerateInstallConfig(ctx context.Context, cluster common.Cluster, cfg []byte) error {
	log := logutil.FromContext(ctx, j.log)
	envVars := append(os.Environ(),
		"S3_ENDPOINT_URL="+j.Config.S3EndpointURL,
		"INSTALLER_CONFIG="+string(cfg),
		"INVENTORY_ENDPOINT="+"http://"+strings.TrimSpace(j.InventoryURL)+":"+strings.TrimSpace(j.InventoryPort)+"/api/assisted-install/v1",
		"IMAGE_NAME="+j.Config.KubeconfigGenerator,
		"S3_BUCKET="+j.Config.S3Bucket,
		"CLUSTER_ID="+cluster.ID.String(),
		"OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE="+j.Config.ReleaseImage,
		"aws_access_key_id="+j.Config.AwsAccessKeyID,
		"aws_secret_access_key="+j.Config.AwsSecretAccessKey,
		"WORK_DIR=/data",
	)
	return j.Execute("python", "./data/process-ignition-manifests-and-kubeconfig.py", envVars, log)
}

func (j *localJob) GenerateISO(ctx context.Context, cluster common.Cluster, jobName string, imageName string, ignitionConfig string, eventsHandler events.Handler) error {
	log := logutil.FromContext(ctx, j.log)
	envVars := append(os.Environ(),
		"S3_ENDPOINT_URL="+j.Config.S3EndpointURL,
		"IGNITION_CONFIG="+ignitionConfig,
		"IMAGE_NAME="+imageName,
		"S3_BUCKET="+j.Config.S3Bucket,
		"aws_access_key_id="+j.Config.AwsAccessKeyID,
		"aws_secret_access_key="+j.Config.AwsSecretAccessKey,
		"WORK_DIR=/data",
	)
	return j.Execute("python", "./data/install_process.py", envVars, log)
}
