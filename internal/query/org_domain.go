package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

const (
	domainTable = "zitadel.projections.org_domains"
)

type Domain struct {
	CreationDate   time.Time
	ChangeDate     time.Time
	Sequence       uint64
	Domain         string
	OrgID          string
	IsVerified     bool
	IsPrimary      bool
	ValidationType domain.OrgDomainValidationType
}

type Domains struct {
	SearchResponse
	Domains []*Domain
}

type OrgDomainSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func NewOrgDomainDomainSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(OrgDomainColumnDomain, value, method)
}

func NewOrgDomainOrgIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(OrgDomainColumnDomain, value, TextEquals)
}

type OrgDomainColumn int8

const (
	OrgDomainColumnCreationDate OrgDomainColumn = iota + 1
	OrgDomainColumnChangeDate
	OrgDomainColumnSequence
	OrgDomainColumnDomain
	OrgDomainColumnOrgID
	OrgDomainColumnIsVerified
	OrgDomainColumnPrimary
	OrgDomainColumnValidationType
)

func (q *Queries) SearchOrgDomains(ctx context.Context, queries *OrgDomainSearchQueries) (domains *Domains, err error) {
	query, scan := prepareDomainsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-ZRfj1", "Errors.domains.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M6mYN", "Errors.orgs.internal")
	}
	domains, err = scan(rows)
	if err != nil {
		return nil, err
	}
	domains.LatestSequence, err = q.latestSequence(ctx, domainTable)
	return domains, err
}

func prepareDomainsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Domains, error)) {
	return sq.Select(
			OrgDomainColumnCreationDate.toColumnName(),
			OrgDomainColumnChangeDate.toColumnName(),
			OrgDomainColumnSequence.toColumnName(),
			OrgDomainColumnDomain.toColumnName(),
			OrgDomainColumnOrgID.toColumnName(),
			OrgDomainColumnIsVerified.toColumnName(),
			OrgDomainColumnPrimary.toColumnName(),
			OrgDomainColumnValidationType.toColumnName(),
			"COUNT(*) OVER ()",
		).From(domainTable).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Domains, error) {
			domains := make([]*Domain, 0)
			var count uint64
			for rows.Next() {
				domain := new(Domain)
				err := rows.Scan(
					&domain.CreationDate,
					&domain.ChangeDate,
					&domain.Sequence,
					&domain.Domain,
					&domain.OrgID,
					&domain.IsVerified,
					&domain.IsPrimary,
					&domain.ValidationType,
					&count,
				)
				if err != nil {
					return nil, err
				}
				domains = append(domains, domain)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-rKd6k", "unable to close rows")
			}

			return &Domains{
				Domains: domains,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func (c OrgDomainColumn) toColumnName() string {
	switch c {
	case OrgDomainColumnCreationDate:
		return "creation_date"
	case OrgDomainColumnChangeDate:
		return "change_date"
	case OrgDomainColumnSequence:
		return "sequence"
	case OrgDomainColumnDomain:
		return "domain"
	case OrgDomainColumnOrgID:
		return "org_id"
	case OrgDomainColumnIsVerified:
		return "is_verified"
	case OrgDomainColumnPrimary:
		return "is_primary"
	case OrgDomainColumnValidationType:
		return "validation_type"
	default:
		return ""
	}
}
