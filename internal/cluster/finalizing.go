package cluster

import (
	context "context"

	"github.com/filanov/bm-inventory/internal/common"
	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

func NewFinalizingState(log logrus.FieldLogger, db *gorm.DB) *finalizingState {
	return &finalizingState{
		log: log,
		db:  db,
	}
}

type finalizingState baseState

var _ StateAPI = (*Manager)(nil)

func (i *finalizingState) RefreshStatus(ctx context.Context, c *common.Cluster, db *gorm.DB) (*UpdateReply, error) {
	return &UpdateReply{
		State:     clusterStatusFinalizing,
		IsChanged: false,
	}, nil
}

//func (i *finalizingState) Install(ctx context.Context, c *common.Cluster) (*UpdateReply, error) {
//	return nil, errors.Errorf("unable to install cluster <%s> in <%s> status",
//		c.ID, swag.StringValue(c.Status))
//}
