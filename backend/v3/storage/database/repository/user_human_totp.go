package repository

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetTOTPSecret implements [domain.HumanUserRepository.SetTOTPSecret].
func (u userHuman) SetTOTPSecret(secret []byte) database.Change {
	return database.NewChange(u.totpSecretColumn(), secret)
}

// SetTOTPVerifiedAt implements [domain.HumanUserRepository.SetTOTPVerifiedAt].
func (u userHuman) SetTOTPVerifiedAt(verifiedAt time.Time) database.Change {
	if verifiedAt.IsZero() {
		return database.NewChange(u.totpVerifiedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(u.totpVerifiedAtColumn(), verifiedAt)
}

// RemoveTOTP implements [domain.HumanUserRepository.RemoveTOTP].
func (u userHuman) RemoveTOTP() database.Change {
	return database.NewChanges(
		database.NewChangeToNull(u.totpSecretColumn()),
		database.NewChangeToNull(u.totpVerifiedAtColumn()),
		database.NewChangeToNull(u.lastSuccessfulTOTPCheckColumn()),
	)
}

// SetLastSuccessfulTOTPCheck implements [domain.HumanUserRepository.SetLastSuccessfulTOTPCheck].
func (u userHuman) SetLastSuccessfulTOTPCheck(checkedAt time.Time) database.Change {
	if checkedAt.IsZero() {
		return database.NewChange(u.lastSuccessfulTOTPCheckColumn(), database.NowInstruction)
	}
	return database.NewChange(u.lastSuccessfulTOTPCheckColumn(), checkedAt)
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
