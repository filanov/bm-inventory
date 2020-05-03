package statemachine

import (
	"context"
	"fmt"
	"github.com/bykof/stateful"
	"github.com/filanov/bm-inventory/models"
	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const (
	Initializing = stateful.DefaultState("Initializing")
	Discovering = stateful.DefaultState("Discovering")
	Known = stateful.DefaultState("Known")
	Insufficient      = stateful.DefaultState("Insufficient")
	Disconnected = stateful.DefaultState("Disconnected")
	Disabled = stateful.DefaultState("Disabled")
	Installing = stateful.DefaultState("Installing")
	Installed = stateful.DefaultState("Installed")
	Error  = stateful.DefaultState("Error")
)


type StatefulHost struct {
	host models.Host
	state stateful.State
	db  *gorm.DB
	ctx context.Context
	log logrus.FieldLogger
}

func (h *StatefulHost) Reload(id, clusterId strfmt.UUID) error {
	if err := h.db.First(&h.host, "id = ? and cluster_id = ?", id, clusterId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.state = Initializing
			return nil
		}
		return err
	}
	switch h.host.StatusInfo {
	case "discovering":
		h.state = Discovering
	case "initializing":
		h.state = Initializing
		// etc.
	}
	return nil
}

func (h *StatefulHost) State() stateful.State {
	return h.state
}

func (h *StatefulHost) SetState(state stateful.State) error {
	// Save in database?
	h.state = state
	return nil
}


func (h *StatefulHost) RegisterHost(params stateful.TransitionArguments) (stateful.State, error){
	hostParams, ok := params.(*models.Host)
	if !ok {
		return nil, fmt.Errorf("Expected host as argument")
	}
	if err := h.Reload(*hostParams.ID, hostParams.ClusterID); err != nil {
		return nil, err
	}
	nextState := h.state
	switch h.state {
	case Initializing:
		// One logic
		nextState = Discovering
	case Discovering:
		//
	}
	return nextState, nil
}


func newStateMachine(h *models.Host) *stateful.StateMachine {
	sh := StatefulHost{}
	sh.Reload(*h.ID, h.ClusterID)
	ret := stateful.StateMachine{StatefulObject:&sh}

	ret.AddTransition(
		sh.RegisterHost,
		stateful.States{Initializing},
		stateful.States{Discovering},
		)
	return &ret
}

func RegisterHost(ctx context.Context, h *models.Host) error {
	sm := newStateMachine(h)
	sh := sm.StatefulObject.(*StatefulHost)
	return sm.Run(sh.RegisterHost, h)
}