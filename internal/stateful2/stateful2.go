package stateful1

import (
	"context"
	"fmt"
	"github.com/bykof/stateful"
	"github.com/filanov/bm-inventory/models"
	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type BaseState struct {
	// Register a new host
	ID string
}

func (b *BaseState) GetID() string {
	return b.ID
}

func (b *BaseState)  IsWildCard() bool {
	return false
}

type StateInterface interface {
	stateful.State
	RegisterHost(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error)
	// Set a new HW information
	UpdateHwInfo(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error)
	// Set host state
	UpdateRole(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error)
	RefreshStatus(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error)
	Install(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error)
	EnableHost(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error)
	DisableHost(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error)
}

// Register a new host
func (b *BaseState) RegisterHost(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Set a new HW information
func (b *BaseState)  UpdateHwInfo(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Set host state
func (b *BaseState) UpdateRole(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// check keep alive
func (b *BaseState) RefreshStatus(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Install host - db is optional, for transactions
func (b *BaseState) Install(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Enable host to get requests (disabled by default)
func (b *BaseState) EnableHost(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Disable host from getting any requests
func (b *BaseState) DisableHost(_ *StatefulHost, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}

var (
	Initializing *InitializingState = &InitializingState{BaseState{ID:"Initializing"}}
	Discovering *DiscoveringState = &DiscoveringState{BaseState{ID:"Discovering"}}
)

type InitializingState struct {
	BaseState
}

func (i *InitializingState)  RegisterHost(h *StatefulHost, params stateful.TransitionArguments) (stateful.State, error){
	// Implement logic here
	return i, nil
}

type DiscoveringState struct {
	BaseState
}

func (d *DiscoveringState)  RegisterHost(h *StatefulHost, params stateful.TransitionArguments) (stateful.State, error){
	// Implement logic here
	return d, nil
}

type StatefulHost struct {
	host models.Host
	state StateInterface
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
	h.state = state.(StateInterface)
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
	return h.state.RegisterHost(h, params)
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

