package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestUserHuman_RecoveryCodeChanges(t *testing.T) {
	userRepo := repository.UserRepository()

	checkedAt := time.Now()

	tests := []struct {
		name              string
		change            database.Change
		expectedStatement string
		expectedArgs      []any
	}{
		{
			name:              "add recovery codes",
			change:            userRepo.Human().AddRecoveryCodes([]string{"code3"}),
			expectedStatement: `recovery_codes = (ARRAY(SELECT DISTINCT unnest(array_cat(COALESCE(users.recovery_codes, $1), $2))))`,
			expectedArgs:      []any{[]string{}, []string{"code3"}},
		},
		{
			name:              "remove recovery code",
			change:            userRepo.Human().RemoveRecoveryCode("code3"),
			expectedStatement: `recovery_codes = (array_remove(users.recovery_codes, $1))`,
			expectedArgs:      []any{"code3"},
		},
		{
			name:              "remove all recovery codes",
			change:            userRepo.Human().RemoveAllRecoveryCodes(),
			expectedStatement: `recovery_codes = $1`,
			expectedArgs:      []any{[]string{}},
		},
		{
			name:              "set recovery code last successful checked at with checkedAt value zero",
			change:            userRepo.Human().SetLastSuccessfulRecoveryCodeCheck(time.Time{}),
			expectedStatement: `recovery_code_last_successful_check = NOW(), recovery_code_failed_attempts = $1`,
			expectedArgs:      []any{0},
		},
		{
			name:              "set recovery code last successful checked at with a valid checkedAt value",
			change:            userRepo.Human().SetLastSuccessfulRecoveryCodeCheck(checkedAt),
			expectedStatement: `recovery_code_last_successful_check = $1, recovery_code_failed_attempts = $2`,
			expectedArgs:      []any{checkedAt, 0},
		},
		{
			name:              "increment recovery code failed attempts",
			change:            userRepo.Human().IncrementRecoveryCodeFailedAttempts(),
			expectedStatement: `recovery_code_failed_attempts = COALESCE(users.recovery_code_failed_attempts, $1) + 1`,
			expectedArgs:      []any{0},
		},
		{
			name:              "reset recovery code failed attempts",
			change:            userRepo.Human().ResetRecoveryCodeFailedAttempts(),
			expectedStatement: `recovery_code_failed_attempts = $1`,
			expectedArgs:      []any{0},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var builder database.StatementBuilder
			err := test.change.Write(&builder)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedStatement, builder.String())
			assert.Equal(t, test.expectedArgs, builder.Args())
		})
	}
}
