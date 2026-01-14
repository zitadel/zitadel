package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
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
}

func (u user) HumanRepository() domain.HumanUserRepository {
	return userHuman{user: u}
}

// existingUser is used to get the columns and conditions for the CTE that selects the existing user in update and delete operations.
var existingUser = user{tableName: "existing_user"}

// Create implements [domain.UserRepository.Create].
func (u user) Create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
	var create func(context.Context, database.QueryExecutor, *domain.User) error
	if user.Human != nil {
		create = userHuman{user: u}.create
	} else if user.Machine != nil {
		create = userMachine{user: u}.create
	}
	return create(ctx, client, user)
}

// Delete implements [domain.UserRepository.Delete].
func (u user) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkPKCondition(u, condition); err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder("DELETE FROM zitadel.users")
	writeCondition(builder, condition)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

const queryUserStmt = "SELECT users.instance_id, users.organization_id, users.id, users.username" +
	", users.state, users.created_at, users.updated_at" +
	// machine columns
	`, CASE WHEN users.type = 'machine' THEN jsonb_build_object('name', users.name` +
	`, 'description', users.description, 'secret', encode(users.secret, 'base64')` +
	`, 'accessTokenType', users.access_token_type) END AS machine` +
	// human columns
	`, CASE WHEN users.type = 'human' THEN jsonb_build_object('firstName', users.first_name, 'lastName', users.last_name` +
	`, 'nickname', users.nickname, 'displayName', users.display_name` +
	`, 'preferredLanguage', users.preferred_language, 'gender', users.gender` +
	`, 'avatarKey', users.avatar_key` +
	`, 'multifactorInitializationSkippedAt', users.multifactor_initialization_skipped_at` +
	`, 'password', encode(users.password, 'base64'), 'passwordChangeRequired', users.password_change_required` +
	`, 'passwordVerifiedAt', users.password_verified_at` +
	// -- users.password_verification_id
	`, 'failedPasswordAttempts', users.failed_password_attempts, 'email', users.email` +
	`, 'emailVerifiedAt', users.email_verified_at` +
	// -- , users.unverified_email_id
	`, 'emailOtpEnabledAt', users.email_otp_enabled_at` +
	`, 'lastSuccessfulEmailOtpCheck', users.last_successful_email_otp_check` +
	// -- , users.email_otp_verification_id
	`, 'phone', users.phone, 'phoneVerifiedAt', users.phone_verified_at` +
	// -- , users.unverified_phone_id
	`, 'smsOtpEnabledAt', users.sms_otp_enabled_at` +
	`, 'lastSuccessfulSmsOtpCheck', users.last_successful_sms_otp_check` +
	// -- , users.sms_otp_verification_id
	// -- , totp_secret_id
	`,'totpVerifiedAt', users.totp_verified_at` +
	// -- , unverified_totp_id
	`,'lastSuccessfulTotpCheck', users.last_successful_totp_check` +
	`) END AS human FROM zitadel.users`

// Get implements [domain.UserRepository.Get].
func (u user) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.User, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(u.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(u.InstanceIDColumn())
	}

	builder := database.NewStatementBuilder(queryUserStmt)
	options.Write(builder)

	return scanUser(ctx, client, builder)
}

// List implements [domain.UserRepository.List].
func (u user) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.User, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(u.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(u.InstanceIDColumn())
	}

	builder := database.NewStatementBuilder(queryUserStmt)
	options.Write(builder)

	return scanUsers(ctx, client, builder)
}

// Update implements [domain.UserRepository.Update].
func (u user) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if err := checkPKCondition(u, condition); err != nil {
		return 0, err
	}
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}

	builder := database.NewStatementBuilder("WITH existing_user AS (SELECT id, instance_id FROM zitadel.users")
	writeCondition(builder, condition)
	builder.WriteRune(')')

	for i, change := range changes {
		if cteChange, ok := change.(database.CTEChange); ok {
			builder.WriteString(", ")
			cteName := fmt.Sprintf("cte_%d", i)
			builder.WriteString(cteName)
			builder.WriteString(" AS (")
			cteChange.SetName(cteName)
			cteChange.WriteCTE(builder)
			builder.WriteRune(')')
		}
	}

	if !database.Changes(changes).IsOnColumn(u.updatedAtColumn()) {
		changes = append(changes, u.clearUpdatedAt())
	}

	builder.WriteString(" UPDATE zitadel.users SET ")
	database.Changes(changes).Write(builder)
	builder.WriteString(" FROM existing_user")
	writeCondition(builder, database.And(
		database.NewColumnCondition(u.idColumn(), existingUser.idColumn()),
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

func (u user) unqualifiedMetadataTableName() string {
	return "user_metadata"
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// AddMetadata implements [domain.UserRepository.AddMetadata].
func (u user) AddMetadata(metadata ...*domain.Metadata) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.user_metadata(instance_id, user_id, key, value, created_at, updated_at)")
			builder.WriteString(" SELECT existing_user.instance_id, existing_user.id, md.key, md.value, md.created_at, md.updated_at")
			builder.WriteString(" FROM existing_user CROSS JOIN (VALUES ")
			for i, md := range metadata {
				if i > 0 {
					builder.WriteString(", ")
				}
				builder.WriteRune('(')
				builder.WriteArgs(md.Key, md.Value)
				var createdAt, updatedAt any = database.DefaultInstruction, database.NullInstruction
				if !md.CreatedAt.IsZero() {
					createdAt = md.CreatedAt
				}
				if !md.UpdatedAt.IsZero() {
					updatedAt = md.UpdatedAt
				}
				builder.WriteArgs(createdAt, updatedAt)
				builder.WriteRune(')')
			}
			builder.WriteString(") AS md(key, value, created_at, updated_at)")
			builder.WriteString("ON CONFLICT (")
			database.Columns{
				u.metadataInstanceIDColumn(),
				u.metadataUserIDColumn(),
				u.metadataKeyColumn(),
			}.WriteQualified(builder)
			builder.WriteString(") DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at")
		},
		nil,
	)
}

// RemoveMetadata implements [domain.UserRepository.RemoveMetadata].
func (u user) RemoveMetadata(condition database.Condition) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM zitadel.user_metadata USING ")
			builder.WriteString(existingUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(existingUser.InstanceIDColumn(), u.metadataInstanceIDColumn()),
				database.NewColumnCondition(existingUser.idColumn(), u.metadataUserIDColumn()),
				condition,
			))
		},
		nil,
	)
}

// SetState implements [domain.UserRepository.SetState].
func (u user) SetState(state domain.UserState) database.Change {
	return database.NewChange(u.StateColumn(), state)
}

// SetUpdatedAt implements [domain.UserRepository.SetUpdatedAt].
func (u user) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(u.updatedAtColumn(), updatedAt)
}

// SetUsername implements [domain.UserRepository.SetUsername].
func (u user) SetUsername(username string) database.Change {
	return database.NewChange(u.UsernameColumn(), username)
}

func (u user) clearUpdatedAt() database.Change {
	return database.NewChange(u.updatedAtColumn(), database.NullInstruction)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// ExistsMetadata implements [domain.UserRepository.ExistsMetadata].
func (u user) ExistsMetadata(condition database.Condition) database.Condition {
	panic("unimplemented")
}

// IDCondition implements [domain.UserRepository.IDCondition].
func (u user) IDCondition(userID string) database.Condition {
	return database.NewTextCondition(u.idColumn(), database.TextOperationEqual, userID)
}

// InstanceIDCondition implements [domain.UserRepository.InstanceIDCondition].
func (u user) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(u.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// LoginNameCondition implements [domain.UserRepository.LoginNameCondition].
func (u user) LoginNameCondition(op database.TextOperation, loginName string) database.Condition {
	panic("unimplemented")
}

func (u user) MetadataConditions() domain.UserMetadataConditions {
	return userMetadataConditions{user: u}
}

// OrganizationIDCondition implements [domain.UserRepository.OrganizationIDCondition].
func (u user) OrganizationIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(u.organizationIDColumn(), database.TextOperationEqual, orgID)
}

// PrimaryKeyCondition implements [domain.UserRepository.PrimaryKeyCondition].
func (u user) PrimaryKeyCondition(instanceID string, userID string) database.Condition {
	return database.And(
		u.InstanceIDCondition(instanceID),
		u.IDCondition(userID),
	)
}

// StateCondition implements [domain.UserRepository.StateCondition].
func (u user) StateCondition(state domain.UserState) database.Condition {
	return database.NewNumberCondition(u.StateColumn(), database.NumberOperationEqual, state)
}

// TypeCondition implements [domain.UserRepository.TypeCondition].
func (u user) TypeCondition(userType domain.UserType) database.Condition {
	return database.NewNumberCondition(u.TypeColumn(), database.NumberOperationEqual, userType)
}

// UsernameCondition implements [domain.UserRepository.UsernameCondition].
func (u user) UsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(u.UsernameColumn(), op, username)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.UserRepository.PrimaryKeyColumns].
func (u user) PrimaryKeyColumns() []database.Column {
	return database.Columns{
		u.InstanceIDColumn(),
		u.idColumn(),
	}
}

// CreatedAtColumn implements [domain.UserRepository.CreatedAtColumn].
func (u user) CreatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "created_at")
}

// StateColumn implements [domain.UserRepository.StateColumn].
func (u user) StateColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "state")
}

// TypeColumn implements [domain.UserRepository.TypeColumn].
func (u user) TypeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "type")
}

// UsernameColumn implements [domain.UserRepository.UsernameColumn].
func (u user) UsernameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "username")
}

func (u user) idColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "id")
}

func (u user) InstanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u user) updatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "updated_at")
}

func (u user) metadataInstanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedMetadataTableName(), "instance_id")
}

func (u user) organizationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "organization_id")
}

func (u user) metadataUserIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedMetadataTableName(), "user_id")
}

func (u user) metadataKeyColumn() database.Column {
	return database.NewColumn(u.unqualifiedMetadataTableName(), "key")
}

func (u user) metadataValueColumn() database.Column {
	return database.NewColumn(u.unqualifiedMetadataTableName(), "value")
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

// Human implements [domain.UserRepository.Human].
func (u user) Human() domain.HumanUserRepository {
	return userHuman{user: u}
}

// LoadIdentityProviderLinks implements [domain.UserRepository.LoadIdentityProviderLinks].
func (u user) LoadIdentityProviderLinks() domain.UserRepository {
	panic("unimplemented")
}

// LoadKeys implements [domain.UserRepository.LoadKeys].
func (u user) LoadKeys() domain.UserRepository {
	panic("unimplemented")
}

// LoadMetadata implements [domain.UserRepository.LoadMetadata].
func (u user) LoadMetadata() domain.UserRepository {
	panic("unimplemented")
}

// LoadPATs implements [domain.UserRepository.LoadPATs].
func (u user) LoadPATs() domain.UserRepository {
	panic("unimplemented")
}

// LoadPasskeys implements [domain.UserRepository.LoadPasskeys].
func (u user) LoadPasskeys() domain.UserRepository {
	panic("unimplemented")
}

// LoadVerifications implements [domain.UserRepository.LoadVerifications].
func (u user) LoadVerifications() domain.UserRepository {
	panic("unimplemented")
}

// Machine implements [domain.UserRepository.Machine].
func (u user) Machine() domain.MachineUserRepository {
	return userMachine{user: u}
}
