package database

type QueryOption func(opts *QueryOpts)

func WithCondition(condition Condition) QueryOption {
	return func(opts *QueryOpts) {
		opts.Condition = condition
	}
}

func WithOrderBy(orderBy ...Column) QueryOption {
	return func(opts *QueryOpts) {
		opts.OrderBy = orderBy
	}
}

func WithLimit(limit uint32) QueryOption {
	return func(opts *QueryOpts) {
		opts.Limit = limit
	}
}

func WithOffset(offset uint32) QueryOption {
	return func(opts *QueryOpts) {
		opts.Offset = offset
	}
}

type QueryOpts struct {
	Condition Condition
	OrderBy   Columns
	Limit     uint32
	Offset    uint32
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
