package query

import (
	"errors"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/domain"
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

type TextQuery struct {
	Column  Column
	Text    string
	Compare TextComparison
}

var (
	ErrInvalidCompare = errors.New("invalid compare")
	ErrMissingColumn  = errors.New("missing column")
	ErrInvalidNumber  = errors.New("value is no number")
)

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
		return &listContains{col: s.Column, args: []interface{}{pq.StringArray{s.Text}}}
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

//Deprecated: Use TextComparison, will be removed as soon as all calls are changed to query
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
		//everything fine
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
		return &listContains{col: s.Column, args: []interface{}{pq.GenericArray{s.Number}}}
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

//Deprecated: Use NumberComparison, will be removed as soon as all calls are changed to query
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

type ListQuery struct {
	Column  Column
	List    []interface{}
	Compare ListComparison
}

func NewListQuery(column Column, value []interface{}, compare ListComparison) (*ListQuery, error) {
	if compare < 0 || compare >= listCompareMax {
		return nil, ErrInvalidCompare
	}
	if column.isZero() {
		return nil, ErrMissingColumn
	}
	return &ListQuery{
		Column:  column,
		List:    value,
		Compare: compare,
	}, nil
}

func (q *ListQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (s *ListQuery) comp() sq.Sqlizer {
	switch s.Compare {
	case ListIn:
		return sq.Eq{s.Column.identifier(): s.List}
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
	name  string
	alias string
}

func (t table) setAlias(a string) table {
	t.alias = a
	return t
}

func (t table) identifier() string {
	if t.alias == "" {
		return t.name
	}
	return t.name + " as " + t.alias
}

func (t table) isZero() bool {
	return t.name == ""
}

type Column struct {
	name  string
	table table
}

type StringColumn Column

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
	return c.identifier()
}

func (c StringColumn) orderBy() string {
	return "LOWER(" + Column(c).identifier() + ")"
}

func (c Column) setTable(t table) Column {
	c.table = t
	return c
}

func (c Column) isZero() bool {
	return c.table.isZero() || c.name == ""
}

func join(join, from Column) string {
	return join.table.identifier() + " ON " + from.identifier() + " = " + join.identifier()
}

type listContains struct {
	col  Column
	args []interface{}
}

func (q *listContains) ToSql() (string, []interface{}, error) {
	return q.col.identifier() + " @> ? ", q.args, nil
}
