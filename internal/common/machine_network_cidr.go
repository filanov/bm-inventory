package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/filanov/bm-inventory/models"
	"github.com/sirupsen/logrus"
)

/*
 * Calculate the machine network CIDR from the one of (ApiVip, IngressVip) and the ip addresses of the hosts.
 * The ip addresses of the host appear with CIDR notation. Therefore, the network can be calculated from it.
 * The goal of this function is to find the first network that one of the vips belongs to it.
 * This network is returned as a result.
 */
func CalculateMachineNetworkCIDR(cluster *models.Cluster) (string, error) {
	var ip string
	if cluster.APIVip != "" {
		ip = cluster.APIVip
	} else if cluster.IngressVip != "" {
		ip = cluster.IngressVip
	} else {
		return "", nil
	}
	parsedVipAddr := net.ParseIP(ip)
	if parsedVipAddr == nil {
		return "", fmt.Errorf("Could not parse VIP ip %s", ip)
	}
	for _, h := range cluster.Hosts {
		var inventory models.Inventory
		err := json.Unmarshal([]byte(h.Inventory), &inventory)
		if err != nil {
			continue
		}
		for _, intf := range inventory.Interfaces {
			for _, ipv4addr := range intf.IPV4Addresses {
				_, ipnet, err := net.ParseCIDR(ipv4addr)
				if err != nil {
					continue
				}
				if ipnet.Contains(parsedVipAddr) {
					return ipnet.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("No suitable matching CIDR found for VIP %s", ip)
}

func ipInCidr(ipStr, cidrStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false
	}
	return ipnet.Contains(ip)
}

func verifyVip(cluster *models.Cluster, vip string, vipName string, mustExist bool) error {
	if !mustExist && vip == "" || ipInCidr(vip, cluster.MachineNetworkCidr) {
		return nil
	}
	return fmt.Errorf("%s <%s> does not belong to machine-network-cidr <%s>", vipName, vip, cluster.MachineNetworkCidr)
}

func verifyDifferentVipAddresses(cluster *models.Cluster) error {
	if cluster.APIVip == cluster.IngressVip && cluster.APIVip != "" {
		return fmt.Errorf("api-vip and ingress-vip cannot have the same value: %s", cluster.APIVip)
	}
	return nil
}

func VerifyVips(cluster *models.Cluster, mustExist bool) error {
	err := verifyVip(cluster, cluster.APIVip, "api-vip", mustExist)
	if err == nil {
		err = verifyVip(cluster, cluster.IngressVip, "ingress-vip", mustExist)
	}
	if err == nil {
		err = verifyDifferentVipAddresses(cluster)
	}
	return err
}

func belongsToNetwork(log logrus.FieldLogger, h *models.Host, machineIpnet *net.IPNet) bool {
	var inventory models.Inventory
	err := json.Unmarshal([]byte(h.Inventory), &inventory)
	if err != nil {
		log.WithError(err).Warnf("Error unmarshalling host %s inventory %s", h.ID, h.Inventory)
		return false
	}
	for _, intf := range inventory.Interfaces {
		for _, ipv4addr := range intf.IPV4Addresses {
			ip, _, err := net.ParseCIDR(ipv4addr)
			if err != nil {
				log.WithError(err).Warnf("Could not parse cidr %s", ipv4addr)
				continue
			}
			if machineIpnet.Contains(ip) {
				return true
			}
		}
	}
	return false
}

func GetMachineCIDRHosts(log logrus.FieldLogger, cluster *models.Cluster) ([]*models.Host, error) {
	if cluster.MachineNetworkCidr == "" {
		return nil, errors.New("Machine network CIDR was not set in cluster")
	}
	_, machineIpnet, err := net.ParseCIDR(cluster.MachineNetworkCidr)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.Host, 0)
	for _, h := range cluster.Hosts {
		if belongsToNetwork(log, h, machineIpnet) {
			ret = append(ret, h)
		}
	}
	return ret, nil
}

func IsHostInMachineNetCidr(log logrus.FieldLogger, cluster *models.Cluster, host *models.Host) bool {
	_, machineIpnet, err := net.ParseCIDR(cluster.MachineNetworkCidr)
	if err != nil {
		return false
	}
	return belongsToNetwork(log, host, machineIpnet)
}
