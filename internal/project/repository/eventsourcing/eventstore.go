package eventsourcing

import (
	"context"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"strconv"

	"github.com/sony/sonyflake"

	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	proj_model "github.com/caos/zitadel/internal/project/model"
)

type ProjectEventstore struct {
	es_int.Eventstore
	projectCache *ProjectCache
	pwGenerator  crypto.Generator
	idGenerator  *sonyflake.Sonyflake
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
	idGenerator := sonyflake.NewSonyflake(sonyflake.Settings{})
	return &ProjectEventstore{
		Eventstore:   conf.Eventstore,
		projectCache: projectCache,
		pwGenerator:  pwGenerator,
		idGenerator:  idGenerator,
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

func (es *ProjectEventstore) CreateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	if !project.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Name is required")
	}
	id, err := es.idGenerator.NextID()
	if err != nil {
		return nil, err
	}
	project.AggregateID = strconv.FormatUint(id, 10)
	project.State = proj_model.PROJECTSTATE_ACTIVE
	repoProject := model.ProjectFromModel(project)

	createAggregate := ProjectCreateAggregate(es.AggregateCreator(), repoProject)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.projectCache.cacheProject(repoProject)
	return model.ProjectToModel(repoProject), nil
}

func (es *ProjectEventstore) UpdateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	if !project.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Name is required")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "project must be active")
	}

	repoExisting := model.ProjectFromModel(existing)
	aggregate := ProjectDeactivateAggregate(es.AggregateCreator(), repoExisting)
	es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)

	es.projectCache.cacheProject(repoExisting)
	return model.ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ReactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	existing, err := es.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "project must be inactive")
	}

	repoExisting := model.ProjectFromModel(existing)
	aggregate := ProjectReactivateAggregate(es.AggregateCreator(), repoExisting)
	es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)

	es.projectCache.cacheProject(repoExisting)
	return model.ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ProjectMemberByIDs(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if member.UserID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ld93d", "userID missing")
	}
	project, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}

	if _, m := project.GetMember(member.UserID); m != nil {
		return m, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-3udjs", "member not found")
}

func (es *ProjectEventstore) AddProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "User is already member of this Project")
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
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) ChangeProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-oe39f", "User is not member of this project")
	}
	repoProject := model.ProjectFromModel(existing)
	repoMember := model.ProjectMemberFromModel(member)

	projectAggregate := ProjectMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)

	if _, m := model.GetProjectMember(repoProject.Members, member.UserID); m != nil {
		return model.ProjectMemberToModel(m), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) RemoveProjectMember(ctx context.Context, member *proj_model.ProjectMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-d43fs", "UserID and Roles are required")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}
	if _, m := existing.GetMember(member.UserID); m == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-swf34", "User is not member of this project")
	}
	repoProject := model.ProjectFromModel(existing)
	repoMember := model.ProjectMemberFromModel(member)

	projectAggregate := ProjectMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	return err
}

func (es *ProjectEventstore) AddProjectRole(ctx context.Context, role *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	if !role.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-idue3", "Key is required")
	}
	existing, err := es.ProjectByID(ctx, role.AggregateID)
	if err != nil {
		return nil, err
	}
	if existing.ContainsRole(role) {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-sk35t", "Project contains role with same key")
	}
	repoProject := model.ProjectFromModel(existing)
	repoRole := model.ProjectRoleFromModel(role)
	projectAggregate := ProjectRoleAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoRole)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}

	es.projectCache.cacheProject(repoProject)

	if _, r := model.GetProjectRole(repoProject.Roles, role.Key); r != nil {
		return model.ProjectRoleToModel(r), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sie83", "Could not find role in list")
}

func (es *ProjectEventstore) ChangeProjectRole(ctx context.Context, role *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	if !role.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9die3", "Key is required")
	}
	existing, err := es.ProjectByID(ctx, role.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existing.ContainsRole(role) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die34", "Role doesn't exist on this project")
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
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sl1or", "Could not find role in list")
}

func (es *ProjectEventstore) RemoveProjectRole(ctx context.Context, role *proj_model.ProjectRole) error {
	if role.Key == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-id823", "Key is required")
	}
	existing, err := es.ProjectByID(ctx, role.AggregateID)
	if err != nil {
		return err
	}
	if !existing.ContainsRole(role) {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-oe823", "Role doesn't exist on project")
	}
	repoProject := model.ProjectFromModel(existing)
	repoRole := model.ProjectRoleFromModel(role)
	projectAggregate := ProjectRoleRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoRole)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}

func (es *ProjectEventstore) ApplicationByIDs(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ld93d", "project oder app AggregateID missing")
	}
	project, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if _, a := project.GetApp(appID); a != nil {
		return a, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-8ei2s", "Could not find app")
}

func (es *ProjectEventstore) AddApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	if app == nil || !app.IsValid(true) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9eidw", "Some required fields are missing")
	}
	existing, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return nil, err
	}
	id, err := es.idGenerator.NextID()
	if err != nil {
		return nil, err
	}
	app.AppID = strconv.FormatUint(id, 10)

	var stringPw string
	var cryptoPw *crypto.CryptoValue
	if app.OIDCConfig != nil {
		app.OIDCConfig.AppID = strconv.FormatUint(id, 10)
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
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		converted := model.AppToModel(a)
		converted.OIDCConfig.ClientSecretString = stringPw
		return converted, nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) ChangeApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	if app == nil || !app.IsValid(false) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dieuw", "some required fields missing")
	}
	existing, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, app := existing.GetApp(app.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die83", "App is not in this project")
	}
	repoProject := model.ProjectFromModel(existing)
	repoApp := model.AppFromModel(app)

	projectAggregate := ApplicationChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoApp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		return model.AppToModel(a), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-dksi8", "Could not find app in list")
}

func (es *ProjectEventstore) RemoveApplication(ctx context.Context, app *proj_model.Application) error {
	if app.AppID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-id832", "AppID is required")
	}
	existing, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return err
	}
	if _, app := existing.GetApp(app.AppID); app == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83s", "Application doesn't exist on project")
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

func (es *ProjectEventstore) DeactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	if appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dlp9e", "appID missing")
	}
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	app := &proj_model.Application{AppID: appID}
	if _, app := existing.GetApp(app.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-slpe9", "App is not in this project")
	}
	repoProject := model.ProjectFromModel(existing)
	repoApp := model.AppFromModel(app)

	projectAggregate := ApplicationDeactivatedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoApp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		return model.AppToModel(a), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sie83", "Could not find app in list")
}

func (es *ProjectEventstore) ReactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	if appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-0odi2", "appID missing")
	}
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	app := &proj_model.Application{AppID: appID}
	if _, app := existing.GetApp(app.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ld92d", "App is not in this project")
	}
	repoProject := model.ProjectFromModel(existing)
	repoApp := model.AppFromModel(app)

	projectAggregate := ApplicationReactivatedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoApp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		return model.AppToModel(a), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sld93", "Could not find app in list")
}

func (es *ProjectEventstore) ChangeOIDCConfig(ctx context.Context, config *proj_model.OIDCConfig) (*proj_model.OIDCConfig, error) {
	if config == nil || !config.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-du834", "invalid oidc config")
	}
	existing, err := es.ProjectByID(ctx, config.AggregateID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existing.GetApp(config.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dkso8", "App is not in this project")
	}
	if app.Type != proj_model.APPTYPE_OIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-98uje", "App is not an oidc application")
	}
	repoProject := model.ProjectFromModel(existing)
	repoConfig := model.OIDCConfigFromModel(config)

	projectAggregate := OIDCConfigChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoConfig)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		return model.OIDCConfigToModel(a.OIDCConfig), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-dk87s", "Could not find app in list")
}

func (es *ProjectEventstore) ChangeOIDCConfigSecret(ctx context.Context, projectID, appID string) (*proj_model.OIDCConfig, error) {
	if appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7ue34", "some required fields missing")
	}
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existing.GetApp(appID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9odi4", "App is not in this project")
	}
	if app.Type != proj_model.APPTYPE_OIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dile4", "App is not an oidc application")
	}
	repoProject := model.ProjectFromModel(existing)

	stringPw, crypto, err := generateNewClientSecret(es.pwGenerator)
	if err != nil {
		return nil, err
	}

	projectAggregate := OIDCConfigSecretChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, appID, crypto)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)

	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		config := model.OIDCConfigToModel(a.OIDCConfig)
		config.ClientSecretString = stringPw
		return config, nil
	}

	return nil, caos_errs.ThrowInternal(nil, "EVENT-dk87s", "Could not find app in list")
}

func (es *ProjectEventstore) ProjectGrantByIDs(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-e8die", "grantID missing")
	}
	project, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if _, g := project.GetGrant(grantID); g != nil {
		return g, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-slo45", "grant not found")
}

func (es *ProjectEventstore) AddProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) (*proj_model.ProjectGrant, error) {
	if grant == nil || !grant.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-37dhs", "Project grant invalid")
	}
	existing, err := es.ProjectByID(ctx, grant.AggregateID)
	if err != nil {
		return nil, err
	}
	if existing.ContainsGrantForOrg(grant.GrantedOrgID) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7ug4g", "Grant for org already exists")
	}
	if !existing.ContainsRoles(grant.RoleKeys) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83d", "One role doesnt exist in Project")
	}
	id, err := es.idGenerator.NextID()
	if err != nil {
		return nil, err
	}
	grant.GrantID = strconv.FormatUint(id, 10)

	repoProject := model.ProjectFromModel(existing)
	repoGrant := model.GrantFromModel(grant)

	addAggregate := ProjectGrantAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, addAggregate)

	if _, g := model.GetProjectGrant(repoProject.Grants, grant.GrantID); g != nil {
		return model.GrantToModel(g), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sk3t5", "Could not find grant in list")
}

func (es *ProjectEventstore) ChangeProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) (*proj_model.ProjectGrant, error) {
	if grant == nil && grant.GrantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-8sie3", "invalid grant")
	}
	existing, err := es.ProjectByID(ctx, grant.AggregateID)
	if err != nil {
		return nil, err
	}
	if _, g := existing.GetGrant(grant.GrantID); g == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die83", "Grant not existing on project")
	}
	if !existing.ContainsRoles(grant.RoleKeys) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-di83d", "One role doesnt exist in Project")
	}
	repoProject := model.ProjectFromModel(existing)
	repoGrant := model.GrantFromModel(grant)

	projectAggregate := ProjectGrantChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)

	if _, g := model.GetProjectGrant(repoProject.Grants, grant.GrantID); g != nil {
		return model.GrantToModel(g), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-dksi8", "Could not find app in list")
}

func (es *ProjectEventstore) RemoveProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) error {
	if grant.GrantID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-8eud6", "GrantId is required")
	}
	existing, err := es.ProjectByID(ctx, grant.AggregateID)
	if err != nil {
		return err
	}
	if _, g := existing.GetGrant(grant.GrantID); g == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-9ie3s", "Grant doesn't exist on project")
	}
	repoProject := model.ProjectFromModel(existing)
	grantRepo := model.GrantFromModel(grant)
	projectAggregate := ProjectGrantRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, grantRepo)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}

func (es *ProjectEventstore) DeactivateProjectGrant(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7due2", "grantID missing")
	}
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	grant := &proj_model.ProjectGrant{GrantID: grantID}
	if _, g := existing.GetGrant(grant.GrantID); g == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-slpe9", "Grant is not in this project")
	}
	repoProject := model.ProjectFromModel(existing)
	repoGrant := model.GrantFromModel(grant)

	projectAggregate := ProjectGrantDeactivatedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	if _, g := model.GetProjectGrant(repoProject.Grants, grant.GrantID); g != nil {
		return model.GrantToModel(g), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sie83", "Could not find grant in list")
}

func (es *ProjectEventstore) ReactivateProjectGrant(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-d7suw", "grantID missing")
	}
	existing, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	grant := &proj_model.ProjectGrant{GrantID: grantID}
	if _, g := existing.GetGrant(grant.GrantID); g == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-0spew", "Grant is not in this project")
	}
	repoProject := model.ProjectFromModel(existing)
	repoGrant := model.GrantFromModel(grant)

	projectAggregate := ProjectGrantReactivatedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)

	if _, g := model.GetProjectGrant(repoProject.Grants, grant.GrantID); g != nil {
		return model.GrantToModel(g), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-9osjw", "Could not find grant in list")
}

func (es *ProjectEventstore) ProjectGrantMemberByIDs(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	if member.GrantID == "" || member.UserID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-8diw2", "userID missing")
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
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-3udjs", "member not found")
}

func (es *ProjectEventstore) AddProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-0dor4", "invalid member")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if existing.ContainsGrantMember(member) {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-8die3", "User is already member of this ProjectGrant")
	}
	repoProject := model.ProjectFromModel(existing)
	repoMember := model.GrantMemberFromModel(member)

	addAggregate := ProjectGrantMemberAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, addAggregate)
	es.projectCache.cacheProject(repoProject)
	if _, g := model.GetProjectGrant(repoProject.Grants, member.GrantID); g != nil {
		if _, m := model.GetProjectGrantMember(g.Members, member.UserID); m != nil {
			return model.GrantMemberToModel(m), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) ChangeProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dkw35", "member is not valid")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existing.ContainsGrantMember(member) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-8dj4s", "User is not member of this grant")
	}
	repoProject := model.ProjectFromModel(existing)
	repoMember := model.GrantMemberFromModel(member)

	projectAggregate := ProjectGrantMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	if _, g := model.GetProjectGrant(repoProject.Grants, member.GrantID); g != nil {
		if _, m := model.GetProjectGrantMember(g.Members, member.UserID); m != nil {
			return model.GrantMemberToModel(m), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-s8ur3", "Could not find member in list")
}

func (es *ProjectEventstore) RemoveProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-8su4r", "member is not valid")
	}
	existing, err := es.ProjectByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}
	if !existing.ContainsGrantMember(member) {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-9ode4", "User is not member of this grant")
	}
	repoProject := model.ProjectFromModel(existing)
	repoMember := model.GrantMemberFromModel(member)

	projectAggregate := ProjectGrantMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	return err
}
