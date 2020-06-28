package subsystem

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/alecthomas/units"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/filanov/bm-inventory/client/installer"
	"github.com/filanov/bm-inventory/models"
)

// #nosec
const (
	clusterInsufficientStateInfo = "cluster is insufficient, exactly 3 known master hosts are needed for installation"
	clusterReadyStateInfo        = "Cluster ready to be installed"
	pullSecret                   = "{\"auths\":{\"cloud.openshift.com\":{\"auth\":\"dXNlcjpwYXNzd29yZAo=\",\"email\":\"r@r.com\"}}}"
	IgnoreStateInfo              = "IgnoreStateInfo"
	installationInfo             = "Installation in progress"
	clusterErrorInfo             = "cluster %s has %d hosts in error"
)

const (
	validDiskSize = int64(128849018880)
)

var (
	validHwInfo = &models.Inventory{
		CPU:    &models.CPU{Count: 16},
		Memory: &models.Memory{PhysicalBytes: int64(32 * units.GiB)},
		Disks: []*models.Disk{
			{DriveType: "SSD", Name: "loop0", SizeBytes: validDiskSize},
			{DriveType: "HDD", Name: "sdb", SizeBytes: validDiskSize}},
		Interfaces: []*models.Interface{
			{
				IPV4Addresses: []string{
					"1.2.3.4/24",
				},
			},
		},
	}
)

func waitForClusterState(ctx context.Context, clusterID strfmt.UUID, state string, timeout time.Duration, stateInfo string) {
	for start := time.Now(); time.Since(start) < timeout; {
		rep, err := bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())
		c := rep.GetPayload()
		if swag.StringValue(c.Status) == state {
			break
		}
		time.Sleep(time.Second)
	}
	rep, err := bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
	Expect(err).NotTo(HaveOccurred())
	c := rep.GetPayload()
	Expect(swag.StringValue(c.Status)).Should(Equal(state))
	if stateInfo != IgnoreStateInfo {
		Expect(swag.StringValue(c.StatusInfo)).Should(Equal(stateInfo))
	}
}

func waitForHostState(ctx context.Context, clusterID strfmt.UUID, hostID strfmt.UUID, state string, timeout time.Duration) {
	for start := time.Now(); time.Since(start) < timeout; {
		rep, err := bmclient.Installer.GetHost(ctx, &installer.GetHostParams{ClusterID: clusterID, HostID: hostID})
		Expect(err).NotTo(HaveOccurred())
		c := rep.GetPayload()
		if swag.StringValue(c.Status) == state {
			break
		}
		time.Sleep(time.Second)
	}
	rep, err := bmclient.Installer.GetHost(ctx, &installer.GetHostParams{ClusterID: clusterID, HostID: hostID})
	Expect(err).NotTo(HaveOccurred())
	c := rep.GetPayload()
	Expect(swag.StringValue(c.Status)).Should(Equal(state))
}

func updateProgress(hostID strfmt.UUID, clusterID strfmt.UUID, progress string) {
	ctx := context.Background()
	installProgress := models.HostInstallProgressParams(progress)
	updateReply, err := bmclient.Installer.UpdateHostInstallProgress(ctx, &installer.UpdateHostInstallProgressParams{
		ClusterID:                 clusterID,
		HostInstallProgressParams: installProgress,
		HostID:                    hostID,
	})
	Expect(err).ShouldNot(HaveOccurred())
	Expect(updateReply).Should(BeAssignableToTypeOf(installer.NewUpdateHostInstallProgressOK()))
}

func installCluster(clusterID strfmt.UUID) {
	ctx := context.Background()
	_, err := bmclient.Installer.InstallCluster(ctx, &installer.InstallClusterParams{ClusterID: clusterID})
	Expect(err).NotTo(HaveOccurred())

	rep, err := bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
	Expect(err).NotTo(HaveOccurred())
	c := rep.GetPayload()
	Expect(swag.StringValue(c.Status)).Should(Equal("installing"))
	Expect(swag.StringValue(c.StatusInfo)).Should(Equal("Installation in progress"))
	Expect(len(c.Hosts)).Should(Equal(4))
	for _, host := range c.Hosts {
		Expect(swag.StringValue(host.Status)).Should(Equal("installing"))
	}

	for _, host := range c.Hosts {
		updateProgress(*host.ID, clusterID, "Done")
	}

	waitForClusterState(ctx, clusterID, "installed", 10*time.Second, "installed")
}

var _ = Describe("cluster install", func() {
	var (
		ctx           = context.Background()
		cluster       *models.Cluster
		validDiskSize = int64(128849018880)
		clusterCIDR   = "10.128.0.0/14"
		serviceCIDR   = "172.30.0.0/16"
	)

	AfterEach(func() {
		clearDB()
	})

	BeforeEach(func() {
		registerClusterReply, err := bmclient.Installer.RegisterCluster(ctx, &installer.RegisterClusterParams{
			NewClusterParams: &models.ClusterCreateParams{
				BaseDNSDomain:            "example.com",
				ClusterNetworkCidr:       &clusterCIDR,
				ClusterNetworkHostPrefix: 23,
				Name:                     swag.String("test-cluster"),
				OpenshiftVersion:         swag.String("4.5"),
				PullSecret:               pullSecret,
				ServiceNetworkCidr:       &serviceCIDR,
				SSHPublicKey:             "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC50TuHS7aYci+U+5PLe/aW/I6maBi9PBDucLje6C6gtArfjy7udWA1DCSIQd+DkHhi57/s+PmvEjzfAfzqo+L+/8/O2l2seR1pPhHDxMR/rSyo/6rZP6KIL8HwFqXHHpDUM4tLXdgwKAe1LxBevLt/yNl8kOiHJESUSl+2QSf8z4SIbo/frDD8OwOvtfKBEG4WCb8zEsEuIPNF/Vo/UxPtS9pPTecEsWKDHR67yFjjamoyLvAzMAJotYgyMoxm8PTyCgEzHk3s3S4iO956d6KVOEJVXnTVhAxrtLuubjskd7N4hVN7h2s4Z584wYLKYhrIBL0EViihOMzY4mH3YE4KZusfIx6oMcggKX9b3NHm0la7cj2zg0r6zjUn6ZCP4gXM99e5q4auc0OEfoSfQwofGi3WmxkG3tEozCB8Zz0wGbi2CzR8zlcF+BNV5I2LESlLzjPY5B4dvv5zjxsYoz94p3rUhKnnPM2zTx1kkilDK5C5fC1k9l/I/r5Qk4ebLQU= oscohen@localhost.localdomain",
			},
		})
		Expect(err).NotTo(HaveOccurred())
		cluster = registerClusterReply.GetPayload()
	})

	generateHWPostStepReply := func(h *models.Host, hwInfo *models.Inventory, hostname string) {
		hwInfo.Hostname = hostname
		hw, err := json.Marshal(&hwInfo)
		Expect(err).NotTo(HaveOccurred())
		_, err = bmclient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
			ClusterID: h.ClusterID,
			HostID:    *h.ID,
			Reply: &models.StepReply{
				ExitCode: 0,
				StepType: models.StepTypeInventory,
				Output:   string(hw),
				StepID:   string(models.StepTypeInventory),
			},
		})
		Expect(err).ShouldNot(HaveOccurred())
	}

	register3nodes := func(clusterID strfmt.UUID) []*models.Host {
		h1 := registerHost(clusterID)
		generateHWPostStepReply(h1, validHwInfo, "h1")
		h2 := registerHost(clusterID)
		generateHWPostStepReply(h2, validHwInfo, "h2")
		h3 := registerHost(clusterID)
		generateHWPostStepReply(h3, validHwInfo, "h3")

		apiVip := "1.2.3.5"
		ingressVip := "1.2.3.6"
		// All hosts are masters, one in discovering state  -> state must be insufficient
		_, err := bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{
				APIVip:     &apiVip,
				IngressVip: &ingressVip,
			},
			ClusterID: clusterID,
		})
		Expect(err).ShouldNot(HaveOccurred())
		return []*models.Host{h1, h2, h3}
	}

	Context("install cluster cases", func() {
		var clusterID strfmt.UUID
		BeforeEach(func() {
			clusterID = *cluster.ID
			registerHostsAndSetRoles(clusterID, 4)
		})

		Context("cancel installation", func() {
			It("[only_k8s]cancel running installation", func() {
				_, err := bmclient.Installer.InstallCluster(ctx, &installer.InstallClusterParams{ClusterID: clusterID})
				Expect(err).NotTo(HaveOccurred())
				waitForClusterState(ctx, clusterID, models.ClusterStatusInstalling, 10*time.Second, installationInfo)
				rep, err := bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
				Expect(err).NotTo(HaveOccurred())
				c := rep.GetPayload()
				for _, host := range c.Hosts {
					waitForHostState(ctx, clusterID, *host.ID, models.HostStatusInstalling, 10*time.Second)
				}
				_, err = bmclient.Installer.CancelInstallation(ctx, &installer.CancelInstallationParams{ClusterID: clusterID})
				Expect(err).NotTo(HaveOccurred())
				rep, err = bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
				Expect(err).NotTo(HaveOccurred())
				c = rep.GetPayload()
				Expect(swag.StringValue(c.Status)).Should(Equal(models.ClusterStatusError))
				for _, host := range c.Hosts {
					Expect(swag.StringValue(host.Status)).Should(Equal(models.HostStatusError))
				}
			})
			It("[only_k8s]cancel installation conflicts", func() {
				_, err := bmclient.Installer.CancelInstallation(ctx, &installer.CancelInstallationParams{ClusterID: clusterID})
				Expect(reflect.TypeOf(err)).Should(Equal(reflect.TypeOf(installer.NewCancelInstallationConflict())))
				rep, err := bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
				Expect(err).NotTo(HaveOccurred())
				c := rep.GetPayload()
				Expect(swag.StringValue(c.Status)).Should(Equal(models.ClusterStatusReady))
			})
			It("[only_k8s]cancel failed cluster", func() {
				FailCluster(ctx, clusterID)
				stateInfo := fmt.Sprintf(clusterErrorInfo, clusterID.String(), 1)
				waitForClusterState(ctx, clusterID, models.ClusterStatusError, 10*time.Second, stateInfo)
				_, err := bmclient.Installer.CancelInstallation(ctx, &installer.CancelInstallationParams{ClusterID: clusterID})
				Expect(err).ShouldNot(HaveOccurred())
				rep, err := bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
				Expect(err).ShouldNot(HaveOccurred())
				c := rep.GetPayload()
				Expect(swag.StringValue(c.Status)).Should(Equal(models.ClusterStatusError))
				rep, err = bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
				Expect(err).ShouldNot(HaveOccurred())
				c = rep.GetPayload()
				Expect(swag.StringValue(c.Status)).Should(Equal(models.ClusterStatusError))
				for _, host := range c.Hosts {
					Expect(swag.StringValue(host.Status)).Should(Equal(models.HostStatusError))
				}
			})
		})
	})

	It("install cluster requirement", func() {
		clusterID := *cluster.ID

		Expect(swag.StringValue(cluster.Status)).Should(Equal("insufficient"))
		Expect(swag.StringValue(cluster.StatusInfo)).Should(Equal(clusterInsufficientStateInfo))

		hosts := register3nodes(clusterID)

		h4 := registerHost(clusterID)

		// All hosts are masters, one in discovering state  -> state must be insufficient
		cluster, err := bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *hosts[0].ID, Role: "master"},
				{ID: *hosts[1].ID, Role: "master"},
				{ID: *h4.ID, Role: "master"},
			},
			},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(swag.StringValue(cluster.GetPayload().Status)).Should(Equal("insufficient"))
		Expect(swag.StringValue(cluster.GetPayload().StatusInfo)).Should(Equal(clusterInsufficientStateInfo))

		// Adding one known host and setting as master -> state must be ready
		_, err = bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *hosts[2].ID, Role: "master"},
			}},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())
		waitForClusterState(ctx, clusterID, models.ClusterStatusReady, 60*time.Second, clusterReadyStateInfo)

	})

	It("install_cluster_states", func() {
		clusterID := *cluster.ID

		Expect(swag.StringValue(cluster.Status)).Should(Equal("insufficient"))
		Expect(swag.StringValue(cluster.StatusInfo)).Should(Equal(clusterInsufficientStateInfo))

		wh1 := registerHost(clusterID)
		generateHWPostStepReply(wh1, validHwInfo, "wh1")
		wh2 := registerHost(clusterID)
		generateHWPostStepReply(wh2, validHwInfo, "wh2")
		wh3 := registerHost(clusterID)
		generateHWPostStepReply(wh3, validHwInfo, "wh3")

		mh1 := registerHost(clusterID)
		generateHWPostStepReply(mh1, validHwInfo, "mh1")
		mh2 := registerHost(clusterID)
		generateHWPostStepReply(mh2, validHwInfo, "mh2")
		mh3 := registerHost(clusterID)
		generateHWPostStepReply(mh3, validHwInfo, "mh3")

		apiVip := "1.2.3.5"
		ingressVip := "1.2.3.6"

		By("All hosts are workers -> state must be insufficient")
		cluster, err := bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *wh1.ID, Role: "worker"},
				{ID: *wh2.ID, Role: "worker"},
				{ID: *wh3.ID, Role: "worker"},
			},
				APIVip:     &apiVip,
				IngressVip: &ingressVip,
			},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(swag.StringValue(cluster.GetPayload().Status)).Should(Equal("insufficient"))
		Expect(swag.StringValue(cluster.GetPayload().StatusInfo)).Should(Equal(clusterInsufficientStateInfo))
		clusterReply, err := bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{
			ClusterID: clusterID,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(clusterReply.Payload.APIVip).To(Equal(apiVip))
		Expect(clusterReply.Payload.MachineNetworkCidr).To(Equal("1.2.3.0/24"))
		Expect(len(clusterReply.Payload.HostNetworks)).To(Equal(1))
		Expect(clusterReply.Payload.HostNetworks[0].Cidr).To(Equal("1.2.3.0/24"))
		hids := make([]interface{}, 0)
		for _, h := range clusterReply.Payload.HostNetworks[0].HostIds {
			hids = append(hids, h)
		}
		Expect(len(hids)).To(Equal(6))
		Expect(*wh1.ID).To(BeElementOf(hids...))
		Expect(*wh2.ID).To(BeElementOf(hids...))
		Expect(*wh3.ID).To(BeElementOf(hids...))
		Expect(*mh1.ID).To(BeElementOf(hids...))
		Expect(*mh2.ID).To(BeElementOf(hids...))
		Expect(*mh3.ID).To(BeElementOf(hids...))

		By("Only two masters -> state must be insufficient")
		_, err = bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *mh1.ID, Role: "master"},
				{ID: *mh2.ID, Role: "master"},
			}},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(swag.StringValue(cluster.GetPayload().Status)).Should(Equal("insufficient"))
		Expect(swag.StringValue(cluster.GetPayload().StatusInfo)).Should(Equal(clusterInsufficientStateInfo))

		By("Three master hosts -> state must be ready")
		_, err = bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *mh3.ID, Role: "master"},
			}},
			ClusterID: clusterID,
		})
		waitForHostState(ctx, clusterID, *mh3.ID, "known", 60*time.Second)

		Expect(err).NotTo(HaveOccurred())
		waitForClusterState(ctx, clusterID, models.ClusterStatusReady, 60*time.Second, clusterReadyStateInfo)

		By("Back to two master hosts -> state must be insufficient")
		cluster, err = bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *mh3.ID, Role: "worker"},
			}},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(swag.StringValue(cluster.GetPayload().Status)).Should(Equal("insufficient"))
		Expect(swag.StringValue(cluster.GetPayload().StatusInfo)).Should(Equal(clusterInsufficientStateInfo))

		By("Three master hosts -> state must be ready")
		_, err = bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *mh3.ID, Role: "master"},
			}},
			ClusterID: clusterID,
		})
		waitForHostState(ctx, clusterID, *mh3.ID, "known", 60*time.Second)

		Expect(err).NotTo(HaveOccurred())
		waitForClusterState(ctx, clusterID, "ready", 60*time.Second, clusterReadyStateInfo)

		By("Back to two master hosts -> state must be insufficient")
		cluster, err = bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *mh3.ID, Role: "worker"},
			}},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(swag.StringValue(cluster.GetPayload().Status)).Should(Equal("insufficient"))
		Expect(swag.StringValue(cluster.GetPayload().StatusInfo)).Should(Equal(clusterInsufficientStateInfo))

		_, err = bmclient.Installer.DeregisterCluster(ctx, &installer.DeregisterClusterParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())

		_, err = bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
		Expect(reflect.TypeOf(err)).To(Equal(reflect.TypeOf(installer.NewGetClusterNotFound())))
	})

	It("install_cluster_insufficient_master", func() {
		clusterID := *cluster.ID

		hwInfo := &models.Inventory{
			CPU:    &models.CPU{Count: 2},
			Memory: &models.Memory{PhysicalBytes: int64(8 * units.GiB)},
			Disks: []*models.Disk{
				{DriveType: "HDD", Name: "sdb", SizeBytes: validDiskSize},
			},
			Interfaces: []*models.Interface{
				{
					IPV4Addresses: []string{
						"1.2.3.4/24",
					},
				},
			},
		}
		h1 := registerHost(clusterID)
		generateHWPostStepReply(h1, hwInfo, "h1")
		apiVip := "1.2.3.8"
		ingressVip := "1.2.3.9"
		_, err := bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{
				APIVip:     &apiVip,
				IngressVip: &ingressVip,
			},
			ClusterID: clusterID,
		})
		Expect(err).To(Not(HaveOccurred()))
		waitForHostState(ctx, clusterID, *h1.ID, "insufficient", 60*time.Second)

		hwInfo = &models.Inventory{
			CPU:    &models.CPU{Count: 16},
			Memory: &models.Memory{PhysicalBytes: int64(32 * units.GiB)},
		}
		h2 := registerHost(clusterID)
		generateHWPostStepReply(h2, hwInfo, "h2")
		h3 := registerHost(clusterID)
		generateHWPostStepReply(h3, hwInfo, "h3")
		h4 := registerHost(clusterID)
		generateHWPostStepReply(h4, hwInfo, "h4")

		_, err = bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *h1.ID, Role: "master"},
				{ID: *h2.ID, Role: "master"},
				{ID: *h3.ID, Role: "master"},
				{ID: *h4.ID, Role: "worker"},
			}},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())
		h1 = getHost(clusterID, *h1.ID)
		waitForHostState(ctx, clusterID, *h1.ID, "insufficient", 60*time.Second)

	})
	It("unique hostname validation", func() {
		clusterID := *cluster.ID
		hosts := register3nodes(clusterID)
		_, err := bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *hosts[0].ID, Role: "master"},
			}},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())

		h1 := getHost(clusterID, *hosts[0].ID)
		waitForHostState(ctx, clusterID, *h1.ID, "known", 60*time.Second)

		By("Registering host with same hostname")
		h4 := registerHost(clusterID)
		generateHWPostStepReply(h4, validHwInfo, "h1")
		h4 = getHost(clusterID, *h4.ID)
		waitForHostState(ctx, clusterID, *h1.ID, "insufficient", 60*time.Second)
		h1 = getHost(clusterID, *h1.ID)
		Expect(*h1.Status).Should(Equal("insufficient"))

		By("Verifying install command")
		_, err = bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *h1.ID, Role: "master"},
				{ID: *hosts[1].ID, Role: "master"},
				{ID: *hosts[2].ID, Role: "master"},
				{ID: *h4.ID, Role: "worker"},
			}},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())
		_, err = bmclient.Installer.InstallCluster(ctx, &installer.InstallClusterParams{ClusterID: clusterID})
		Expect(err).Should(HaveOccurred())

		By("Registering one more host with same hostname")
		disabledHost := registerHost(clusterID)
		generateHWPostStepReply(disabledHost, validHwInfo, "h1")
		disabledHost = getHost(clusterID, *disabledHost.ID)
		waitForHostState(ctx, clusterID, *disabledHost.ID, "insufficient", 60*time.Second)
		_, err = bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *disabledHost.ID, Role: "worker"},
			}},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())

		By("Changing hostname, verify host is known now")
		generateHWPostStepReply(h4, validHwInfo, "h4")
		waitForHostState(ctx, clusterID, *h4.ID, "known", 60*time.Second)

		By("Disable host with the same hostname and verify h1 is known")
		_, err = bmclient.Installer.DisableHost(ctx, &installer.DisableHostParams{
			ClusterID: clusterID,
			HostID:    *disabledHost.ID,
		})
		Expect(err).NotTo(HaveOccurred())
		disabledHost = getHost(clusterID, *disabledHost.ID)
		Expect(*disabledHost.Status).Should(Equal("disabled"))
		waitForHostState(ctx, clusterID, *h1.ID, "known", 60*time.Second)

		By("waiting for cluster to be in ready state")
		waitForClusterState(ctx, clusterID, models.ClusterStatusReady, 60*time.Second, clusterReadyStateInfo)

		By("Verify install after disabling the host with same hostname")
		_, err = bmclient.Installer.InstallCluster(ctx, &installer.InstallClusterParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())

	})
})

func FailCluster(ctx context.Context, clusterID strfmt.UUID) strfmt.UUID {
	c, err := bmclient.Installer.InstallCluster(ctx, &installer.InstallClusterParams{ClusterID: clusterID})
	Expect(err).NotTo(HaveOccurred())
	var masterHostID strfmt.UUID
	for _, host := range c.GetPayload().Hosts {
		if host.Role == "master" {
			masterHostID = *host.ID
			break
		}
	}
	updateProgress(masterHostID, clusterID, "Failed because some error")
	masterHost := getHost(clusterID, masterHostID)
	Expect(*masterHost.Status).Should(Equal("error"))
	return masterHostID
}

var _ = Describe("cluster install, with default network params", func() {
	var (
		ctx     = context.Background()
		cluster *models.Cluster
	)

	AfterEach(func() {
		clearDB()
	})

	BeforeEach(func() {
		By("Register cluster")
		registerClusterReply, err := bmclient.Installer.RegisterCluster(ctx, &installer.RegisterClusterParams{
			NewClusterParams: &models.ClusterCreateParams{
				BaseDNSDomain:    "example.com",
				Name:             swag.String("test-cluster"),
				OpenshiftVersion: swag.String("4.5"),
				PullSecret:       pullSecret,
				SSHPublicKey:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC50TuHS7aYci+U+5PLe/aW/I6maBi9PBDucLje6C6gtArfjy7udWA1DCSIQd+DkHhi57/s+PmvEjzfAfzqo+L+/8/O2l2seR1pPhHDxMR/rSyo/6rZP6KIL8HwFqXHHpDUM4tLXdgwKAe1LxBevLt/yNl8kOiHJESUSl+2QSf8z4SIbo/frDD8OwOvtfKBEG4WCb8zEsEuIPNF/Vo/UxPtS9pPTecEsWKDHR67yFjjamoyLvAzMAJotYgyMoxm8PTyCgEzHk3s3S4iO956d6KVOEJVXnTVhAxrtLuubjskd7N4hVN7h2s4Z584wYLKYhrIBL0EViihOMzY4mH3YE4KZusfIx6oMcggKX9b3NHm0la7cj2zg0r6zjUn6ZCP4gXM99e5q4auc0OEfoSfQwofGi3WmxkG3tEozCB8Zz0wGbi2CzR8zlcF+BNV5I2LESlLzjPY5B4dvv5zjxsYoz94p3rUhKnnPM2zTx1kkilDK5C5fC1k9l/I/r5Qk4ebLQU= oscohen@localhost.localdomain",
			},
		})
		Expect(err).NotTo(HaveOccurred())
		cluster = registerClusterReply.GetPayload()
	})

	It("[only_k8s]install cluster", func() {
		clusterID := *cluster.ID
		registerHostsAndSetRoles(clusterID, 3)
		rep, err := bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())
		c := rep.GetPayload()
		startTimeInstalling := c.InstallStartedAt
		startTimeInstalled := c.InstallCompletedAt

		_, err = bmclient.Installer.InstallCluster(ctx, &installer.InstallClusterParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())

		rep, err = bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())
		c = rep.GetPayload()
		Expect(swag.StringValue(c.Status)).Should(Equal("installing"))
		Expect(swag.StringValue(c.StatusInfo)).Should(Equal("Installation in progress"))
		Expect(len(c.Hosts)).Should(Equal(3))
		Expect(c.InstallStartedAt).ShouldNot(Equal(startTimeInstalling))
		for _, host := range c.Hosts {
			Expect(swag.StringValue(host.Status)).Should(Equal("installing"))
		}
		// fake installation completed
		for _, host := range c.Hosts {
			updateProgress(*host.ID, clusterID, "Done")
		}

		waitForClusterState(ctx, clusterID, "installed", 10*time.Second, "installed")
		rep, err = bmclient.Installer.GetCluster(ctx, &installer.GetClusterParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())
		c = rep.GetPayload()
		Expect(c.InstallCompletedAt).ShouldNot(Equal(startTimeInstalled))
	})
})

func registerHostsAndSetRoles(clusterID strfmt.UUID, numHosts int) {
	ctx := context.Background()

	generateHWPostStepReply := func(h *models.Host, hwInfo *models.Inventory, hostname string) {
		hwInfo.Hostname = hostname
		hw, err := json.Marshal(&hwInfo)
		Expect(err).NotTo(HaveOccurred())
		_, err = bmclient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
			ClusterID: h.ClusterID,
			HostID:    *h.ID,
			Reply: &models.StepReply{
				ExitCode: 0,
				Output:   string(hw),
				StepID:   string(models.StepTypeInventory),
				StepType: models.StepTypeInventory,
			},
		})
		Expect(err).ShouldNot(HaveOccurred())
	}
	for i := 0; i < numHosts; i++ {
		hostname := fmt.Sprintf("h%d", i)
		host := registerHost(clusterID)
		generateHWPostStepReply(host, validHwInfo, hostname)
		var role string
		if i < 3 {
			role = "master"
		} else {
			role = "worker"
		}
		_, err := bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterUpdateParams: &models.ClusterUpdateParams{HostsRoles: []*models.ClusterUpdateParamsHostsRolesItems0{
				{ID: *host.ID, Role: role},
			}},
			ClusterID: clusterID,
		})
		Expect(err).NotTo(HaveOccurred())

	}
	apiVip := ""
	ingressVip := ""
	_, err := bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
		ClusterUpdateParams: &models.ClusterUpdateParams{
			APIVip:     &apiVip,
			IngressVip: &ingressVip,
		},
		ClusterID: clusterID,
	})
	Expect(err).NotTo(HaveOccurred())
	apiVip = "1.2.3.8"
	ingressVip = "1.2.3.9"
	c, err := bmclient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
		ClusterUpdateParams: &models.ClusterUpdateParams{
			APIVip:     &apiVip,
			IngressVip: &ingressVip,
		},
		ClusterID: clusterID,
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(swag.StringValue(c.GetPayload().Status)).Should(Equal("ready"))
	Expect(swag.StringValue(c.GetPayload().StatusInfo)).Should(Equal(clusterReadyStateInfo))
}
