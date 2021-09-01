package query

import (
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type SearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn string
	Asc           bool
}

func (req *SearchRequest) ToQuery(query sq.SelectBuilder) sq.SelectBuilder {
	if req.Offset > 0 {
		query = query.Offset(req.Offset)
	}
	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}

	if req.SortingColumn != "" {
		clause := "LOWER(?)"
		if !req.Asc {
			clause += " DESC"
		}
		query.OrderByClause(clause, req.SortingColumn)
	}

	return query
}

const sqlPlaceholder = "?"

type SearchQuery interface {
	ToQuery(sq.SelectBuilder) sq.SelectBuilder
}

type TextQuery struct {
	Column  string
	Text    string
	Compare TextComparison
}

func NewTextQuery(column, value string, compare TextComparison) (*TextQuery, error) {
	if compare < 0 || compare >= textMax {
		return nil, errors.New("invalid compare")
	}
	if column == "" {
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

func (s *TextQuery) comp() map[string]interface{} {
	switch s.Compare {
	case TextEquals:
		return sq.Eq{s.Column: s.Text}
	case TextEqualsIgnore:
		return sq.Eq{"LOWER(" + s.Column + ")": strings.ToLower(s.Text)}
	case TextStartsWith:
		return sq.Like{s.Column: s.Text + sqlPlaceholder}
	case TextStartsWithIgnore:
		return sq.Like{"LOWER(" + s.Column + ")": strings.ToLower(s.Text) + sqlPlaceholder}
	case TextEndsWith:
		return sq.Like{s.Column: sqlPlaceholder + s.Text}
	case TextEndsWithIgnore:
		return sq.Like{"LOWER(" + s.Column + ")": sqlPlaceholder + strings.ToLower(s.Text)}
	case TextContains:
		return sq.Like{s.Column: sqlPlaceholder + s.Text + sqlPlaceholder}
	case TextContainsIgnore:
		return sq.Like{"LOWER(" + s.Column + ")": sqlPlaceholder + strings.ToLower(s.Text) + sqlPlaceholder}
	}
	return nil
}

type TextComparison int

const (
	TextEquals TextComparison = iota
	TextEqualsIgnore
	TextStartsWith
	TextStartsWithIgnore
	TextEndsWith
	TextEndsWithIgnore
	TextContains
	TextContainsIgnore

	textMax
)
