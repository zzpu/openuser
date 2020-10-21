// Code generated by go-swagger; DO NOT EDIT.

package admin

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/zzpu/ums/internal/httpclient/models"
)

// NewCreateRecoveryLinkParams creates a new CreateRecoveryLinkParams object
// with the default values initialized.
func NewCreateRecoveryLinkParams() *CreateRecoveryLinkParams {
	var ()
	return &CreateRecoveryLinkParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateRecoveryLinkParamsWithTimeout creates a new CreateRecoveryLinkParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateRecoveryLinkParamsWithTimeout(timeout time.Duration) *CreateRecoveryLinkParams {
	var ()
	return &CreateRecoveryLinkParams{

		timeout: timeout,
	}
}

// NewCreateRecoveryLinkParamsWithContext creates a new CreateRecoveryLinkParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateRecoveryLinkParamsWithContext(ctx context.Context) *CreateRecoveryLinkParams {
	var ()
	return &CreateRecoveryLinkParams{

		Context: ctx,
	}
}

// NewCreateRecoveryLinkParamsWithHTTPClient creates a new CreateRecoveryLinkParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateRecoveryLinkParamsWithHTTPClient(client *http.Client) *CreateRecoveryLinkParams {
	var ()
	return &CreateRecoveryLinkParams{
		HTTPClient: client,
	}
}

/*CreateRecoveryLinkParams contains all the parameters to send to the API endpoint
for the create recovery link operation typically these are written to a http.Request
*/
type CreateRecoveryLinkParams struct {

	/*Body*/
	Body *models.CreateRecoveryLink

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create recovery link params
func (o *CreateRecoveryLinkParams) WithTimeout(timeout time.Duration) *CreateRecoveryLinkParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create recovery link params
func (o *CreateRecoveryLinkParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create recovery link params
func (o *CreateRecoveryLinkParams) WithContext(ctx context.Context) *CreateRecoveryLinkParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create recovery link params
func (o *CreateRecoveryLinkParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create recovery link params
func (o *CreateRecoveryLinkParams) WithHTTPClient(client *http.Client) *CreateRecoveryLinkParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create recovery link params
func (o *CreateRecoveryLinkParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create recovery link params
func (o *CreateRecoveryLinkParams) WithBody(body *models.CreateRecoveryLink) *CreateRecoveryLinkParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create recovery link params
func (o *CreateRecoveryLinkParams) SetBody(body *models.CreateRecoveryLink) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateRecoveryLinkParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
