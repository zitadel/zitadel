package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type GroupMetadataList struct {
	SearchResponse
	Metadata []*GroupMetadata
}

type GroupMetadata struct {
	CreationDate  time.Time `json:"creation_date,omitempty"`
	GroupID       string    `json:"-"`
	ChangeDate    time.Time `json:"change_date,omitempty"`
	ResourceOwner string    `json:"resource_owner,omitempty"`
	Sequence      uint64    `json:"sequence,omitempty"`
	Key           string    `json:"key,omitempty"`
	Value         []byte    `json:"value,omitempty"`
}

type GroupMetadataSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

var (
	groupMetadataTable = table{
		name:          projection.GroupMetadataProjectionTable,
		instanceIDCol: projection.GroupMetadataColumnInstanceID,
	}
	GroupMetadataGroupIDCol = Column{
		name:  projection.GroupMetadataColumnGroupID,
		table: groupMetadataTable,
	}
	GroupMetadataCreationDateCol = Column{
		name:  projection.GroupMetadataColumnCreationDate,
		table: groupMetadataTable,
	}
	GroupMetadataChangeDateCol = Column{
		name:  projection.GroupMetadataColumnChangeDate,
		table: groupMetadataTable,
	}
	GroupMetadataResourceOwnerCol = Column{
		name:  projection.GroupMetadataColumnResourceOwner,
		table: groupMetadataTable,
	}
	GroupMetadataInstanceIDCol = Column{
		name:  projection.GroupMetadataColumnInstanceID,
		table: groupMetadataTable,
	}
	GroupMetadataSequenceCol = Column{
		name:  projection.GroupMetadataColumnSequence,
		table: groupMetadataTable,
	}
	GroupMetadataKeyCol = Column{
		name:  projection.GroupMetadataColumnKey,
		table: groupMetadataTable,
	}
	GroupMetadataValueCol = Column{
		name:  projection.GroupMetadataColumnValue,
		table: groupMetadataTable,
	}
)

func (q *Queries) GetGroupMetadataByKey(ctx context.Context, shouldTriggerBulk bool, groupID, key string, withOwnerRemoved bool, queries ...SearchQuery) (metadata *GroupMetadata, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerGroupMetadataProjection")
		ctx, err = projection.GroupMetadataProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareGroupMetadataQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		GroupMetadataGroupIDCol.identifier():    groupID,
		GroupMetadataKeyCol.identifier():        key,
		GroupMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-bCGG2", "Errors.Query.SQLStatment")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		metadata, err = scan(row)
		return err
	}, stmt, args...)
	return metadata, err
}

func (q *Queries) SearchGroupMetadataForGroups(ctx context.Context, shouldTriggerBulk bool, groupIDs []string, queries *GroupMetadataSearchQueries) (metadata *GroupMetadataList, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerGroupMetadataProjection")
		ctx, err = projection.GroupMetadataProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareGroupMetadataListQuery()
	eq := sq.Eq{
		GroupMetadataGroupIDCol.identifier():    groupIDs,
		GroupMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Egbgd", "Errors.Query.SQLStatment")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		metadata, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	metadata.State, err = q.latestState(ctx, groupMetadataTable)
	return metadata, err
}

func (q *Queries) SearchGroupMetadata(ctx context.Context, shouldTriggerBulk bool, groupID string, queries *GroupMetadataSearchQueries, withOwnerRemoved bool) (metadata *GroupMetadataList, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerGroupMetadataProjection")
		ctx, err = projection.GroupMetadataProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareGroupMetadataListQuery()
	eq := sq.Eq{
		GroupMetadataGroupIDCol.identifier():    groupID,
		GroupMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Egbgd", "Errors.Query.SQLStatment")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		metadata, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	metadata.State, err = q.latestState(ctx, groupMetadataTable)
	return metadata, err
}

func (q *GroupMetadataSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (r *GroupMetadataSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewGroupMetadataResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func NewGroupMetadataResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(GroupMetadataResourceOwnerCol, value, TextEquals)
}

func NewGroupMetadataKeySearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(GroupMetadataKeyCol, value, comparison)
}

func NewGroupMetadataExistsQuery(key string, value []byte, keyComparison TextComparison, valueComparison BytesComparison) (SearchQuery, error) {
	// linking queries for the subselect
	instanceQuery, err := NewColumnComparisonQuery(GroupMetadataInstanceIDCol, GroupColumnInstanceID, ColumnEquals)
	if err != nil {
		return nil, err
	}

	groupIDQuery, err := NewColumnComparisonQuery(GroupMetadataGroupIDCol, GroupColumnID, ColumnEquals)
	if err != nil {
		return nil, err
	}

	// text query to select data from the linked sub select
	metadataKeyQuery, err := NewTextQuery(GroupMetadataKeyCol, key, keyComparison)
	if err != nil {
		return nil, err
	}

	// text query to select data from the linked sub select
	metadataValueQuery, err := NewBytesQuery(GroupMetadataValueCol, value, valueComparison)
	if err != nil {
		return nil, err
	}

	// full definition of the sub select
	subSelect, err := NewSubSelect(GroupMetadataGroupIDCol, []SearchQuery{instanceQuery, groupIDQuery, metadataKeyQuery, metadataValueQuery})
	if err != nil {
		return nil, err
	}

	// "WHERE * IN (*)" query with subquery as list-data provider
	return NewListQuery(
		GroupColumnID,
		subSelect,
		ListIn,
	)
}

func prepareGroupMetadataQuery() (sq.SelectBuilder, func(*sql.Row) (*GroupMetadata, error)) {
	return sq.Select(
			GroupMetadataCreationDateCol.identifier(),
			GroupMetadataChangeDateCol.identifier(),
			GroupMetadataResourceOwnerCol.identifier(),
			GroupMetadataSequenceCol.identifier(),
			GroupMetadataKeyCol.identifier(),
			GroupMetadataValueCol.identifier(),
		).
			From(groupMetadataTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*GroupMetadata, error) {
			m := new(GroupMetadata)
			err := row.Scan(
				&m.CreationDate,
				&m.ChangeDate,
				&m.ResourceOwner,
				&m.Sequence,
				&m.Key,
				&m.Value,
			)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-Rgh32", "Errors.Metadata.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-Hhjt2", "Errors.Internal")
			}
			return m, nil
		}
}

func prepareGroupMetadataListQuery() (sq.SelectBuilder, func(*sql.Rows) (*GroupMetadataList, error)) {
	return sq.Select(
			GroupMetadataCreationDateCol.identifier(),
			GroupMetadataChangeDateCol.identifier(),
			GroupMetadataGroupIDCol.identifier(),
			GroupMetadataResourceOwnerCol.identifier(),
			GroupMetadataSequenceCol.identifier(),
			GroupMetadataKeyCol.identifier(),
			GroupMetadataValueCol.identifier(),
			countColumn.identifier()).
			From(groupMetadataTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*GroupMetadataList, error) {
			metadata := make([]*GroupMetadata, 0)
			var count uint64
			for rows.Next() {
				m := new(GroupMetadata)
				err := rows.Scan(
					&m.CreationDate,
					&m.ChangeDate,
					&m.GroupID,
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
				return nil, zerrors.ThrowInternal(err, "QUERY-sd3gh", "Errors.Query.CloseRows")
			}

			return &GroupMetadataList{
				Metadata: metadata,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
