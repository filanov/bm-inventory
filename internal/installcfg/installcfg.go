package installcfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/filanov/bm-inventory/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type bmc struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type host struct {
	Name            string `yaml:"name"`
	Role            string `yaml:"role"`
	Bmc             bmc    `yaml:"bmc"`
	BootMACAddress  string `yaml:"bootMACAddress"`
	BootMode        string `yaml:"bootMode"`
	HardwareProfile string `yaml:"hardwareProfile"`
}

type baremetal struct {
	ProvisioningNetworkInterface string `yaml:"provisioningNetworkInterface"`
	APIVIP                       string `yaml:"apiVIP"`
	IngressVIP                   string `yaml:"ingressVIP"`
	DNSVIP                       string `yaml:"dnsVIP"`
	Hosts                        []host `yaml:"hosts"`
}

type platform struct {
	Baremetal baremetal `yaml:"baremetal"`
}

type InstallerConfigBaremetal struct {
	APIVersion string `yaml:"apiVersion"`
	BaseDomain string `yaml:"baseDomain"`
	Networking struct {
		NetworkType    string `yaml:"networkType"`
		ClusterNetwork []struct {
			Cidr       string `yaml:"cidr"`
			HostPrefix int    `yaml:"hostPrefix"`
		} `yaml:"clusterNetwork"`
		MachineNetwork []struct {
			Cidr string `yaml:"cidr"`
		} `yaml:"machineNetwork"`
		ServiceNetwork []string `yaml:"serviceNetwork"`
	} `yaml:"networking"`
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Compute []struct {
		Name     string `yaml:"name"`
		Replicas int    `yaml:"replicas"`
	} `yaml:"compute"`
	ControlPlane struct {
		Name     string `yaml:"name"`
		Replicas int    `yaml:"replicas"`
	} `yaml:"controlPlane"`
	Platform   platform `yaml:"platform"`
	PullSecret string   `yaml:"pullSecret"`
	SSHKey     string   `yaml:"sshKey"`
}

func countHostsByRole(cluster *models.Cluster, role string) int {
	var count int
	for _, host := range cluster.Hosts {
		if host.Role == role {
			count += 1
		}
	}
	return count
}

func getMachineCIDR(cluster *models.Cluster) (string, error) {
	parsedVipAddr := net.ParseIP(string(cluster.APIVip))
	if parsedVipAddr == nil {
		errStr := fmt.Sprintf("Could not parse VIP ip %s", cluster.APIVip)
		logrus.Warn(errStr)
		return "", errors.New(errStr)
	}
	for _, h := range cluster.Hosts {
		var inventory models.Inventory
		err := json.Unmarshal([]byte(h.Inventory), &inventory)
		if err != nil {
			logrus.WithError(err).Warnf("Error unmarshalling host inventory %s", h.Inventory)
			continue
		}
		for _, intf := range inventory.Interfaces {
			for _, ipv4addr := range intf.IPV4Addresses {
				_, ipnet, err := net.ParseCIDR(ipv4addr)
				if err != nil {
					logrus.WithError(err).Warnf("Could not parse cidr %s", ipv4addr)
					continue
				}
				if ipnet.Contains(parsedVipAddr) {
					return ipnet.String(), nil
				}
			}
		}
	}
	errStr := fmt.Sprintf("No suitable matching CIDR found for VIP %s", cluster.APIVip)
	logrus.Warn(errStr)
	return "", errors.New(errStr)
}

func getBasicInstallConfig(cluster *models.Cluster, machineCIDR string) *InstallerConfigBaremetal {
	return &InstallerConfigBaremetal{
		APIVersion: "v1",
		BaseDomain: cluster.BaseDNSDomain,
		Networking: struct {
			NetworkType    string `yaml:"networkType"`
			ClusterNetwork []struct {
				Cidr       string `yaml:"cidr"`
				HostPrefix int    `yaml:"hostPrefix"`
			} `yaml:"clusterNetwork"`
			MachineNetwork []struct {
				Cidr string `yaml:"cidr"`
			} `yaml:"machineNetwork"`
			ServiceNetwork []string `yaml:"serviceNetwork"`
		}{
			NetworkType: "OpenShiftSDN",
			ClusterNetwork: []struct {
				Cidr       string `yaml:"cidr"`
				HostPrefix int    `yaml:"hostPrefix"`
			}{
				{Cidr: cluster.ClusterNetworkCidr, HostPrefix: int(cluster.ClusterNetworkHostPrefix)},
			},
			MachineNetwork: []struct {
				Cidr string `yaml:"cidr"`
			}{
				{Cidr: machineCIDR},
			},
			ServiceNetwork: []string{cluster.ServiceNetworkCidr},
		},
		Metadata: struct {
			Name string `yaml:"name"`
		}{
			Name: cluster.Name,
		},
		Compute: []struct {
			Name     string `yaml:"name"`
			Replicas int    `yaml:"replicas"`
		}{
			{Name: "worker", Replicas: countHostsByRole(cluster, "worker")},
		},
		ControlPlane: struct {
			Name     string `yaml:"name"`
			Replicas int    `yaml:"replicas"`
		}{
			Name:     "master",
			Replicas: countHostsByRole(cluster, "master"),
		},
		PullSecret: cluster.PullSecret,
		SSHKey:     cluster.SSHPublicKey,
	}
}

// [TODO] - remove once we decide to use specific values from the hosts of the cluster
func getDummyMAC(dummyMAC string, count int) (string, error) {
	hwMac, err := net.ParseMAC(dummyMAC)
	if err != nil {
		logrus.Warn("Failed to parse dummyMac")
		return "", err
	}
	hwMac[len(hwMac)-1] = hwMac[len(hwMac)-1] + byte(count)
	return hwMac.String(), nil
}

func setPlatformInstallconfig(cluster *models.Cluster, cfg *InstallerConfigBaremetal) error {
	// set hosts
	numMasters := countHostsByRole(cluster, "master")
	numWorkers := countHostsByRole(cluster, "worker")
	masterCount := 0
	workerCount := 0
	hosts := make([]host, numWorkers+numMasters)

	// dummy MAC and port, once we start using real BMH, those values should be set from cluster
	dummyMAC := "00:aa:39:b3:51:10"
	dummyPort := 6230

	for i := range hosts {
		logrus.Infof("Setting master, host %d, master count %d", i, masterCount)
		if i >= numMasters {
			hosts[i].Name = fmt.Sprintf("openshift-worker-%d", workerCount)
			hosts[i].Role = "worker"
			workerCount += 1
		} else {
			hosts[i].Name = fmt.Sprintf("openshift-master-%d", masterCount)
			hosts[i].Role = "master"
			masterCount += 1
		}
		hosts[i].Bmc = bmc{
			Address:  fmt.Sprintf("ipmi://192.168.111.1:%d", dummyPort+i),
			Username: "admin",
			Password: "rackattack",
		}
		hwMac, err := getDummyMAC(dummyMAC, i)
		if err != nil {
			logrus.Warn("Failed to parse dummyMac")
			return err
		}
		hosts[i].BootMACAddress = hwMac
		hosts[i].BootMode = "UEFI"
		hosts[i].HardwareProfile = "unknown"
	}
	cfg.Platform = platform{
		Baremetal: baremetal{
			ProvisioningNetworkInterface: "ethh0",
			APIVIP:                       cluster.APIVip.String(),
			IngressVIP:                   cluster.IngressVip.String(),
			DNSVIP:                       cluster.APIVip.String(),
			Hosts:                        hosts,
		},
	}
	return nil
}

func GetInstallConfig(cluster *models.Cluster) ([]byte, error) {
	machineCidr, err := getMachineCIDR(cluster)
	if err != nil {
		return nil, err
	}
	cfg := getBasicInstallConfig(cluster, machineCidr)
	err = setPlatformInstallconfig(cluster, cfg)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(*cfg)
}
