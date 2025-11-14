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
	panic("unimplemented")
}

// DisableEmailOTP implements [domain.HumanUserRepository.DisableEmailOTP].
func (u userHuman) DisableEmailOTP() database.Change {
	return database.NewChangeToNull(u.emailOTPEnabledAtColumn())
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
