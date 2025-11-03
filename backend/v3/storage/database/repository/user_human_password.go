package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func (u userHuman) SetPassword(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification domain.VerificationType) (int64, error) {
	switch v := verification.(type) {
	case *domain.VerificationTypeSuccessful:
		return u.setPasswordFromSuccessfulVerification(ctx, client, condition, v)
	case *domain.VerificationTypeSkipVerification:
		return u.setPasswordSkipVerification(ctx, client, condition, v)
	case *domain.VerificationTypeInitCode:
		return u.initPasswordVerification(ctx, client, condition, v)
	case *domain.VerificationTypeUpdate:
		return u.updatePasswordVerification(ctx, client, condition, v)
	}
	panic("unknown verification type")
}

func (u userHuman) GetPasswordVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Verification, error) {
	return u.verification.get(ctx, client, database.Exists(
		u.unqualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
			database.NewColumnCondition(u.passwordVerificationIDColumn(), u.verification.IDColumn()),
			condition,
		),
	))
}

func (u userHuman) IncrementPasswordVerificationAttempts(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	var builder database.StatementBuilder
	builder.WriteString("UPDATE zitadel.verifications SET ")
	database.NewIncrementColumnChange(u.verification.FailedAttemptsColumn()).Write(&builder)
	builder.WriteString(" FROM zitadel.human_users WHERE ")
	database.And(
		database.NewColumnCondition(u.verification.InstanceIDColumn(), u.InstanceIDColumn()),
		database.NewColumnCondition(u.verification.IDColumn(), u.passwordVerificationIDColumn()),
		condition,
	).Write(&builder)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) setPasswordFromSuccessfulVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeSuccessful) (int64, error) {
	var builder database.StatementBuilder
	builder.WriteString("WITH found_verification AS (SELECT ")
	database.Columns{
		u.verification.InstanceIDColumn(),
		u.verification.IDColumn(),
		u.verification.ValueColumn(),
	}.WriteQualified(&builder)
	builder.WriteString(" FROM ")
	builder.WriteString(u.qualifiedTableName())
	builder.WriteString(" JOIN ")
	builder.WriteString(u.verification.qualifiedTableName())
	builder.WriteString(" ON ")
	database.And(
		database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
		database.NewColumnCondition(u.passwordVerificationIDColumn(), u.verification.IDColumn()),
	).Write(&builder)
	writeCondition(&builder, condition)

	builder.WriteString(") UPDATE zitadel.human_users SET ")
	database.NewChanges(
		u.SetUpdatedAt(verification.VerifiedAt),
		database.NewChangeToColumn(u.PasswordColumn(), database.NewColumn("found_verification", "value")),
		u.setPasswordVerifiedAt(verification.VerifiedAt),
		u.clearPasswordVerificationID(),
	).Write(&builder)
	builder.WriteString(" FROM zitadel.found_verification")
	writeCondition(&builder, database.And(
		database.NewColumnCondition(u.InstanceIDColumn(), database.NewColumn("found_verification", "instance_id")),
		database.NewColumnCondition(u.passwordVerificationIDColumn(), database.NewColumn("found_verification", "id")),
	))

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) setPasswordSkipVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeSkipVerification) (int64, error) {
	var builder database.StatementBuilder

	var verifiedAt time.Time
	if verification.VerifiedAt != nil {
		verifiedAt = *verification.VerifiedAt
	}

	builder.WriteString("UPDATE zitadel.human_users SET ")
	database.NewChanges(
		u.SetUpdatedAt(verifiedAt),
		u.setPassword(verification.Value),
		u.setPasswordVerifiedAt(verifiedAt),
		u.clearPasswordVerificationID(),
	).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) initPasswordVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeInitCode) (int64, error) {
	var builder database.StatementBuilder

	var createdAt any = database.NowInstruction
	if !verification.CreatedAt.IsZero() {
		createdAt = verification.CreatedAt
	}

	builder.WriteString("WITH found_user AS ( SELECT * FROM zitadel.human_users")
	writeCondition(&builder, condition)
	builder.WriteString(" FOR NO KEY UPDATE), verification AS (" +
		" INSERT INTO zitadel.verifications (instance_id, code, value, created_at, expiry)" +
		" SELECT u.instance_id, ",
	)
	builder.WriteArgs(
		verification.Code,
		verification.Value,
		createdAt,
		verification.Expiry,
	)
	builder.WriteString(
		" FROM found_user u" +
			" RETURNING id" +
			") " +
			"UPDATE zitadel.human_users SET password_verification_id = (SELECT id FROM verification) WHERE (instance_id, id) IN (SELECT instance_id, id FROM found_user)",
	)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) updatePasswordVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeUpdate) (int64, error) {
	var builder database.StatementBuilder

	changes := make(database.Changes, 0, 3)
	changes = append(changes, u.verification.SetCode(verification.Code))
	if verification.Value != nil {
		changes = append(changes, u.verification.SetValue(*verification.Value))
	}
	if verification.Expiry != nil {
		changes = append(changes, u.verification.setExpiry(*verification.Expiry))
	}

	builder.WriteString("WITH found_verification AS ( SELECT verifications.* FROM zitadel.human_users JOIN zitadel.verifications ON ")
	database.And(
		database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
		database.NewColumnCondition(u.passwordVerificationIDColumn(), u.verification.IDColumn()),
	).Write(&builder)
	writeCondition(&builder, condition)

	builder.WriteString(") UPDATE zitadel.verifications SET")
	changes.Write(&builder)
	builder.WriteString(" FROM found_verification")
	writeCondition(&builder, database.And(
		database.NewColumnCondition(u.verification.InstanceIDColumn(), database.NewColumn("found_verification", "instance_id")),
		database.NewColumnCondition(u.verification.IDColumn(), database.NewColumn("found_verification", "id")),
	))

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetPasswordChangeRequired implements [domain.HumanUserRepository].
func (u userHuman) SetPasswordChangeRequired(required bool) database.Change {
	return database.NewChange(u.passwordChangeRequiredColumn(), required)
}

// IncrementFailedPasswordAttempts implements [domain.HumanUserRepository].
func (u userHuman) IncrementFailedPasswordAttempts() database.Change {
	return database.NewIncrementColumnChange(u.FailedPasswordAttemptsColumn())
}

// ResetFailedPasswordAttempts implements [domain.HumanUserRepository].
func (u userHuman) ResetFailedPasswordAttempts() database.Change {
	return database.NewChange(u.FailedPasswordAttemptsColumn(), 0)
}

func (u userHuman) setPassword(password string) database.Change {
	return database.NewChange(u.PasswordColumn(), password)
}

func (u userHuman) setPasswordVerifiedAt(verifiedAt time.Time) database.Change {
	if verifiedAt.IsZero() {
		return database.NewChange(u.PasswordVerifiedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(u.PasswordVerifiedAtColumn(), verifiedAt)
}

func (u userHuman) clearPasswordVerificationID() database.Change {
	return database.NewChangeToNull(u.passwordVerificationIDColumn())
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) PasswordVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_verified_at")
}

func (u userHuman) FailedPasswordAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "failed_password_attempts")
}

func (u userHuman) PasswordColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password")
}

func (u userHuman) passwordChangeRequiredColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_change_required")
}

func (u userHuman) passwordVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_verification_id")
}
