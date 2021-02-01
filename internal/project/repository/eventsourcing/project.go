package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func ProjectByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "Errors.Project.ProjectIDMissing")
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
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-doe93", "Errors.Internal")
	}
	return aggCreator.NewAggregate(ctx, project.AggregateID, model.ProjectAggregate, model.ProjectVersion, project.Sequence)
}

func ProjectAggregateOverwriteContext(ctx context.Context, aggCreator *es_models.AggregateCreator, project *model.Project, resourceOwnerID string, userID string) (*es_models.Aggregate, error) {
	if project == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ADv2r", "Errors.Internal")
	}

	return aggCreator.NewAggregate(ctx, project.AggregateID, model.ProjectAggregate, model.ProjectVersion, project.Sequence, es_models.OverwriteResourceOwner(resourceOwnerID), es_models.OverwriteEditorUser(userID))
}

func ProjectCreateAggregate(aggCreator *es_models.AggregateCreator, project *model.Project, member *model.ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if project == nil || member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "Errors.Internal")
		}

		agg, err := ProjectAggregate(ctx, aggCreator, project)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.ProjectAggregate).
			ResourceOwnerFilter(authz.GetCtxData(ctx).OrgID).
			EventTypesFilter(model.ProjectAdded, model.ProjectChanged, model.ProjectRemoved)

		validation := addProjectValidation(project.Name)
		agg, err = agg.SetPrecondition(validationQuery, validation).AppendEvent(model.ProjectAdded, project)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectMemberAdded, member)
	}
}

func ProjectUpdateAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, newProject *model.Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if newProject == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		changes := existingProject.Changes(newProject)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9soPE", "Errors.NoChangesFound")
		}
		if existingProject.Name != newProject.Name {
			validationQuery := es_models.NewSearchQuery().
				AggregateTypeFilter(model.ProjectAggregate).
				EventTypesFilter(model.ProjectAdded, model.ProjectChanged, model.ProjectRemoved)

			validation := addProjectValidation(newProject.Name)
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

func ProjectRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existingProject *model.Project) (*es_models.Aggregate, error) {
	if existingProject == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Cj7lb", "Errors.Internal")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.ProjectRemoved, existingProject)
}

func ProjectMemberAddedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, member *model.ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ie34f", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
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

func ProjectMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, member *model.ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d34fs", "Errors.Internal")
		}

		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectMemberChanged, member)
	}
}

func ProjectMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, member *model.ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dieu7", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectMemberRemoved, member)
	}
}

func ProjectRoleAddedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, roles ...*model.ProjectRole) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if roles == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sleo9", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			agg, err = agg.AppendEvent(model.ProjectRoleAdded, role)
			if err != nil {
				return nil, err
			}
		}
		return agg, nil
	}
}

func ProjectRoleChangedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, role *model.ProjectRole) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if role == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-oe8sf", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectRoleChanged, role)
	}
}

func ProjectRoleRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existingProject *model.Project, role *model.ProjectRole, grants []*model.ProjectGrant) (*es_models.Aggregate, error) {
	if role == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8eis", "Errors.Internal")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.ProjectRoleRemoved, role)
	if err != nil {
		return nil, err
	}
	for _, grant := range grants {
		var changes map[string]interface{}
		if _, g := model.GetProjectGrant(existingProject.Grants, grant.GrantID); grant != nil {
			changes = g.Changes(grant)
			agg, err = agg.AppendEvent(model.ProjectGrantCascadeChanged, changes)
			if err != nil {
				return nil, err
			}
		}
	}
	return agg, nil
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
		return agg, nil
	}
}

func ApplicationChangedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sleo9", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, a := range existingProject.Applications {
			if a.AppID == app.AppID {
				changes = a.Changes(app)
			}
		}
		agg.AppendEvent(model.ApplicationChanged, changes)

		return agg, nil
	}
}

func ApplicationRemovedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-se23g", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationRemoved, &model.ApplicationID{AppID: app.AppID})

		return agg, nil
	}
}

func ApplicationDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slfi3", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationDeactivated, &model.ApplicationID{AppID: app.AppID})

		return agg, nil
	}
}

func ApplicationReactivatedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, app *model.Application) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if app == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slf32", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ApplicationReactivated, &model.ApplicationID{AppID: app.AppID})

		return agg, nil
	}
}

func OIDCConfigChangedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, config *model.OIDCConfig) func(ctx context.Context) (*es_models.Aggregate, error) {
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
				if a.OIDCConfig != nil {
					changes = a.OIDCConfig.Changes(config)
				}
			}
		}
		agg.AppendEvent(model.OIDCConfigChanged, changes)

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

func OIDCApplicationTokenAddedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, token *model.Token) es_sdk.AggregateFunc {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := ProjectAggregateOverwriteContext(ctx, aggCreator, existingProject, existingProject.ResourceOwner, existingProject.AggregateID) //TODO: !!!!
		if err != nil {
			return nil, err
		}

		agg.AppendEvent(model.TokenAdded, token)

		return agg, nil
	}
}

func ProjectGrantAddedAggregate(aggCreator *es_models.AggregateCreator, project *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kd89w", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, project)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(org_es_model.OrgAggregate).
			AggregateIDFilter(grant.GrantedOrgID)

		validation := addProjectGrantValidation()
		agg.SetPrecondition(validationQuery, validation).AppendEvent(model.ProjectGrantAdded, grant)
		return agg, nil
	}
}

func ProjectGrantChangedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d9ie2", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		var changes map[string]interface{}
		for _, g := range existingProject.Grants {
			if g.GrantID == grant.GrantID {
				changes = g.Changes(grant)
			}
		}
		agg.AppendEvent(model.ProjectGrantChanged, changes)

		return agg, nil
	}
}

func ProjectGrantRemovedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kci8d", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ProjectGrantRemoved, &model.ProjectGrantID{GrantID: grant.GrantID})

		return agg, nil
	}
}

func ProjectGrantDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-id832", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(model.ProjectGrantDeactivated, &model.ProjectGrantID{GrantID: grant.GrantID})

		return agg, nil
	}
}

func ProjectGrantReactivatedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, grant *model.ProjectGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8diw2", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectGrantReactivated, &model.ProjectGrantID{GrantID: grant.GrantID})
	}
}

func ProjectGrantMemberAddedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, member *model.ProjectGrantMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4ufh6", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
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

func ProjectGrantMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, member *model.ProjectGrantMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8i4h", "Errors.Internal")
		}

		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
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

func ProjectGrantMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, existingProject *model.Project, member *model.ProjectGrantMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slp0r", "Errors.Internal")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existingProject)
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
			case model.ProjectRemoved:
				for i := len(projects) - 1; i >= 0; i-- {
					if projects[i].AggregateID == event.AggregateID {
						projects[i] = projects[len(projects)-1]
						projects[len(projects)-1] = nil
						projects = projects[:len(projects)-1]
					}
				}
			}
		}
		for _, p := range projects {
			if p.Name == projectName {
				return errors.ThrowPreconditionFailed(nil, "EVENT-s9oPw", "Errors.Project.AlreadyExists")
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
			case org_es_model.OrgAggregate:
				switch event.Type {
				case org_es_model.OrgAdded:
					existsOrg = true
				case org_es_model.OrgRemoved:
					existsOrg = false
				}
			}
		}
		if existsOrg {
			return nil
		}
		return errors.ThrowPreconditionFailed(nil, "EVENT-3OfIm", "Errors.Project.OrgNotExisting")
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
			case usr_model.UserAdded, usr_model.UserRegistered, usr_model.HumanRegistered, usr_model.MachineAdded, usr_model.HumanAdded:
				existsUser = true
			case usr_model.UserRemoved:
				existsUser = false
			}
		}
	}
	if existsUser {
		return nil
	}
	return errors.ThrowPreconditionFailed(nil, "EVENT-3OfIm", "Errors.Project.UserNotExisting")
}
