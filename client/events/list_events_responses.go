// Code generated by go-swagger; DO NOT EDIT.

package events

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/filanov/bm-inventory/models"
)

// ListEventsReader is a Reader for the ListEvents structure.
type ListEventsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListEventsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListEventsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewListEventsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListEventsOK creates a ListEventsOK with default headers values
func NewListEventsOK() *ListEventsOK {
	return &ListEventsOK{}
}

/*ListEventsOK handles this case with default header values.

Success.
*/
type ListEventsOK struct {
	Payload models.EventList
}

func (o *ListEventsOK) Error() string {
	return fmt.Sprintf("[GET /events/{entity_id}][%d] listEventsOK  %+v", 200, o.Payload)
}

func (o *ListEventsOK) GetPayload() models.EventList {
	return o.Payload
}

func (o *ListEventsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListEventsInternalServerError creates a ListEventsInternalServerError with default headers values
func NewListEventsInternalServerError() *ListEventsInternalServerError {
	return &ListEventsInternalServerError{}
}

/*ListEventsInternalServerError handles this case with default header values.

Error.
*/
type ListEventsInternalServerError struct {
	Payload *models.Error
}

func (o *ListEventsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /events/{entity_id}][%d] listEventsInternalServerError  %+v", 500, o.Payload)
}

func (o *ListEventsInternalServerError) GetPayload() *models.Error {
	return o.Payload
}

func (o *ListEventsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
