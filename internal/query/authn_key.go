package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
)

var (
	authNKeyTable = table{
		name: projection.AuthNKeyTable,
	}
	AuthNKeyColumnID = Column{
		name:  projection.AuthNKeyIDCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnCreationDate = Column{
		name:  projection.AuthNKeyCreationDateCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnResourceOwner = Column{
		name:  projection.AuthNKeyResourceOwnerCol,
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
	ResourceOwner string
	Sequence      uint64

	Expiration time.Time
	Type       domain.AuthNKeyType
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

func (q *Queries) SearchAuthNKeys(ctx context.Context, queries *AuthNKeySearchQueries) (authNKeys *AuthNKeys, err error) {
	query, scan := prepareAuthNKeysQuery()
	query = queries.toQuery(query)
	stmt, args, err := query.Where(
		sq.Eq{
			AuthNKeyColumnEnabled.identifier(): true,
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-SAf3f", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dbg53", "Errors.Internal")
	}
	authNKeys, err = scan(rows)
	if err != nil {
		return nil, err
	}
	authNKeys.LatestSequence, err = q.latestSequence(ctx, authNKeyTable)
	return authNKeys, err
}

func (q *Queries) GetAuthNKeyByID(ctx context.Context, shouldRealTime bool, id string, queries ...SearchQuery) (*AuthNKey, error) {
	if shouldRealTime {
		projection.AuthNKeyProjection.TriggerBulk(ctx)
	}
	query, scan := prepareAuthNKeyQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	stmt, args, err := query.Where(
		sq.Eq{
			AuthNKeyColumnID.identifier():      id,
			AuthNKeyColumnEnabled.identifier(): true,
		}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-AGhg4", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) GetAuthNKeyPublicKeyByIDAndIdentifier(ctx context.Context, id string, identifier string) ([]byte, error) {
	stmt, scan := prepareAuthNKeyPublicKeyQuery()
	query, args, err := stmt.Where(
		sq.And{
			sq.Eq{
				AuthNKeyColumnID.identifier():         id,
				AuthNKeyColumnIdentifier.identifier(): identifier,
				AuthNKeyColumnEnabled.identifier():    true,
			},
			sq.Gt{
				AuthNKeyColumnExpiration.identifier(): time.Now(),
			},
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-DAb32", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
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

func prepareAuthNKeysQuery() (sq.SelectBuilder, func(rows *sql.Rows) (*AuthNKeys, error)) {
	return sq.Select(
			AuthNKeyColumnID.identifier(),
			AuthNKeyColumnCreationDate.identifier(),
			AuthNKeyColumnResourceOwner.identifier(),
			AuthNKeyColumnSequence.identifier(),
			AuthNKeyColumnExpiration.identifier(),
			AuthNKeyColumnType.identifier(),
			countColumn.identifier(),
		).From(authNKeyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*AuthNKeys, error) {
			authNKeys := make([]*AuthNKey, 0)
			var count uint64
			for rows.Next() {
				authNKey := new(AuthNKey)
				err := rows.Scan(
					&authNKey.ID,
					&authNKey.CreationDate,
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
				return nil, errors.ThrowInternal(err, "QUERY-Dgfn3", "Errors.Query.CloseRows")
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
			AuthNKeyColumnResourceOwner.identifier(),
			AuthNKeyColumnSequence.identifier(),
			AuthNKeyColumnExpiration.identifier(),
			AuthNKeyColumnType.identifier(),
		).From(authNKeyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*AuthNKey, error) {
			authNKey := new(AuthNKey)
			err := row.Scan(
				&authNKey.ID,
				&authNKey.CreationDate,
				&authNKey.ResourceOwner,
				&authNKey.Sequence,
				&authNKey.Expiration,
				&authNKey.Type,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-Dgr3g", "Errors.AuthNKey.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-BGnbr", "Errors.Internal")
			}
			return authNKey, nil
		}
}

func prepareAuthNKeyPublicKeyQuery() (sq.SelectBuilder, func(row *sql.Row) ([]byte, error)) {
	return sq.Select(
			AuthNKeyColumnPublicKey.identifier(),
		).From(authNKeyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) ([]byte, error) {
			var publicKey []byte
			err := row.Scan(
				&publicKey,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-SDf32", "Errors.AuthNKey.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-Bfs2a", "Errors.Internal")
			}
			return publicKey, nil
		}
}
