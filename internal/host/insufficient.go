package host

import (
	"context"
	"encoding/json"

	"github.com/filanov/bm-inventory/internal/hardware"
	"github.com/filanov/bm-inventory/models"
	logutil "github.com/filanov/bm-inventory/pkg/log"
	"github.com/go-openapi/swag"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewInsufficientState(log logrus.FieldLogger, db *gorm.DB, hwValidator hardware.Validator) *insufficientState {
	return &insufficientState{
		baseState: baseState{
			log: log,
			db:  db,
		},
		hwValidator: hwValidator,
	}
}

type insufficientState struct {
	baseState
	hwValidator hardware.Validator
}

func (i *insufficientState) UpdateHwInfo(ctx context.Context, h *models.Host, hwInfo string) (*UpdateReply, error) {
	h.HardwareInfo = hwInfo
	return updateHwInfo(logutil.FromContext(ctx, i.log), i.hwValidator, h, i.db)
}

func (d *insufficientState) UpdateInventory(ctx context.Context, h *models.Host, inventory string) (*UpdateReply, error) {
	h.Inventory = inventory
	return updateInventory(logutil.FromContext(ctx, d.log), d.hwValidator, h, d.db)
}

func (i *insufficientState) UpdateRole(ctx context.Context, h *models.Host, role string, db *gorm.DB) (*UpdateReply, error) {
	h.Role = role
	cdb := i.db
	if db != nil {
		cdb = db
	}
	return updateRole(logutil.FromContext(ctx, i.log), h, cdb)
}

func (i *insufficientState) RefreshStatus(ctx context.Context, h *models.Host) (*UpdateReply, error) {
	//checking if need to change state to disconnect
	stateReply, err := updateByKeepAlive(logutil.FromContext(ctx, i.log), h, i.db)
	if err != nil || stateReply.IsChanged {
		return stateReply, err
	}
	var statusInfoDetails = make(map[string]string)
	//checking inventory isInsufficient
	inventoryReply, _ := i.hwValidator.IsSufficient(h)
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
		i.log.Infof("refresh status host: %s role reply %+v inventory reply %+v", h.ID, roleReply, inventoryReply)
		return updateState(i.log, HostStatusKnown, "", h, i.db)
	} else {
		statusInfo, err := json.Marshal(statusInfoDetails)
		if err != nil {
			return nil, err
		}
		return updateState(i.log, HostStatusInsufficient, string(statusInfo), h, i.db)
	}
}

func (i *insufficientState) Install(ctx context.Context, h *models.Host, db *gorm.DB) (*UpdateReply, error) {
	return nil, errors.Errorf("unable to install host <%s> in <%s> status",
		h.ID, swag.StringValue(h.Status))
}

func (i *insufficientState) EnableHost(ctx context.Context, h *models.Host) (*UpdateReply, error) {
	// State in the same state
	return &UpdateReply{
		State:     HostStatusInsufficient,
		IsChanged: false,
	}, nil
}

func (i *insufficientState) DisableHost(ctx context.Context, h *models.Host) (*UpdateReply, error) {
	return updateState(logutil.FromContext(ctx, i.log), HostStatusDisabled, statusInfoDisabled, h, i.db)
}
