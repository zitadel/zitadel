package repository

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetTOTPSecret implements [domain.HumanUserRepository].
func (u userHuman) SetTOTPSecret(secret *crypto.CryptoValue) database.Change {
	return database.NewChange(u.totpSecretColumn(), secret)
}

// SetTOTPVerifiedAt implements [domain.HumanUserRepository].
func (u userHuman) SetTOTPVerifiedAt(verifiedAt time.Time) database.Change {
	if verifiedAt.IsZero() {
		return database.NewChange(u.totpVerifiedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(u.totpVerifiedAtColumn(), verifiedAt)
}

// RemoveTOTP implements [domain.HumanUserRepository].
func (u userHuman) RemoveTOTP() database.Change {
	return database.NewChanges(
		database.NewChangeToNull(u.totpSecretColumn()),
		database.NewChangeToNull(u.totpVerifiedAtColumn()),
		database.NewChangeToNull(u.lastSuccessfulTOTPCheckColumn()),
	)
}

// SetLastSuccessfulTOTPCheck implements [domain.HumanUserRepository].
func (u userHuman) SetLastSuccessfulTOTPCheck(checkedAt time.Time) database.Change {
	if checkedAt.IsZero() {
		return database.NewChanges(
			database.NewChange(u.lastSuccessfulTOTPCheckColumn(), database.NowInstruction),
			u.ResetTOTPFailedAttempts(),
		)
	}
	return database.NewChanges(
		database.NewChange(u.lastSuccessfulTOTPCheckColumn(), checkedAt),
		u.ResetTOTPFailedAttempts(),
	)
}

// IncrementTOTPFailedAttempts implements [domain.HumanUserRepository].
func (u userHuman) IncrementTOTPFailedAttempts() database.Change {
	return database.NewIncrementColumnChange(u.totpFailedAttemptsColumn(), database.Coalesce(u.totpFailedAttemptsColumn(), 0))
}

// ResetTOTPFailedAttempts implements [domain.HumanUserRepository].
func (u userHuman) ResetTOTPFailedAttempts() database.Change {
	return database.NewChange(u.totpFailedAttemptsColumn(), 0)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) totpSecretColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_secret")
}

func (u userHuman) totpVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_verified_at")
}

func (u userHuman) lastSuccessfulTOTPCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_last_successful_check")
}

func (u userHuman) totpFailedAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_failed_attempts")
}
