package sdk

import (
	"fmt"

	"github.com/caos/zitadel/internal/errors"
)

var (
	_ AppendEventError = (*appendEventError)(nil)
	_ errors.Error     = (*appendEventError)(nil)
)

type AppendEventError interface {
	error
	IsAppendEventError()
}

type appendEventError struct {
	*errors.CaosError
}

func ThrowAppendEventError(parent error, id, message string) error {
	return &appendEventError{errors.CreateCaosError(parent, id, message)}
}

func ThrowAggregaterf(parent error, id, format string, a ...interface{}) error {
	return ThrowAppendEventError(parent, id, fmt.Sprintf(format, a...))
}

func (err *appendEventError) IsAppendEventError() {}

func IsAppendEventError(err error) bool {
	_, ok := err.(AppendEventError)
	return ok
}
