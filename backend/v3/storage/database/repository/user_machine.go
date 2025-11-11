package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.MachineUserRepository = (*userMachine)(nil)

type userMachine struct {
	shouldLoadMetadata bool
	metadata           userMetadata

	shouldLoadKeys bool
	keys           userMachineKey

	shouldLoadPATs bool
	pats           userPersonalAccessToken
}

func MachineUserRepository() domain.MachineUserRepository {
	return new(userMachine)
}

func (m userMachine) unqualifiedTableName() string {
	return "users"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

const insertMachineStmt = "INSERT INTO zitadel.machine_users (" +
	"instance_id, organization_id, id, username, username_org_unique, state, created_at, updated_at, name, description, secret, access_token_type" +
	") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING created_at, updated_at"

// create inserts a new machine user into the database.
// the type of the user must be checked before calling this method.
func (m userMachine) create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
	var createdAt, updatedAt any = database.DefaultInstruction, database.DefaultInstruction
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt
	}
	if !user.UpdatedAt.IsZero() {
		updatedAt = user.UpdatedAt
	}

	return client.QueryRow(
		ctx, insertMachineStmt,
		user.InstanceID,
		user.OrgID,
		user.ID,
		user.Username,
		user.IsUsernameOrgUnique,
		user.State,
		createdAt,
		updatedAt,
		user.Machine.Name,
		user.Machine.Description,
		user.Machine.Secret,
		user.Machine.AccessTokenType,
	).Scan(
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

// Update implements [domain.MachineUserRepository].
func (m userMachine) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if err := checkPKCondition(m, condition); err != nil {
		return 0, err
	}
	if !condition.IsRestrictingColumn(m.typeColumn()) {
		condition = database.And(
			condition,
			m.TypeCondition(domain.UserTypeMachine),
		)
	}
	if !database.Changes(changes).IsOnColumn(m.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(m.UpdatedAtColumn(), database.NullInstruction))
	}
	builder := database.NewStatementBuilder(`UPDATE zitadel.machine_users SET `)
	database.Changes(changes).Write(builder)
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// Keys implements [domain.MachineUserRepository].
func (m userMachine) Keys() domain.MachineKeyRepository {
	return m.keys
}

// LoadKeys implements [domain.MachineUserRepository].
func (m userMachine) LoadKeys() domain.MachineUserRepository {
	return &userMachine{
		shouldLoadMetadata: m.shouldLoadMetadata,
		metadata:           m.metadata,
		shouldLoadKeys:     true,
		keys:               m.keys,
	}
}

// PersonalAccessTokens implements [domain.MachineUserRepository].
func (m userMachine) PersonalAccessTokens() domain.PersonalAccessTokenRepository {
	return m.pats
}

// LoadPersonalAccessTokens implements [domain.MachineUserRepository].
func (m userMachine) LoadPersonalAccessTokens() domain.MachineUserRepository {
	return &userMachine{
		shouldLoadMetadata: m.shouldLoadMetadata,
		metadata:           m.metadata,
		shouldLoadKeys:     m.shouldLoadKeys,
		keys:               m.keys,
		shouldLoadPATs:     true,
		pats:               m.pats,
	}
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetAccessTokenType implements [domain.MachineUserRepository.SetAccessTokenType].
func (m userMachine) SetAccessTokenType(accessTokenType domain.AccessTokenType) database.Change {
	return database.NewChange(m.AccessTokenTypeColumn(), accessTokenType)
}

// SetDescription implements [domain.MachineUserRepository.SetDescription].
func (m userMachine) SetDescription(description *string) database.Change {
	return database.NewChangePtr(m.DescriptionColumn(), description)
}

// SetName implements [domain.MachineUserRepository.SetName].
func (m userMachine) SetName(name string) database.Change {
	return database.NewChange(m.NameColumn(), name)
}

// SetSecret implements [domain.MachineUserRepository.SetSecret].
func (m userMachine) SetSecret(secret *string) database.Change {
	return database.NewChangePtr(m.SecretColumn(), secret)
}

// SetState implements [domain.MachineUserRepository.SetState].
func (m userMachine) SetState(state domain.UserState) database.Change {
	return database.NewChange(m.StateColumn(), state)
}

// SetUpdatedAt implements [domain.MachineUserRepository.SetUpdatedAt].
func (m userMachine) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(m.UpdatedAtColumn(), updatedAt)
}

// SetUsername implements [domain.MachineUserRepository.SetUsername].
func (m userMachine) SetUsername(username string) database.Change {
	return database.NewChange(m.UsernameColumn(), username)
}

// SetUsernameOrgUnique implements [domain.MachineUserRepository.SetUsernameOrgUnique].
func (m userMachine) SetUsernameOrgUnique(usernameOrgUnique bool) database.Change {
	return database.NewChange(m.UsernameOrgUniqueColumn(), usernameOrgUnique)
}

func (m userMachine) typeColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "type")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PrimaryKeyCondition implements [domain.MachineUserRepository].
func (m userMachine) PrimaryKeyCondition(instanceID string, userID string) database.Condition {
	return database.And(
		m.InstanceIDCondition(instanceID),
		m.IDCondition(userID),
	)
}

// AccessTokenTypeCondition implements [domain.MachineUserRepository].
func (m userMachine) AccessTokenTypeCondition(accessTokenType domain.AccessTokenType) database.Condition {
	panic("unimplemented")
}

// CreatedAtCondition implements [domain.MachineUserRepository].
func (m userMachine) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(m.CreatedAtColumn(), op, createdAt)
}

// DescriptionCondition implements [domain.MachineUserRepository].
func (m userMachine) DescriptionCondition(op database.TextOperation, description string) database.Condition {
	return database.NewTextCondition(m.DescriptionColumn(), op, description)
}

// IDCondition implements [domain.MachineUserRepository].
func (m userMachine) IDCondition(userID string) database.Condition {
	return database.NewTextCondition(m.IDColumn(), database.TextOperationEqual, userID)
}

// InstanceIDCondition implements [domain.MachineUserRepository].
func (m userMachine) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(m.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// NameCondition implements [domain.MachineUserRepository].
func (m userMachine) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(m.NameColumn(), op, name)
}

// OrgIDCondition implements [domain.MachineUserRepository].
func (m userMachine) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(m.OrgIDColumn(), database.TextOperationEqual, orgID)
}

// StateCondition implements [domain.MachineUserRepository].
func (m userMachine) StateCondition(state domain.UserState) database.Condition {
	return database.NewTextCondition(m.StateColumn(), database.TextOperationEqual, state.String())
}

// UpdatedAtCondition implements [domain.MachineUserRepository].
func (m userMachine) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(m.UpdatedAtColumn(), op, updatedAt)
}

// UsernameCondition implements [domain.MachineUserRepository].
func (m userMachine) UsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(m.UsernameColumn(), op, username)
}

// UsernameOrgUniqueCondition implements [domain.MachineUserRepository].
func (m userMachine) UsernameOrgUniqueCondition(condition bool) database.Condition {
	return database.NewBooleanCondition(m.UsernameOrgUniqueColumn(), condition)
}

// TypeCondition implements [domain.MachineUserRepository].
func (m userMachine) TypeCondition(userType domain.UserType) database.Condition {
	// TODO(adlerhurst): it doesn't make sense to have this method on userMachine
	return user{}.TypeCondition(userType)
}

func (u userMachine) ExistsMetadata(cond database.Condition) database.Condition {
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

// PrimaryKeyColumns implements [domain.MachineUserRepository].
func (m userMachine) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		m.InstanceIDColumn(),
		m.IDColumn(),
	}
}

// AccessTokenTypeColumn implements [domain.MachineUserRepository].
func (m userMachine) AccessTokenTypeColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "access_token_type")
}

// CreatedAtColumn implements [domain.MachineUserRepository].
func (m userMachine) CreatedAtColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "created_at")
}

// DescriptionColumn implements [domain.MachineUserRepository].
func (m userMachine) DescriptionColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "description")
}

// IDColumn implements [domain.MachineUserRepository].
func (m userMachine) IDColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "id")
}

// InstanceIDColumn implements [domain.MachineUserRepository].
func (m userMachine) InstanceIDColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "instance_id")
}

// NameColumn implements [domain.MachineUserRepository].
func (m userMachine) NameColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "name")
}

// OrganizationIDColumn implements [domain.MachineUserRepository].
func (m userMachine) OrgIDColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "organization_id")
}

// SecretColumn implements [domain.MachineUserRepository].
func (m userMachine) SecretColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "secret")
}

// StateColumn implements [domain.MachineUserRepository].
func (m userMachine) StateColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "state")
}

// UpdatedAtColumn implements [domain.MachineUserRepository].
func (m userMachine) UpdatedAtColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "updated_at")
}

// UsernameColumn implements [domain.MachineUserRepository].
func (m userMachine) UsernameColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "username")
}

// UsernameOrgUniqueColumn implements [domain.MachineUserRepository].
func (m userMachine) UsernameOrgUniqueColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "username_org_unique")
}
