// Code generated by go-swagger; DO NOT EDIT.

package public

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/zzpu/ums/internal/httpclient/models"
)

// InitializeSelfServiceRecoveryViaBrowserFlowReader is a Reader for the InitializeSelfServiceRecoveryViaBrowserFlow structure.
type InitializeSelfServiceRecoveryViaBrowserFlowReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *InitializeSelfServiceRecoveryViaBrowserFlowReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 302:
		result := NewInitializeSelfServiceRecoveryViaBrowserFlowFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewInitializeSelfServiceRecoveryViaBrowserFlowInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewInitializeSelfServiceRecoveryViaBrowserFlowFound creates a InitializeSelfServiceRecoveryViaBrowserFlowFound with default headers values
func NewInitializeSelfServiceRecoveryViaBrowserFlowFound() *InitializeSelfServiceRecoveryViaBrowserFlowFound {
	return &InitializeSelfServiceRecoveryViaBrowserFlowFound{}
}

/*InitializeSelfServiceRecoveryViaBrowserFlowFound handles this case with default header values.

Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.
*/
type InitializeSelfServiceRecoveryViaBrowserFlowFound struct {
}

func (o *InitializeSelfServiceRecoveryViaBrowserFlowFound) Error() string {
	return fmt.Sprintf("[GET /self-service/recovery/browser][%d] initializeSelfServiceRecoveryViaBrowserFlowFound ", 302)
}

func (o *InitializeSelfServiceRecoveryViaBrowserFlowFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewInitializeSelfServiceRecoveryViaBrowserFlowInternalServerError creates a InitializeSelfServiceRecoveryViaBrowserFlowInternalServerError with default headers values
func NewInitializeSelfServiceRecoveryViaBrowserFlowInternalServerError() *InitializeSelfServiceRecoveryViaBrowserFlowInternalServerError {
	return &InitializeSelfServiceRecoveryViaBrowserFlowInternalServerError{}
}

/*InitializeSelfServiceRecoveryViaBrowserFlowInternalServerError handles this case with default header values.

genericError
*/
type InitializeSelfServiceRecoveryViaBrowserFlowInternalServerError struct {
	Payload *models.GenericError
}

func (o *InitializeSelfServiceRecoveryViaBrowserFlowInternalServerError) Error() string {
	return fmt.Sprintf("[GET /self-service/recovery/browser][%d] initializeSelfServiceRecoveryViaBrowserFlowInternalServerError  %+v", 500, o.Payload)
}

func (o *InitializeSelfServiceRecoveryViaBrowserFlowInternalServerError) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *InitializeSelfServiceRecoveryViaBrowserFlowInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
