package eventstore

import (
	"context"

	"github.com/caos/logging"

	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	features_model "github.com/caos/zitadel/internal/features/model"
	"github.com/caos/zitadel/internal/features/repository/view/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
)

type FeaturesRepo struct {
	Eventstore v1.Eventstore

	View *admin_view.View

	SearchLimit    uint64
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *FeaturesRepo) GetDefaultFeatures(ctx context.Context) (*features_model.FeaturesView, error) {
	features, viewErr := repo.View.FeaturesByAggregateID(domain.IAMID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		features = new(model.FeaturesView)
	}
	events, esErr := repo.getIAMEvents(ctx, features.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-Lsoj7", "Errors.Org.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-PSoc3").WithError(esErr).Debug("error retrieving new events")
		return model.FeaturesToModel(features), nil
	}
	featuresCopy := *features
	for _, event := range events {
		if err := featuresCopy.AppendEvent(event); err != nil {
			return model.FeaturesToModel(&featuresCopy), nil
		}
	}
	return model.FeaturesToModel(&featuresCopy), nil
}

func (repo *FeaturesRepo) GetOrgFeatures(ctx context.Context, orgID string) (*features_model.FeaturesView, error) {
	features, err := repo.View.FeaturesByAggregateID(orgID)
	if errors.IsNotFound(err) {
		return repo.GetDefaultFeatures(ctx)
	}
	if err != nil {
		return nil, err
	}
	return model.FeaturesToModel(features), nil
}

func (repo *FeaturesRepo) getIAMEvents(ctx context.Context, sequence uint64) ([]*models.Event, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}
