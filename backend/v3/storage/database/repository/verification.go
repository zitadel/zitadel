package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type verification struct{}

func newVerification() *verification {
	return &verification{}
}

func (v verification) unqualifiedTableName() string {
	return "verifications"
}

func (v verification) tableName() string {
	return "zitadel." + v.unqualifiedTableName()
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

const getVerificationStmt = "SELECT verifications.value, verifications.code, verifications.expires_at, verifications.failed_attempts FROM zitadel.verifications"

func (v verification) get(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Verification, error) {
	var builder database.StatementBuilder
	builder.WriteString(getVerificationStmt)
	condition.Write(&builder)
	return getOne[domain.Verification](ctx, client, &builder)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (v verification) SetValue(value string) database.Change {
	return database.NewChange(v.CodeColumn(), value)
}

func (v verification) SetCode(code []byte) database.Change {
	return database.NewChange(v.CodeColumn(), code)
}

func (v verification) setExpiresAt(expiry time.Duration) database.Change {
	if expiry == 0 {
		return database.NewChangeToNull(v.ExpiresAtColumn())
	}
	return database.NewChange(v.ExpiresAtColumn(), expiry)
}

func (v verification) IncreaseFailedAttempts() database.Change {
	return database.NewIncrementColumnChange(v.FailedAttemptsColumn())
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (v verification) PrimaryKeyCondition(instanceID, id string) database.Condition {
	return database.And(
		v.InstanceIDCondition(instanceID),
		v.IDCondition(id),
	)
}

func (v verification) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(v.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (v verification) IDCondition(id string) database.Condition {
	return database.NewTextCondition(v.IDColumn(), database.TextOperationEqual, id)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (v verification) PrimaryKeyColumns() []database.Column {
	return database.Columns{
		v.IDColumn(),
		v.InstanceIDColumn(),
	}
}

func (v verification) InstanceIDColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "instance_id")
}

func (v verification) IDColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "id")
}

func (v verification) ValueColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "value")
}

func (v verification) CodeColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "code")
}

func (v verification) ExpiresAtColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "expires_at")
}

func (v verification) FailedAttemptsColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "failed_attempts")
}

func (v verification) CreationDateColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "creation_date")
}

func (v verification) UpdatedAtColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "change_date")
}
