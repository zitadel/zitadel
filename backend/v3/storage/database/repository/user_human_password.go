package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func (u userHuman) SetPassword(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification domain.VerificationType) (int64, error) {
	switch v := verification.(type) {
	case *domain.VerificationTypeIsVerified:
		return u.setPasswordVerified(ctx, client, condition, v)
	case *domain.VerificationTypeSet:
		return u.setPasswordVerificationCode(ctx, client, condition, v)
	case *domain.VerificationTypeVerifiedValue:
		return u.updatePasswordVerifiedValue(ctx, client, condition, v)
	}
	panic("unknown verification type")
}

func (u userHuman) setPasswordVerified(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeIsVerified) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString("UPDATE zitadel.human_users SET ")
	database.NewChanges(
		u.SetUpdatedAt(verification.VerifiedAt),
		u.setPasswordVerifiedAt(verification.VerifiedAt),
		u.clearPasswordVerificationID(),
	).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) setPasswordVerificationCode(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeSet) (int64, error) {
	var builder database.StatementBuilder
	builder.WriteString("WITH found_user AS ( SELECT * FROM zitadel.human_users")
	writeCondition(&builder, condition)
	builder.WriteString(" FOR NO KEY UPDATE), verification AS (" +
		" INSERT INTO zitadel.verifications (instance_id, code, expires_at, created_at, value)" +
		" SELECT u.instance_id, ",
	)
	builder.WriteArgs(
		verification.Code,
		verification.ExpiresAt,
		verification.CreatedAt,
		verification.Value,
	)
	builder.WriteString(
		" FROM found_user u" +
			" RETURNING id" +
			") " +
			"UPDATE zitadel.human_users SET password_verification_id = (SELECT id FROM verification) WHERE (instance_id, id) IN (SELECT instance_id, id FROM found_user)",
	)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) updatePasswordVerifiedValue(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeVerifiedValue) (int64, error) {
	var builder database.StatementBuilder

	var (
		updatedAt  database.Change
		verifiedAt database.Change
	)

	builder.WriteString("UPDATE zitadel.human_users SET ")
	database.NewChanges(
		updatedAt,
		verifiedAt,
		u.clearPasswordVerificationID(),
		u.setPassword(verification.Value),
	).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// // GetPasswordVerification implements [domain.HumanUserRepository].
// func (u userHuman) GetPasswordVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Verification, error) {
// 	return u.verification.get(ctx, client, database.Exists(
// 		u.unqualifiedTableName(),
// 		database.And(
// 			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
// 			database.NewColumnCondition(u.passwordVerificationIDColumn(), u.verification.IDColumn()),
// 			condition,
// 		),
// 	))
// }

// // ResetPassword implements [domain.HumanUserRepository].
// func (u userHuman) ResetPassword(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.Verification) (int64, error) {
// 	var builder database.StatementBuilder
// 	builder.WriteString("WITH found_user AS ( SELECT * FROM zitadel.human_users")
// 	writeCondition(&builder, condition)
// 	builder.WriteString(" FOR NO KEY UPDATE), verification AS (" +
// 		" INSERT INTO zitadel.verifications (instance_id, code, expires_at)" +
// 		" SELECT u.instance_id, ",
// 	)
// 	builder.WriteArgs(
// 		verification.Code,
// 		verification.ExpiresAt,
// 	)
// 	builder.WriteString(
// 		" FROM found_user u" +
// 			" RETURNING id" +
// 			") " +
// 			"UPDATE zitadel.human_users SET password_verification_id = (SELECT id FROM verification) WHERE (instance_id, id) IN (SELECT instance_id, id FROM found_user)",
// 	)

// 	return client.Exec(ctx, builder.String(), builder.Args()...)
// }

// func (u userHuman) IncrementFailedPasswordVerificationAttempts(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
// 	var builder database.StatementBuilder
// 	builder.WriteString("UPDATE zitadel.verifications SET failed_attempts = failed_attempts + 1 WHERE (instance_id, id) = (SELECT instance_id, unverified_password_id FROM zitadel.human_users")
// 	writeCondition(&builder, condition)
// 	builder.WriteRune(')')
// 	return client.Exec(ctx, builder.String(), builder.Args()...)
// }

// // SetPasswordVerified implements [domain.HumanUserRepository].
// func (u userHuman) VerifyPassword(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
// 	return u.VerifyPasswordAt(ctx, client, condition, time.Time{})
// }

// // SetPasswordVerifiedAt implements [domain.HumanUserRepository].
// func (u userHuman) VerifyPasswordAt(ctx context.Context, client database.QueryExecutor, condition database.Condition, verifiedAt time.Time) (int64, error) {
// 	return u.verifyPassword(ctx, client, condition, u.setPasswordVerifiedAt(verifiedAt))
// }

// func (u userHuman) verifyPassword(ctx context.Context, client database.QueryExecutor, condition database.Condition, verifiedAt database.Change) (int64, error) {
// 	var builder database.StatementBuilder
// 	builder.WriteString("UPDATE zitadel.human_users SET ")
// 	database.NewChanges(
// 		verifiedAt,
// 		u.clearPasswordVerificationID(),
// 	).Write(&builder)
// 	builder.WriteString(" FROM SELECT value FROM zitadel.verifications WHERE ")
// 	database.Exists(
// 		u.unqualifiedTableName(),
// 		database.And(
// 			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
// 			database.NewColumnCondition(u.passwordVerificationIDColumn(), u.verification.IDColumn()),
// 			condition,
// 		),
// 	).Write(&builder)

// 	return client.Exec(ctx, builder.String(), builder.Args()...)
// }

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetPasswordChangeRequired implements [domain.HumanUserRepository].
func (u userHuman) SetPasswordChangeRequired(required bool) database.Change {
	return database.NewChange(u.passwordChangeRequiredColumn(), required)
}

// IncrementFailedPasswordAttempts implements [domain.HumanUserRepository].
func (u userHuman) IncrementFailedPasswordAttempts() database.Change {
	return database.NewIncrementColumnChange(u.failedPasswordAttemptsColumn())
}

// ResetFailedPasswordAttempts implements [domain.HumanUserRepository].
func (u userHuman) ResetFailedPasswordAttempts() database.Change {
	return database.NewChange(u.failedPasswordAttemptsColumn(), 0)
}

func (u userHuman) setPassword(password string) database.Change {
	return database.NewChange(u.passwordColumn(), password)
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

func (u userHuman) passwordChangeRequiredColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_change_required")
}

func (u userHuman) failedPasswordAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "failed_password_attempts")
}

func (u userHuman) passwordVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_verification_id")
}

func (u userHuman) passwordColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password")
}
