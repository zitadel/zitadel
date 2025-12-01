package repository

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// CheckEmailOTP implements [domain.HumanUserRepository.CheckEmailOTP].
func (u userHuman) CheckEmailOTP(check domain.CheckType) database.Change {
	switch typ := check.(type) {
	case *domain.CheckTypeFailed:
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.verification.qualifiedTableName())
			builder.WriteString(" SET ")
			database.NewIncrementColumnChange(u.verification.FailedAttemptsColumn())
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(u.verification.InstanceIDColumn(), existingHumanUser.instanceIDColumn()),
				database.NewColumnCondition(u.verification.IDColumn(), existingHumanUser.emailOTPVerificationIDColumn()),
			))
		}, nil)
	case *domain.CheckTypeSucceeded:
		lastSucceededChange := database.NewChange(u.lastSuccessfulEmailOTPCheckColumn(), database.NowInstruction)
		if !typ.SucceededAt.IsZero() {
			lastSucceededChange = database.NewChange(u.lastSuccessfulEmailOTPCheckColumn(), typ.SucceededAt)
		}
		return database.NewChanges(
			lastSucceededChange,
			database.NewChangeToNull(u.emailOTPVerificationIDColumn()),
		)
	case *domain.CheckTypeInit:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.verifications ()")
			}
		)
	}

}

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
	panic("unimplemented")
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

func (u userHuman) emailOTPVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_otp_verification_id")
}

func (u userHuman) lastSuccessfulEmailOTPCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_otp_last_successfully_checked_at")
}
