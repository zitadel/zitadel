package database

type QueryOption func(opts *QueryOpts)

// WithCondition sets the condition for the query.
func WithCondition(condition Condition) QueryOption {
	return func(opts *QueryOpts) {
		opts.Condition = condition
	}
}

// WithOrderBy sets the columns to order the results by.
func WithOrderBy(ordering OrderDirection, orderBy ...Column) QueryOption {
	return func(opts *QueryOpts) {
		opts.OrderBy = orderBy
		opts.Ordering = ordering
	}
}

func WithOrderByAscending(columns ...Column) QueryOption {
	return WithOrderBy(OrderDirectionAsc, columns...)
}

func WithOrderByDescending(columns ...Column) QueryOption {
	return WithOrderBy(OrderDirectionDesc, columns...)
}

// WithLimit sets the maximum number of results to return.
func WithLimit(limit uint32) QueryOption {
	return func(opts *QueryOpts) {
		opts.Limit = limit
	}
}

// WithOffset sets the number of results to skip before returning the results.
func WithOffset(offset uint32) QueryOption {
	return func(opts *QueryOpts) {
		opts.Offset = offset
	}
}

// WithGroupBy sets the columns to group the results by.
func WithGroupBy(groupBy ...Column) QueryOption {
	return func(opts *QueryOpts) {
		opts.GroupBy = groupBy
	}
}

// WithLeftJoin adds a LEFT JOIN to the query.
func WithLeftJoin(table string, columns Condition) QueryOption {
	return func(opts *QueryOpts) {
		opts.Joins = append(opts.Joins, join{
			table:   table,
			typ:     JoinTypeLeft,
			columns: columns,
		})
	}
}

func WithPermissionCheck(permission string) QueryOption {
	return func(opts *QueryOpts) {
		opts.Permission = permission
	}
}

type joinType string

const (
	JoinTypeLeft joinType = "LEFT"
)

type join struct {
	table   string
	typ     joinType
	columns Condition
}

type OrderDirection uint8

const (
	OrderDirectionAsc OrderDirection = iota
	OrderDirectionDesc
)

// QueryOpts holds the options for a query.
// It is used to build the SQL SELECT statement.
type QueryOpts struct {
	// Condition is the condition to filter the results.
	// It is used to build the WHERE clause of the SQL statement.
	Condition Condition
	// OrderBy is the columns to order the results by.
	// It is used to build the ORDER BY clause of the SQL statement.
	OrderBy Columns
	// Ordering defines if the columns should be ordered ascending or descending.
	// Default is ascending.
	Ordering OrderDirection
	// Limit is the maximum number of results to return.
	// It is used to build the LIMIT clause of the SQL statement.
	Limit uint32
	// Offset is the number of results to skip before returning the results.
	// It is used to build the OFFSET clause of the SQL statement.
	Offset uint32
	// GroupBy is the columns to group the results by.
	// It is used to build the GROUP BY clause of the SQL statement.
	GroupBy Columns
	// Joins is a list of joins to be applied to the query.
	// It is used to build the JOIN clauses of the SQL statement.
	Joins []join
	// Permission required to read or write the resource.
	// When unset, no permission check is made.
	Permission string
}

func (opts *QueryOpts) Write(builder *StatementBuilder) {
	opts.WriteLeftJoins(builder)
	opts.WriteCondition(builder)
	opts.WriteGroupBy(builder)
	opts.WriteOrderBy(builder)
	opts.WriteLimit(builder)
	opts.WriteOffset(builder)
}

func (opts *QueryOpts) WriteCondition(builder *StatementBuilder) {
	if opts.Condition == nil {
		return
	}
	builder.WriteString(" WHERE ")
	opts.Condition.Write(builder)
}

func (opts *QueryOpts) WriteOrderBy(builder *StatementBuilder) {
	if len(opts.OrderBy) == 0 {
		return
	}
	builder.WriteString(" ORDER BY ")
	for i, col := range opts.OrderBy {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.WriteQualified(builder)
		if opts.Ordering == OrderDirectionDesc {
			builder.WriteString(" DESC")
		}
	}
}

func (opts *QueryOpts) WriteLimit(builder *StatementBuilder) {
	if opts.Limit == 0 {
		return
	}
	builder.WriteString(" LIMIT ")
	builder.WriteArg(opts.Limit)
}

func (opts *QueryOpts) WriteOffset(builder *StatementBuilder) {
	if opts.Offset == 0 {
		return
	}
	builder.WriteString(" OFFSET ")
	builder.WriteArg(opts.Offset)
}

func (opts *QueryOpts) WriteGroupBy(builder *StatementBuilder) {
	if len(opts.GroupBy) == 0 {
		return
	}
	builder.WriteString(" GROUP BY ")
	opts.GroupBy.WriteQualified(builder)
}

func (opts *QueryOpts) WriteLeftJoins(builder *StatementBuilder) {
	if len(opts.Joins) == 0 {
		return
	}
	for _, join := range opts.Joins {
		builder.WriteString(" ")
		builder.WriteString(string(join.typ))
		builder.WriteString(" JOIN ")
		builder.WriteString(join.table)
		builder.WriteString(" ON ")
		join.columns.Write(builder)
	}
}
