package query

import (
	sq "github.com/Masterminds/squirrel"
)

// loginNamePatternFilter is a login name filter other than EQUALS /
// EQUALS_IGNORE_CASE: STARTS_WITH, ENDS_WITH and CONTAINS, case-sensitive or
// IGNORE_CASE. The login_names3 view computes login names, so a pattern cannot
// use an index.
//
// planUserSearch binds the filter to the instance, and to the organization of an
// organization filter that all users must match (see findConjunctTextEquals). Bound, it reads the matching user IDs from
// login_names3 once per statement. Unbound, it falls back to the view query
// correlated to the outer user, whose plan Postgres can turn into one
// evaluation of the view per user.
type loginNamePatternFilter struct {
	match      *textQuery
	correlated SearchQuery

	binding *loginNamePatternBinding
}

// loginNamePatternBinding is what planUserSearch knows about the statement a
// loginNamePatternFilter is part of.
type loginNamePatternBinding struct {
	instanceID string
	// resourceOwner restricts the matches to one organization if not empty.
	resourceOwner string
	// conjunct is true if the filter is ANDed with the rest of the WHERE
	// clause, false if it is nested in OR or NOT.
	conjunct bool
}

func newLoginNamePatternFilter(value string, comparison TextComparison) (*loginNamePatternFilter, error) {
	match, err := NewTextQuery(LoginNameNameCol, value, comparison)
	if err != nil {
		return nil, err
	}
	correlated, err := newLoginNameExistsViewQuery(value, comparison)
	if err != nil {
		return nil, err
	}
	return &loginNamePatternFilter{
		match:      match,
		correlated: correlated,
	}, nil
}

func (q *loginNamePatternFilter) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return query.Where(q.comp())
}

func (q *loginNamePatternFilter) Col() Column {
	return UserIDCol
}

func (q *loginNamePatternFilter) comp() sq.Sqlizer {
	if q.binding == nil {
		return q.correlated.comp()
	}
	matches := sq.Select(LoginNameUserIDCol.identifier()).
		From(loginNameTable.identifier()).
		Where(sq.Eq{LoginNameInstanceIDCol.identifier(): q.binding.instanceID})
	if q.binding.resourceOwner != "" {
		matches = matches.Where(sq.Eq{LoginNameResourceOwnerCol.identifier(): q.binding.resourceOwner})
	}
	matches = q.match.toQuery(matches)
	if q.binding.conjunct {
		// The IDs are collected once into an array (an InitPlan) and looked up
		// by primary key, as in loginEqualitySeeksCondition, whatever the
		// planner estimates for the pattern.
		return sq.ConcatExpr(UserIDCol.identifier()+" = ANY(ARRAY(", matches, "))")
	}
	// Under OR or NOT the array would be a filter that Postgres searches
	// element by element for every user: it hashes only constant arrays.
	// An uncorrelated IN is hashed once (a hashed SubPlan).
	return sq.ConcatExpr(UserIDCol.identifier()+" IN (", matches, ")")
}

// bindLoginNamePatterns binds every login name pattern filter in queries,
// including those nested in AND, OR and NOT, to the instance and, if not empty,
// to resourceOwner. Pass the organization only if queries also restrict users
// to it: a pattern restricted to the organization then finds the same users,
// because a user's login names carry the user's organization.
// The caller's queries are not modified; changed parts are copied.
func bindLoginNamePatterns(queries []SearchQuery, instanceID, resourceOwner string) []SearchQuery {
	bound, _ := bindLoginNamePatternList(queries, loginNamePatternBinding{
		instanceID:    instanceID,
		resourceOwner: resourceOwner,
		conjunct:      true,
	})
	return bound
}

func bindLoginNamePatternList(queries []SearchQuery, binding loginNamePatternBinding) ([]SearchQuery, bool) {
	var bound []SearchQuery
	for i, qry := range queries {
		replacement, changed := bindLoginNamePattern(qry, binding)
		if !changed {
			continue
		}
		if bound == nil {
			bound = append([]SearchQuery(nil), queries...)
		}
		bound[i] = replacement
	}
	if bound == nil {
		return queries, false
	}
	return bound, true
}

func bindLoginNamePattern(qry SearchQuery, binding loginNamePatternBinding) (SearchQuery, bool) {
	nested := binding
	nested.conjunct = false
	switch v := qry.(type) {
	case *loginNamePatternFilter:
		bound := *v
		bound.binding = &binding
		return &bound, true
	case *AndQuery:
		queries, changed := bindLoginNamePatternList(v.queries, binding)
		if !changed {
			return qry, false
		}
		return &AndQuery{queries: queries}, true
	case *OrQuery:
		queries, changed := bindLoginNamePatternList(v.queries, nested)
		if !changed {
			return qry, false
		}
		return &OrQuery{queries: queries}, true
	case *NotQuery:
		inner, changed := bindLoginNamePattern(v.query, nested)
		if !changed {
			return qry, false
		}
		return &NotQuery{query: inner}, true
	default:
		return qry, false
	}
}

// findConjunctTextEquals returns the value of an EQUALS filter on column that
// every result must match: one in queries or, recursively, in an AND among
// them. Filters under OR or NOT do not count.
func findConjunctTextEquals(column Column, queries []SearchQuery) (string, bool) {
	if text, ok := findTextEqualsQuery(column, queries); ok {
		return text, true
	}
	for _, qry := range queries {
		if and, ok := qry.(*AndQuery); ok {
			if text, ok := findConjunctTextEquals(column, and.queries); ok {
				return text, true
			}
		}
	}
	return "", false
}
