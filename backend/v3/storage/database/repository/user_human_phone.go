package repository

import (
	"fmt"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// RemovePhone implements [domain.HumanUserRepository.RemovePhone].
func (u userHuman) RemovePhone() database.Change {
	return database.NewChanges(
		database.NewChangeToNull(u.phoneColumn()),
		database.NewChangeToNull(u.phoneVerifiedAtColumn()),
		database.NewChangeToNull(u.smsOTPEnabledAtColumn()),
		database.NewChangeToNull(u.lastSuccessfulSMSOTPCheckColumn()),
		database.NewChangeToNull(u.smsOTPVerificationIDColumn()),
	)
}

// SetPhone implements [domain.HumanUserRepository.SetPhone].
func (u userHuman) SetPhone(verification domain.VerificationType) database.Change {
	switch typ := verification.(type) {
	case *domain.VerificationTypeSkipped:
		skippedAt := database.NewChange(u.phoneVerifiedAtColumn(), database.NowInstruction)
		if !typ.SkippedAt.IsZero() {
			skippedAt = database.NewChange(u.phoneVerifiedAtColumn(), typ.SkippedAt)
		}
		return database.NewChanges(
			database.NewChange(u.phoneColumn(), *typ.Value),
			skippedAt,
		)
	case *domain.VerificationTypeFailed:
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.verification.qualifiedTableName())
			builder.WriteString(" SET ")
			database.NewIncrementColumnChange(u.verification.FailedAttemptsColumn())
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(u.verification.InstanceIDColumn(), existingHumanUser.InstanceIDColumn()),
				database.NewColumnCondition(u.verification.IDColumn(), existingHumanUser.smsOTPVerificationIDColumn()),
			))
		}, nil)
	case *domain.VerificationTypeInit:
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.verifications(instance_id, user_id, value, code, created_at, expiry) SELECT instance_id, id, ")
			builder.WriteArgs(typ.Value, typ.Code, typ.CreatedAt, typ.Expiry)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			builder.WriteString(" RETURNING id")
		}, func(name string) database.Change {
			return database.NewChangeToStatement(u.phoneVerificationIDColumn(), func(builder *database.StatementBuilder) {
				builder.WriteString("SELECT id FROM ")
				builder.WriteString(name)
			})
		})
	case *domain.VerificationTypeUpdate:
		changes := make(database.Changes, 0, 3)
		if typ.Value != nil {
			changes = append(changes, database.NewChange(u.phoneColumn(), *typ.Value))
		}
		if typ.Code != nil {
			changes = append(changes, database.NewChange(u.verification.CodeColumn(), typ.Code))
		}
		if typ.Expiry != nil {
			changes = append(changes, database.NewChange(u.verification.ExpiryColumn(), *typ.Expiry))
		}
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.verification.qualifiedTableName())
			builder.WriteString(" SET ")
			changes.Write(builder)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(u.verification.InstanceIDColumn(), existingHumanUser.InstanceIDColumn()),
				database.NewColumnCondition(u.verification.IDColumn(), existingHumanUser.phoneVerificationIDColumn()),
			))
		}, nil)
	}
	panic(fmt.Sprintf("type not allowed for phone verification change %T", verification))
}

// CheckSMSOTP implements [domain.HumanUserRepository.CheckSMSOTP].
func (u userHuman) CheckSMSOTP(check domain.CheckType) database.Change {
	panic("unimplemented")
}

// DisableSMSOTP implements [domain.HumanUserRepository.DisableSMSOTP].
func (u userHuman) DisableSMSOTP() database.Change {
	return database.NewChangeToNull(u.smsOTPEnabledAtColumn())
}

// EnableSMSOTP implements [domain.HumanUserRepository.EnableSMSOTP].
func (u userHuman) EnableSMSOTP() database.Change {
	return database.NewChange(u.smsOTPEnabledAtColumn(), database.NowInstruction)
}

// EnableSMSOTPAt implements [domain.HumanUserRepository.EnableSMSOTPAt].
func (u userHuman) EnableSMSOTPAt(enabledAt time.Time) database.Change {
	return database.NewChange(u.smsOTPEnabledAtColumn(), enabledAt)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PhoneCondition implements [domain.HumanUserRepository.PhoneCondition].
func (u userHuman) PhoneCondition(op database.TextOperation, phone string) database.Condition {
	return database.NewTextCondition(u.phoneColumn(), op, phone)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) phoneColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "phone")
}

func (u userHuman) phoneVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "phone_verified_at")
}

func (u userHuman) phoneVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "unverified_phone_id")
}

func (u userHuman) smsOTPEnabledAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "sms_otp_enabled_at")
}

func (u userHuman) lastSuccessfulSMSOTPCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "last_successful_sms_otp_check")
}

func (u userHuman) smsOTPVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "sms_otp_verification_id")
}
