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

// SetPhone implements [domain.HumanUserRepository].
func (u userHuman) SetPhone(phone string) database.Change {
	return database.NewChange(u.phoneColumn(), phone)
}

// SetUnverifiedPhone implements [domain.HumanUserRepository].
func (u userHuman) SetUnverifiedPhone(phone string) database.Change {
	return database.NewChange(u.unverifiedPhoneColumn(), phone)
}

// SetPhoneVerification implements [domain.HumanUserRepository].
func (u userHuman) SetPhoneVerification(verification domain.VerificationType) database.Change {
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		return u.verification.init(typ, existingHumanUser.unqualifiedTableName(), existingHumanUser.phoneVerificationIDColumn())
	case *domain.VerificationTypeUpdate:
		return u.verification.update(typ, existingHumanUser.unqualifiedTableName(),
			existingHumanUser.InstanceIDColumn(), existingHumanUser.phoneVerificationIDColumn(),
		)
	case *domain.VerificationTypeSucceeded:
		return database.NewChanges(
			u.verification.verified(typ, existingHumanUser.unqualifiedTableName(), existingHumanUser.InstanceIDColumn(),
				existingHumanUser.phoneVerificationIDColumn(), u.phoneVerifiedAtColumn(), u.failedSMSOTPAttemptsColumn(),
			),
			database.NewChangeToColumn(u.phoneColumn(), u.unverifiedPhoneColumn()),
		)
	case *domain.VerificationTypeFailed:
		return u.verification.failed(existingHumanUser.unqualifiedTableName(), existingHumanUser.InstanceIDColumn(), existingHumanUser.phoneVerificationIDColumn())
	case *domain.VerificationTypeSkipped:
		return u.verification.skipped(typ, u.phoneVerifiedAtColumn(), u.phoneVerificationIDColumn(), u.failedSMSOTPAttemptsColumn())
	}
	panic(fmt.Sprintf("type not allowed for phone verification change %T", verification))
}

// RemovePhone implements [domain.HumanUserRepository].
func (u userHuman) RemovePhone() database.Change {
	return database.NewChanges(
		database.NewChangeToNull(u.phoneColumn()),
		database.NewChangeToNull(u.phoneVerifiedAtColumn()),
		database.NewChangeToNull(u.phoneVerificationIDColumn()),
		database.NewChangeToNull(u.smsOTPEnabledAtColumn()),
		database.NewChangeToNull(u.lastSuccessfulSMSOTPCheckColumn()),
		database.NewChangeToNull(u.failedSMSOTPAttemptsColumn()),
	)
}

// EnableSMSOTPAt implements [domain.HumanUserRepository].
func (u userHuman) EnableSMSOTPAt(enabledAt time.Time) database.Change {
	return database.NewChange(u.smsOTPEnabledAtColumn(), enabledAt)
}

// EnableSMSOTP implements [domain.HumanUserRepository].
func (u userHuman) EnableSMSOTP() database.Change {
	return database.NewChange(u.smsOTPEnabledAtColumn(), database.NowInstruction)
}

// DisableSMSOTP implements [domain.HumanUserRepository].
func (u userHuman) DisableSMSOTP() database.Change {
	return database.NewChangeToNull(u.smsOTPEnabledAtColumn())
}

func (u userHuman) SetLastSuccessfulSMSOTPCheck(checkedAt time.Time) database.Change {
	if checkedAt.IsZero() {
		return database.NewChanges(
			database.NewChange(u.lastSuccessfulSMSOTPCheckColumn(), database.NowInstruction),
			u.ResetSMSOTPFailedAttempts(),
		)
	}
	return database.NewChanges(
		database.NewChange(u.lastSuccessfulSMSOTPCheckColumn(), checkedAt),
		u.ResetSMSOTPFailedAttempts(),
	)
}

// IncrementSMSOTPFailedAttempts implements [domain.HumanUserRepository].
func (u userHuman) IncrementSMSOTPFailedAttempts() database.Change {
	return database.NewIncrementColumnChange(u.failedSMSOTPAttemptsColumn(), database.Coalesce(u.failedSMSOTPAttemptsColumn(), 0))
}

// ResetSMSOTPFailedAttempts implements [domain.HumanUserRepository].
func (u userHuman) ResetSMSOTPFailedAttempts() database.Change {
	return database.NewChange(u.failedSMSOTPAttemptsColumn(), 0)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PhoneCondition implements [domain.HumanUserRepository].
func (u userHuman) PhoneCondition(op database.TextOperation, phone string) database.Condition {
	return database.Or(u.UnverifiedPhoneCondition(op, phone), u.VerifiedPhoneCondition(op, phone))
}

// UnverifiedPhoneCondition implements [domain.HumanUserRepository].
func (u userHuman) UnverifiedPhoneCondition(op database.TextOperation, phone string) database.Condition {
	return database.NewTextCondition(u.unverifiedPhoneColumn(), op, phone)
}

// VerifiedPhoneCondition implements [domain.HumanUserRepository].
func (u userHuman) VerifiedPhoneCondition(op database.TextOperation, phone string) database.Condition {
	return database.NewTextCondition(u.phoneColumn(), op, phone)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) phoneColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "phone")
}

func (u userHuman) unverifiedPhoneColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "unverified_phone")
}

func (u userHuman) phoneVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "phone_verified_at")
}

func (u userHuman) phoneVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "phone_verification_id")
}

func (u userHuman) smsOTPEnabledAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "sms_otp_enabled_at")
}

func (u userHuman) lastSuccessfulSMSOTPCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "sms_otp_last_successful_check")
}

func (u userHuman) failedSMSOTPAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "sms_otp_failed_attempts")
}
