package err

import (
	"fmt"

	bpluse "github.com/MenaEnergyVentures/bplus/err"
)

// It is recommended that each module define its own error file

func internalMakeBplusError(ll bpluse.LogLevels, e BPlusErrorCode, args ...interface{}) bpluse.BPlusError {
	return bpluse.BPlusError{
		ErrorCode:    e,
		ErrorMessage: fmt.Sprintf(ErrMessages[e], args...),
		LogLevel:     ll,
	}

}

// MakeBplusError - returns a customized CAFUError for BPlus
func MakeBplusError(e BPlusErrorCode, args ...interface{}) bpluse.BPlusError {
	return internalMakeBplusError(bpluse.Error, e, args...)

}

// MakeBplusWarning - returns a customized CAFUError for BPlus
func MakeBplusWarning(e BPlusErrorCode, args ...interface{}) bpluse.BPlusError {
	return internalMakeBplusError(bpluse.Warning, e, args...)

}

// BPlusErrorCode - A B Plus error code
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
}
