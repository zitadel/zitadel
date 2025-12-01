package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// SetEmail implements [domain.HumanUserRepository].
func (u userHuman) SetEmail(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification domain.VerificationType) (int64, error) {
	switch v := verification.(type) {
	case *domain.VerificationTypeVerified:
		return u.setEmailVerified(ctx, client, condition, v)
	case *domain.VerificationTypeSkipped:
		return u.setEmailSkipVerification(ctx, client, condition, v)
	case *domain.VerificationTypeInit:
		return u.initEmail(ctx, client, condition, v)
	case *domain.VerificationTypeUpdate:
		return u.updateEmailVerification(ctx, client, condition, v)
	}
	panic("unknown verification type")
}

// GetEmailVerification implements [domain.HumanUserRepository].
func (u userHuman) GetEmailVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Verification, error) {
	return u.verification.getVerification(ctx, client, database.Exists(
		u.unqualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
			database.NewColumnCondition(u.emailVerificationIDColumn(), u.verification.IDColumn()),
			condition,
		),
	))
}

// SetEmailOTPCheck implements [domain.HumanUserRepository].
func (u userHuman) SetEmailOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition, check domain.CheckType) (int64, error) {
	var builder database.StatementBuilder
	switch c := check.(type) {
	case *domain.CheckTypeSucceeded:
		builder.WriteString("UPDATE zitadel.human_users SET ")
		database.NewChanges(
			u.SetUpdatedAt(c.SucceededAt),
			u.setLastSuccessfulEmailOTPCheck(c.SucceededAt),
			u.clearEmailOTPVerificationID(),
		).Write(&builder)
		writeCondition(&builder, condition)
	case *domain.CheckTypeFailed:
		builder.WriteString("UPDATE zitadel.verifications SET ")
		database.NewIncrementColumnChange(u.verification.FailedAttemptsColumn()).Write(&builder)
		builder.WriteString(" FROM zitadel.human_users WHERE ")
		database.And(
			database.NewColumnCondition(u.verification.InstanceIDColumn(), u.InstanceIDColumn()),
			database.NewColumnCondition(u.verification.IDColumn(), u.emailOTPVerificationIDColumn()),
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
		builder.WriteString(") UPDATE zitadel.human_users SET email_otp_verification_id = (SELECT id FROM created_verification) ")
		builder.WriteString("WHERE (instance_id, id) IN (SELECT instance_id, user_id FROM found_user)")
	}

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// GetEmailOTPCheck implements [domain.HumanUserRepository].
func (u userHuman) GetEmailOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Check, error) {
	return u.verification.getCheck(ctx, client, database.Exists(
		u.unqualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
			database.NewColumnCondition(u.emailOTPVerificationIDColumn(), u.verification.IDColumn()),
			condition,
		),
	))
}

func (u userHuman) setEmailVerified(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeVerified) (int64, error) {
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
		database.NewColumnCondition(u.emailVerificationIDColumn(), u.verification.IDColumn()),
	).Write(&builder)
	writeCondition(&builder, condition)

	builder.WriteString(") UPDATE zitadel.human_users SET ")
	database.NewChanges(
		u.SetUpdatedAt(verification.VerifiedAt),
		database.NewChangeToColumn(u.EmailColumn(), database.NewColumn("found_verification", "value")),
		u.setEmailVerifiedAt(verification.VerifiedAt),
		u.clearEmailVerificationID(),
	).Write(&builder)
	builder.WriteString(" FROM zitadel.found_verification")
	writeCondition(&builder, database.And(
		database.NewColumnCondition(u.InstanceIDColumn(), database.NewColumn("found_verification", "instance_id")),
		database.NewColumnCondition(u.emailVerificationIDColumn(), database.NewColumn("found_verification", "id")),
	))

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) setEmailSkipVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeSkipped) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString("UPDATE zitadel.human_users SET ")
	database.NewChanges(
		u.SetUpdatedAt(verification.SkippedAt),
		u.setEmail(*verification.Value),
		u.setEmailVerifiedAt(verification.SkippedAt),
		u.clearEmailVerificationID(),
	).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) initEmail(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeInit) (int64, error) {
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
			"UPDATE zitadel.human_users SET email_verification_id = (SELECT id FROM verification) WHERE (instance_id, id) IN (SELECT instance_id, id FROM found_user)",
	)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) updateEmailVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeUpdate) (int64, error) {
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
		database.NewColumnCondition(u.emailVerificationIDColumn(), u.verification.IDColumn()),
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

func (u userHuman) setEmail(email string) database.Change {
	return database.NewChange(u.EmailColumn(), email)
}

func (u userHuman) setEmailVerifiedAt(verifiedAt time.Time) database.Change {
	if verifiedAt.IsZero() {
		return database.NewChange(u.EmailVerifiedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(u.EmailVerifiedAtColumn(), verifiedAt)
}

func (u userHuman) clearEmailVerificationID() database.Change {
	return database.NewChangeToNull(u.emailVerificationIDColumn())
}

func (u userHuman) clearEmailOTPVerificationID() database.Change {
	return database.NewChangeToNull(u.emailOTPVerificationIDColumn())
}

func (u userHuman) setLastSuccessfulEmailOTPCheck(succeededAt time.Time) database.Change {
	return database.NewChange(u.lastSuccessfulEmailOTPCheckColumn(), succeededAt)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (u userHuman) EmailCondition(op database.TextOperation, email string) database.Condition {
	return database.NewTextCondition(u.EmailColumn(), op, email)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) EmailVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_verified_at")
}
func (u userHuman) EmailColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email")
}

func (u userHuman) emailVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_verification_id")
}

func (u userHuman) emailOTPVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_otp_verification_id")
}

func (u userHuman) lastSuccessfulEmailOTPCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "last_successful_email_otp_check")
}
