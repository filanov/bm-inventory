// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/filanov/bm-inventory/models"
)

// UpdateHostClusterReader is a Reader for the UpdateHostCluster structure.
type UpdateHostClusterReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateHostClusterReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateHostClusterOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewUpdateHostClusterNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateHostClusterInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUpdateHostClusterOK creates a UpdateHostClusterOK with default headers values
func NewUpdateHostClusterOK() *UpdateHostClusterOK {
	return &UpdateHostClusterOK{}
}

/*UpdateHostClusterOK handles this case with default header values.

Success.
*/
type UpdateHostClusterOK struct {
}

func (o *UpdateHostClusterOK) Error() string {
	return fmt.Sprintf("[POST /clusters/{cluster_id}/hosts/{host_id}/actions/move][%d] updateHostClusterOK ", 200)
}

func (o *UpdateHostClusterOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateHostClusterNotFound creates a UpdateHostClusterNotFound with default headers values
func NewUpdateHostClusterNotFound() *UpdateHostClusterNotFound {
	return &UpdateHostClusterNotFound{}
}

/*UpdateHostClusterNotFound handles this case with default header values.

Error.
*/
type UpdateHostClusterNotFound struct {
	Payload *models.Error
}

func (o *UpdateHostClusterNotFound) Error() string {
	return fmt.Sprintf("[POST /clusters/{cluster_id}/hosts/{host_id}/actions/move][%d] updateHostClusterNotFound  %+v", 404, o.Payload)
}

func (o *UpdateHostClusterNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *UpdateHostClusterNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateHostClusterInternalServerError creates a UpdateHostClusterInternalServerError with default headers values
func NewUpdateHostClusterInternalServerError() *UpdateHostClusterInternalServerError {
	return &UpdateHostClusterInternalServerError{}
}

/*UpdateHostClusterInternalServerError handles this case with default header values.

Error.
*/
type UpdateHostClusterInternalServerError struct {
	Payload *models.Error
}

func (o *UpdateHostClusterInternalServerError) Error() string {
	return fmt.Sprintf("[POST /clusters/{cluster_id}/hosts/{host_id}/actions/move][%d] updateHostClusterInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateHostClusterInternalServerError) GetPayload() *models.Error {
	return o.Payload
}

func (o *UpdateHostClusterInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}