package err

import "encoding/json"

// CAFUError - defines the error structure of all return values
type CAFUError struct {
	ErrorCode     int
	ErrorMessage  string
	HTTPErrorCode int
	LogLevel      LogLevels
	TrajectoryID  string
}

// LogLevels - the different log levels
type LogLevels int

// Values for log levels
const (
	Error LogLevels = iota + 1
)

func (val CAFUError) Error() string {
	ret, _ := json.Marshal(val)
	return string(ret)
}

func make403(code int, message string) CAFUError {
	return CAFUError{HTTPErrorCode: 403, ErrorCode: code, ErrorMessage: message, LogLevel: Error}
}
