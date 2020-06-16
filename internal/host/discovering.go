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

func NewDiscoveringState(log logrus.FieldLogger, db *gorm.DB, hwValidator hardware.Validator) *discoveringState {
	return &discoveringState{
		baseState: baseState{
			log: log,
			db:  db,
		},
		hwValidator: hwValidator,
	}
}

type discoveringState struct {
	baseState
	hwValidator hardware.Validator
}

func (d *discoveringState) UpdateHwInfo(ctx context.Context, h *models.Host, hwInfo string) (*UpdateReply, error) {
	h.HardwareInfo = hwInfo
	return updateHwInfo(logutil.FromContext(ctx, d.log), d.hwValidator, h, d.db)
}

func (d *discoveringState) UpdateInventory(ctx context.Context, h *models.Host, inventory string) (*UpdateReply, error) {
	h.Inventory = inventory
	return updateInventory(logutil.FromContext(ctx, d.log), d.hwValidator, h, d.db)
}

func (d *discoveringState) UpdateRole(ctx context.Context, h *models.Host, role string, db *gorm.DB) (*UpdateReply, error) {
	h.Role = role
	cdb := d.db
	if db != nil {
		cdb = db
	}
	return updateRole(logutil.FromContext(ctx, d.log), h, cdb)
}

func (d *discoveringState) RefreshStatus(ctx context.Context, h *models.Host) (*UpdateReply, error) {
	//checking if need to change state to disconnect
	stateReply, err := updateByKeepAlive(logutil.FromContext(ctx, d.log), h, d.db)
	if err != nil || stateReply.IsChanged {
		return stateReply, err
	}
	var statusInfoDetails = make(map[string]string)
	//checking inventory isInsufficient
	inventoryReply, _ := d.hwValidator.IsSufficient(h)
	if inventoryReply != nil {
		statusInfoDetails[inventoryReply.Type] = inventoryReply.Reason
	} else {
		statusInfoDetails["hardware"] = "parsing error"
	}

	//TODO: checking connectivity isInsufficient

	//checking role
	roleReply := isSufficientRole(h)
	statusInfoDetails[roleReply.Type] = roleReply.Reason

	d.log.Infof("refresh status host: %s role reply %+v inventory reply %+v", h.ID, roleReply, inventoryReply)

	if inventoryReply != nil && inventoryReply.IsSufficient && roleReply.IsSufficient {
		return updateState(d.log, HostStatusKnown, "", h, d.db)
	} else {
		statusInfo, err := json.Marshal(statusInfoDetails)
		if err != nil {
			return nil, err
		}
		return updateState(d.log, HostStatusInsufficient, string(statusInfo), h, d.db)
	}

}

func (d *discoveringState) Install(ctx context.Context, h *models.Host, db *gorm.DB) (*UpdateReply, error) {
	return nil, errors.Errorf("unable to install host <%s> in <%s> status",
		h.ID, swag.StringValue(h.Status))
}

func (d *discoveringState) EnableHost(ctx context.Context, h *models.Host) (*UpdateReply, error) {
	// State in the same state
	return &UpdateReply{
		State:     HostStatusDiscovering,
		IsChanged: false,
	}, nil
}

func (d *discoveringState) DisableHost(ctx context.Context, h *models.Host) (*UpdateReply, error) {
	return updateState(logutil.FromContext(ctx, d.log), HostStatusDisabled, statusInfoDisabled, h, d.db)
}
