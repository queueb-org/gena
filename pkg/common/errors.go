package common

import (
	"fmt"
	"strings"
)

// ErrStack is a simple accumulator for errors.
// Accumulates non-nil errors. Use [ErrStack.Collect] method
// to return wrapped collected errors or nil.
type ErrStack struct {
	errs []error
}

// Add error if it's not blank.
func (e *ErrStack) Add(errs ...error) {
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
}

// Collect returns error, if error is considered as blank return nil.
func (e *ErrStack) Collect() error {
	if e == nil || len(e.errs) == 0 {
		return nil
	}

	return e
}

// Unwrap implements wrapping error interface.
func (e *ErrStack) Unwrap() (err error) {
	if len(e.errs) != 0 {
		err = e.errs[0]
	}
	return
}

// Error implements [error] interface.
func (e *ErrStack) Error() string {
	builder := &strings.Builder{}
	last := len(e.errs) - 1

	for idx, err := range e.errs {
		_, _ = builder.WriteString(fmt.Sprintf("[%d] %v", idx, err))
		if idx < last {
			builder.WriteString("; ")
		}
	}

	return builder.String()
}

// helpers

// ErrorString returns error string, if error is blank returns "".
func ErrorString(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
