package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

func ProjectByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "id should be filled")
	}
	return ProjectQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func ProjectQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.ProjectAggregate).
		LatestSequenceFilter(latestSequence)
}

func ProjectAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, project *Project) (*es_models.Aggregate, error) {
	if project == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-doe93", "existing project should not be nil")
	}
	return aggCreator.NewAggregate(ctx, project.ID, model.ProjectAggregate, projectVersion, project.Sequence)
}

func ProjectCreateAggregate(aggCreator *es_models.AggregateCreator, project *Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if project == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "project should not be nil")
		}

		agg, err := ProjectAggregate(ctx, aggCreator, project)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.ProjectAdded, project)
	}
}

func ProjectUpdateAggregate(aggCreator *es_models.AggregateCreator, existing *Project, new *Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if new == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new project should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.Changes(new)
		return agg.AppendEvent(model.ProjectChanged, changes)
	}
}

func ProjectDeactivateAggregate(aggCreator *es_models.AggregateCreator, project *Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return projectStateAggregate(aggCreator, project, model.ProjectDeactivated)
}

func ProjectReactivateAggregate(aggCreator *es_models.AggregateCreator, project *Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return projectStateAggregate(aggCreator, project, model.ProjectReactivated)
}

func projectStateAggregate(aggCreator *es_models.AggregateCreator, project *Project, state models.EventType) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, project)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(state, nil)
	}
}

func ProjectMemberAddedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, member *ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ie34f", "member should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectMemberAdded, member)
	}
}

func ProjectMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, member *ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d34fs", "member should not be nil")
		}

		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectMemberChanged, member)
	}
}

func ProjectMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, member *ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dieu7", "member should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectMemberRemoved, member)
	}
}

func ProjectRoleAddedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, role *ProjectRole) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if role == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sleo9", "role should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectRoleAdded, role)
	}
}

func ProjectRoleChangedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, role *ProjectRole) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if role == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-oe8sf", "member should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectRoleChanged, role)
	}
}

func ProjectRoleRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, role *ProjectRole) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if role == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8eis", "member should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectRoleRemoved, role)
	}
}

func ApplicationAddedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, app *Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-09du7", "app should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationAdded, app)
		if app.OIDCConfig != nil {
			agg.AppendEvent(model.OIDCConfigAdded, app.OIDCConfig)
		}
		return agg, nil
	}
}

func ApplicationChangedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, app *Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sleo9", "app should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, a := range existing.Applications {
			if a.AppID == app.AppID {
				changes = a.Changes(app)
			}
		}
		agg.AppendEvent(model.ApplicationChanged, changes)

		return agg, nil
	}
}

func ApplicationRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, app *Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-se23g", "app should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationRemoved, &ApplicationID{AppID: app.AppID})

		return agg, nil
	}
}

func ApplicationDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, app *Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slfi3", "app should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationDeactivated, &ApplicationID{AppID: app.AppID})

		return agg, nil
	}
}

func ApplicationReactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, app *Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slf32", "app should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationReactivated, &ApplicationID{AppID: app.AppID})

		return agg, nil
	}
}

func OIDCConfigChangedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, config *OIDCConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if config == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slf32", "config should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, a := range existing.Applications {
			if a.AppID == config.AppID {
				if a.OIDCConfig != nil {
					changes = a.OIDCConfig.Changes(config)
				}
			}
		}
		agg.AppendEvent(model.OIDCConfigChanged, changes)

		return agg, nil
	}
}

func OIDCConfigSecretChangedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, appID string, crypto *crypto.CryptoValue) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 1)
		changes["appId"] = appID
		changes["clientSecret"] = crypto

		agg.AppendEvent(model.OIDCConfigSecretChanged, changes)

		return agg, nil
	}
}
