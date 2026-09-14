package query

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type loginEqualitySeek struct {
	sql  string
	args []interface{}
}

// loginNameEqualsFilter marks login-name equality for an indexed ID-seek in
// prepareUsersQuery. Unextracted markers fall back to the view-based exists query.
type loginNameEqualsFilter struct {
	username   string
	domain     string
	loginName  string
	ignoreCase bool
}

func newLoginNameEqualsFilter(value string, ignoreCase bool) (*loginNameEqualsFilter, error) {
	if ignoreCase {
		value = strings.ToLower(value)
	}
	username := value
	domainIndex := strings.LastIndex(value, "@")
	var domainSuffix string
	// split between the last @ (so ignore it if the login name ends with it)
	if domainIndex > 0 && domainIndex != len(value)-1 {
		domainSuffix = value[domainIndex+1:]
		username = value[:domainIndex]
	}
	return &loginNameEqualsFilter{
		username:   username,
		domain:     domainSuffix,
		loginName:  value,
		ignoreCase: ignoreCase,
	}, nil
}

func (q *loginNameEqualsFilter) comparison() TextComparison {
	if q.ignoreCase {
		return TextEqualsIgnoreCase
	}
	return TextEquals
}

func (q *loginNameEqualsFilter) fallback() SearchQuery {
	fallback, _ := newLoginNameExistsViewQuery(q.loginName, q.comparison())
	return fallback
}

func (q *loginNameEqualsFilter) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return q.fallback().toQuery(query)
}

func (q *loginNameEqualsFilter) Col() Column {
	return UserIDCol
}

func (q *loginNameEqualsFilter) comp() sq.Sqlizer {
	return q.fallback().comp()
}

func (q *loginNameEqualsFilter) matchesArgs(instanceID string) []interface{} {
	return []interface{}{
		instanceID,
		instanceID,
		q.domain,
		instanceID,
		q.username,
		q.loginName,
		q.username,
		q.domain,
		q.loginName,
	}
}

func (q *loginNameEqualsFilter) idSeek(instanceID string) loginEqualitySeek {
	// userLoginNameMatchesQuery uses *_lower columns; the _case_sensitive embed uses raw columns.
	subQuery := userLoginNameMatchesQuery
	if !q.ignoreCase {
		subQuery = userLoginNameMatchesCaseSensitiveQuery
	}
	return loginEqualitySeek{
		sql:  "SELECT user_id AS id FROM (" + subQuery + ") AS login_name_matches",
		args: q.matchesArgs(instanceID),
	}
}

// lowerEqualsFilter is username, email, or phone EQUALS / EQUALS_IGNORE_CASE.
// Unextracted filters compile to the same WHERE clause as textQuery.
type lowerEqualsFilter struct {
	textQuery
	tbl   table
	idCol Column
}

func newLowerEqualsSearchQuery(valueCol Column, tbl table, idCol Column, value string, comparison TextComparison) (SearchQuery, error) {
	if comparison != TextEquals && comparison != TextEqualsIgnoreCase {
		return NewTextQuery(valueCol, value, comparison)
	}
	tq, err := NewTextQuery(valueCol, value, comparison)
	if err != nil {
		return nil, err
	}
	return &lowerEqualsFilter{textQuery: *tq, tbl: tbl, idCol: idCol}, nil
}

func (q *lowerEqualsFilter) idSeek(instanceID string) loginEqualitySeek {
	sql := "SELECT " + q.idCol.identifier() + " AS id FROM " + q.tbl.identifier() +
		" WHERE " + q.tbl.InstanceIDIdentifier() + " = ? AND LOWER(" + q.Column.identifier() + ") = "
	args := []interface{}{instanceID}
	if q.Compare == TextEquals {
		sql += "LOWER(?) AND " + q.Column.identifier() + " = ?"
		args = append(args, q.Text, q.Text)
	} else {
		sql += "?"
		args = append(args, strings.ToLower(q.Text))
	}
	return loginEqualitySeek{sql: sql, args: args}
}

func loginEqualitySeekFromLeaf(instanceID string, qry SearchQuery) (loginEqualitySeek, bool) {
	switch v := qry.(type) {
	case *loginNameEqualsFilter:
		if v.loginName == "" {
			return loginEqualitySeek{}, false
		}
		return v.idSeek(instanceID), true
	case *lowerEqualsFilter:
		if v.Text == "" {
			return loginEqualitySeek{}, false
		}
		return v.idSeek(instanceID), true
	default:
		return loginEqualitySeek{}, false
	}
}

func loginEqualitySeeksFromQuery(instanceID string, qry SearchQuery) ([]loginEqualitySeek, bool) {
	if seek, ok := loginEqualitySeekFromLeaf(instanceID, qry); ok {
		return []loginEqualitySeek{seek}, true
	}
	or, ok := qry.(*OrQuery)
	if !ok {
		return nil, false
	}
	seeks := make([]loginEqualitySeek, 0, len(or.queries))
	for _, inner := range or.queries {
		seek, ok := loginEqualitySeekFromLeaf(instanceID, inner)
		if !ok {
			return nil, false
		}
		seeks = append(seeks, seek)
	}
	if len(seeks) == 0 {
		return nil, false
	}
	return seeks, true
}

// extractLoginEqualitySeeks pulls one login-equality filter from queries.
// A match may be a top-level equals leaf, a flat OrQuery of those leaves, or
// the same nested under AndQuery. Remaining AND conjuncts stay as WHERE.
// Equality nested under OrQuery or NotQuery is not extracted.
func extractLoginEqualitySeeks(instanceID string, queries []SearchQuery) ([]loginEqualitySeek, []SearchQuery, bool) {
	for i, qry := range queries {
		if and, ok := qry.(*AndQuery); ok {
			seeks, andRest, ok := extractLoginEqualitySeeks(instanceID, and.queries)
			if !ok {
				continue
			}
			return seeks, spliceRemainingQuery(queries, i, remainingAndQuery(andRest)), true
		}
		seeks, ok := loginEqualitySeeksFromQuery(instanceID, qry)
		if !ok {
			continue
		}
		return seeks, spliceRemainingQuery(queries, i, nil), true
	}
	return nil, queries, false
}

func remainingAndQuery(rest []SearchQuery) SearchQuery {
	switch len(rest) {
	case 0:
		return nil
	case 1:
		return rest[0]
	default:
		and, err := NewAndQuery(rest...)
		if err != nil {
			return nil
		}
		return and
	}
}

func spliceRemainingQuery(queries []SearchQuery, i int, replacement SearchQuery) []SearchQuery {
	n := len(queries) - 1
	if replacement != nil {
		n++
	}
	remaining := make([]SearchQuery, 0, n)
	remaining = append(remaining, queries[:i]...)
	if replacement != nil {
		remaining = append(remaining, replacement)
	}
	remaining = append(remaining, queries[i+1:]...)
	return remaining
}

func joinLoginEqualitySeeks(seeks []loginEqualitySeek) sq.Sqlizer {
	parts := make([]string, len(seeks))
	args := make([]interface{}, 0, len(seeks)*2)
	for i, seek := range seeks {
		parts[i] = seek.sql
		args = append(args, seek.args...)
	}
	return sq.Expr(
		"INNER JOIN ("+strings.Join(parts, " UNION ")+") AS matches ON "+UserIDCol.identifier()+" = matches.id",
		args...,
	)
}
