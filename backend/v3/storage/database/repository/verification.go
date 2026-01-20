package repository

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type verification struct{}

func (v verification) unqualifiedTableName() string {
	return "verifications"
}

func (v verification) qualifiedTableName() string {
	return "zitadel." + v.unqualifiedTableName()
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (v verification) instanceIDColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "instance_id")
}

func (v verification) idColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "id")
}

func (v verification) valueColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "value")
}

func (v verification) codeColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "code")
}

func (v verification) expiryColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "expires_at")
}

func (v verification) failedAttemptsColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "failed_attempts")
}

func (v verification) creationDateColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "created_at")
}

func (v verification) userIDColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "user_id")
}
