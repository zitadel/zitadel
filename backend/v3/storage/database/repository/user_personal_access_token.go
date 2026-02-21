package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userPersonalAccessToken struct{}

// AddPersonalAccessToken implements [domain.MachineUserRepository].
func (userPersonalAccessToken) AddPersonalAccessToken(pat *domain.PersonalAccessToken) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.user_personal_access_tokens (" +
				"instance_id, user_id, id, created_at, expires_at, scopes" +
				") SELECT instance_id, id, ",
			)
			var createdAt any = database.NowInstruction
			if !pat.CreatedAt.IsZero() {
				createdAt = pat.CreatedAt
			}
			var expiresAt any = database.NullInstruction
			if !pat.ExpiresAt.IsZero() {
				expiresAt = pat.ExpiresAt
			}
			builder.WriteArgs(pat.ID, createdAt, expiresAt, pat.Scopes)
			builder.WriteString(" FROM existing_user")
		},
		nil,
	)
}

// RemovePersonalAccessToken implements [domain.MachineUserRepository].
func (userPersonalAccessToken) RemovePersonalAccessToken(id string) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM zitadel.user_personal_access_tokens WHERE " +
				"(instance_id, user_id, id) = (SELECT instance_id, id, ")
			builder.WriteArg(id)
			builder.WriteString(" FROM existing_user)")
		},
		nil,
	)
}

func (userPersonalAccessToken) qualifiedTableName() string {
	return "zitadel.user_personal_access_tokens"
}

func (userPersonalAccessToken) unqualifiedTableName() string {
	return "user_personal_access_tokens"
}

func (u userPersonalAccessToken) instanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u userPersonalAccessToken) userIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "user_id")
}
