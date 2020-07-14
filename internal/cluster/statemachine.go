package cluster

import (
	"github.com/filanov/stateswitch"
)

const (
	TransitionTypeCancelInstallation   = "CancelInstallation"
	TransitionTypeResetCluster         = "ResetCluster"
	TransitionTypeCompleteInstallation = "CompleteInstallation"
)

func NewClusterStateMachine(th *transitionHandler) stateswitch.StateMachine {
	sm := stateswitch.NewStateMachine()

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType: TransitionTypeCancelInstallation,
		SourceStates: []stateswitch.State{
			clusterStatusInstalling,
			clusterStatusError,
		},
		DestinationState: clusterStatusError,
		PostTransition:   th.PostCancelInstallation,
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType: TransitionTypeResetCluster,
		SourceStates: []stateswitch.State{
			clusterStatusError,
		},
		DestinationState: clusterStatusInsufficient,
		PostTransition:   th.PostResetCluster,
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType: TransitionTypeCompleteInstallation,
		Condition:      th.isSuccess,
		Transition: func(stateSwitch stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
			params, _ := args.(*TransitionArgsCompleteInstallation)
			params.reason = statusInfoInstalled
			return nil
		},
		SourceStates: []stateswitch.State{
			clusterStatusFinalizing,
		},
		DestinationState: clusterStatusInstalled,
		PostTransition:   th.PostCompleteInstallation,
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType: TransitionTypeCompleteInstallation,
		Condition:      th.notSuccess,
		SourceStates: []stateswitch.State{
			clusterStatusFinalizing,
		},
		DestinationState: clusterStatusError,
		PostTransition:   th.PostCompleteInstallation,
	})

	return sm
}
