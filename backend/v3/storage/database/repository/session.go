package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

var _ domain.SessionRepository = (*session)(nil)

type session struct {
	factorRepo    sessionFactor
	metadataRepo  sessionMetadata
	userAgentRepo sessionUserAgent
}

func (s session) qualifiedTableName() string {
	return "zitadel.sessions"
}

func (s session) unqualifiedTableName() string {
	return "sessions"
}

func SessionRepository() domain.SessionRepository {
	return new(session)
}

const querySessionStmt = `
SELECT 
	sessions.instance_id
	, sessions.id
	, sessions.token_id
	, sessions.lifetime
	, sessions.expiration
	, sessions.user_id
	, sessions.creator_id
	, sessions.created_at
	, sessions.updated_at
	, jsonb_agg(distinct 
		jsonb_build_object(
			'instanceId', session_factors.instance_id
			, 'sessionId', session_factors.session_id
			, 'type', session_factors.type
			, 'lastChallengedAt', session_factors.last_challenged_at
			, 'challengedPayload', session_factors.challenged_payload
			, 'lastVerifiedAt', session_factors.last_verified_at
			, 'verifiedPayload', session_factors.verified_payload
		)
	)
	FILTER (WHERE session_factors.session_id IS NOT NULL) AS factors
	, jsonb_agg(distinct
		jsonb_build_object(
			'instanceId', session_metadata.instance_id
			, 'sessionId', session_metadata.session_id
			, 'key', session_metadata.key
			, 'value', encode(session_metadata.value, 'base64')
			, 'createdAt', session_metadata.created_at
			, 'updatedAt', session_metadata.updated_at
		)
	)
	FILTER (WHERE session_metadata.session_id IS NOT NULL) AS metadata
	, CASE WHEN session_user_agents.fingerprint_id IS NOT NULL THEN
		jsonb_build_object(
			'fingerprintId', session_user_agents.fingerprint_id
			, 'description', session_user_agents.description
			, 'ip', session_user_agents.ip
			, 'headers', session_user_agents.headers
		)
	END as user_agent
	FROM zitadel.sessions`

// Get implements [domain.SessionRepository].
func (s session) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Session, error) {
	opts = append(opts,
		s.joinFactors(),
		s.joinMetadata(),
		s.joinUserAgent(),
		database.WithGroupBy(
			s.InstanceIDColumn(),
			s.IDColumn(),
			s.userAgentRepo.fingerprintIDColumn(),
			s.userAgentRepo.descriptionColumn(),
			s.userAgentRepo.ipColumn(),
			s.userAgentRepo.headersColumn(),
		),
	)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(s.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(s.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(querySessionStmt)
	options.Write(&builder)

	return scanSession(ctx, client, &builder)
}

// List implements [domain.SessionRepository].
func (s session) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Session, error) {
	opts = append(opts,
		s.joinFactors(),
		s.joinMetadata(),
		s.joinUserAgent(),
		database.WithGroupBy(
			s.InstanceIDColumn(),
			s.IDColumn(),
			s.userAgentRepo.fingerprintIDColumn(),
			s.userAgentRepo.descriptionColumn(),
			s.userAgentRepo.ipColumn(),
			s.userAgentRepo.headersColumn(),
		),
	)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if options.Condition == nil || !options.Condition.IsRestrictingColumn(s.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(s.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(querySessionStmt)
	options.Write(&builder)

	return scanSessions(ctx, client, &builder)
}

const upsertSessionUserAgentStmt = `
WITH user_agent AS (
	INSERT INTO zitadel.session_user_agents(
		instance_id, fingerprint_id, description, ip, headers
	)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (instance_id, fingerprint_id)
	DO UPDATE SET description = EXCLUDED.description, ip = EXCLUDED.ip, headers = excluded.headers
) `

// Create implements [domain.SessionRepository].
func (s session) Create(ctx context.Context, client database.QueryExecutor, session *domain.Session) error {
	var createdAt any = database.DefaultInstruction
	if !session.CreatedAt.IsZero() {
		createdAt = session.CreatedAt
	}
	var updatedAt = createdAt
	if !session.UpdatedAt.IsZero() {
		updatedAt = session.UpdatedAt
	}

	builder := new(database.StatementBuilder)
	var fingerprintID *string
	if session.UserAgent != nil {
		builder = database.NewStatementBuilder(upsertSessionUserAgentStmt,
			session.InstanceID,
			session.UserAgent.FingerprintID,
			session.UserAgent.Description,
			session.UserAgent.IP,
			session.UserAgent.Header,
		)
		fingerprintID = session.UserAgent.FingerprintID
	}
	builder.WriteString(`INSERT INTO `)
	builder.WriteString(s.qualifiedTableName())
	builder.WriteString(` (instance_id, id, lifetime, creator_id, user_agent_id, created_at, updated_at) VALUES ( `)
	builder.WriteArgs(session.InstanceID, session.ID, session.Lifetime, session.CreatorID, fingerprintID, createdAt, updatedAt)
	builder.WriteString(` ) RETURNING created_at, updated_at`)
	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&session.CreatedAt, &session.UpdatedAt)
}

// Update implements [domain.SessionRepository].
func (s session) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if !condition.IsRestrictingColumn(s.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(s.InstanceIDColumn())
	}
	if !database.Changes(changes).IsOnColumn(s.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(s.UpdatedAtColumn(), database.DefaultInstruction))
	}

	var builder database.StatementBuilder
	builder.WriteString("WITH existing_session AS (SELECT * FROM zitadel.sessions ")
	writeCondition(&builder, condition)
	builder.WriteString(") ")
	for i, change := range changes {
		sessionCTE(change, i, 0, &builder)
	}
	builder.WriteString("UPDATE zitadel.sessions SET ")
	if err := database.Changes(changes).Write(&builder); err != nil {
		return 0, err
	}
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func sessionCTE(change database.Change, i, j int, builder *database.StatementBuilder) {
	if multi, ok := change.(database.Changes); ok {
		for j, ch := range multi {
			sessionCTE(ch, i, j, builder)
		}
	}
	if cte, ok := change.(database.CTEChange); ok {
		name := fmt.Sprintf("cte_%d_%d", i, j)
		fmt.Fprintf(builder, ", %s as (", name)
		cte.SetName(name)
		cte.WriteCTE(builder)
		builder.WriteString(") ")
	}
}

// Delete implements [domain.SessionRepository].
func (s session) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if !condition.IsRestrictingColumn(s.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(s.InstanceIDColumn())
	}
	var builder database.StatementBuilder
	builder.WriteString("DELETE FROM zitadel.sessions")
	writeCondition(&builder, condition)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetUpdatedAt implements [domain.sessionChanges].
func (s session) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(s.UpdatedAtColumn(), updatedAt)
}

// SetToken implements [domain.sessionChanges].
func (s session) SetToken(token string) database.Change {
	return database.NewChange(s.TokenIDColumn(), token)
}

// SetLifetime implements [domain.sessionChanges].
func (s session) SetLifetime(lifetime time.Duration) database.Change {
	return database.NewChange(s.LifetimeColumn(), lifetime)
}

// SetChallenge implements [domain.sessionChanges].
func (s session) SetChallenge(challenge domain.SessionChallenge) database.Change {
	switch c := challenge.(type) {
	case *domain.SessionChallengePasskey:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_challenged_at, challenged_payload) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypePasskey, c.LastChallengedAt, c)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_challenged_at = EXCLUDED.last_challenged_at, challenged_payload = EXCLUDED.challenged_payload")
			}, nil,
		)
	case *domain.SessionChallengeOTPSMS:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_challenged_at, challenged_payload) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeOTPSMS, c.LastChallengedAt, c)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_challenged_at = EXCLUDED.last_challenged_at, challenged_payload = EXCLUDED.challenged_payload")
			}, nil,
		)
	case *domain.SessionChallengeOTPEmail:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_challenged_at, challenged_payload) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeOTPEmail, c.LastChallengedAt, c)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_challenged_at = EXCLUDED.last_challenged_at, challenged_payload = EXCLUDED.challenged_payload")
			}, nil,
		)
	}
	return nil
}

// SetFactor implements [domain.sessionChanges].
func (s session) SetFactor(factor domain.SessionFactor) database.Change {
	switch f := factor.(type) {
	case *domain.SessionFactorUser:
		return database.NewChanges(
			database.NewChange(s.UserIDColumn(), f.UserID),
			database.NewCTEChange(
				func(builder *database.StatementBuilder) {
					builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at, verified_payload) SELECT instance_id, id, ")
					builder.WriteArgs(domain.SessionFactorTypeUser, f.LastVerifiedAt, f)
					builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at, verified_payload = EXCLUDED.verified_payload")
				}, nil,
			),
		)
	case *domain.SessionFactorPassword:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypePassword, f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	case *domain.SessionFactorIdentityProviderIntent:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeIdentityProviderIntent, f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	case *domain.SessionFactorPasskey:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at, verified_payload) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypePasskey, f.LastVerifiedAt, f)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at, verified_payload = EXCLUDED.verified_payload")
			}, nil,
		)
	case *domain.SessionFactorTOTP:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeTOTP, f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	case *domain.SessionFactorOTPSMS:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeOTPSMS, f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	case *domain.SessionFactorOTPEmail:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeOTPEmail, f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	case *domain.SessionFactorRecoveryCode:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeRecoveryCode, f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	default:
		return nil
	}
}

// ClearFactor implements [domain.sessionChanges].
func (s session) ClearFactor(factorType domain.SessionFactorType) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE zitadel.session_factors sf SET last_verified_at = NULL FROM existing_session es WHERE sf.instance_id = es.instance_id AND sf.session_id = es.id AND sf.type = ")
			builder.WriteArg(factorType)
		}, nil)
}

// SetMetadata implements [domain.sessionChanges].
func (s session) SetMetadata(metadata []*domain.SessionMetadata) database.Change {
	changes := make([]database.Change, len(metadata)+1)
	keys := make([]any, len(metadata))
	for i, md := range metadata {
		keys[i] = md.Key
		changes[i] = database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_metadata (instance_id, session_id, key, value) SELECT instance_id, id, ")
				builder.WriteArgs(md.Key, md.Value)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, key) DO UPDATE SET value = EXCLUDED.value")
			}, nil,
		)
	}
	changes[len(metadata)] = database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM zitadel.session_metadata WHERE instance_id = (SELECT instance_id FROM existing_session) AND session_id = (SELECT id from existing_session) AND key NOT IN (")
			builder.WriteArgs(keys...)
			builder.WriteString(")")
		}, nil)
	return database.NewChanges(changes...)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PrimaryKeyCondition implements [domain.sessionConditions].
func (s session) PrimaryKeyCondition(instanceID, sessionID string) database.Condition {
	return database.And(
		s.InstanceIDCondition(instanceID),
		s.IDCondition(sessionID),
	)
}

// InstanceIDCondition implements [domain.sessionConditions].
func (s session) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(s.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// IDCondition implements [domain.sessionConditions].
func (s session) IDCondition(sessionID string) database.Condition {
	return database.NewTextCondition(s.IDColumn(), database.TextOperationEqual, sessionID)
}

// UserAgentIDCondition implements [domain.sessionConditions].
func (s session) UserAgentIDCondition(userAgentID string) database.Condition {
	return database.NewTextCondition(s.userAgentIDColumn(), database.TextOperationEqual, userAgentID)
}

// UserIDCondition implements [domain.sessionConditions].
func (s session) UserIDCondition(userID string) database.Condition {
	return database.NewTextCondition(s.UserIDColumn(), database.TextOperationEqual, userID)
}

// CreatorIDCondition implements [domain.sessionConditions].
func (s session) CreatorIDCondition(creatorID string) database.Condition {
	return database.NewTextCondition(s.CreatorIDColumn(), database.TextOperationEqual, creatorID)
}

// ExpirationCondition implements [domain.sessionConditions].
func (s session) ExpirationCondition(op database.NumberOperation, expiration time.Time) database.Condition {
	return database.NewNumberCondition(s.ExpirationColumn(), op, expiration)
}

// CreatedAtCondition implements [domain.sessionConditions].
func (s session) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(s.CreatedAtColumn(), op, createdAt)
}

// UpdatedAtCondition implements [domain.sessionConditions].
func (s session) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(s.UpdatedAtColumn(), op, updatedAt)
}

// ExistsFactor implements [domain.sessionConditions].
func (s session) ExistsFactor(cond database.Condition) database.Condition {
	return database.Exists(
		s.factorRepo.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(s.InstanceIDColumn(), s.factorRepo.instanceIDColumn()),
			database.NewColumnCondition(s.IDColumn(), s.factorRepo.sessionIDColumn()),
			cond,
		),
	)
}

// FactorConditions implements [domain.sessionConditions].
func (s session) FactorConditions() domain.SessionFactorConditions {
	return s.factorRepo
}

// ExistsMetadata implements [domain.sessionConditions].
func (s session) ExistsMetadata(cond database.Condition) database.Condition {
	return database.Exists(
		s.metadataRepo.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(s.InstanceIDColumn(), s.metadataRepo.instanceIDColumn()),
			database.NewColumnCondition(s.IDColumn(), s.metadataRepo.sessionIDColumn()),
			cond,
		),
	)
}

// MetadataConditions implements [domain.sessionConditions].
func (s session) MetadataConditions() domain.SessionMetadataConditions {
	return s.metadataRepo
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.Repository].
func (s session) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		s.InstanceIDColumn(),
		s.IDColumn(),
	}
}

// InstanceIDColumn implements [domain.sessionColumns].
func (s session) InstanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

// IDColumn implements [domain.sessionColumns].
func (s session) IDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "id")
}

// TokenIDColumn implements [domain.sessionColumns].
func (s session) TokenIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "token_id")
}

// LifetimeColumn implements [domain.sessionColumns].
func (s session) LifetimeColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "lifetime")
}

// ExpirationColumn implements [domain.sessionColumns].
func (s session) ExpirationColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "expiration")
}

// UserIDColumn implements [domain.sessionColumns].
func (s session) UserIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "user_id")
}

// UserAgentIDColumn implements [domain.sessionColumns].
func (s session) userAgentIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "user_agent_id")
}

// CreatorIDColumn implements [domain.sessionColumns].
func (s session) CreatorIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "creator_id")
}

// CreatedAtColumn implements [domain.sessionColumns].
func (s session) CreatedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "created_at")
}

// UpdatedAtColumn implements [domain.sessionColumns].
func (s session) UpdatedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "updated_at")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

type rawSession struct {
	*domain.Session
	TokenID    *string                           `json:"tokenID" db:"token_id"`
	Lifetime   *time.Duration                    `json:"lifetime" db:"lifetime"`
	Expiration *time.Time                        `json:"expiration" db:"expiration"`
	UserID     *string                           `json:"userID" db:"user_id"`
	CreatorID  *string                           `json:"creatorID" db:"creator_id"`
	Factors    JSONArray[rawFactor]              `json:"factors,omitempty" db:"factors"`
	Metadata   JSONArray[domain.SessionMetadata] `json:"metadata,omitempty" db:"metadata"`
}

func scanSession(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Session, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	raw := new(rawSession)
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(raw); err != nil {
		return nil, err
	}
	return rawSessionToDomain(raw)
}

func scanSessions(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Session, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var sessions []*rawSession
	if err := rows.(database.CollectableRows).Collect(&sessions); err != nil {
		return nil, err
	}

	result := make([]*domain.Session, len(sessions))
	for i, session := range sessions {
		result[i], err = rawSessionToDomain(session)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func rawSessionToDomain(raw *rawSession) (*domain.Session, error) {
	raw.Session.TokenID = gu.Value(raw.TokenID)
	raw.Session.Lifetime = gu.Value(raw.Lifetime)
	raw.Session.Expiration = gu.Value(raw.Expiration)
	raw.Session.UserID = gu.Value(raw.UserID)
	raw.Session.CreatorID = gu.Value(raw.CreatorID)

	for _, factor := range raw.Factors {
		f, ch, err := factor.ToDomain()
		if err != nil {
			return nil, err
		}
		if f != nil {
			raw.Session.Factors.AppendTo(f)
		}
		if ch != nil {
			raw.Challenges.AppendTo(ch)
		}
	}
	raw.Session.Metadata = raw.Metadata
	return raw.Session, nil
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

func (s session) joinUserAgent() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(s.InstanceIDColumn(), s.userAgentRepo.instanceIDColumn()),
		database.NewColumnCondition(s.userAgentIDColumn(), s.userAgentRepo.fingerprintIDColumn()),
	)
	return database.WithLeftJoin(
		s.userAgentRepo.qualifiedTableName(),
		database.And(columns...),
	)
}

func (s session) joinFactors() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(s.InstanceIDColumn(), s.factorRepo.instanceIDColumn()),
		database.NewColumnCondition(s.IDColumn(), s.factorRepo.sessionIDColumn()),
	)
	return database.WithLeftJoin(
		s.factorRepo.qualifiedTableName(),
		database.And(columns...),
	)
}

func (s session) joinMetadata() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(s.InstanceIDColumn(), s.metadataRepo.instanceIDColumn()),
		database.NewColumnCondition(s.IDColumn(), s.metadataRepo.sessionIDColumn()),
	)
	return database.WithLeftJoin(
		s.metadataRepo.qualifiedTableName(),
		database.And(columns...),
	)
}
