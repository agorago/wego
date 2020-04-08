package context

import (
	"fmt"
	"context"
	"github.com/google/uuid"
	"net/http"

	"github.com/gorilla/mux"
)

// Bunch of constants defining the stuff that is stored in the Context
// There are convenience functions below to set and access these values
// These should not used in the service or middleware
// Instead use the associated convenience functions. Hence all of them are private
const (
	responseError        = "BPLUS-RESPONSE-ERROR"
	requestPayload       = "BPLUS-REQUEST-PAYLOAD"
	responsePayload      = "BPLUS-RESPONSE-PAYLOAD"
	proxyParams          = "BPLUS-PROXY-PARAMS"
	proxyResponsePayload = "BPLUS-PROXY-RESPONSE-PAYLOAD"
	proxyResponseError   = "BPLUS-PROXY-RESPONSE-ERROR"

	RemoteAddr       = "BPLUS-REMOTE-ADDRESS"
	RequestURI       = "BPLUS-REQUEST-URI"
	URL              = "BPLUS-URL"
	Method           = "BPLUS-METHOD"
	TransferEncoding = "BPLUS-TRANSFER-ENCODING"
	ContentLength    = "BPLUS-CONTENT-LENGTH"
	Host             = "BPLUS-HOST"
	TraceID          = "BPLUS-TRACE_ID"
	TrajectoryID     = "BPLUS-TRAJECTORY-ID"

	allKeys = "BPLUS-ALL-KEYS"
)

// utilities to manipulate the context

// Add - Add a key to context. Store this key in a separate array within context so
// that we can iterate through all available keys if required
// context.Context does not provide a way to iterate through the keys added in it
func Add(ctx context.Context, key string, val interface{}) context.Context {
	keys, _ := Value(ctx, allKeys).([]string)
	keys = append(keys, key)
	ctx = context.WithValue(ctx, Key(allKeys), keys)

	return context.WithValue(ctx, Key(key), val)
}

//GetAllKeys - get a list of all the keys in the given context
func GetAllKeys(ctx context.Context) []string {
	ret, ok := Value(ctx, allKeys).([]string)
	if ok {
		return ret
	}
	return nil
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

// CopyHeadersToHTTPRequest - copy the context parameters to the http request as
// headers
func CopyHeadersToHTTPRequest(ctx context.Context, req *http.Request) {
	for _, s := range GetAllKeys(ctx) {
		val, ok := Value(ctx, s).(string)
		if ok {
			req.Header.Set(s, val)
		}
	}
}

// Enhance - enhance the context with things from HTTP request
func Enhance(ctx context.Context, r *http.Request) context.Context {
	ctx = copyPathParams(ctx, r)
	ctx = copyQueryParams(ctx, r)
	ctx = copyHTTPHeaders(ctx, r)
	ctx = copyStandardHTTPHeaders(ctx,r)
	ctx = generateTraceID(ctx)
	return ctx
}

func copyStandardHTTPHeaders(ctx context.Context, r *http.Request) context.Context {
	ctx = Add(ctx,RemoteAddr,r.RemoteAddr)
	ctx = Add(ctx,RequestURI,r.RequestURI)
	ctx = Add(ctx,URL,r.URL)
	ctx = Add(ctx,Method,r.Method)
	ctx = Add(ctx,TransferEncoding,r.TransferEncoding)
	ctx = Add(ctx,ContentLength,r.ContentLength)
	ctx = Add(ctx,Host,r.Host)

	return ctx
}

func generateTraceID(ctx context.Context)context.Context{
	tr := GetTraceId(ctx)
	if tr == ""  {
		t := uuid.New().String()
		fmt.Printf("Setting traceId to %s\n",t)
		return setTraceId(ctx,t)
	}
	return ctx
}

// GetTraceId - returns the trace ID stored in the context
func GetTraceId(ctx context.Context) string{
	tr,ok :=  Value(ctx,TraceID).(string)
	if ok {
		return tr
	}
	return ""
}

// setTraceId - returns a context with a traceID set
func setTraceId(ctx context.Context,traceId string)context.Context{
	return Add(ctx,TraceID,traceId)
}

// GetTrajectoryId - returns the trace ID stored in the context
func GetTrajectoryID(ctx context.Context) string{
	tr,ok :=  Value(ctx,TrajectoryID).(string)
	if ok {
		return tr
	}
	return ""
}

// setTraceId - returns a context with a traceID set
func SetTrajectoryID(ctx context.Context,t string)context.Context{
	return Add(ctx,TrajectoryID,t)
}

// SetError - sets the error into the context and returns the enhanced context
func SetError(ctx context.Context, err error) context.Context {
	return Add(ctx, responseError, err)
}

// GetError - gets the error from the context
func GetError(ctx context.Context) error {
	err := Value(ctx, responseError)
	if err != nil {
		return err.(error)
	}
	return nil
}

// GetPayload - gets the payload from the context
func GetPayload(ctx context.Context) interface{} {
	return Value(ctx, requestPayload)
}

// SetPayload - sets the payload and returns the enhance context
func SetPayload(ctx context.Context, payload interface{}) context.Context {
	return Add(ctx, requestPayload, payload)
}

// GetResponsePayload - gets the payload from the context
func GetResponsePayload(ctx context.Context) interface{} {
	return Value(ctx, responsePayload)
}

// SetResponsePayload - sets the payload and returns the enhance context
func SetResponsePayload(ctx context.Context, payload interface{}) context.Context {
	return Add(ctx, responsePayload, payload)
}

// GetProxyRequestParams - gets the operation descriptor from the context
func GetProxyRequestParams(ctx context.Context) []interface{} {
	return Value(ctx, proxyParams).([]interface{})
}

// SetProxyRequestParams - sets the operation descriptor and returns the enhance context
func SetProxyRequestParams(ctx context.Context, params []interface{}) context.Context {
	return Add(ctx, proxyParams, params)
}

// GetProxyResponsePayload - gets the payload from the context
func GetProxyResponsePayload(ctx context.Context) interface{} {
	return Value(ctx, proxyResponsePayload)
}

// SetProxyResponsePayload - sets the payload and returns the enhance context
func SetProxyResponsePayload(ctx context.Context, payload interface{}) context.Context {
	return Add(ctx, proxyResponsePayload, payload)
}

// SetProxyResponseError - sets the error into the context and returns the enhanced context
func SetProxyResponseError(ctx context.Context, err error) context.Context {
	return Add(ctx, proxyResponseError, err)
}

// GetProxyResponseError - gets the error from the context
func GetProxyResponseError(ctx context.Context) error {
	err := Value(ctx, proxyResponseError)
	if err != nil {
		return err.(error)
	}
	return nil
}
