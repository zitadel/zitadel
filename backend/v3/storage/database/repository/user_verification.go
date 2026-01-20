package repository

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type userVerificationRepo struct{}

func (userVerificationRepo) unqualifiedTableName() string {
	return "verifications"
}
func (userVerificationRepo) qualifiedTableName() string {
	return "zitadel.verifications"
}

func (u userVerificationRepo) instanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u userVerificationRepo) userIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "user_id")
}

func (u userVerificationRepo) idColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "id")
}
