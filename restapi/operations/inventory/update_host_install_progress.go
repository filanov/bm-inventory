// Code generated by go-swagger; DO NOT EDIT.

package inventory

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateHostInstallProgressHandlerFunc turns a function with the right signature into a update host install progress handler
type UpdateHostInstallProgressHandlerFunc func(UpdateHostInstallProgressParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateHostInstallProgressHandlerFunc) Handle(params UpdateHostInstallProgressParams) middleware.Responder {
	return fn(params)
}

// UpdateHostInstallProgressHandler interface for that can handle valid update host install progress params
type UpdateHostInstallProgressHandler interface {
	Handle(UpdateHostInstallProgressParams) middleware.Responder
}

// NewUpdateHostInstallProgress creates a new http.Handler for the update host install progress operation
func NewUpdateHostInstallProgress(ctx *middleware.Context, handler UpdateHostInstallProgressHandler) *UpdateHostInstallProgress {
	return &UpdateHostInstallProgress{Context: ctx, Handler: handler}
}

/*UpdateHostInstallProgress swagger:route PUT /clusters/{clusterId}/hosts/{hostId}/progress inventory updateHostInstallProgress

Update installation progress

*/
type UpdateHostInstallProgress struct {
	Context *middleware.Context
	Handler UpdateHostInstallProgressHandler
}

func (o *UpdateHostInstallProgress) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateHostInstallProgressParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
