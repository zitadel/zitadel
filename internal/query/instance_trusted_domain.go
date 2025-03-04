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

type InstanceTrustedDomain struct {
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
	Domain       string
	InstanceID   string
}

type InstanceTrustedDomains struct {
	SearchResponse
	Domains []*InstanceTrustedDomain
}

type InstanceTrustedDomainSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *InstanceTrustedDomainSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func NewInstanceTrustedDomainDomainSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(InstanceTrustedDomainDomainCol, value, method)
}

func (q *Queries) SearchInstanceTrustedDomains(ctx context.Context, queries *InstanceTrustedDomainSearchQueries) (domains *InstanceTrustedDomains, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareInstanceTrustedDomainsQuery()
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			InstanceTrustedDomainInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-SGrt4", "Errors.Query.SQLStatement")
	}

	return q.queryInstanceTrustedDomains(ctx, stmt, scan, args...)
}

func (q *Queries) queryInstanceTrustedDomains(ctx context.Context, stmt string, scan func(*sql.Rows) (*InstanceTrustedDomains, error), args ...interface{}) (domains *InstanceTrustedDomains, err error) {
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

func prepareInstanceTrustedDomainsQuery() (sq.SelectBuilder, func(*sql.Rows) (*InstanceTrustedDomains, error)) {
	return sq.Select(
			InstanceTrustedDomainCreationDateCol.identifier(),
			InstanceTrustedDomainChangeDateCol.identifier(),
			InstanceTrustedDomainSequenceCol.identifier(),
			InstanceTrustedDomainDomainCol.identifier(),
			InstanceTrustedDomainInstanceIDCol.identifier(),
			countColumn.identifier(),
		).From(instanceTrustedDomainsTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*InstanceTrustedDomains, error) {
			domains := make([]*InstanceTrustedDomain, 0)
			var count uint64
			for rows.Next() {
				domain := new(InstanceTrustedDomain)
				err := rows.Scan(
					&domain.CreationDate,
					&domain.ChangeDate,
					&domain.Sequence,
					&domain.Domain,
					&domain.InstanceID,
					&count,
				)
				if err != nil {
					return nil, err
				}
				domains = append(domains, domain)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-SDg4h", "Errors.Query.CloseRows")
			}

			return &InstanceTrustedDomains{
				Domains: domains,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

var (
	instanceTrustedDomainsTable = table{
		name:          projection.InstanceTrustedDomainTable,
		instanceIDCol: projection.InstanceTrustedDomainInstanceIDCol,
	}
	InstanceTrustedDomainCreationDateCol = Column{
		name:  projection.InstanceTrustedDomainCreationDateCol,
		table: instanceTrustedDomainsTable,
	}
	InstanceTrustedDomainChangeDateCol = Column{
		name:  projection.InstanceTrustedDomainChangeDateCol,
		table: instanceTrustedDomainsTable,
	}
	InstanceTrustedDomainSequenceCol = Column{
		name:  projection.InstanceTrustedDomainSequenceCol,
		table: instanceTrustedDomainsTable,
	}
	InstanceTrustedDomainDomainCol = Column{
		name:  projection.InstanceTrustedDomainDomainCol,
		table: instanceTrustedDomainsTable,
	}
	InstanceTrustedDomainInstanceIDCol = Column{
		name:  projection.InstanceTrustedDomainInstanceIDCol,
		table: instanceTrustedDomainsTable,
	}
)
