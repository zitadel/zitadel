package query

import (
	"errors"

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

type Column interface {
	toFullColumnName() string
	toColumnName() string
}

func (req *SearchRequest) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	if req.Offset > 0 {
		query = query.Offset(req.Offset)
	}
	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}

	if req.SortingColumn != nil {
		clause := "LOWER(" + sqlPlaceholder + ")"
		if !req.Asc {
			clause += " DESC"
		}
		query = query.OrderByClause(clause, req.SortingColumn.toFullColumnName())
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
)

func NewTextQuery(column Column, value string, compare TextComparison) (*TextQuery, error) {
	if compare < 0 || compare >= textCompareMax {
		return nil, ErrInvalidCompare
	}
	if column == nil || column.toFullColumnName() == "" {
		return nil, ErrMissingColumn
	}
	return &TextQuery{
		Column:  column,
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
		return sq.Eq{s.Column.toFullColumnName(): s.Text}, nil
	case TextEqualsIgnoreCase:
		return sq.ILike{s.Column.toFullColumnName(): s.Text}, nil
	case TextStartsWith:
		return sq.Like{s.Column.toFullColumnName(): s.Text + "%"}, nil
	case TextStartsWithIgnoreCase:
		return sq.ILike{s.Column.toFullColumnName(): s.Text + "%"}, nil
	case TextEndsWith:
		return sq.Like{s.Column.toFullColumnName(): "%" + s.Text}, nil
	case TextEndsWithIgnoreCase:
		return sq.ILike{s.Column.toFullColumnName(): "%" + s.Text}, nil
	case TextContains:
		return sq.Like{s.Column.toFullColumnName(): "%" + s.Text + "%"}, nil
	case TextContainsIgnoreCase:
		return sq.ILike{s.Column.toFullColumnName(): "%" + s.Text + "%"}, nil
	case TextListContains:
		return s.Column.toFullColumnName() + " @> ? ", []interface{}{pq.StringArray{s.Text}}
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
