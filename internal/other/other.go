package other

import (
	"context"
	"encoding/json"
	"github.com/filanov/bm-inventory/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type State string
type States []State

type StatefulHost struct {
	host models.Host
	state State
	db  *gorm.DB
	ctx context.Context
	log logrus.FieldLogger
}



const (
	Initializing = State("Initializing")
	Discovering = State("Discovering")
	Known = State("Known")
	Insufficient = State("Insufficient")
)

type TransitionType string

const (
	RegisterHostType = TransitionType("RegisterHost")
	HardwareInfoType = TransitionType("HardwareInfo")
)

type ValidatorFunction func(h *StatefulHost,args ...interface{}) bool
type ActionFunction func(h *StatefulHost, args ...interface{}) error

type Transition struct {
	Sources States
	Destination State
	TransitionType TransitionType
	Validator ValidatorFunction
	Action ActionFunction
}

func registerHost1(h *StatefulHost, args ...interface{}) error {
	// RegisterHost logic 1
	return nil
}

func registerHost2(h *StatefulHost, args ...interface{}) error {
	// RegisterHost logic 2
	return nil
}

func sufficient(h *StatefulHost, args ...interface{}) bool {
	// Is hardware sufficient
	return true
}

func notSufficient(h *StatefulHost, args ...interface{}) bool {
	return !sufficient(h, args...)
}


func getTransitions() []Transition {
	ret := []Transition {
		{
			Sources:States{Initializing, Discovering},
			Destination: Discovering,
			TransitionType:RegisterHostType,
			Action: registerHost1,
		},
		{
			Sources:States{Known},
			Destination: Discovering,
			TransitionType:RegisterHostType,
			Action: registerHost2,
		},
		{
			Sources:States{Discovering},
			Destination:Insufficient,
			TransitionType:HardwareInfoType,
			Validator: notSufficient,
		},
		{
			Sources:States{Discovering},
			Destination:Known,
			TransitionType:HardwareInfoType,
			Validator: sufficient,
		},
	}
	return ret
}

type StateMachine struct {
	// state machine data here
}

func newStateMachine() *StateMachine {
	ret := StateMachine{}
	for _, t := range getTransitions() {
		// Build state machine here
	}
	return &ret
}

var (
	myStateMachine = newStateMachine()
)

func reload(h *models.Host) *StatefulHost {
	ret := StatefulHost{}
	// load from db
	return &ret
}

func (sm *StateMachine)Run(ctx context.Context, transitionType TransitionType, h *models.Host, args ...interface{}) error {
	sh := reload(h)
	// Perform transitions
	return nil
}

func RegisterHost(ctx context.Context, h *models.Host) error {
	return myStateMachine.Run(ctx, RegisterHostType, h)
}


func HardwareInfo(ctx context.Context, h *models.Host, hardwareInfo string) error {
	var hwInfo models.Introspection
	json.Unmarshal([]byte(hardwareInfo), &hwInfo)
	return myStateMachine.Run(ctx, HardwareInfoType, h, &hwInfo)
}