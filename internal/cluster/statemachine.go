package cluster

import (
	"github.com/filanov/stateswitch"
)

const (
	TransitionTypeCancelInstallation = "CancelInstallation"
)

func NewClusterStateMachine(th *transitionHandler) stateswitch.StateMachine {
	sm := stateswitch.NewStateMachine()

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType: TransitionTypeCancelInstallation,
		SourceStates: []stateswitch.State{
			clusterStatusInstalling,
		},
		DestinationState: clusterStatusError,
		PostTransition:   th.PostCancelInstallation,
	})

	return sm
}
