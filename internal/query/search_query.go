package query

import (
	"errors"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/lib/pq"
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
		clause := "LOWER(" + sqlPlaceholder + ")"
		if !req.Asc {
			clause += " DESC"
		}
		query = query.OrderByClause(clause, req.SortingColumn.identifier())
	}

	return query
}

const sqlPlaceholder = "?"

type SearchQuery interface {
	ToQuery(sq.SelectBuilder) sq.SelectBuilder
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

func (q *TextQuery) ToQuery(query sq.SelectBuilder) sq.SelectBuilder {
	where, args := q.comp()
	return query.Where(where, args...)
}

func (s *TextQuery) comp() (comparison interface{}, args []interface{}) {
	switch s.Compare {
	case TextEquals:
		return sq.Eq{s.Column.identifier(): s.Text}, nil
	case TextEqualsIgnoreCase:
		return sq.ILike{s.Column.identifier(): s.Text}, nil
	case TextStartsWith:
		return sq.Like{s.Column.identifier(): s.Text + "%"}, nil
	case TextStartsWithIgnoreCase:
		return sq.ILike{s.Column.identifier(): s.Text + "%"}, nil
	case TextEndsWith:
		return sq.Like{s.Column.identifier(): "%" + s.Text}, nil
	case TextEndsWithIgnoreCase:
		return sq.ILike{s.Column.identifier(): "%" + s.Text}, nil
	case TextContains:
		return sq.Like{s.Column.identifier(): "%" + s.Text + "%"}, nil
	case TextContainsIgnoreCase:
		return sq.ILike{s.Column.identifier(): "%" + s.Text + "%"}, nil
	case TextListContains:
		return s.Column.identifier() + " @> ? ", []interface{}{pq.StringArray{s.Text}}
	}
	return nil, nil
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

	textCompareMax
)

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

func (q *NumberQuery) ToQuery(query sq.SelectBuilder) sq.SelectBuilder {
	where, args := q.comp()
	return query.Where(where, args...)
}

func (s *NumberQuery) comp() (comparison interface{}, args []interface{}) {
	switch s.Compare {
	case NumberEquals:
		return sq.Eq{s.Column.identifier(): s.Number}, nil
	case NumberNotEquals:
		return sq.NotEq{s.Column.identifier(): s.Number}, nil
	case NumberLess:
		return sq.Lt{s.Column.identifier(): s.Number}, nil
	case NumberGreater:
		return sq.Gt{s.Column.identifier(): s.Number}, nil
	case NumberListContains:
		return s.Column.identifier() + " @> ? ", []interface{}{pq.Array(s.Number)}
	}
	return nil, nil
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

func (c Column) identifier() string {
	if c.table.alias == "" {
		return c.name
	}
	return c.table.alias + "." + c.name
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
