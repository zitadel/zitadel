package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type InstanceAllowedDomain struct {
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
	Domain       string
	InstanceID   string
}

type InstanceAllowedDomains struct {
	SearchResponse
	Domains []*InstanceAllowedDomain
}

type InstanceAllowedDomainSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *InstanceAllowedDomainSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func NewInstanceAllowedDomainDomainSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(InstanceAllowedDomainDomainCol, value, method)
}

func (q *Queries) SearchInstanceAllowedDomains(ctx context.Context, queries *InstanceAllowedDomainSearchQueries) (domains *InstanceAllowedDomains, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareInstanceAllowedDomainsQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			InstanceAllowedDomainInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-SGrt4", "Errors.Query.SQLStatement")
	}

	return q.queryInstanceAllowedDomains(ctx, stmt, scan, args...)
}

//
//func (q *Queries) SearchInstanceDomainsGlobal(ctx context.Context, queries *InstanceDomainSearchQueries) (domains *InstanceDomains, err error) {
//	ctx, span := tracing.NewSpan(ctx)
//	defer func() { span.EndWithError(err) }()
//
//	query, scan := prepareInstanceDomainsQuery(ctx, q.client)
//	stmt, args, err := queries.toQuery(query).ToSql()
//	if err != nil {
//		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-IHhLR", "Errors.Query.SQLStatement")
//	}
//
//	return q.queryInstanceDomains(ctx, stmt, scan, args...)
//}

func (q *Queries) queryInstanceAllowedDomains(ctx context.Context, stmt string, scan func(*sql.Rows) (*InstanceAllowedDomains, error), args ...interface{}) (domains *InstanceAllowedDomains, err error) {
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

func prepareInstanceAllowedDomainsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*InstanceAllowedDomains, error)) {
	return sq.Select(
			InstanceAllowedDomainCreationDateCol.identifier(),
			InstanceAllowedDomainChangeDateCol.identifier(),
			InstanceAllowedDomainSequenceCol.identifier(),
			InstanceAllowedDomainDomainCol.identifier(),
			InstanceAllowedDomainInstanceIDCol.identifier(),
			countColumn.identifier(),
		).From(instanceAllowedDomainsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*InstanceAllowedDomains, error) {
			domains := make([]*InstanceAllowedDomain, 0)
			var count uint64
			for rows.Next() {
				domain := new(InstanceAllowedDomain)
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

			return &InstanceAllowedDomains{
				Domains: domains,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

var (
	instanceAllowedDomainsTable = table{
		name:          projection.InstanceAllowedDomainTable,
		instanceIDCol: projection.InstanceAllowedDomainInstanceIDCol,
	}
	InstanceAllowedDomainCreationDateCol = Column{
		name:  projection.InstanceAllowedDomainCreationDateCol,
		table: instanceAllowedDomainsTable,
	}
	InstanceAllowedDomainChangeDateCol = Column{
		name:  projection.InstanceAllowedDomainChangeDateCol,
		table: instanceAllowedDomainsTable,
	}
	InstanceAllowedDomainSequenceCol = Column{
		name:  projection.InstanceAllowedDomainSequenceCol,
		table: instanceAllowedDomainsTable,
	}
	InstanceAllowedDomainDomainCol = Column{
		name:  projection.InstanceAllowedDomainDomainCol,
		table: instanceAllowedDomainsTable,
	}
	InstanceAllowedDomainInstanceIDCol = Column{
		name:  projection.InstanceAllowedDomainInstanceIDCol,
		table: instanceAllowedDomainsTable,
	}
)
