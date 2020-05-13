package err

import (
	"context"
	"encoding/json"
	wegocontext "github.com/agorago/wego/context"
	"net/http"

	"github.com/agorago/wego/i18n"
)

type HttpCodeProvider interface {
	GetHttpCode() int
}

// WeGOError - defines the error structure of all return values
type WeGOError struct {
	ErrorCode     int
	ErrorMessage  string
	HTTPErrorCode int
	LogLevel      LogLevel
	TrajectoryID  string
	UserID        string
	TraceId       string
}

// LogLevel - the different log levels
type LogLevel int

// Values for log levels
const (
	Error LogLevel = iota + 1
	Warning
)

func (val WeGOError) GetHttpCode() int {
	return val.HTTPErrorCode
}

func (val WeGOError) Error() string {
	ret, _ := json.Marshal(val)
	return string(ret)
}

// Make403 - make a 403 error
func Make403(code int, message string) WeGOError {
	return WeGOError{HTTPErrorCode: 403, ErrorCode: code, ErrorMessage: message, LogLevel: Error}
}

// MakeErr - Make a generic error
func MakeErr(ctx context.Context, ll LogLevel, code int, msgCode string, args map[string]interface{}) WeGOError {
	return MakeErrWithHTTPCode(ctx, ll, code, msgCode, http.StatusInternalServerError, args)
}

// MakeErrWithHTTPCode - Make a generic error with http error code
func MakeErrWithHTTPCode(ctx context.Context, ll LogLevel, code int, msgCode string, hTTPError int, args map[string]interface{}) WeGOError {
	msg := msgCode

	//msg = fmt.Sprintf(message, params...)
	m := i18n.Translate(ctx, msg, args)
	if m != "" {
		msg = m
	}

	return WeGOError{
		ErrorCode:     code,
		ErrorMessage:  msg,
		LogLevel:      ll,
		HTTPErrorCode: hTTPError,
		TraceId:       wegocontext.GetTraceId(ctx),
		TrajectoryID:  wegocontext.GetTrajectoryID(ctx),
	}
}
