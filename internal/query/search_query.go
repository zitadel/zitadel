package query

import (
	"errors"
	"reflect"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/domain"
)

type SearchResponse struct {
	Count uint64
	*LatestSequence
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

type orQuery struct {
	queries []SearchQuery
}

func newOrQuery(queries ...SearchQuery) (*orQuery, error) {
	if len(queries) == 0 {
		return nil, ErrMissingColumn
	}
	return &orQuery{queries: queries}, nil
}

func (q *orQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *orQuery) comp() sq.Sqlizer {
	or := make(sq.Or, len(q.queries))
	for i, query := range q.queries {
		or[i] = query.comp()
	}
	return or
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
type TextQuery struct {
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

func NewTextQuery(col Column, value string, compare TextComparison) (*TextQuery, error) {
	if compare < 0 || compare >= textCompareMax {
		return nil, ErrInvalidCompare
	}
	if col.isZero() {
		return nil, ErrMissingColumn
	}
	return &TextQuery{
		Column:  col,
		Text:    value,
		Compare: compare,
	}, nil
}

func (q *InTextQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (s *InTextQuery) comp() sq.Sqlizer {
	// This translates to an IN query
	return sq.Eq{s.Column.identifier(): s.Values}
}

func (q *TextQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (s *TextQuery) comp() sq.Sqlizer {
	switch s.Compare {
	case TextEquals:
		return sq.Eq{s.Column.identifier(): s.Text}
	case TextEqualsIgnoreCase:
		return sq.ILike{s.Column.identifier(): s.Text}
	case TextStartsWith:
		return sq.Like{s.Column.identifier(): s.Text + "%"}
	case TextStartsWithIgnoreCase:
		return sq.ILike{s.Column.identifier(): s.Text + "%"}
	case TextEndsWith:
		return sq.Like{s.Column.identifier(): "%" + s.Text}
	case TextEndsWithIgnoreCase:
		return sq.ILike{s.Column.identifier(): "%" + s.Text}
	case TextContains:
		return sq.Like{s.Column.identifier(): "%" + s.Text + "%"}
	case TextContainsIgnoreCase:
		return sq.ILike{s.Column.identifier(): "%" + s.Text + "%"}
	case TextListContains:
		return &listContains{col: s.Column, args: []interface{}{s.Text}}
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

// Deprecated: Use TextComparison, will be removed as soon as all calls are changed to query
func TextComparisonFromMethod(m domain.SearchMethod) TextComparison {
	switch m {
	case domain.SearchMethodEquals:
		return TextEquals
	case domain.SearchMethodEqualsIgnoreCase:
		return TextEqualsIgnoreCase
	case domain.SearchMethodStartsWith:
		return TextStartsWith
	case domain.SearchMethodStartsWithIgnoreCase:
		return TextStartsWithIgnoreCase
	case domain.SearchMethodContains:
		return TextContains
	case domain.SearchMethodContainsIgnoreCase:
		return TextContainsIgnoreCase
	case domain.SearchMethodEndsWith:
		return TextEndsWith
	case domain.SearchMethodEndsWithIgnoreCase:
		return TextEndsWithIgnoreCase
	case domain.SearchMethodListContains:
		return TextListContains
	default:
		return textCompareMax
	}
}

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

func (s *NumberQuery) comp() sq.Sqlizer {
	switch s.Compare {
	case NumberEquals:
		return sq.Eq{s.Column.identifier(): s.Number}
	case NumberNotEquals:
		return sq.NotEq{s.Column.identifier(): s.Number}
	case NumberLess:
		return sq.Lt{s.Column.identifier(): s.Number}
	case NumberGreater:
		return sq.Gt{s.Column.identifier(): s.Number}
	case NumberListContains:
		return &listContains{col: s.Column, args: []interface{}{s.Number}}
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

type ListQuery struct {
	Column  Column
	Data    interface{}
	Compare ListComparison
}

func NewListQuery(column Column, value interface{}, compare ListComparison) (*ListQuery, error) {
	if compare < 0 || compare >= listCompareMax {
		return nil, ErrInvalidCompare
	}
	if column.isZero() {
		return nil, ErrMissingColumn
	}
	return &ListQuery{
		Column:  column,
		Data:    value,
		Compare: compare,
	}, nil
}

func (q *ListQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (s *ListQuery) comp() sq.Sqlizer {
	switch s.Compare {
	case ListIn:
		if subSelect, ok := s.Data.(*SubSelect); ok {
			subSelect, args, err := subSelect.comp().ToSql()
			if err != nil {
				return nil
			}
			return sq.Expr(s.Column.identifier()+" IN ( "+subSelect+" )", args...)
		}
		return sq.Eq{s.Column.identifier(): s.Data}
	}
	return nil
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

func (q *BoolQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (s *BoolQuery) comp() sq.Sqlizer {
	return sq.Eq{s.Column.identifier(): s.Value}
}

var (
	//countColumn represents the default counter for search responses
	countColumn = Column{
		name: "COUNT(*) OVER ()",
	}
	//uniqueColumn shows if there are any results
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

func (q *listContains) ToSql() (string, []interface{}, error) {
	return q.col.identifier() + " @> ? ", []interface{}{q.args}, nil
}
