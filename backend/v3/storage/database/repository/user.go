package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

func UserRepository() domain.UserRepository {
	return new(user)
}

func MachineUserRepository() domain.MachineUserRepository {
	return &userMachine{user: user{}}
}

type user struct {
	verification verification
	tableName    string
	machineKey
	userPersonalAccessToken
	userMetadata
	userPasskey
	userIdentityProviderLink
}

func (u user) HumanRepository() domain.HumanUserRepository {
	return userHuman{user: u}
}

// existingUser is used to get the columns and conditions for the CTE that selects the existing user in update and delete operations.
var existingUser = user{tableName: "existing_user"}

// Create implements [domain.UserRepository].
func (u user) Create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
	var create func(context.Context, *database.StatementBuilder, database.QueryExecutor, *domain.User) error
	switch {
	case user.Human != nil:
		create = userHuman{user: u}.create
	case user.Machine != nil:
		create = userMachine{user: u}.create
	default:
		return zerrors.ThrowInvalidArgument(nil, "REPOS-KxMnG", "no user type defined")
	}
	// we need to cheat a bit here to be able to use CTEChanges
	// we pretend we have an existing user to be able to use the existing change mechanisms
	builder := database.NewStatementBuilder("WITH existing_user AS (SELECT $1 AS instance_id, $2 AS organization_id, $3 AS id)", user.InstanceID, user.OrganizationID, user.ID)
	return create(ctx, builder, client, user)
}

// Delete implements [domain.UserRepository].
func (u user) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkPKCondition(u, condition); err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder("DELETE FROM zitadel.users")
	writeCondition(builder, condition)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

var queryUserStmt = "SELECT users.instance_id, users.organization_id, users.id, users.username" +
	", users.state, users.created_at, users.updated_at" +
	// metadata
	`, jsonb_agg(DISTINCT jsonb_build_object('instanceId', user_metadata.instance_id, 'key', user_metadata.key, 'value', encode(user_metadata.value, 'base64'), 'createdAt', user_metadata.created_at, 'updatedAt', user_metadata.updated_at)) FILTER (WHERE user_metadata.user_id IS NOT NULL) AS metadata` +
	// machine
	`, CASE WHEN users.type = 'machine' THEN jsonb_build_object('name', users.name` +
	`, 'description', users.description, 'secret', users.secret` +
	`, 'accessTokenType', users.access_token_type` +
	`, 'keys', jsonb_agg(DISTINCT jsonb_build_object('id', machine_keys.id, 'publicKey', encode(machine_keys.public_key, 'base64'), 'createdAt', machine_keys.created_at, 'expiresAt', machine_keys.expires_at, 'type', machine_keys.type)) FILTER (WHERE machine_keys.user_id IS NOT NULL)` +
	`, 'pats', jsonb_agg(DISTINCT jsonb_build_object('id', user_personal_access_tokens.id, 'createdAt', user_personal_access_tokens.created_at, 'expiresAt', user_personal_access_tokens.expires_at, 'scopes', user_personal_access_tokens.scopes)) FILTER (WHERE user_personal_access_tokens.user_id IS NOT NULL)` +
	`) END AS machine` +
	// human
	`, CASE WHEN users.type = 'human' THEN jsonb_build_object('firstName', users.first_name, 'lastName', users.last_name` +
	`, 'nickname', users.nickname, 'displayName', users.display_name` +
	`, 'preferredLanguage', users.preferred_language, 'gender', users.gender` +
	`, 'avatarKey', users.avatar_key` +
	`, 'multifactorInitializationSkippedAt', users.multifactor_initialization_skipped_at` +
	`, 'password', jsonb_build_object('hash', users.password_hash, 'isChangeRequired', users.password_change_required, 'changedAt', users.password_changed_at, 'failedAttempts', users.password_failed_attempts, 'lastSuccessfullyCheckedAt', users.password_last_successful_check, 'pendingVerification', ` + verificationQuery(userHuman{}.passwordVerificationIDColumn()) + `)` +
	`, 'email', jsonb_build_object('address', users.email, 'unverifiedAddress', users.unverified_email, 'verifiedAt', users.email_verified_at, 'otp', jsonb_build_object('enabledAt', users.email_otp_enabled_at, 'lastSuccessfullyCheckedAt', users.email_otp_last_successful_check, 'failedAttempts', users.email_otp_failed_attempts), 'pendingVerification', ` + verificationQuery(userHuman{}.emailVerificationIDColumn()) + `)` +
	`, 'phone', CASE WHEN users.phone IS NOT NULL OR users.phone_verification_id IS NOT NULL THEN jsonb_build_object('number', users.phone, 'unverifiedNumber', users.unverified_phone, 'verifiedAt', users.phone_verified_at, 'otp', jsonb_build_object('enabledAt', users.sms_otp_enabled_at, 'lastSuccessfullyCheckedAt', users.sms_otp_last_successful_check, 'failedAttempts', users.sms_otp_failed_attempts), 'pendingVerification', ` + verificationQuery(userHuman{}.phoneVerificationIDColumn()) + `) ELSE NULL END` +
	`, 'totp', CASE WHEN users.totp_secret IS NOT NULL THEN jsonb_build_object('secret', encode(users.totp_secret, 'escape')::JSONB, 'verifiedAt', users.totp_verified_at, 'lastSuccessfullyCheckedAt', users.totp_last_successful_check, 'failedAttempts', users.totp_failed_attempts) ELSE NULL END` +
	`, 'recoveryCodes', CASE WHEN users.recovery_codes IS NOT NULL THEN jsonb_build_object('codes', users.recovery_codes, 'lastSuccessfullyCheckedAt', users.recovery_code_last_successful_check, 'failedAttempts', users.recovery_code_failed_attempts) ELSE NULL END` +
	`, 'passkeys', jsonb_agg(DISTINCT jsonb_build_object('id', user_passkeys.token_id, 'keyId', encode(user_passkeys.key_id, 'base64'), 'type', user_passkeys.type, 'name', user_passkeys.name, 'signCount', user_passkeys.sign_count, 'challenge', encode(user_passkeys.challenge, 'base64'), 'publicKey', encode(user_passkeys.public_key, 'base64'), 'attestationType', user_passkeys.attestation_type, 'aaGuid', encode(user_passkeys.authenticator_attestation_guid, 'base64'), 'rpId', user_passkeys.relying_party_id, 'createdAt', user_passkeys.created_at, 'updatedAt', user_passkeys.updated_at, 'verifiedAt', user_passkeys.verified_at)) FILTER (WHERE user_passkeys.user_id IS NOT NULL)` +
	`, 'verifications', (SELECT jsonb_agg(jsonb_build_object('id', verifications.id, 'code', encode(verifications.code, 'escape')::JSONB, 'createdAt', verifications.created_at, 'updatedAt', verifications.updated_at, 'expiresAt', verifications.created_at+verifications.expiry, 'failedAttempts', verifications.failed_attempts)) FROM zitadel.verifications WHERE verifications.instance_id = users.instance_id AND verifications.user_id = users.id AND verifications.id NOT IN (COALESCE(users.password_verification_id, ''), COALESCE(users.email_verification_id, ''), COALESCE(users.phone_verification_id, ''), COALESCE(users.invite_verification_id, '')))` +
	`, 'identityProviderLinks', jsonb_agg(DISTINCT jsonb_build_object('providerId', user_identity_provider_links.identity_provider_id, 'providedUserId', user_identity_provider_links.provided_user_id, 'providedUsername', user_identity_provider_links.provided_username, 'createdAt', user_identity_provider_links.created_at, 'updatedAt', user_identity_provider_links.updated_at)) FILTER (WHERE user_identity_provider_links.user_id IS NOT NULL)` +
	`, 'invite', CASE WHEN users.invite_verification_id IS NOT NULL OR users.invite_accepted_at IS NOT NULL THEN jsonb_build_object('acceptedAt', users.invite_accepted_at, 'failedAttempts', users.invite_failed_attempts, 'pendingVerification', ` + verificationQuery(userHuman{}.inviteVerificationIDColumn()) + `) ELSE NULL END` +
	`) END AS human FROM zitadel.users`

func verificationQuery(column database.Column) string {
	var builder database.StatementBuilder
	builder.WriteString("CASE WHEN ")
	column.WriteQualified(&builder)
	builder.WriteString(` IS NOT NULL THEN (SELECT row_to_json(res.*) FROM (SELECT verifications.id, encode(verifications.code, 'escape')::JSONB AS code, verifications.created_at as "createdAt", verifications.updated_at as "updatedAt", verifications.created_at+verifications.expiry AS "expiresAt", verifications.failed_attempts as "failedAttempts" FROM zitadel.verifications WHERE verifications.instance_id = users.instance_id AND verifications.id = `)
	column.WriteQualified(&builder)
	builder.WriteString(`) AS res) ELSE NULL END`)
	return builder.String()
}

// Get implements [domain.UserRepository].
func (u user) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.User, error) {
	options := new(database.QueryOpts)
	for _, opt := range u.appendQueryOpts(opts) {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(u.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(u.InstanceIDColumn())
	}

	builder := database.NewStatementBuilder(queryUserStmt)
	options.Write(builder)

	return scanUser(ctx, client, builder)
}

// List implements [domain.UserRepository].
func (u user) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.User, error) {
	options := new(database.QueryOpts)
	for _, opt := range u.appendQueryOpts(opts) {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(u.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(u.InstanceIDColumn())
	}

	builder := database.NewStatementBuilder(queryUserStmt)
	options.Write(builder)

	return scanUsers(ctx, client, builder)
}

func (u user) appendQueryOpts(opts []database.QueryOption) []database.QueryOption {
	return append(opts,
		u.joinMachineKeys(),
		u.joinPATs(),
		u.joinMetadata(),
		u.joinVerifications(),
		u.joinPasskeys(),
		u.joinIdentityProviderLinks(),
		database.WithGroupBy(u.PrimaryKeyColumns()...),
	)
}

// Update implements [domain.UserRepository].
func (u user) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if err := checkPKCondition(u, condition); err != nil {
		return 0, err
	}
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if !database.Changes(changes).IsOnColumn(u.updatedAtColumn()) {
		changes = append(changes, u.clearUpdatedAt())
	}

	builder := database.NewStatementBuilder("WITH existing_user AS (SELECT * FROM zitadel.users")
	writeCondition(builder, condition)
	builder.WriteString(") ")
	for i, change := range changes {
		sessionCTE(change, i, 0, builder)
	}

	builder.WriteString("UPDATE zitadel.users SET ")
	if err := database.Changes(changes).Write(builder); err != nil {
		return 0, err
	}

	builder.WriteString(" FROM existing_user")
	writeCondition(builder, database.And(
		database.NewColumnCondition(u.IDColumn(), existingUser.IDColumn()),
		database.NewColumnCondition(u.InstanceIDColumn(), existingUser.InstanceIDColumn()),
	))

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u user) unqualifiedTableName() string {
	if u.tableName != "" {
		return u.tableName
	}
	return "users"
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetState implements [domain.UserRepository].
func (u user) SetState(state domain.UserState) database.Change {
	return database.NewChange(u.StateColumn(), state)
}

// SetUpdatedAt implements [domain.UserRepository].
func (u user) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(u.updatedAtColumn(), updatedAt)
}

// SetUsername implements [domain.UserRepository].
func (u user) SetUsername(username string) database.Change {
	return database.NewChange(u.UsernameColumn(), username)
}

func (u user) clearUpdatedAt() database.Change {
	return database.NewChange(u.updatedAtColumn(), database.NullInstruction)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// IDCondition implements [domain.UserRepository].
func (u user) IDCondition(userID string) database.Condition {
	return database.NewTextCondition(u.IDColumn(), database.TextOperationEqual, userID)
}

// InstanceIDCondition implements [domain.UserRepository].
func (u user) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(u.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// LoginNameCondition implements [domain.UserRepository].
func (u user) LoginNameCondition(op database.TextOperation, loginName string) database.Condition {
	panic("unimplemented")
}

// OrganizationIDCondition implements [domain.UserRepository].
func (u user) OrganizationIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(u.organizationIDColumn(), database.TextOperationEqual, orgID)
}

// PrimaryKeyCondition implements [domain.UserRepository].
func (u user) PrimaryKeyCondition(instanceID string, userID string) database.Condition {
	return database.And(
		u.InstanceIDCondition(instanceID),
		u.IDCondition(userID),
	)
}

// StateCondition implements [domain.UserRepository].
func (u user) StateCondition(state domain.UserState) database.Condition {
	return database.NewNumberCondition(u.StateColumn(), database.NumberOperationEqual, state)
}

// TypeCondition implements [domain.UserRepository].
func (u user) TypeCondition(userType domain.UserType) database.Condition {
	return database.NewNumberCondition(u.TypeColumn(), database.NumberOperationEqual, userType)
}

// UsernameCondition implements [domain.UserRepository].
func (u user) UsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(u.UsernameColumn(), op, username)
}

// ExistsMetadata implements [domain.UserRepository].
func (u user) ExistsMetadata(condition database.Condition) database.Condition {
	return database.Exists(
		u.userMetadata.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.userMetadata.instanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.userMetadata.userIDColumn()),
			condition,
		),
	)
}

// ExistsPasskey implements [domain.HumanUserRepository].
func (u user) ExistsPasskey(condition database.Condition) database.Condition {
	return database.Exists(
		u.userPasskey.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.userPasskey.instanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.userPasskey.userIDColumn()),
			condition,
		),
	)
}

// ExistsIdentityProviderLink implements [domain.HumanUserRepository].
func (u user) ExistsIdentityProviderLink(condition database.Condition) database.Condition {
	return database.Exists(
		u.userIdentityProviderLink.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.userIdentityProviderLink.instanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.userIdentityProviderLink.userIDColumn()),
			condition,
		),
	)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.UserRepository].
func (u user) PrimaryKeyColumns() []database.Column {
	return database.Columns{
		u.InstanceIDColumn(),
		u.IDColumn(),
	}
}

// CreatedAtColumn implements [domain.UserRepository].
func (u user) CreatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "created_at")
}

// StateColumn implements [domain.UserRepository].
func (u user) StateColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "state")
}

// TypeColumn implements [domain.UserRepository].
func (u user) TypeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "type")
}

// UsernameColumn implements [domain.UserRepository].
func (u user) UsernameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "username")
}

// IDColumn implements [domain.UserRepository].
func (u user) IDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "id")
}

func (u user) InstanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u user) updatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "updated_at")
}

func (u user) organizationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "organization_id")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

type rawUser struct {
	domain.User
	Machine *json.RawMessage `json:"machine,omitempty" db:"machine"`
	Human   *json.RawMessage `json:"human,omitempty" db:"human"`
}

func (u *rawUser) toDomain() (*domain.User, error) {
	if u.Machine != nil {
		err := json.Unmarshal(*u.Machine, &u.User.Machine)
		if err != nil {
			return nil, err
		}
	} else if u.Human != nil {
		err := json.Unmarshal(*u.Human, &u.User.Human)
		if err != nil {
			return nil, err
		}
	}
	return &u.User, nil
}

func scanUser(ctx context.Context, client database.QueryExecutor, builder *database.StatementBuilder) (*domain.User, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var user rawUser
	if err = rows.(database.CollectableRows).CollectExactlyOneRow(&user); err != nil {
		return nil, err
	}
	return user.toDomain()
}

func scanUsers(ctx context.Context, client database.QueryExecutor, builder *database.StatementBuilder) ([]*domain.User, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var users []rawUser
	if err = rows.(database.CollectableRows).Collect(&users); err != nil {
		return nil, err
	}
	result := make([]*domain.User, len(users))
	for i, user := range users {
		domainUser, err := user.toDomain()
		if err != nil {
			return nil, err
		}
		result[i] = domainUser
	}
	return result, nil
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

// Human implements [domain.UserRepository].
func (u user) Human() domain.HumanUserRepository {
	return userHuman{user: u}
}

// Machine implements [domain.UserRepository].
func (u user) Machine() domain.MachineUserRepository {
	return userMachine{user: u}
}

func (u user) joinMachineKeys() database.QueryOption {
	return database.WithLeftJoin(
		u.machineKey.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.machineKey.instanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.machineKey.userIDColumn()),
		),
	)
}

func (u user) joinPATs() database.QueryOption {
	return database.WithLeftJoin(
		u.userPersonalAccessToken.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.userPersonalAccessToken.instanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.userPersonalAccessToken.userIDColumn()),
		),
	)
}

func (u user) joinMetadata() database.QueryOption {
	return database.WithLeftJoin(
		u.userMetadata.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.userMetadata.instanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.userMetadata.userIDColumn()),
		),
	)
}

func (u user) joinPasskeys() database.QueryOption {
	return database.WithLeftJoin(
		u.userPasskey.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.userPasskey.instanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.userPasskey.userIDColumn()),
		),
	)
}

func (u user) joinVerifications() database.QueryOption {
	return database.WithLeftJoin(
		u.verification.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.instanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.verification.userIDColumn()),
		),
	)
}

func (u user) joinIdentityProviderLinks() database.QueryOption {
	return database.WithLeftJoin(
		u.userIdentityProviderLink.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.userIdentityProviderLink.instanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.userIdentityProviderLink.userIDColumn()),
		),
	)
}
