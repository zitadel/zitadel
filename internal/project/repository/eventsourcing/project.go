package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
)

func ProjectAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, project *model.Project) (*es_models.Aggregate, error) {
	if project == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-doe93", "Errors.Internal")
	}
	return aggCreator.NewAggregate(ctx, project.AggregateID, model.ProjectAggregate, model.ProjectVersion, project.Sequence)
}

func ApplicationAddedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-09du7", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationAdded, app)
		if app.OIDCConfig != nil {
			agg.AppendEvent(model.OIDCConfigAdded, app.OIDCConfig)
		}
		if app.APIConfig != nil {
			agg.AppendEvent(model.APIConfigAdded, app.APIConfig)
		}
		return agg, nil
	}
}

func APIConfigChangedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, config *model.APIConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if config == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slf32", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, a := range existingProject.Applications {
			if a.AppID == config.AppID {
				if a.APIConfig != nil {
					changes = a.APIConfig.Changes(config)
				}
			}
		}
		agg.AppendEvent(model.APIConfigChanged, changes)

		return agg, nil
	}
}

func OIDCConfigSecretChangedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, appID string, secret *crypto.CryptoValue) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 2)
		changes["appId"] = appID
		changes["clientSecret"] = secret

		agg.AppendEvent(model.OIDCConfigSecretChanged, changes)

		return agg, nil
	}
}

func APIConfigSecretChangedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, appID string, secret *crypto.CryptoValue) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 2)
		changes["appId"] = appID
		changes["clientSecret"] = secret

		agg.AppendEvent(model.APIConfigSecretChanged, changes)

		return agg, nil
	}
}

func OIDCClientSecretCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, appID string) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 1)
		changes["appId"] = appID

		agg.AppendEvent(model.OIDCClientSecretCheckSucceeded, changes)

		return agg, nil
	}
}

func OIDCClientSecretCheckFailedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, appID string) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 1)
		changes["appId"] = appID

		agg.AppendEvent(model.OIDCClientSecretCheckFailed, changes)

		return agg, nil
	}
}

func OIDCApplicationKeyAddedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, key *model.ClientKey) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ClientKeyAdded, key)

		return agg, nil
	}
}

func OIDCApplicationKeyRemovedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, keyID string) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 1)
		changes["keyId"] = keyID

		agg.AppendEvent(model.ClientKeyRemoved, changes)

		return agg, nil
	}
}
