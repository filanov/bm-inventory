package host

import (
	"github.com/filanov/bm-inventory/models"
	"github.com/filanov/stateswitch"
	"github.com/go-openapi/swag"
	"github.com/jinzhu/gorm"
)

type stateHost struct {
	srcState string
	host     *models.Host
	db       *gorm.DB
}

func newStateHost(h *models.Host, db *gorm.DB) *stateHost {
	return &stateHost{
		srcState: swag.StringValue(h.Status),
		host:     h,
		db:       db,
	}
}

func (sh *stateHost) State() stateswitch.State {
	return stateswitch.State(swag.StringValue(sh.host.Status))
}

func (sh *stateHost) SetState(state stateswitch.State) error {
	err := sh.db.First(&sh.host, "id = ? and cluster_id = ?", sh.host.ID, sh.host.ClusterID).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	sh.host.Status = swag.String(string(state))
	return nil
}
