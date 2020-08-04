// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeregisterHostHandlerFunc turns a function with the right signature into a deregister host handler
type DeregisterHostHandlerFunc func(DeregisterHostParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeregisterHostHandlerFunc) Handle(params DeregisterHostParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeregisterHostHandler interface for that can handle valid deregister host params
type DeregisterHostHandler interface {
	Handle(DeregisterHostParams, interface{}) middleware.Responder
}

// NewDeregisterHost creates a new http.Handler for the deregister host operation
func NewDeregisterHost(ctx *middleware.Context, handler DeregisterHostHandler) *DeregisterHost {
	return &DeregisterHost{Context: ctx, Handler: handler}
}

/*DeregisterHost swagger:route DELETE /clusters/{cluster_id}/hosts/{host_id} installer deregisterHost

Deregisters an OpenShift bare metal host.

*/
type DeregisterHost struct {
	Context *middleware.Context
	Handler DeregisterHostHandler
}

func (o *DeregisterHost) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeregisterHostParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
