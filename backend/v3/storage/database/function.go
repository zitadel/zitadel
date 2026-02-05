package database

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
func (c coalesce) Write(builder *StatementBuilder) error {
	builder.WriteString("COALESCE(")
	for i, value := range c.values {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteArg(value)
	}
	builder.WriteString(")")
	return nil
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
