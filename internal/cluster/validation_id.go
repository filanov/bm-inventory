package cluster

import (
	"net/http"

	"github.com/filanov/bm-inventory/internal/common"
	"github.com/pkg/errors"

	"github.com/filanov/bm-inventory/models"
)

type validationID models.ClusterValidationID

const (
	IsMachineCidrDefined                       = validationID(models.ClusterValidationIDMachineCidrDefined)
	isMachineCidrEqualsToCalculatedCidr        = validationID(models.ClusterValidationIDMachineCidrEqualsToCalculatedCidr)
	isApiVipBelongToMachineCidrAndNotInUse     = validationID(models.ClusterValidationIDAPIVipBelongsToMachineCidrAndNotInUse)
	isIngressVipBelongToMachineCidrAndNotInUse = validationID(models.ClusterValidationIDIngressVipBelongsToMachineCidrAndNotInUse)
	NoPendingForInputHost                      = validationID(models.ClusterValidationIDNoPendingForInputHost)
	AllHostsAreReadyToInstall                  = validationID(models.ClusterValidationIDAllHostsAreReadyToInstall)
	HasExactlyThreeMasters                     = validationID(models.ClusterValidationIDHasExactlyThreeMasters)
)

func (v validationID) category() (string, error) {
	switch v {
	case IsMachineCidrDefined, isMachineCidrEqualsToCalculatedCidr, isApiVipBelongToMachineCidrAndNotInUse, isIngressVipBelongToMachineCidrAndNotInUse:
		return "network", nil
	case NoPendingForInputHost, AllHostsAreReadyToInstall, HasExactlyThreeMasters:
		return "hosts-data", nil
	}
	return "", common.NewApiError(http.StatusInternalServerError, errors.Errorf("Unexpected cluster validation id %s", string(v)))
}

func (v validationID) String() string {
	return string(v)
}
