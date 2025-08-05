package query

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UserSchemas struct {
	SearchResponse
	UserSchemas []*UserSchema
}

func (e *UserSchemas) SetState(s *State) {
	e.State = s
}

type UserSchema struct {
	domain.ObjectDetails
	State                  domain.UserSchemaState
	Type                   string
	Revision               uint32
	Schema                 json.RawMessage
	PossibleAuthenticators database.NumberArray[domain.AuthenticatorType]
}

type UserSchemaSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

var (
	userSchemaTable = table{
		name:          projection.UserSchemaTable,
		instanceIDCol: projection.UserSchemaInstanceIDCol,
	}
	UserSchemaIDCol = Column{
		name:  projection.UserSchemaIDCol,
		table: userSchemaTable,
	}
	UserSchemaCreationDateCol = Column{
		name:  projection.UserSchemaCreationDateCol,
		table: userSchemaTable,
	}
	UserSchemaChangeDateCol = Column{
		name:  projection.UserSchemaChangeDateCol,
		table: userSchemaTable,
	}
	UserSchemaInstanceIDCol = Column{
		name:  projection.UserSchemaInstanceIDCol,
		table: userSchemaTable,
	}
	UserSchemaSequenceCol = Column{
		name:  projection.UserSchemaSequenceCol,
		table: userSchemaTable,
	}
	UserSchemaStateCol = Column{
		name:  projection.UserSchemaStateCol,
		table: userSchemaTable,
	}
	UserSchemaTypeCol = Column{
		name:  projection.UserSchemaTypeCol,
		table: userSchemaTable,
	}
	UserSchemaRevisionCol = Column{
		name:  projection.UserSchemaRevisionCol,
		table: userSchemaTable,
	}
	UserSchemaSchemaCol = Column{
		name:  projection.UserSchemaSchemaCol,
		table: userSchemaTable,
	}
	UserSchemaPossibleAuthenticatorsCol = Column{
		name:  projection.UserSchemaPossibleAuthenticatorsCol,
		table: userSchemaTable,
	}
)

func (q *Queries) GetUserSchemaByID(ctx context.Context, id string) (userSchema *UserSchema, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	eq := sq.Eq{
		UserSchemaIDCol.identifier():         id,
		UserSchemaInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}

	query, scan := prepareUserSchemaQuery()
	return genericRowQuery(ctx, q.client, query.Where(eq), scan)
}

func (q *Queries) SearchUserSchema(ctx context.Context, queries *UserSchemaSearchQueries) (userSchemas *UserSchemas, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	eq := sq.Eq{
		UserSchemaInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}

	query, scan := prepareUserSchemasQuery()
	return genericRowsQueryWithState(ctx, q.client, userSchemaTable, combineToWhereStmt(query, queries.toQuery, eq), scan)
}

func (q *UserSchemaSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func NewUserSchemaTypeSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(UserSchemaTypeCol, value, comparison)
}

func NewUserSchemaIDSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(UserSchemaIDCol, value, comparison)
}

func NewUserSchemaStateSearchQuery(value domain.UserSchemaState) (SearchQuery, error) {
	return NewNumberQuery(UserSchemaStateCol, value, NumberEquals)
}

func prepareUserSchemaQuery() (sq.SelectBuilder, func(*sql.Row) (*UserSchema, error)) {
	return sq.Select(
			UserSchemaIDCol.identifier(),
			UserSchemaCreationDateCol.identifier(),
			UserSchemaChangeDateCol.identifier(),
			UserSchemaSequenceCol.identifier(),
			UserSchemaInstanceIDCol.identifier(),
			UserSchemaStateCol.identifier(),
			UserSchemaTypeCol.identifier(),
			UserSchemaRevisionCol.identifier(),
			UserSchemaSchemaCol.identifier(),
			UserSchemaPossibleAuthenticatorsCol.identifier(),
		).
			From(userSchemaTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*UserSchema, error) {
			u := new(UserSchema)
			var schema database.ByteArray[byte]
			err := row.Scan(
				&u.ID,
				&u.CreationDate,
				&u.EventDate,
				&u.Sequence,
				&u.ResourceOwner,
				&u.State,
				&u.Type,
				&u.Revision,
				&schema,
				&u.PossibleAuthenticators,
			)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-SAF3t", "Errors.UserSchema.NotExists")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-WRB2Q", "Errors.Internal")
			}

			u.Schema = json.RawMessage(schema)

			return u, nil
		}
}

func prepareUserSchemasQuery() (sq.SelectBuilder, func(*sql.Rows) (*UserSchemas, error)) {
	return sq.Select(
			UserSchemaIDCol.identifier(),
			UserSchemaCreationDateCol.identifier(),
			UserSchemaChangeDateCol.identifier(),
			UserSchemaSequenceCol.identifier(),
			UserSchemaInstanceIDCol.identifier(),
			UserSchemaStateCol.identifier(),
			UserSchemaTypeCol.identifier(),
			UserSchemaRevisionCol.identifier(),
			UserSchemaSchemaCol.identifier(),
			UserSchemaPossibleAuthenticatorsCol.identifier(),
			countColumn.identifier()).
			From(userSchemaTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*UserSchemas, error) {
			schemas := make([]*UserSchema, 0)
			var (
				schema database.ByteArray[byte]
				count  uint64
			)

			for rows.Next() {
				u := new(UserSchema)
				err := rows.Scan(
					&u.ID,
					&u.CreationDate,
					&u.EventDate,
					&u.Sequence,
					&u.ResourceOwner,
					&u.State,
					&u.Type,
					&u.Revision,
					&schema,
					&u.PossibleAuthenticators,
					&count,
				)
				if err != nil {
					return nil, err
				}

				u.Schema = json.RawMessage(schema)
				schemas = append(schemas, u)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-Lwj2e", "Errors.Query.CloseRows")
			}

			return &UserSchemas{
				UserSchemas: schemas,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
