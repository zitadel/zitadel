package eventsourcing

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	key_model "github.com/caos/zitadel/internal/key/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

const (
	projectOwnerRole       = "PROJECT_OWNER"
	projectOwnerGlobalRole = "PROJECT_OWNER_GLOBAL"
)

type ProjectEventstore struct {
	es_int.Eventstore
	projectCache  *ProjectCache
	passwordAlg   crypto.HashAlgorithm
	pwGenerator   crypto.Generator
	idGenerator   id.Generator
	ClientKeySize int
}

type ProjectConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartProject(conf ProjectConfig, systemDefaults sd.SystemDefaults) (*ProjectEventstore, error) {
	projectCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	passwordAlg := crypto.NewBCrypt(systemDefaults.SecretGenerators.PasswordSaltCost)
	pwGenerator := crypto.NewHashGenerator(systemDefaults.SecretGenerators.ClientSecretGenerator, passwordAlg)
	return &ProjectEventstore{
		Eventstore:    conf.Eventstore,
		projectCache:  projectCache,
		passwordAlg:   passwordAlg,
		pwGenerator:   pwGenerator,
		idGenerator:   id.SonyFlakeGenerator,
		ClientKeySize: int(systemDefaults.SecretGenerators.ClientKeySize),
	}, nil
}

func (es *ProjectEventstore) ProjectByID(ctx context.Context, id string) (*proj_model.Project, error) {
	project := es.projectCache.getProject(id)

	query, err := ProjectByIDQuery(project.AggregateID, project.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, project.AppendEvents, query)
	if err != nil && !(caos_errs.IsNotFound(err) && project.Sequence != 0) {
		return nil, err
	}
	if project.State == int32(proj_model.ProjectStateRemoved) {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-dG8ie", "Errors.Project.NotFound")
	}
	es.projectCache.cacheProject(project)
	return model.ProjectToModel(project), nil
}

func (es *ProjectEventstore) ProjectEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	query, err := ProjectByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return es.FilterEvents(ctx, query)
}

func (es *ProjectEventstore) CreateProject(ctx context.Context, project *proj_model.Project, global bool) (*proj_model.Project, error) {
	if !project.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-IOVCC", "Errors.Project.Invalid")
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	project.AggregateID = id
	project.State = proj_model.ProjectStateActive
	repoProject := model.ProjectFromModel(project)
	projectRole := projectOwnerRole
	if global {
		projectRole = projectOwnerGlobalRole
	}
	member := &model.ProjectMember{
		UserID: authz.GetCtxData(ctx).UserID,
		Roles:  []string{projectRole},
	}

	createAggregate := ProjectCreateAggregate(es.AggregateCreator(), repoProject, member)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.projectCache.cacheProject(repoProject)
	return model.ProjectToModel(repoProject), nil
}

func (es *ProjectEventstore) UpdateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	if !project.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7eBD6", "Errors.Project.Invalid")
	}
	existingProject, err := es.ProjectByID(ctx, project.AggregateID)
	if err != nil {
		return nil, err
	}
	repoExisting := model.ProjectFromModel(existingProject)
	repoNew := model.ProjectFromModel(project)

	updateAggregate := ProjectUpdateAggregate(es.AggregateCreator(), repoExisting, repoNew)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.projectCache.cacheProject(repoExisting)
	return model.ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) DeactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	existingProject, err := es.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !existingProject.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "Errors.Project.NotActive")
	}

	repoExisting := model.ProjectFromModel(existingProject)
	aggregate := ProjectDeactivateAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoExisting)
	return model.ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ReactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	existingProject, err := es.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingProject.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "Errors.Project.NotInactive")
	}

	repoExisting := model.ProjectFromModel(existingProject)
	aggregate := ProjectReactivateAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoExisting)
	return model.ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) RemoveProject(ctx context.Context, proj *proj_model.Project) error {
	project, aggregate, err := es.PrepareRemoveProject(ctx, proj)
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, project.AppendEvents, aggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(project)
	return nil
}

func (es *ProjectEventstore) PrepareRemoveProject(ctx context.Context, proj *proj_model.Project) (*model.Project, *es_models.Aggregate, error) {
	if proj.AggregateID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Akbov", "Errors.ProjectInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, proj.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	repoProject := model.ProjectFromModel(existingProject)
	projectAggregate, err := ProjectRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoProject)
	if err != nil {
		return nil, nil, err
	}
	return repoProject, projectAggregate, nil
}

func (es *ProjectEventstore) ProjectMemberByIDs(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if member.UserID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ld93d", "Errors.Project.UserIDMissing")
	}
	project, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}

	if _, m := project.GetMember(member.UserID); m != nil {
		return m, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-3udjs", "Errors.Project.MemberNotFound")
}

func (es *ProjectEventstore) AddProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-2OWkC", "Errors.Project.MemberInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existingProject.GetMember(member.UserID); m != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "Errors.Project.MemberAlreadyExists")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoMember := model.ProjectMemberFromModel(member)

	addAggregate := ProjectMemberAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)

	if _, m := model.GetProjectMember(repoProject.Members, member.UserID); m != nil {
		return model.ProjectMemberToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-rfQWv", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Buh04", "Errors.Project.MemberInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existingProject.GetMember(member.UserID); m == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-oe39f", "Errors.Project.MemberNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoMember := model.ProjectMemberFromModel(member)

	projectAggregate := ProjectMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)

	if _, m := model.GetProjectMember(repoProject.Members, member.UserID); m != nil {
		return model.ProjectMemberToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-pLyzi", "Errors.Internal")
}

func (es *ProjectEventstore) RemoveProjectMember(ctx context.Context, member *proj_model.ProjectMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-d43fs", "Errors.Project.MemberInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}
	if _, m := existingProject.GetMember(member.UserID); m == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-swf34", "Errors.Project.MemberNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoMember := model.ProjectMemberFromModel(member)

	projectAggregate := ProjectMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return err
}

func (es *ProjectEventstore) PrepareRemoveProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*model.ProjectMember, *es_models.Aggregate, error) {
	if member.UserID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-tCXHE", "Errors.Project.MemberInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	if _, m := existingProject.GetMember(member.UserID); m == nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-wPcg5", "Errors.Project.MemberNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoMember := model.ProjectMemberFromModel(member)

	projectAggregate := ProjectMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	agg, err := projectAggregate(ctx)
	if err != nil {
		return nil, nil, err
	}

	return repoMember, agg, err
}

func (es *ProjectEventstore) AddProjectRoles(ctx context.Context, roles ...*proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	if roles == nil || len(roles) == 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-uOJAs", "Errors.Project.MinimumOneRoleNeeded")
	}
	for _, role := range roles {
		if !role.IsValid() {
			return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-iduG4", "Errors.Project.RoleInvalid")
		}
	}
	existingProject, err := es.ProjectByID(ctx, roles[0].AggregateID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if existingProject.ContainsRole(role) {
			return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-sk35t", "Errors.Project.RoleAlreadyExists")
		}
	}

	repoProject := model.ProjectFromModel(existingProject)
	repoRoles := model.ProjectRolesFromModel(roles)
	projectAggregate := ProjectRoleAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoRoles...)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	if len(repoRoles) > 1 {
		return nil, nil
	}
	es.projectCache.cacheProject(repoProject)
	if _, r := model.GetProjectRole(repoProject.Roles, repoRoles[0].Key); r != nil {
		return model.ProjectRoleToModel(r), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sie83", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeProjectRole(ctx context.Context, role *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	if !role.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9die3", "Errors.Project.RoleInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, role.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existingProject.ContainsRole(role) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die34", "Errors.Project.RoleNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoRole := model.ProjectRoleFromModel(role)
	projectAggregate := ProjectRoleChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoRole)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}

	es.projectCache.cacheProject(repoProject)

	if _, r := model.GetProjectRole(repoProject.Roles, role.Key); r != nil {
		return model.ProjectRoleToModel(r), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sl1or", "Errors.Internal")
}

func (es *ProjectEventstore) PrepareRemoveProjectRole(ctx context.Context, role *proj_model.ProjectRole) (*model.Project, *es_models.Aggregate, error) {
	if role.Key == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-id823", "Errors.Project.RoleInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, role.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	if !existingProject.ContainsRole(role) {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-oe823", "Errors.Project.RoleNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoRole := model.ProjectRoleFromModel(role)
	grants := es.RemoveRoleFromGrants(repoProject, role.Key)
	projectAggregate, err := ProjectRoleRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoProject, repoRole, grants)
	if err != nil {
		return nil, nil, err
	}
	return repoProject, projectAggregate, nil
}

func (es *ProjectEventstore) RemoveRoleFromGrants(existingProject *model.Project, roleKey string) []*model.ProjectGrant {
	grants := make([]*model.ProjectGrant, len(existingProject.Grants))
	for i, grant := range existingProject.Grants {
		newGrant := *grant
		roles := make([]string, 0)
		for _, role := range newGrant.RoleKeys {
			if role != roleKey {
				roles = append(roles, role)
			}
		}
		newGrant.RoleKeys = roles
		grants[i] = &newGrant
	}
	return grants
}

func (es *ProjectEventstore) RemoveProjectRole(ctx context.Context, role *proj_model.ProjectRole) error {
	project, aggregate, err := es.PrepareRemoveProjectRole(ctx, role)
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, project.AppendEvents, aggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(project)
	return nil
}

func (es *ProjectEventstore) ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*proj_model.ProjectChanges, error) {
	query := ChangesQuery(id, lastSequence, limit, sortAscending)

	events, err := es.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-328b1", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-FpQqK", "Errors.Changes.NotFound")
	}

	changes := make([]*proj_model.ProjectChange, len(events))

	for i, event := range events {
		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-qxIR7").OnError(err).Debug("unable to parse timestamp")
		change := &proj_model.ProjectChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
		}

		if event.Data != nil {
			var data interface{}
			if strings.Contains(change.EventType, "application") {
				data = new(model.Application)
			} else {
				data = new(model.Project)
			}
			err = json.Unmarshal(event.Data, data)
			logging.Log("EVENT-NCkpN").OnError(err).Debug("unable to unmarshal data")
			change.Data = data
		}

		changes[i] = change
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &proj_model.ProjectChanges{
		Changes:      changes,
		LastSequence: lastSequence,
	}, nil
}

func ChangesQuery(projectID string, latestSequence, limit uint64, sortAscending bool) *es_models.SearchQuery {
	query := es_models.NewSearchQuery().
		AggregateTypeFilter(model.ProjectAggregate)
	if !sortAscending {
		query.OrderDesc()
	}

	query.LatestSequenceFilter(latestSequence).
		AggregateIDFilter(projectID).
		SetLimit(limit)
	return query
}

func (es *ProjectEventstore) ApplicationByIDs(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ld93d", "Errors.Project.IDMissing")
	}
	project, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if _, a := project.GetApp(appID); a != nil {
		return a, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-8ei2s", "Errors.Project.AppNotFound")
}

func (es *ProjectEventstore) AddApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	if app == nil || !app.IsValid(true) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9eidw", "Errors.Project.AppInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return nil, err
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	app.AppID = id

	var stringPw string
	if app.OIDCConfig != nil {
		app.OIDCConfig.AppID = id
		err := app.OIDCConfig.GenerateNewClientID(es.idGenerator, existingProject)
		if err != nil {
			return nil, err
		}
		stringPw, err = app.OIDCConfig.GenerateClientSecretIfNeeded(es.pwGenerator)
		if err != nil {
			return nil, err
		}
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoApp := model.AppFromModel(app)

	addAggregate := ApplicationAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoApp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		converted := model.AppToModel(a)
		converted.OIDCConfig.ClientSecretString = stringPw
		converted.OIDCConfig.FillCompliance()
		return converted, nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-GvPct", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	if app == nil || !app.IsValid(false) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dieuw", "Errors.Project.AppInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, app := existingProject.GetApp(app.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die83", "Errors.Project.AppNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoApp := model.AppFromModel(app)

	projectAggregate := ApplicationChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoApp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		return model.AppToModel(a), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-dksi8", "Errors.Internal")
}

func (es *ProjectEventstore) RemoveApplication(ctx context.Context, app *proj_model.Application) error {
	if app.AppID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-id832", "Errors.Project.IDMissing")
	}
	existingProject, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return err
	}
	if _, app := existingProject.GetApp(app.AppID); app == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83s", "Errors.Project.AppNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	appRepo := model.AppFromModel(app)
	projectAggregate := ApplicationRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, appRepo)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}

func (es *ProjectEventstore) PrepareRemoveApplication(ctx context.Context, app *proj_model.Application) (*model.Application, *es_models.Aggregate, error) {
	if app.AppID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-xu0Wy", "Errors.Project.IDMissing")
	}
	existingProject, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	if _, app := existingProject.GetApp(app.AppID); app == nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-gaOD2", "Errors.Project.AppNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	appRepo := model.AppFromModel(app)
	projectAggregate := ApplicationRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, appRepo)
	agg, err := projectAggregate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return appRepo, agg, nil
}

func (es *ProjectEventstore) ApplicationChanges(ctx context.Context, projectID string, appID string, lastSequence uint64, limit uint64, sortAscending bool) (*proj_model.ApplicationChanges, error) {
	query := ChangesQuery(projectID, lastSequence, limit, sortAscending)

	events, err := es.Eventstore.FilterEvents(ctx, query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-sw6Ku", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-9IHLP", "Errors.Changes.NotFound")
	}

	result := make([]*proj_model.ApplicationChange, 0)
	for _, event := range events {
		if !strings.Contains(event.Type.String(), "application") || event.Data == nil {
			continue
		}

		app := new(model.Application)
		err := json.Unmarshal(event.Data, app)
		logging.Log("EVENT-GIiKD").OnError(err).Debug("unable to unmarshal data")
		if app.AppID != appID {
			continue
		}

		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-MJzeN").OnError(err).Debug("unable to parse timestamp")

		result = append(result, &proj_model.ApplicationChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
			Data:       app,
		})
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &proj_model.ApplicationChanges{
		Changes:      result,
		LastSequence: lastSequence,
	}, nil
}

func (es *ProjectEventstore) DeactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	if appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dlp9e", "Errors.Project.IDMissing")
	}
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	app := &proj_model.Application{AppID: appID}
	if _, app := existingProject.GetApp(app.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-slpe9", "Errors.Project.AppNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoApp := model.AppFromModel(app)

	projectAggregate := ApplicationDeactivatedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoApp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		return model.AppToModel(a), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sie83", "Errors.Internal")
}

func (es *ProjectEventstore) ReactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	if appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-0odi2", "Errors.Project.IDMissing")
	}
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	app := &proj_model.Application{AppID: appID}
	if _, app := existingProject.GetApp(app.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ld92d", "Errors.Project.AppNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoApp := model.AppFromModel(app)

	projectAggregate := ApplicationReactivatedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoApp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		return model.AppToModel(a), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sld93", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeOIDCConfig(ctx context.Context, config *proj_model.OIDCConfig) (*proj_model.OIDCConfig, error) {
	if config == nil || !config.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-du834", "Errors.Project.OIDCConfigInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, config.AggregateID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(config.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dkso8", "Errors.Project.AppNoExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-98uje", "Errors.Project.AppIsNotOIDC")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoConfig := model.OIDCConfigFromModel(config)

	projectAggregate := OIDCConfigChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoConfig)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		return model.OIDCConfigToModel(a.OIDCConfig), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-dk87s", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeOIDCConfigSecret(ctx context.Context, projectID, appID string) (*proj_model.OIDCConfig, error) {
	if appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7ue34", "Errors.Project.OIDCConfigInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(appID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9odi4", "Errors.Project.AppNotExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dile4", "Errors.Project.AppIsNotOIDC")
	}
	if app.OIDCConfig.AuthMethodType == proj_model.OIDCAuthMethodTypeNone {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-GDrg2", "Errors.Project.OIDCAuthMethodNoneSecret")
	}
	repoProject := model.ProjectFromModel(existingProject)

	stringPw, err := app.OIDCConfig.GenerateNewClientSecret(es.pwGenerator)
	if err != nil {
		return nil, err
	}

	projectAggregate := OIDCConfigSecretChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, appID, app.OIDCConfig.ClientSecret)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)

	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		config := model.OIDCConfigToModel(a.OIDCConfig)
		config.ClientSecretString = stringPw
		return config, nil
	}

	return nil, caos_errs.ThrowInternal(nil, "EVENT-dk87s", "Errors.Internal")
}

func (es *ProjectEventstore) VerifyOIDCClientSecret(ctx context.Context, projectID, appID string, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if appID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-H3RT2", "Errors.Project.RequiredFieldsMissing")
	}
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(appID); app == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-D6hba", "Errors.Project.AppNoExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-huywq", "Errors.Project.AppIsNotOIDC")
	}

	ctx, spanHash := tracing.NewSpan(ctx)
	err = crypto.CompareHash(app.OIDCConfig.ClientSecret, []byte(secret), es.passwordAlg)
	spanHash.EndWithError(err)
	if err == nil {
		err = es.setOIDCClientSecretCheckResult(ctx, existingProject, app.AppID, OIDCClientSecretCheckSucceededAggregate)
		logging.Log("EVENT-AE1vf").OnError(err).Warn("could not push event OIDCClientSecretCheckSucceeded")
		return nil
	}
	err = es.setOIDCClientSecretCheckResult(ctx, existingProject, app.AppID, OIDCClientSecretCheckFailedAggregate)
	logging.Log("EVENT-GD1gh").OnError(err).Warn("could not push event OIDCClientSecretCheckFailed")
	return caos_errs.ThrowInvalidArgument(nil, "EVENT-wg24q", "Errors.Project.OIDCSecretInvalid")
}

func (es *ProjectEventstore) setOIDCClientSecretCheckResult(ctx context.Context, project *proj_model.Project, appID string, check func(*es_models.AggregateCreator, *model.Project, string) es_sdk.AggregateFunc) error {
	repoProject := model.ProjectFromModel(project)
	agg := check(es.AggregateCreator(), repoProject, appID)
	err := es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}

func (es *ProjectEventstore) AddClientKey(ctx context.Context, key *proj_model.ClientKey) (*proj_model.ClientKey, error) {
	existingProject, err := es.ProjectByID(ctx, key.AggregateID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(key.ApplicationID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Dbf32", "Errors.Project.AppNoExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Dff54", "Errors.Project.AppIsNotOIDC")
	}
	key.ClientID = app.OIDCConfig.ClientID
	key.KeyID, err = es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	if key.ExpirationDate.IsZero() {
		key.ExpirationDate, err = key_model.DefaultExpiration()
		if err != nil {
			logging.Log("EVENT-Adgf2").WithError(err).Warn("unable to set default date")
			return nil, errors.ThrowInternal(err, "EVENT-j68fg", "Errors.Internal")
		}
	}
	if key.ExpirationDate.Before(time.Now()) {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-C6YV5", "Errors.MachineKey.ExpireBeforeNow")
	}

	repoProject := model.ProjectFromModel(existingProject)
	repoKey := model.ClientKeyFromModel(key)
	err = repoKey.GenerateClientKeyPair(es.ClientKeySize)
	if err != nil {
		return nil, err
	}
	agg := OIDCApplicationKeyAddedAggregate(es.AggregateCreator(), repoProject, repoKey)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, agg)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)

	return model.ClientKeyToModel(repoKey), nil
}

func (es *ProjectEventstore) RemoveApplicationKey(ctx context.Context, projectID, applicationID, keyID string) error {
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(applicationID); app == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-ADfzz", "Errors.Project.AppNoExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-ADffh", "Errors.Project.AppIsNotOIDC")
	}
	if _, key := app.GetKey(keyID); key == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-D2Sff", "Errors.Project.AppKeyNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	agg := OIDCApplicationKeyRemovedAggregate(es.AggregateCreator(), repoProject, keyID)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}

func (es *ProjectEventstore) TokenAdded(ctx context.Context, token *proj_model.Token) (*proj_model.Token, error) {
	existingProject, err := es.ProjectByID(ctx, token.AggregateID)
	if err != nil {
		return nil, err
	}
	token.TokenID, err = es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoToken := model.TokenFromModel(token)
	agg := OIDCApplicationTokenAddedAggregate(es.AggregateCreator(), repoProject, repoToken)
	err = es_sdk.Push(ctx, es.PushAggregates, repoToken.AppendEvents, agg)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	return model.TokenToModel(repoToken), nil
}

func (es *ProjectEventstore) ProjectGrantByIDs(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-e8die", "Errors.Project.IDMissing")
	}
	project, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if _, g := project.GetGrant(grantID); g != nil {
		return g, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-slo45", "Errors.Project.GrantNotFound")
}

func (es *ProjectEventstore) AddProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) (*proj_model.ProjectGrant, error) {
	if grant == nil || !grant.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-37dhs", "Errors.Project.GrantInvalid")
	}
	project, err := es.ProjectByID(ctx, grant.AggregateID)
	if err != nil {
		return nil, err
	}
	if project.ContainsGrantForOrg(grant.GrantedOrgID) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7ug4g", "Errors.Project.GrantAlreadyExists")
	}
	if !project.ContainsRoles(grant.RoleKeys) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83d", "Errors.Project.GrantHasNotExistingRole")
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	grant.GrantID = id

	repoProject := model.ProjectFromModel(project)
	repoGrant := model.GrantFromModel(grant)

	addAggregate := ProjectGrantAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	if _, g := model.GetProjectGrant(repoProject.Grants, grant.GrantID); g != nil {
		return model.GrantToModel(g), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sk3t5", "Errors.Internal")
}

func (es *ProjectEventstore) PrepareChangeProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) (*model.Project, func(ctx context.Context) (*es_models.Aggregate, error), []string, error) {
	if grant == nil && grant.GrantID == "" {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-8sie3", "Errors.Project.GrantInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, grant.AggregateID)
	if err != nil {
		return nil, nil, nil, err
	}
	_, existingGrant := existingProject.GetGrant(grant.GrantID)
	if existingGrant == nil {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die83", "Errors.Project.GrantNotExisting")
	}
	if !existingProject.ContainsRoles(grant.RoleKeys) {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83d", "Error.Project.GrantHasNotExistingRole")
	}
	removedRoles := existingGrant.GetRemovedRoles(grant.RoleKeys)
	repoProject := model.ProjectFromModel(existingProject)
	repoGrant := model.GrantFromModel(grant)

	projectAggregate := ProjectGrantChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoGrant)
	return repoProject, projectAggregate, removedRoles, nil
}

func (es *ProjectEventstore) RemoveProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) error {
	repoProject, projectAggregate, err := es.PrepareRemoveProjectGrant(ctx, grant)
	if err != nil {
		return err
	}
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}

func (es *ProjectEventstore) RemoveProjectGrants(ctx context.Context, grants ...*proj_model.ProjectGrant) error {
	aggregates := make([]*es_models.Aggregate, len(grants))
	for i, grant := range grants {
		_, projectAggregate, err := es.PrepareRemoveProjectGrant(ctx, grant)
		if err != nil {
			return err
		}
		agg, err := projectAggregate(ctx)
		if err != nil {
			return err
		}
		aggregates[i] = agg
	}
	return es_sdk.PushAggregates(ctx, es.PushAggregates, nil, aggregates...)
}

func (es *ProjectEventstore) PrepareRemoveProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) (*model.Project, func(ctx context.Context) (*es_models.Aggregate, error), error) {
	if grant.GrantID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-8eud6", "Errors.Project.IDMissing")
	}
	existingProject, err := es.ProjectByID(ctx, grant.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	if _, g := existingProject.GetGrant(grant.GrantID); g == nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9ie3s", "Errors.Project.GrantNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	grantRepo := model.GrantFromModel(grant)
	projectAggregate := ProjectGrantRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, grantRepo)
	return repoProject, projectAggregate, nil
}

func (es *ProjectEventstore) DeactivateProjectGrant(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7due2", "Errors.Project.IDMissing")
	}
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	grant := &proj_model.ProjectGrant{GrantID: grantID}
	if _, g := existingProject.GetGrant(grant.GrantID); g == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-slpe9", "Errors.Project.GrantNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoGrant := model.GrantFromModel(grant)

	projectAggregate := ProjectGrantDeactivatedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, g := model.GetProjectGrant(repoProject.Grants, grant.GrantID); g != nil {
		return model.GrantToModel(g), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sie83", "Errors.Internal")
}

func (es *ProjectEventstore) ReactivateProjectGrant(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-d7suw", "Errors.Project.IDMissing")
	}
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	grant := &proj_model.ProjectGrant{GrantID: grantID}
	if _, g := existingProject.GetGrant(grant.GrantID); g == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-0spew", "Errors.Project.GrantNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoGrant := model.GrantFromModel(grant)

	projectAggregate := ProjectGrantReactivatedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)

	if _, g := model.GetProjectGrant(repoProject.Grants, grant.GrantID); g != nil {
		return model.GrantToModel(g), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-9osjw", "Errors.Internal")
}

func (es *ProjectEventstore) ProjectGrantMemberByIDs(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	if member.GrantID == "" || member.UserID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-8diw2", "Errors.Project.UserIDMissing")
	}
	project, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, g := project.GetGrant(member.GrantID); g != nil {
		if _, m := g.GetMember(member.UserID); m != nil {
			return m, nil
		}
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-LxiBI", "Errors.Project.MemberNotFound")
}

func (es *ProjectEventstore) AddProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-0dor4", "Errors.Project.MemberInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if existingProject.ContainsGrantMember(member) {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-8die3", "Errors.Project.MemberAlreadyExists")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoMember := model.GrantMemberFromModel(member)

	addAggregate := ProjectGrantMemberAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, g := model.GetProjectGrant(repoProject.Grants, member.GrantID); g != nil {
		if _, m := model.GetProjectGrantMember(g.Members, member.UserID); m != nil {
			return model.GrantMemberToModel(m), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-BBcGD", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dkw35", "Errors.Project.MemberInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existingProject.ContainsGrantMember(member) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-8dj4s", "Errors.Project.MemberNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoMember := model.GrantMemberFromModel(member)

	projectAggregate := ProjectGrantMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, g := model.GetProjectGrant(repoProject.Grants, member.GrantID); g != nil {
		if _, m := model.GetProjectGrantMember(g.Members, member.UserID); m != nil {
			return model.GrantMemberToModel(m), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-s8ur3", "Errors.Internal")
}

func (es *ProjectEventstore) RemoveProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-8su4r", "Errors.Project.MemberInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}
	if !existingProject.ContainsGrantMember(member) {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-9ode4", "Errors.Project.MemberNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoMember := model.GrantMemberFromModel(member)

	projectAggregate := ProjectGrantMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return err
}
