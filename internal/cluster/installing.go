package cluster

import (
	context "context"

	host2 "github.com/filanov/bm-inventory/internal/host"

	logutil "github.com/filanov/bm-inventory/pkg/log"

	"github.com/sirupsen/logrus"

	"github.com/filanov/bm-inventory/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func NewInstallingState(log logrus.FieldLogger, db *gorm.DB) *installingState {
	return &installingState{
		log: log,
		db:  db,
	}
}

type installingState baseState

var _ StateAPI = (*Manager)(nil)

func (i *installingState) RefreshStatus(ctx context.Context, c *models.Cluster, db *gorm.DB) (*UpdateReply, error) {
	log := logutil.FromContext(ctx, i.log)
	installationState, err := i.getClusterInstallationState(ctx, c)
	if err != nil {
		return nil, errors.Errorf("couldn't determine cluster %s installation state", c.ID)
	}

	switch installationState {
	case clusterStatusInstalled:
		return updateState(clusterStatusInstalled, c, i.db, log)
	case clusterStatusError:
		return updateState(clusterStatusError, c, i.db, log)
	case clusterStatusInstalling:
		return &UpdateReply{
			State:     clusterStatusInstalling,
			IsChanged: false,
		}, nil
	}
	return nil, errors.Errorf("cluster % state transaction is not clear, installation state: %s ", c.ID, installationState)
}

func (i *installingState) getClusterInstallationState(ctx context.Context, c *models.Cluster) (string, error) {
	log := logutil.FromContext(ctx, i.log)

	if err := i.db.Preload("Hosts").First(&c, "id = ?", c.ID).Error; err != nil {
		return "", errors.Errorf("cluster %s not found", c.ID)
	}

	mappedMastersByRole := mapMasterHostsByStatus(c)

	// Cluster is in installed
	mastersInInstalled, ok := mappedMastersByRole[host2.HostStatusInstalled]
	if ok && len(mastersInInstalled) >= minHostsNeededForInstallation {
		log.Infof("Cluster %s has at least %d installed hosts, cluster is installed.", c.ID, len(mastersInInstalled))
		return host2.HostStatusInstalled, nil
	}

	// Cluster is installing
	mastersInInstalling, ok := mappedMastersByRole[host2.HostStatusInstalling]
	if ok && len(mastersInInstalling) > 0 &&
		(len(mastersInInstalling)+len(mastersInInstalled)) >= minHostsNeededForInstallation {
		return host2.HostStatusInstalling, nil
	}

	// Cluster is in error
	mastersInError := mappedMastersByRole[host2.HostStatusError]
	log.Warningf("Cluster %s has %d hosts in error.", c.ID, len(mastersInError))
	return host2.HostStatusError, nil
}
