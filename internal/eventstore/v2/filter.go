package eventstore

import (
	"github.com/caos/zitadel/internal/errors"
)

type Filter struct {
	field     Field
	value     interface{}
	operation Operation
}

type Operation int32

const (
	Operation_Equals Operation = 1 + iota
	Operation_Greater
	Operation_Less
	Operation_In
)

type Field int32

const (
	Field_AggregateType Field = 1 + iota
	Field_AggregateID
	Field_LatestSequence
	Field_ResourceOwner
	Field_EditorService
	Field_EditorUser
	Field_EventType
)

//NewFilter is used in tests. Use searchQuery.*Filter() instead
func NewFilter(field Field, value interface{}, operation Operation) *Filter {
	return &Filter{
		field:     field,
		value:     value,
		operation: operation,
	}
}

func (f *Filter) Field() Field {
	return f.field
}
func (f *Filter) Operation() Operation {
	return f.operation
}
func (f *Filter) Value() interface{} {
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
