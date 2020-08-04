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

// GetHostReader is a Reader for the GetHost structure.
type GetHostReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetHostReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetHostOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetHostNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetHostInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetHostOK creates a GetHostOK with default headers values
func NewGetHostOK() *GetHostOK {
	return &GetHostOK{}
}

/*GetHostOK handles this case with default header values.

Success.
*/
type GetHostOK struct {
	Payload *models.Host
}

func (o *GetHostOK) Error() string {
	return fmt.Sprintf("[GET /clusters/{cluster_id}/hosts/{host_id}][%d] getHostOK  %+v", 200, o.Payload)
}

func (o *GetHostOK) GetPayload() *models.Host {
	return o.Payload
}

func (o *GetHostOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Host)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetHostNotFound creates a GetHostNotFound with default headers values
func NewGetHostNotFound() *GetHostNotFound {
	return &GetHostNotFound{}
}

/*GetHostNotFound handles this case with default header values.

Error.
*/
type GetHostNotFound struct {
	Payload *models.Error
}

func (o *GetHostNotFound) Error() string {
	return fmt.Sprintf("[GET /clusters/{cluster_id}/hosts/{host_id}][%d] getHostNotFound  %+v", 404, o.Payload)
}

func (o *GetHostNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetHostNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetHostInternalServerError creates a GetHostInternalServerError with default headers values
func NewGetHostInternalServerError() *GetHostInternalServerError {
	return &GetHostInternalServerError{}
}

/*GetHostInternalServerError handles this case with default header values.

Error.
*/
type GetHostInternalServerError struct {
	Payload *models.Error
}

func (o *GetHostInternalServerError) Error() string {
	return fmt.Sprintf("[GET /clusters/{cluster_id}/hosts/{host_id}][%d] getHostInternalServerError  %+v", 500, o.Payload)
}

func (o *GetHostInternalServerError) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetHostInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
