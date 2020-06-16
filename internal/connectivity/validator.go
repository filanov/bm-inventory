package connectivity

import (
	"encoding/json"
	"fmt"

	"github.com/filanov/bm-inventory/internal/common"
	"github.com/filanov/bm-inventory/models"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=validator.go -package=connectivity -destination=mock_connectivity_validator.go
type Validator interface {
	IsSufficient(host *models.Host, cluster *common.Cluster) (*common.IsSufficientReply, error)
	GetHostValidInterfaces(host *models.Host) ([]*models.Interface, error)
}

func NewValidator(log logrus.FieldLogger) Validator {
	return &validator{
		log: log,
	}
}

type validator struct {
	log logrus.FieldLogger
}

func (v *validator) IsSufficient(host *models.Host, cluster *common.Cluster) (*common.IsSufficientReply, error) {
	var reason string
	isSufficient := true

	_, err := v.GetHostValidInterfaces(host)
	if err != nil {
		isSufficient = false
		reason = "Waiting to receive connectivity information"
	}

	if cluster.MachineNetworkCidr == "" {
		isSufficient = false
		reason += ", Could not determine connectivity because API VIP not set"
	}

	if !common.IsHostInMachineNetCidr(v.log, cluster, host) {
		isSufficient = false
		reason += fmt.Sprintf(", host %s does not belong to cluster machine network %s, The machine network is set by configuring the API-VIP", *host.ID, cluster.MachineNetworkCidr)
	}

	return &common.IsSufficientReply{
		Type:         "connectivity",
		IsSufficient: isSufficient,
		Reason:       reason,
	}, nil
}

func (v *validator) GetHostValidInterfaces(host *models.Host) ([]*models.Interface, error) {
	var inventory models.Inventory
	if err := json.Unmarshal([]byte(host.Inventory), &inventory); err != nil {
		return nil, err
	}
	if len(inventory.Interfaces) == 0 {
		return nil, fmt.Errorf("host %s doesn't have interfaces", host.ID)
	}
	return inventory.Interfaces, nil
}
