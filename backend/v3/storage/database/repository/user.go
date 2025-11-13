package repository

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.UserRepository = (*user)(nil)

type user struct {
	shouldLoadMetadata bool
	metadata           userMetadata

	shouldLoadIdentityProviderLinks bool
	identityProviderLinks           userIdentityProviderLink

	shouldLoadKeys bool
	keys           userMachineKey

	shouldLoadPATs bool
	pats           userPersonalAccessToken
}

func UserRepository() domain.UserRepository {
	return new(user)
}

func (u user) unqualifiedTableName() string {
	return "users"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

// Create implements [domain.UserRepository].
func (u user) Create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
	if user.Machine != nil {
		return userMachine{user: &u}.create(ctx, client, user)
	}
	if user.Human != nil {
		return userHuman{user: &u}.create(ctx, client, user)
	}
	// TODO(adlerhurst): return a proper error here
	return database.NewCheckError(u.unqualifiedTableName(), "type", errors.New("no type specified"))
}

// Update implements [domain.UserRepository].
func (u user) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if err := checkPKCondition(u, condition); err != nil {
		return 0, err
	}
	if !database.Changes(changes).IsOnColumn(u.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(u.UpdatedAtColumn(), database.NullInstruction))
	}
	builder := database.NewStatementBuilder(`UPDATE zitadel.users SET `)
	database.Changes(changes).Write(builder)
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// Delete implements [domain.UserRepository].
func (u user) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkPKCondition(u, condition); err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder(`DELETE FROM zitadel.users `)
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

const userQuery = "SELECT" +
	" users.instance_id" +
	" , users.organization_id" +
	" , users.id" +
	" , users.username" +
	" , users.username_org_unique" +
	" , users.state" +
	" , users.type" +
	" , users.created_at" +
	" , users.updated_at" +
	" , users.first_name AS human.first_name" +
	" , users.last_name AS human.last_name" +
	" , users.nickname AS human.nickname" +
	" , users.display_name AS human.display_name" +
	" , users.preferred_language AS human.preferred_language" +
	" , users.gender AS human.gender" +
	" , users.avatar_key AS human.avatar_key" +
	// " , users.multi_factor_initialization_skipped_at" +
	// " , users.password" +
	// " , users.password_change_required" +
	// " , users.password_verified_at" +
	// " , users.failed_password_attempts" +
	// " , users.email" +
	// " , users.email_verified_at" +
	// " , users.email_otp_enabled" +
	// " , users.last_successful_email_otp_check" +
	// " , users.phone" +
	// " , users.phone_verified_at" +
	// " , users.sms_otp_enabled" +
	// " , users.last_successful_sms_otp_check" +
	// " , users.last_successful_totp_check" +
	" , users.name AS machine.name" +
	" , users.description AS machine.description" +
	" , users.secret AS machine.Secret" +
	" , users.access_token_type AS machine.AccessTokenType" +
	` , jsonb_agg(json_build_object('instanceId', user_metadata.instance_id, 'orgId', user_metadata.organization_id, 'key', user_metadata.key, 'value', encode(user_metadata.value, 'base64'), 'createdAt', user_metadata.created_at, 'updatedAt', user_metadata.updated_at)) FILTER (WHERE user_metadata.organization_id IS NOT NULL) AS metadata` +
	" FROM zitadel.users "

// Get implements [domain.UserRepository].
func (u user) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.User, error) {
	builder := database.NewStatementBuilder(userQuery)
	opts = append(opts,
		u.joinMetadata(),
		database.WithGroupBy(u.PrimaryKeyColumns()...),
	)

	for _, option := range opts {
		option(&database.QueryOpts{})
	}

	user, err := getOne[rawUser](ctx, client, builder)
	if err != nil {
		return nil, err
	}
	user.User.Metadata = user.Metadata

	return user.User, nil
}

// List implements [domain.UserRepository].
func (u user) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.User, error) {
	builder := database.NewStatementBuilder(userQuery)
	opts = append(opts,
		u.joinMetadata(),
		database.WithGroupBy(u.PrimaryKeyColumns()...),
	)

	for _, option := range opts {
		option(&database.QueryOpts{})
	}

	users, err := getMany[rawUser](ctx, client, builder)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.User, len(users))
	for i, user := range users {
		user.User.Metadata = user.Metadata
		result[i] = user.User
	}
	return result, nil
}

type rawUser struct {
	*domain.User
	Metadata JSONArray[domain.UserMetadata] `json:"metadata,omitempty" db:"metadata"`
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetState implements [domain.UserRepository].
func (u user) SetState(state domain.UserState) database.Change {
	return database.NewChange(u.StateColumn(), state.String())
}

// SetUpdatedAt implements [domain.UserRepository].
func (u user) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(u.UpdatedAtColumn(), updatedAt)
}

// SetUsername implements [domain.UserRepository].
func (u user) SetUsername(username string) database.Change {
	return database.NewChange(u.UsernameColumn(), username)
}

// SetUsernameOrgUnique implements [domain.UserRepository].
func (u user) SetUsernameOrgUnique(usernameOrgUnique bool) database.Change {
	return database.NewChange(u.UsernameOrgUniqueColumn(), usernameOrgUnique)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (u user) PrimaryKeyCondition(instanceID, userID string) database.Condition {
	return database.And(
		u.InstanceIDCondition(instanceID),
		u.IDCondition(userID),
	)
}

// CreatedAtCondition implements [domain.UserRepository].
func (u user) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(u.CreatedAtColumn(), op, createdAt)
}

// IDCondition implements [domain.UserRepository].
func (u user) IDCondition(userID string) database.Condition {
	return database.NewTextCondition(u.IDColumn(), database.TextOperationEqual, userID)
}

// InstanceIDCondition implements [domain.UserRepository].
func (u user) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(u.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// OrgIDCondition implements [domain.UserRepository].
func (u user) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(u.OrgIDColumn(), database.TextOperationEqual, orgID)
}

// StateCondition implements [domain.UserRepository].
func (u user) StateCondition(state domain.UserState) database.Condition {
	return database.NewTextCondition(u.StateColumn(), database.TextOperationEqual, state.String())
}

// UpdatedAtCondition implements [domain.UserRepository].
func (u user) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(u.UpdatedAtColumn(), op, updatedAt)
}

// UsernameCondition implements [domain.UserRepository].
func (u user) UsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(u.UsernameColumn(), op, username)
}

// UsernameOrgUniqueCondition implements [domain.UserRepository].
func (u user) UsernameOrgUniqueCondition(condition bool) database.Condition {
	return database.NewBooleanCondition(u.UsernameOrgUniqueColumn(), condition)
}

// TypeCondition implements [domain.UserRepository].
func (u user) TypeCondition(userType domain.UserType) database.Condition {
	return database.NewNumberCondition(u.typeColumn(), database.NumberOperationEqual, userType)
}

func (u user) ExistsMetadata(cond database.Condition) database.Condition {
	return database.Exists(
		u.metadata.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.metadata.InstanceIDColumn()),
			database.NewColumnCondition(u.IDColumn(), u.metadata.UserIDColumn()),
			cond,
		),
	)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u user) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		u.InstanceIDColumn(),
		u.IDColumn(),
	}
}

// CreatedAtColumn implements [domain.UserRepository].
func (u user) CreatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "created_at")
}

// IDColumn implements [domain.UserRepository].
func (u user) IDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "id")
}

// InstanceIDColumn implements [domain.UserRepository].
func (u user) InstanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

// OrgIDColumn implements [domain.UserRepository].
func (u user) OrgIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "org_id")
}

// StateColumn implements [domain.UserRepository].
func (u user) StateColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "state")
}

// UpdatedAtColumn implements [domain.UserRepository].
func (u user) UpdatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "updated_at")
}

// UsernameColumn implements [domain.UserRepository].
func (u user) UsernameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "username")
}

// UsernameOrgUniqueColumn implements [domain.UserRepository].
func (u user) UsernameOrgUniqueColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "username_org_unique")
}

func (u user) typeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "type")
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

// Human implements [domain.UserRepository].
func (u user) Human() domain.HumanUserRepository {
	return &userHuman{user: &u}
}

// Machine implements [domain.UserRepository].
func (u user) Machine() domain.MachineUserRepository {
	return &userMachine{user: &u}
}

func (u user) LoadMetadata() domain.UserRepository {
	return &user{
		shouldLoadMetadata: true,
	}
}

func (u user) copy() user {
	return user{
		shouldLoadMetadata:              u.shouldLoadMetadata,
		metadata:                        u.metadata,
		shouldLoadIdentityProviderLinks: u.shouldLoadIdentityProviderLinks,
		identityProviderLinks:           u.identityProviderLinks,
		shouldLoadKeys:                  u.shouldLoadKeys,
		keys:                            u.keys,
		shouldLoadPATs:                  u.shouldLoadPATs,
		pats:                            u.pats,
	}
}

func (u user) joinMetadata() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(u.InstanceIDColumn(), u.metadata.InstanceIDColumn()),
		database.NewColumnCondition(u.IDColumn(), u.metadata.UserIDColumn()),
	)

	// If metadata should not be joined, we make sure to return null for the metadata columns
	// the query optimizer of the dialect should optimize this away if no metadata are requested
	if !u.shouldLoadMetadata {
		columns = append(columns, database.IsNull(u.metadata.UserIDColumn()))
	}

	return database.WithLeftJoin(
		u.metadata.qualifiedTableName(),
		database.And(columns...),
	)
}
