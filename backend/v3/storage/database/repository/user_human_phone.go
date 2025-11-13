package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// SetPhone implements [domain.HumanUserRepository.SetPhone].
func (u userHuman) SetPhone(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification domain.VerificationType) (int64, error) {
	switch v := verification.(type) {
	case *domain.VerificationTypeVerified:
		return u.setPhoneVerified(ctx, client, condition, v)
	case *domain.VerificationTypeSkipped:
		return u.setPhoneSkipVerification(ctx, client, condition, v)
	case *domain.VerificationTypeInit:
		return u.initPhone(ctx, client, condition, v)
	case *domain.VerificationTypeUpdate:
		return u.updatePhoneVerification(ctx, client, condition, v)
	}
	panic("unknown verification type")
}

// GetPhoneVerification implements [domain.HumanUserRepository.GetPhoneVerification].
func (u userHuman) GetPhoneVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Verification, error) {
	return u.verification.getVerification(ctx, client, database.Exists(
		u.unqualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
			database.NewColumnCondition(u.phoneVerificationIDColumn(), u.verification.IDColumn()),
			condition,
		),
	))
}

// SetSMSOTPCheck implements [domain.HumanUserRepository.SetSMSOTPCheck].
func (u userHuman) SetSMSOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition, check domain.CheckType) (int64, error) {
	var builder database.StatementBuilder
	switch c := check.(type) {
	case *domain.CheckTypeSucceeded:
		builder.WriteString("UPDATE zitadel.human_users SET ")
		database.NewChanges(
			u.SetUpdatedAt(c.SucceededAt),
			u.setLastSuccessfulSMSOTPCheck(c.SucceededAt),
			u.clearSMSOTPVerificationID(),
		).Write(&builder)
		writeCondition(&builder, condition)
	case *domain.CheckTypeFailed:
		builder.WriteString("UPDATE zitadel.verifications SET ")
		database.NewIncrementColumnChange(u.verification.FailedAttemptsColumn()).Write(&builder)
		builder.WriteString(" FROM zitadel.human_users WHERE ")
		database.And(
			database.NewColumnCondition(u.verification.InstanceIDColumn(), u.InstanceIDColumn()),
			database.NewColumnCondition(u.verification.IDColumn(), u.smsOTPVerificationIDColumn()),
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
		builder.WriteString(") UPDATE zitadel.human_users SET sms_otp_verification_id = (SELECT id FROM created_verification) ")
		builder.WriteString("WHERE (instance_id, id) IN (SELECT instance_id, user_id FROM found_user)")
	}

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// GetSMSOTPCheck implements [domain.HumanUserRepository.GetSMSOTPCheck].
func (u userHuman) GetSMSOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Check, error) {
	return u.verification.getCheck(ctx, client, database.Exists(
		u.unqualifiedTableName(),
		database.And(
			database.NewColumnCondition(u.InstanceIDColumn(), u.verification.InstanceIDColumn()),
			database.NewColumnCondition(u.smsOTPVerificationIDColumn(), u.verification.IDColumn()),
			condition,
		),
	))
}

func (u userHuman) setPhoneVerified(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeVerified) (int64, error) {
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
		database.NewColumnCondition(u.phoneVerificationIDColumn(), u.verification.IDColumn()),
	).Write(&builder)
	writeCondition(&builder, condition)

	builder.WriteString(") UPDATE zitadel.human_users SET ")
	database.NewChanges(
		u.SetUpdatedAt(verification.VerifiedAt),
		database.NewChangeToColumn(u.PhoneColumn(), database.NewColumn("found_verification", "value")),
		u.setPhoneVerifiedAt(verification.VerifiedAt),
		u.clearPhoneVerificationID(),
	).Write(&builder)
	builder.WriteString(" FROM zitadel.found_verification")
	writeCondition(&builder, database.And(
		database.NewColumnCondition(u.InstanceIDColumn(), database.NewColumn("found_verification", "instance_id")),
		database.NewColumnCondition(u.phoneVerificationIDColumn(), database.NewColumn("found_verification", "id")),
	))

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) setPhoneSkipVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeSkipped) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString("UPDATE zitadel.human_users SET ")
	database.NewChanges(
		u.SetUpdatedAt(verification.SkippedAt),
		u.setPhone(*verification.Value),
		u.setPhoneVerifiedAt(verification.SkippedAt),
		u.clearPhoneVerificationID(),
	).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) initPhone(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeInit) (int64, error) {
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
			"UPDATE zitadel.human_users SET phone_verification_id = (SELECT id FROM verification) WHERE (instance_id, id) IN (SELECT instance_id, id FROM found_user)",
	)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (u userHuman) updatePhoneVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeUpdate) (int64, error) {
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
		database.NewColumnCondition(u.phoneVerificationIDColumn(), u.verification.IDColumn()),
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

// RemovePhone implements [domain.HumanUserRepository.RemovePhone].
func (u userHuman) RemovePhone() database.Change {
	return database.NewChanges(
		database.NewChangeToNull(u.PhoneColumn()),
		database.NewChangeToNull(u.PhoneVerifiedAtColumn()),
		u.clearPhoneVerificationID(),
		u.SetSMSOTPEnabled(false),
		database.NewChangeToNull(u.lastSuccessfulSMSOTPCheckColumn()),
		u.clearSMSOTPVerificationID(),
	)
}

// SetSMSOTPEnabled implements [domain.HumanUserRepository.SetSMSOTPEnabled].
func (u userHuman) SetSMSOTPEnabled(enabled bool) database.Change {
	return database.NewChange(u.smsOTPEnabledColumn(), enabled)
}

func (u userHuman) setPhone(phone string) database.Change {
	return database.NewChange(u.PhoneColumn(), phone)
}

func (u userHuman) setPhoneVerifiedAt(verifiedAt time.Time) database.Change {
	if verifiedAt.IsZero() {
		return database.NewChange(u.PhoneVerifiedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(u.PhoneVerifiedAtColumn(), verifiedAt)
}

func (u userHuman) clearPhoneVerificationID() database.Change {
	return database.NewChangeToNull(u.phoneVerificationIDColumn())
}

func (u userHuman) clearSMSOTPVerificationID() database.Change {
	return database.NewChangeToNull(u.smsOTPVerificationIDColumn())
}

func (u userHuman) setLastSuccessfulSMSOTPCheck(succeededAt time.Time) database.Change {
	return database.NewChange(u.lastSuccessfulSMSOTPCheckColumn(), succeededAt)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PhoneCondition implements [domain.HumanUserRepository.PhoneCondition].
func (u userHuman) PhoneCondition(op database.TextOperation, phone string) database.Condition {
	return database.NewTextCondition(u.PhoneColumn(), op, phone)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PhoneVerifiedAtColumn implements [domain.HumanUserRepository.PhoneVerifiedAtColumn].
func (u userHuman) PhoneVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "phone_verified_at")
}

// PhoneColumn implements [domain.HumanUserRepository.PhoneColumn].
func (u userHuman) PhoneColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "phone")
}

func (u userHuman) phoneVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "phone_verification_id")
}

func (u userHuman) smsOTPVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "sms_otp_verification_id")
}

func (u userHuman) lastSuccessfulSMSOTPCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "last_successful_sms_otp_check")
}

func (u userHuman) smsOTPEnabledColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "sms_otp_enabled")
}
