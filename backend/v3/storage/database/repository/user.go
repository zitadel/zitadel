package repository

import (
	"context"
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
	// TODO: implement machine user repository
	return nil
}

type user struct {
	verification verification
	tableName    string
}

// existingUser is used to get the columns and conditions for the CTE that selects the existing user in update and delete operations.
var existingUser = user{tableName: "existing_user"}

// Create implements [domain.UserRepository.Create].
func (u user) Create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
	panic("unimplemented")
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

// Get implements [domain.UserRepository.Get].
func (u user) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.User, error) {
	panic("unimplemented")
}

// List implements [domain.UserRepository.List].
func (u user) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.User, error) {
	panic("unimplemented")
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
		database.NewColumnCondition(u.instanceIDColumn(), existingUser.instanceIDColumn()),
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
				database.NewColumnCondition(existingUser.instanceIDColumn(), u.metadataInstanceIDColumn()),
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
	return database.NewTextCondition(u.instanceIDColumn(), database.TextOperationEqual, instanceID)
}

// LoginNameCondition implements [domain.UserRepository.LoginNameCondition].
func (u user) LoginNameCondition(op database.TextOperation, loginName string) database.Condition {
	panic("unimplemented")
}

// MetadataKeyCondition implements [domain.UserRepository.MetadataKeyCondition].
func (u user) MetadataKeyCondition(op database.TextOperation, key string) database.Condition {
	return database.NewTextCondition(u.metadataKeyColumn(), op, key)
}

// MetadataValueCondition implements [domain.UserRepository.MetadataValueCondition].
func (u user) MetadataValueCondition(op database.BytesOperation, value []byte) database.Condition {
	return database.NewBytesCondition[[]byte](database.SHA256Column(u.metadataValueColumn()), op, database.SHA256Value(value))
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
		u.instanceIDColumn(),
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

func (u user) instanceIDColumn() database.Column {
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
	panic("unimplemented")
}
