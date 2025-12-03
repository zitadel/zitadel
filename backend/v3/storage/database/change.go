package database

import (
	"reflect"
	"slices"

	"go.uber.org/mock/gomock"
)

// Change represents a change to a column in a database table.
// Its written in the SET clause of an UPDATE statement.
type Change interface {
	gomock.Matcher
	// Write writes the change to the given statement builder.
	Write(builder *StatementBuilder)
	// IsOnColumn checks if the change is on the given column.
	IsOnColumn(col Column) bool
}

type change[V Value] struct {
	column Column
	value  V
}

// Matches implements [gomock.Matcher].
func (c *change[V]) Matches(x any) bool {
	toMatch, ok := x.(*change[V])
	if !ok {
		return false
	}
	colMatch := c.column.Equals(toMatch.column)
	valueMatch := reflect.DeepEqual(c.value, toMatch.value)
	return colMatch && valueMatch
}

// String implements [gomock.Matcher].
func (c *change[V]) String() string {
	return "database.change"
}

var (
	_ Change         = (*change[string])(nil)
	_ gomock.Matcher = (*change[string])(nil)
)

// NewChange creates a new Change for the given column and value.
// If you want to set a column to NULL, use [NewChangePtr].
func NewChange[V Value](col Column, value V) Change {
	return &change[V]{
		column: col,
		value:  value,
	}
}

// NewChangePtr creates a new Change for the given column and value pointer.
// If the value pointer is nil, the column will be set to NULL.
func NewChangePtr[V Value](col Column, value *V) Change {
	if value == nil {
		return NewChange(col, NullInstruction)
	}
	return NewChange(col, *value)
}

// Write implements [Change].
func (c change[V]) Write(builder *StatementBuilder) {
	c.column.WriteUnqualified(builder)
	builder.WriteString(" = ")
	builder.WriteArg(c.value)
}

// IsOnColumn implements [Change].
func (c change[V]) IsOnColumn(col Column) bool {
	return c.column.Equals(col)
}

type Changes []Change

func NewChanges(cols ...Change) Change {
	return Changes(cols)
}

// IsOnColumn implements [Change].
func (c Changes) IsOnColumn(col Column) bool {
	return slices.ContainsFunc(c, func(change Change) bool {
		return change.IsOnColumn(col)
	})
}

// Write implements [Change].
func (m Changes) Write(builder *StatementBuilder) {
	for i, change := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		change.Write(builder)
	}
}

// Matches implements [gomock.Matcher].
func (c Changes) Matches(x any) bool {
	toMatch, ok := x.(*Changes)
	if !ok {
		return false
	}
	if len(c) != len(*toMatch) {
		return false
	}
	for i := range c {
		if !c[i].Matches((*toMatch)[i]) {
			return false
		}
	}
	return true
}

// String implements [gomock.Matcher].
func (c Changes) String() string {
	return "database.Changes"
}

var _ Change = Changes(nil)

func NewChangeToNull(col Column) Change {
	return NewChange(col, NullInstruction)
}

func NewChangeToColumn(to, from Column) Change {
	return &changeToColumn{to: to, from: from}
}

type changeToColumn struct {
	to   Column
	from Column
}

// IsOnColumn implements [Change].
func (c *changeToColumn) IsOnColumn(col Column) bool {
	return c.to.Equals(col)
}

// Matches implements [Change].
func (c *changeToColumn) Matches(x any) bool {
	toMatch, ok := x.(*changeToColumn)
	if !ok {
		return false
	}
	return c.to.Equals(toMatch.to) && c.from.Equals(toMatch.from)
}

// String implements [Change].
func (c *changeToColumn) String() string {
	return "database.changeToColumn"
}

// Write implements [Change].
func (c *changeToColumn) Write(builder *StatementBuilder) {
	c.to.WriteUnqualified(builder)
	builder.WriteString(" = ")
	c.from.WriteQualified(builder)
}

var _ Change = (*changeToColumn)(nil)

type incrementColumnChange struct {
	column Column
}

func NewIncrementColumnChange(col Column) Change {
	return &incrementColumnChange{
		column: col,
	}
}

// IsOnColumn implements [Change].
func (i *incrementColumnChange) IsOnColumn(col Column) bool {
	return i.column.Equals(col)
}

// Matches implements [Change].
func (i *incrementColumnChange) Matches(x any) bool {
	toMatch, ok := x.(*incrementColumnChange)
	if !ok {
		return false
	}
	return i.column.Equals(toMatch.column)
}

// String implements [Change].
func (i *incrementColumnChange) String() string {
	return "database.incrementColumnChange"
}

// Write implements [Change].
func (i *incrementColumnChange) Write(builder *StatementBuilder) {
	i.column.WriteUnqualified(builder)
	builder.WriteString(" = ")
	i.column.WriteUnqualified(builder)
	builder.WriteString(" + 1")
}

var _ Change = (*incrementColumnChange)(nil)

func NewChangeToStatement(col Column, stmt func(builder *StatementBuilder)) Change {
	return &changeToStatement{
		column: col,
		stmt:   stmt,
	}
}

type changeToStatement struct {
	column Column
	stmt   func(builder *StatementBuilder)
}

// IsOnColumn implements [Change].
func (c *changeToStatement) IsOnColumn(col Column) bool {
	return c.column.Equals(col)
}

// Matches implements [Change].
func (c *changeToStatement) Matches(x any) bool {
	toMatch, ok := x.(*changeToStatement)
	if !ok {
		return false
	}
	var expectedBuilder, actualBuilder StatementBuilder
	c.stmt(&expectedBuilder)
	toMatch.stmt(&actualBuilder)

	if expectedBuilder.String() != actualBuilder.String() {
		return false
	}
	if !slices.Equal(expectedBuilder.Args(), actualBuilder.Args()) {
		return false
	}
	return c.column.Equals(toMatch.column)
}

// String implements [Change].
func (c *changeToStatement) String() string {
	return "database.changeToStatement"
}

// Write implements [Change].
func (c *changeToStatement) Write(builder *StatementBuilder) {
	_, ok := c.column.(Columns)
	if ok {
		builder.WriteRune('(')
	}
	c.column.WriteUnqualified(builder)
	if ok {
		builder.WriteRune(')')
	}
	builder.WriteString(" = (")
	c.stmt(builder)
	builder.WriteString(")")
}

var _ Change = (*changeToStatement)(nil)

// CTEChange represents a change that uses a Common Table Expression (CTE).
// It intercepts the Write process to first write the CTE part, and then the main change part.
type CTEChange interface {
	Change
	// WriteCTE writes the CTE part of the change to the given statement builder.
	// It writes the part inside the brackets meaning without the "WITH cte_name AS (" and the ending ")".
	WriteCTE(builder *StatementBuilder)
	// SetName sets the name of the CTE.
	// This is defined by the caller to ensure uniqueness.
	// The name is used to reference the CTE in the main change.
	SetName(name string)
}

func NewCTEChange(cte func(builder *StatementBuilder), change func(name string) Change) CTEChange {
	return &cteChange{
		cte:    cte,
		change: change,
	}
}

type cteChange struct {
	name   string
	cte    func(builder *StatementBuilder)
	change func(name string) Change
}

// IsOnColumn implements [CTEChange].
func (c *cteChange) IsOnColumn(col Column) bool {
	return false
}

// Matches implements [CTEChange].
func (c *cteChange) Matches(x any) bool {
	toMatch, ok := x.(*cteChange)
	if !ok {
		return false
	}
	var expectedCTEBuilder, actualCTEBuilder StatementBuilder
	c.cte(&expectedCTEBuilder)
	if c.change != nil {
		c.change(c.name).Write(&expectedCTEBuilder)
	}
	toMatch.cte(&actualCTEBuilder)
	if toMatch.change != nil {
		toMatch.change(toMatch.name).Write(&actualCTEBuilder)
	}

	if expectedCTEBuilder.String() != actualCTEBuilder.String() {
		return false
	}
	return slices.Equal(expectedCTEBuilder.Args(), actualCTEBuilder.Args())
}

// Name implements [CTEChange].
func (c *cteChange) SetName(name string) {
	c.name = name
}

// String implements [CTEChange].
func (c *cteChange) String() string {
	return "database.cteChange"
}

// Write implements [CTEChange].
func (c *cteChange) Write(builder *StatementBuilder) {
	if c.change == nil {
		return
	}
	c.change(c.name).Write(builder)
}

// WriteCTE implements [CTEChange].
func (c *cteChange) WriteCTE(builder *StatementBuilder) {
	c.cte(builder)
}
