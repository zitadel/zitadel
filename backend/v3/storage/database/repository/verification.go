package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type verification struct{}

func (v verification) unqualifiedTableName() string {
	return "verifications"
}

func (v verification) qualifiedTableName() string {
	return "zitadel." + v.unqualifiedTableName()
}

// verified constructs the following changes:
// - sets the verifiedAt column to the current time or the provided time
// - sets the pending verification ID column to null
// - deletes the verification and returns its value to be set
//
//	The caller is responsible for resetting the failed attempts counter on the object.
func (v verification) verified(verified *domain.VerificationTypeVerified, existingTableName string, instanceID, pendingVerificationID, verifiedAt, value database.Column) database.Changes {
	verifiedAtChange := database.NewChange(verifiedAt, database.NowInstruction)
	if !verified.VerifiedAt.IsZero() {
		verifiedAtChange = database.NewChange(verifiedAt, verified.VerifiedAt)
	}

	return database.NewChanges(
		verifiedAtChange,
		database.NewChangeToNull(pendingVerificationID),
		database.NewChangeToStatement(value, func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM ")
			builder.WriteString(v.qualifiedTableName())
			builder.WriteString(" USING ")
			builder.WriteString(existingTableName)
			writeCondition(builder, database.And(
				database.NewColumnCondition(v.instanceIDColumn(), instanceID),
				database.NewColumnCondition(v.idColumn(), pendingVerificationID),
			))
			builder.WriteString(" RETURNING value")
		}),
	)
}

func (v verification) skipped(skipped *domain.VerificationTypeSkipped, verifiedAt, value database.Column) database.Change {
	skippedAt := database.NewChange(verifiedAt, database.NowInstruction)
	if !skipped.SkippedAt.IsZero() {
		skippedAt = database.NewChange(verifiedAt, skipped.SkippedAt)
	}
	return database.NewChanges(
		database.NewChange(value, *skipped.Value),
		skippedAt,
	)
}

func (v verification) failed(existingTableName string, instanceID, verificationID database.Column) database.CTEChange {
	return database.NewCTEChange(func(builder *database.StatementBuilder) {
		builder.WriteString("UPDATE ")
		builder.WriteString(v.qualifiedTableName())
		builder.WriteString(" SET ")
		database.NewIncrementColumnChange(v.failedAttemptsColumn())
		builder.WriteString(" FROM ")
		builder.WriteString(existingTableName)
		writeCondition(builder, database.And(
			database.NewColumnCondition(v.instanceIDColumn(), instanceID),
			database.NewColumnCondition(v.idColumn(), verificationID),
		))
	}, nil)
}

func (v verification) init(init *domain.VerificationTypeInit, existingTableName string, verificationID database.Column) database.CTEChange {
	return database.NewCTEChange(func(builder *database.StatementBuilder) {
		var (
			createdAt any = database.NowInstruction
			expiry    any = database.NullInstruction
		)
		if !init.CreatedAt.IsZero() {
			createdAt = init.CreatedAt
		}
		if init.Expiry != nil {
			expiry = *init.Expiry
		}
		builder.WriteString("INSERT INTO zitadel.verifications(instance_id, user_id, value, code, created_at, expiry) SELECT instance_id, id, ")
		builder.WriteArgs(init.Value, init.Code, createdAt, expiry)
		builder.WriteString(" FROM ")
		builder.WriteString(existingTableName)
		builder.WriteString(" RETURNING id")
	}, func(name string) database.Change {
		return database.NewChangeToStatement(verificationID, func(builder *database.StatementBuilder) {
			builder.WriteString("SELECT id FROM ")
			builder.WriteString(name)
		})
	})
}

func (v verification) initCheck(init *domain.CheckTypeInit, existingTableName string, verificationID database.Column) database.CTEChange {
	return database.NewCTEChange(func(builder *database.StatementBuilder) {
		var (
			createdAt any = database.NowInstruction
			expiry    any = database.NullInstruction
		)
		if !init.CreatedAt.IsZero() {
			createdAt = init.CreatedAt
		}
		if init.Expiry != nil {
			expiry = *init.Expiry
		}
		builder.WriteString("INSERT INTO zitadel.verifications(instance_id, user_id, code, created_at, expiry) SELECT instance_id, id, ")
		builder.WriteArgs(init.Code, createdAt, expiry)
		builder.WriteString(" FROM ")
		builder.WriteString(existingTableName)
		builder.WriteString(" RETURNING id")
	}, func(name string) database.Change {
		return database.NewChangeToStatement(verificationID, func(builder *database.StatementBuilder) {
			builder.WriteString("SELECT id FROM ")
			builder.WriteString(name)
		})
	})
}

func (v verification) update(update *domain.VerificationTypeUpdate, existingTableName string, instanceID, verificationID database.Column) database.CTEChange {
	changes := make(database.Changes, 0, 3)
	if update.Value != nil {
		changes = append(changes, database.NewChange(v.valueColumn(), *update.Value))
	}
	if update.Code != nil {
		changes = append(changes, database.NewChange(v.codeColumn(), *update.Code))
	}
	if update.Expiry != nil {
		changes = append(changes, database.NewChange(v.expiryColumn(), *update.Expiry))
	}
	return database.NewCTEChange(func(builder *database.StatementBuilder) {
		builder.WriteString("UPDATE ")
		builder.WriteString(v.qualifiedTableName())
		builder.WriteString(" SET ")
		changes.Write(builder)
		builder.WriteString(" FROM ")
		builder.WriteString(existingTableName)
		writeCondition(builder, database.And(
			database.NewColumnCondition(v.instanceIDColumn(), instanceID),
			database.NewColumnCondition(v.idColumn(), verificationID),
		))
	}, nil)
}

func (v verification) succeeded(succeeded *domain.CheckTypeSucceeded, lastSuccessfulCheck, verificationID database.Column) database.Changes {
	lastSucceededChange := database.NewChange(lastSuccessfulCheck, database.NowInstruction)
	if !succeeded.SucceededAt.IsZero() {
		lastSucceededChange = database.NewChange(lastSuccessfulCheck, succeeded.SucceededAt)
	}
	return database.NewChanges(
		lastSucceededChange,
		database.NewChangeToNull(verificationID),
	)
}

// func (v verification) userVerification(builder *database.StatementBuilder, verificationType domain.VerificationType) database.Change {
// 	switch typ := verificationType.(type) {
// 	case *domain.VerificationTypeSkipped:
// 		return v.skipped(typ, userHuman)
// 	case *domain.VerificationTypeFailed:
// 		return database.NewCTEChange(func(builder *database.StatementBuilder) {
// 			builder.WriteString("UPDATE ")
// 			builder.WriteString(u.verification.qualifiedTableName())
// 			builder.WriteString(" SET ")
// 			database.NewIncrementColumnChange(u.verification.failedAttemptsColumn())
// 			builder.WriteString(" FROM ")
// 			builder.WriteString(existingHumanUser.unqualifiedTableName())
// 			writeCondition(builder, database.And(
// 				database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
// 				database.NewColumnCondition(u.verification.idColumn(), existingHumanUser.emailOTPVerificationIDColumn()),
// 			))
// 		}, nil)
// 	case *domain.VerificationTypeInit:
// 		return database.NewCTEChange(func(builder *database.StatementBuilder) {
// 			builder.WriteString("INSERT INTO zitadel.verifications(instance_id, user_id, value, code, created_at, expiry) SELECT instance_id, id, ")
// 			builder.WriteArgs(typ.Value, typ.Code, typ.CreatedAt, typ.Expiry)
// 			builder.WriteString(" FROM ")
// 			builder.WriteString(existingHumanUser.unqualifiedTableName())
// 			builder.WriteString(" RETURNING id")
// 		}, func(name string) database.Change {
// 			return database.NewChangeToStatement(u.emailVerificationIDColumn(), func(builder *database.StatementBuilder) {
// 				builder.WriteString("SELECT id FROM ")
// 				builder.WriteString(name)
// 			})
// 		})
// 	case *domain.VerificationTypeUpdate:
// 		changes := make(database.Changes, 0, 3)
// 		if typ.Value != nil {
// 			changes = append(changes, database.NewChange(u.EmailColumn(), *typ.Value))
// 		}
// 		if typ.Code != nil {
// 			changes = append(changes, database.NewChange(u.verification.codeColumn(), typ.Code))
// 		}
// 		if typ.Expiry != nil {
// 			changes = append(changes, database.NewChange(u.verification.expiryColumn(), *typ.Expiry))
// 		}
// 		return database.NewCTEChange(func(builder *database.StatementBuilder) {
// 			builder.WriteString("UPDATE ")
// 			builder.WriteString(u.verification.qualifiedTableName())
// 			builder.WriteString(" SET ")
// 			changes.Write(builder)
// 			builder.WriteString(" FROM ")
// 			builder.WriteString(existingHumanUser.unqualifiedTableName())
// 			writeCondition(builder, database.And(
// 				database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
// 				database.NewColumnCondition(u.verification.idColumn(), existingHumanUser.emailVerificationIDColumn()),
// 			))
// 		}, nil)
// 	}
// 	panic(fmt.Sprintf("type not allowed for email verification change %T", verification))
// }

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (v verification) instanceIDColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "instance_id")
}

func (v verification) idColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "id")
}

func (v verification) valueColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "value")
}

func (v verification) codeColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "code")
}

func (v verification) expiryColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "expires_at")
}

func (v verification) failedAttemptsColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "failed_attempts")
}

func (v verification) creationDateColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "created_at")
}

func (v verification) userIDColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "user_id")
}
