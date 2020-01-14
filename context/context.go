package context

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// utilities to manipulate the context

func Add(ctx context.Context, key string, val interface{}) context.Context {
	return context.WithValue(ctx, Key(key), val)
}

// Key - the key for storing information in the context
type Key string

func copyPathParams(ctx context.Context, r *http.Request) context.Context {
	for name, val := range mux.Vars(r) {
		ctx = Add(ctx, name, val)
	}
	return ctx
}

func copyQueryParams(ctx context.Context, r *http.Request) context.Context {
	for name, val := range r.URL.Query() {
		ctx = Add(ctx, name, val[0])
	}
	return ctx
}

func copyHTTPHeaders(ctx context.Context, r *http.Request) context.Context {
	for name, val := range r.Header {
		ctx = Add(ctx, name, val[0])
	}
	return ctx
}

// Value - extracts the value of key k from context
func Value(ctx context.Context, k string) interface{} {
	return ctx.Value(Key(k))
}

// Enhance - enhance the context with things from HTTP request
func Enhance(ctx context.Context, r *http.Request) context.Context {
	ctx = copyPathParams(ctx, r)
	ctx = copyQueryParams(ctx, r)
	ctx = copyHTTPHeaders(ctx, r)

	return ctx
}

// SetError - sets the error into the context and returns the enhanced context
func SetError(ctx context.Context, err error) context.Context {
	return Add(ctx, "RESPONSE_ERROR", err)
}

// GetError - gets the error from the context
func GetError(ctx context.Context) error {
	err := Value(ctx, "RESPONSE_ERROR")
	if err != nil {
		return err.(error)
	}
	return nil
}

// GetPayload - gets the payload from the context
func GetPayload(ctx context.Context) interface{} {
	return Value(ctx, "REQUEST_PAYLOAD")
}

// SetPayload - sets the payload and returns the enhance context
func SetPayload(ctx context.Context, payload interface{}) context.Context {
	return Add(ctx, "REQUEST_PAYLOAD", payload)
}

// GetResponsePayload - gets the payload from the context
func GetResponsePayload(ctx context.Context) interface{} {
	return Value(ctx, "RESPONSE_PAYLOAD")
}

// SetResponsePayload - sets the payload and returns the enhance context
func SetResponsePayload(ctx context.Context, payload interface{}) context.Context {
	return Add(ctx, "RESPONSE_PAYLOAD", payload)
}

// GetProxyRequestParams - gets the operation descriptor from the context
func GetProxyRequestParams(ctx context.Context) []interface{} {
	return Value(ctx, "PROXY_PARAMS").([]interface{})
}

// SetProxyRequestParams - sets the operation descriptor and returns the enhance context
func SetProxyRequestParams(ctx context.Context, params []interface{}) context.Context {
	return Add(ctx, "PROXY_PARAMS", params)
}

// GetProxyResponsePayload - gets the payload from the context
func GetProxyResponsePayload(ctx context.Context) interface{} {
	return Value(ctx, "PROXY_RESPONSE_PAYLOAD")
}

// SetProxyResponsePayload - sets the payload and returns the enhance context
func SetProxyResponsePayload(ctx context.Context, payload interface{}) context.Context {
	return Add(ctx, "PROXY_RESPONSE_PAYLOAD", payload)
}

// SetProxyResponseError - sets the error into the context and returns the enhanced context
func SetProxyResponseError(ctx context.Context, err error) context.Context {
	return Add(ctx, "PROXY_RESPONSE_ERROR", err)
}

// GetProxyResponseError - gets the error from the context
func GetProxyResponseError(ctx context.Context) error {
	err := Value(ctx, "PROXY_RESPONSE_ERROR")
	if err != nil {
		return err.(error)
	}
	return nil
}
