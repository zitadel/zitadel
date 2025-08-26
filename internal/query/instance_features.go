package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/feature"
)

type InstanceFeatures struct {
	Details                        *domain.ObjectDetails
	LoginDefaultOrg                FeatureSource[bool]
	UserSchema                     FeatureSource[bool]
	TokenExchange                  FeatureSource[bool]
	ImprovedPerformance            FeatureSource[[]feature.ImprovedPerformanceType]
	DebugOIDCParentError           FeatureSource[bool]
	OIDCSingleV1SessionTermination FeatureSource[bool]
	EnableBackChannelLogout        FeatureSource[bool]
	LoginV2                        FeatureSource[*feature.LoginV2]
	PermissionCheckV2              FeatureSource[bool]
	ConsoleUseV2UserApi            FeatureSource[bool]
}

func (q *Queries) GetInstanceFeatures(ctx context.Context, cascade bool) (_ *InstanceFeatures, err error) {
	var system *SystemFeatures
	if cascade {
		system, err = q.GetSystemFeatures(ctx)
		if err != nil {
			return nil, err
		}
	}
	m := NewInstanceFeaturesReadModel(ctx, system)
	if err = q.eventstore.FilterToQueryReducer(ctx, m); err != nil {
		return nil, err
	}
	m.instance.Details = readModelToObjectDetails(m.ReadModel)
	return m.instance, nil
}
