package database

type Columns []Column

// WriteQualified implements [Column].
// Columns are separated by ", ".
func (m Columns) WriteQualified(builder *StatementBuilder) {
	for i, col := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.WriteQualified(builder)
	}
}

// WriteUnqualified implements [Column].
// Columns are separated by ", ".
func (m Columns) WriteUnqualified(builder *StatementBuilder) {
	for i, col := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.WriteUnqualified(builder)
	}
}

// Column represents a column in a database table.
type Column interface {
	// WriteQualified writes the column with the table name as prefix.
	WriteQualified(builder *StatementBuilder)
	// WriteUnqualified writes the column without the table name as prefix.
	WriteUnqualified(builder *StatementBuilder)
	// Equals checks if two columns are equal.
	Equals(col Column) bool
}

type column struct {
	table string
	name  string
}

func NewColumn(table, name string) Column {
	return &column{table: table, name: name}
}

// WriteQualified implements [Column].
func (c column) WriteQualified(builder *StatementBuilder) {
	builder.Grow(len(c.table) + len(c.name) + 1)
	builder.WriteString(c.table)
	builder.WriteRune('.')
	builder.WriteString(c.name)
}

// WriteUnqualified implements [Column].
func (c column) WriteUnqualified(builder *StatementBuilder) {
	builder.WriteString(c.name)
}

func (c *column) Equals(col Column) bool {
	if col == nil {
		return c == nil
	}
	toMatch, ok := col.(*column)
	if !ok {
		return false
	}
	return c.table == toMatch.table && c.name == toMatch.name
}

// LowerColumn returns a column that represents LOWER(col).
func LowerColumn(col Column) Column {
	return &functionColumn{fn: functionLower, col: col}
}

// SHA256Column returns a column that represents SHA256(col).
func SHA256Column(col Column) Column {
	return &functionColumn{fn: functionSHA256, col: col}
}

type functionColumn struct {
	fn  function
	col Column
}

type function string

const (
	_              function = ""
	functionLower  function = "LOWER"
	functionSHA256 function = "SHA256"
)

func (c functionColumn) WriteQualified(builder *StatementBuilder) {
	builder.Grow(len(c.fn) + 2)
	builder.WriteString(string(c.fn))
	builder.WriteRune('(')
	c.col.WriteQualified(builder)
	builder.WriteRune(')')
}

func (c functionColumn) WriteUnqualified(builder *StatementBuilder) {
	builder.Grow(len(c.fn) + 2)
	builder.WriteString(string(c.fn))
	builder.WriteRune('(')
	c.col.WriteUnqualified(builder)
	builder.WriteRune(')')
}

func (c *functionColumn) Equals(col Column) bool {
	if col == nil {
		return c == nil
	}
	toMatch, ok := col.(*functionColumn)
	if !ok || toMatch.fn != c.fn {
		return false
	}
	return c.col.Equals(toMatch.col)
}
