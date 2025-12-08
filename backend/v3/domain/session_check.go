package domain

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// tarpitFn represents a tarpit function
//
// The input is the number of failed attempts after which the tarpit is started
type tarpitFn func(failedAttempts uint64)

// totpValidateFn represents a function to validate TOTP
type totpValidateFn func(toValidate, verifier string) bool

func getLockoutPolicy(ctx context.Context, db database.QueryExecutor, lockoutSettingsRepo LockoutSettingsRepository, instanceID, orgID string) (*LockoutSettings, error) {
	// We need the organization lockout policy first, and if not available, the instance (default) policy.
	// So we retrieve all records with a matching instance ID and organization ID OR
	// all records with a matching instance ID and NULL (or empty) organization ID.
	// Then we assume NULLs are sorted as largest numbers (that's the case in Postgres),
	// so we sort ascending by organization ID.
	// We limit the result to 1 so that we get either the org policy or the instance one.
	settings, err := lockoutSettingsRepo.List(ctx, db,
		listLockoutSettingCondition(lockoutSettingsRepo, instanceID, orgID),
		database.WithOrderByAscending(lockoutSettingsRepo.OrganizationIDColumn(), lockoutSettingsRepo.InstanceIDColumn()),
		database.WithLimit(1),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOM-3B8Z6s", "failed fetching lockout settings")
	}

	if rowsReturned := len(settings); rowsReturned != 1 {
		return nil, zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, int64(rowsReturned)), "DOM-mmsrCt", "unexpected number of rows returned")
	}

	return settings[0], nil
}

func listLockoutSettingCondition(repo LockoutSettingsRepository, instanceID, orgID string) database.QueryOption {
	instanceAndOrg := database.And(repo.InstanceIDCondition(instanceID), repo.OrganizationIDCondition(&orgID))
	orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
	onlyInstance := database.And(repo.InstanceIDCondition(instanceID), orgNullOrEmpty)

	return database.WithCondition(database.Or(instanceAndOrg, onlyInstance))
}

func shouldLockUser(lockoutSetting *LockoutSettings, userFailedAttempts uint64) bool {
	return lockoutSetting != nil &&
		lockoutSetting.MaxOTPAttempts != nil && *lockoutSetting.MaxOTPAttempts > 0 &&
		userFailedAttempts+1 >= *lockoutSetting.MaxOTPAttempts
}
