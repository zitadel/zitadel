package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

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
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

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

func (q *Queries) SearchAuthNKeys(ctx context.Context, queries *AuthNKeySearchQueries, withOwnerRemoved bool) (authNKeys *AuthNKeys, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareAuthNKeysQuery(ctx, q.client)
	query = queries.toQuery(query)
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

	query, scan := prepareAuthNKeysDataQuery(ctx, q.client)
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

func (q *Queries) GetAuthNKeyByID(ctx context.Context, shouldTriggerBulk bool, id string, queries ...SearchQuery) (key *AuthNKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerAuthNKeyProjection")
		ctx, err = projection.AuthNKeyProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareAuthNKeyQuery(ctx, q.client)
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

func (q *Queries) GetAuthNKeyPublicKeyByIDAndIdentifier(ctx context.Context, id string, identifier string, withOwnerRemoved bool) (key []byte, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAuthNKeyPublicKeyQuery(ctx, q.client)
	eq := sq.And{
		sq.Eq{
			AuthNKeyColumnID.identifier():         id,
			AuthNKeyColumnIdentifier.identifier(): identifier,
			AuthNKeyColumnEnabled.identifier():    true,
			AuthNKeyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		},
		sq.Gt{
			AuthNKeyColumnExpiration.identifier(): time.Now(),
		},
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-DAb32", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		key, err = scan(row)
		return err
	}, query, args...)
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

func prepareAuthNKeysQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (*AuthNKeys, error)) {
	return sq.Select(
			AuthNKeyColumnID.identifier(),
			AuthNKeyColumnCreationDate.identifier(),
			AuthNKeyColumnChangeDate.identifier(),
			AuthNKeyColumnResourceOwner.identifier(),
			AuthNKeyColumnSequence.identifier(),
			AuthNKeyColumnExpiration.identifier(),
			AuthNKeyColumnType.identifier(),
			countColumn.identifier(),
		).From(authNKeyTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*AuthNKeys, error) {
			authNKeys := make([]*AuthNKey, 0)
			var count uint64
			for rows.Next() {
				authNKey := new(AuthNKey)
				err := rows.Scan(
					&authNKey.ID,
					&authNKey.CreationDate,
					&authNKey.ChangeDate,
					&authNKey.ResourceOwner,
					&authNKey.Sequence,
					&authNKey.Expiration,
					&authNKey.Type,
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

func prepareAuthNKeyQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(row *sql.Row) (*AuthNKey, error)) {
	return sq.Select(
			AuthNKeyColumnID.identifier(),
			AuthNKeyColumnCreationDate.identifier(),
			AuthNKeyColumnChangeDate.identifier(),
			AuthNKeyColumnResourceOwner.identifier(),
			AuthNKeyColumnSequence.identifier(),
			AuthNKeyColumnExpiration.identifier(),
			AuthNKeyColumnType.identifier(),
		).From(authNKeyTable.identifier() + db.Timetravel(call.Took(ctx))).
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

func prepareAuthNKeyPublicKeyQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(row *sql.Row) ([]byte, error)) {
	return sq.Select(
			AuthNKeyColumnPublicKey.identifier(),
		).From(authNKeyTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) ([]byte, error) {
			var publicKey []byte
			err := row.Scan(
				&publicKey,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-SDf32", "Errors.AuthNKey.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-Bfs2a", "Errors.Internal")
			}
			return publicKey, nil
		}
}

func prepareAuthNKeysDataQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (*AuthNKeysData, error)) {
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
		).From(authNKeyTable.identifier() + db.Timetravel(call.Took(ctx))).
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
