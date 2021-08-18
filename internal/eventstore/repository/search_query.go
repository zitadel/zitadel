package repository

import "github.com/caos/zitadel/internal/errors"

//SearchQuery defines the which and how data are queried
type SearchQuery struct {
	Columns Columns
	Limit   uint64
	Desc    bool
	Filters [][]*Filter
}

//Columns defines which fields of the event are needed for the query
type Columns int32

const (
	//ColumnsEvent represents all fields of an event
	ColumnsEvent = iota + 1
	//ColumnsMaxSequence represents the latest sequence of the filtered events
	ColumnsMaxSequence

	columnsCount
)

func (c Columns) Validate() error {
	if c <= 0 || c >= columnsCount {
		return errors.ThrowPreconditionFailed(nil, "REPOS-x8R35", "column out of range")
	}
	return nil
}

//Filter represents all fields needed to compare a field of an event with a value
type Filter struct {
	Field     Field
	Value     interface{}
	Operation Operation
}

//Operation defines how fields are compared
type Operation int32

const (
	// OperationEquals compares two values for equality
	OperationEquals Operation = iota + 1
	// OperationGreater compares if the given values is greater than the stored one
	OperationGreater
	// OperationLess compares if the given values is less than the stored one
	OperationLess
	//OperationIn checks if a stored value matches one of the passed value list
	OperationIn
	//OperationJSONContains checks if a stored value matches the given json
	OperationJSONContains

	operationCount
)

//Field is the representation of a field from the event
type Field int32

const (
	//FieldAggregateType represents the aggregate type field
	FieldAggregateType Field = iota + 1
	//FieldAggregateID represents the aggregate id field
	FieldAggregateID
	//FieldSequence represents the sequence field
	FieldSequence
	//FieldResourceOwner represents the resource owner field
	FieldResourceOwner
	//FieldEditorService represents the editor service field
	FieldEditorService
	//FieldEditorUser represents the editor user field
	FieldEditorUser
	//FieldEventType represents the event type field
	FieldEventType
	//FieldEventData represents the event data field
	FieldEventData

	fieldCount
)

//NewFilter is used in tests. Use searchQuery.*Filter() instead
func NewFilter(field Field, value interface{}, operation Operation) *Filter {
	return &Filter{
		Field:     field,
		Value:     value,
		Operation: operation,
	}
}

//Validate checks if the fields of the filter have valid values
func (f *Filter) Validate() error {
	if f == nil {
		return errors.ThrowPreconditionFailed(nil, "REPO-z6KcG", "filter is nil")
	}
	if f.Field <= 0 || f.Field >= fieldCount {
		return errors.ThrowPreconditionFailed(nil, "REPO-zw62U", "field not definded")
	}
	if f.Value == nil {
		return errors.ThrowPreconditionFailed(nil, "REPO-GJ9ct", "no value definded")
	}
	if f.Operation <= 0 || f.Operation >= operationCount {
		return errors.ThrowPreconditionFailed(nil, "REPO-RrQTy", "operation not definded")
	}
	return nil
}
