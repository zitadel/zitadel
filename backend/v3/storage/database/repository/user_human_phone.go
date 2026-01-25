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
		database.NewChangeToNull(u.failedSMSOTPAttemptsColumn()),
	)
}

// SetPhone implements [domain.HumanUserRepository.SetPhone].
func (u userHuman) SetPhone(verification domain.VerificationType) database.Change {
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		return u.verification.init(typ, existingHumanUser.unqualifiedTableName(), existingHumanUser.phoneVerificationIDColumn())
	case *domain.VerificationTypeSkipped:
		return u.verification.skipped(typ, u.phoneVerifiedAtColumn(), u.phoneColumn())
	case *domain.VerificationTypeUpdate:
		return u.verification.update(typ, existingHumanUser.unqualifiedTableName(),
			existingHumanUser.InstanceIDColumn(), existingHumanUser.phoneVerificationIDColumn(),
		)
	case *domain.VerificationTypeVerified:
		return u.verification.verified(typ, existingUser.unqualifiedTableName(), u.InstanceIDColumn(),
			u.phoneVerificationIDColumn(), u.phoneVerifiedAtColumn(), u.phoneColumn())
	case *domain.VerificationTypeFailed:
		return u.verification.failed(existingHumanUser.unqualifiedTableName(), existingHumanUser.InstanceIDColumn(), existingHumanUser.phoneVerificationIDColumn())
	}
	panic(fmt.Sprintf("type not allowed for phone verification change %T", verification))
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
