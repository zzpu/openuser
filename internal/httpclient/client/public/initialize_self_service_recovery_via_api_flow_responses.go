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

// InitializeSelfServiceRecoveryViaAPIFlowReader is a Reader for the InitializeSelfServiceRecoveryViaAPIFlow structure.
type InitializeSelfServiceRecoveryViaAPIFlowReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *InitializeSelfServiceRecoveryViaAPIFlowReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewInitializeSelfServiceRecoveryViaAPIFlowOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewInitializeSelfServiceRecoveryViaAPIFlowBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewInitializeSelfServiceRecoveryViaAPIFlowInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewInitializeSelfServiceRecoveryViaAPIFlowOK creates a InitializeSelfServiceRecoveryViaAPIFlowOK with default headers values
func NewInitializeSelfServiceRecoveryViaAPIFlowOK() *InitializeSelfServiceRecoveryViaAPIFlowOK {
	return &InitializeSelfServiceRecoveryViaAPIFlowOK{}
}

/*InitializeSelfServiceRecoveryViaAPIFlowOK handles this case with default header values.

recoveryFlow
*/
type InitializeSelfServiceRecoveryViaAPIFlowOK struct {
	Payload *models.RecoveryFlow
}

func (o *InitializeSelfServiceRecoveryViaAPIFlowOK) Error() string {
	return fmt.Sprintf("[GET /self-service/recovery/api][%d] initializeSelfServiceRecoveryViaApiFlowOK  %+v", 200, o.Payload)
}

func (o *InitializeSelfServiceRecoveryViaAPIFlowOK) GetPayload() *models.RecoveryFlow {
	return o.Payload
}

func (o *InitializeSelfServiceRecoveryViaAPIFlowOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RecoveryFlow)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInitializeSelfServiceRecoveryViaAPIFlowBadRequest creates a InitializeSelfServiceRecoveryViaAPIFlowBadRequest with default headers values
func NewInitializeSelfServiceRecoveryViaAPIFlowBadRequest() *InitializeSelfServiceRecoveryViaAPIFlowBadRequest {
	return &InitializeSelfServiceRecoveryViaAPIFlowBadRequest{}
}

/*InitializeSelfServiceRecoveryViaAPIFlowBadRequest handles this case with default header values.

genericError
*/
type InitializeSelfServiceRecoveryViaAPIFlowBadRequest struct {
	Payload *models.GenericError
}

func (o *InitializeSelfServiceRecoveryViaAPIFlowBadRequest) Error() string {
	return fmt.Sprintf("[GET /self-service/recovery/api][%d] initializeSelfServiceRecoveryViaApiFlowBadRequest  %+v", 400, o.Payload)
}

func (o *InitializeSelfServiceRecoveryViaAPIFlowBadRequest) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *InitializeSelfServiceRecoveryViaAPIFlowBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInitializeSelfServiceRecoveryViaAPIFlowInternalServerError creates a InitializeSelfServiceRecoveryViaAPIFlowInternalServerError with default headers values
func NewInitializeSelfServiceRecoveryViaAPIFlowInternalServerError() *InitializeSelfServiceRecoveryViaAPIFlowInternalServerError {
	return &InitializeSelfServiceRecoveryViaAPIFlowInternalServerError{}
}

/*InitializeSelfServiceRecoveryViaAPIFlowInternalServerError handles this case with default header values.

genericError
*/
type InitializeSelfServiceRecoveryViaAPIFlowInternalServerError struct {
	Payload *models.GenericError
}

func (o *InitializeSelfServiceRecoveryViaAPIFlowInternalServerError) Error() string {
	return fmt.Sprintf("[GET /self-service/recovery/api][%d] initializeSelfServiceRecoveryViaApiFlowInternalServerError  %+v", 500, o.Payload)
}

func (o *InitializeSelfServiceRecoveryViaAPIFlowInternalServerError) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *InitializeSelfServiceRecoveryViaAPIFlowInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
