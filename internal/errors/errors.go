package errors

import (
	stderrors "errors"
	"fmt"
)

type ErrorCode string

const (
	ErrParse      ErrorCode = "PARSE_ERROR"
	ErrValidation ErrorCode = "VALIDATION_ERROR"
	ErrCloud      ErrorCode = "CLOUD_ERROR"
	ErrPlugin     ErrorCode = "PLUGIN_ERROR"
	ErrAuth       ErrorCode = "AUTH_ERROR"
	ErrPipeline   ErrorCode = "PIPELINE_ERROR"
	ErrConfig     ErrorCode = "CONFIG_ERROR"
	ErrDatabase   ErrorCode = "DATABASE_ERROR"
	ErrNetwork    ErrorCode = "NETWORK_ERROR"
	ErrInternal   ErrorCode = "INTERNAL_ERROR"
	ErrNotFound   ErrorCode = "NOT_FOUND"
	ErrConflict   ErrorCode = "CONFLICT"
)

type NaeosError struct {
	Code    ErrorCode
	Message string
	Inner   error
}

func (e *NaeosError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *NaeosError) Unwrap() error {
	return e.Inner
}

func New(code ErrorCode, msg string) *NaeosError {
	return &NaeosError{Code: code, Message: msg}
}

func Wrap(code ErrorCode, msg string, err error) *NaeosError {
	return &NaeosError{Code: code, Message: msg, Inner: err}
}

func Is(err error, code ErrorCode) bool {
	for err != nil {
		var ne *NaeosError
		if stderrors.As(err, &ne) && ne.Code == code {
			return true
		}
		uw, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = uw.Unwrap()
	}
	return false
}

var (
	ErrNotConnected   = New(ErrNetwork, "not connected")
	ErrInvalidSpec    = New(ErrValidation, "invalid spec")
	ErrPluginNotFound = New(ErrPlugin, "plugin not found")
	ErrDeployFailed   = New(ErrCloud, "deploy failed")
)
