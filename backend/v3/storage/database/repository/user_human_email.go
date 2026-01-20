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

// CheckEmailOTP implements [domain.HumanUserRepository.CheckEmailOTP].
func (u userHuman) CheckEmailOTP(check domain.CheckType) database.Change {
	switch typ := check.(type) {
	case *domain.CheckTypeFailed:
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.verification.qualifiedTableName())
			builder.WriteString(" SET ")
			database.NewIncrementColumnChange(u.verification.failedAttemptsColumn())
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
				database.NewColumnCondition(u.verification.idColumn(), existingHumanUser.emailOTPVerificationIDColumn()),
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
				var (
					createdAt any = database.NowInstruction
					expiry    any = database.NullInstruction
				)
				if !typ.CreatedAt.IsZero() {
					createdAt = typ.CreatedAt
				}
				if typ.Expiry != nil {
					expiry = *typ.Expiry
				}
				builder.WriteString("INSERT INTO zitadel.verifications (instance_id, code, created_at, expiry) SELECT ")
				existingHumanUser.InstanceIDColumn().WriteQualified(builder)
				builder.WriteString(", ")
				builder.WriteArgs(typ.Code, createdAt, expiry)
				builder.WriteString(" FROM ")
				builder.WriteString(existingHumanUser.unqualifiedTableName())
				builder.WriteString(" RETURNING id")
			}, func(name string) database.Change {
				return database.NewChangeToStatement(u.emailOTPVerificationIDColumn(), func(builder *database.StatementBuilder) {
					builder.WriteString("SELECT id FROM ")
					builder.WriteString(name)
				})
			},
		)
	}
	panic(fmt.Sprintf("type not allowed for email OTP check change %T", check))
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
	switch typ := verification.(type) {
	case *domain.VerificationTypeSkipped:
		skippedAt := database.NewChange(u.emailVerifiedAtColumn(), database.NowInstruction)
		if !typ.SkippedAt.IsZero() {
			skippedAt = database.NewChange(u.emailVerifiedAtColumn(), typ.SkippedAt)
		}
		return database.NewChanges(
			database.NewChange(u.EmailColumn(), *typ.Value),
			skippedAt,
		)
	case *domain.VerificationTypeFailed:
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.verification.qualifiedTableName())
			builder.WriteString(" SET ")
			database.NewIncrementColumnChange(u.verification.failedAttemptsColumn())
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
				database.NewColumnCondition(u.verification.idColumn(), existingHumanUser.emailOTPVerificationIDColumn()),
			))
		}, nil)
	case *domain.VerificationTypeInit:
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.verifications(instance_id, user_id, value, code, created_at, expiry) SELECT instance_id, id, ")
			builder.WriteArgs(typ.Value, typ.Code, typ.CreatedAt, typ.Expiry)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			builder.WriteString(" RETURNING id")
		}, func(name string) database.Change {
			return database.NewChangeToStatement(u.emailVerificationIDColumn(), func(builder *database.StatementBuilder) {
				builder.WriteString("SELECT id FROM ")
				builder.WriteString(name)
			})
		})
	case *domain.VerificationTypeUpdate:
		changes := make(database.Changes, 0, 3)
		if typ.Value != nil {
			changes = append(changes, database.NewChange(u.EmailColumn(), *typ.Value))
		}
		if typ.Code != nil {
			changes = append(changes, database.NewChange(u.verification.codeColumn(), typ.Code))
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
				database.NewColumnCondition(u.verification.idColumn(), existingHumanUser.emailVerificationIDColumn()),
			))
		}, nil)
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

func (u userHuman) emailOTPVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_otp_verification_id")
}

func (u userHuman) lastSuccessfulEmailOTPCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_otp_last_successfully_checked_at")
}

func (u userHuman) emailVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "email_verified_at")
}

func (u userHuman) emailVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "unverified_email_id")
}
