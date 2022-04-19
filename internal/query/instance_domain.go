package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type InstanceDomain struct {
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
	Domain       string
	InstanceID   string
	IsGenerated  bool
	IsPrimary    bool
}

type InstanceDomains struct {
	SearchResponse
	Domains []*InstanceDomain
}

type InstanceDomainSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *InstanceDomainSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func NewInstanceDomainDomainSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(InstanceDomainDomainCol, value, method)
}

func NewInstanceDomainInstanceIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(InstanceDomainInstanceIDCol, value, TextEquals)
}

func NewInstanceDomainGeneratedSearchQuery(generated bool) (SearchQuery, error) {
	return NewBoolQuery(InstanceDomainIsGeneratedCol, generated)
}

func NewInstanceDomainPrimarySearchQuery(primary bool) (SearchQuery, error) {
	return NewBoolQuery(InstanceDomainIsPrimaryCol, primary)
}

func (q *Queries) SearchInstanceDomains(ctx context.Context, queries *InstanceDomainSearchQueries) (domains *InstanceDomains, err error) {
	query, scan := prepareInstanceDomainsQuery()
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			InstanceDomainInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-inlsF", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-38Fni", "Errors.Internal")
	}
	domains, err = scan(rows)
	if err != nil {
		return nil, err
	}
	domains.LatestSequence, err = q.latestSequence(ctx, instanceDomainsTable)
	return domains, err
}

func prepareInstanceDomainsQuery() (sq.SelectBuilder, func(*sql.Rows) (*InstanceDomains, error)) {
	return sq.Select(
			InstanceDomainCreationDateCol.identifier(),
			InstanceDomainChangeDateCol.identifier(),
			InstanceDomainSequenceCol.identifier(),
			InstanceDomainDomainCol.identifier(),
			InstanceDomainInstanceIDCol.identifier(),
			InstanceDomainIsGeneratedCol.identifier(),
			InstanceDomainIsPrimaryCol.identifier(),
			countColumn.identifier(),
		).From(instanceDomainsTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*InstanceDomains, error) {
			domains := make([]*InstanceDomain, 0)
			var count uint64
			for rows.Next() {
				domain := new(InstanceDomain)
				err := rows.Scan(
					&domain.CreationDate,
					&domain.ChangeDate,
					&domain.Sequence,
					&domain.Domain,
					&domain.InstanceID,
					&domain.IsGenerated,
					&domain.IsPrimary,
					&count,
				)
				if err != nil {
					return nil, err
				}
				domains = append(domains, domain)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-8nlWW", "Errors.Query.CloseRows")
			}

			return &InstanceDomains{
				Domains: domains,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

var (
	instanceDomainsTable = table{
		name: projection.InstanceDomainTable,
	}

	InstanceDomainCreationDateCol = Column{
		name:  projection.InstanceDomainCreationDateCol,
		table: instanceDomainsTable,
	}
	InstanceDomainChangeDateCol = Column{
		name:  projection.InstanceDomainChangeDateCol,
		table: instanceDomainsTable,
	}
	InstanceDomainSequenceCol = Column{
		name:  projection.InstanceDomainSequenceCol,
		table: instanceDomainsTable,
	}
	InstanceDomainDomainCol = Column{
		name:  projection.InstanceDomainDomainCol,
		table: instanceDomainsTable,
	}
	InstanceDomainInstanceIDCol = Column{
		name:  projection.InstanceDomainInstanceIDCol,
		table: instanceDomainsTable,
	}
	InstanceDomainIsGeneratedCol = Column{
		name:  projection.InstanceDomainIsGeneratedCol,
		table: instanceDomainsTable,
	}
	InstanceDomainIsPrimaryCol = Column{
		name:  projection.InstanceDomainIsPrimaryCol,
		table: instanceDomainsTable,
	}
)
