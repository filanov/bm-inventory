package host

import (
	"context"
	"time"

	"github.com/filanov/bm-inventory/internal/hardware"
	"github.com/filanov/bm-inventory/models"
	logutil "github.com/filanov/bm-inventory/pkg/log"
	"github.com/go-openapi/swag"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func NewDisconnectedState(log logrus.FieldLogger, db *gorm.DB, hwValidator hardware.Validator) *disconnectedState {
	return &disconnectedState{
		baseState: baseState{
			log: log,
			db:  db,
		},
		hwValidator: hwValidator,
	}
}

type disconnectedState struct {
	baseState
	hwValidator hardware.Validator
}

func (d *disconnectedState) RefreshStatus(ctx context.Context, h *models.Host, db *gorm.DB) (*models.Host, error) {
	log := logutil.FromContext(ctx, d.log)
	if time.Since(time.Time(h.CheckedInAt)) < 3*time.Minute {
		return updateHostStatus(log, db, h.ClusterID, *h.ID, swag.StringValue(h.Status), HostStatusDiscovering, statusInfoDiscovering)
	}

	// Stay in the same state
	return h, nil
}
