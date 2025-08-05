package query

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UserMetadataList struct {
	SearchResponse
	Metadata []*UserMetadata
}

type UserMetadata struct {
	CreationDate  time.Time `json:"creation_date,omitempty"`
	UserID        string    `json:"-"`
	ChangeDate    time.Time `json:"change_date,omitempty"`
	ResourceOwner string    `json:"resource_owner,omitempty"`
	Sequence      uint64    `json:"sequence,omitempty"`
	Key           string    `json:"key,omitempty"`
	Value         []byte    `json:"value,omitempty"`
}

type UserMetadataSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func userMetadataCheckPermission(ctx context.Context, userMetadataList *UserMetadataList, permissionCheck domain.PermissionCheck) {
	userMetadataList.Metadata = slices.DeleteFunc(userMetadataList.Metadata,
		func(userMetadata *UserMetadata) bool {
			return userCheckPermission(ctx, userMetadata.ResourceOwner, userMetadata.UserID, permissionCheck) != nil
		},
	)
}

func userMetadataPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool, queries *UserMetadataSearchQueries) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		UserMetadataResourceOwnerCol,
		domain.PermissionUserRead,
		SingleOrgPermissionOption(queries.Queries),
		OwnedRowsPermissionOption(UserMetadataUserIDCol),
	)
	return query.JoinClause(join, args...)
}

var (
	userMetadataTable = table{
		name:          projection.UserMetadataProjectionTable,
		instanceIDCol: projection.UserMetadataColumnInstanceID,
	}
	UserMetadataUserIDCol = Column{
		name:  projection.UserMetadataColumnUserID,
		table: userMetadataTable,
	}
	UserMetadataCreationDateCol = Column{
		name:  projection.UserMetadataColumnCreationDate,
		table: userMetadataTable,
	}
	UserMetadataChangeDateCol = Column{
		name:  projection.UserMetadataColumnChangeDate,
		table: userMetadataTable,
	}
	UserMetadataResourceOwnerCol = Column{
		name:  projection.UserMetadataColumnResourceOwner,
		table: userMetadataTable,
	}
	UserMetadataInstanceIDCol = Column{
		name:  projection.UserMetadataColumnInstanceID,
		table: userMetadataTable,
	}
	UserMetadataSequenceCol = Column{
		name:  projection.UserMetadataColumnSequence,
		table: userMetadataTable,
	}
	UserMetadataKeyCol = Column{
		name:  projection.UserMetadataColumnKey,
		table: userMetadataTable,
	}
	UserMetadataValueCol = Column{
		name:  projection.UserMetadataColumnValue,
		table: userMetadataTable,
	}
)

func (q *Queries) GetUserMetadataByKey(ctx context.Context, shouldTriggerBulk bool, userID, key string, withOwnerRemoved bool, queries ...SearchQuery) (metadata *UserMetadata, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerUserMetadataProjection")
		ctx, err = projection.UserMetadataProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareUserMetadataQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		UserMetadataUserIDCol.identifier():     userID,
		UserMetadataKeyCol.identifier():        key,
		UserMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-aDGG2", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		metadata, err = scan(row)
		return err
	}, stmt, args...)
	return metadata, err
}

func (q *Queries) SearchUserMetadataForUsers(ctx context.Context, shouldTriggerBulk bool, userIDs []string, queries *UserMetadataSearchQueries) (metadata *UserMetadataList, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerUserMetadataProjection")
		ctx, err = projection.UserMetadataProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareUserMetadataListQuery()
	eq := sq.Eq{
		UserMetadataUserIDCol.identifier():     userIDs,
		UserMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Egbgd", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		metadata, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	metadata.State, err = q.latestState(ctx, userMetadataTable)
	return metadata, err
}

func (q *Queries) SearchUserMetadata(ctx context.Context, shouldTriggerBulk bool, userID string, queries *UserMetadataSearchQueries, permissionCheck domain.PermissionCheck) (metadata *UserMetadataList, err error) {
	permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	users, err := q.searchUserMetadata(ctx, shouldTriggerBulk, userID, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && !authz.GetFeatures(ctx).PermissionCheckV2 {
		userMetadataCheckPermission(ctx, users, permissionCheck)
	}
	return users, nil
}

func (q *Queries) searchUserMetadata(ctx context.Context, shouldTriggerBulk bool, userID string, queries *UserMetadataSearchQueries, permissionCheckV2 bool) (metadata *UserMetadataList, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerUserMetadataProjection")
		ctx, err = projection.UserMetadataProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareUserMetadataListQuery()
	query = userMetadataPermissionCheckV2(ctx, query, permissionCheckV2, queries)
	eq := sq.Eq{
		UserMetadataUserIDCol.identifier():     userID,
		UserMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Egbgd", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		metadata, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	metadata.State, err = q.latestState(ctx, userMetadataTable)
	return metadata, err
}

func (q *UserMetadataSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (r *UserMetadataSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewUserMetadataResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func NewUserMetadataResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(UserMetadataResourceOwnerCol, value, TextEquals)
}

func NewUserMetadataKeySearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(UserMetadataKeyCol, value, comparison)
}

func NewUserMetadataExistsQuery(key string, value []byte, keyComparison TextComparison, valueComparison BytesComparison) (SearchQuery, error) {
	// linking queries for the subselect
	instanceQuery, err := NewColumnComparisonQuery(UserMetadataInstanceIDCol, UserInstanceIDCol, ColumnEquals)
	if err != nil {
		return nil, err
	}

	userIDQuery, err := NewColumnComparisonQuery(UserMetadataUserIDCol, UserIDCol, ColumnEquals)
	if err != nil {
		return nil, err
	}

	// text query to select data from the linked sub select
	metadataKeyQuery, err := NewTextQuery(UserMetadataKeyCol, key, keyComparison)
	if err != nil {
		return nil, err
	}

	// text query to select data from the linked sub select
	metadataValueQuery, err := NewBytesQuery(UserMetadataValueCol, value, valueComparison)
	if err != nil {
		return nil, err
	}

	// full definition of the sub select
	subSelect, err := NewSubSelect(UserMetadataUserIDCol, []SearchQuery{instanceQuery, userIDQuery, metadataKeyQuery, metadataValueQuery})
	if err != nil {
		return nil, err
	}

	// "WHERE * IN (*)" query with subquery as list-data provider
	return NewListQuery(
		UserIDCol,
		subSelect,
		ListIn,
	)
}

func prepareUserMetadataQuery() (sq.SelectBuilder, func(*sql.Row) (*UserMetadata, error)) {
	return sq.Select(
			UserMetadataCreationDateCol.identifier(),
			UserMetadataChangeDateCol.identifier(),
			UserMetadataResourceOwnerCol.identifier(),
			UserMetadataSequenceCol.identifier(),
			UserMetadataKeyCol.identifier(),
			UserMetadataValueCol.identifier(),
		).
			From(userMetadataTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*UserMetadata, error) {
			m := new(UserMetadata)
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

func prepareUserMetadataListQuery() (sq.SelectBuilder, func(*sql.Rows) (*UserMetadataList, error)) {
	return sq.Select(
			UserMetadataCreationDateCol.identifier(),
			UserMetadataChangeDateCol.identifier(),
			UserMetadataUserIDCol.identifier(),
			UserMetadataResourceOwnerCol.identifier(),
			UserMetadataSequenceCol.identifier(),
			UserMetadataKeyCol.identifier(),
			UserMetadataValueCol.identifier(),
			countColumn.identifier()).
			From(userMetadataTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*UserMetadataList, error) {
			metadata := make([]*UserMetadata, 0)
			var count uint64
			for rows.Next() {
				m := new(UserMetadata)
				err := rows.Scan(
					&m.CreationDate,
					&m.ChangeDate,
					&m.UserID,
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

			return &UserMetadataList{
				Metadata: metadata,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
