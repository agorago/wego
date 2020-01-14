package err

import (
	bpluse "github.com/MenaEnergyVentures/bplus/err"
)

// It is recommended that each module define its own error file

// MakeBplusError - returns a customized CAFUError for BPlus
func MakeBplusError(e BPlusErrorCode) bpluse.CAFUError {
	return bpluse.CAFUError{
		ErrorCode:    e,
		ErrorMessage: ErrMessages[e],
		LogLevel:     bpluse.Error,
	}

}

// BPlusErrorCode - A B Plus error code
type BPlusErrorCode = int

// enumeration for B Plus Error codes
const (
	SERVICE_NOT_FOUND BPlusErrorCode = iota + 1000
)

// ErrMessages - list of all messages corresponding to this code
var ErrMessages = map[BPlusErrorCode]string{
	SERVICE_NOT_FOUND: "Service %s is not found",
}
