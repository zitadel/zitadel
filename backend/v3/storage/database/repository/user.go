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
	machine userMachine
	human   userHuman

	shouldLoadMetadata bool
	metadata           userMetadata
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
		return u.machine.create(ctx, client, user)
	}
	if user.Human != nil {
		return u.human.create(ctx, client, user)
	}
	// TODO(adlerhurst): return a proper error here
	return database.NewCheckError(u.unqualifiedTableName(), "type", errors.New("no type specified"))
}

// Update implements [domain.UserRepository].
func (u *user) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
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

// Get implements [domain.UserRepository].
func (u user) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.User, error) {
	opts = append(opts,
		u.joinMetadata(),
		database.WithGroupBy(u.PrimaryKeyColumns()...),
	)
	panic("unimplemented")
}

// List implements [domain.UserRepository].
func (u user) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.User, error) {
	opts = append(opts,
		u.joinMetadata(),
		database.WithGroupBy(u.PrimaryKeyColumns()...),
	)
	panic("unimplemented")
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
	switch userType {
	case domain.UserTypeHuman:
		return database.IsNull(u.machine.IDColumn())
	case domain.UserTypeMachine:
		return database.IsNull(u.human.IDColumn())
	}
	// TODO(adlerhurst): add a log line here to indicate an invalid user type was provided
	return database.IsNull(u.IDColumn())
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

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

// Human implements [domain.UserRepository].
func (u user) Human() domain.HumanUserRepository {
	return u.human
}

// Machine implements [domain.UserRepository].
func (u user) Machine() domain.MachineUserRepository {
	return u.machine
}

func (u user) LoadMetadata() domain.UserRepository {
	return &user{
		shouldLoadMetadata: true,
		machine: userMachine{
			shouldLoadMetadata: true,
		},
		human: userHuman{
			shouldLoadMetadata: true,
		},
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
