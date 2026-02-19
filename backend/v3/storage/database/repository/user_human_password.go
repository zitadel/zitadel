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

func (u userHuman) SetPassword(password string) database.Change {
	return database.NewChange(u.passwordHashColumn(), password)
}

// SetResetPasswordVerification implements [domain.HumanUserRepository].
func (u userHuman) SetResetPasswordVerification(verification domain.VerificationType) database.Change {
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		return u.verification.init(typ, existingUser.unqualifiedTableName(), existingHumanUser.passwordVerificationIDColumn())
	case *domain.VerificationTypeUpdate:
		return u.verification.update(typ, existingHumanUser.unqualifiedTableName(),
			existingHumanUser.InstanceIDColumn(), existingHumanUser.passwordVerificationIDColumn(),
		)
	case *domain.VerificationTypeSucceeded:
		return u.verification.verified(typ, existingHumanUser.unqualifiedTableName(), existingHumanUser.InstanceIDColumn(),
			existingHumanUser.passwordVerificationIDColumn(), u.passwordVerifiedAtColumn(), u.failedPasswordAttemptsColumn(),
		)
	case *domain.VerificationTypeSkipped:
		return u.verification.skipped(typ, u.passwordVerifiedAtColumn(), u.passwordVerificationIDColumn(), u.failedPasswordAttemptsColumn())
	case *domain.VerificationTypeFailed:
		return u.verification.failed(existingHumanUser.unqualifiedTableName(), existingHumanUser.InstanceIDColumn(), existingHumanUser.passwordVerificationIDColumn())
	}
	panic(fmt.Sprintf("undefined verification type %T", verification))
}

// SetPasswordChangeRequired implements [domain.HumanUserRepository].
func (u userHuman) SetPasswordChangeRequired(required bool) database.Change {
	return database.NewChange(u.passwordChangeRequiredColumn(), required)
}

func (u userHuman) SetLastSuccessfulPasswordCheck(checkedAt time.Time) database.Change {
	if checkedAt.IsZero() {
		return database.NewChanges(
			database.NewChange(u.lastSuccessfulPasswordCheckColumn(), database.NowInstruction),
			u.ResetPasswordFailedAttempts(),
		)
	}
	return database.NewChanges(
		database.NewChange(u.lastSuccessfulPasswordCheckColumn(), checkedAt),
		u.ResetPasswordFailedAttempts(),
	)
}

func (u userHuman) IncrementPasswordFailedAttempts() database.Change {
	return database.NewIncrementColumnChange(u.failedPasswordAttemptsColumn(), database.Coalesce(u.failedPasswordAttemptsColumn(), 0))
}

func (u userHuman) ResetPasswordFailedAttempts() database.Change {
	return database.NewChange(u.failedPasswordAttemptsColumn(), 0)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) passwordHashColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_hash")
}

func (u userHuman) passwordChangeRequiredColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_change_required")
}

func (u userHuman) passwordVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_changed_at")
}

func (u userHuman) passwordVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_verification_id")
}

func (u userHuman) lastSuccessfulPasswordCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_last_successful_check")
}

func (u userHuman) failedPasswordAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_failed_attempts")
}
