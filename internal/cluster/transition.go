package cluster

import (
	"context"

	logutil "github.com/filanov/bm-inventory/pkg/log"

	"github.com/sirupsen/logrus"

	"github.com/filanov/stateswitch"
	"github.com/pkg/errors"
)

type transitionHandler struct {
	log logrus.FieldLogger
}

////////////////////////////////////////////////////////////////////////////
// CancelInstallation
////////////////////////////////////////////////////////////////////////////

type TransitionArgsCancelInstallation struct {
	ctx    context.Context
	reason string
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
	return updateClusterStateWithParams(logutil.FromContext(params.ctx, th.log), sCluster.srcState,
		params.reason, sCluster.cluster, sCluster.db)
}
