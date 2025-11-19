package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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

const querySessionStmt = `SELECT sessions.instance_id, sessions.id, sessions.token, sessions.lifetime, sessions.expiration, sessions.user_id, sessions.creator_id, sessions.created_at, sessions.updated_at` +
	` , jsonb_agg(json_build_object('instanceId', session_factors.instance_id, 'sessionId', session_factors.session_id, 'type', session_factors.type, 'lastChallengeAt', session_factors.last_challenged_at, 'challenged_payload', session_factors.challenged_payload, 'lastVerifiedAt', session_factors.last_verified_at, 'verified_payload', session_factors.verified_payload)) FILTER (WHERE session_factors.session_id IS NOT NULL) AS factors` +
	` , jsonb_agg(json_build_object('instanceId', session_metadata.instance_id, 'sessionId', session_metadata.session_id, 'key', session_metadata.key, 'value', encode(session_metadata.value, 'base64'), 'createdAt', session_metadata.created_at, 'updatedAt', session_metadata.updated_at)) FILTER (WHERE session_metadata.session_id IS NOT NULL) AS metadata` +
	` FROM zitadel.sessions`

// Get implements [domain.SessionRepository].
func (s session) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Session, error) {
	opts = append(opts,
		s.joinFactors(),
		s.joinMetadata(),
		database.WithGroupBy(s.InstanceIDColumn(), s.IDColumn()),
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
		s.joinUserAgent(),
		s.joinFactors(),
		s.joinMetadata(),
		database.WithGroupBy(s.InstanceIDColumn(), s.IDColumn()),
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

	return scanSessions(ctx, client, &builder)
}

const upsertSessionUserAgentStmt = `WITH user_agent AS (
	INSERT INTO zitadel.session_user_agents(
		instance_id, fingerprint_id, description, ip, headers
	)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (instance_id, fingerprint_id)
	DO UPDATE SET description = EXCLUDED.description, ip = EXCLUDED.ip, headers = excluded.headers
) `

// Create implements [domain.SessionRepository].
func (s session) Create(ctx context.Context, client database.QueryExecutor, session *domain.Session) error {
	var (
		createdAt, updatedAt any = database.DefaultInstruction, database.DefaultInstruction
	)
	if !session.CreatedAt.IsZero() {
		createdAt = session.CreatedAt
	}
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
	builder.WriteString(`INSERT INTO ` + s.qualifiedTableName() + ` (instance_id, id, creator_id, user_agent_id, created_at, updated_at) VALUES ( `)
	builder.WriteArgs(session.InstanceID, session.ID, session.CreatorID, fingerprintID, createdAt, updatedAt)
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
		changes = append(changes, database.NewChange(s.UpdatedAtColumn(), database.NullInstruction))
	}

	var builder database.StatementBuilder
	builder.WriteString("WITH existing_session AS (SELECT * FROM zitadel.sessions ")
	writeCondition(&builder, condition)
	builder.WriteString(") ")
	for i, change := range changes {
		if multi, ok := change.(database.Changes); ok {
			for j, ch := range multi {
				sessionCTE(ch, i, j, &builder)
			}
		}
		sessionCTE(change, i, 0, &builder)
	}
	builder.WriteString("UPDATE zitadel.sessions SET ")
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func sessionCTE(change database.Change, i, j int, builder *database.StatementBuilder) {
	if cte, ok := change.(database.CTEChange); ok {
		name := fmt.Sprintf("cte_%d_%d", i, j)
		builder.WriteString(fmt.Sprintf(", %s as (", name))
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
	return database.NewChange(s.TokenColumn(), token)
}

// SetLifetime implements [domain.sessionChanges].
func (s session) SetLifetime(lifetime time.Duration) database.Change {
	return database.NewChange(s.LifetimeColumn(), lifetime)
}

// SetChallenge implements [domain.sessionChanges].
func (s session) SetChallenge(challenge domain.SessionChallenge) database.Change {
	switch c := challenge.(type) {
	case *domain.SessionChallengePasskey:
		payload, err := json.Marshal(c)
		_ = err // TODO: handle error
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_challenged_at, challenged_payload) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypePasskey, c.LastChallengedAt, payload) //TODO: json?
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_challenged_at = EXCLUDED.last_challenged_at, challenged_payload = EXCLUDED.challenged_payload")
			}, nil,
		)
	case *domain.SessionChallengeOTPSMS:
		payload, err := json.Marshal(c)
		_ = err // TODO: handle error
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_challenged_at, challenged_payload) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeOTPSMS, c.LastChallengedAt, payload) //TODO: json?
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_challenged_at = EXCLUDED.last_challenged_at, challenged_payload = EXCLUDED.challenged_payload")
			}, nil,
		)
	case *domain.SessionChallengeOTPEmail:
		payload, err := json.Marshal(c)
		_ = err // TODO: handle error
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_challenged_at, challenged_payload) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeOTPEmail, c.LastChallengedAt, payload) //TODO: json?
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_challenged_at = EXCLUDED.last_challenged_at, challenged_payload = EXCLUDED.challenged_payload")
			}, nil,
		)
	}
	return nil //TODO: error?
}

// SetFactor implements [domain.sessionChanges].
func (s session) SetFactor(factor domain.SessionFactor) database.Change {
	switch f := factor.(type) {
	case *domain.SessionFactorUser:
		return database.NewChanges(
			database.NewChange(s.UserIDColumn(), f.UserID),
			database.NewCTEChange(
				func(builder *database.StatementBuilder) {
					builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
					builder.WriteArgs(domain.SessionFactorTypeUser.String(), f.LastVerifiedAt)
					builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
				}, nil,
			),
		)
	case *domain.SessionFactorPassword:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypePassword.String(), f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	case *domain.SessionFactorIdentityProviderIntent:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeIdentityProviderIntent.String(), f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	case *domain.SessionFactorPasskey:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at, verified_payload) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypePasskey.String(), f.LastVerifiedAt, f.UserVerified)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at, verified_payload = EXCLUDED.verified_payload")
			}, nil,
		)
	case *domain.SessionFactorTOTP:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeTOTP.String(), f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	case *domain.SessionFactorOTPSMS:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeOTPSMS.String(), f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	case *domain.SessionFactorOTPEmail:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, type, last_verified_at) SELECT instance_id, id, ")
				builder.WriteArgs(domain.SessionFactorTypeOTPEmail.String(), f.LastVerifiedAt)
				builder.WriteString(" FROM existing_session ON CONFLICT (instance_id, session_id, type) DO UPDATE SET last_verified_at = EXCLUDED.last_verified_at")
			}, nil,
		)
	default:
		return nil //TODO: error?
	}
}

// ClearFactor implements [domain.sessionChanges].
func (s session) ClearFactor(factorType domain.SessionFactorType) database.Change {
	return database.NewChange(s.factorRepo.LastVerifiedAtColumn(), database.NullInstruction)
}

// SetMetadata implements [domain.sessionChanges].
func (s session) SetMetadata(metadata []domain.SessionMetadata) database.Change {
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
			builder.WriteString("DELETE FROM zitadel.session_metadata WHERE instance_id = (SELECT instance_id FROM existing_session) AND session_id = (SELECT id FROM existing_session) AND key NOT IN (")
			builder.WriteArgs(keys...)
			builder.WriteString(")")
		}, nil)
	return database.NewChanges(changes...)
}

// SetUserAgent implements [domain.sessionChanges].
func (s session) SetUserAgent(userAgent domain.SessionUserAgent) database.Change {
	//TODO: upsert user agent?
	return database.NewChanges(
		database.NewChange(s.UserAgentIDColumn(), *userAgent.FingerprintID),
	)
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
	return database.NewTextCondition(s.UserAgentIDColumn(), database.TextOperationEqual, userAgentID)
}

// UserIDCondition implements [domain.sessionConditions].
func (s session) UserIDCondition(userID string) database.Condition {
	return database.NewTextCondition(s.UserIDColumn(), database.TextOperationEqual, userID)
}

// CreatorIDCondition implements [domain.sessionConditions].
func (s session) CreatorIDCondition(creatorID string) database.Condition {
	return database.NewTextCondition(s.CreatorIDColumn(), database.TextOperationEqual, creatorID)
}

// CreatedAtCondition implements [domain.sessionConditions].
func (s session) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(s.CreatedAtColumn(), op, createdAt)
}

// UpdatedAtCondition implements [domain.sessionConditions].
func (s session) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(s.UpdatedAtColumn(), op, updatedAt)
}

func (s session) ExistsFactor(cond database.Condition) database.Condition {
	return database.Exists(
		s.factorRepo.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(s.InstanceIDColumn(), s.factorRepo.InstanceIDColumn()),
			database.NewColumnCondition(s.IDColumn(), s.factorRepo.SessionIDColumn()),
			cond,
		),
	)
}

func (s session) FactorConditions() domain.SessionFactorConditions {
	return s.factorRepo
}

func (s session) ExistsMetadata(cond database.Condition) database.Condition {
	return database.Exists(
		s.metadataRepo.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(s.InstanceIDColumn(), s.metadataRepo.InstanceIDColumn()),
			database.NewColumnCondition(s.IDColumn(), s.metadataRepo.SessionIDColumn()),
			cond,
		),
	)
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

// TokenColumn implements [domain.sessionColumns].
func (s session) TokenColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "token")
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
func (s session) UserAgentIDColumn() database.Column {
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
	Token      sql.NullString                    `db:"token"`
	Lifetime   sql.Null[time.Duration]           `db:"lifetime"`
	Expiration sql.NullTime                      `db:"expiration"`
	UserID     sql.NullString                    `db:"user_id"`
	CreatorID  sql.NullString                    `db:"creator_id"`
	Factors    JSONArray[rawFactor]              `json:"factorRepo,omitempty" db:"factors"`
	Metadata   JSONArray[domain.SessionMetadata] `json:"metadata,omitempty" db:"metadata"`
	UserAgent  rawUserAgent                      `json:"userAgent,omitempty" db:"userAgentRepo"`
}

type rawFactor struct {
	Type              string    `db:"type"`
	LastChallengedAt  time.Time `db:"last_challenged_at"`
	ChallengedPayload []byte    `db:"challenged_payload"`
	LastVerifiedAt    time.Time `db:"last_verified_at"`
	VerifiedPayload   []byte    `db:"verified_payload"`
}

func (f *rawFactor) ToDomain() (domain.SessionFactor, error) {
	factorType, err := domain.SessionFactorTypeString(f.Type)
	if err != nil {
		return nil, err
	}
	switch factorType {
	case domain.SessionFactorTypeUser:
		return &domain.SessionFactorUser{
			UserID:         "",
			LastVerifiedAt: f.LastVerifiedAt,
		}, nil
	case domain.SessionFactorTypePassword:
		return &domain.SessionFactorPassword{
			LastVerifiedAt: f.LastVerifiedAt,
			//LastFailedAt:   f.,
		}, nil
	case domain.SessionFactorTypePasskey:
		passkey := new(domain.SessionFactorPasskey)
		json.Unmarshal(f.VerifiedPayload, passkey)
		return &domain.SessionFactorPasskey{
			LastVerifiedAt: f.LastVerifiedAt,
			UserVerified:   passkey.UserVerified,
		}, nil
	case domain.SessionFactorTypeIdentityProviderIntent:
		return &domain.SessionFactorIdentityProviderIntent{
			LastVerifiedAt: f.LastVerifiedAt,
		}, nil
	case domain.SessionFactorTypeTOTP:
		return &domain.SessionFactorTOTP{
			LastVerifiedAt: f.LastVerifiedAt,
			//LastFailedAt:   time.Time{},
		}, nil
	case domain.SessionFactorTypeOTPSMS:
		return &domain.SessionFactorOTPSMS{
			LastVerifiedAt: f.LastVerifiedAt,
			//LastFailedAt:   time.Time{},
		}, nil
	case domain.SessionFactorTypeOTPEmail:
		return &domain.SessionFactorOTPEmail{
			LastVerifiedAt: f.LastVerifiedAt,
			//LastFailedAt:   time.Time{},
		}, nil
	}
	return nil, nil // TODO: !
}

type rawUserAgent struct {
	FingerprintID sql.NullString `json:"fingerprintId,omitempty" db:"fingerprint_id"`
	Description   sql.NullString `json:"description,omitempty" db:"description"`
	//IP            net.IP         `json:"ip,omitempty" db:"ip"`
	//Header        http.Header    `json:"header,omitempty" db:"headers"`
}

func scanSession(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Session, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var raw rawSession
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(&raw); err != nil {
		return nil, err
	}

	raw.Session.Token = raw.Token.String
	raw.Session.Lifetime = raw.Lifetime.V
	raw.Session.Expiration = raw.Expiration.Time
	raw.Session.UserID = raw.UserID.String
	raw.Session.CreatorID = raw.CreatorID.String

	raw.Session.Factors = make([]domain.SessionFactor, len(raw.Factors))
	for i, factor := range raw.Factors {
		raw.Session.Factors[i], err = factor.ToDomain()
		_ = err //TODO: ?
	}
	raw.Session.Metadata = make([]domain.SessionMetadata, len(raw.Metadata))
	for i, metadata := range raw.Metadata {
		raw.Session.Metadata[i] = *metadata
	}
	raw.Session.UserAgent = &domain.SessionUserAgent{
		FingerprintID: &raw.UserAgent.FingerprintID.String,
		Description:   &raw.UserAgent.Description.String,
	}
	return raw.Session, nil
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
		result[i] = session.Session
		//for i, factor := range session.Factors {
		//	session.Session.Factors[i] = *factor
		//}
		for i, metadata := range session.Metadata {
			session.Session.Metadata[i] = *metadata
		}
	}

	return result, nil
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

func (s session) joinUserAgent() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(s.InstanceIDColumn(), s.userAgentRepo.InstanceIDColumn()),
		database.NewColumnCondition(s.IDColumn(), s.userAgentRepo.FingerprintIDColumn()),
	)
	//
	//// If domains should not be joined, we make sure to return null for the domain columns
	//// the query optimizer of the dialect should optimize this away if no domains are requested
	//if !s.shouldLoadDomains {
	//	columns = append(columns, database.IsNull(s.factorRepo.SessionIDColumn()))
	//}

	return database.WithLeftJoin(
		s.userAgentRepo.qualifiedTableName(),
		database.And(columns...),
	)
}

//func (s session) LoadDomains() domain.OrganizationRepository {
//	return &session{
//		shouldLoadDomains:  true,
//		shouldLoadMetadata: s.shouldLoadMetadata,
//	}
//}

func (s session) joinFactors() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(s.InstanceIDColumn(), s.factorRepo.InstanceIDColumn()),
		database.NewColumnCondition(s.IDColumn(), s.factorRepo.SessionIDColumn()),
	)
	//
	//// If domains should not be joined, we make sure to return null for the domain columns
	//// the query optimizer of the dialect should optimize this away if no domains are requested
	//if !s.shouldLoadDomains {
	//	columns = append(columns, database.IsNull(s.factorRepo.SessionIDColumn()))
	//}

	return database.WithLeftJoin(
		s.factorRepo.qualifiedTableName(),
		database.And(columns...),
	)
}

//
//func (s session) LoadMetadata() domain.OrganizationRepository {
//	return &session{
//		shouldLoadDomains:  s.shouldLoadDomains,
//		shouldLoadMetadata: true,
//	}
//}

func (s session) joinMetadata() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(s.InstanceIDColumn(), s.metadataRepo.InstanceIDColumn()),
		database.NewColumnCondition(s.IDColumn(), s.metadataRepo.SessionIDColumn()),
	)
	//
	//// If metadata should not be joined, we make sure to return null for the metadata columns
	//// the query optimizer of the dialect should optimize this away if no metadata are requested
	//if !o.shouldLoadMetadata {
	//	columns = append(columns, database.IsNull(s.metadataRepo.SessionIDColumn()))
	//}

	return database.WithLeftJoin(
		s.metadataRepo.qualifiedTableName(),
		database.And(columns...),
	)
}
