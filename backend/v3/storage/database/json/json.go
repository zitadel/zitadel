package json

import (
	"encoding/json"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ database.Change = (*change)(nil)

type JSONFieldChange struct {
	path  string
	value any
}

func NewChange(path string, value any) JSONFieldChange {
	return JSONFieldChange{
		path:  path,
		value: value,
	}
}

type change struct {
	column  database.Column
	changes []JSONFieldChange
}

var _ database.Change = (*change)(nil)

func NewJsonChange(col database.Column, changes ...JSONFieldChange) database.Change {
	return &change{
		column:  col,
		changes: changes,
	}
}

// func (c change) Write(builder *database.StatementBuilder) error {
// 	return c.WriteUpdate(builder)
// }

func (c change) Write(builder *database.StatementBuilder) error {
	c.column.WriteUnqualified(builder)
	builder.WriteString(" = ")

	return c.writeUpdate(builder, len(c.changes)-1)
}

func (c change) writeUpdate(builder *database.StatementBuilder, i int) error {
	k, v := c.changes[i].path, c.changes[i].value

	value, err := json.Marshal(v)
	if err != nil {
		return err
	}
	builder.WriteString("jsonb_set_lax(")
	if i == 0 {
		c.column.WriteUnqualified(builder)
	} else {
		c.writeUpdate(builder, i-1)
	}
	builder.WriteString(", " + k)
	if value == nil {
		builder.WriteString(", " + strings.ToUpper(string(value)))
	} else {
		// builder.WriteString(", '" + string(value) + "'")
		builder.WriteString(", '" + string(value) + "'")
	}
	builder.WriteString(", " + "true")
	builder.WriteString(", 'delete_key'")

	builder.WriteString(")")

	return nil
}

// IsOnColumn implements [JSONFieldChange].
func (c change) IsOnColumn(col database.Column) bool {
	return c.column.Equals(col)
}

// type Changes []Change

// func NewChanges(cols ...Change) Change {
// 	return Changes(cols)
// }

// // IsOnColumn implements [Change].
// func (c Changes) IsOnColumn(col Column) bool {
// 	return slices.ContainsFunc(c, func(change Change) bool {
// 		return change.IsOnColumn(col)
// 	})
// }

// // Write implements [Change].
// func (m Changes) Write(builder *StatementBuilder) {
// 	for i, change := range m {
// 		if i > 0 {
// 			builder.WriteString(", ")
// 		}
// 		change.Write(builder)
// 	}
// }

// var _ Change = Changes(nil)
