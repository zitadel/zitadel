package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

const orgTable = "zitadel.projections.orgs"

var (
	orgsQuery = sq.Select(
		OrgColumnID.toColumnName(),
		OrgColumnCreationDate.toColumnName(),
		OrgColumnChangeDate.toColumnName(),
		OrgColumnResourceOwner.toColumnName(),
		OrgColumnState.toColumnName(),
		OrgColumnSequence.toColumnName(),
		OrgColumnName.toColumnName(),
		OrgColumnDomain.toColumnName(),
		"COUNT(name) OVER ()").
		From(orgTable).PlaceholderFormat(sq.Dollar)

	orgQuery = sq.Select(
		OrgColumnCreationDate.toColumnName(),
		OrgColumnChangeDate.toColumnName(),
		OrgColumnResourceOwner.toColumnName(),
		OrgColumnState.toColumnName(),
		OrgColumnSequence.toColumnName(),
		OrgColumnName.toColumnName(),
		OrgColumnDomain.toColumnName()).
		From(orgTable).PlaceholderFormat(sq.Dollar)

	orgUniqueQuery = sq.Select("COUNT(*) = 0").
			From(orgTable).PlaceholderFormat(sq.Dollar)
)

func (q *Queries) OrgByID(ctx context.Context, id string) (*Org, error) {
	query, args, err := orgQuery.Where(sq.Eq{
		OrgColumnID.toColumnName(): id,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-L39bz", "Errors.orgs.invalid.request")
	}

	org := &Org{ID: id}
	row := q.client.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, errors.ThrowNotFound(row.Err(), "QUERY-P3Cmb", "error.internal")
	}
	err = row.Scan(
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

func (q *Queries) OrgByDomainGlobal(ctx context.Context, domain string) (*Org, error) {
	query, args, err := orgQuery.Where(sq.Eq{
		OrgColumnDomain.toColumnName(): domain,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-Me8DB", "Errors.orgs.invalid.request")
	}

	org := new(Org)
	row := q.client.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, errors.ThrowNotFound(row.Err(), "QUERY-5plRd", "error.internal")
	}
	err = row.Scan(
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
			return nil, errors.ThrowNotFound(err, "QUERY-1RBM2", "errors.orgs.not_found")
		}
		return nil, errors.ThrowInternal(err, "QUERY-kBjpb", "errors.internal")
	}

	return org, nil
}

func (q *Queries) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	query, args, err := orgUniqueQuery.Where(sq.Eq{
		OrgColumnName.toColumnName():   name,
		OrgColumnDomain.toColumnName(): domain,
	}).ToSql()
	if err != nil {
		return false, errors.ThrowInvalidArgument(err, "QUERY-OrAo1", "Errors.orgs.invalid.request")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
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
			&org.ID,
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

func NewOrgDomainSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(OrgColumnDomain, value, method)
}

func NewOrgNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(OrgColumnName, value, method)
}

func (q *OrgSearchQueries) ToQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.ToQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

type OrgColumn int32

const (
	OrgColumnCreationDate OrgColumn = iota + 1
	OrgColumnChangeDate
	OrgColumnResourceOwner
	OrgColumnState
	OrgColumnSequence
	OrgColumnName
	OrgColumnDomain
	OrgColumnID
)

func (c OrgColumn) toColumnName() string {
	switch c {
	case OrgColumnCreationDate:
		return "creation_date"
	case OrgColumnChangeDate:
		return "change_date"
	case OrgColumnResourceOwner:
		return "resource_owner"
	case OrgColumnState:
		return "org_state"
	case OrgColumnSequence:
		return "sequence"
	case OrgColumnName:
		return "name"
	case OrgColumnDomain:
		return "domain"
	case OrgColumnID:
		return "id"
	default:
		return ""
	}
}
