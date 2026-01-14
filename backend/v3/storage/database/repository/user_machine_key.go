package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type machineKeyRepo struct{}

// AddKey implements [domain.MachineUserRepository].
func (u machineKeyRepo) AddKey(key *domain.MachineKey) database.Change {
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

// RemoveKey implements [domain.MachineUserRepository].
func (u machineKeyRepo) RemoveKey(id string) database.Change {
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

func (u machineKeyRepo) qualifiedTableName() string {
	return "zitadel.machine_keys"
}

func (u machineKeyRepo) unqualifiedTableName() string {
	return "machine_keys"
}

func (u machineKeyRepo) instanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u machineKeyRepo) userIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "user_id")
}
