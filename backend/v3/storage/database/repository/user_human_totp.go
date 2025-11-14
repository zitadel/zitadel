package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// CheckTOTP implements [domain.HumanUserRepository.CheckTOTP].
func (u userHuman) CheckTOTP(check domain.CheckType) database.Change {
	panic("unimplemented")
}

// RemoveTOTP implements [domain.HumanUserRepository.RemoveTOTP].
func (u userHuman) RemoveTOTP() database.Change {
	return database.NewChanges(
		database.NewChangeToNull(u.totpSecretIDColumn()),
		database.NewChangeToNull(u.totpVerifiedAtColumn()),
		database.NewChangeToNull(u.lastSuccessfulTOTPCheckColumn()),
	)
}

// SetTOTP implements [domain.HumanUserRepository.SetTOTP].
func (u userHuman) SetTOTP(verification domain.VerificationType) database.Change {
	panic("unimplemented")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) totpSecretIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_secret_id")
}

func (u userHuman) totpVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "totp_verified_at")
}

func (u userHuman) lastSuccessfulTOTPCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "last_successful_totp_check")
}
