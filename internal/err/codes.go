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
)

// ErrMessages - list of all messages corresponding to this code
var ErrMessages = map[BPlusErrorCode]string{
	ServiceNotFound: "Service %s is not found",
}
