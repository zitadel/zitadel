package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareInstanceDomainsQuery()
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			InstanceDomainInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-inlsF", "Errors.Query.SQLStatement")
	}

	return q.queryInstanceDomains(ctx, stmt, scan, args...)
}

func (q *Queries) SearchInstanceDomainsGlobal(ctx context.Context, queries *InstanceDomainSearchQueries) (domains *InstanceDomains, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareInstanceDomainsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-IHhLR", "Errors.Query.SQLStatement")
	}

	return q.queryInstanceDomains(ctx, stmt, scan, args...)
}

func (q *Queries) queryInstanceDomains(ctx context.Context, stmt string, scan func(*sql.Rows) (*InstanceDomains, error), args ...interface{}) (domains *InstanceDomains, err error) {
	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		domains, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	domains.State, err = q.latestState(ctx, instanceDomainsTable)
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
		).From(instanceDomainsTable.identifier()).
			PlaceholderFormat(sq.Dollar),
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
				return nil, zerrors.ThrowInternal(err, "QUERY-8nlWW", "Errors.Query.CloseRows")
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
		name:          projection.InstanceDomainTable,
		instanceIDCol: projection.InstanceDomainInstanceIDCol,
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
