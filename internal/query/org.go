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

func prepareOrgQuery() (sq.SelectBuilder, func(*sql.Row) (*Org, error)) {
	return sq.Select(
			OrgColumnID.toColumnName(),
			OrgColumnCreationDate.toColumnName(),
			OrgColumnChangeDate.toColumnName(),
			OrgColumnResourceOwner.toColumnName(),
			OrgColumnState.toColumnName(),
			OrgColumnSequence.toColumnName(),
			OrgColumnName.toColumnName(),
			OrgColumnDomain.toColumnName()).
			From(orgTable).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Org, error) {
			o := new(Org)
			err := row.Scan(
				&o.ID,
				&o.CreationDate,
				&o.ChangeDate,
				&o.ResourceOwner,
				&o.State,
				&o.Sequence,
				&o.Name,
				&o.Domain,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-iTTGJ", "errors.orgs.not_found")
				}
				return nil, errors.ThrowInternal(err, "QUERY-pWS5H", "errors.internal")
			}
			return o, nil
		}
}

func (q *Queries) prepareOrgsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Orgs, error)) {
	return sq.Select(
			OrgColumnID.toColumnName(),
			OrgColumnCreationDate.toColumnName(),
			OrgColumnChangeDate.toColumnName(),
			OrgColumnResourceOwner.toColumnName(),
			OrgColumnState.toColumnName(),
			OrgColumnSequence.toColumnName(),
			OrgColumnName.toColumnName(),
			OrgColumnDomain.toColumnName(),
			"COUNT(name) OVER ()").
			From(orgTable).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Orgs, error) {
			orgs := make([]*Org, 0)
			var count uint64
			for rows.Next() {
				org := new(Org)
				err := rows.Scan(
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
				if err != nil {
					return nil, err
				}
				orgs = append(orgs, org)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-QMXJv", "unable to close rows")
			}

			return &Orgs{
				Orgs:  orgs,
				Count: count,
			}, nil
		}
}

func (q *Queries) prepareOrgUniqueQuery() (sq.SelectBuilder, func(*sql.Row) (bool, error)) {
	return sq.Select("COUNT(*) = 0").
			From(orgTable).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (isUnique bool, err error) {
			err = row.Scan(&isUnique)
			if err != nil {
				return false, errors.ThrowInternal(err, "QUERY-pWS5H", "errors.internal")
			}
			return isUnique, err
		}
}

func (q *Queries) OrgByID(ctx context.Context, id string) (*Org, error) {
	stmt, scan := prepareOrgQuery()
	query, args, err := stmt.Where(sq.Eq{
		OrgColumnID.toColumnName(): id,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-AWx52", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) OrgByDomainGlobal(ctx context.Context, domain string) (*Org, error) {
	stmt, scan := prepareOrgQuery()
	query, args, err := stmt.Where(sq.Eq{
		OrgColumnDomain.toColumnName(): domain,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-TYUCE", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	query, scan := q.prepareOrgUniqueQuery()
	stmt, args, err := query.Where(sq.Eq{
		OrgColumnDomain.toColumnName(): domain,
	}).ToSql()
	if err != nil {
		return false, errors.ThrowInternal(err, "QUERY-TYUCE", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) ExistsOrg(ctx context.Context, id string) (err error) {
	_, err = q.OrgByID(ctx, id)
	return err
}

func (q *Queries) SearchOrgs(ctx context.Context, queries *OrgSearchQueries) (orgs *Orgs, err error) {
	query, scan := q.prepareOrgsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-wQ3by", "Errors.orgs.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M6mYN", "Errors.orgs.internal")
	}
	return scan(rows)
}

type Orgs struct {
	Count uint64
	Orgs  []*Org
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

func (q *OrgSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
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
