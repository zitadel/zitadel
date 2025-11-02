package json

import (
	"encoding/json"
	"fmt"
	"strings"

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
	return nil
}

func (c *ArrayChange) removeFromArray(builder *database.StatementBuilder, changes jsonChanges, i int) error {
	fmt.Println("\033[45m[DBUGPRINT]\033[0m[:1]\033[45m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m REMOVE")
	fmt.Printf("\033[43m[DBUGPRINT]\033[0m[:1]\033[43m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m i = %+v\n", i)
	path, value, err := c.getPathValue()
	if err != nil {
		return err
	}

	// outt = jsonb_set(source, path,
	//  (CASE WHEN (SELECT source ?& path) THEN
	//    COALESCE(
	//    (SELECT jsonb_agg(v)
	//    FROM jsonb_array_elements(source #>path) AS elem(v)
	//    WHERE v::TEXT <> value::TEXT),
	//    jsonb_build_array())
	//  ELSE
	//    jsonb_build_array()
	//  END)::JSONB, true);

	// builder.WriteString(">>>>> " + strconv.Itoa(i) + " Start <<<<<")
	builder.WriteString("jsonb_set(")
	//
	if i < 0 {
		changes.column.WriteQualified(builder)
	} else {
		changes.changes[i].writeUpdate(builder, changes, i-1)
	}
	//
	builder.WriteString(", '" + writePath(path))
	//
	builder.WriteString("', " + "(CASE WHEN (SELECT ")
	changes.column.WriteUnqualified(builder)
	builder.WriteString(" ?& '" + writePath(path) + "') THEN")
	builder.WriteString(" COALESCE(")
	builder.WriteString(" (SELECT jsonb_agg(v)")
	builder.WriteString(" FROM jsonb_array_elements(")
	changes.column.WriteUnqualified(builder)
	builder.WriteString(" #>")
	builder.WriteString(" '" + writePath(path) + "') AS elem(v)")
	builder.WriteString(" WHERE v::TEXT <> " + value + "::TEXT)")
	builder.WriteString(", jsonb_build_array()) ")
	builder.WriteString(" ELSE jsonb_build_array() ")
	builder.WriteString(" END)::JSONB, true")
	builder.WriteString(")")
	// builder.WriteString(">>>>> " + strconv.Itoa(i) + " End <<<<<")

	return nil
}

func (c *ArrayChange) addToArray(builder *database.StatementBuilder, changes jsonChanges, i int) error {
	fmt.Println("\033[45m[DBUGPRINT]\033[0m[:1]\033[45m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m ADD")
	fmt.Printf("\033[43m[DBUGPRINT]\033[0m[:1]\033[43m>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\033[0m i = %+v\n", i)
	path, value, err := c.getPathValue()
	if err != nil {
		return err
	}

	// 	SET properties = jsonb_insert(
	//   properties,
	//   (CASE WHEN (SELECT properties ? 'attributes') THEN
	//   '{attributes, -1}'
	//   ELSE
	//   '{attributes}'
	//   END)::TEXT[],
	//   (CASE WHEN (SELECT properties ? 'attributes') THEN
	//   '"3sdfdf8"'::JSONB
	//   ELSE
	//   jsonb_build_array('value')
	//   END),
	//   true
	// )

	// builder.WriteString(">>>>> " + strconv.Itoa(i) + " Start <<<<<")
	builder.WriteString("jsonb_insert(")
	//
	if i < 0 {
		changes.column.WriteQualified(builder)
	} else {
		changes.changes[i].writeUpdate(builder, changes, i-1)
	}
	//
	builder.WriteString(", " + "(CASE WHEN (SELECT ")
	changes.column.WriteQualified(builder)
	builder.WriteString(" ?& '" + writePath(path) + "') THEN")
	builder.WriteString(" '{" + strings.Join(path, "'") + ", -1}'")
	builder.WriteString(" ELSE '{" + strings.Join(path, ",") + "}'")
	builder.WriteString(" END)::TEXT[]")
	//
	builder.WriteString(", " + "(CASE WHEN (SELECT ")
	changes.column.WriteQualified(builder)
	builder.WriteString(" ?& '" + writePath(path) + "') THEN")
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
	// builder.WriteString(">>>>> " + strconv.Itoa(i) + " End <<<<<")

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

func (c jsonChanges) Write(builder *database.StatementBuilder) error {
	c.column.WriteUnqualified(builder)
	builder.WriteString(" = ")

	if c.changes == nil {
		return nil
	}
	// return c.writeUpdate(builder, len(c.changes)-1)
	// return c.changes[len(c.changes)-1].writeUpdate(builder, c, len(c.changes)-2)
	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> len(c.changes) = %+v\n", len(c.changes))
	return c.changes[len(c.changes)-1].writeUpdate(builder, c, len(c.changes)-2)
}

func (c jsonChanges) IsOnColumn(col database.Column) bool {
	return c.column.Equals(col)
}
