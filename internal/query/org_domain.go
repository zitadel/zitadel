package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
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

func (q *OrgDomainSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func NewOrgDomainDomainSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(OrgDomainDomainCol, value, method)
}

func NewOrgDomainOrgIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(OrgDomainOrgIDCol, value, TextEquals)
}

func NewOrgDomainVerifiedSearchQuery(verified bool) (SearchQuery, error) {
	return NewBoolQuery(OrgDomainIsVerifiedCol, verified)
}

func (q *Queries) SearchOrgDomains(ctx context.Context, queries *OrgDomainSearchQueries) (domains *Domains, err error) {
	query, scan := prepareDomainsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-ZRfj1", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M6mYN", "Errors.Internal")
	}
	domains, err = scan(rows)
	if err != nil {
		return nil, err
	}
	domains.LatestSequence, err = q.latestSequence(ctx, orgDomainsTable)
	return domains, err
}

func prepareDomainsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Domains, error)) {
	return sq.Select(
			OrgDomainCreationDateCol.identifier(),
			OrgDomainChangeDateCol.identifier(),
			OrgDomainSequenceCol.identifier(),
			OrgDomainDomainCol.identifier(),
			OrgDomainOrgIDCol.identifier(),
			OrgDomainIsVerifiedCol.identifier(),
			OrgDomainIsPrimaryCol.identifier(),
			OrgDomainValidationTypeCol.identifier(),
			countColumn.identifier(),
		).From(orgDomainsTable.identifier()).PlaceholderFormat(sq.Dollar),
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
				return nil, errors.ThrowInternal(err, "QUERY-rKd6k", "Errors.Query.CloseRows")
			}

			return &Domains{
				Domains: domains,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

var (
	orgDomainsTable = table{
		name: projection.OrgDomainTable,
	}

	OrgDomainCreationDateCol = Column{
		name:  projection.OrgDomainCreationDateCol,
		table: orgDomainsTable,
	}
	OrgDomainChangeDateCol = Column{
		name:  projection.OrgDomainChangeDateCol,
		table: orgDomainsTable,
	}
	OrgDomainSequenceCol = Column{
		name:  projection.OrgDomainSequenceCol,
		table: orgDomainsTable,
	}
	OrgDomainDomainCol = Column{
		name:  projection.OrgDomainDomainCol,
		table: orgDomainsTable,
	}
	OrgDomainOrgIDCol = Column{
		name:  projection.OrgDomainOrgIDCol,
		table: orgDomainsTable,
	}
	OrgDomainIsVerifiedCol = Column{
		name:  projection.OrgDomainIsVerifiedCol,
		table: orgDomainsTable,
	}
	OrgDomainIsPrimaryCol = Column{
		name:  projection.OrgDomainIsPrimaryCol,
		table: orgDomainsTable,
	}
	OrgDomainValidationTypeCol = Column{
		name:  projection.OrgDomainValidationTypeCol,
		table: orgDomainsTable,
	}
)
