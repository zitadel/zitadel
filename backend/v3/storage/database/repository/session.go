package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

var _ domain.SessionRepository = (*session)(nil)

type session struct {
	factorRepo   sessionFactor
	metadataRepo sessionMetadata
	userAgentRepo sessionUserAgent
}

func (s session) unqualifiedTableName() string {
	return "sessions"
}

func SessionRepository() domain.SessionRepository {
	return new(session)
}

const querySessionStmt = ""

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
	builder.WriteString(queryOrganizationStmt)
	options.Write(&builder)

	return scanSession(ctx, client, &builder)
}

// List implements [domain.SessionRepository].
func (s session) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Session, error) {
}

const createSessionStmt = "INSERT INTO zitadel.sessions(instance_id, id, creator_id, fingerprint_id) VALUES ($1, $2, $3, $4)"
const upsertSessionUserAgentStmt = "INSERT INTO zitadel.session_user_agents(instancce_id, fingerprint_id, description, ip, headers) VALUES ($1, $2, $3, $4) ON CONFLICT (instance_id, fingerprint_id) DO UPDATE SET description = EXCLUDED.description, ip = EXCLUDED.ip, headers = excluded.headers"

// Create implements [domain.SessionRepository].
func (s session) Create(ctx context.Context, client database.QueryExecutor, session *domain.Session) error {
	builder := database.StatementBuilder{}
	var fingerprintID *string
	if session.UserAgent != nil {
		builder.AppendArgs(
			session.UserAgent.InstanceID,
			session.UserAgent.FingerprintID,
			session.UserAgent.Description,
			session.UserAgent.IP,
			session.UserAgent.Header,
		)
		builder.WriteString(upsertSessionUserAgentStmt)
		_, err := client.Exec(ctx, builder.String(), builder.Args()...)
		if err != nil {
			return err
		}
		fingerprintID = session.UserAgent.FingerprintID
	}
	builder.AppendArgs(session.InstanceID, session.ID, session.CreatorID, fingerprintID)
	builder.WriteString(createSessionStmt)
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

	factorChanges := make([]database.Change, 0, len(changes))
	var builder database.StatementBuilder
	builder.WriteString("UPDATE zitadel.sessions SET ")
	for _, change := range changes {
		if isFactorChange(change, s.factorRepo) {
			factorChanges = append(factorChanges, change)
		}
	}
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, condition)

	if len(factorChanges) > 0 {
		builder.WriteString("INSERT INTO zitadel.session_factors (instance_id, session_id, factor_type, last_verified_at, last_challenged_at, payload) VALUES ")
	}

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func isFactorChange(change database.Change, factorRepo sessionFactor) bool {
	return change.IsOnColumn(factorRepo.FactorTypeColumn()) ||
		change.IsOnColumn(factorRepo.LastVerifiedAtColumn()) ||
		change.IsOnColumn(factorRepo.LastChallengedAtColumn()) ||
		change.IsOnColumn(factorRepo.PayloadColumn())
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
		return database.NewChanges(
			database.NewChange(s.factorRepo.FactorTypeColumn(), domain.SessionFactorTypePasskey),
			database.NewChange(s.factorRepo.LastChallengedAtColumn(), c.LastChallengedAt),
			database.NewChange(s.factorRepo.PayloadColumn(), payload),
		)
	case *domain.SessionChallengeOTPSMS:
		payload, err := json.Marshal(c)
		return database.NewChanges(
			database.NewChange(s.factorRepo.FactorTypeColumn(), domain.SessionFactorTypeOTPSMS),
			database.NewChange(s.factorRepo.LastChallengedAtColumn(), c.LastChallengedAt),
			database.NewChange(s.factorRepo.PayloadColumn(), payload),
	case *domain.SessionChallengeOTPEmail:
		payload, err := json.Marshal(c)
		return database.NewChanges(
			database.NewChange(s.factorRepo.FactorTypeColumn(), domain.SessionFactorTypeOTPEmail),
			database.NewChange(s.factorRepo.LastChallengedAtColumn(), c.LastChallengedAt),
			database.NewChange(s.factorRepo.PayloadColumn(), payload),
	}
}

// SetFactor implements [domain.sessionChanges].
func (s session) SetFactor(factor domain.SessionFactor) database.Change {
	switch f := factor.(type) {
	case *domain.SessionFactorUser:
		return database.NewChanges(
			database.NewChange(s.UserIDColumn(), f.UserID),
			database.NewChange(s.factorRepo.InstanceIDColumn(), f.InstanceID),
			database.NewChange(s.factorRepo.SessionIDColumn(), f.SessionID),
			database.NewChange(s.factorRepo.FactorTypeColumn(), domain.SessionFactorTypeUser),
			database.NewChange(s.factorRepo.LastVerifiedAtColumn(), f.LastVerifiedAt),
		)
	case *domain.SessionFactorPassword:
		return database.NewChanges(
			database.NewChange(s.factorRepo.InstanceIDColumn(), f.InstanceID),
			database.NewChange(s.factorRepo.SessionIDColumn(), f.SessionID),
			database.NewChange(s.factorRepo.FactorTypeColumn(), domain.SessionFactorTypePassword),
			database.NewChange(s.factorRepo.LastVerifiedAtColumn(), f.LastVerifiedAt),
			)
	case *domain.SessionFactorTOTP:
		return database.NewChanges(
			database.NewChange(s.factorRepo.InstanceIDColumn(), f.InstanceID),
			database.NewChange(s.factorRepo.SessionIDColumn(), f.SessionID),
			database.NewChange(s.factorRepo.FactorTypeColumn(), domain.SessionFactorTypeTOTP),
			database.NewChange(s.factorRepo.LastVerifiedAtColumn(), f.LastVerifiedAt),
		)
	case *domain.SessionFactorWebAuthn:
		return database.NewChange(s.factorRepo.FactorsColumn(), factor)
	default:
		return database.NewChange(s.factorRepo.FactorsColumn(), factor)
	}
}

// ClearFactor implements [domain.sessionChanges].
func (s session) ClearFactor() database.Change {
	return database.NewChange(s.factorRepo.LastVerifiedAtColumn(), database.NullInstruction)
}

// SetMetadata implements [domain.sessionChanges].
func (s session) SetMetadata(metadata []domain.SessionMetadata) database.Change {
	changes := make([]database.Change, len(metadata))
	for i, md := range metadata {
		changes[i] = database.NewChange(s.metadataRepo.InstanceIDColumn(), md.InstanceID)
		changes[i] = database.NewChange(s.metadataRepo.SessionIDColumn(), md.SessionID)
		changes[i] = database.NewChange(s.metadataRepo.KeyColumn(), md.Key)
		changes[i] = database.NewChange(s.metadataRepo.ValueColumn(), md.Value)
	}
	return database.NewChanges(changes...)
}

// SetUserAgent implements [domain.sessionChanges].
func (s session) SetUserAgent(userAgent domain.SessionUserAgent) database.Change {
	//TODO: upsert user agent?
	return database.NewChanges(
	database.NewChange(s.UserAgentIDColumn(), userAgent.FingerprintID),


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

//
//// UpdatedAtColumn implements [domain.sessionColumns].
//func (s session) FactorColumns() sessionFactor {
//	return sessionFactor{}
//}
//
//// UpdatedAtColumn implements [domain.sessionColumns].
//func (s session) MetadataColumns() sessionMetadataColumns
//
//// UpdatedAtColumn implements [domain.sessionColumns].
//func (s session) UserAgentColumns() sessionUserAgentColumns

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

type rawSession struct {
	*domain.Session
	Factors  JSONArray[domain.SessionFactor]   `json:"factorRepo,omitempty" db:"factorRepo"`
	Metadata JSONArray[domain.SessionMetadata] `json:"metadata,omitempty" db:"metadata"`
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
	raw.Session.Factors = raw.Factors
	raw.Session.Metadata = raw.Metadata

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
		result[i].Factors = session.Factors
		result[i].Metadata = session.Metadata
	}

	return result, nil
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------
//
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
