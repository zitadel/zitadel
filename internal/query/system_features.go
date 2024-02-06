package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/feature"
)

type FeatureSource[T any] struct {
	Level feature.Level
	Value T
}

type SystemFeatures struct {
	Details *domain.ObjectDetails

	LoginDefaultOrg                 FeatureSource[bool]
	TriggerIntrospectionProjections FeatureSource[bool]
	LegacyIntrospection             FeatureSource[bool]
}

func (q *Queries) GetSystemFeatures(ctx context.Context, cascade bool) (_ *SystemFeatures, err error) {
	var defaults *feature.Features
	if cascade {
		defaults = new(feature.Features)
		*defaults = q.defaultFeatures // make sure we copy the defaults.
	}
	m := NewSystemFeaturesReadModel(defaults)
	if err := q.eventstore.FilterToQueryReducer(ctx, m); err != nil {
		return nil, err
	}
	m.system.Details = readModelToObjectDetails(m.ReadModel)
	return m.system, nil
}
