package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type verification struct{}

func (v verification) unqualifiedTableName() string {
	return "verifications"
}

func (v verification) qualifiedTableName() string {
	return "zitadel." + v.unqualifiedTableName()
}

// init creates a new verification
func (v verification) init(init *domain.VerificationTypeInit, existingTableName string, verificationID database.Column) database.CTEChange {
	return database.NewCTEChange(func(builder *database.StatementBuilder) {
		var (
			createdAt any = database.NowInstruction
			expiry    any = database.NullInstruction
			id        any = database.GenRandomUUIDInstruction
		)
		if !init.CreatedAt.IsZero() {
			createdAt = init.CreatedAt
		}
		if init.Expiry != nil {
			expiry = *init.Expiry
		}
		if init.ID != nil {
			id = *init.ID
		}
		builder.WriteString("INSERT INTO zitadel.verifications(instance_id, user_id, id, code, created_at, updated_at, expiry) SELECT instance_id, id, ")
		builder.WriteArgs(id, init.Code, createdAt, createdAt, expiry)
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

// updates the given fields of the verification with the given ID. Only non-nil fields get updated.
func (v verification) update(update *domain.VerificationTypeUpdate, existingTableName string, instanceID, verificationID database.Column) database.CTEChange {
	changes := make(database.Changes, 0, 3)
	if update.Code != nil {
		changes = append(changes, database.NewChange(v.codeColumn(), update.Code))
	}
	if update.Expiry != nil {
		changes = append(changes, database.NewChange(v.expiryColumn(), *update.Expiry))
	}
	if update.UpdatedAt.IsZero() {
		changes = append(changes, database.NewChange(v.updatedAtColumn(), database.NowInstruction))
	} else {
		changes = append(changes, database.NewChange(v.updatedAtColumn(), update.UpdatedAt))
	}
	return database.NewCTEChange(func(builder *database.StatementBuilder) {
		builder.WriteString("UPDATE ")
		builder.WriteString(v.qualifiedTableName())
		builder.WriteString(" SET ")
		err := changes.Write(builder)
		logging.New(logging.StreamRuntime).Debug("write changes in cte failed", "error", err)
		builder.WriteString(" FROM ")
		builder.WriteString(existingTableName)
		writeCondition(builder, database.And(
			database.NewColumnCondition(v.instanceIDColumn(), instanceID),
			database.NewColumnCondition(v.idColumn(), verificationID),
		))
	}, nil)
}

// verified constructs the following changes:
// - sets the verifiedAt column to the current time or the provided time
// - clears the pending verification ID
// - deletes the verification
// - resets the failed attempts counter
func (v verification) verified(verified *domain.VerificationTypeSucceeded, existingTableName string, instanceID, pendingVerificationID, verifiedAt, failedAttempts database.Column) database.Changes {
	verifiedAtChange := database.NewChange(verifiedAt, database.NowInstruction)
	if !verified.VerifiedAt.IsZero() {
		verifiedAtChange = database.NewChange(verifiedAt, verified.VerifiedAt)
	}

	return database.NewChanges(
		verifiedAtChange,
		database.NewChangeToNull(pendingVerificationID),
		database.NewChange(failedAttempts, 0),
		database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("DELETE FROM ")
				builder.WriteString(v.qualifiedTableName())
				builder.WriteString(" USING ")
				builder.WriteString(existingTableName)
				writeCondition(builder, database.And(
					database.NewColumnCondition(v.instanceIDColumn(), instanceID),
					database.NewColumnCondition(v.idColumn(), pendingVerificationID),
				))
			}, nil,
		),
	)
}

// skipped constructs the following changes:
// - sets the verifiedAt column to the current time or the provided time
// - clears the pending verification ID
// - resets the failed attempts counter
func (v verification) skipped(skipped *domain.VerificationTypeSkipped, verifiedAt, pendingVerificationID, failedAttempts database.Column) database.Changes {
	verifiedAtChange := database.NewChange(verifiedAt, database.NowInstruction)
	if !skipped.SkippedAt.IsZero() {
		verifiedAtChange = database.NewChange(verifiedAt, skipped.SkippedAt)
	}
	return database.NewChanges(
		verifiedAtChange,
		database.NewChange(failedAttempts, 0),
		database.NewChangeToNull(pendingVerificationID),
	)
}

// failed increments the failed attempts counter.
func (v verification) failed(existingTableName string, instanceID, verificationID database.Column) database.CTEChange {
	return database.NewCTEChange(func(builder *database.StatementBuilder) {
		builder.WriteString("UPDATE ")
		builder.WriteString(v.qualifiedTableName())
		builder.WriteString(" SET ")
		err := database.NewIncrementColumnChange(v.failedAttemptsColumn(), database.Coalesce(v.failedAttemptsColumn(), 0)).Write(builder)
		logging.New(logging.StreamRuntime).Debug("write cte failed", "error", err)
		builder.WriteString(" FROM ")
		builder.WriteString(existingTableName)
		writeCondition(builder, database.And(
			database.NewColumnCondition(v.instanceIDColumn(), instanceID),
			database.NewColumnCondition(v.idColumn(), verificationID),
		))
	}, nil)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (v verification) instanceIDColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "instance_id")
}

func (v verification) idColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "id")
}

func (v verification) codeColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "code")
}

func (v verification) expiryColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "expiry")
}

func (v verification) failedAttemptsColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "failed_attempts")
}

func (v verification) userIDColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "user_id")
}

func (v verification) updatedAtColumn() database.Column {
	return database.NewColumn(v.unqualifiedTableName(), "updated_at")
}
