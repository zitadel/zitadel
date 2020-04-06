package models

import (
	"github.com/caos/zitadel/internal/errors"
)

type Filter struct {
	field     Field
	value     interface{}
	operation Operation
}

//NewFilter is used in tests. Use searchQuery.*Filter() instead
func NewFilter(field Field, value interface{}, operation Operation) *Filter {
	return &Filter{
		field:     field,
		value:     value,
		operation: operation,
	}
}

func (f *Filter) GetField() Field {
	return f.field
}
func (f *Filter) GetOperation() Operation {
	return f.operation
}
func (f *Filter) GetValue() interface{} {
	return f.value
}

func (f *Filter) Validate() error {
	if f == nil {
		return errors.ThrowPreconditionFailed(nil, "MODEL-z6KcG", "filter is nil")
	}
	if f.field <= 0 {
		return errors.ThrowPreconditionFailed(nil, "MODEL-zw62U", "field not definded")
	}
	if f.value == nil {
		return errors.ThrowPreconditionFailed(nil, "MODEL-GJ9ct", "no value definded")
	}
	if f.operation <= 0 {
		return errors.ThrowPreconditionFailed(nil, "MODEL-RrQTy", "operation not definded")
	}
	return nil
}
