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

func (f *FeatureSource[T]) set(level feature.Level, value any) {
	f.Level = level
	f.Value = value.(T)
}

type SystemFeatures struct {
	Details *domain.ObjectDetails

	LoginDefaultOrg                 FeatureSource[bool]
	TriggerIntrospectionProjections FeatureSource[bool]
	LegacyIntrospection             FeatureSource[bool]
	UserSchema                      FeatureSource[bool]
	TokenExchange                   FeatureSource[bool]
	Actions                         FeatureSource[bool]
	ImprovedPerformance             FeatureSource[[]feature.ImprovedPerformanceType]
	OIDCSingleV1SessionTermination  FeatureSource[bool]
	DisableUserTokenEvent           FeatureSource[bool]
	EnableBackChannelLogout         FeatureSource[bool]
}

func (q *Queries) GetSystemFeatures(ctx context.Context) (_ *SystemFeatures, err error) {
	m := NewSystemFeaturesReadModel()
	if err := q.eventstore.FilterToQueryReducer(ctx, m); err != nil {
		return nil, err
	}
	m.system.Details = readModelToObjectDetails(m.ReadModel)
	return m.system, nil
}
