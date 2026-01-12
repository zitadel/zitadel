package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userMachine struct {
	user
}

const createMachineUserStmt = "INSERT INTO zitadel.users (" +
	"instance_id" +
	", organization_id" +
	", id" +
	", username" +
	", username_org_unique" +
	", state" +
	", type" +
	", created_at" +
	", updated_at" +
	", name" +
	", description" +
	", secret" +
	", access_token_type" +
	") VALUES ("

func (u userMachine) create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
	builder := database.NewStatementBuilder(createMachineUserStmt,
		user.InstanceID, user.OrganizationID, user.ID, user.Username, database.NullInstruction, user.State, "machine")
	var createdAt any = database.NowInstruction
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt
	}
	builder.AppendArgs(createdAt, createdAt, user.Machine.Name, user.Machine.Description, user.Machine.Secret, user.Machine.AccessTokenType)
	builder.WriteString(") RETURNING created_at, updated_at")
	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// AddKey implements domain.MachineUserRepository.
func (u *userMachine) AddKey(key *domain.MachineKey) database.Change {
	panic("unimplemented")
}

// AddPersonalAccessToken implements domain.MachineUserRepository.
func (u *userMachine) AddPersonalAccessToken(pat *domain.PersonalAccessToken) database.Change {
	panic("unimplemented")
}

// RemoveKey implements domain.MachineUserRepository.
func (u *userMachine) RemoveKey(id string) database.Change {
	panic("unimplemented")
}

// RemovePersonalAccessToken implements domain.MachineUserRepository.
func (u *userMachine) RemovePersonalAccessToken(id string) database.Change {
	panic("unimplemented")
}

// SetAccessTokenType implements domain.MachineUserRepository.
func (u *userMachine) SetAccessTokenType(tokenType domain.PersonalAccessTokenType) database.Change {
	return database.NewChange(u.accessTokenTypeColumn(), tokenType)
}

// SetDescription implements domain.MachineUserRepository.
func (u *userMachine) SetDescription(description string) database.Change {
	return database.NewChange(u.descriptionColumn(), description)
}

// SetName implements domain.MachineUserRepository.
func (u *userMachine) SetName(name string) database.Change {
	return database.NewChange(u.nameColumn(), name)
}

// SetSecret implements domain.MachineUserRepository.
func (u *userMachine) SetSecret(secret *string) database.Change {
	panic("unimplemented")
}

// Update implements domain.MachineUserRepository.
func (u *userMachine) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	panic("unimplemented")
}

var _ domain.MachineUserRepository = (*userMachine)(nil)

func (u userMachine) nameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "name")
}

func (u userMachine) descriptionColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "description")
}

func (u userMachine) accessTokenTypeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "access_token_type")
}
