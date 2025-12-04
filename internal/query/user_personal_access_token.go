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
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func patsCheckPermission(ctx context.Context, tokens *PersonalAccessTokens, permissionCheck domain.PermissionCheck) {
	tokens.PersonalAccessTokens = slices.DeleteFunc(tokens.PersonalAccessTokens,
		func(token *PersonalAccessToken) bool {
			return userCheckPermission(ctx, token.ResourceOwner, token.UserID, permissionCheck) != nil
		},
	)
}

var (
	personalAccessTokensTable = table{
		name:          projection.PersonalAccessTokenProjectionTable,
		instanceIDCol: projection.PersonalAccessTokenColumnInstanceID,
	}
	PersonalAccessTokenColumnID = Column{
		name:  projection.PersonalAccessTokenColumnID,
		table: personalAccessTokensTable,
	}
	PersonalAccessTokenColumnUserID = Column{
		name:  projection.PersonalAccessTokenColumnUserID,
		table: personalAccessTokensTable,
	}
	PersonalAccessTokenColumnExpiration = Column{
		name:  projection.PersonalAccessTokenColumnExpiration,
		table: personalAccessTokensTable,
	}
	PersonalAccessTokenColumnScopes = Column{
		name:  projection.PersonalAccessTokenColumnScopes,
		table: personalAccessTokensTable,
	}
	PersonalAccessTokenColumnCreationDate = Column{
		name:  projection.PersonalAccessTokenColumnCreationDate,
		table: personalAccessTokensTable,
	}
	PersonalAccessTokenColumnChangeDate = Column{
		name:  projection.PersonalAccessTokenColumnChangeDate,
		table: personalAccessTokensTable,
	}
	PersonalAccessTokenColumnResourceOwner = Column{
		name:  projection.PersonalAccessTokenColumnResourceOwner,
		table: personalAccessTokensTable,
	}
	PersonalAccessTokenColumnInstanceID = Column{
		name:  projection.PersonalAccessTokenColumnInstanceID,
		table: personalAccessTokensTable,
	}
	PersonalAccessTokenColumnSequence = Column{
		name:  projection.PersonalAccessTokenColumnSequence,
		table: personalAccessTokensTable,
	}
	PersonalAccessTokenColumnOwnerRemoved = Column{
		name:  projection.PersonalAccessTokenColumnOwnerRemoved,
		table: personalAccessTokensTable,
	}
)

type PersonalAccessTokens struct {
	SearchResponse
	PersonalAccessTokens []*PersonalAccessToken
}

type PersonalAccessToken struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	UserID     string
	Expiration time.Time
	Scopes     database.TextArray[string]
}

type PersonalAccessTokenSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) PersonalAccessTokenByID(ctx context.Context, shouldTriggerBulk bool, id string, queries ...SearchQuery) (pat *PersonalAccessToken, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerPersonalAccessTokenProjection")
		ctx, err = projection.PersonalAccessTokenProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := preparePersonalAccessTokenQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		PersonalAccessTokenColumnID.identifier():           id,
		PersonalAccessTokenColumnInstanceID.identifier():   authz.GetInstance(ctx).InstanceID(),
		PersonalAccessTokenColumnOwnerRemoved.identifier(): false,
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dgfb4", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		pat, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	return pat, nil
}

// SearchPersonalAccessTokens returns personal access token resources.
// If permissionCheck is nil, the PATs are not filtered.
// If permissionCheck is not nil and the PermissionCheckV2 feature flag is false, the returned PATs are filtered in-memory by the given permission check.
// If permissionCheck is not nil and the PermissionCheckV2 feature flag is true, the returned PATs are filtered in the database.
func (q *Queries) SearchPersonalAccessTokens(ctx context.Context, queries *PersonalAccessTokenSearchQueries, permissionCheck domain.PermissionCheck) (authNKeys *PersonalAccessTokens, err error) {
	permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	keys, err := q.searchPersonalAccessTokens(ctx, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && !authz.GetFeatures(ctx).PermissionCheckV2 {
		patsCheckPermission(ctx, keys, permissionCheck)
	}
	return keys, nil
}

func (q *Queries) searchPersonalAccessTokens(ctx context.Context, queries *PersonalAccessTokenSearchQueries, permissionCheckV2 bool) (personalAccessTokens *PersonalAccessTokens, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := preparePersonalAccessTokensQuery()
	query = queries.toQuery(query)
	query = userPermissionCheckV2WithCustomColumns(ctx, query, permissionCheckV2, queries.Queries, PersonalAccessTokenColumnResourceOwner, PersonalAccessTokenColumnUserID)
	eq := sq.Eq{
		PersonalAccessTokenColumnInstanceID.identifier():   authz.GetInstance(ctx).InstanceID(),
		PersonalAccessTokenColumnOwnerRemoved.identifier(): false,
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-Hjw2w", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		personalAccessTokens, err = scan(rows)
		return err

	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Bmz63", "Errors.Internal")
	}

	personalAccessTokens.State, err = q.latestState(ctx, personalAccessTokensTable)
	return personalAccessTokens, err
}

func NewPersonalAccessTokenResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(PersonalAccessTokenColumnResourceOwner, value, TextEquals)
}

func NewPersonalAccessTokenUserIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(PersonalAccessTokenColumnUserID, value, TextEquals)
}

func NewPersonalAccessTokenIDQuery(id string) (SearchQuery, error) {
	return NewTextQuery(PersonalAccessTokenColumnID, id, TextEquals)
}

func NewPersonalAccessTokenCreationDateQuery(ts time.Time, compare TimestampComparison) (SearchQuery, error) {
	return NewTimestampQuery(PersonalAccessTokenColumnCreationDate, ts, compare)
}

func NewPersonalAccessTokenExpirationDateDateQuery(ts time.Time, compare TimestampComparison) (SearchQuery, error) {
	return NewTimestampQuery(PersonalAccessTokenColumnExpiration, ts, compare)
}

func (r *PersonalAccessTokenSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewPersonalAccessTokenResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (q *PersonalAccessTokenSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func preparePersonalAccessTokenQuery() (sq.SelectBuilder, func(*sql.Row) (*PersonalAccessToken, error)) {
	return sq.Select(
			PersonalAccessTokenColumnID.identifier(),
			PersonalAccessTokenColumnCreationDate.identifier(),
			PersonalAccessTokenColumnChangeDate.identifier(),
			PersonalAccessTokenColumnResourceOwner.identifier(),
			PersonalAccessTokenColumnSequence.identifier(),
			PersonalAccessTokenColumnUserID.identifier(),
			PersonalAccessTokenColumnExpiration.identifier(),
			PersonalAccessTokenColumnScopes.identifier()).
			From(personalAccessTokensTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*PersonalAccessToken, error) {
			p := new(PersonalAccessToken)
			err := row.Scan(
				&p.ID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.Sequence,
				&p.UserID,
				&p.Expiration,
				&p.Scopes,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-fRunu", "Errors.PersonalAccessToken.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-dj2FF", "Errors.Internal")
			}
			return p, nil
		}
}

func preparePersonalAccessTokensQuery() (sq.SelectBuilder, func(*sql.Rows) (*PersonalAccessTokens, error)) {
	return sq.Select(
			PersonalAccessTokenColumnID.identifier(),
			PersonalAccessTokenColumnCreationDate.identifier(),
			PersonalAccessTokenColumnChangeDate.identifier(),
			PersonalAccessTokenColumnResourceOwner.identifier(),
			PersonalAccessTokenColumnSequence.identifier(),
			PersonalAccessTokenColumnUserID.identifier(),
			PersonalAccessTokenColumnExpiration.identifier(),
			PersonalAccessTokenColumnScopes.identifier(),
			countColumn.identifier()).
			From(personalAccessTokensTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*PersonalAccessTokens, error) {
			personalAccessTokens := make([]*PersonalAccessToken, 0)
			var count uint64
			for rows.Next() {
				token := new(PersonalAccessToken)
				err := rows.Scan(
					&token.ID,
					&token.CreationDate,
					&token.ChangeDate,
					&token.ResourceOwner,
					&token.Sequence,
					&token.UserID,
					&token.Expiration,
					&token.Scopes,
					&count,
				)
				if err != nil {
					return nil, err
				}
				personalAccessTokens = append(personalAccessTokens, token)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-QMXJv", "Errors.Query.CloseRows")
			}

			return &PersonalAccessTokens{
				PersonalAccessTokens: personalAccessTokens,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
