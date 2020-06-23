package cluster

import (
	"github.com/filanov/bm-inventory/internal/common"
	"github.com/filanov/stateswitch"
	"github.com/go-openapi/swag"
	"github.com/jinzhu/gorm"
)

type stateCluster struct {
	srcState string
	cluster  *common.Cluster
	db       *gorm.DB
}

func newStateCluster(c *common.Cluster, db *gorm.DB) *stateCluster {
	return &stateCluster{
		srcState: swag.StringValue(c.Status),
		cluster:  c,
		db:       db,
	}
}

func (sh *stateCluster) State() stateswitch.State {
	return stateswitch.State(swag.StringValue(sh.cluster.Status))
}

func (sh *stateCluster) SetState(state stateswitch.State) error {
	err := sh.db.Preload("Hosts").First(&sh.cluster, "id = ?", sh.cluster.ID).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	sh.cluster.Status = swag.String(string(state))
	return nil
}
