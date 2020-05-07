// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Cluster cluster
//
// swagger:model cluster
type Cluster struct {

	// Virtual IP used to reach the OpenShift cluster API.
	// Format: hostname
	APIVip strfmt.Hostname `json:"api_vip,omitempty"`

	// Base domain of the cluster. All DNS records must be sub-domains of this base and include the cluster name.
	BaseDNSDomain string `json:"base_dns_domain,omitempty"`

	// IP address block from which Pod IPs are allocated This block must not overlap with existing physical networks. These IP addresses are used for the Pod network, and if you need to access the Pods from an external network, configure load balancers and routers to manage the traffic.
	// Pattern: ^([0-9]{1,3}\.){3}[0-9]{1,3}\/[0-9]|[1-2][0-9]|3[0-2]?$
	ClusterNetworkCidr string `json:"cluster_network_cidr,omitempty"`

	// The subnet prefix length to assign to each individual node. For example, if clusterNetworkHostPrefix is set to 23, then each node is assigned a /23 subnet out of the given cidr (clusterNetworkCIDR), which allows for 510 (2^(32 - 23) - 2) pod IPs addresses. If you are required to provide access to nodes from an external network, configure load balancers and routers to manage the traffic.
	// Maximum: 32
	// Minimum: 1
	ClusterNetworkHostPrefix int64 `json:"cluster_network_host_prefix,omitempty"`

	// The time that this cluster was created.
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"created_at,omitempty" gorm:"type:datetime"`

	// Virtual IP used internally by the cluster for automating internal DNS requirements.
	// Format: hostname
	DNSVip strfmt.Hostname `json:"dns_vip,omitempty"`

	// Hosts that are associated with this cluster.
	Hosts []*Host `json:"hosts" gorm:"foreignkey:ClusterID;association_foreignkey:ID"`

	// Self link.
	// Required: true
	Href *string `json:"href"`

	// Unique identifier of the object.
	// Required: true
	// Format: uuid
	ID *strfmt.UUID `json:"id" gorm:"primary_key"`

	// Virtual IP used for cluster ingress traffic.
	// Format: hostname
	IngressVip strfmt.Hostname `json:"ingress_vip,omitempty"`

	// The time that this cluster completed installation.
	// Format: date-time
	InstallCompletedAt strfmt.DateTime `json:"install_completed_at,omitempty" gorm:"type:datetime;default:0"`

	// The time that this cluster began installation.
	// Format: date-time
	InstallStartedAt strfmt.DateTime `json:"install_started_at,omitempty" gorm:"type:datetime;default:0"`

	// Indicates the type of this object. Will be 'Cluster' if this is a complete object or 'ClusterLink' if it is just a link.
	// Required: true
	// Enum: [Cluster]
	Kind *string `json:"kind"`

	// Name of the OpenShift cluster.
	Name string `json:"name,omitempty"`

	// Version of the OpenShift cluster.
	OpenshiftVersion string `json:"openshift_version,omitempty"`

	// The pull secret that obtained from the Pull Secret page on the Red Hat OpenShift Cluster Manager site.
	PullSecret string `json:"pull_secret,omitempty" gorm:"type:varchar(4096)"`

	// The IP address pool to use for service IP addresses. You can enter only one IP address pool. If you need to access the services from an external network, configure load balancers and routers to manage the traffic.
	// Pattern: ^([0-9]{1,3}\.){3}[0-9]{1,3}\/[0-9]|[1-2][0-9]|3[0-2]?$
	ServiceNetworkCidr string `json:"service_network_cidr,omitempty"`

	// SSH public key for debugging OpenShift nodes.
	SSHPublicKey string `json:"ssh_public_key,omitempty" gorm:"type:varchar(1024)"`

	// Status of the OpenShift cluster.
	// Required: true
	// Enum: [insufficient ready error installing installed]
	Status *string `json:"status"`

	// Additional information pertaining to the status of the OpenShift cluster.
	// Required: true
	StatusInfo *string `json:"status_info"`

	// The last time that this cluster was updated.
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updated_at,omitempty" gorm:"type:datetime"`
}

// Validate validates this cluster
func (m *Cluster) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAPIVip(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateClusterNetworkCidr(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateClusterNetworkHostPrefix(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDNSVip(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHosts(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHref(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIngressVip(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInstallCompletedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInstallStartedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKind(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateServiceNetworkCidr(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatusInfo(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Cluster) validateAPIVip(formats strfmt.Registry) error {

	if swag.IsZero(m.APIVip) { // not required
		return nil
	}

	if err := validate.FormatOf("api_vip", "body", "hostname", m.APIVip.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateClusterNetworkCidr(formats strfmt.Registry) error {

	if swag.IsZero(m.ClusterNetworkCidr) { // not required
		return nil
	}

	if err := validate.Pattern("cluster_network_cidr", "body", string(m.ClusterNetworkCidr), `^([0-9]{1,3}\.){3}[0-9]{1,3}\/[0-9]|[1-2][0-9]|3[0-2]?$`); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateClusterNetworkHostPrefix(formats strfmt.Registry) error {

	if swag.IsZero(m.ClusterNetworkHostPrefix) { // not required
		return nil
	}

	if err := validate.MinimumInt("cluster_network_host_prefix", "body", int64(m.ClusterNetworkHostPrefix), 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("cluster_network_host_prefix", "body", int64(m.ClusterNetworkHostPrefix), 32, false); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("created_at", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateDNSVip(formats strfmt.Registry) error {

	if swag.IsZero(m.DNSVip) { // not required
		return nil
	}

	if err := validate.FormatOf("dns_vip", "body", "hostname", m.DNSVip.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateHosts(formats strfmt.Registry) error {

	if swag.IsZero(m.Hosts) { // not required
		return nil
	}

	for i := 0; i < len(m.Hosts); i++ {
		if swag.IsZero(m.Hosts[i]) { // not required
			continue
		}

		if m.Hosts[i] != nil {
			if err := m.Hosts[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("hosts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Cluster) validateHref(formats strfmt.Registry) error {

	if err := validate.Required("href", "body", m.Href); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateIngressVip(formats strfmt.Registry) error {

	if swag.IsZero(m.IngressVip) { // not required
		return nil
	}

	if err := validate.FormatOf("ingress_vip", "body", "hostname", m.IngressVip.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateInstallCompletedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.InstallCompletedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("install_completed_at", "body", "date-time", m.InstallCompletedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateInstallStartedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.InstallStartedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("install_started_at", "body", "date-time", m.InstallStartedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

var clusterTypeKindPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["Cluster"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		clusterTypeKindPropEnum = append(clusterTypeKindPropEnum, v)
	}
}

const (

	// ClusterKindCluster captures enum value "Cluster"
	ClusterKindCluster string = "Cluster"
)

// prop value enum
func (m *Cluster) validateKindEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, clusterTypeKindPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Cluster) validateKind(formats strfmt.Registry) error {

	if err := validate.Required("kind", "body", m.Kind); err != nil {
		return err
	}

	// value enum
	if err := m.validateKindEnum("kind", "body", *m.Kind); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateServiceNetworkCidr(formats strfmt.Registry) error {

	if swag.IsZero(m.ServiceNetworkCidr) { // not required
		return nil
	}

	if err := validate.Pattern("service_network_cidr", "body", string(m.ServiceNetworkCidr), `^([0-9]{1,3}\.){3}[0-9]{1,3}\/[0-9]|[1-2][0-9]|3[0-2]?$`); err != nil {
		return err
	}

	return nil
}

var clusterTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["insufficient","ready","error","installing","installed"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		clusterTypeStatusPropEnum = append(clusterTypeStatusPropEnum, v)
	}
}

const (

	// ClusterStatusInsufficient captures enum value "insufficient"
	ClusterStatusInsufficient string = "insufficient"

	// ClusterStatusReady captures enum value "ready"
	ClusterStatusReady string = "ready"

	// ClusterStatusError captures enum value "error"
	ClusterStatusError string = "error"

	// ClusterStatusInstalling captures enum value "installing"
	ClusterStatusInstalling string = "installing"

	// ClusterStatusInstalled captures enum value "installed"
	ClusterStatusInstalled string = "installed"
)

// prop value enum
func (m *Cluster) validateStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, clusterTypeStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Cluster) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Status); err != nil {
		return err
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", *m.Status); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateStatusInfo(formats strfmt.Registry) error {

	if err := validate.Required("status_info", "body", m.StatusInfo); err != nil {
		return err
	}

	return nil
}

func (m *Cluster) validateUpdatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.UpdatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("updated_at", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Cluster) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Cluster) UnmarshalBinary(b []byte) error {
	var res Cluster
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
