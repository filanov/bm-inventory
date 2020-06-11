package host

import (
	"context"

	"github.com/go-openapi/swag"

	"github.com/filanov/bm-inventory/models"
	logutil "github.com/filanov/bm-inventory/pkg/log"

	"github.com/sirupsen/logrus"

	"github.com/filanov/stateswitch"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type transitionHandler struct {
	db                  *gorm.DB
	log                 logrus.FieldLogger
	externalValidations ExternalValidations
}

////////////////////////////////////////////////////////////////////////////
// RegisterHost
////////////////////////////////////////////////////////////////////////////

type TransitionArgsRegisterHost struct {
	ctx context.Context
}

func (th *transitionHandler) CanHostRegister(sw stateswitch.StateSwitch, _ stateswitch.TransitionArgs) (bool, error) {
	sHost, ok := sw.(*stateHost)
	if !ok {
		return false, errors.New("PostRegisterHost incompatible type of StateSwitch")
	}
	cluster := models.Cluster{}
	if err := th.db.First(&cluster, "id = ?", sHost.host.ClusterID).Error; err != nil {
		return false, errors.Wrapf(err, "failed to get cluster <%s>", sHost.host.ClusterID)
	}

	if err := th.externalValidations.AcceptRegistration(&cluster); err != nil {
		return false, err
	}
	return true, nil
}

func (th *transitionHandler) PostRegisterHost(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
	sHost, ok := sw.(*stateHost)
	if !ok {
		return errors.New("PostRegisterHost incompatible type of StateSwitch")
	}
	params, ok := args.(*TransitionArgsRegisterHost)
	if !ok {
		return errors.New("PostRegisterHost invalid argument")
	}

	host := models.Host{}
	log := logutil.FromContext(params.ctx, th.log)

	// if already exists, reset role and hw info
	if err := th.db.First(&host, "id = ? and cluster_id = ?", sHost.host.ID, sHost.host.ClusterID).Error; err == nil {
		currentState := swag.StringValue(host.Status)
		host.Status = sHost.host.Status
		return updateHostStateWithParams(log, currentState, statusInfoDiscovering, &host, th.db,
			"hardware_info", "", "role", "")
	}

	log.Infof("Register new host %s cluster %s", sHost.host.ID.String(), sHost.host.ClusterID)
	return th.db.Create(sHost.host).Error
}

func (th *transitionHandler) PostRegisterDuringInstallation(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
	sHost, ok := sw.(*stateHost)
	if !ok {
		return errors.New("RegisterNewHost incompatible type of StateSwitch")
	}
	params, ok := args.(*TransitionArgsRegisterHost)
	if !ok {
		return errors.New("PostRegisterDuringInstallation invalid argument")
	}
	return updateHostStateWithParams(logutil.FromContext(params.ctx, th.log), sHost.srcState,
		"Tried to register during installation", sHost.host, th.db)
}

////////////////////////////////////////////////////////////////////////////
// Installation failure
////////////////////////////////////////////////////////////////////////////

type TransitionArgsHostInstallationFailed struct {
	ctx context.Context
}

func (th *transitionHandler) PostHostInstallationFailed(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
	sHost, ok := sw.(*stateHost)
	if !ok {
		return errors.New("HostInstallationFailed incompatible type of StateSwitch")
	}
	params, ok := args.(*TransitionArgsHostInstallationFailed)
	if !ok {
		return errors.New("HostInstallationFailed invalid argument")
	}
	return updateHostStateWithParams(logutil.FromContext(params.ctx, th.log), sHost.srcState,
		"installation command failed", sHost.host, th.db)
}
