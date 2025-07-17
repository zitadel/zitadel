package database

type QueryOption func(opts *QueryOpts)

// WithCondition sets the condition for the query.
func WithCondition(condition Condition) QueryOption {
	return func(opts *QueryOpts) {
		opts.Condition = condition
	}
}

// WithOrderBy sets the columns to order the results by.
func WithOrderBy(orderBy ...Column) QueryOption {
	return func(opts *QueryOpts) {
		opts.OrderBy = orderBy
	}
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

// QueryOpts holds the options for a query.
// It is used to build the SQL SELECT statement.
type QueryOpts struct {
	// Condition is the condition to filter the results.
	// It is used to build the WHERE clause of the SQL statement.
	Condition Condition
	// OrderBy is the columns to order the results by.
	// It is used to build the ORDER BY clause of the SQL statement.
	OrderBy   Columns
	// Limit is the maximum number of results to return.
	// It is used to build the LIMIT clause of the SQL statement.
	Limit     uint32
	// Offset is the number of results to skip before returning the results.
	// It is used to build the OFFSET clause of the SQL statement.
	Offset    uint32
}

func (opts *QueryOpts) Write(builder *StatementBuilder) {
	opts.WriteCondition(builder)
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
	opts.OrderBy.Write(builder)
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
