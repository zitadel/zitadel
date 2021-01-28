package eventstore

import (
	"context"
	"strings"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	key_model "github.com/caos/zitadel/internal/key/model"
	key_view_model "github.com/caos/zitadel/internal/key/repository/view/model"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_proj_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_grant_model "github.com/caos/zitadel/internal/usergrant/model"
	usr_grant_event "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing"
)

type ProjectRepo struct {
	es_int.Eventstore
	SearchLimit     uint64
	ProjectEvents   *proj_event.ProjectEventstore
	UserGrantEvents *usr_grant_event.UserGrantEventStore
	UserEvents      *usr_event.UserEventstore
	IAMEvents       *iam_event.IAMEventstore
	View            *view.View
	Roles           []string
	IAMID           string
}

func (repo *ProjectRepo) ProjectByID(ctx context.Context, id string) (*proj_model.ProjectView, error) {
	project, viewErr := repo.View.ProjectByID(id)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		project = new(model.ProjectView)
	}

	events, esErr := repo.ProjectEvents.ProjectEventsByID(ctx, id, project.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-8yfKu", "Errors.Project.NotFound")
	}

	if esErr != nil {
		logging.Log("EVENT-V9x1V").WithError(viewErr).Debug("error retrieving new events")
		return model.ProjectToModel(project), nil
	}

	viewProject := *project
	for _, event := range events {
		err := project.AppendEvent(event)
		if err != nil {
			return model.ProjectToModel(&viewProject), nil
		}
	}
	if viewProject.State == int32(proj_model.ProjectStateRemoved) {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-3Mo0s", "Errors.Project.NotFound")
	}
	return model.ProjectToModel(project), nil
}

func (repo *ProjectRepo) CreateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	ctxData := authz.GetCtxData(ctx)
	iam, err := repo.IAMEvents.IAMByID(ctx, repo.IAMID)
	if err != nil {
		return nil, err
	}
	return repo.ProjectEvents.CreateProject(ctx, project, iam.GlobalOrgID == ctxData.OrgID)
}

func (repo *ProjectRepo) UpdateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	return repo.ProjectEvents.UpdateProject(ctx, project)
}

func (repo *ProjectRepo) DeactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	return repo.ProjectEvents.DeactivateProject(ctx, id)
}

func (repo *ProjectRepo) ReactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	return repo.ProjectEvents.ReactivateProject(ctx, id)
}

func (repo *ProjectRepo) RemoveProject(ctx context.Context, projectID string) error {
	proj := proj_model.NewProject(projectID)
	aggregates := make([]*es_models.Aggregate, 0)
	project, agg, err := repo.ProjectEvents.PrepareRemoveProject(ctx, proj)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, agg)

	// remove user_grants
	usergrants, err := repo.View.UserGrantsByProjectID(projectID)
	if err != nil {
		return err
	}
	for _, grant := range usergrants {
		_, aggs, err := repo.UserGrantEvents.PrepareRemoveUserGrant(ctx, grant.ID, true)
		if err != nil {
			return err
		}
		for _, agg := range aggs {
			aggregates = append(aggregates, agg)
		}
	}

	return es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, project.AppendEvents, aggregates...)
}

func (repo *ProjectRepo) SearchProjects(ctx context.Context, request *proj_model.ProjectViewSearchRequest) (*proj_model.ProjectViewSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestProjectSequence("")
	logging.Log("EVENT-Edc56").OnError(sequenceErr).Warn("could not read latest project sequence")

	permissions := authz.GetRequestPermissionsFromCtx(ctx)
	if !authz.HasGlobalPermission(permissions) {
		ids := authz.GetAllPermissionCtxIDs(permissions)
		if _, q := request.GetSearchQuery(proj_model.ProjectViewSearchKeyProjectID); q != nil {
			containsID := false
			for _, id := range ids {
				if id == q.Value {
					containsID = true
					break
				}
			}
			if !containsID {
				result := &proj_model.ProjectViewSearchResponse{
					Offset:      request.Offset,
					Limit:       request.Limit,
					TotalResult: uint64(0),
					Result:      []*proj_model.ProjectView{},
				}
				if sequenceErr == nil {
					result.Sequence = sequence.CurrentSequence
					result.Timestamp = sequence.LastSuccessfulSpoolerRun
				}
				return result, nil
			}
		} else {
			request.Queries = append(request.Queries, &proj_model.ProjectViewSearchQuery{Key: proj_model.ProjectViewSearchKeyProjectID, Method: global_model.SearchMethodIsOneOf, Value: ids})
		}
	}

	projects, count, err := repo.View.SearchProjects(request)
	if err != nil {
		return nil, err
	}
	result := &proj_model.ProjectViewSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectsToModel(projects),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) ProjectGrantViewByID(ctx context.Context, grantID string) (project *proj_model.ProjectGrantView, err error) {
	p, err := repo.View.ProjectGrantByID(grantID)
	if err != nil {
		return nil, err
	}
	return model.ProjectGrantToModel(p), nil
}

func (repo *ProjectRepo) ProjectMemberByID(ctx context.Context, projectID, userID string) (*proj_model.ProjectMemberView, error) {
	member, err := repo.View.ProjectMemberByIDs(projectID, userID)
	if err != nil {
		return nil, err
	}
	return model.ProjectMemberToModel(member), nil
}

func (repo *ProjectRepo) AddProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	return repo.ProjectEvents.AddProjectMember(ctx, member)
}

func (repo *ProjectRepo) ChangeProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	return repo.ProjectEvents.ChangeProjectMember(ctx, member)
}

func (repo *ProjectRepo) RemoveProjectMember(ctx context.Context, projectID, userID string) error {
	member := proj_model.NewProjectMember(projectID, userID)
	return repo.ProjectEvents.RemoveProjectMember(ctx, member)
}

func (repo *ProjectRepo) SearchProjectMembers(ctx context.Context, request *proj_model.ProjectMemberSearchRequest) (*proj_model.ProjectMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestProjectMemberSequence("")
	logging.Log("EVENT-3dgt6").OnError(sequenceErr).Warn("could not read latest project member sequence")
	members, count, err := repo.View.SearchProjectMembers(request)
	if err != nil {
		return nil, err
	}
	result := &proj_model.ProjectMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectMembersToModel(members),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) AddProjectRole(ctx context.Context, role *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	return repo.ProjectEvents.AddProjectRoles(ctx, role)
}

func (repo *ProjectRepo) BulkAddProjectRole(ctx context.Context, roles []*proj_model.ProjectRole) error {
	_, err := repo.ProjectEvents.AddProjectRoles(ctx, roles...)
	return err
}

func (repo *ProjectRepo) ChangeProjectRole(ctx context.Context, member *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	return repo.ProjectEvents.ChangeProjectRole(ctx, member)
}

func (repo *ProjectRepo) RemoveProjectRole(ctx context.Context, projectID, key string) error {
	role := proj_model.NewProjectRole(projectID, key)
	aggregates := make([]*es_models.Aggregate, 0)
	project, agg, err := repo.ProjectEvents.PrepareRemoveProjectRole(ctx, role)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, agg)

	usergrants, err := repo.View.UserGrantsByProjectIDAndRoleKey(projectID, key)
	if err != nil {
		return err
	}
	for _, grant := range usergrants {
		changed := &usr_grant_model.UserGrant{
			ObjectRoot: models.ObjectRoot{AggregateID: grant.ID, Sequence: grant.Sequence, ResourceOwner: grant.ResourceOwner},
			RoleKeys:   grant.RoleKeys,
			ProjectID:  grant.ProjectID,
			UserID:     grant.UserID,
		}
		changed.RemoveRoleKeyIfExisting(key)
		_, agg, err := repo.UserGrantEvents.PrepareChangeUserGrant(ctx, changed, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, agg)
	}
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, project.AppendEvents, aggregates...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProjectRepo) SearchProjectRoles(ctx context.Context, projectID string, request *proj_model.ProjectRoleSearchRequest) (*proj_model.ProjectRoleSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.AppendProjectQuery(projectID)
	sequence, sequenceErr := repo.View.GetLatestProjectRoleSequence("")
	logging.Log("LSp0d-47suf").OnError(sequenceErr).Warn("could not read latest project role sequence")
	roles, count, err := repo.View.SearchProjectRoles(request)
	if err != nil {
		return nil, err
	}

	result := &proj_model.ProjectRoleSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.ProjectRolesToModel(roles),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*proj_model.ProjectChanges, error) {
	changes, err := repo.ProjectEvents.ProjectChanges(ctx, id, lastSequence, limit, sortAscending)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierId
		user, _ := repo.UserEvents.UserByID(ctx, change.ModifierId)
		if user != nil {
			if user.Human != nil {
				change.ModifierName = user.DisplayName
			}
			if user.Machine != nil {
				change.ModifierName = user.Machine.Name
			}
		}
	}
	return changes, nil
}

func (repo *ProjectRepo) ApplicationByID(ctx context.Context, projectID, appID string) (*proj_model.ApplicationView, error) {
	app, viewErr := repo.View.ApplicationByID(projectID, appID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		app = new(model.ApplicationView)
		app.ID = appID
	}

	events, esErr := repo.ProjectEvents.ProjectEventsByID(ctx, projectID, app.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-Fshu8", "Errors.Application.NotFound")
	}

	if esErr != nil {
		logging.Log("EVENT-SLCo9").WithError(viewErr).Debug("error retrieving new events")
		return model.ApplicationViewToModel(app), nil
	}

	viewApp := *app
	for _, event := range events {
		err := app.AppendEventIfMyApp(event)
		if err != nil {
			return model.ApplicationViewToModel(&viewApp), nil
		}
		if app.State == int32(proj_model.AppStateRemoved) {
			return nil, caos_errs.ThrowNotFound(nil, "EVENT-Msl96", "Errors.Application.NotFound")
		}
	}
	return model.ApplicationViewToModel(app), nil
}

func (repo *ProjectRepo) AddApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	return repo.ProjectEvents.AddApplication(ctx, app)
}

func (repo *ProjectRepo) ChangeApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	return repo.ProjectEvents.ChangeApplication(ctx, app)
}

func (repo *ProjectRepo) DeactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	return repo.ProjectEvents.DeactivateApplication(ctx, projectID, appID)
}

func (repo *ProjectRepo) ReactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	return repo.ProjectEvents.ReactivateApplication(ctx, projectID, appID)
}

func (repo *ProjectRepo) RemoveApplication(ctx context.Context, projectID, appID string) error {
	app := proj_model.NewApplication(projectID, appID)
	return repo.ProjectEvents.RemoveApplication(ctx, app)
}

func (repo *ProjectRepo) SearchApplications(ctx context.Context, request *proj_model.ApplicationSearchRequest) (*proj_model.ApplicationSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestApplicationSequence("")
	logging.Log("EVENT-SKe8s").OnError(sequenceErr).Warn("could not read latest application sequence")
	apps, count, err := repo.View.SearchApplications(request)
	if err != nil {
		return nil, err
	}
	result := &proj_model.ApplicationSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.ApplicationViewsToModel(apps),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) ApplicationChanges(ctx context.Context, id string, appId string, lastSequence uint64, limit uint64, sortAscending bool) (*proj_model.ApplicationChanges, error) {
	changes, err := repo.ProjectEvents.ApplicationChanges(ctx, id, appId, lastSequence, limit, sortAscending)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierId
		user, _ := repo.UserEvents.UserByID(ctx, change.ModifierId)
		if user != nil {
			if user.Human != nil {
				change.ModifierName = user.DisplayName
			}
			if user.Machine != nil {
				change.ModifierName = user.Machine.Name
			}
		}
	}
	return changes, nil
}

func (repo *ProjectRepo) SearchClientKeys(ctx context.Context, request *key_model.AuthNKeySearchRequest) (*key_model.AuthNKeySearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestAuthNKeySequence("")
	logging.Log("EVENT-ADwgw").OnError(sequenceErr).Warn("could not read latest authn key sequence")
	keys, count, err := repo.View.SearchAuthNKeys(request)
	if err != nil {
		return nil, err
	}
	result := &key_model.AuthNKeySearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      key_view_model.AuthNKeysToModel(keys),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) GetClientKey(ctx context.Context, projectID, applicationID, keyID string) (*key_model.AuthNKeyView, error) {
	key, err := repo.View.AuthNKeyByIDs(applicationID, keyID)
	if err != nil {
		return nil, err
	}
	return key_view_model.AuthNKeyToModel(key), nil
}

func (repo *ProjectRepo) AddClientKey(ctx context.Context, key *proj_model.ClientKey) (*proj_model.ClientKey, error) {
	return repo.ProjectEvents.AddClientKey(ctx, key)
}

func (repo *ProjectRepo) RemoveClientKey(ctx context.Context, projectID, applicationID, keyID string) error {
	return repo.ProjectEvents.RemoveApplicationKey(ctx, projectID, applicationID, keyID)
}

func (repo *ProjectRepo) ChangeOIDCConfig(ctx context.Context, config *proj_model.OIDCConfig) (*proj_model.OIDCConfig, error) {
	return repo.ProjectEvents.ChangeOIDCConfig(ctx, config)
}

func (repo *ProjectRepo) ChangeOIDConfigSecret(ctx context.Context, projectID, appID string) (*proj_model.OIDCConfig, error) {
	return repo.ProjectEvents.ChangeOIDCConfigSecret(ctx, projectID, appID)
}

func (repo *ProjectRepo) ProjectGrantByID(ctx context.Context, grantID string) (*proj_model.ProjectGrantView, error) {
	grant, err := repo.View.ProjectGrantByID(grantID)
	if err != nil {
		return nil, err
	}
	return model.ProjectGrantToModel(grant), nil
}

func (repo *ProjectRepo) SearchProjectGrants(ctx context.Context, request *proj_model.ProjectGrantViewSearchRequest) (*proj_model.ProjectGrantViewSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestProjectGrantSequence("")
	logging.Log("EVENT-Skw9f").OnError(sequenceErr).Warn("could not read latest project grant sequence")
	projects, count, err := repo.View.SearchProjectGrants(request)
	if err != nil {
		return nil, err
	}
	result := &proj_model.ProjectGrantViewSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.ProjectGrantsToModel(projects),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) SearchGrantedProjects(ctx context.Context, request *proj_model.ProjectGrantViewSearchRequest) (*proj_model.ProjectGrantViewSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestProjectGrantSequence("")
	logging.Log("EVENT-Skw9f").OnError(sequenceErr).Warn("could not read latest project grant sequence")

	permissions := authz.GetRequestPermissionsFromCtx(ctx)
	if !authz.HasGlobalPermission(permissions) {
		ids := authz.GetAllPermissionCtxIDs(permissions)
		if _, q := request.GetSearchQuery(proj_model.GrantedProjectSearchKeyGrantID); q != nil {
			containsID := false
			for _, id := range ids {
				if id == q.Value {
					containsID = true
					break
				}
			}
			if !containsID {
				result := &proj_model.ProjectGrantViewSearchResponse{
					Offset:      request.Offset,
					Limit:       request.Limit,
					TotalResult: uint64(0),
					Result:      []*proj_model.ProjectGrantView{},
				}
				if sequenceErr == nil {
					result.Sequence = sequence.CurrentSequence
					result.Timestamp = sequence.LastSuccessfulSpoolerRun
				}
				return result, nil
			}
		} else {
			request.Queries = append(request.Queries, &proj_model.ProjectGrantViewSearchQuery{Key: proj_model.GrantedProjectSearchKeyGrantID, Method: global_model.SearchMethodIsOneOf, Value: ids})
		}
	}

	projects, count, err := repo.View.SearchProjectGrants(request)
	if err != nil {
		return nil, err
	}
	result := &proj_model.ProjectGrantViewSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.ProjectGrantsToModel(projects),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) AddProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.AddProjectGrant(ctx, grant)
}

func (repo *ProjectRepo) ChangeProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) (*proj_model.ProjectGrant, error) {
	project, aggFunc, removedRoles, err := repo.ProjectEvents.PrepareChangeProjectGrant(ctx, grant)
	if err != nil {
		return nil, err
	}
	agg, err := aggFunc(ctx)
	if err != nil {
		return nil, err
	}
	aggregates := make([]*es_models.Aggregate, 0)
	aggregates = append(aggregates, agg)

	usergrants, err := repo.View.UserGrantsByProjectID(grant.AggregateID)
	if err != nil {
		return nil, err
	}
	for _, grant := range usergrants {
		changed := &usr_grant_model.UserGrant{
			ObjectRoot: models.ObjectRoot{AggregateID: grant.ID, Sequence: grant.Sequence, ResourceOwner: grant.ResourceOwner},
			RoleKeys:   grant.RoleKeys,
			ProjectID:  grant.ProjectID,
			UserID:     grant.UserID,
		}
		roleDeleted := changed.RemoveRoleKeysIfExisting(removedRoles)
		if roleDeleted {
			_, agg, err := repo.UserGrantEvents.PrepareChangeUserGrant(ctx, changed, true)
			if err != nil {
				return nil, err
			}
			aggregates = append(aggregates, agg)
		}
	}
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, project.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}
	if _, g := es_proj_model.GetProjectGrant(project.Grants, grant.GrantID); g != nil {
		return es_proj_model.GrantToModel(g), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-dksi8", "Could not find app in list")
}

func (repo *ProjectRepo) DeactivateProjectGrant(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.DeactivateProjectGrant(ctx, projectID, grantID)
}

func (repo *ProjectRepo) ReactivateProjectGrant(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.ReactivateProjectGrant(ctx, projectID, grantID)
}

func (repo *ProjectRepo) RemoveProjectGrant(ctx context.Context, projectID, grantID string) error {
	grant, err := repo.ProjectEvents.ProjectGrantByIDs(ctx, projectID, grantID)
	if err != nil {
		return err
	}
	aggregates := make([]*es_models.Aggregate, 0)
	project, aggFunc, err := repo.ProjectEvents.PrepareRemoveProjectGrant(ctx, grant)
	if err != nil {
		return err
	}
	agg, err := aggFunc(ctx)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, agg)

	usergrants, err := repo.View.UserGrantsByOrgIDAndProjectID(grant.GrantedOrgID, projectID)
	if err != nil {
		return err
	}
	for _, grant := range usergrants {
		_, grantAggregates, err := repo.UserGrantEvents.PrepareRemoveUserGrant(ctx, grant.ID, true)
		if err != nil {
			return err
		}
		for _, agg := range grantAggregates {
			aggregates = append(aggregates, agg)
		}
	}
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, project.AppendEvents, aggregates...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProjectRepo) ProjectGrantMemberByID(ctx context.Context, projectID, userID string) (*proj_model.ProjectGrantMemberView, error) {
	member, err := repo.View.ProjectGrantMemberByIDs(projectID, userID)
	if err != nil {
		return nil, err
	}
	return model.ProjectGrantMemberToModel(member), nil
}

func (repo *ProjectRepo) AddProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	return repo.ProjectEvents.AddProjectGrantMember(ctx, member)
}

func (repo *ProjectRepo) ChangeProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	return repo.ProjectEvents.ChangeProjectGrantMember(ctx, member)
}

func (repo *ProjectRepo) RemoveProjectGrantMember(ctx context.Context, projectID, grantID, userID string) error {
	member := proj_model.NewProjectGrantMember(projectID, grantID, userID)
	return repo.ProjectEvents.RemoveProjectGrantMember(ctx, member)
}

func (repo *ProjectRepo) SearchProjectGrantMembers(ctx context.Context, request *proj_model.ProjectGrantMemberSearchRequest) (*proj_model.ProjectGrantMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestProjectGrantMemberSequence("")
	logging.Log("EVENT-Du8sk").OnError(sequenceErr).Warn("could not read latest project grant sequence")
	members, count, err := repo.View.SearchProjectGrantMembers(request)
	if err != nil {
		return nil, err
	}
	result := &proj_model.ProjectGrantMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectGrantMembersToModel(members),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) GetProjectMemberRoles(ctx context.Context) ([]string, error) {
	iam, err := repo.IAMEvents.IAMByID(ctx, repo.IAMID)
	if err != nil {
		return nil, err
	}
	roles := make([]string, 0)
	global := authz.GetCtxData(ctx).OrgID == iam.GlobalOrgID
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "PROJECT") && !strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			if global && !strings.HasSuffix(roleMap, "GLOBAL") {
				continue
			}
			roles = append(roles, roleMap)
		}
	}
	return roles, nil
}

func (repo *ProjectRepo) GetProjectGrantMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}
