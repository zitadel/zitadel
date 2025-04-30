package v4

type queryOpts struct {
	condition Condition
	orderBy   Columns
	limit     uint32
	offset    uint32
}

func (opts *queryOpts) writeCondition(builder *statementBuilder) {
	if opts.condition == nil {
		return
	}
	builder.WriteString(" WHERE ")
	opts.condition.writeTo(builder)
}

func (opts *queryOpts) writeOrderBy(builder *statementBuilder) {
	if len(opts.orderBy) == 0 {
		return
	}
	builder.WriteString(" ORDER BY ")
	opts.orderBy.writeTo(builder)
}

func (opts *queryOpts) writeLimit(builder *statementBuilder) {
	if opts.limit == 0 {
		return
	}
	builder.WriteString(" LIMIT ")
	builder.writeArg(opts.limit)
}

func (opts *queryOpts) writeOffset(builder *statementBuilder) {
	if opts.offset == 0 {
		return
	}
	builder.WriteString(" OFFSET ")
	builder.writeArg(opts.offset)
}

type QueryOption func(*queryOpts)

func WithCondition(condition Condition) QueryOption {
	return func(opts *queryOpts) {
		opts.condition = condition
	}
}

func WithOrderBy(orderBy ...Column) QueryOption {
	return func(opts *queryOpts) {
		opts.orderBy = orderBy
	}
}

func WithLimit(limit uint32) QueryOption {
	return func(opts *queryOpts) {
		opts.limit = limit
	}
}

func WithOffset(offset uint32) QueryOption {
	return func(opts *queryOpts) {
		opts.offset = offset
	}
}
