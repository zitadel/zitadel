//go:build integration

package query

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/feature"
)

// instanceFeaturesForIntegrationTests retrieves the instance feature state from the eventstore and populates target.
// This ensure strong consistency for the instance features in integration tests.
// It comes at the cost of 2 additional eventstore queries per call, so this won't be used in production code.
// In order to keep the production code clean, this function panics on error, which should then fail any running test.
func (q *Queries) instanceFeaturesForIntegrationTests(ctx context.Context, target *feature.Features) {
	f, err := q.GetInstanceFeatures(ctx, true)
	if err != nil {
		logging.OnError(ctx, err).Panic("could not get features from eventstore")
	}
	target.LoginDefaultOrg = f.LoginDefaultOrg.Value
	target.UserSchema = f.UserSchema.Value
	target.ImprovedPerformance = f.ImprovedPerformance.Value
	target.DebugOIDCParentError = f.DebugOIDCParentError.Value
	target.OIDCSingleV1SessionTermination = f.OIDCSingleV1SessionTermination.Value
	if f.LoginV2.Value != nil {
		target.LoginV2 = *f.LoginV2.Value
	}
	target.PermissionCheckV2 = f.PermissionCheckV2.Value
	target.ConsoleUseV2UserApi = f.ManagementConsoleUseV2UserApi.Value
	target.EnableRelationalTables = f.EnableRelationalTables.Value
}
