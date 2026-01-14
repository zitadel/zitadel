package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userMachine struct {
	user
}

func (u userMachine) create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
	var createdAt any = database.NowInstruction
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt
	}

	builder := database.NewStatementBuilder(
		"WITH existing_user AS (INSERT INTO zitadel.users ("+
			"instance_id, organization_id, id, username, state, type"+
			", name, description, secret, access_token_type, created_at, updated_at)"+
			" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, ",
		user.InstanceID, user.OrganizationID, user.ID, user.Username, user.State, "machine",
		user.Machine.Name, user.Machine.Description, user.Machine.Secret, uint8(user.Machine.AccessTokenType),
	)
	builder.WriteArgs(createdAt, createdAt)
	builder.WriteString(") RETURNING *)")

	changes := make(database.Changes, 0, len(user.Machine.Keys)+len(user.Machine.PATs))
	for _, key := range user.Machine.Keys {
		changes = append(changes, u.AddKey(key))
	}
	for _, pat := range user.Machine.PATs {
		changes = append(changes, u.AddPersonalAccessToken(pat))
	}
	changes.Write(builder)
	builder.WriteString("SELECT created_at, updated_at FROM existing_user")

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// AddKey implements [domain.MachineUserRepository].
func (u userMachine) AddKey(key *domain.MachineKey) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.machine_keys (" +
				"instance_id, user_id, id, created_at, expires_at, type, public_key" +
				") SELECT instance_id, id, ",
			)
			var createdAt any = database.NowInstruction
			if !key.CreatedAt.IsZero() {
				createdAt = key.CreatedAt
			}
			builder.WriteArgs(key.ID, createdAt, key.ExpiresAt, key.Type, key.PublicKey)
			builder.WriteString(" FROM existing_user")
		},
		nil,
	)
}

// AddPersonalAccessToken implements [domain.MachineUserRepository].
func (u userMachine) AddPersonalAccessToken(pat *domain.PersonalAccessToken) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.machine_user_personal_access_tokens (" +
				"instance_id, user_id, token_id, created_at, expiration, scopes" +
				") SELECT instance_id, id, ",
			)
			var createdAt any = database.NowInstruction
			if !pat.CreatedAt.IsZero() {
				createdAt = pat.CreatedAt
			}
			builder.WriteArgs(pat.ID, createdAt, pat.ExpiresAt, pat.Scopes)
			builder.WriteString(" FROM existing_user")
		},
		nil,
	)
}

// RemoveKey implements [domain.MachineUserRepository].
func (u userMachine) RemoveKey(id string) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM zitadel.machine_keys WHERE " +
				"(instance_id, user_id, id) = (SELECT instance_id, id, ")
			builder.WriteArg(id)
			builder.WriteString(" FROM existing_user)")
		},
		nil,
	)
}

// RemovePersonalAccessToken implements [domain.MachineUserRepository].
func (u userMachine) RemovePersonalAccessToken(id string) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM zitadel.machine_user_personal_access_tokens WHERE " +
				"(instance_id, user_id, token_id) = (SELECT instance_id, id, ")
			builder.WriteArg(id)
			builder.WriteString(" FROM existing_user)")
		},
		nil,
	)
}

// SetAccessTokenType implements [domain.MachineUserRepository].
func (u userMachine) SetAccessTokenType(tokenType domain.PersonalAccessTokenType) database.Change {
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

// Update implements [domain.MachineUserRepository].
func (u userMachine) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
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

func (u userMachine) secretColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "secret")
}
