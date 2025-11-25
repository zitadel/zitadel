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
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		var createdAt any = database.NowInstruction
		if !typ.CreatedAt.IsZero() {
			createdAt = typ.CreatedAt
		}
		return database.NewChanges(
			database.NewCTEChange(func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.verifications (instance_id, user_id, value, code, created_at, expiry) SELECT")
				builder.WriteArgs(
					existingHumanUser.instanceIDColumn(),
					existingHumanUser.idColumn(),
					typ.Value,
					typ.Code,
					createdAt,
					typ.Expiry,
				)
				builder.WriteString(" FROM ")
				builder.WriteString(existingHumanUser.unqualifiedTableName())
				builder.WriteString(" RETURNING verification.*")
			},
				func(name string) database.Change {
					return database.NewChangeToStatement(
						u.unverifiedPasswordIDColumn(),
						func(builder *database.StatementBuilder) {
							builder.WriteString(" SELECT ")
							existingHumanUser.verification.IDColumn().WriteQualified(builder)
							builder.WriteString(" FROM ")
							builder.WriteString(name)
							writeCondition(builder, database.And(
								database.NewColumnCondition(u.instanceIDColumn(), database.NewColumn(name, "instance_id")),
								database.NewColumnCondition(u.idColumn(), database.NewColumn(name, "user_id")),
							))
						},
					)
				},
			),
		)
	}
	panic(fmt.Sprintf("unhandled verification type %T", verification))
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
