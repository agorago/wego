package err

import (
	"context"
	"net/http"

	wegoe "github.com/agorago/wego/err"
)

// It is recommended that each module define its own error file

func internalMakeWegoError(ctx context.Context, ll wegoe.LogLevel, e WegoErrorCode, httpErrorCode int, args map[string]interface{}) wegoe.WeGOError {
	return wegoe.MakeErrWithHTTPCode(ctx, ll, int(e), e.String(), httpErrorCode, args)
}

// MakeBplusError - returns a customized CAFUError for BPlus
func MakeBplusError(ctx context.Context, e WegoErrorCode, args map[string]interface{}) wegoe.WeGOError {
	return internalMakeWegoError(ctx, wegoe.Error, e, http.StatusInternalServerError, args)

}

// MakeBplusWarning - returns a customized CAFUError for BPlus
func MakeBplusWarning(ctx context.Context, e WegoErrorCode, args map[string]interface{}) wegoe.WeGOError {
	return internalMakeWegoError(ctx, wegoe.Warning, e, http.StatusInternalServerError, args)

}

// MakeBplusErrorWithErrorCode - returns a customized CAFUError for BPlus
func MakeBplusErrorWithErrorCode(ctx context.Context, httpErrorCode int, e WegoErrorCode, args map[string]interface{}) wegoe.WeGOError {
	return internalMakeWegoError(ctx, wegoe.Error, e, httpErrorCode, args)

}

// MakeBplusWarningWithErrorCode - returns a customized CAFUError for BPlus
func MakeBplusWarningWithErrorCode(ctx context.Context, httpErrorCode int, e WegoErrorCode, args map[string]interface{}) wegoe.WeGOError {
	return internalMakeWegoError(ctx, wegoe.Warning, e, httpErrorCode, args)

}

// WegoErrorCode - A BPlus error code
type WegoErrorCode int

// enumeration for B Plus Error codes
const (
	ServiceNotFound                     WegoErrorCode = iota + 1000 // wego.errors.ServiceNotFound
	OperationNotFound                                               // wego.errors.OperationNotFound
	DecodingError                                                   // wego.errors.DecodingError
	CannotGenerateHTTPRequest                                       // wego.errors.CannotGenerateHTTPRequest
	CannotGenerateHTTPRequest1                                      // wego.errors.CannotGenerateHTTPRequest1
	CannotGenerateHTTPRequestForPayload                             // wego.errors
	ResponseUnmarshalException                                      // wego.errors.CannotGenerateHTTPRequestForPayload
	ParamsNotExpected                                               // wego.errors.ParamsNotExpected
	HTTPCallFailed                                                  // wego.errors.HTTPCallFailed
	CannotReadResponseBody                                          // wego.errors.CannotReadResponseBody
	CannotMakeStateEntity                                           // wego.errors.CannotMakeStateEntity
	ErrorInDecoding                                                 // wego.errors.ErrorInDecoding
	ErrorInAutoState                                                // wego.errors.ErrorInAutoState
	AutoStateNotConfigured                                          // wego.errors.AutoStateNotConfigured
	InvalidState                                                    // wego.errors.InvalidState
	InvalidEvent                                                    // wego.errors.InvalidEvent
	CannotReadFile                                                  // wego.errors.CannotReadFile
	EventNotFoundInRequest                                          // wego.errors.EventNotFoundInRequest
	ParameterMissingInRequest                                       // wego.errors.ParameterMissingInRequest
	ErrorInObtainingSTM                                             // wego.errors.ErrorInObtainingSTM
	Non200StatusCodeReturned                                        // wego.errors.Non200StatusCodeReturned
	ValidationError                                                 // wego.errors.ValidationError
	UnparseableFile                                                 //wego.error.UnparseableFile
	ErrorInInvokingService                                          //wego.error.ErrorInInvokingService
)

//go:generate stringer -linecomment -type=WegoErrorCode