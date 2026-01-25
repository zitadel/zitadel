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

// DisableEmailOTP implements [domain.HumanUserRepository.DisableEmailOTP].
func (u userHuman) DisableEmailOTP() database.Change {
	return database.NewChange(u.emailOTPEnabledAtColumn(), false)
}

// EnableEmailOTP implements [domain.HumanUserRepository.EnableEmailOTP].
func (u userHuman) EnableEmailOTP() database.Change {
	return database.NewChange(u.emailOTPEnabledAtColumn(), database.NowInstruction)
}

// EnableEmailOTPAt implements [domain.HumanUserRepository.EnableEmailOTPAt].
func (u userHuman) EnableEmailOTPAt(enabledAt time.Time) database.Change {
	return database.NewChange(u.emailOTPEnabledAtColumn(), enabledAt)
}

// SetEmail implements [domain.HumanUserRepository.SetEmail].
func (u userHuman) SetEmail(verification domain.VerificationType) database.Change {
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		return u.verification.init(typ, existingHumanUser.unqualifiedTableName(), existingHumanUser.emailVerificationIDColumn())
	case *domain.VerificationTypeSkipped:
		return u.verification.skipped(typ, u.emailVerifiedAtColumn(), u.EmailColumn())
	case *domain.VerificationTypeUpdate:
		return u.verification.update(typ, existingHumanUser.unqualifiedTableName(),
			existingHumanUser.InstanceIDColumn(), existingHumanUser.emailVerificationIDColumn(),
		)
	case *domain.VerificationTypeVerified:
		return u.verification.verified(typ, existingUser.unqualifiedTableName(), u.InstanceIDColumn(),
			u.emailVerificationIDColumn(), u.emailVerifiedAtColumn(), u.EmailColumn())
	case *domain.VerificationTypeFailed:
		return u.verification.failed(existingHumanUser.unqualifiedTableName(), existingHumanUser.InstanceIDColumn(), existingHumanUser.emailVerificationIDColumn())
	}
	panic(fmt.Sprintf("type not allowed for email verification change %T", verification))
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// EmailCondition implements [domain.HumanUserRepository.EmailCondition].
func (u userHuman) EmailCondition(op database.TextOperation, email string) database.Condition {
	return database.NewTextCondition(u.EmailColumn(), op, email)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// EmailColumn implements [domain.HumanUserRepository.EmailColumn].
func (u userHuman) EmailColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email")
}

func (u userHuman) emailOTPEnabledAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_otp_enabled_at")
}

func (u userHuman) emailVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_verified_at")
}

func (u userHuman) emailVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_verification_id")
}
