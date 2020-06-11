package cluster

import (
	context "context"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/go-openapi/swag"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/filanov/bm-inventory/models"
	installer2 "github.com/filanov/bm-inventory/restapi/operations/installer"
)

func NewRegistrar(log logrus.FieldLogger, db *gorm.DB) *registrar {
	return &registrar{
		log: log,
		db:  db,
	}
}

type registrar struct {
	log logrus.FieldLogger
	db  *gorm.DB
}

func (r *registrar) RegisterCluster(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	id := strfmt.UUID(uuid.New().String())

	cluster.ID = &id
	url := installer2.GetClusterURL{ClusterID: id}
	cluster.Href = swag.String(url.String())

	cluster.Status = swag.String(clusterStatusInsufficient)
	cluster.StatusInfo = swag.String(statusInfoInsufficient)
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			r.log.Error("update cluster failed")
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		r.log.WithError(tx.Error).Error("failed to start transaction")
	}

	if err := tx.Preload("Hosts").Create(cluster).Error; err != nil {
		r.log.Errorf("Error registering cluster %s", cluster.Name)
		tx.Rollback()
		return cluster, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return cluster, err
	}

	return cluster, nil
}

func (r *registrar) DeregisterCluster(ctx context.Context, cluster *models.Cluster) error {
	var txErr error
	tx := r.db.Begin()

	defer func() {
		if txErr != nil {
			tx.Rollback()
		}
	}()

	if swag.StringValue(cluster.Status) == clusterStatusInstalling {
		tx.Rollback()
		return errors.Errorf("cluster %s can not be removed while being installed", cluster.ID)
	}

	if txErr = tx.Where("cluster_id = ?", cluster.ID).Delete(&models.Host{}).Error; txErr != nil {
		tx.Rollback()
		return errors.Errorf("failed to deregister host while unregistering cluster %s", cluster.ID)
	}

	if txErr = tx.Delete(cluster).Error; txErr != nil {
		tx.Rollback()
		return errors.Errorf("failed to delete cluster %s", cluster.ID)
	}

	if tx.Commit().Error != nil {
		tx.Rollback()
		return errors.Errorf("failed to delete cluster %s, commit tx", cluster.ID)
	}
	return nil
}
