package cluster

import (
	"context"

	"github.com/jinzhu/gorm"

	logutil "github.com/filanov/bm-inventory/pkg/log"

	"github.com/sirupsen/logrus"

	"github.com/filanov/stateswitch"
	"github.com/pkg/errors"
)

type transitionHandler struct {
	log logrus.FieldLogger
	db  *gorm.DB
}

////////////////////////////////////////////////////////////////////////////
// CancelInstallation
////////////////////////////////////////////////////////////////////////////

type TransitionArgsCancelInstallation struct {
	ctx    context.Context
	reason string
	db     *gorm.DB
}

func (th *transitionHandler) PostCancelInstallation(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
	sCluster, ok := sw.(*stateCluster)
	if !ok {
		return errors.New("PostCancelInstallation incompatible type of StateSwitch")
	}
	params, ok := args.(*TransitionArgsCancelInstallation)
	if !ok {
		return errors.New("PostCancelInstallation invalid argument")
	}
	if sCluster.srcState == clusterStatusError {
		return nil
	}
	return updateClusterStateWithParams(logutil.FromContext(params.ctx, th.log), sCluster.srcState,
		params.reason, sCluster.cluster, params.db)
}

////////////////////////////////////////////////////////////////////////////
// ResetCluster
////////////////////////////////////////////////////////////////////////////

type TransitionArgsResetCluster struct {
	ctx    context.Context
	reason string
	db     *gorm.DB
}

func (th *transitionHandler) PostResetCluster(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
	sCluster, ok := sw.(*stateCluster)
	if !ok {
		return errors.New("PostResetCluster incompatible type of StateSwitch")
	}
	params, ok := args.(*TransitionArgsResetCluster)
	if !ok {
		return errors.New("PostResetCluster invalid argument")
	}
	return updateClusterStateWithParams(logutil.FromContext(params.ctx, th.log), sCluster.srcState,
		params.reason, sCluster.cluster, params.db)
}

////////////////////////////////////////////////////////////////////////////
// Complete installation
////////////////////////////////////////////////////////////////////////////

type TransitionArgsCompleteInstallation struct {
	ctx       context.Context
	isSuccess bool
	reason    string
	db        *gorm.DB
}

func (th *transitionHandler) PostCompleteInstallation(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
	sCluster, ok := sw.(*stateCluster)
	if !ok {
		return errors.New("PostCompleteInstallation incompatible type of StateSwitch")
	}
	params, ok := args.(*TransitionArgsCompleteInstallation)
	if !ok {
		return errors.New("PostCompleteInstallation invalid argument")
	}
	sCluster.State()

	return updateClusterStateWithParams(logutil.FromContext(params.ctx, th.log), sCluster.srcState,
		params.reason, sCluster.cluster, params.db)
}

func (th *transitionHandler) isSuccess(stateSwitch stateswitch.StateSwitch, args stateswitch.TransitionArgs) (b bool, err error) {
	params, _ := args.(*TransitionArgsCompleteInstallation)
	return params.isSuccess, nil
}

func (th *transitionHandler) notSuccess(stateSwitch stateswitch.StateSwitch, args stateswitch.TransitionArgs) (b bool, err error) {
	params, _ := args.(*TransitionArgsCompleteInstallation)
	return !params.isSuccess, nil
}
