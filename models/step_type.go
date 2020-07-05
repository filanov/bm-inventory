// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// StepType step type
//
// swagger:model step-type
type StepType string

const (

	// StepTypeHardwareInfo captures enum value "hardware-info"
	StepTypeHardwareInfo StepType = "hardware-info"

	// StepTypeConnectivityCheck captures enum value "connectivity-check"
	StepTypeConnectivityCheck StepType = "connectivity-check"

	// StepTypeExecute captures enum value "execute"
	StepTypeExecute StepType = "execute"

	// StepTypeInventory captures enum value "inventory"
	StepTypeInventory StepType = "inventory"

	// StepTypeInstall captures enum value "install"
	StepTypeInstall StepType = "install"

	// StepTypeFreeNetworkAddresses captures enum value "free-network-addresses"
	StepTypeFreeNetworkAddresses StepType = "free-network-addresses"

	// StepTypeResetAgent captures enum value "reset-agent"
	StepTypeResetAgent StepType = "reset-agent"

	// StepTypeResetPodman captures enum value "reset-podman"
	StepTypeResetPodman StepType = "reset-podman"
)

// for schema
var stepTypeEnum []interface{}

func init() {
	var res []StepType
	if err := json.Unmarshal([]byte(`["hardware-info","connectivity-check","execute","inventory","install","free-network-addresses","reset-agent","reset-podman"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		stepTypeEnum = append(stepTypeEnum, v)
	}
}

func (m StepType) validateStepTypeEnum(path, location string, value StepType) error {
	if err := validate.EnumCase(path, location, value, stepTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this step type
func (m StepType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateStepTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
