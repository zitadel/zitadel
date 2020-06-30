package eventsourcing

import (
	"context"
	"encoding/json"
	"log"
	"strings"

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
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
)

const (
	projectOwnerRole = "PROJECT_OWNER"
)

type ProjectEventstore struct {
	es_int.Eventstore
	projectCache *ProjectCache
	passwordAlg  crypto.HashAlgorithm
	pwGenerator  crypto.Generator
	idGenerator  id.Generator
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
		Eventstore:   conf.Eventstore,
		projectCache: projectCache,
		passwordAlg:  passwordAlg,
		pwGenerator:  pwGenerator,
		idGenerator:  id.SonyFlakeGenerator,
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

func (es *ProjectEventstore) CreateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	if !project.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Errors.Project.Invalid")
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	project.AggregateID = id
	project.State = proj_model.ProjectStateActive
	repoProject := model.ProjectFromModel(project)
	member := &model.ProjectMember{
		UserID: authz.GetCtxData(ctx).UserID,
		Roles:  []string{projectOwnerRole},
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Errors.Project.Invalid")
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
	existing, err := es.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !existing.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "Errors.Project.NotActive")
	}

	repoExisting := model.ProjectFromModel(existing)
	aggregate := ProjectDeactivateAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoExisting)
	return model.ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ReactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	existing, err := es.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "Errors.Project.NotInactive")
	}

	repoExisting := model.ProjectFromModel(existing)
	aggregate := ProjectReactivateAggregate(es.AggregateCreator(), repoExisting)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoExisting)
	return model.ProjectToModel(repoExisting), nil
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Errors.Project.MemberInvalid")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "Errors.Project.MemberAlreadyExists")
	}
	repoProject := model.ProjectFromModel(existing)
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
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Errors.Project.MemberInvalid")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-oe39f", "Errors.Project.MemberNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
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
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Errors.Internal")
}

func (es *ProjectEventstore) RemoveProjectMember(ctx context.Context, member *proj_model.ProjectMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-d43fs", "Errors.Project.MemberInvalid")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-swf34", "Errors.Project.MemberNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
	repoMember := model.ProjectMemberFromModel(member)

	projectAggregate := ProjectMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return err
}

func (es *ProjectEventstore) AddProjectRoles(ctx context.Context, roles ...*proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	if roles == nil || len(roles) == 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-uOJAs", "Errors.Project.MinimumOneRoleNeeded")
	}
	for _, role := range roles {
		if !role.IsValid() {
			return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-idue3", "Errors.Project.MemberInvalid")
		}
	}
	existing, err := es.ProjectByID(ctx, roles[0].AggregateID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if existing.ContainsRole(role) {
			return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-sk35t", "Errors.Project.RoleAlreadyExists")
		}
	}

	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, role.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existing.ContainsRole(role) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die34", "Errors.Project.RoleNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, role.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	if !existing.ContainsRole(role) {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-oe823", "Errors.Project.RoleNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
	repoRole := model.ProjectRoleFromModel(role)
	grants := es.RemoveRoleFromGrants(repoProject, role.Key)
	projectAggregate, err := ProjectRoleRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoProject, repoRole, grants)
	if err != nil {
		return nil, nil, err
	}
	return repoProject, projectAggregate, nil
}

func (es *ProjectEventstore) RemoveRoleFromGrants(existing *model.Project, roleKey string) []*model.ProjectGrant {
	grants := make([]*model.ProjectGrant, 0)
	for _, grant := range existing.Grants {
		for i, role := range grant.RoleKeys {
			if role == roleKey {
				grant.RoleKeys[i] = grant.RoleKeys[len(grant.RoleKeys)-1]
				grant.RoleKeys[len(grant.RoleKeys)-1] = ""
				grant.RoleKeys = grant.RoleKeys[:len(grant.RoleKeys)-1]
				grants = append(grants, grant)
			}
		}
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

	result := make([]*proj_model.ProjectChange, 0)

	for _, u := range events {
		creationDate, err := ptypes.TimestampProto(u.CreationDate)
		logging.Log("EVENT-qxIR7").OnError(err).Debug("unable to parse timestamp")
		change := &proj_model.ProjectChange{
			ChangeDate: creationDate,
			EventType:  u.Type.String(),
			ModifierId: u.EditorUser,
			Sequence:   u.Sequence,
		}

		projectDummy := proj_model.Project{}
		appDummy := model.Application{}
		change.Data = projectDummy
		if u.Data != nil {
			if strings.Contains(change.EventType, "application") {
				if err := json.Unmarshal(u.Data, &appDummy); err != nil {
					log.Println("Error getting data!", err.Error())
				}
				change.Data = appDummy
			} else {
				if err := json.Unmarshal(u.Data, &projectDummy); err != nil {
					log.Println("Error getting data!", err.Error())
				}
				change.Data = projectDummy
			}
		}

		result = append(result, change)
		if lastSequence < u.Sequence {
			lastSequence = u.Sequence

		}
	}

	changes := &proj_model.ProjectChanges{
		Changes:      result,
		LastSequence: lastSequence,
	}

	return changes, nil
}

func ChangesQuery(projID string, latestSequence, limit uint64, sortAscending bool) *es_models.SearchQuery {
	query := es_models.NewSearchQuery().
		AggregateTypeFilter(model.ProjectAggregate)
	if !sortAscending {
		query.OrderDesc()
	}

	query.LatestSequenceFilter(latestSequence).
		AggregateIDFilter(projID).
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
	existing, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return nil, err
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	app.AppID = id

	var stringPw string
	var cryptoPw *crypto.CryptoValue
	if app.OIDCConfig != nil {
		app.OIDCConfig.AppID = id
		stringPw, cryptoPw, err = generateNewClientSecret(es.pwGenerator)
		if err != nil {
			return nil, err
		}
		app.OIDCConfig.ClientSecret = cryptoPw
		clientID, err := generateNewClientID(es.idGenerator, existing)
		if err != nil {
			return nil, err
		}
		app.OIDCConfig.ClientID = clientID
	}
	repoProject := model.ProjectFromModel(existing)
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
		return converted, nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	if app == nil || !app.IsValid(false) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dieuw", "Errors.Project.AppInvalid")
	}
	existing, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, app := existing.GetApp(app.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die83", "Errors.Project.AppNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return err
	}
	if _, app := existing.GetApp(app.AppID); app == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83s", "Errors.Project.AppNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
	appRepo := model.AppFromModel(app)
	projectAggregate := ApplicationRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, appRepo)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}

func (es *ProjectEventstore) ApplicationChanges(ctx context.Context, id string, secId string, lastSequence uint64, limit uint64, sortAscending bool) (*proj_model.ApplicationChanges, error) {
	query := ChangesQuery(id, lastSequence, limit, sortAscending)

	events, err := es.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-sw6Ku", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-9IHLP", "Errors.Changes.NotFound")
	}

	result := make([]*proj_model.ApplicationChange, 0)

	for _, u := range events {
		creationDate, err := ptypes.TimestampProto(u.CreationDate)
		logging.Log("EVENT-MJzeN").OnError(err).Debug("unable to parse timestamp")
		change := &proj_model.ApplicationChange{
			ChangeDate: creationDate,
			EventType:  u.Type.String(),
			ModifierId: u.EditorUser,
			Sequence:   u.Sequence,
		}
		appendChanges := true

		if change.EventType == model.ApplicationAdded.String() ||
			change.EventType == model.ApplicationChanged.String() ||
			change.EventType == model.OIDCConfigAdded.String() ||
			change.EventType == model.OIDCConfigChanged.String() {
			appDummy := model.Application{}
			if u.Data != nil {
				if err := json.Unmarshal(u.Data, &appDummy); err != nil {
					log.Println("Error getting data!", err.Error())
				}
			}
			change.Data = appDummy
			if appDummy.AppID != secId {
				appendChanges = false
			}
		} else {
			appendChanges = false
		}

		if appendChanges {
			result = append(result, change)
			if lastSequence < u.Sequence {
				lastSequence = u.Sequence
			}
		}
	}

	changes := &proj_model.ApplicationChanges{
		Changes:      result,
		LastSequence: lastSequence,
	}

	return changes, nil
}

func (es *ProjectEventstore) DeactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	if appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dlp9e", "Errors.Project.IDMissing")
	}
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	app := &proj_model.Application{AppID: appID}
	if _, app := existing.GetApp(app.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-slpe9", "Errors.Project.AppNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	app := &proj_model.Application{AppID: appID}
	if _, app := existing.GetApp(app.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ld92d", "Errors.Project.AppNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, config.AggregateID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existing.GetApp(config.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dkso8", "Errors.Project.AppNoExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-98uje", "Errors.Project.AppIsNotOIDC")
	}
	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existing.GetApp(appID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9odi4", "Errors.Project.AppNotExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dile4", "Errors.Project.AppIsNotOIDC")
	}
	repoProject := model.ProjectFromModel(existing)

	stringPw, crypto, err := generateNewClientSecret(es.pwGenerator)
	if err != nil {
		return nil, err
	}

	projectAggregate := OIDCConfigSecretChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, appID, crypto)
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

func (es *ProjectEventstore) VerifyOIDCClientSecret(ctx context.Context, projectID, appID string, secret string) error {
	if appID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-H3RT2", "Errors.Project.RequiredFieldsMissing")
	}
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}
	var app *proj_model.Application
	if _, app = existing.GetApp(appID); app == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-D6hba", "Errors.Project.AppNoExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-huywq", "Errors.Project.AppIsNotOIDC")
	}

	if err := crypto.CompareHash(app.OIDCConfig.ClientSecret, []byte(secret), es.passwordAlg); err == nil {
		return es.setOIDCClientSecretCheckResult(ctx, existing, app.AppID, OIDCClientSecretCheckSucceededAggregate)
	}
	if err := es.setOIDCClientSecretCheckResult(ctx, existing, app.AppID, OIDCClientSecretCheckFailedAggregate); err != nil {
		return err
	}
	return caos_errs.ThrowInvalidArgument(nil, "EVENT-wg24q", "Errors.Internal")
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
	existing, err := es.ProjectByID(ctx, grant.AggregateID)
	if err != nil {
		return nil, err
	}
	if existing.ContainsGrantForOrg(grant.GrantedOrgID) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7ug4g", "Errors.Project.GrantAlreadyExists")
	}
	if !existing.ContainsRoles(grant.RoleKeys) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83d", "Errors.Project.GrantHasNotExistingRole")
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	grant.GrantID = id

	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, grant.AggregateID)
	if err != nil {
		return nil, nil, nil, err
	}
	_, existingGrant := existing.GetGrant(grant.GrantID)
	if existingGrant == nil {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die83", "Errors.Project.GrantNotExisting")
	}
	if !existing.ContainsRoles(grant.RoleKeys) {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83d", "Error.Project.GrantHasNotExistingRole")
	}
	removedRoles := existingGrant.GetRemovedRoles(grant.RoleKeys)
	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, grant.AggregateID)
	if err != nil {
		return nil, nil, err
	}
	if _, g := existing.GetGrant(grant.GrantID); g == nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9ie3s", "Errors.Project.GrantNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
	grantRepo := model.GrantFromModel(grant)
	projectAggregate := ProjectGrantRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, grantRepo)
	return repoProject, projectAggregate, nil
}

func (es *ProjectEventstore) DeactivateProjectGrant(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7due2", "Errors.Project.IDMissing")
	}
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	grant := &proj_model.ProjectGrant{GrantID: grantID}
	if _, g := existing.GetGrant(grant.GrantID); g == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-slpe9", "Errors.Project.GrantNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	grant := &proj_model.ProjectGrant{GrantID: grantID}
	if _, g := existing.GetGrant(grant.GrantID); g == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-0spew", "Errors.Project.GrantNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
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
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-3udjs", "Errors.Project.MemberNotFound")
}

func (es *ProjectEventstore) AddProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-0dor4", "Errors.Project.MemberInvalid")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if existing.ContainsGrantMember(member) {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-8die3", "Errors.Project.MemberAlreadyExists")
	}
	repoProject := model.ProjectFromModel(existing)
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
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dkw35", "Errors.Project.MemberInvalid")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existing.ContainsGrantMember(member) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-8dj4s", "Errors.Project.MemberNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
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
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}
	if !existing.ContainsGrantMember(member) {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-9ode4", "Errors.Project.MemberNotExisting")
	}
	repoProject := model.ProjectFromModel(existing)
	repoMember := model.GrantMemberFromModel(member)

	projectAggregate := ProjectGrantMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return err
}
