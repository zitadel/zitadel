package query

import (
	"context"
	"database/sql"
	"encoding/json"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Sessions struct {
	SearchResponse
	Sessions []*Session
}

type Session struct {
	ID             string
	CreationDate   time.Time
	ChangeDate     time.Time
	Sequence       uint64
	State          domain.SessionState
	ResourceOwner  string
	Creator        string
	UserFactor     SessionUserFactor
	PasswordFactor SessionPasswordFactor
	Metadata       map[string][]byte
}

type SessionUserFactor struct {
	UserID        string
	UserCheckedAt time.Time
	LoginName     string
	DisplayName   string
}

type SessionPasswordFactor struct {
	PasswordCheckedAt time.Time
}

type SessionsSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *SessionsSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

var (
	sessionsTable = table{
		name:          projection.SessionsProjectionTable,
		instanceIDCol: projection.SessionColumnInstanceID,
	}
	SessionColumnID = Column{
		name:  projection.SessionColumnID,
		table: sessionsTable,
	}
	SessionColumnCreationDate = Column{
		name:  projection.SessionColumnCreationDate,
		table: sessionsTable,
	}
	SessionColumnChangeDate = Column{
		name:  projection.SessionColumnChangeDate,
		table: sessionsTable,
	}
	SessionColumnSequence = Column{
		name:  projection.SessionColumnSequence,
		table: sessionsTable,
	}
	SessionColumnState = Column{
		name:  projection.SessionColumnState,
		table: sessionsTable,
	}
	SessionColumnResourceOwner = Column{
		name:  projection.SessionColumnResourceOwner,
		table: sessionsTable,
	}
	SessionColumnInstanceID = Column{
		name:  projection.SessionColumnInstanceID,
		table: sessionsTable,
	}
	SessionColumnCreator = Column{
		name:  projection.SessionColumnCreator,
		table: sessionsTable,
	}
	SessionColumnUserID = Column{
		name:  projection.SessionColumnUserID,
		table: sessionsTable,
	}
	SessionColumnUserCheckedAt = Column{
		name:  projection.SessionColumnUserCheckedAt,
		table: sessionsTable,
	}
	SessionColumnPasswordCheckedAt = Column{
		name:  projection.SessionColumnPasswordCheckedAt,
		table: sessionsTable,
	}
	SessionColumnMetadata = Column{
		name:  projection.SessionColumnMetadata,
		table: sessionsTable,
	}
	SessionColumnToken = Column{
		name:  projection.SessionColumnToken,
		table: sessionsTable,
	}
)

func (q *Queries) SessionByID(ctx context.Context, id string) (_ *Session, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSessionQuery(ctx, q.client)
	stmt, args, err := query.Where(
		sq.Eq{
			SessionColumnID.identifier():         id,
			SessionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-dn9JW", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) SearchSessions(ctx context.Context, queries *SessionsSearchQueries) (_ *Sessions, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSessionsQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			SessionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-sn9Jf", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil || rows.Err() != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Sfg42", "Errors.Internal")
	}
	sessions, err := scan(rows)
	if err != nil {
		return nil, err
	}
	sessions.LatestSequence, err = q.latestSequence(ctx, sessionsTable)
	return sessions, err
}

func NewSessionIDsSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(SessionColumnID, list, ListIn)
}

func prepareSessionQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Session, error)) {
	return sq.Select(
			SessionColumnID.identifier(),
			SessionColumnCreationDate.identifier(),
			SessionColumnChangeDate.identifier(),
			SessionColumnSequence.identifier(),
			SessionColumnState.identifier(),
			SessionColumnResourceOwner.identifier(),
			SessionColumnCreator.identifier(),
			SessionColumnUserID.identifier(),
			SessionColumnUserCheckedAt.identifier(),
			LoginNameNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			SessionColumnPasswordCheckedAt.identifier(),
			SessionColumnMetadata.identifier(),
		).From(sessionsTable.identifier()).
			LeftJoin(join(LoginNameUserIDCol, SessionColumnUserID)).
			LeftJoin(join(HumanUserIDCol, SessionColumnUserID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*Session, error) {
			session := new(Session)

			var (
				userID            sql.NullString
				userCheckedAt     sql.NullTime
				loginName         sql.NullString
				displayName       sql.NullString
				passwordCheckedAt sql.NullTime
				metadata          []byte
			)

			err := row.Scan(
				&session.ID,
				&session.CreationDate,
				&session.ChangeDate,
				&session.Sequence,
				&session.State,
				&session.ResourceOwner,
				&session.Creator,
				&userID,
				&userCheckedAt,
				&loginName,
				&displayName,
				&passwordCheckedAt,
				&metadata,
			)

			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-SFeaa", "Errors.Session.NotExisting")
				}
				return nil, errors.ThrowInternal(err, "QUERY-SAder", "Errors.Internal")
			}
			session.UserFactor.UserID = userID.String
			session.UserFactor.UserCheckedAt = userCheckedAt.Time
			session.UserFactor.LoginName = loginName.String
			session.UserFactor.DisplayName = displayName.String
			session.PasswordFactor.PasswordCheckedAt = passwordCheckedAt.Time
			if len(metadata) != 0 {
				if err := json.Unmarshal(metadata, &session.Metadata); err != nil {
					return nil, err
				}
			}

			return session, nil
		}
}

func prepareSessionsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Sessions, error)) {
	return sq.Select(
			SessionColumnID.identifier(),
			SessionColumnCreationDate.identifier(),
			SessionColumnChangeDate.identifier(),
			SessionColumnSequence.identifier(),
			SessionColumnState.identifier(),
			SessionColumnResourceOwner.identifier(),
			SessionColumnCreator.identifier(),
			SessionColumnUserID.identifier(),
			SessionColumnUserCheckedAt.identifier(),
			LoginNameNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			SessionColumnPasswordCheckedAt.identifier(),
			SessionColumnMetadata.identifier(),
			countColumn.identifier(),
		).From(sessionsTable.identifier()).
			LeftJoin(join(LoginNameUserIDCol, SessionColumnUserID)).
			LeftJoin(join(HumanUserIDCol, SessionColumnUserID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(rows *sql.Rows) (*Sessions, error) {
			sessions := &Sessions{Sessions: []*Session{}}

			for rows.Next() {
				session := new(Session)

				var (
					userID            sql.NullString
					userCheckedAt     sql.NullTime
					loginName         sql.NullString
					displayName       sql.NullString
					passwordCheckedAt sql.NullTime
					metadata          []byte
				)

				err := rows.Scan(
					&session.ID,
					&session.CreationDate,
					&session.ChangeDate,
					&session.Sequence,
					&session.State,
					&session.ResourceOwner,
					&session.Creator,
					&userID,
					&userCheckedAt,
					&loginName,
					&displayName,
					&passwordCheckedAt,
					&metadata,
					&sessions.Count,
				)

				if err != nil {
					return nil, errors.ThrowInternal(err, "QUERY-SAfeg", "Errors.Internal")
				}
				session.UserFactor.UserID = userID.String
				session.UserFactor.UserCheckedAt = userCheckedAt.Time
				session.UserFactor.LoginName = loginName.String
				session.UserFactor.DisplayName = displayName.String
				session.PasswordFactor.PasswordCheckedAt = passwordCheckedAt.Time
				if len(metadata) != 0 {
					if err := json.Unmarshal(metadata, &session.Metadata); err != nil {
						return nil, err
					}
				}

				sessions.Sessions = append(sessions.Sessions, session)
			}

			return sessions, nil
		}
}
