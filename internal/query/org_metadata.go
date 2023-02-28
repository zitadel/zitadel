package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type OrgMetadataList struct {
	SearchResponse
	Metadata []*OrgMetadata
}

type OrgMetadata struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	Key           string
	Value         []byte
}

type OrgMetadataSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

var (
	orgMetadataTable = table{
		name:          projection.OrgMetadataProjectionTable,
		instanceIDCol: projection.OrgMetadataColumnInstanceID,
	}
	OrgMetadataOrgIDCol = Column{
		name:  projection.OrgMetadataColumnOrgID,
		table: orgMetadataTable,
	}
	OrgMetadataCreationDateCol = Column{
		name:  projection.OrgMetadataColumnCreationDate,
		table: orgMetadataTable,
	}
	OrgMetadataChangeDateCol = Column{
		name:  projection.OrgMetadataColumnChangeDate,
		table: orgMetadataTable,
	}
	OrgMetadataResourceOwnerCol = Column{
		name:  projection.OrgMetadataColumnResourceOwner,
		table: orgMetadataTable,
	}
	OrgMetadataInstanceIDCol = Column{
		name:  projection.OrgMetadataColumnInstanceID,
		table: orgMetadataTable,
	}
	OrgMetadataSequenceCol = Column{
		name:  projection.OrgMetadataColumnSequence,
		table: orgMetadataTable,
	}
	OrgMetadataKeyCol = Column{
		name:  projection.OrgMetadataColumnKey,
		table: orgMetadataTable,
	}
	OrgMetadataValueCol = Column{
		name:  projection.OrgMetadataColumnValue,
		table: orgMetadataTable,
	}
	OrgMetadataOwnerRemovedCol = Column{
		name:  projection.OrgMetadataColumnOwnerRemoved,
		table: orgMetadataTable,
	}
)

func (q *Queries) GetOrgMetadataByKey(ctx context.Context, shouldTriggerBulk bool, orgID string, key string, withOwnerRemoved bool, queries ...SearchQuery) (_ *OrgMetadata, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.OrgMetadataProjection.Trigger(ctx)
	}

	query, scan := prepareOrgMetadataQuery(ctx, q.client)
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		OrgMetadataOrgIDCol.identifier():      orgID,
		OrgMetadataKeyCol.identifier():        key,
		OrgMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[OrgMetadataOwnerRemovedCol.identifier()] = false
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-aDaG2", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) SearchOrgMetadata(ctx context.Context, shouldTriggerBulk bool, orgID string, queries *OrgMetadataSearchQueries, withOwnerRemoved bool) (_ *OrgMetadataList, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.OrgMetadataProjection.Trigger(ctx)
	}
	eq := sq.Eq{
		OrgMetadataOrgIDCol.identifier():      orgID,
		OrgMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[OrgMetadataOwnerRemovedCol.identifier()] = false
	}
	query, scan := prepareOrgMetadataListQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Egbld", "Errors.Query.SQLStatment")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Ho2wf", "Errors.Internal")
	}
	metadata, err := scan(rows)
	if err != nil {
		return nil, err
	}
	metadata.LatestSequence, err = q.latestSequence(ctx, orgMetadataTable)
	return metadata, err
}

func (q *OrgMetadataSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (r *OrgMetadataSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewOrgMetadataResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func NewOrgMetadataResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(OrgMetadataResourceOwnerCol, value, TextEquals)
}

func NewOrgMetadataKeySearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(OrgMetadataKeyCol, value, comparison)
}

func prepareOrgMetadataQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*OrgMetadata, error)) {
	return sq.Select(
			OrgMetadataCreationDateCol.identifier(),
			OrgMetadataChangeDateCol.identifier(),
			OrgMetadataResourceOwnerCol.identifier(),
			OrgMetadataSequenceCol.identifier(),
			OrgMetadataKeyCol.identifier(),
			OrgMetadataValueCol.identifier(),
		).
			From(orgMetadataTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*OrgMetadata, error) {
			m := new(OrgMetadata)
			err := row.Scan(
				&m.CreationDate,
				&m.ChangeDate,
				&m.ResourceOwner,
				&m.Sequence,
				&m.Key,
				&m.Value,
			)

			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-Rph32", "Errors.Metadata.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-Hajt2", "Errors.Internal")
			}
			return m, nil
		}
}

func prepareOrgMetadataListQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*OrgMetadataList, error)) {
	return sq.Select(
			OrgMetadataCreationDateCol.identifier(),
			OrgMetadataChangeDateCol.identifier(),
			OrgMetadataResourceOwnerCol.identifier(),
			OrgMetadataSequenceCol.identifier(),
			OrgMetadataKeyCol.identifier(),
			OrgMetadataValueCol.identifier(),
			countColumn.identifier()).
			From(orgMetadataTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*OrgMetadataList, error) {
			metadata := make([]*OrgMetadata, 0)
			var count uint64
			for rows.Next() {
				m := new(OrgMetadata)
				err := rows.Scan(
					&m.CreationDate,
					&m.ChangeDate,
					&m.ResourceOwner,
					&m.Sequence,
					&m.Key,
					&m.Value,
					&count,
				)
				if err != nil {
					return nil, err
				}

				metadata = append(metadata, m)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-dd3gh", "Errors.Query.CloseRows")
			}

			return &OrgMetadataList{
				Metadata: metadata,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
