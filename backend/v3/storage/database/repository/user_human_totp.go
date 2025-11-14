package repository

import (
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// CheckTOTP implements [domain.HumanUserRepository.CheckTOTP].
func (u userHuman) CheckTOTP(check domain.CheckType) database.Change {
	switch typ := check.(type) {
	case *domain.CheckTypeFailed:
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.verification.qualifiedTableName())
			builder.WriteString(" SET ")
			database.NewIncrementColumnChange(u.verification.FailedAttemptsColumn())

		}, nil)
	case *domain.CheckTypeSucceeded:
		lastSucceededChange := database.NewChange(u.lastSuccessfulTOTPCheckColumn(), database.NowInstruction)
		if !typ.SucceededAt.IsZero() {
			lastSucceededChange = database.NewChange(u.lastSuccessfulTOTPCheckColumn(), typ.SucceededAt)
		}
		return database.NewChanges(
			// database.NewChange(u.totpFailedAttemptsColumn()),
			database.NewCTEChange(
				func(builder *database.StatementBuilder) {
					builder.WriteString("UPDATE ")
					builder.WriteString(u.verification.qualifiedTableName())
					builder.WriteString(" SET ")
					database.NewChange(u.verification.FailedAttemptsColumn(), 0)
					builder.WriteString(" FROM ")
					builder.WriteString(existingHumanUser.unqualifiedTableName())
					writeCondition(builder, database.And(
						database.NewColumnCondition(u.verification.InstanceIDColumn(), existingHumanUser.instanceIDColumn()),
						database.NewColumnCondition(u.verification.IDColumn(), existingHumanUser.totpSecretIDColumn()),
					))
				},
				nil,
			),
			lastSucceededChange,
		)
	}
	panic(fmt.Sprintf("unhandled check type %T", check))
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
