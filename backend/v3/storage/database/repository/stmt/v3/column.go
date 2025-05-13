package v3

type Column interface {
	Name() string
	Write(builder statementBuilder)
}

type ignoreCaseColumn interface {
	Column
	WriteIgnoreCase(builder statementBuilder)
}

var (
	columnNameID        = "id"
	columnNameName      = "name"
	columnNameCreatedAt = "created_at"
	columnNameUpdatedAt = "updated_at"
	columnNameDeletedAt = "deleted_at"

	columnNameInstanceID = "instance_id"

	columnNameOrgID = "org_id"
)

type column struct {
	table Table
	name  string
}

// Write implements Column.
func (c *column) Write(builder statementBuilder) {
	c.table.writeOn(builder)
	builder.writeRune('.')
	builder.writeString(c.name)
}

// Name implements [Column].
func (c *column) Name() string {
	return c.name
}

var _ Column = (*column)(nil)

type columnIgnoreCase struct {
	column
	suffix string
}

// WriteIgnoreCase implements ignoreCaseColumn.
func (c *columnIgnoreCase) WriteIgnoreCase(builder statementBuilder) {
	c.Write(builder)
	builder.writeString(c.suffix)
}

var _ ignoreCaseColumn = (*columnIgnoreCase)(nil)
