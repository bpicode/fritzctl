// Package errors is a simplified variant of github.com/pkg/errors.
package errors

import (
	"fmt"
)

// Wrapf creates a new error with cause, and an inherent error message specified by the format parameters.
// If the cause is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withCause{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

type withCause struct {
	cause error
	msg   string
}

// Error makes *withCause an error type. The whole causal chain is concatenated with ": ".
func (wc *withCause) Error() string {
	return wc.msg + ": " + wc.cause.Error()
}

// Cause returns the wrapped error, see Wrapf.
func (wc *withCause) Cause() error {
	return wc.cause
}

// Msg returns the "bare" error message. It differs from Error in that the causal chain is omitted.
func (wc *withCause) Msg() string {
	return wc.msg
}
