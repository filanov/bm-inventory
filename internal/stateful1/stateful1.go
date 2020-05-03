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
	Known = stateful.DefaultState("Known")
	Insufficient      = stateful.DefaultState("Insufficient")
	Disconnected = stateful.DefaultState("Disconnected")
	Disabled = stateful.DefaultState("Disabled")
	Installing = stateful.DefaultState("Installing")
	Installed = stateful.DefaultState("Installed")
	Error  = stateful.DefaultState("Error")
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
	RegisterHost(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error)
	// Set a new HW information
	UpdateHwInfo(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error)
	// Set host state
	UpdateRole(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error)
	RefreshStatus(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error)
	Install(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error)
	EnableHost(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error)
	DisableHost(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error)
}

// Register a new host
func (b *BaseState) RegisterHost(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Set a new HW information
func (b *BaseState)  UpdateHwInfo(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Set host state
func (b *BaseState) UpdateRole(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// check keep alive
func (b *BaseState) RefreshStatus(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Install host - db is optional, for transactions
func (b *BaseState) Install(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Enable host to get requests (disabled by default)
func (b *BaseState) EnableHost(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}
// Disable host from getting any requests
func (b *BaseState) DisableHost(_ *HostStatemachine, _ stateful.TransitionArguments) (stateful.State, error){
	return b, nil
}

var (
	Initializing *InitializingState = &InitializingState{BaseState{ID:"Initializing"}}
	Discovering *DiscoveringState = &DiscoveringState{BaseState{ID:"Discovering"}}
)

type InitializingState struct {
	BaseState
}

func (i *InitializingState)  RegisterHost(h *HostStatemachine, params stateful.TransitionArguments) (stateful.State, error){
	// Implement logic here
	return i, nil
}

type DiscoveringState struct {
	BaseState
}

func (d *DiscoveringState)  RegisterHost(h *HostStatemachine, params stateful.TransitionArguments) (stateful.State, error){
	// Implement logic here
	return d, nil
}

type HostStatemachine struct {
	host models.Host
	state StateInterface
	db  *gorm.DB
	ctx context.Context
	log logrus.FieldLogger
}

func (h *HostStatemachine) Reload(id, clusterId strfmt.UUID) error {
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


func (h *HostStatemachine) RegisterHost(params stateful.TransitionArguments) (stateful.State, error){
	hostParams, ok := params.(*models.Host)
	if !ok {
		return nil, fmt.Errorf("Expected host as argument")
	}
	if err := h.Reload(*hostParams.ID, hostParams.ClusterID); err != nil {
		return nil, err
	}
	return h.state.RegisterHost(h, params)
}

