package cluster

import (
	"fmt"
	"net"

	"github.com/go-openapi/swag"

	"github.com/go-openapi/strfmt"

	"github.com/filanov/bm-inventory/internal/host"

	"github.com/filanov/bm-inventory/internal/network"

	"github.com/sirupsen/logrus"

	"github.com/filanov/bm-inventory/internal/common"
	"github.com/jinzhu/gorm"

	"github.com/filanov/bm-inventory/models"
)

type validationStatus string

const (
	ValidationSuccess validationStatus = "success"
	ValidationFailure validationStatus = "failure"
	ValidationPending validationStatus = "pending"
)

const (
	ApiVipName     = "api vip"
	IngressVipName = "ingress vip"
)

func (v validationStatus) String() string {
	return string(v)
}

type clusterValidationContext struct {
	clusterId     strfmt.UUID
	cluster       *common.Cluster
	db            *gorm.DB
	calculateCidr string
}

type validationConditon func(context *clusterValidationContext) validationStatus
type validationStringFormatter func(context *clusterValidationContext, status validationStatus) string

type validation struct {
	id        validationID
	condition validationConditon
	formatter validationStringFormatter
}

func (c *clusterValidationContext) loadCluster() error {
	var cluster common.Cluster
	err := c.db.Preload("Hosts", "status <> ?", host.HostStatusDisabled).Take(&cluster, "id = ?", c.clusterId.String()).Error
	if err == nil {
		c.cluster = &cluster
	}
	return err
}

func newClusterValidationContext(clusterId strfmt.UUID, db *gorm.DB) (*clusterValidationContext, error) {
	ret := &clusterValidationContext{
		clusterId: clusterId,
		db:        db,
	}
	err := ret.loadCluster()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func boolValue(b bool) validationStatus {
	if b {
		return ValidationSuccess
	} else {
		return ValidationFailure
	}
}

type clusterValidator struct {
	log logrus.FieldLogger
}

func (v *clusterValidator) isMachineCidrDefined(c *clusterValidationContext) validationStatus {
	if c.cluster.MachineNetworkCidr == "" {
		return ValidationPending
	}
	_, _, err := net.ParseCIDR(c.cluster.MachineNetworkCidr)
	return boolValue(err == nil)
}

func (v *clusterValidator) printIsMachineCidrDefined(context *clusterValidationContext, status validationStatus) string {
	switch status {
	case ValidationPending:
		return "Machine network CIDR is undefined"
	case ValidationSuccess:
		return "Machine network CIDR is defined"
	case ValidationFailure:
		return fmt.Sprintf("Machine network CIDR has invalid value %s", context.cluster.MachineNetworkCidr)
	default:
		return fmt.Sprintf("Unexpected status %s", status)
	}
}

func (v *clusterValidator) isMachineCidrEqualsToCalculatedCidr(c *clusterValidationContext) validationStatus {
	cidr, err := network.CalculateMachineNetworkCIDR(c.cluster.APIVip, c.cluster.IngressVip, c.cluster.Hosts)
	c.calculateCidr = cidr
	return boolValue(err == nil && cidr == c.cluster.MachineNetworkCidr)
}

func (v *clusterValidator) printIsMachineCidrEqualsToCalculatedCidr(context *clusterValidationContext, status validationStatus) string {
	switch status {
	case ValidationSuccess:
		return "Cluster machine CIDR equals to the calculated CIDR "
	case ValidationFailure:
		return fmt.Sprintf("Cluster machine CIDR %s is different than the calculated CIDR %s", context.cluster.MachineNetworkCidr, context.calculateCidr)
	default:
		return fmt.Sprintf("Unexpected status %s", status)
	}
}

func (v *clusterValidator) isApiVipBelongToMachineCidrAndNotInUse(c *clusterValidationContext) validationStatus {
	err := network.VerifyVip(c.cluster.Hosts, c.cluster.MachineNetworkCidr, c.cluster.APIVip, ApiVipName,
		true, v.log)
	return boolValue(err == nil)
}

func (v *clusterValidator) printIsApiVipBelongToMachineCidrAndNotInUse(context *clusterValidationContext, status validationStatus) string {
	switch status {
	case ValidationSuccess:
		return fmt.Sprintf("%s %s belongs to machine CIDR and not in use ", ApiVipName, context.cluster.APIVip)
	case ValidationFailure:
		return fmt.Sprintf("%s %s does not belong to machine CIDR or already in use ", ApiVipName, context.cluster.APIVip)
	default:
		return fmt.Sprintf("Unexpected status %s", status)
	}
}

func (v *clusterValidator) isIngressVipBelongToMachineCidrAndNotInUse(c *clusterValidationContext) validationStatus {
	err := network.VerifyVip(c.cluster.Hosts, c.cluster.MachineNetworkCidr, c.cluster.IngressVip, IngressVipName,
		true, v.log)
	return boolValue(err == nil)
}

func (v *clusterValidator) printIsIngressApiVipBelongToMachineCidrAndNotInUse(context *clusterValidationContext, status validationStatus) string {
	switch status {
	case ValidationSuccess:
		return fmt.Sprintf("%s %s belongs to machine CIDR and not in use ", IngressVipName, context.cluster.IngressVip)
	case ValidationFailure:
		return fmt.Sprintf("%s %s does not belong to machine CIDR or already in use ", IngressVipName, context.cluster.IngressVip)
	default:
		return fmt.Sprintf("Unexpected status %s", status)
	}
}

func (v *clusterValidator) hasExactlyThreeKnownMasters(c *clusterValidationContext) validationStatus {
	mappedMastersByRole := MapMasterHostsByStatus(c.cluster)

	mastersInKnown, ok := mappedMastersByRole[models.HostStatusKnown]
	return boolValue(ok && len(mastersInKnown) == MinHostsNeededForInstallation)
}

func (v *clusterValidator) printHasExactlyThreeKnownMasters(context *clusterValidationContext, status validationStatus) string {
	switch status {
	case ValidationSuccess:
		return fmt.Sprintf("Cluster has %d known master hosts", MinHostsNeededForInstallation)
	case ValidationFailure:
		return fmt.Sprintf("cluster is expected to have exactly %d known master to be installed", MinHostsNeededForInstallation)
	default:
		return fmt.Sprintf("Unexpected status %s", status)
	}
}

func (v *clusterValidator) noHostInPendingForInput(c *clusterValidationContext) validationStatus {
	mappedHostsByStatus := MapHostsByStatus(c.cluster)

	hostsInPendingForInput, ok := mappedHostsByStatus[models.HostStatusPendingForInput]
	return boolValue(!ok || (ok && len(hostsInPendingForInput) == 0))
}

func (v *clusterValidator) printNoHostInPendingForInput(context *clusterValidationContext, status validationStatus) string {
	switch status {
	case ValidationSuccess:
		return "Cluster does not have pending for input hosts"
	case ValidationFailure:
		return "Cluster has pending for input hosts"
	default:
		return fmt.Sprintf("Unexpected status %s", status)
	}
}

func (v *clusterValidator) allHostsAreReadyToInstall(c *clusterValidationContext) validationStatus {
	foundNotKnownHost := false
	for _, host := range c.cluster.Hosts {
		if swag.StringValue(host.Status) != models.HostStatusDisabled && swag.StringValue(host.Status) != models.HostStatusKnown {
			foundNotKnownHost = true
			break
		}
	}
	return boolValue(!foundNotKnownHost)
}

func (v *clusterValidator) printAllHostsAreReadyToInstall(context *clusterValidationContext, status validationStatus) string {
	switch status {
	case ValidationSuccess:
		return "All hosts in the cluster are ready to install"
	case ValidationFailure:
		return "Cluster has hosts that are not ready to install"
	default:
		return fmt.Sprintf("Unexpected status %s", status)
	}
}
