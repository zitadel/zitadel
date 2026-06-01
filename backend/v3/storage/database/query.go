package database

import (
	"fmt"
	"reflect"
	"slices"

	"go.uber.org/mock/gomock"
)

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

// WithPermissionCheck enables a check if the authenticated user has the
// passed permission to read or write a resource.
func WithPermissionCheck(permission string) QueryOption {
	return func(opts *QueryOpts) {
		opts.Permission = permission
	}
}

// WithResultLock locks the results of the query during the transaction.
func WithResultLock() QueryOption {
	return func(opts *QueryOpts) {
		opts.ShouldLock = true
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
	// ShouldLock the results during the transaction.
	ShouldLock bool
}

// Matches implements [gomock.Matcher].
func (q *QueryOpts) Matches(x any) bool {
	// first check if the x is a [QueryOpt]
	inputOpts, ok := x.(*QueryOpts)
	if !ok {
		// second possibility is a [QueryOption]
		optFunc, ok := x.(QueryOption)
		if !ok {
			return false
		}

		// QueryOption is a function that takes a *QueryOpts in input and fills it with its data.
		// QueryOption data is not accessible because it's a function, so we exploit this "hack"
		// to read the data.
		inputOpts = new(QueryOpts)
		optFunc(inputOpts)
	}

	// inputOpts now contains the actual data but made of other interfaces and functions/decorator.
	// Doing a reflect.DeepEqual() will fail because you will end up into comparing functions
	// which is not possible (comparison is successful only if both functions are nil).
	// So we exploit the Write() method of QueryOpts to fill up the StatementBuilders:
	// these objects are basically string builders, so we can leverage their String()
	// method (implementing Stringer interface) to make an easy and safe comparison.
	inputBuilder, expectedBuilder := &StatementBuilder{}, &StatementBuilder{}
	inputOpts.Write(inputBuilder)
	q.Write(expectedBuilder)

	deepEq := func(input, expected any) int {
		if reflect.DeepEqual(input, expected) {
			return 0
		}
		return -1
	}

	if slices.CompareFunc(inputBuilder.Args(), expectedBuilder.Args(), deepEq) != 0 {
		return false
	}
	return inputBuilder.String() == expectedBuilder.String()
}

// String implements [gomock.Matcher].
func (q *QueryOpts) String() string {
	return fmt.Sprintf("QueryOpts: {%v,%v,%v,%v,%v,%v,%v}\n",
		q.Condition,
		q.OrderBy,
		q.Ordering,
		q.Limit,
		q.Offset,
		q.GroupBy,
		q.Joins,
	)
}

var _ (gomock.Matcher) = (*QueryOpts)(nil)

func (opts *QueryOpts) Write(builder *StatementBuilder) {
	opts.WriteLeftJoins(builder)
	opts.WriteCondition(builder)
	opts.WriteGroupBy(builder)
	opts.WriteOrderBy(builder)
	opts.WriteLimit(builder)
	opts.WriteOffset(builder)
	opts.WriteLock(builder)
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

func (opts *QueryOpts) WriteLock(builder *StatementBuilder) {
	if !opts.ShouldLock {
		return
	}
	builder.WriteString(" FOR UPDATE")
}
