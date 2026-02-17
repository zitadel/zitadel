package database

import (
	"reflect"
	"slices"
)

type coalesce struct {
	values []any
}

func Coalesce(values ...any) *coalesce {
	return &coalesce{values: values}
}

// WriteArg implements [Change].
func (c *coalesce) WriteArg(builder *StatementBuilder) {
	c.WriteQualified(builder)
}

// IsOnColumn implements [Change].
func (c coalesce) IsOnColumn(col Column) bool {
	for _, value := range c.values {
		if change, ok := value.(Change); ok {
			if change.IsOnColumn(col) {
				return true
			}
		}
	}
	return false
}

// Write implements [Change].
func (c coalesce) Write(builder *StatementBuilder) {
	builder.WriteString("COALESCE(")
	for i, value := range c.values {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteArg(value)
	}
	builder.WriteString(")")
}

// Equals implements [Column].
func (c *coalesce) Equals(col Column) bool {
	if col == nil {
		return c == nil
	}
	other, ok := col.(*coalesce)
	if !ok || len(other.values) != len(c.values) {
		return false
	}
	for i, value := range c.values {
		if value != other.values[i] {
			return false
		}
	}
	return true
}

// Matches implements [Column].
func (c *coalesce) Matches(x any) bool {
	toMatch, ok := x.(*coalesce)
	if !ok || len(toMatch.values) != len(c.values) {
		return false
	}
	for i, value := range c.values {
		if value != toMatch.values[i] {
			return false
		}
	}
	return true
}

// String implements [Column].
func (c *coalesce) String() string {
	return "database.coalesce"
}

// WriteQualified implements [Column].
func (c *coalesce) WriteQualified(builder *StatementBuilder) {
	builder.WriteString("COALESCE(")
	for i, value := range c.values {
		if i > 0 {
			builder.WriteString(", ")
		}
		if col, ok := value.(Column); ok {
			col.WriteQualified(builder)
			continue
		}
		builder.WriteArg(value)
	}
	builder.WriteString(")")
}

// WriteUnqualified implements [Column].
func (c *coalesce) WriteUnqualified(builder *StatementBuilder) {
	builder.WriteString("COALESCE(")
	for i, value := range c.values {
		if i > 0 {
			builder.WriteString(", ")
		}
		if col, ok := value.(Column); ok {
			col.WriteUnqualified(builder)
			continue
		}
		builder.WriteArg(value)
	}
	builder.WriteString(")")
}

var (
	_ Column = (*coalesce)(nil)
	_ Change = (*coalesce)(nil)
)

type PatchJSONB interface {
	Change
	Column
}

// SetJSONValue creates a new JSON change that sets the given value at the given path in the JSON column.
//
// Example to set a new value at path {"a", "b"} in column "data":
//
//	SetJSONValue(NewColumn("table", "data"), []string{"a", "b"}, "newValue")
//
// Example to remove or reset a value at path {"a", "b"} in column "data":
//
//	SetJSONValue[any](NewColumn("table", "data"), []string{"a", "b"}, nil)
//
// Please see [AppendToJSONArray] and [RemoveFromJSONArray] for changes that add or remove values from JSON arrays.
func SetJSONValue[V Value](col Column, path []string, value V) Change {
	return &setJSONB[V]{
		column: col,
		value:  value,
		path:   path,
	}
}

// AppendToJSONArray creates a new JSON change that appends the given value to the array at the given path in the JSON column.
// If the path does not exist, it will be created as an array with the given value as its first element.
// If the path exists but is not an array, the change will fail.
func AppendToJSONArray[V Value](col Column, path []string, value V) Change {
	return &setJSONB[V]{
		column:  col,
		value:   value,
		path:    path,
		isArray: true,
	}
}

// RemoveFromJSONArray creates a new JSON change that removes the given value from the array at the given path in the JSON column.
// If the path does not exist or the value does not exist in the array, the change will do nothing.
// If the path exists but is not an array, the change will fail.
func RemoveFromJSONArray[V Value](col Column, path []string, value V) Change {
	return &setJSONB[V]{
		column:             col,
		value:              value,
		path:               path,
		isArray:            true,
		deleteArrayElement: true,
	}
}

type setJSONB[V Value] struct {
	column             Column
	path               []string
	value              V
	isArray            bool
	deleteArrayElement bool
}

// WriteArg implements [Change].
func (s setJSONB[V]) WriteArg(builder *StatementBuilder) {
	s.Write(builder)
}

// IsOnColumn implements [Change].
func (s setJSONB[V]) IsOnColumn(col Column) bool {
	return s.column.Equals(col)
}

// Write implements [Change].
func (s setJSONB[V]) Write(builder *StatementBuilder) {
	builder.WriteString("zitadel.jsonb_patch(")
	s.column.WriteQualified(builder)
	builder.WriteString(", ")
	builder.WriteArgs(s.path, s.value, s.isArray, s.deleteArrayElement)
	builder.WriteString(")")
}

// Equals implements [Column].
func (s *setJSONB[V]) Equals(col Column) bool {
	if col == nil {
		return s == nil
	}
	other, ok := col.(*setJSONB[V])
	if !ok {
		return false
	}
	if !s.column.Equals(other.column) {
		return false
	}
	if !slices.Equal(s.path, other.path) {
		return false
	}
	if !reflect.DeepEqual(s.value, other.value) {
		return false
	}
	if s.isArray != other.isArray {
		return false
	}
	if s.deleteArrayElement != other.deleteArrayElement {
		return false
	}
	return true
}

// Matches implements [Change].
func (s *setJSONB[V]) Matches(x any) bool {
	toMatch, ok := x.(*setJSONB[V])
	if !ok {
		return false
	}
	return s.Equals(toMatch)
}

// String implements [Change].
func (s *setJSONB[V]) String() string {
	return "database.setJSONB"
}

// WriteQualified implements [Column].
func (s setJSONB[V]) WriteQualified(builder *StatementBuilder) {
	s.Write(builder)
}

// WriteUnqualified implements [Column].
func (s setJSONB[V]) WriteUnqualified(builder *StatementBuilder) {
	s.Write(builder)
}

var (
	_ PatchJSONB = (*setJSONB[Value])(nil)
)
