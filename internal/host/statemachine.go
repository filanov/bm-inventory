package host

import "github.com/filanov/stateswitch"

const (
	TransitionTypeRegisterHost = "RegisterHost"
)

func NewHostStateMachine(th *transitionHandler) stateswitch.StateMachine {
	sm := stateswitch.NewStateMachine()

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType: TransitionTypeRegisterHost,
		SourceStates: []stateswitch.State{
			"",
			HostStatusDiscovering,
			HostStatusKnown,
			HostStatusDisconnected,
			HostStatusInsufficient,
		},
		DestinationState: HostStatusDiscovering,
		PostTransition:   th.PostRegisterHost,
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeRegisterHost,
		SourceStates:     []stateswitch.State{HostStatusInstalling, HostStatusInstallingInProgress},
		DestinationState: HostStatusError,
		PostTransition:   th.PostRegisterDuringInstallation,
	})

	return sm
}
