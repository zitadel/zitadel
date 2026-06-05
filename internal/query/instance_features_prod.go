//go:build !integration

package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/feature"
)

// instanceFeaturesForIntegrationTests is no-op in production builds.
func (*Queries) instanceFeaturesForIntegrationTests(context.Context, *feature.Features) {}
