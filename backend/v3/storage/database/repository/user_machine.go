package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userMachine struct {
	user
	machineKey
	userPersonalAccessToken
}

func (u userMachine) create(ctx context.Context, builder *database.StatementBuilder, client database.QueryExecutor, user *domain.User) error {
	changes := make(database.Changes, 0, len(user.Machine.Keys)+len(user.Machine.PATs))
	for _, key := range user.Machine.Keys {
		changes = append(changes, u.AddKey(key))
	}
	for _, pat := range user.Machine.PATs {
		changes = append(changes, u.AddPersonalAccessToken(pat))
	}
	for _, metadata := range user.Metadata {
		changes = append(changes, u.SetMetadata(metadata))
	}
	for i, change := range changes {
		sessionCTE(change, i, 0, builder)
	}

	builder.WriteString("INSERT INTO zitadel.users (" +
		"instance_id, organization_id, id, username, state, type" +
		", name, description, secret, access_token_type, created_at, updated_at) VALUES (",
	)

	var createdAt any = database.NowInstruction
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt
	}
	builder.WriteArgs(
		user.InstanceID, user.OrganizationID, user.ID, user.Username, user.State, "machine",
		user.Machine.Name, user.Machine.Description, user.Machine.Secret, uint8(user.Machine.AccessTokenType),
		createdAt, createdAt,
	)

	builder.WriteString(") RETURNING created_at, updated_at")

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// Update implements [domain.MachineUserRepository].
func (u userMachine) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if !condition.IsRestrictingColumn(u.TypeColumn()) {
		condition = database.And(condition, u.TypeCondition(domain.UserTypeMachine))
	}
	return u.user.Update(ctx, client, condition, changes...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetAccessTokenType implements [domain.MachineUserRepository].
func (u userMachine) SetAccessTokenType(tokenType domain.AccessTokenType) database.Change {
	return database.NewChange(u.accessTokenTypeColumn(), tokenType)
}

// SetDescription implements [domain.MachineUserRepository].
func (u userMachine) SetDescription(description string) database.Change {
	return database.NewChange(u.descriptionColumn(), description)
}

// SetName implements [domain.MachineUserRepository].
func (u userMachine) SetName(name string) database.Change {
	return database.NewChange(u.nameColumn(), name)
}

// SetSecret implements [domain.MachineUserRepository].
func (u userMachine) SetSecret(secret *string) database.Change {
	return database.NewChangePtr(u.secretColumn(), secret)
}

var _ domain.MachineUserRepository = (*userMachine)(nil)

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userMachine) nameColumn() database.Column {
	return database.NewColumn(u.user.unqualifiedTableName(), "name")
}

func (u userMachine) descriptionColumn() database.Column {
	return database.NewColumn(u.user.unqualifiedTableName(), "description")
}

func (u userMachine) accessTokenTypeColumn() database.Column {
	return database.NewColumn(u.user.unqualifiedTableName(), "access_token_type")
}

func (u userMachine) secretColumn() database.Column {
	return database.NewColumn(u.user.unqualifiedTableName(), "secret")
}
