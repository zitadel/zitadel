package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// SetTOTP implements [domain.HumanUserRepository.SetTOTP].
func (u userHuman) SetTOTP(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification domain.VerificationType) (int64, error) {
	switch v := verification.(type) {
	case *domain.VerificationTypeVerified:
		return u.setTOTPVerified(ctx, client, condition, v)
	case *domain.VerificationTypeSkipped:
		return u.setTOTPSkipVerification(ctx, client, condition, v)
	case *domain.VerificationTypeInit:
		return u.initTOTP(ctx, client, condition, v)
	case *domain.VerificationTypeUpdate:
		return u.updateTOTPVerification(ctx, client, condition, v)
	}
	panic("unknown verification type")
}

// GetTOTPVerification implements [domain.HumanUserRepository.GetTOTPVerification].
func (u userHuman) GetTOTPVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Verification, error) {
	return u.verification.getVerification(ctx, client, database.Exists(
		u.unqualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
			database.NewColumnCondition(u.totpVerificationIDColumn(), u.verification.IDColumn()),
			condition,
		),
	))
}

// SetTOTPCheck implements [domain.HumanUserRepository.SetTOTPCheck].
func (u userHuman) SetTOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition, check domain.CheckType) (int64, error) {
	var builder database.StatementBuilder
	switch c := check.(type) {
	case *domain.CheckTypeSucceeded:
		builder.WriteString("UPDATE zitadel.human_users SET ")
		database.NewChanges(
			u.SetUpdatedAt(c.SucceededAt),
			u.setLastSuccessfulTOTPCheck(c.SucceededAt),
			u.clearTOTPVerificationID(),
		).Write(&builder)
		writeCondition(&builder, condition)
	case *domain.CheckTypeFailed:
		builder.WriteString("UPDATE zitadel.verifications SET ")
		database.NewIncrementColumnChange(u.verification.FailedAttemptsColumn()).Write(&builder)
		builder.WriteString(" FROM zitadel.human_users WHERE ")
		database.And(
			database.NewColumnCondition(u.verification.InstanceIDColumn(), u.InstanceIDColumn()),
			database.NewColumnCondition(u.verification.IDColumn(), u.totpVerificationIDColumn()),
			condition,
		).Write(&builder)
	case *domain.CheckTypeInit:
		var createdAt any = database.NowInstruction
		if !c.CreatedAt.IsZero() {
			createdAt = c.CreatedAt
		}
		builder.WriteString("WITH found_user AS (SELECT")
		builder.WriteString(" instance_id, user_id FROM zitadel.human_users")
		writeCondition(&builder, condition)
		builder.WriteString("), created_verification AS (")
		builder.WriteString("INSERT INTO zitadel.verifications (instance_id, code, created_at, expiry) ")
		builder.WriteString("SELECT u.instance_id, ")
		builder.WriteArgs(c.Code, createdAt, c.Expiry)
		builder.WriteString(" FROM found_user u RETURNING id")
		builder.WriteString(") UPDATE zitadel.human_users SET totp_verification_id = (SELECT id FROM created_verification) ")
		builder.WriteString("WHERE (instance_id, id) IN (SELECT instance_id, user_id FROM found_user)")
	}

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// GetTOTPCheck implements [domain.HumanUserRepository.GetTOTPCheck].
func (u userHuman) GetTOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Check, error) {
	return u.verification.getCheck(ctx, client, database.Exists(
		u.unqualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
			database.NewColumnCondition(u.totpVerificationIDColumn(), u.verification.IDColumn()),
			condition,
		),
	))
}

func (u userHuman) setTOTPVerified(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeVerified) (int64, error) {
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
		database.NewColumnCondition(u.totpVerificationIDColumn(), u.verification.IDColumn()),
	).Write(&builder)
	writeCondition(&builder, condition)

	builder.WriteString(") UPDATE zitadel.human_users SET ")
	database.NewChanges(
		u.SetUpdatedAt(verification.VerifiedAt),
		database.NewChangeToColumn(u.TOTPColumn(), database.NewColumn("found_verification", "value")),
		u.setTOTPVerifiedAt(verification.VerifiedAt),
		u.clearTOTPVerificationID(),
	).Write(&builder)
	builder.WriteString(" FROM zitadel.found_verification")
	writeCondition(&builder, database.And(
		database.NewColumnCondition(u.InstanceIDColumn(), database.NewColumn("found_verification", "instance_id")),
		database.NewColumnCondition(u.totpVerificationIDColumn(), database.NewColumn("found_verification", "id")),
	))

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) setTOTPSkipVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeSkipped) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString("UPDATE zitadel.human_users SET ")
	database.NewChanges(
		u.SetUpdatedAt(verification.VerifiedAt),
		u.setTOTP(*verification.Value),
		u.setTOTPVerifiedAt(verification.VerifiedAt),
		u.clearTOTPVerificationID(),
	).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) initTOTP(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeInit) (int64, error) {
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
			"UPDATE zitadel.human_users SET totp_verification_id = (SELECT id FROM verification) WHERE (instance_id, id) IN (SELECT instance_id, id FROM found_user)",
	)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) updateTOTPVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeUpdate) (int64, error) {
	var builder database.StatementBuilder

	changes := make(database.Changes, 0, 3)
	if verification.Code != nil {
		changes = append(changes, u.verification.SetCode(*verification.Code))
	}
	if verification.Value != nil {
		changes = append(changes, u.verification.SetValue(*verification.Value))
	}
	if verification.Expiry != nil {
		changes = append(changes, u.verification.setExpiry(*verification.Expiry))
	}

	builder.WriteString("WITH found_verification AS ( SELECT verifications.* FROM zitadel.human_users JOIN zitadel.verifications ON ")
	database.And(
		database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
		database.NewColumnCondition(u.totpVerificationIDColumn(), u.verification.IDColumn()),
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

// SetMFAInitSkippedAt implements [domain.HumanUserRepository].
func (u userHuman) RemoveTOTP() database.Change {
	return database.NewChanges(
		database.NewChangeToNull(u.totpSecretColumn()),
		database.NewChangeToNull(u.totpVerifiedAtColumn()),
		database.NewChangeToNull(u.failedTOTPAttemptsColumn()),
	)
}

// SetTOTPEnabled implements [domain.HumanUserRepository.SetTOTPEnabled].
func (u userHuman) SetTOTPEnabled(enabled bool) database.Change {
	return database.NewChange(u.totpEnabledColumn(), enabled)
}

func (u userHuman) setTOTP(totp string) database.Change {
	return database.NewChange(u.TOTPColumn(), totp)
}

func (u userHuman) setTOTPVerifiedAt(verifiedAt time.Time) database.Change {
	if verifiedAt.IsZero() {
		return database.NewChange(u.TOTPVerifiedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(u.TOTPVerifiedAtColumn(), verifiedAt)
}

func (u userHuman) clearTOTPVerificationID() database.Change {
	return database.NewChangeToNull(u.totpVerificationIDColumn())
}

func (u userHuman) setLastSuccessfulTOTPCheck(succeededAt time.Time) database.Change {
	return database.NewChange(u.lastSuccessfulTOTPCheckColumn(), succeededAt)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// TOTPVerifiedAtColumn implements [domain.HumanUserRepository.TOTPVerifiedAtColumn].
func (u userHuman) TOTPVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_verified_at")
}

// TOTPColumn implements [domain.HumanUserRepository.TOTPColumn].
func (u userHuman) TOTPColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp")
}

func (u userHuman) totpVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_verification_id")
}

func (u userHuman) lastSuccessfulTOTPCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "last_successful_totp_check")
}

func (u userHuman) totpEnabledColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_enabled")
}

func (u userHuman) totpSecretColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_secret")
}

func (u userHuman) totpVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_verified_at")
}

func (u userHuman) failedTOTPAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "failed_totp_attempts")
}
