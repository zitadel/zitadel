package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

const orgByIDStmt = "SELECT creation_date, change_date, resource_owner, state, sequence, name, domain FROM zitadel.projections.orgs WHERE id = $1"

func (q *Queries) OrgByID(ctx context.Context, id string) (*Org, error) {
	org := &Org{ID: id}
	row := q.client.QueryRowContext(ctx, orgByIDStmt, id)
	if row.Err() != nil {
		return nil, errors.ThrowNotFound(row.Err(), "QUERY-P3Cmb", "error.internal")
	}
	err := row.Scan(
		&org.CreationDate,
		&org.ChangeDate,
		&org.ResourceOwner,
		&org.State,
		&org.Sequence,
		&org.Name,
		&org.Domain,
	)

	if err != nil {
		if errs.Is(err, sql.ErrNoRows) {
			return nil, errors.ThrowNotFound(err, "QUERY-VQhIc", "errors.orgs.not_found")
		}
		return nil, errors.ThrowInternal(err, "QUERY-yCNV9", "errors.internal")
	}

	return org, nil
}

const orgByDomainStmt = "SELECT id, creation_date, change_date, resource_owner, state, sequence, name, domain FROM zitadel.projections.orgs WHERE domain = $1"

func (q *Queries) OrgByDomainGlobal(ctx context.Context, domain string) (*Org, error) {
	org := new(Org)
	row := q.client.QueryRowContext(ctx, orgByDomainStmt, domain)
	if row.Err() != nil {
		return nil, errors.ThrowNotFound(row.Err(), "QUERY-P3Cmb", "error.internal")
	}
	err := row.Scan(
		&org.ID,
		&org.CreationDate,
		&org.ChangeDate,
		&org.ResourceOwner,
		&org.State,
		&org.Sequence,
		&org.Name,
		&org.Domain,
	)

	if err != nil {
		if errs.Is(err, sql.ErrNoRows) {
			return nil, errors.ThrowNotFound(err, "QUERY-VQhIc", "errors.orgs.not_found")
		}
		return nil, errors.ThrowInternal(err, "QUERY-yCNV9", "errors.internal")
	}

	return org, nil
}

const orgUniqueStmt = "SELECT COUNT(*) = 0 FROM zitadel.projections.orgs WHERE name = $1 AND domain = $2"

func (q *Queries) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	row := q.client.QueryRowContext(ctx, orgUniqueStmt, name, domain)
	if row.Err() != nil {
		return false, errors.ThrowNotFound(row.Err(), "QUERY-eAvMO", "error.internal")
	}
	err = row.Scan(&isUnique)
	if err != nil {
		return false, errors.ThrowInternal(err, "QUERY-uCLym", "errors.internal")
	}

	return isUnique, nil
}

func (q *Queries) ExistsOrg(ctx context.Context, id string) (err error) {
	_, err = q.OrgByID(ctx, id)
	return err
}

var orgsQuery = squirrel.Select("creation_date", "change_date", "resource_owner", "state", "sequence", "name", "domain", "COUNT(name) OVER ()").
	From("zitadel.projections.orgs")

func (q *Queries) SearchOrgs(ctx context.Context, query *OrgSearchQueries) (orgs []*Org, count uint64, err error) {
	stmt, args, err := query.ToQuery(orgsQuery).ToSql()
	if err != nil {
		return nil, 0, errors.ThrowInvalidArgument(err, "QUERY-wQ3by", "Errors.orgs.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, 0, errors.ThrowInternal(err, "QUERY-M6mYN", "Errors.orgs.internal")
	}
	orgs = make([]*Org, 0, query.Limit)
	for rows.Next() {
		org := new(Org)
		rows.Scan(
			&org.CreationDate,
			&org.ChangeDate,
			&org.ResourceOwner,
			&org.State,
			&org.Sequence,
			&org.Name,
			&org.Domain,
			&count,
		)
		orgs = append(orgs, org)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, errors.ThrowInternal(err, "QUERY-pA0Wj", "Errors.orgs.internal")
	}

	return orgs, count, nil
}

type Org struct {
	ID            string          `col:"id"`
	CreationDate  time.Time       `col:"creation_date"`
	ChangeDate    time.Time       `col:"change_date"`
	ResourceOwner string          `col:"resource_owner"`
	State         domain.OrgState `col:"org_state"`
	Sequence      uint64          `col:"sequence"`

	Name   string `col:"name"`
	Domain string `col:"domain"`
}

type OrgSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *OrgSearchQueries) ToQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	query = q.SearchRequest.ToQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

func NewOrgDomainSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery("domain", value, method)
}

func NewOrgNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery("name", value, method)
}

func newOrgIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery("id", id, TextEquals)
}
