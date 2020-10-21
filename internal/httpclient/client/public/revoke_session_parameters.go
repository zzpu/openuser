// Code generated by go-swagger; DO NOT EDIT.

package public

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

// NewRevokeSessionParams creates a new RevokeSessionParams object
// with the default values initialized.
func NewRevokeSessionParams() *RevokeSessionParams {
	var ()
	return &RevokeSessionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRevokeSessionParamsWithTimeout creates a new RevokeSessionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRevokeSessionParamsWithTimeout(timeout time.Duration) *RevokeSessionParams {
	var ()
	return &RevokeSessionParams{

		timeout: timeout,
	}
}

// NewRevokeSessionParamsWithContext creates a new RevokeSessionParams object
// with the default values initialized, and the ability to set a context for a request
func NewRevokeSessionParamsWithContext(ctx context.Context) *RevokeSessionParams {
	var ()
	return &RevokeSessionParams{

		Context: ctx,
	}
}

// NewRevokeSessionParamsWithHTTPClient creates a new RevokeSessionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRevokeSessionParamsWithHTTPClient(client *http.Client) *RevokeSessionParams {
	var ()
	return &RevokeSessionParams{
		HTTPClient: client,
	}
}

/*RevokeSessionParams contains all the parameters to send to the API endpoint
for the revoke session operation typically these are written to a http.Request
*/
type RevokeSessionParams struct {

	/*Body*/
	Body *models.RevokeSession

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the revoke session params
func (o *RevokeSessionParams) WithTimeout(timeout time.Duration) *RevokeSessionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the revoke session params
func (o *RevokeSessionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the revoke session params
func (o *RevokeSessionParams) WithContext(ctx context.Context) *RevokeSessionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the revoke session params
func (o *RevokeSessionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the revoke session params
func (o *RevokeSessionParams) WithHTTPClient(client *http.Client) *RevokeSessionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the revoke session params
func (o *RevokeSessionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the revoke session params
func (o *RevokeSessionParams) WithBody(body *models.RevokeSession) *RevokeSessionParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the revoke session params
func (o *RevokeSessionParams) SetBody(body *models.RevokeSession) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *RevokeSessionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
