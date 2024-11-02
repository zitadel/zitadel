package query

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
)

type SearchResponse struct {
	Count uint64
	*State
}

type SearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn Column
	Asc           bool
}

func (req *SearchRequest) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	if req.Offset > 0 {
		query = query.Offset(req.Offset)
	}
	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}

	if !req.SortingColumn.isZero() {
		clause := req.SortingColumn.orderBy()
		if !req.Asc {
			clause += " DESC"
		}
		query = query.OrderByClause(clause)
	}

	return query
}

type SearchQuery interface {
	toQuery(sq.SelectBuilder) sq.SelectBuilder
	comp() sq.Sqlizer
	Col() Column
}

type NotNullQuery struct {
	Column Column
}

func NewNotNullQuery(col Column) (*NotNullQuery, error) {
	if col.isZero() {
		return nil, ErrMissingColumn
	}
	return &NotNullQuery{
		Column: col,
	}, nil
}

func (q *NotNullQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *NotNullQuery) comp() sq.Sqlizer {
	return sq.NotEq{q.Column.identifier(): nil}
}

func (q *NotNullQuery) Col() Column {
	return q.Column
}

type IsNullQuery struct {
	Column Column
}

func NewIsNullQuery(col Column) (*IsNullQuery, error) {
	if col.isZero() {
		return nil, ErrMissingColumn
	}
	return &IsNullQuery{
		Column: col,
	}, nil
}

func (q *IsNullQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *IsNullQuery) comp() sq.Sqlizer {
	return sq.Eq{q.Column.identifier(): nil}
}
func (q *IsNullQuery) Col() Column {
	return q.Column
}

type OrQuery struct {
	queries []SearchQuery
}

func NewOrQuery(queries ...SearchQuery) (*OrQuery, error) {
	if len(queries) == 0 {
		return nil, ErrMissingColumn
	}
	return &OrQuery{queries: queries}, nil
}

func (q *OrQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *OrQuery) comp() sq.Sqlizer {
	or := make(sq.Or, len(q.queries))
	for i, query := range q.queries {
		or[i] = query.comp()
	}
	return or
}

type AndQuery struct {
	queries []SearchQuery
}

func (q *AndQuery) Col() Column {
	return Column{}
}
func NewAndQuery(queries ...SearchQuery) (*AndQuery, error) {
	if len(queries) == 0 {
		return nil, ErrMissingColumn
	}
	return &AndQuery{queries: queries}, nil
}

func (q *AndQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *AndQuery) comp() sq.Sqlizer {
	and := make(sq.And, len(q.queries))
	for i, query := range q.queries {
		and[i] = query.comp()
	}
	return and
}

type NotQuery struct {
	query SearchQuery
}

func (q *NotQuery) Col() Column {
	return q.query.Col()
}
func NewNotQuery(query SearchQuery) (*NotQuery, error) {
	if query == nil {
		return nil, ErrMissingColumn
	}
	return &NotQuery{query: query}, nil
}

func (q *NotQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (notQ NotQuery) ToSql() (sql string, args []interface{}, err error) {
	querySql, queryArgs, queryErr := notQ.query.comp().ToSql()
	// Handle the error from the query's ToSql() function.
	if queryErr != nil {
		return "", queryArgs, queryErr
	}
	// Construct the SQL statement.
	sql = fmt.Sprintf("NOT (%s)", querySql)
	return sql, queryArgs, nil
}

func (q *NotQuery) comp() sq.Sqlizer {
	return q
}

func (q *OrQuery) Col() Column {
	return Column{}
}

type ColumnComparisonQuery struct {
	Column1 Column
	Compare ColumnComparison
	Column2 Column
}

func NewColumnComparisonQuery(col1 Column, col2 Column, compare ColumnComparison) (*ColumnComparisonQuery, error) {
	if compare < 0 || compare >= columnCompareMax {
		return nil, ErrInvalidCompare
	}
	if col1.isZero() {
		return nil, ErrMissingColumn
	}
	if col2.isZero() {
		return nil, ErrMissingColumn
	}
	return &ColumnComparisonQuery{
		Column1: col1,
		Column2: col2,
		Compare: compare,
	}, nil
}

func (q *ColumnComparisonQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *ColumnComparisonQuery) Col() Column {
	return Column{}
}

func (s *ColumnComparisonQuery) comp() sq.Sqlizer {
	switch s.Compare {
	case ColumnEquals:
		return sq.Expr(s.Column1.identifier() + " = " + s.Column2.identifier())
	case ColumnNotEquals:
		return sq.Expr(s.Column1.identifier() + " != " + s.Column2.identifier())
	}
	return nil
}

type ColumnComparison int

const (
	ColumnEquals ColumnComparison = iota
	ColumnNotEquals

	columnCompareMax
)

type InTextQuery struct {
	Column Column
	Values []string
}

func (q *InTextQuery) Col() Column {
	return q.Column
}

func NewInTextQuery(col Column, values []string) (*InTextQuery, error) {
	if len(values) == 0 {
		return nil, ErrEmptyValues
	}
	if col.isZero() {
		return nil, ErrMissingColumn
	}
	return &InTextQuery{
		Column: col,
		Values: values,
	}, nil
}

type textQuery struct {
	Column  Column
	Text    string
	Compare TextComparison
}

var (
	ErrNothingSelected = errors.New("nothing selected")
	ErrInvalidCompare  = errors.New("invalid compare")
	ErrMissingColumn   = errors.New("missing column")
	ErrInvalidNumber   = errors.New("value is no number")
	ErrEmptyValues     = errors.New("values array must not be empty")
)

func NewTextQuery(col Column, value string, compare TextComparison) (*textQuery, error) {
	if compare < 0 || compare >= textCompareMax {
		return nil, ErrInvalidCompare
	}
	if col.isZero() {
		return nil, ErrMissingColumn
	}
	// handle the comparisons which use (i)like and therefore need to escape potential wildcards in the value
	switch compare {
	case TextEqualsIgnoreCase,
		TextStartsWith,
		TextStartsWithIgnoreCase,
		TextEndsWith,
		TextEndsWithIgnoreCase,
		TextContains,
		TextContainsIgnoreCase:
		value = database.EscapeLikeWildcards(value)
	case TextEquals,
		TextListContains,
		TextNotEquals,
		textCompareMax:
		// do nothing
	}

	return &textQuery{
		Column:  col,
		Text:    value,
		Compare: compare,
	}, nil
}

func (q *textQuery) Col() Column {
	return q.Column
}

func (q *InTextQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *InTextQuery) comp() sq.Sqlizer {
	// This translates to an IN query
	return sq.Eq{q.Column.identifier(): q.Values}
}

func (q *textQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *textQuery) comp() sq.Sqlizer {
	switch q.Compare {
	case TextEquals:
		return sq.Eq{q.Column.identifier(): q.Text}
	case TextNotEquals:
		return sq.NotEq{q.Column.identifier(): q.Text}
	case TextEqualsIgnoreCase:
		return sq.ILike{q.Column.identifier(): q.Text}
	case TextStartsWith:
		return sq.Like{q.Column.identifier(): q.Text + "%"}
	case TextStartsWithIgnoreCase:
		return sq.ILike{q.Column.identifier(): q.Text + "%"}
	case TextEndsWith:
		return sq.Like{q.Column.identifier(): "%" + q.Text}
	case TextEndsWithIgnoreCase:
		return sq.ILike{q.Column.identifier(): "%" + q.Text}
	case TextContains:
		return sq.Like{q.Column.identifier(): "%" + q.Text + "%"}
	case TextContainsIgnoreCase:
		return sq.ILike{q.Column.identifier(): "%" + q.Text + "%"}
	case TextListContains:
		return &listContains{col: q.Column, args: []interface{}{q.Text}}
	case textCompareMax:
		return nil
	}

	return nil
}

type TextComparison int

const (
	TextEquals TextComparison = iota
	TextEqualsIgnoreCase
	TextStartsWith
	TextStartsWithIgnoreCase
	TextEndsWith
	TextEndsWithIgnoreCase
	TextContains
	TextContainsIgnoreCase
	TextListContains
	TextNotEquals

	textCompareMax
)

type NumberQuery struct {
	Column  Column
	Number  interface{}
	Compare NumberComparison
}

func NewNumberQuery(c Column, value interface{}, compare NumberComparison) (*NumberQuery, error) {
	if compare < 0 || compare >= numberCompareMax {
		return nil, ErrInvalidCompare
	}
	if c.isZero() {
		return nil, ErrMissingColumn
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		// everything fine
	default:
		return nil, ErrInvalidNumber
	}
	return &NumberQuery{
		Column:  c,
		Number:  value,
		Compare: compare,
	}, nil
}

func (q *NumberQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *NumberQuery) Col() Column {
	return q.Column
}

func (q *NumberQuery) comp() sq.Sqlizer {
	switch q.Compare {
	case NumberEquals:
		return sq.Eq{q.Column.identifier(): q.Number}
	case NumberNotEquals:
		return sq.NotEq{q.Column.identifier(): q.Number}
	case NumberLess:
		return sq.Lt{q.Column.identifier(): q.Number}
	case NumberGreater:
		return sq.Gt{q.Column.identifier(): q.Number}
	case NumberListContains:
		return &listContains{col: q.Column, args: []interface{}{q.Number}}
	case numberCompareMax:
		return nil
	}

	return nil
}

type NumberComparison int

const (
	NumberEquals NumberComparison = iota
	NumberNotEquals
	NumberLess
	NumberGreater
	NumberListContains

	numberCompareMax
)

// Deprecated: Use NumberComparison, will be removed as soon as all calls are changed to query
func NumberComparisonFromMethod(m domain.SearchMethod) NumberComparison {
	switch m {
	case domain.SearchMethodEquals:
		return NumberEquals
	case domain.SearchMethodNotEquals:
		return NumberNotEquals
	case domain.SearchMethodGreaterThan:
		return NumberGreater
	case domain.SearchMethodLessThan:
		return NumberLess
	case domain.SearchMethodListContains:
		return NumberListContains
	default:
		return numberCompareMax
	}
}

type SubSelect struct {
	Column  Column
	Queries []SearchQuery
}

func NewSubSelect(c Column, queries []SearchQuery) (*SubSelect, error) {
	if len(queries) == 0 {
		return nil, ErrNothingSelected
	}
	if c.isZero() {
		return nil, ErrMissingColumn
	}

	return &SubSelect{
		Column:  c,
		Queries: queries,
	}, nil
}

func (q *SubSelect) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *SubSelect) comp() sq.Sqlizer {
	selectQuery := sq.Select(q.Column.identifier()).From(q.Column.table.identifier())
	for _, query := range q.Queries {
		selectQuery = query.toQuery(selectQuery)
	}
	return selectQuery
}

type listQuery struct {
	Column  Column
	Data    interface{}
	Compare ListComparison
}

func NewListQuery(column Column, value interface{}, compare ListComparison) (*listQuery, error) {
	if compare < 0 || compare >= listCompareMax {
		return nil, ErrInvalidCompare
	}
	if column.isZero() {
		return nil, ErrMissingColumn
	}
	return &listQuery{
		Column:  column,
		Data:    value,
		Compare: compare,
	}, nil
}

func (q *listQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *listQuery) comp() sq.Sqlizer {
	if q.Compare != ListIn {
		return nil
	}

	if subSelect, ok := q.Data.(*SubSelect); ok {
		subSelect, args, err := subSelect.comp().ToSql()
		if err != nil {
			return nil
		}
		return sq.Expr(q.Column.identifier()+" IN ( "+subSelect+" )", args...)
	}
	return sq.Eq{q.Column.identifier(): q.Data}
}

func (q *listQuery) Col() Column {
	return q.Column
}

type ListComparison int

const (
	ListIn ListComparison = iota

	listCompareMax
)

func ListComparisonFromMethod(m domain.SearchMethod) ListComparison {
	switch m {
	case domain.SearchMethodEquals:
		return ListIn
	default:
		return listCompareMax
	}
}

type or struct {
	queries []SearchQuery
}

func Or(queries ...SearchQuery) *or {
	return &or{
		queries: queries,
	}
}

func (q *or) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *or) comp() sq.Sqlizer {
	queries := make([]sq.Sqlizer, 0)
	for _, query := range q.queries {
		queries = append(queries, query.comp())
	}
	return sq.Or(queries)
}

func (q *or) Col() Column {
	return Column{}
}

type BoolQuery struct {
	Column Column
	Value  bool
}

func NewBoolQuery(c Column, value bool) (*BoolQuery, error) {
	return &BoolQuery{
		Column: c,
		Value:  value,
	}, nil
}

func (q *BoolQuery) Col() Column {
	return q.Column
}

func (q *BoolQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *BoolQuery) comp() sq.Sqlizer {
	return sq.Eq{q.Column.identifier(): q.Value}
}

type TimestampComparison int

const (
	TimestampEquals TimestampComparison = iota
	TimestampGreater
	TimestampGreaterOrEquals
	TimestampLess
	TimestampLessOrEquals
)

type TimestampQuery struct {
	Column  Column
	Compare TimestampComparison
	Value   time.Time
}

func NewTimestampQuery(c Column, value time.Time, compare TimestampComparison) (*TimestampQuery, error) {
	return &TimestampQuery{
		Column:  c,
		Compare: compare,
		Value:   value,
	}, nil
}

func (q *TimestampQuery) Col() Column {
	return q.Column
}

func (q *TimestampQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *TimestampQuery) comp() sq.Sqlizer {
	switch q.Compare {
	case TimestampEquals:
		return sq.Eq{q.Column.identifier(): q.Value}
	case TimestampGreater:
		return sq.Gt{q.Column.identifier(): q.Value}
	case TimestampGreaterOrEquals:
		return sq.GtOrEq{q.Column.identifier(): q.Value}
	case TimestampLess:
		return sq.Lt{q.Column.identifier(): q.Value}
	case TimestampLessOrEquals:
		return sq.LtOrEq{q.Column.identifier(): q.Value}
	}
	return nil
}

var (
	// countColumn represents the default counter for search responses
	countColumn = Column{
		name: "COUNT(*) OVER ()",
	}
	// uniqueColumn shows if there are any results
	uniqueColumn = Column{
		name: "COUNT(*) = 0",
	}
)

type table struct {
	name          string
	alias         string
	instanceIDCol string
}

func (t table) setAlias(a string) table {
	t.alias = a
	return t
}

func (t table) identifier() string {
	if t.alias == "" {
		return t.name
	}
	return t.name + " AS " + t.alias
}

func (t table) isZero() bool {
	return t.name == ""
}

func (t table) InstanceIDIdentifier() string {
	if t.alias != "" {
		return t.alias + "." + t.instanceIDCol
	}
	return t.name + "." + t.instanceIDCol
}

type Column struct {
	name           string
	table          table
	isOrderByLower bool
}

func (c Column) identifier() string {
	if c.table.alias != "" {
		return c.table.alias + "." + c.name
	}
	if c.table.name != "" {
		return c.table.name + "." + c.name
	}
	return c.name
}

func (c Column) orderBy() string {
	if !c.isOrderByLower {
		return c.identifier()
	}
	return "LOWER(" + c.identifier() + ")"
}

func (c Column) setTable(t table) Column {
	c.table = t
	return c
}

func (c Column) isZero() bool {
	return c.table.isZero() || c.name == ""
}

func join(join, from Column) string {
	if join.identifier() == join.table.InstanceIDIdentifier() {
		return join.table.identifier() + " ON " + from.identifier() + " = " + join.identifier()
	}
	return join.table.identifier() + " ON " + from.identifier() + " = " + join.identifier() + " AND " + from.table.InstanceIDIdentifier() + " = " + join.table.InstanceIDIdentifier()
}

type listContains struct {
	col  Column
	args interface{}
}

func NewListContains(c Column, value interface{}) (*listContains, error) {
	return &listContains{
		col:  c,
		args: value,
	}, nil
}

func (q *listContains) Col() Column {
	return q.col
}

func (q *listContains) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *listContains) ToSql() (string, []interface{}, error) {
	return q.col.identifier() + " @> ? ", []interface{}{q.args}, nil
}

func (q *listContains) comp() sq.Sqlizer {
	return q
}
