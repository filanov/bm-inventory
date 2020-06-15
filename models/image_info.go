// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ImageInfo image info
//
// swagger:model image_info
type ImageInfo struct {

	// created at
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"created_at,omitempty" gorm:"type:datetime"`

	// Image generator version
	GeneratorVersion string `json:"generator_version,omitempty"`

	// The URL of the HTTP/S proxy that agents should use to access the discovery service
	// http://\<user\>:\<password\>@\<server\>:\<port\>/
	//
	ProxyURL string `json:"proxy_url,omitempty"`

	// SSH public key for debugging the installation
	SSHPublicKey string `json:"ssh_public_key,omitempty" gorm:"type:varchar(1024)"`
}

// Validate validates this image info
func (m *ImageInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ImageInfo) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("created_at", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ImageInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ImageInfo) UnmarshalBinary(b []byte) error {
	var res ImageInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
