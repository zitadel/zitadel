package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
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

func ProjectAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, project *model.Project) (*es_models.Aggregate, error) {
	if project == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-doe93", "existing project should not be nil")
	}
	return aggCreator.NewAggregate(ctx, project.AggregateID, model.ProjectAggregate, model.ProjectVersion, project.Sequence)
}

func ProjectCreateAggregate(aggCreator *es_models.AggregateCreator, project *model.Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if project == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "project should not be nil")
		}

		agg, err := ProjectAggregate(ctx, aggCreator, project)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.ProjectAggregate).
			EventTypesFilter(model.ProjectAdded, model.ProjectChanged, model.ProjectRemoved)

		validation := addProjectValidation(project.Name)
		return agg.SetPrecondition(validationQuery, validation).AppendEvent(model.ProjectAdded, project)
	}
}

func ProjectUpdateAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, new *model.Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if new == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new project should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.Changes(new)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9soPE", "no changes found")
		}
		if existing.Name != new.Name {
			validationQuery := es_models.NewSearchQuery().
				AggregateTypeFilter(model.ProjectAggregate).
				EventTypesFilter(model.ProjectAdded, model.ProjectChanged, model.ProjectRemoved)

			validation := addProjectValidation(new.Name)
			agg.SetPrecondition(validationQuery, validation)
		}
		return agg.AppendEvent(model.ProjectChanged, changes)
	}
}

func ProjectDeactivateAggregate(aggCreator *es_models.AggregateCreator, project *model.Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return projectStateAggregate(aggCreator, project, model.ProjectDeactivated)
}

func ProjectReactivateAggregate(aggCreator *es_models.AggregateCreator, project *model.Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return projectStateAggregate(aggCreator, project, model.ProjectReactivated)
}

func projectStateAggregate(aggCreator *es_models.AggregateCreator, project *model.Project, state models.EventType) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, project)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(state, nil)
	}
}

func ProjectMemberAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, member *model.ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ie34f", "member should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(usr_model.UserAggregate).
			AggregateIDFilter(member.UserID)

		validation := addProjectMemberValidation()
		return agg.SetPrecondition(validationQuery, validation).AppendEvent(model.ProjectMemberAdded, member)
	}
}

func ProjectMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, member *model.ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
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

func ProjectMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, member *model.ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
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

func ProjectRoleAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, role *model.ProjectRole) func(ctx context.Context) (*es_models.Aggregate, error) {
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

func ProjectRoleChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, role *model.ProjectRole) func(ctx context.Context) (*es_models.Aggregate, error) {
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

func ProjectRoleRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, role *model.ProjectRole) func(ctx context.Context) (*es_models.Aggregate, error) {
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

func ApplicationAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
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

func ApplicationChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
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

func ApplicationRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-se23g", "app should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationRemoved, &model.ApplicationID{AppID: app.AppID})

		return agg, nil
	}
}

func ApplicationDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slfi3", "app should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationDeactivated, &model.ApplicationID{AppID: app.AppID})

		return agg, nil
	}
}

func ApplicationReactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slf32", "app should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationReactivated, &model.ApplicationID{AppID: app.AppID})

		return agg, nil
	}
}

func OIDCConfigChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, config *model.OIDCConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
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

func OIDCConfigSecretChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, appID string, secret *crypto.CryptoValue) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
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

func OIDCClientSecretCheckSucceededAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, appID string) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 1)
		changes["appId"] = appID

		agg.AppendEvent(model.OIDCClientSecretCheckSucceeded, changes)

		return agg, nil
	}
}

func OIDCClientSecretCheckFailedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, appID string) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 1)
		changes["appId"] = appID

		agg.AppendEvent(model.OIDCClientSecretCheckFailed, changes)

		return agg, nil
	}
}

func ProjectGrantAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kd89w", "grant should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(org_model.OrgAggregate).
			AggregateIDFilter(grant.GrantedOrgID)

		validation := addProjectGrantValidation()
		agg.SetPrecondition(validationQuery, validation).AppendEvent(model.ProjectGrantAdded, grant)
		return agg, nil
	}
}

func ProjectGrantChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d9ie2", "grant should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, g := range existing.Grants {
			if g.GrantID == grant.GrantID {
				changes = g.Changes(grant)
			}
		}
		agg.AppendEvent(model.ProjectGrantChanged, changes)

		return agg, nil
	}
}

func ProjectGrantRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kci8d", "grant should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ProjectGrantRemoved, &model.ProjectGrantID{GrantID: grant.GrantID})

		return agg, nil
	}
}

func ProjectGrantDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-id832", "grant should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ProjectGrantDeactivated, &model.ProjectGrantID{GrantID: grant.GrantID})

		return agg, nil
	}
}

func ProjectGrantReactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8diw2", "grant should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectGrantReactivated, &model.ProjectGrantID{GrantID: grant.GrantID})
	}
}

func ProjectGrantMemberAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, member *model.ProjectGrantMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4ufh6", "grant should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(usr_model.UserAggregate).
			AggregateIDFilter(member.UserID)

		validation := addProjectGrantMemberValidation()
		return agg.SetPrecondition(validationQuery, validation).AppendEvent(model.ProjectGrantMemberAdded, member)
	}
}

func ProjectGrantMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, member *model.ProjectGrantMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8i4h", "member should not be nil")
		}

		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 1)
		changes["grantId"] = member.GrantID
		changes["userId"] = member.UserID
		changes["roles"] = member.Roles

		return agg.AppendEvent(model.ProjectGrantMemberChanged, changes)
	}
}

func ProjectGrantMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Project, member *model.ProjectGrantMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slp0r", "member should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectGrantMemberRemoved, member)
	}
}

func addProjectValidation(projectName string) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		projects := make([]*model.Project, 0)
		for _, event := range events {
			switch event.Type {
			case model.ProjectAdded:
				project := &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: event.AggregateID}}
				project.AppendAddProjectEvent(event)
				projects = append(projects, project)
			case model.ProjectChanged:
				_, project := model.GetProject(projects, event.AggregateID)
				project.AppendAddProjectEvent(event)
			case model.ProjectRoleRemoved:
				for i, project := range projects {
					if project.AggregateID == event.AggregateID {
						projects[i] = projects[len(projects)-1]
						projects[len(projects)-1] = nil
						projects = projects[:len(projects)-1]
					}
				}
			}
		}
		for _, p := range projects {
			if p.Name == projectName {
				return errors.ThrowPreconditionFailed(nil, "EVENT-s9oPw", "conditions not met")
			}
		}
		return nil
	}
}

func addProjectMemberValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		return checkExistsUser(events...)
	}
}

func addProjectGrantValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existsOrg := false
		for _, event := range events {
			switch event.AggregateType {
			case org_model.OrgAggregate:
				switch event.Type {
				case org_model.OrgAdded:
					existsOrg = true
				case org_model.OrgRemoved:
					existsOrg = false
				}
			}
		}
		if existsOrg {
			return nil
		}
		return errors.ThrowPreconditionFailed(nil, "EVENT-3OfIm", "conditions not met")
	}
}

func addProjectGrantMemberValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		return checkExistsUser(events...)
	}
}

func checkExistsUser(events ...*es_models.Event) error {
	existsUser := false
	for _, event := range events {
		switch event.AggregateType {
		case usr_model.UserAggregate:
			switch event.Type {
			case usr_model.UserAdded, usr_model.UserRegistered:
				existsUser = true
			case usr_model.UserRemoved:
				existsUser = false
			}
		}
	}
	if existsUser {
		return nil
	}
	return errors.ThrowPreconditionFailed(nil, "EVENT-3OfIm", "conditions not met")
}
