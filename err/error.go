package err

import (
	"context"
	"encoding/json"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	"net/http"

	"gitlab.intelligentb.com/devops/bplus/i18n"
)

type HttpCodeProvider interface {
	GetHttpCode() int
}

// BPlusError - defines the error structure of all return values
type BPlusError struct {
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

func (val BPlusError) GetHttpCode() int {
	return val.HTTPErrorCode
}

func (val BPlusError) Error() string {
	ret, _ := json.Marshal(val)
	return string(ret)
}

// Make403 - make a 403 error
func Make403(code int, message string) BPlusError {
	return BPlusError{HTTPErrorCode: 403, ErrorCode: code, ErrorMessage: message, LogLevel: Error}
}

// MakeErr - Make a generic error
func MakeErr(ctx context.Context, ll LogLevel, code int, msgCode string, args map[string]interface{}) BPlusError {
	return MakeErrWithHTTPCode(ctx, ll, code, msgCode, http.StatusInternalServerError, args)
}

// MakeErrWithHTTPCode - Make a generic error with http error code
func MakeErrWithHTTPCode(ctx context.Context, ll LogLevel, code int, msgCode string, hTTPError int, args map[string]interface{}) BPlusError {
	msg := msgCode

	//msg = fmt.Sprintf(message, params...)
	m := i18n.Translate(ctx, msg, args)
	if m != "" {
		msg = m
	}

	return BPlusError{
		ErrorCode:     code,
		ErrorMessage:  msg,
		LogLevel:      ll,
		HTTPErrorCode: hTTPError,
		TraceId:       bplusc.GetTraceId(ctx),
		TrajectoryID:  bplusc.GetTrajectoryID(ctx),
	}
}
