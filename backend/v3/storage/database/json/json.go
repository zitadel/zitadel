package json

import (
	"encoding/json"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ database.Change = (*jsonChanges)(nil)

type JsonUpdate interface {
	writeUpdate(builder *database.StatementBuilder, changes jsonChanges, i int) error
	getPathValue() ([]string, string, error)
}

type FieldChange struct {
	path  []string
	value any
}

func NewFieldChange(path []string, value any) JsonUpdate {
	return &FieldChange{
		path:  path,
		value: value,
	}
}

func (c *FieldChange) getPathValue() ([]string, string, error) {
	path, v := c.path, c.value

	value, err := json.Marshal(v)
	if err != nil {
		return nil, "", err
	}

	return path, string(value), nil
}

func writePath(path []string) string {
	if len(path) == 0 {
		return "{}"
	}

	out := "{"

	for _, p := range path[:len(path)-1] {
		// out += "'" + p + "',"
		out += p
	}
	// out += "'" + path[len(path)-1] + "'"
	out += path[len(path)-1]

	out += "}"

	return out
}

func (c *FieldChange) writeUpdate(builder *database.StatementBuilder, changes jsonChanges, i int) error {
	path, value, err := c.getPathValue()
	if err != nil {
		return err
	}
	fmt.Println("\033[45m[DBUGPRINT]\033[0m[:1]\033[45m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m inside writeUpdate()")

	builder.WriteString("jsonb_set_lax(")
	if i < 0 {
		changes.column.WriteQualified(builder)
	} else {
		changes.changes[i].writeUpdate(builder, changes, i-1)
	}
	builder.WriteString(", '" + writePath(path))
	if value == "null" {
		builder.WriteString("', " + string(value))
	} else {
		builder.WriteString("', '" + string(value) + "'")
	}
	builder.WriteString(", " + "true")
	builder.WriteString(", 'delete_key'")

	builder.WriteString(")")

	return nil
}

type jsonArrayValues interface {
	~int64 | ~string
}

type ArrayChange struct {
	path   []string
	value  any
	remove bool
}

func NewArrayChange(path []string, value any, remove bool) []JsonUpdate {
	// if remove {
	return []JsonUpdate{
		&ArrayChange{
			path:   path,
			value:  value,
			remove: remove,
		},
	}
	// } else {
	// 	return []JsonUpdate{
	// 		// first remove so we don't have duplicates in the array
	// 		&ArrayChange{
	// 			path:  path,
	// 			value: value,
	// 			// remove: !remove,
	// 			remove: true,
	// 		},
	// 		&ArrayChange{
	// 			path:   path,
	// 			value:  value,
	// 			remove: remove,
	// 		},
	// 	}
	// }
}

func (c *ArrayChange) getPathValue() ([]string, string, error) {
	path, v := c.path, c.value

	value, err := json.Marshal(v)
	if err != nil {
		return nil, "", err
	}

	return path, string(value), nil
}

func (c *ArrayChange) writeUpdate(builder *database.StatementBuilder, changes jsonChanges, i int) error {
	fmt.Printf("\033[43m[DBUGPRINT]\033[0m[:1]\033[43m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m c.remove = %p\n", &c.remove)
	if c.remove {
		return c.removeFromArray(builder, changes, i)
	} else {
		return c.addToArray(builder, changes, i)
	}
}

func (c *ArrayChange) addToArray(builder *database.StatementBuilder, changes jsonChanges, i int) error {
	fmt.Println("\033[45m[DBUGPRINT]\033[0m[:1]\033[45m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m ADD")
	fmt.Printf("\033[43m[DBUGPRINT]\033[0m[:1]\033[43m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m i = %+v\n", i)
	path, value, err := c.getPathValue()
	if err != nil {
		return err
	}
	fmt.Printf("\033[43m[DBUGPRINT]\033[0m[:1]\033[43m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m value = %+v\n", value)

	builder.WriteString("zitadel.jsonb_array_append(")
	if i < 0 {
		changes.column.WriteQualified(builder)
	} else {
		changes.changes[i].writeUpdate(builder, changes, i-1)
	}
	builder.WriteString(", ")
	builder.WriteArg(path)
	builder.WriteString(", ")
	builder.WriteArg(value)
	if len(value) > 2 && value[0] == '"' && value[len(value)-1] == '"' {
		builder.WriteString("::TEXT")
	} else {
		builder.WriteString("::NUMERIC")
	}

	builder.WriteString(")")

	return nil
}

func (c *ArrayChange) removeFromArray(builder *database.StatementBuilder, changes jsonChanges, i int) error {
	fmt.Println("\033[45m[DBUGPRINT]\033[0m[:1]\033[45m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m REMOVE")
	fmt.Printf("\033[43m[DBUGPRINT]\033[0m[:1]\033[43m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m i = %+v\n", i)
	path, value, err := c.getPathValue()
	if err != nil {
		return err
	}

	builder.WriteString("zitadel.jsonb_array_remove(")
	if i < 0 {
		changes.column.WriteQualified(builder)
	} else {
		changes.changes[i].writeUpdate(builder, changes, i-1)
	}
	builder.WriteString(", ")
	builder.WriteArg(path)
	builder.WriteString(", ")
	builder.WriteArg(value)
	builder.WriteString("::TEXT")
	builder.WriteString(")")

	return nil
}

var _ JsonUpdate = (*ArrayChange)(nil)

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

func (c *jsonChanges) Matches(x any) bool {
	return false
}

func (c *jsonChanges) String() string {
	return ""
}

func (c jsonChanges) Write(builder *database.StatementBuilder) error {
	c.column.WriteUnqualified(builder)
	builder.WriteString(" = ")

	if c.changes == nil {
		return nil
	}
	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> len(c.changes) = %+v\n", len(c.changes))
	return c.changes[len(c.changes)-1].writeUpdate(builder, c, len(c.changes)-2)
}

func (c jsonChanges) IsOnColumn(col database.Column) bool {
	return c.column.Equals(col)
}
