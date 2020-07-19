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

// PostStepReplyReader is a Reader for the PostStepReply structure.
type PostStepReplyReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostStepReplyReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewPostStepReplyNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostStepReplyBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewPostStepReplyNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPostStepReplyInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPostStepReplyNoContent creates a PostStepReplyNoContent with default headers values
func NewPostStepReplyNoContent() *PostStepReplyNoContent {
	return &PostStepReplyNoContent{}
}

/*PostStepReplyNoContent handles this case with default header values.

Success.
*/
type PostStepReplyNoContent struct {
}

func (o *PostStepReplyNoContent) Error() string {
	return fmt.Sprintf("[POST /clusters/{cluster_id}/hosts/{host_id}/instructions][%d] postStepReplyNoContent ", 204)
}

func (o *PostStepReplyNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPostStepReplyBadRequest creates a PostStepReplyBadRequest with default headers values
func NewPostStepReplyBadRequest() *PostStepReplyBadRequest {
	return &PostStepReplyBadRequest{}
}

/*PostStepReplyBadRequest handles this case with default header values.

Error.
*/
type PostStepReplyBadRequest struct {
	Payload *models.Error
}

func (o *PostStepReplyBadRequest) Error() string {
	return fmt.Sprintf("[POST /clusters/{cluster_id}/hosts/{host_id}/instructions][%d] postStepReplyBadRequest  %+v", 400, o.Payload)
}

func (o *PostStepReplyBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostStepReplyBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostStepReplyNotFound creates a PostStepReplyNotFound with default headers values
func NewPostStepReplyNotFound() *PostStepReplyNotFound {
	return &PostStepReplyNotFound{}
}

/*PostStepReplyNotFound handles this case with default header values.

Error.
*/
type PostStepReplyNotFound struct {
	Payload *models.Error
}

func (o *PostStepReplyNotFound) Error() string {
	return fmt.Sprintf("[POST /clusters/{cluster_id}/hosts/{host_id}/instructions][%d] postStepReplyNotFound  %+v", 404, o.Payload)
}

func (o *PostStepReplyNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostStepReplyNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostStepReplyInternalServerError creates a PostStepReplyInternalServerError with default headers values
func NewPostStepReplyInternalServerError() *PostStepReplyInternalServerError {
	return &PostStepReplyInternalServerError{}
}

/*PostStepReplyInternalServerError handles this case with default header values.

Error.
*/
type PostStepReplyInternalServerError struct {
	Payload *models.Error
}

func (o *PostStepReplyInternalServerError) Error() string {
	return fmt.Sprintf("[POST /clusters/{cluster_id}/hosts/{host_id}/instructions][%d] postStepReplyInternalServerError  %+v", 500, o.Payload)
}

func (o *PostStepReplyInternalServerError) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostStepReplyInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
