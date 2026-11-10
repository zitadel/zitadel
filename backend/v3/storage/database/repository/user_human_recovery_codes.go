package repository

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (u userHuman) AddRecoveryCodes(codes []string) database.Change {
	return database.NewChangeToStatement(u.recoveryCodesColumn(), func(builder *database.StatementBuilder) {
		// the following lines build the value expression
		// `ARRAY(SELECT DISTINCT unnest(array_cat(recovery_codes, $1)))`
		// for the SET recovery_codes part of the user UPDATE query
		builder.WriteString("ARRAY(SELECT DISTINCT unnest(array_cat(")
		database.Coalesce(u.recoveryCodesColumn(), []string{}).WriteQualified(builder)
		builder.WriteString(", ")
		builder.WriteArg(codes)
		builder.WriteString(")))")
	})
}

func (u userHuman) RemoveRecoveryCode(code string) database.Change {
	return database.NewChangeToStatement(u.recoveryCodesColumn(), func(builder *database.StatementBuilder) {
		// the following lines build the value expression
		// `array_remove(recovery_codes, $1)`
		// for the SET recovery_codes part of the user UPDATE query
		builder.WriteString("array_remove(")
		u.recoveryCodesColumn().WriteQualified(builder)
		builder.WriteString(", ")
		builder.WriteArg(code)
		builder.WriteString(")")
	})
}

func (u userHuman) RemoveAllRecoveryCodes() database.Change {
	return database.NewChange(u.recoveryCodesColumn(), []string{})
}

func (u userHuman) SetLastSuccessfulRecoveryCodeCheck(checkedAt time.Time) database.Change {
	if checkedAt.IsZero() {
		return database.NewChanges(
			database.NewChange(u.lastSuccessfulRecoveryCodeCheckColumn(), database.NowInstruction),
			u.ResetRecoveryCodeFailedAttempts(),
		)
	}
	return database.NewChanges(
		database.NewChange(u.lastSuccessfulRecoveryCodeCheckColumn(), checkedAt),
		u.ResetRecoveryCodeFailedAttempts(),
	)
}

func (u userHuman) IncrementRecoveryCodeFailedAttempts() database.Change {
	return database.NewIncrementColumnChange(u.recoveryCodeFailedAttemptsColumn(), database.Coalesce(u.recoveryCodeFailedAttemptsColumn(), 0))
}

func (u userHuman) ResetRecoveryCodeFailedAttempts() database.Change {
	return database.NewChange(u.recoveryCodeFailedAttemptsColumn(), 0)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) recoveryCodesColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "recovery_codes")
}

func (u userHuman) lastSuccessfulRecoveryCodeCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "recovery_code_last_successful_check")
}

func (u userHuman) recoveryCodeFailedAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "recovery_code_failed_attempts")
}
