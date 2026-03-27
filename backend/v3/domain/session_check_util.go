package domain

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type verifierFn func(encoded, password string) (updated string, err error)

func GetLockoutPolicy(ctx context.Context, opts *InvokeOpts, instanceID, orgID string) (*LockoutSettings, error) {
	lockoutSettingRepo := opts.lockoutSettingRepo

	// We need the organization lockout policy first, and if not available, the instance (default) policy.
	// So we retrieve all records with a matching instance ID and organization ID OR
	// all records with a matching instance ID and NULL (or empty) organization ID.
	// Then we assume NULLs are sorted as largest numbers (that's the case in Postgres),
	// so we sort ascending by organization ID.
	// We limit the result to 1 so that we get either the org policy or the instance one.
	settings, err := lockoutSettingRepo.List(ctx, opts.DB(),
		listLockoutSettingCondition(lockoutSettingRepo, instanceID, orgID),
		database.WithOrderByAscending(lockoutSettingRepo.OrganizationIDColumn(), lockoutSettingRepo.InstanceIDColumn()),
		database.WithLimit(1),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOM-3B8Z6s", "failed fetching lockout settings")
	}

	if rowsReturned := len(settings); rowsReturned != 1 {
		return nil, zerrors.ThrowInternal(NewRowsReturnedMismatchError(1, int64(rowsReturned)), "DOM-mmsrCt", "unexpected number of rows returned")
	}

	return settings[0], nil
}

func listLockoutSettingCondition(repo LockoutSettingsRepository, instanceID, orgID string) database.QueryOption {
	instanceAndOrg := database.And(repo.InstanceIDCondition(instanceID), repo.OrganizationIDCondition(&orgID))
	orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
	onlyInstance := database.And(repo.InstanceIDCondition(instanceID), orgNullOrEmpty)

	return database.WithCondition(database.Or(instanceAndOrg, onlyInstance))
}
