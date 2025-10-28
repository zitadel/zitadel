package json

import (
	"encoding/json"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ database.Change = (*jsonChanges)(nil)

type JsonUpdate interface {
	writeUpdate(builder *database.StatementBuilder, changes jsonChanges, i int) error
	getPathValue() (string, string, error)
}

type FieldChange struct {
	path  string
	value any
}

func NewFieldChange(path string, value any) JsonUpdate {
	return &FieldChange{
		path:  path,
		value: value,
	}
}

func (c *FieldChange) getPathValue() (string, string, error) {
	path, v := c.path, c.value

	value, err := json.Marshal(v)
	if err != nil {
		return "", "", err
	}

	// var out string
	// if valueString, ok := v.(string); ok {
	// 	out = "\"" + string(valueString) + "\""
	// } else {
	// 	out = string(value)
	// }

	// return path, out, nil

	return path, string(value), nil
}

func (c *FieldChange) writeUpdate(builder *database.StatementBuilder, changes jsonChanges, i int) error {
	path, value, err := c.getPathValue()
	if err != nil {
		return err
	}

	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> path = %+v\n", path)
	builder.WriteString("jsonb_set_lax(")
	// if i == 0 {
	if i < 0 {
		changes.column.WriteQualified(builder)
	} else {
		changes.changes[i].writeUpdate(builder, changes, i-1)
	}
	builder.WriteString(", " + path)
	if value == "null" {
		builder.WriteString(", " + string(value))
	} else {
		builder.WriteString(", '" + string(value) + "'")
	}
	builder.WriteString(", " + "true")
	builder.WriteString(", 'delete_key'")

	builder.WriteString(")")

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ArrayChange struct {
	path   string
	value  any
	remove bool
}

func NewArrayChange(path string, value any, remove bool) JsonUpdate {
	return &ArrayChange{
		path:   path,
		value:  value,
		remove: true,
	}
}

func (c *ArrayChange) getPathValue() (string, string, error) {
	path, v := c.path, c.value

	value, err := json.Marshal(v)
	if err != nil {
		return "", "", err
	}

	return path, string(value), nil
}

func (c *ArrayChange) writeUpdate(builder *database.StatementBuilder, changes jsonChanges, i int) error {
	path, value, err := c.getPathValue()
	if err != nil {
		return err
	}

	// jsonb_set(properties, '{attributes}', (properties->'attributes') - 'is_new');
	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> path = %+v\n", path)

	if c.remove {
		builder.WriteString("jsonb_set(")
		//
		if i < 0 {
			changes.column.WriteQualified(builder)
		} else {
			// changes.changes[i-1].writeUpdate(builder, changes, i-1)
			c.writeUpdate(builder, changes, i-1)
		}
		//
		builder.WriteString(", '{" + path + "}'")
		//
		builder.WriteString(", " + "(")
		changes.column.WriteUnqualified(builder)
		builder.WriteString("->'" + path + "')")
		builder.WriteString(" - ")
		if len(value) > 2 && value[0] == '"' && value[len(value)-1] == '"' {
			builder.WriteString("'" + value[1:len(value)-1] + "'")
		} else {
			builder.WriteString(value)
		}
		builder.WriteString(")")
	} else {

		builder.WriteString("jsonb_insert(")
		//
		if i == 0 {
			changes.column.WriteUnqualified(builder)
		} else {
			changes.changes[i-1].writeUpdate(builder, changes, i-1)
		}
		//
		builder.WriteString(", " + "(CASE WHEN (SELECT ")
		changes.column.WriteUnqualified(builder)
		builder.WriteString(" ? '" + path + "') THEN")
		builder.WriteString(" '{" + path + ", -1}'")
		builder.WriteString(" ELSE '{" + path + "}'")
		builder.WriteString(" END)::TEXT[]")
		//
		builder.WriteString(", " + "(CASE WHEN (SELECT ")
		changes.column.WriteUnqualified(builder)
		builder.WriteString(" ? '" + path + "') THEN")
		builder.WriteString(" '" + value + "'::JSONB")
		builder.WriteString(" ELSE jsonb_build_array(")
		if len(value) > 2 && value[0] == '"' && value[len(value)-1] == '"' {
			builder.WriteString("'" + value[1:len(value)-1] + "'")
		} else {
			builder.WriteString(value)
		}
		builder.WriteString(")")
		builder.WriteString(" END)::JSONB")
		//
		builder.WriteString(", " + "true")

		builder.WriteString(")")
	}

	return nil
}

var _ JsonUpdate = (*ArrayChange)(nil)

// func JSONArrayRemove(path string, value any) JSONFieldChange {
// 			path:  path,
// 		value: value,
// 	}
// }

type jsonChanges struct {
	column  database.Column
	changes []JsonUpdate
}

var _ database.Change = (*jsonChanges)(nil)

func NewJsonChanges(col database.Column, changes ...JsonUpdate) database.Change {
	return &jsonChanges{
		column:  col,
		changes: changes,
	}
}

// func (c change) Write(builder *database.StatementBuilder) error {
// 	return c.WriteUpdate(builder)
// }

func (c jsonChanges) Write(builder *database.StatementBuilder) error {
	c.column.WriteUnqualified(builder)
	builder.WriteString(" = ")

	if c.changes == nil {
		return nil
	}
	// return c.writeUpdate(builder, len(c.changes)-1)
	// return c.changes[len(c.changes)-1].writeUpdate(builder, c, len(c.changes)-2)
	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> len(c.changes) = %+v\n", len(c.changes))
	return c.changes[len(c.changes)-1].writeUpdate(builder, c, len(c.changes)-1)
}

// func (c jsonChanges) writeUpdate(builder *database.StatementBuilder, i int) error {
// 	path, v := c.changes[i].path, c.changes[i].value

// 	value, err := json.Marshal(v)
// 	if err != nil {
// 		return err
// 	}
// 	builder.WriteString("jsonb_set_lax(")
// 	if i == 0 {
// 		c.column.WriteUnqualified(builder)
// 	} else {
// 		c.writeUpdate(builder, i-1)
// 	}
// 	builder.WriteString(", " + path)
// 	if value == nil {
// 		builder.WriteString(", " + string(value))
// 	} else {
// 		builder.WriteString(", '" + string(value) + "'")
// 	}
// 	builder.WriteString(", " + "true")
// 	builder.WriteString(", 'delete_key'")

// 	builder.WriteString(")")

// 	return nil
// }

func (c jsonChanges) IsOnColumn(col database.Column) bool {
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
