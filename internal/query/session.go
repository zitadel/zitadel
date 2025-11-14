package query

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
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
	IntentFactor   SessionIntentFactor
	WebAuthNFactor SessionWebAuthNFactor
	TOTPFactor     SessionTOTPFactor
	OTPSMSFactor   SessionOTPFactor
	OTPEmailFactor SessionOTPFactor
	Metadata       map[string][]byte
	UserAgent      domain.UserAgent
	Expiration     time.Time
}

type SessionUserFactor struct {
	UserID        string
	ResourceOwner string
	UserCheckedAt time.Time
	LoginName     string
	DisplayName   string
}

type SessionPasswordFactor struct {
	PasswordCheckedAt time.Time
}

type SessionIntentFactor struct {
	IntentCheckedAt time.Time
}

type SessionWebAuthNFactor struct {
	WebAuthNCheckedAt time.Time
	UserVerified      bool
}

type SessionTOTPFactor struct {
	TOTPCheckedAt time.Time
}

type SessionOTPFactor struct {
	OTPCheckedAt time.Time
}

type SessionsSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func sessionsCheckPermission(ctx context.Context, sessions *Sessions, permissionCheck domain.PermissionCheck) {
	sessions.Sessions = slices.DeleteFunc(sessions.Sessions,
		func(session *Session) bool {
			return sessionCheckPermission(ctx, session.ResourceOwner, session.Creator, session.UserAgent, session.UserFactor, permissionCheck) != nil
		},
	)
}

func sessionCheckPermission(ctx context.Context, resourceOwner string, creator string, useragent domain.UserAgent, userFactor SessionUserFactor, permissionCheck domain.PermissionCheck) error {
	data := authz.GetCtxData(ctx)
	// no permission check necessary if user is creator
	if data.UserID == creator {
		return nil
	}
	// no permission check necessary if session belongs to the user
	if userFactor.UserID != "" && data.UserID == userFactor.UserID {
		return nil
	}
	// no permission check necessary if session belongs to the same useragent as used
	if data.AgentID != "" && useragent.FingerprintID != nil && *useragent.FingerprintID != "" && data.AgentID == *useragent.FingerprintID {
		return nil
	}
	// if session belongs to a user, check for permission on the user resource
	if userFactor.ResourceOwner != "" {
		if err := permissionCheck(ctx, domain.PermissionSessionRead, userFactor.ResourceOwner, userFactor.UserID); err != nil {
			return err
		}
		return nil
	}
	// default, check for permission on instance
	return permissionCheck(ctx, domain.PermissionSessionRead, resourceOwner, "")
}

func sessionsPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		SessionColumnResourceOwner,
		domain.PermissionSessionRead,
		// Allow if user is creator
		OwnedRowsPermissionOption(SessionColumnCreator),
		// Allow if session belongs to the user
		OwnedRowsPermissionOption(SessionColumnUserID),
		// Allow if session belongs to the same useragent
		ConnectionPermissionOption(SessionColumnUserAgentFingerprintID, authz.GetCtxData(ctx).AgentID),
	)
	return query.JoinClause(join, args...)
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
	SessionColumnUserResourceOwner = Column{
		name:  projection.SessionColumnUserResourceOwner,
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
	SessionColumnIntentCheckedAt = Column{
		name:  projection.SessionColumnIntentCheckedAt,
		table: sessionsTable,
	}
	SessionColumnWebAuthNCheckedAt = Column{
		name:  projection.SessionColumnWebAuthNCheckedAt,
		table: sessionsTable,
	}
	SessionColumnWebAuthNUserVerified = Column{
		name:  projection.SessionColumnWebAuthNUserVerified,
		table: sessionsTable,
	}
	SessionColumnTOTPCheckedAt = Column{
		name:  projection.SessionColumnTOTPCheckedAt,
		table: sessionsTable,
	}
	SessionColumnOTPSMSCheckedAt = Column{
		name:  projection.SessionColumnOTPSMSCheckedAt,
		table: sessionsTable,
	}
	SessionColumnOTPEmailCheckedAt = Column{
		name:  projection.SessionColumnOTPEmailCheckedAt,
		table: sessionsTable,
	}
	SessionColumnMetadata = Column{
		name:  projection.SessionColumnMetadata,
		table: sessionsTable,
	}
	SessionColumnToken = Column{
		name:  projection.SessionColumnTokenID,
		table: sessionsTable,
	}
	SessionColumnUserAgentFingerprintID = Column{
		name:  projection.SessionColumnUserAgentFingerprintID,
		table: sessionsTable,
	}
	SessionColumnUserAgentIP = Column{
		name:  projection.SessionColumnUserAgentIP,
		table: sessionsTable,
	}
	SessionColumnUserAgentDescription = Column{
		name:  projection.SessionColumnUserAgentDescription,
		table: sessionsTable,
	}
	SessionColumnUserAgentHeader = Column{
		name:  projection.SessionColumnUserAgentHeader,
		table: sessionsTable,
	}
	SessionColumnExpiration = Column{
		name:  projection.SessionColumnExpiration,
		table: sessionsTable,
	}
)

func (q *Queries) SessionByID(ctx context.Context, shouldTriggerBulk bool, id, sessionToken string, permissionCheck domain.PermissionCheck) (session *Session, err error) {
	session, tokenID, err := q.sessionByID(ctx, shouldTriggerBulk, id)
	if err != nil {
		return nil, err
	}
	if sessionToken == "" {
		// for internal calls, no token or permission check is necessary
		if permissionCheck == nil {
			return session, nil
		}
		if err := sessionCheckPermission(ctx, session.ResourceOwner, session.Creator, session.UserAgent, session.UserFactor, permissionCheck); err != nil {
			return nil, err
		}
		return session, nil
	}
	if err := q.sessionTokenVerifier(ctx, sessionToken, session.ID, tokenID); err != nil {
		return nil, zerrors.ThrowPermissionDenied(nil, "QUERY-dsfr3", "Errors.PermissionDenied")
	}
	return session, nil
}

func (q *Queries) sessionByID(ctx context.Context, shouldTriggerBulk bool, id string) (session *Session, tokenID string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerSessionProjection")
		ctx, err = projection.SessionProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("unable to trigger")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareSessionQuery()
	stmt, args, err := query.Where(
		sq.Eq{
			SessionColumnID.identifier():         id,
			SessionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		},
	).ToSql()
	if err != nil {
		return nil, "", zerrors.ThrowInternal(err, "QUERY-dn9JW", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		session, tokenID, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, "", err
	}
	return session, tokenID, nil
}

func (q *Queries) SearchSessions(ctx context.Context, queries *SessionsSearchQueries, permissionCheck domain.PermissionCheck) (*Sessions, error) {
	permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	sessions, err := q.searchSessions(ctx, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && !permissionCheckV2 {
		sessionsCheckPermission(ctx, sessions, permissionCheck)
	}
	return sessions, nil
}

func (q *Queries) searchSessions(ctx context.Context, queries *SessionsSearchQueries, permissionCheckV2 bool) (sessions *Sessions, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSessionsQuery()
	query = sessionsPermissionCheckV2(ctx, query, permissionCheckV2)
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			SessionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).
		ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-sn9Jf", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		sessions, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Sfg42", "Errors.Internal")
	}

	sessions.State, err = q.latestState(ctx, sessionsTable)
	return sessions, err
}

func NewSessionIDsSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(SessionColumnID, list, ListIn)
}

func NewSessionCreatorSearchQuery(creator string) (SearchQuery, error) {
	return NewTextQuery(SessionColumnCreator, creator, TextEquals)
}

func NewSessionUserAgentFingerprintIDSearchQuery(fingerprintID string) (SearchQuery, error) {
	return NewTextQuery(SessionColumnUserAgentFingerprintID, fingerprintID, TextEquals)
}

func NewUserIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(SessionColumnUserID, id, TextEquals)
}

func NewCreationDateQuery(datetime time.Time, compare TimestampComparison) (SearchQuery, error) {
	return NewTimestampQuery(SessionColumnCreationDate, datetime, compare)
}

func NewExpirationDateQuery(datetime time.Time, compare TimestampComparison) (SearchQuery, error) {
	return NewTimestampQuery(SessionColumnExpiration, datetime, compare)
}

func prepareSessionQuery() (sq.SelectBuilder, func(*sql.Row) (*Session, string, error)) {
	return sq.Select(
			SessionColumnID.identifier(),
			SessionColumnCreationDate.identifier(),
			SessionColumnChangeDate.identifier(),
			SessionColumnSequence.identifier(),
			SessionColumnState.identifier(),
			SessionColumnResourceOwner.identifier(),
			SessionColumnCreator.identifier(),
			SessionColumnUserID.identifier(),
			SessionColumnUserResourceOwner.identifier(),
			SessionColumnUserCheckedAt.identifier(),
			LoginNameNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			SessionColumnPasswordCheckedAt.identifier(),
			SessionColumnIntentCheckedAt.identifier(),
			SessionColumnWebAuthNCheckedAt.identifier(),
			SessionColumnWebAuthNUserVerified.identifier(),
			SessionColumnTOTPCheckedAt.identifier(),
			SessionColumnOTPSMSCheckedAt.identifier(),
			SessionColumnOTPEmailCheckedAt.identifier(),
			SessionColumnMetadata.identifier(),
			SessionColumnToken.identifier(),
			SessionColumnUserAgentFingerprintID.identifier(),
			SessionColumnUserAgentIP.identifier(),
			SessionColumnUserAgentDescription.identifier(),
			SessionColumnUserAgentHeader.identifier(),
			SessionColumnExpiration.identifier(),
		).From(sessionsTable.identifier()).
			LeftJoin(join(LoginNameUserIDCol, SessionColumnUserID)).
			LeftJoin(join(HumanUserIDCol, SessionColumnUserID)).
			LeftJoin(join(UserIDCol, SessionColumnUserID)).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*Session, string, error) {
			session := new(Session)

			var (
				userID              sql.NullString
				userResourceOwner   sql.NullString
				userCheckedAt       sql.NullTime
				loginName           sql.NullString
				displayName         sql.NullString
				passwordCheckedAt   sql.NullTime
				intentCheckedAt     sql.NullTime
				webAuthNCheckedAt   sql.NullTime
				webAuthNUserPresent sql.NullBool
				totpCheckedAt       sql.NullTime
				otpSMSCheckedAt     sql.NullTime
				otpEmailCheckedAt   sql.NullTime
				metadata            database.Map[[]byte]
				token               sql.NullString
				userAgentIP         sql.NullString
				userAgentHeader     database.Map[[]string]
				expiration          sql.NullTime
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
				&userResourceOwner,
				&userCheckedAt,
				&loginName,
				&displayName,
				&passwordCheckedAt,
				&intentCheckedAt,
				&webAuthNCheckedAt,
				&webAuthNUserPresent,
				&totpCheckedAt,
				&otpSMSCheckedAt,
				&otpEmailCheckedAt,
				&metadata,
				&token,
				&session.UserAgent.FingerprintID,
				&userAgentIP,
				&session.UserAgent.Description,
				&userAgentHeader,
				&expiration,
			)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, "", zerrors.ThrowNotFound(err, "QUERY-SFeaa", "Errors.Session.NotExisting")
				}
				return nil, "", zerrors.ThrowInternal(err, "QUERY-SAder", "Errors.Internal")
			}

			session.UserFactor.UserID = userID.String
			session.UserFactor.ResourceOwner = userResourceOwner.String
			session.UserFactor.UserCheckedAt = userCheckedAt.Time
			session.UserFactor.LoginName = loginName.String
			session.UserFactor.DisplayName = displayName.String
			session.PasswordFactor.PasswordCheckedAt = passwordCheckedAt.Time
			session.IntentFactor.IntentCheckedAt = intentCheckedAt.Time
			session.WebAuthNFactor.WebAuthNCheckedAt = webAuthNCheckedAt.Time
			session.WebAuthNFactor.UserVerified = webAuthNUserPresent.Bool
			session.TOTPFactor.TOTPCheckedAt = totpCheckedAt.Time
			session.OTPSMSFactor.OTPCheckedAt = otpSMSCheckedAt.Time
			session.OTPEmailFactor.OTPCheckedAt = otpEmailCheckedAt.Time
			session.Metadata = metadata
			session.UserAgent.Header = http.Header(userAgentHeader)
			if userAgentIP.Valid {
				session.UserAgent.IP = net.ParseIP(userAgentIP.String)
			}
			session.Expiration = expiration.Time
			return session, token.String, nil
		}
}

func prepareSessionsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Sessions, error)) {
	return sq.Select(
			SessionColumnID.identifier(),
			SessionColumnCreationDate.identifier(),
			SessionColumnChangeDate.identifier(),
			SessionColumnSequence.identifier(),
			SessionColumnState.identifier(),
			SessionColumnResourceOwner.identifier(),
			SessionColumnCreator.identifier(),
			SessionColumnUserID.identifier(),
			SessionColumnUserResourceOwner.identifier(),
			SessionColumnUserCheckedAt.identifier(),
			LoginNameNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			SessionColumnPasswordCheckedAt.identifier(),
			SessionColumnIntentCheckedAt.identifier(),
			SessionColumnWebAuthNCheckedAt.identifier(),
			SessionColumnWebAuthNUserVerified.identifier(),
			SessionColumnTOTPCheckedAt.identifier(),
			SessionColumnOTPSMSCheckedAt.identifier(),
			SessionColumnOTPEmailCheckedAt.identifier(),
			SessionColumnMetadata.identifier(),
			SessionColumnUserAgentFingerprintID.identifier(),
			SessionColumnUserAgentIP.identifier(),
			SessionColumnUserAgentDescription.identifier(),
			SessionColumnUserAgentHeader.identifier(),
			SessionColumnExpiration.identifier(),
			countColumn.identifier(),
		).From(sessionsTable.identifier()).
			LeftJoin(join(LoginNameUserIDCol, SessionColumnUserID)).
			LeftJoin(join(HumanUserIDCol, SessionColumnUserID)).
			LeftJoin(join(UserIDCol, SessionColumnUserID)).
			PlaceholderFormat(sq.Dollar), func(rows *sql.Rows) (*Sessions, error) {
			sessions := &Sessions{Sessions: []*Session{}}

			for rows.Next() {
				session := new(Session)

				var (
					userID              sql.NullString
					userResourceOwner   sql.NullString
					userCheckedAt       sql.NullTime
					loginName           sql.NullString
					displayName         sql.NullString
					passwordCheckedAt   sql.NullTime
					intentCheckedAt     sql.NullTime
					webAuthNCheckedAt   sql.NullTime
					webAuthNUserPresent sql.NullBool
					totpCheckedAt       sql.NullTime
					otpSMSCheckedAt     sql.NullTime
					otpEmailCheckedAt   sql.NullTime
					metadata            database.Map[[]byte]
					userAgentIP         sql.NullString
					userAgentHeader     database.Map[[]string]
					expiration          sql.NullTime
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
					&userResourceOwner,
					&userCheckedAt,
					&loginName,
					&displayName,
					&passwordCheckedAt,
					&intentCheckedAt,
					&webAuthNCheckedAt,
					&webAuthNUserPresent,
					&totpCheckedAt,
					&otpSMSCheckedAt,
					&otpEmailCheckedAt,
					&metadata,
					&session.UserAgent.FingerprintID,
					&userAgentIP,
					&session.UserAgent.Description,
					&userAgentHeader,
					&expiration,
					&sessions.Count,
				)

				if err != nil {
					return nil, zerrors.ThrowInternal(err, "QUERY-SAfeg", "Errors.Internal")
				}
				session.UserFactor.UserID = userID.String
				session.UserFactor.ResourceOwner = userResourceOwner.String
				session.UserFactor.UserCheckedAt = userCheckedAt.Time
				session.UserFactor.LoginName = loginName.String
				session.UserFactor.DisplayName = displayName.String
				session.PasswordFactor.PasswordCheckedAt = passwordCheckedAt.Time
				session.IntentFactor.IntentCheckedAt = intentCheckedAt.Time
				session.WebAuthNFactor.WebAuthNCheckedAt = webAuthNCheckedAt.Time
				session.WebAuthNFactor.UserVerified = webAuthNUserPresent.Bool
				session.TOTPFactor.TOTPCheckedAt = totpCheckedAt.Time
				session.OTPSMSFactor.OTPCheckedAt = otpSMSCheckedAt.Time
				session.OTPEmailFactor.OTPCheckedAt = otpEmailCheckedAt.Time
				session.Metadata = metadata
				session.UserAgent.Header = http.Header(userAgentHeader)
				if userAgentIP.Valid {
					session.UserAgent.IP = net.ParseIP(userAgentIP.String)
				}
				session.Expiration = expiration.Time

				sessions.Sessions = append(sessions.Sessions, session)
			}

			return sessions, nil
		}
}
