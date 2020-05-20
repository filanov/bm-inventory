package cluster

import (
	"context"

	"github.com/sirupsen/logrus"

	host2 "github.com/filanov/bm-inventory/internal/host"
	"github.com/filanov/bm-inventory/models"
	logutil "github.com/filanov/bm-inventory/pkg/log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func NewInsufficientState(log logrus.FieldLogger, db *gorm.DB) *insufficientState {
	return &insufficientState{
		log: log,
		db:  db,
	}
}

type insufficientState baseState

func (i *insufficientState) RefreshStatus(ctx context.Context, c *models.Cluster, db *gorm.DB) (*UpdateReply, error) {

	log := logutil.FromContext(ctx, i.log)

	if err := db.Preload("Hosts").First(&c, "id = ?", c.ID).Error; err != nil {
		return &UpdateReply{
			State:     clusterStatusInsufficient,
			IsChanged: false}, errors.Errorf("cluster %s not found", c.ID)
	}
	mappedMastersByRole := mapMasterHostsByStatus(c)

	// Cluster is ready
	mastersInKnown, ok := mappedMastersByRole[host2.HostStatusKnown]
	if ok && len(mastersInKnown) >= minHostsNeededForInstallation {
		log.Infof("Cluster %s has at least %d known master hosts, cluster is ready.", c.ID, minHostsNeededForInstallation)
		return updateState(clusterStatusReady, c, db, log)

		//cluster is still insufficient
	} else {
		return &UpdateReply{State: clusterStatusInsufficient,
			IsChanged: false}, nil
	}
}
