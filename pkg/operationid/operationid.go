package requestid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type operationIDKey string

const (
	headerKey              = "X-Operation-ID"
	ctxKey    operationIDKey = "operation-id"
)

// Middleware wraps an http handler.
// Before the inner handler is called, the X-Request-ID header is extracted
// and injected into the request context.
// It is safe to be called if the header does not exists
func Middleware(inner http.Handler) http.Handler {
	return handler{inner: inner}
}

// Transport wraps an inner transport.
// It adds the request-id from the request context to the request header.
func Transport(inner http.RoundTripper) http.RoundTripper {
	return transport{inner: inner}
}

// ApplyTransport injects the request-id transport to an http client
func ApplyTransport(client *http.Client) {
	inner := client.Transport
	if inner == nil {
		inner = http.DefaultTransport
	}
	client.Transport = Transport(inner)
}

// FromContext returns the request id stored in the the context
func FromContext(ctx context.Context) string {
	operationID := ctx.Value(ctxKey)
	if operationID == nil {
		operationID = ""
	}
	return operationID.(string)
}

func ToContext(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, ctxKey, operationID)
}

func OperationIDLogger(logger logrus.FieldLogger, operationID string) logrus.FieldLogger {
	return logger.WithField("operation_io", operationID)
}

type handler struct {
	inner http.Handler
}

func FromRequest(r *http.Request) string {
	return r.Header.Get(headerKey)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	operationID := NewID()
	ctx := ToContext(r.Context(), operationID)
	r = r.WithContext(ctx)
	w.Header().Set(h.id.HeaderKey, id)
	h.inner.ServeHTTP(w, r)
}

type transport struct {
	inner http.RoundTripper
}

func (t transport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err :=  t.inner.RoundTrip(r)
	if operationID := r.Context().Value(ctxKey); operationID != nil {
		r.Header.Set(headerKey, operationID.(string))
	}
	// log the operation ID from the response
	return resp, err
}

func NewID() string {
	return uuid.New().String()
}
