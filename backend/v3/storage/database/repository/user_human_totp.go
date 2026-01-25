package repository

import (
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

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
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.verifications (instance_id, user_id, value, code, created_at, expiry) SELECT existing_user.instance_id, existing_user.id, ")
			builder.WriteArgs(
				typ.Value,
				typ.Code,
				createdAt,
				typ.Expiry,
			)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			builder.WriteString(" RETURNING verifications.*")
		},
			func(name string) database.Change {
				return database.NewChangeToStatement(
					u.passwordVerificationIDColumn(),
					func(builder *database.StatementBuilder) {
						builder.WriteString(" SELECT ")
						existingHumanUser.verification.idColumn().WriteQualified(builder)
						builder.WriteString(" FROM ")
						builder.WriteString(name)
						writeCondition(builder, database.And(
							database.NewColumnCondition(u.InstanceIDColumn(), database.NewColumn(name, "instance_id")),
							database.NewColumnCondition(u.IDColumn(), database.NewColumn(name, "user_id")),
						))
					},
				)
			},
		)
	case *domain.VerificationTypeVerified:
		verifiedAtChange := database.NewChange(u.totpVerifiedAtColumn(), database.NowInstruction)
		if !typ.VerifiedAt.IsZero() {
			verifiedAtChange = database.NewChange(u.totpVerifiedAtColumn(), typ.VerifiedAt)
		}
		return verifiedAtChange
	case *domain.VerificationTypeUpdate:
		changes := make(database.Changes, 0, 3)
		if typ.Code != nil {
			changes = append(changes, database.NewChange(u.verification.codeColumn(), typ.Code))
		}
		if typ.Value != nil {
			changes = append(changes, database.NewChange(u.verification.valueColumn(), *typ.Value))
		}
		if typ.Expiry != nil {
			changes = append(changes, database.NewChange(u.verification.expiryColumn(), *typ.Expiry))
		}
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.verification.qualifiedTableName())
			builder.WriteString(" SET ")
			changes.Write(builder)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
				// database.NewColumnCondition(u.verification.idColumn(), existingHumanUser.unverifiedTOTPIDColumn()),
			))
		}, nil)
	case *domain.VerificationTypeSkipped:
		var skippedAt any = database.NowInstruction
		skippedAtChange := database.NewChange(u.totpVerifiedAtColumn(), database.NowInstruction)
		if !typ.SkippedAt.IsZero() {
			skippedAt = typ.SkippedAt
			skippedAtChange = database.NewChange(u.totpVerifiedAtColumn(), typ.SkippedAt)
		}
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO ")
				builder.WriteString(u.verification.qualifiedTableName())
				builder.WriteString(" (")
				database.Columns{
					u.verification.instanceIDColumn(),
					u.verification.valueColumn(),
					u.verification.creationDateColumn(),
				}.WriteUnqualified(builder)
				builder.WriteString(") SELECT ")
				existingHumanUser.InstanceIDColumn().WriteQualified(builder)
				builder.WriteString(", ")
				builder.WriteArgs(typ.Value, skippedAt)
				builder.WriteString(" FROM ")
				builder.WriteString(existingHumanUser.unqualifiedTableName())
				builder.WriteString(" RETURNING verifications.*")
			},
			func(name string) database.Change {
				return database.NewChanges(
					database.NewChangeToStatement(u.totpSecretIDColumn(), func(builder *database.StatementBuilder) {
						builder.WriteString(" SELECT id FROM ")
						builder.WriteString(name)
					}),
					skippedAtChange,
				)
			},
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
