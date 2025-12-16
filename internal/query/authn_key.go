package query

import (
	"context"
	"database/sql"
	_ "embed"
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

func keysCheckPermission(ctx context.Context, keys *AuthNKeys, permissionCheck domain.PermissionCheck) {
	keys.AuthNKeys = slices.DeleteFunc(keys.AuthNKeys,
		func(key *AuthNKey) bool {
			return userCheckPermission(ctx, key.ResourceOwner, key.AggregateID, permissionCheck) != nil
		},
	)
}

var (
	authNKeyTable = table{
		name:          projection.AuthNKeyTable,
		instanceIDCol: projection.AuthNKeyInstanceIDCol,
	}
	AuthNKeyColumnID = Column{
		name:  projection.AuthNKeyIDCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnCreationDate = Column{
		name:  projection.AuthNKeyCreationDateCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnChangeDate = Column{
		name:  projection.AuthNKeyChangeDateCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnResourceOwner = Column{
		name:  projection.AuthNKeyResourceOwnerCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnInstanceID = Column{
		name:  projection.AuthNKeyInstanceIDCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnAggregateID = Column{
		name:  projection.AuthNKeyAggregateIDCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnSequence = Column{
		name:  projection.AuthNKeySequenceCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnObjectID = Column{
		name:  projection.AuthNKeyObjectIDCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnExpiration = Column{
		name:  projection.AuthNKeyExpirationCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnIdentifier = Column{
		name:  projection.AuthNKeyIdentifierCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnPublicKey = Column{
		name:  projection.AuthNKeyPublicKeyCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnType = Column{
		name:  projection.AuthNKeyTypeCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnEnabled = Column{
		name:  projection.AuthNKeyEnabledCol,
		table: authNKeyTable,
	}
)

type AuthNKeys struct {
	SearchResponse
	AuthNKeys []*AuthNKey
}

type AuthNKey struct {
	ID            string
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	ApplicationID string

	Expiration time.Time
	Type       domain.AuthNKeyType
}

type AuthNKeysData struct {
	SearchResponse
	AuthNKeysData []*AuthNKeyData
}

type AuthNKeyData struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	Expiration time.Time
	Type       domain.AuthNKeyType
	Identifier string
	PublicKey  []byte
}

type AuthNKeySearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *AuthNKeySearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

type JoinFilter int

const (
	JoinFilterUnspecified JoinFilter = iota
	JoinFilterApp
	JoinFilterUserMachine
)

// SearchAuthNKeys returns machine or app keys, depending on the join filter.
// If permissionCheck is nil, the keys are not filtered.
// If permissionCheck is not nil and the PermissionCheckV2 feature flag is false, the returned keys are filtered in-memory by the given permission check.
// If permissionCheck is not nil and the PermissionCheckV2 feature flag is true, the returned keys are filtered in the database.
func (q *Queries) SearchAuthNKeys(ctx context.Context, queries *AuthNKeySearchQueries, filter JoinFilter, permissionCheck domain.PermissionCheck) (authNKeys *AuthNKeys, err error) {
	permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	keys, err := q.searchAuthNKeys(ctx, queries, filter, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && !authz.GetFeatures(ctx).PermissionCheckV2 {
		keysCheckPermission(ctx, keys, permissionCheck)
	}
	return keys, nil
}

func (q *Queries) searchAuthNKeys(ctx context.Context, queries *AuthNKeySearchQueries, joinFilter JoinFilter, permissionCheckV2 bool) (authNKeys *AuthNKeys, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareAuthNKeysQuery()
	query = queries.toQuery(query)
	switch joinFilter {
	case JoinFilterUnspecified:
		// Select all authN keys
	case JoinFilterApp:
		query = query.Join(join(AppColumnID, AuthNKeyColumnObjectID))
	case JoinFilterUserMachine:
		query = query.Join(join(MachineUserIDCol, AuthNKeyColumnIdentifier))
		query = userPermissionCheckV2WithCustomColumns(ctx, query, permissionCheckV2, queries.Queries, AuthNKeyColumnResourceOwner, AuthNKeyColumnIdentifier)
	}
	eq := sq.Eq{
		AuthNKeyColumnEnabled.identifier():    true,
		AuthNKeyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}

	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-SAf3f", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		authNKeys, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dbg53", "Errors.Internal")
	}

	authNKeys.State, err = q.latestState(ctx, authNKeyTable)
	return authNKeys, err
}

func (q *Queries) SearchAuthNKeysData(ctx context.Context, queries *AuthNKeySearchQueries) (authNKeys *AuthNKeysData, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareAuthNKeysDataQuery()
	query = queries.toQuery(query)
	eq := sq.Eq{
		AuthNKeyColumnEnabled.identifier():    true,
		AuthNKeyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-SAg3f", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		authNKeys, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dbi53", "Errors.Internal")
	}
	authNKeys.State, err = q.latestState(ctx, authNKeyTable)
	return authNKeys, err
}

func (q *Queries) GetAuthNKeyByIDWithPermission(ctx context.Context, shouldTriggerBulk bool, id string, permissionCheck domain.PermissionCheck, queries ...SearchQuery) (*AuthNKey, error) {
	key, err := q.GetAuthNKeyByID(ctx, shouldTriggerBulk, id, queries...)
	if err != nil {
		return nil, err
	}

	if err := appCheckPermission(ctx, key.ResourceOwner, key.AggregateID, permissionCheck); err != nil {
		return nil, err
	}

	return key, nil
}

func (q *Queries) GetAuthNKeyByID(ctx context.Context, shouldTriggerBulk bool, id string, queries ...SearchQuery) (key *AuthNKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerAuthNKeyProjection")
		ctx, err = projection.AuthNKeyProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareAuthNKeyQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		AuthNKeyColumnID.identifier():         id,
		AuthNKeyColumnEnabled.identifier():    true,
		AuthNKeyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-AGhg4", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		key, err = scan(row)
		return err
	}, stmt, args...)
	return key, err
}

func NewAuthNKeyResourceOwnerQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AuthNKeyColumnResourceOwner, id, TextEquals)
}

func NewAuthNKeyAggregateIDQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AuthNKeyColumnAggregateID, id, TextEquals)
}

func NewAuthNKeyObjectIDQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AuthNKeyColumnObjectID, id, TextEquals)
}

func NewAuthNKeyIDQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AuthNKeyColumnID, id, TextEquals)
}

func NewAuthNKeyIdentifyerQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AuthNKeyColumnIdentifier, id, TextEquals)
}

func NewAuthNKeyCreationDateQuery(ts time.Time, compare TimestampComparison) (SearchQuery, error) {
	return NewTimestampQuery(AuthNKeyColumnCreationDate, ts, compare)
}

func NewAuthNKeyExpirationDateDateQuery(ts time.Time, compare TimestampComparison) (SearchQuery, error) {
	return NewTimestampQuery(AuthNKeyColumnExpiration, ts, compare)
}

//go:embed authn_key_user.sql
var authNKeyUserQuery string

type AuthNKeyUser struct {
	UserID        string
	ResourceOwner string
	Username      string
	TokenType     domain.OIDCTokenType
	PublicKey     []byte
}

func (q *Queries) GetAuthNKeyUser(ctx context.Context, keyID, userID string) (_ *AuthNKeyUser, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	dst := new(AuthNKeyUser)
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(
			&dst.UserID,
			&dst.ResourceOwner,
			&dst.Username,
			&dst.TokenType,
			&dst.PublicKey,
		)
	},
		authNKeyUserQuery,
		authz.GetInstance(ctx).InstanceID(),
		keyID, userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowNotFound(err, "QUERY-Tha6f", "Errors.AuthNKey.NotFound")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-aen2A", "Errors.Internal")
	}
	return dst, nil
}

func prepareAuthNKeysQuery() (sq.SelectBuilder, func(rows *sql.Rows) (*AuthNKeys, error)) {
	query := sq.Select(
		AuthNKeyColumnID.identifier(),
		AuthNKeyColumnAggregateID.identifier(),
		AuthNKeyColumnCreationDate.identifier(),
		AuthNKeyColumnChangeDate.identifier(),
		AuthNKeyColumnResourceOwner.identifier(),
		AuthNKeyColumnSequence.identifier(),
		AuthNKeyColumnExpiration.identifier(),
		AuthNKeyColumnType.identifier(),
		AuthNKeyColumnObjectID.identifier(),
		countColumn.identifier(),
	).From(authNKeyTable.identifier()).
		PlaceholderFormat(sq.Dollar)

	return query, func(rows *sql.Rows) (*AuthNKeys, error) {
		authNKeys := make([]*AuthNKey, 0)
		var count uint64
		for rows.Next() {
			authNKey := new(AuthNKey)
			err := rows.Scan(
				&authNKey.ID,
				&authNKey.AggregateID,
				&authNKey.CreationDate,
				&authNKey.ChangeDate,
				&authNKey.ResourceOwner,
				&authNKey.Sequence,
				&authNKey.Expiration,
				&authNKey.Type,
				&authNKey.ApplicationID,
				&count,
			)
			if err != nil {
				return nil, err
			}
			authNKeys = append(authNKeys, authNKey)
		}

		if err := rows.Close(); err != nil {
			return nil, zerrors.ThrowInternal(err, "QUERY-Dgfn3", "Errors.Query.CloseRows")
		}

		return &AuthNKeys{
			AuthNKeys: authNKeys,
			SearchResponse: SearchResponse{
				Count: count,
			},
		}, nil
	}
}

func prepareAuthNKeyQuery() (sq.SelectBuilder, func(row *sql.Row) (*AuthNKey, error)) {
	return sq.Select(
			AuthNKeyColumnID.identifier(),
			AuthNKeyColumnCreationDate.identifier(),
			AuthNKeyColumnChangeDate.identifier(),
			AuthNKeyColumnResourceOwner.identifier(),
			AuthNKeyColumnSequence.identifier(),
			AuthNKeyColumnExpiration.identifier(),
			AuthNKeyColumnType.identifier(),
		).From(authNKeyTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*AuthNKey, error) {
			authNKey := new(AuthNKey)
			err := row.Scan(
				&authNKey.ID,
				&authNKey.CreationDate,
				&authNKey.ChangeDate,
				&authNKey.ResourceOwner,
				&authNKey.Sequence,
				&authNKey.Expiration,
				&authNKey.Type,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-Dgr3g", "Errors.AuthNKey.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-BGnbr", "Errors.Internal")
			}
			return authNKey, nil
		}
}

func prepareAuthNKeysDataQuery() (sq.SelectBuilder, func(rows *sql.Rows) (*AuthNKeysData, error)) {
	return sq.Select(
			AuthNKeyColumnID.identifier(),
			AuthNKeyColumnCreationDate.identifier(),
			AuthNKeyColumnChangeDate.identifier(),
			AuthNKeyColumnResourceOwner.identifier(),
			AuthNKeyColumnSequence.identifier(),
			AuthNKeyColumnExpiration.identifier(),
			AuthNKeyColumnType.identifier(),
			AuthNKeyColumnIdentifier.identifier(),
			AuthNKeyColumnPublicKey.identifier(),
			countColumn.identifier(),
		).From(authNKeyTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*AuthNKeysData, error) {
			authNKeys := make([]*AuthNKeyData, 0)
			var count uint64
			for rows.Next() {
				authNKey := new(AuthNKeyData)
				err := rows.Scan(
					&authNKey.ID,
					&authNKey.CreationDate,
					&authNKey.ChangeDate,
					&authNKey.ResourceOwner,
					&authNKey.Sequence,
					&authNKey.Expiration,
					&authNKey.Type,
					&authNKey.Identifier,
					&authNKey.PublicKey,
					&count,
				)
				if err != nil {
					return nil, err
				}
				authNKeys = append(authNKeys, authNKey)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-Dgfn3", "Errors.Query.CloseRows")
			}

			return &AuthNKeysData{
				AuthNKeysData: authNKeys,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
