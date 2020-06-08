package cluster

import (
	context "context"

	logutil "github.com/filanov/bm-inventory/pkg/log"

	"github.com/go-openapi/strfmt"

	"github.com/pkg/errors"

	"github.com/go-openapi/swag"

	"github.com/filanov/bm-inventory/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func NewInstaller(log logrus.FieldLogger, db *gorm.DB) *installer {
	return &installer{
		log: log,
		db:  db,
	}
}

type installer struct {
	log logrus.FieldLogger
	db  *gorm.DB
}

func (i *installer) Install(ctx context.Context, c *models.Cluster, db *gorm.DB) error {
	log := logutil.FromContext(ctx, i.log)

	switch swag.StringValue(c.Status) {
	case "":
	case clusterStatusReady:
		log.Infof("cluster %s is starting installation", c.ID)
	case clusterStatusInsufficient:
		return errors.Errorf("cluster %s is missing the resources to be installed", c.ID)
	case clusterStatusInstalling:
		return errors.Errorf("cluster %s is already installing", c.ID)
	case clusterStatusInstalled:
		return errors.Errorf("cluster %s is already installed", c.ID)
	case clusterStatusError:
		return errors.Errorf("cluster %s has a error", c.ID)
	default:
		return errors.Errorf("cluster %s state is unclear - cluster state: %s", c.ID, swag.StringValue(c.Status))
	}

	_, err := updateState(clusterStatusInstalling, statusInfoInstalling, c, db, i.log)
	if err != nil {
		return err
	}

	return nil
}

func (i *installer) GetMasterNodesIds(ctx context.Context, cluster *models.Cluster, db *gorm.DB) ([]*strfmt.UUID, error) {
	return getKnownMastersNodesIds(cluster, db)
}
