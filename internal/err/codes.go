package err

import (
	"context"

	bpluse "github.com/MenaEnergyVentures/bplus/err"
)

// It is recommended that each module define its own error file

func internalMakeBplusError(ctx context.Context, ll bpluse.LogLevel, e BPlusErrorCode, args map[string]interface{}) bpluse.BPlusError {
	return bpluse.MakeErr(ctx, ll, int(e), e.String(), args)
}

// MakeBplusError - returns a customized CAFUError for BPlus
func MakeBplusError(ctx context.Context, e BPlusErrorCode, args map[string]interface{}) bpluse.BPlusError {
	return internalMakeBplusError(ctx, bpluse.Error, e, args)

}

// MakeBplusWarning - returns a customized CAFUError for BPlus
func MakeBplusWarning(ctx context.Context, e BPlusErrorCode, args map[string]interface{}) bpluse.BPlusError {
	return internalMakeBplusError(ctx, bpluse.Warning, e, args)

}

// BPlusErrorCode - A BPlus error code
type BPlusErrorCode int

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

//go:generate stringer -type=BPlusErrorCode
