package query

import (
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
)

type SearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn Column
	Asc           bool
}

type Column interface{ toColumnName() string }

func (req *SearchRequest) ToQuery(query sq.SelectBuilder) sq.SelectBuilder {
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
		query.OrderByClause(clause, req.SortingColumn.toColumnName())
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

func NewTextQuery(column Column, value string, compare TextComparison) (*TextQuery, error) {
	if compare < 0 || compare >= textMax {
		return nil, errors.New("invalid compare")
	}
	if column == nil || column.toColumnName() == "" {
		return nil, errors.New("missing column")
	}
	return &TextQuery{
		Column:  column,
		Text:    value,
		Compare: compare,
	}, nil
}

func (q *TextQuery) ToQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = query.Where(q.comp())
	return query
}

func (s *TextQuery) comp() sq.Sqlizer {
	switch s.Compare {
	case TextEquals:
		return sq.Eq{s.Column.toColumnName(): s.Text}
	case TextEqualsIgnoreCase:
		return sq.Eq{"LOWER(" + s.Column.toColumnName() + ")": strings.ToLower(s.Text)}
	case TextStartsWith:
		return sq.Like{s.Column.toColumnName(): s.Text + "%"}
	case TextStartsWithIgnoreCase:
		return sq.Like{"LOWER(" + s.Column.toColumnName() + ")": strings.ToLower(s.Text) + "%"}
	case TextEndsWith:
		return sq.Like{s.Column.toColumnName(): "%" + s.Text}
	case TextEndsWithIgnoreCase:
		return sq.Like{"LOWER(" + s.Column.toColumnName() + ")": "%" + strings.ToLower(s.Text)}
	case TextContains:
		return sq.Like{s.Column.toColumnName(): "%" + s.Text + "%"}
	case TextContainsIgnoreCase:
		return sq.Like{"LOWER(" + s.Column.toColumnName() + ")": "%" + strings.ToLower(s.Text) + "%"}
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

	textMax
)

func TextCompareFromMethod(m domain.SearchMethod) TextComparison {
	switch m {
	case domain.SearchMethodEquals:
		return TextEquals
	case domain.SearchMethodEqualsIgnoreCase:
		return TextEndsWithIgnoreCase
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
	default:
		return textMax
	}
}
