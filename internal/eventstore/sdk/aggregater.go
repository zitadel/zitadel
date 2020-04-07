package sdk

import (
	"fmt"

	"github.com/caos/zitadel/internal/errors"
)

var (
	_ AggregaterError = (*aggregaterError)(nil)
	_ errors.
		Error = (*aggregaterError)(nil)
)

type AggregaterError interface {
	error
	IsAggregater()
}

type aggregaterError struct {
	*errors.CaosError
}

func ThrowAggregater(parent error, id, message string) error {
	return &aggregaterError{errors.CreateCaosError(parent, id, message)}
}

func ThrowAggregaterf(parent error, id, format string, a ...interface{}) error {
	return ThrowAggregater(parent, id, fmt.Sprintf(format, a...))
}

func (err *aggregaterError) IsAggregater() {}

func IsAggregater(err error) bool {
	_, ok := err.(AggregaterError)
	return ok
}
