package host

import (
	"context"
	"encoding/json"

	"github.com/filanov/bm-inventory/internal/hardware"
	"github.com/filanov/bm-inventory/models"
	logutil "github.com/filanov/bm-inventory/pkg/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewKnownState(log logrus.FieldLogger, db *gorm.DB, hwValidator hardware.Validator) *knownState {
	return &knownState{
		baseState: baseState{
			log: log,
			db:  db,
		},
		hwValidator: hwValidator,
	}
}

type knownState struct {
	baseState
	hwValidator hardware.Validator
}

func (k *knownState) UpdateHwInfo(ctx context.Context, h *models.Host, hwInfo string) (*UpdateReply, error) {
	h.HardwareInfo = hwInfo
	return updateHwInfo(logutil.FromContext(ctx, k.log), k.hwValidator, h, k.db)
}

func (k *knownState) UpdateInventory(ctx context.Context, h *models.Host, inventory string) (*UpdateReply, error) {
	h.Inventory = inventory
	return updateInventory(logutil.FromContext(ctx, k.log), k.hwValidator, h, k.db)
}

func (k *knownState) UpdateRole(ctx context.Context, h *models.Host, role string, db *gorm.DB) (*UpdateReply, error) {
	h.Role = role
	cdb := k.db
	if db != nil {
		cdb = db
	}
	return updateRole(logutil.FromContext(ctx, k.log), h, cdb)
}

func (k *knownState) RefreshStatus(ctx context.Context, h *models.Host) (*UpdateReply, error) {
	//checking if need to change state to disconnect
	stateReply, err := updateByKeepAlive(logutil.FromContext(ctx, k.log), h, k.db)
	if err != nil || stateReply.IsChanged {
		return stateReply, err
	}
	var statusInfoDetails = make(map[string]string)
	//checking inventory isInsufficient
	inventoryReply, _ := k.hwValidator.IsSufficient(h)
	if inventoryReply != nil {
		statusInfoDetails[inventoryReply.Type] = inventoryReply.Reason
	} else {
		statusInfoDetails["hardware"] = "parsing error"
	}
	//TODO: checking connectivity isInsufficient

	//checking role
	roleReply := isSufficientRole(h)
	statusInfoDetails[roleReply.Type] = roleReply.Reason

	if inventoryReply != nil && inventoryReply.IsSufficient && roleReply.IsSufficient {
		return updateState(k.log, HostStatusKnown, "", h, k.db)
	} else {
		k.log.Infof("refresh status host: %s role reply %+v inventory reply %+v", h.ID, roleReply, inventoryReply)
		statusInfo, err := json.Marshal(statusInfoDetails)
		if err != nil {
			return nil, err
		}
		return updateState(k.log, HostStatusInsufficient, string(statusInfo), h, k.db)
	}
}

func (k *knownState) Install(ctx context.Context, h *models.Host, db *gorm.DB) (*UpdateReply, error) {
	if h.Role == "" {
		return nil, errors.Errorf("unable to install host <%s> without a role", h.ID)
	}
	cdb := k.db
	if db != nil {
		cdb = db
	}
	return updateState(logutil.FromContext(ctx, k.log), HostStatusInstalling, statusInfoInstalling, h, cdb)
}

func (k *knownState) EnableHost(ctx context.Context, h *models.Host) (*UpdateReply, error) {
	// State in the same state
	return &UpdateReply{
		State:     HostStatusKnown,
		IsChanged: false,
	}, nil
}

func (k *knownState) DisableHost(ctx context.Context, h *models.Host) (*UpdateReply, error) {
	return updateState(logutil.FromContext(ctx, k.log), HostStatusDisabled, statusInfoDisabled, h, k.db)
}
