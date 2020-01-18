package err

import (
	"context"

	bpluse "github.com/MenaEnergyVentures/bplus/err"
)

// It is recommended that each module define its own error file

func internalMakeBplusError(ctx context.Context, ll bpluse.LogLevel, e BPlusErrorCode, args ...interface{}) bpluse.BPlusError {
	return bpluse.MakeErr(ctx, ll, e, ErrMessages[e], args...)
}

// MakeBplusError - returns a customized CAFUError for BPlus
func MakeBplusError(ctx context.Context, e BPlusErrorCode, args ...interface{}) bpluse.BPlusError {
	return internalMakeBplusError(ctx, bpluse.Error, e, args...)

}

// MakeBplusWarning - returns a customized CAFUError for BPlus
func MakeBplusWarning(ctx context.Context, e BPlusErrorCode, args ...interface{}) bpluse.BPlusError {
	return internalMakeBplusError(ctx, bpluse.Warning, e, args...)

}

// BPlusErrorCode - A BPlus error code
type BPlusErrorCode = int

// enumeration for B Plus Error codes
const (
	ServiceNotFound BPlusErrorCode = iota + 1000
	OperationNotFound
	DecodingError
	CannotGenerateHTTPRequest
	CannotGenerateHTTPRequest1
	CannotGenerateHTTPRequestForPayload
	ResponseUnmarshalException
	ParamsNotExpected
	HTTPCallFailed
	CannotReadResponseBody
	CannotMakeStateEntity
	ErrorInDecoding
	ErrorInAutoState
	AutoStateNotConfigured
	InvalidState
	InvalidEvent
	CannotReadFile
	EventNotFoundInRequest
	ParameterMissingInRequest
	ErrorInObtainingSTM
)

// ErrMessages - list of all messages corresponding to this code
var ErrMessages = map[BPlusErrorCode]string{
	ServiceNotFound:                     "Service %s is not found",
	OperationNotFound:                   "Operation %s in service %s not found",
	DecodingError:                       "Error in decoding the request. error = %s",
	CannotGenerateHTTPRequest:           "unable to generate HTTP request for param %#v. error is %s",
	CannotGenerateHTTPRequest1:          "unable to generate HTTP request. error is %s",
	CannotGenerateHTTPRequestForPayload: "Cannot construct the message from payload %s",
	ResponseUnmarshalException:          "Unable to unmarshal response payload.Error = %s",
	ParamsNotExpected:                   "#params passed does not match expected. Actual = %d. Expected = %d",
	HTTPCallFailed:                      "http call failed. err = %s",
	CannotReadResponseBody:              "cannot read response body err = %s",
	CannotMakeStateEntity:               "cannot make the state entity. error = %s",
	ErrorInDecoding:                     "Error in decoding request",
	ErrorInAutoState:                    "automatic state for %s threw error %s",
	AutoStateNotConfigured:              "autostate %s is not configured in the action catalog",
	InvalidState:                        "invalid state %s returned by the entity",
	InvalidEvent:                        "invalid event %s for the current state %s",
	CannotReadFile:                      "Cannot read file %s. Erorr = %s",
	EventNotFoundInRequest:              "Cannot find event in the request.",
	ParameterMissingInRequest:           "Parameter %s is missing in request",
	ErrorInObtainingSTM:                 "Error in Obtaining STM",
}
