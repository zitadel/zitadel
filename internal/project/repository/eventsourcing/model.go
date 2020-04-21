package eventsourcing

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"reflect"
)

const (
	projectVersion = "v1"
)

type Project struct {
	es_models.ObjectRoot
	Name         string           `json:"name,omitempty"`
	State        int32            `json:"-"`
	Members      []*ProjectMember `json:"-"`
	Roles        []*ProjectRole   `json:"-"`
	Applications []*Application   `json:"-"`
	Grants       []*ProjectGrant  `json:"-"`
}

type ProjectMember struct {
	es_models.ObjectRoot
	UserID string   `json:"userId,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

type ProjectRole struct {
	es_models.ObjectRoot
	Key         string `json:"key,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Group       string `json:"group,omitempty"`
}

type Application struct {
	es_models.ObjectRoot
	AppID      string      `json:"appId"`
	State      int32       `json:"-"`
	Name       string      `json:"name,omitempty"`
	Type       int32       `json:"appType,omitempty"`
	OIDCConfig *OIDCConfig `json:"-"`
}

type ApplicationID struct {
	es_models.ObjectRoot
	AppID string `json:"appId"`
}

type OIDCConfig struct {
	es_models.ObjectRoot
	AppID                  string              `json:"appId"`
	ClientID               string              `json:"clientId,omitempty"`
	ClientSecret           *crypto.CryptoValue `json:"clientSecret,omitempty"`
	RedirectUris           []string            `json:"redirectUris,omitempty"`
	ResponseTypes          []int32             `json:"responseTypes,omitempty"`
	GrantTypes             []int32             `json:"grantTypes,omitempty"`
	ApplicationType        int32               `json:"applicationType,omitempty"`
	AuthMethodType         int32               `json:"authMethodType,omitempty"`
	PostLogoutRedirectUris []string            `json:"postLogoutRedirectUris,omitempty"`
}

type ProjectGrant struct {
	es_models.ObjectRoot
	State        int32                 `json:"-"`
	GrantID      string                `json:"grantId,omitempty"`
	GrantedOrgID string                `json:"grantedOrgId,omitempty"`
	RoleKeys     []string              `json:"roleKeys,omitempty"`
	Members      []*ProjectGrantMember `json:"-"`
}

type ProjectGrantID struct {
	es_models.ObjectRoot
	GrantID string `json:"grantId"`
}

type ProjectGrantMember struct {
	es_models.ObjectRoot
	GrantID string   `json:"grantId,omitempty"`
	UserID  string   `json:"userId,omitempty"`
	Roles   []string `json:"roles,omitempty"`
}

func (p *Project) Changes(changed *Project) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func (a *Application) Changes(changed *Application) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["appId"] = a.AppID
	if changed.Name != "" && a.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func (c *OIDCConfig) Changes(changed *OIDCConfig) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["appId"] = c.AppID
	if !reflect.DeepEqual(c.RedirectUris, changed.RedirectUris) {
		changes["redirectUris"] = changed.RedirectUris
	}
	if !reflect.DeepEqual(c.ResponseTypes, changed.ResponseTypes) {
		changes["responseTypes"] = changed.ResponseTypes
	}
	if !reflect.DeepEqual(c.GrantTypes, changed.GrantTypes) {
		changes["grantTypes"] = changed.GrantTypes
	}
	if c.ApplicationType != changed.ApplicationType {
		changes["applicationType"] = changed.ApplicationType
	}
	if c.AuthMethodType != changed.AuthMethodType {
		changes["authMethodType"] = changed.AuthMethodType
	}
	if !reflect.DeepEqual(c.PostLogoutRedirectUris, changed.PostLogoutRedirectUris) {
		changes["postLogoutRedirectUris"] = changed.PostLogoutRedirectUris
	}
	return changes
}

func (g *ProjectGrant) Changes(changed *ProjectGrant) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["grantId"] = g.GrantID
	if !reflect.DeepEqual(g.RoleKeys, changed.RoleKeys) {
		changes["roleKeys"] = changed.RoleKeys
	}
	return changes
}

func ProjectFromModel(project *model.Project) *Project {
	members := ProjectMembersFromModel(project.Members)
	roles := ProjectRolesFromModel(project.Roles)
	apps := AppsFromModel(project.Applications)
	grants := GrantsFromModel(project.Grants)
	return &Project{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  project.ObjectRoot.AggregateID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		Name:         project.Name,
		State:        model.ProjectStateToInt(project.State),
		Members:      members,
		Roles:        roles,
		Applications: apps,
		Grants:       grants,
	}
}

func ProjectToModel(project *Project) *model.Project {
	members := ProjectMembersToModel(project.Members)
	roles := ProjectRolesToModel(project.Roles)
	apps := AppsToModel(project.Applications)
	grants := GrantsToModel(project.Grants)
	return &model.Project{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  project.AggregateID,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
			Sequence:     project.Sequence,
		},
		Name:         project.Name,
		State:        model.ProjectStateFromInt(project.State),
		Members:      members,
		Roles:        roles,
		Applications: apps,
		Grants:       grants,
	}
}

func ProjectMembersToModel(members []*ProjectMember) []*model.ProjectMember {
	convertedMembers := make([]*model.ProjectMember, len(members))
	for i, m := range members {
		convertedMembers[i] = ProjectMemberToModel(m)
	}
	return convertedMembers
}

func ProjectMembersFromModel(members []*model.ProjectMember) []*ProjectMember {
	convertedMembers := make([]*ProjectMember, len(members))
	for i, m := range members {
		convertedMembers[i] = ProjectMemberFromModel(m)
	}
	return convertedMembers
}

func ProjectMemberFromModel(member *model.ProjectMember) *ProjectMember {
	return &ProjectMember{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  member.ObjectRoot.AggregateID,
			Sequence:     member.Sequence,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
		},
		UserID: member.UserID,
		Roles:  member.Roles,
	}
}

func ProjectMemberToModel(member *ProjectMember) *model.ProjectMember {
	return &model.ProjectMember{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  member.AggregateID,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
			Sequence:     member.Sequence,
		},
		UserID: member.UserID,
		Roles:  member.Roles,
	}
}

func ProjectRolesToModel(roles []*ProjectRole) []*model.ProjectRole {
	convertedRoles := make([]*model.ProjectRole, len(roles))
	for i, r := range roles {
		convertedRoles[i] = ProjectRoleToModel(r)
	}
	return convertedRoles
}

func ProjectRolesFromModel(roles []*model.ProjectRole) []*ProjectRole {
	convertedRoles := make([]*ProjectRole, len(roles))
	for i, r := range roles {
		convertedRoles[i] = ProjectRoleFromModel(r)
	}
	return convertedRoles
}

func ProjectRoleFromModel(role *model.ProjectRole) *ProjectRole {
	return &ProjectRole{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  role.ObjectRoot.AggregateID,
			Sequence:     role.Sequence,
			ChangeDate:   role.ChangeDate,
			CreationDate: role.CreationDate,
		},
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}

func ProjectRoleToModel(role *ProjectRole) *model.ProjectRole {
	return &model.ProjectRole{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  role.AggregateID,
			ChangeDate:   role.ChangeDate,
			CreationDate: role.CreationDate,
			Sequence:     role.Sequence,
		},
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}

func AppsToModel(apps []*Application) []*model.Application {
	convertedApps := make([]*model.Application, len(apps))
	for i, a := range apps {
		convertedApps[i] = AppToModel(a)
	}
	return convertedApps
}

func AppsFromModel(apps []*model.Application) []*Application {
	convertedApps := make([]*Application, len(apps))
	for i, a := range apps {
		convertedApps[i] = AppFromModel(a)
	}
	return convertedApps
}

func AppFromModel(app *model.Application) *Application {
	converted := &Application{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  app.ObjectRoot.AggregateID,
			Sequence:     app.Sequence,
			ChangeDate:   app.ChangeDate,
			CreationDate: app.CreationDate,
		},
		AppID: app.AppID,
		Name:  app.Name,
		State: int32(app.State),
		Type:  int32(app.Type),
	}
	if app.OIDCConfig != nil {
		converted.OIDCConfig = OIDCConfigFromModel(app.OIDCConfig)
	}
	return converted
}

func AppToModel(app *Application) *model.Application {
	converted := &model.Application{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  app.AggregateID,
			ChangeDate:   app.ChangeDate,
			CreationDate: app.CreationDate,
			Sequence:     app.Sequence,
		},
		AppID: app.AppID,
		Name:  app.Name,
		State: model.AppState(app.State),
		Type:  model.AppType(app.Type),
	}
	if app.OIDCConfig != nil {
		converted.OIDCConfig = OIDCConfigToModel(app.OIDCConfig)
	}
	return converted
}

func OIDCConfigFromModel(config *model.OIDCConfig) *OIDCConfig {
	responseTypes := make([]int32, len(config.ResponseTypes))
	for i, rt := range config.ResponseTypes {
		responseTypes[i] = int32(rt)
	}
	grantTypes := make([]int32, len(config.GrantTypes))
	for i, rt := range config.GrantTypes {
		grantTypes[i] = int32(rt)
	}
	return &OIDCConfig{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  config.ObjectRoot.AggregateID,
			Sequence:     config.Sequence,
			ChangeDate:   config.ChangeDate,
			CreationDate: config.CreationDate,
		},
		AppID:                  config.AppID,
		ClientID:               config.ClientID,
		ClientSecret:           config.ClientSecret,
		RedirectUris:           config.RedirectUris,
		ResponseTypes:          responseTypes,
		GrantTypes:             grantTypes,
		ApplicationType:        int32(config.ApplicationType),
		AuthMethodType:         int32(config.AuthMethodType),
		PostLogoutRedirectUris: config.PostLogoutRedirectUris,
	}
}

func OIDCConfigToModel(config *OIDCConfig) *model.OIDCConfig {
	responseTypes := make([]model.OIDCResponseType, len(config.ResponseTypes))
	for i, rt := range config.ResponseTypes {
		responseTypes[i] = model.OIDCResponseType(rt)
	}
	grantTypes := make([]model.OIDCGrantType, len(config.GrantTypes))
	for i, rt := range config.GrantTypes {
		grantTypes[i] = model.OIDCGrantType(rt)
	}
	return &model.OIDCConfig{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  config.ObjectRoot.AggregateID,
			Sequence:     config.Sequence,
			ChangeDate:   config.ChangeDate,
			CreationDate: config.CreationDate,
		},
		AppID:                  config.AppID,
		ClientID:               config.ClientID,
		ClientSecret:           config.ClientSecret,
		RedirectUris:           config.RedirectUris,
		ResponseTypes:          responseTypes,
		GrantTypes:             grantTypes,
		ApplicationType:        model.OIDCApplicationType(config.ApplicationType),
		AuthMethodType:         model.OIDCAuthMethodType(config.AuthMethodType),
		PostLogoutRedirectUris: config.PostLogoutRedirectUris,
	}
}

func GrantsToModel(grants []*ProjectGrant) []*model.ProjectGrant {
	convertedGrants := make([]*model.ProjectGrant, len(grants))
	for i, g := range grants {
		convertedGrants[i] = GrantToModel(g)
	}
	return convertedGrants
}

func GrantsFromModel(grants []*model.ProjectGrant) []*ProjectGrant {
	convertedGrants := make([]*ProjectGrant, len(grants))
	for i, g := range grants {
		convertedGrants[i] = GrantFromModel(g)
	}
	return convertedGrants
}

func GrantFromModel(grant *model.ProjectGrant) *ProjectGrant {
	members := GrantMembersFromModel(grant.Members)
	return &ProjectGrant{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  grant.ObjectRoot.AggregateID,
			Sequence:     grant.Sequence,
			ChangeDate:   grant.ChangeDate,
			CreationDate: grant.CreationDate,
		},
		GrantID:      grant.GrantID,
		GrantedOrgID: grant.GrantedOrgID,
		State:        int32(grant.State),
		RoleKeys:     grant.RoleKeys,
		Members:      members,
	}
}

func GrantToModel(grant *ProjectGrant) *model.ProjectGrant {
	members := GrantMembersToModel(grant.Members)
	return &model.ProjectGrant{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  grant.AggregateID,
			ChangeDate:   grant.ChangeDate,
			CreationDate: grant.CreationDate,
			Sequence:     grant.Sequence,
		},
		GrantID:      grant.GrantID,
		GrantedOrgID: grant.GrantedOrgID,
		State:        model.ProjectGrantState(grant.State),
		RoleKeys:     grant.RoleKeys,
		Members:      members,
	}
}

func GrantMembersToModel(members []*ProjectGrantMember) []*model.ProjectGrantMember {
	convertedMembers := make([]*model.ProjectGrantMember, len(members))
	for i, g := range members {
		convertedMembers[i] = GrantMemberToModel(g)
	}
	return convertedMembers
}

func GrantMembersFromModel(members []*model.ProjectGrantMember) []*ProjectGrantMember {
	convertedMembers := make([]*ProjectGrantMember, len(members))
	for i, g := range members {
		convertedMembers[i] = GrantMemberFromModel(g)
	}
	return convertedMembers
}

func GrantMemberFromModel(member *model.ProjectGrantMember) *ProjectGrantMember {
	return &ProjectGrantMember{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  member.ObjectRoot.AggregateID,
			Sequence:     member.Sequence,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
		},
		GrantID: member.GrantID,
		UserID:  member.UserID,
		Roles:   member.Roles,
	}
}

func GrantMemberToModel(member *ProjectGrantMember) *model.ProjectGrantMember {
	return &model.ProjectGrantMember{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  member.AggregateID,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
			Sequence:     member.Sequence,
		},
		GrantID: member.GrantID,
		UserID:  member.UserID,
		Roles:   member.Roles,
	}
}
func ProjectFromEvents(project *Project, events ...*es_models.Event) (*Project, error) {
	if project == nil {
		project = &Project{}
	}

	return project, project.AppendEvents(events...)
}

func (p *Project) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) AppendEvent(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case model.ProjectAdded, model.ProjectChanged:
		if err := json.Unmarshal(event.Data, p); err != nil {
			logging.Log("EVEN-idl93").WithError(err).Error("could not unmarshal event data")
			return err
		}
		p.State = int32(model.PROJECTSTATE_ACTIVE)
		return nil
	case model.ProjectDeactivated:
		return p.appendDeactivatedEvent()
	case model.ProjectReactivated:
		return p.appendReactivatedEvent()
	case model.ProjectMemberAdded:
		return p.appendAddMemberEvent(event)
	case model.ProjectMemberChanged:
		return p.appendChangeMemberEvent(event)
	case model.ProjectMemberRemoved:
		return p.appendRemoveMemberEvent(event)
	case model.ProjectRoleAdded:
		return p.appendAddRoleEvent(event)
	case model.ProjectRoleChanged:
		return p.appendChangeRoleEvent(event)
	case model.ProjectRoleRemoved:
		return p.appendRemoveRoleEvent(event)
	case model.ApplicationAdded:
		return p.appendAddAppEvent(event)
	case model.ApplicationChanged:
		return p.appendChangeAppEvent(event)
	case model.ApplicationRemoved:
		return p.appendRemoveAppEvent(event)
	case model.ApplicationDeactivated:
		return p.appendAppStateEvent(event, model.APPSTATE_INACTIVE)
	case model.ApplicationReactivated:
		return p.appendAppStateEvent(event, model.APPSTATE_ACTIVE)
	case model.OIDCConfigAdded:
		return p.appendAddOIDCConfigEvent(event)
	case model.OIDCConfigChanged, model.OIDCConfigSecretChanged:
		return p.appendChangeOIDCConfigEvent(event)
	case model.ProjectGrantAdded:
		return p.appendAddGrantEvent(event)
	case model.ProjectGrantChanged:
		return p.appendChangeGrantEvent(event)
	case model.ProjectGrantDeactivated:
		return p.appendGrantStateEvent(event, model.PROJECTGRANTSTATE_INACTIVE)
	case model.ProjectGrantReactivated:
		return p.appendGrantStateEvent(event, model.PROJECTGRANTSTATE_ACTIVE)
	case model.ProjectGrantRemoved:
		return p.appendRemoveGrantEvent(event)
	case model.ProjectGrantMemberAdded:
		return p.appendAddGrantMemberEvent(event)
	case model.ProjectGrantMemberChanged:
		return p.appendChangeGrantMemberEvent(event)
	case model.ProjectGrantMemberRemoved:
		return p.appendRemoveGrantMemberEvent(event)
	}
	return nil
}

func (p *Project) appendDeactivatedEvent() error {
	p.State = int32(model.PROJECTSTATE_INACTIVE)
	return nil
}

func (p *Project) appendReactivatedEvent() error {
	p.State = int32(model.PROJECTSTATE_ACTIVE)
	return nil
}

func (p *Project) appendAddMemberEvent(event *es_models.Event) error {
	member := &ProjectMember{}
	err := member.setData(event)
	if err != nil {
		return err
	}
	member.ObjectRoot.CreationDate = event.CreationDate
	p.Members = append(p.Members, member)
	return nil
}

func (p *Project) appendChangeMemberEvent(event *es_models.Event) error {
	member := &ProjectMember{}
	err := member.setData(event)
	if err != nil {
		return err
	}
	for i, m := range p.Members {
		if m.UserID == member.UserID {
			p.Members[i] = member
		}
	}
	return nil
}

func (p *Project) appendRemoveMemberEvent(event *es_models.Event) error {
	member := &ProjectMember{}
	err := member.setData(event)
	if err != nil {
		return err
	}
	for i, m := range p.Members {
		if m.UserID == member.UserID {
			p.Members[i] = p.Members[len(p.Members)-1]
			p.Members[len(p.Members)-1] = nil
			p.Members = p.Members[:len(p.Members)-1]
		}
	}
	return nil
}

func (m *ProjectMember) setData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-e4dkp").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *Project) appendAddRoleEvent(event *es_models.Event) error {
	role := new(ProjectRole)
	err := role.setData(event)
	if err != nil {
		return err
	}
	role.ObjectRoot.CreationDate = event.CreationDate
	p.Roles = append(p.Roles, role)
	return nil
}

func (p *Project) appendChangeRoleEvent(event *es_models.Event) error {
	role := new(ProjectRole)
	err := role.setData(event)
	if err != nil {
		return err
	}
	for i, r := range p.Roles {
		if r.Key == role.Key {
			p.Roles[i] = role
		}
	}
	return nil
}

func (p *Project) appendRemoveRoleEvent(event *es_models.Event) error {
	role := new(ProjectRole)
	err := role.setData(event)
	if err != nil {
		return err
	}
	for i, r := range p.Roles {
		if r.Key == role.Key {
			p.Roles[i] = p.Roles[len(p.Roles)-1]
			p.Roles[len(p.Roles)-1] = nil
			p.Roles = p.Roles[:len(p.Roles)-1]
		}
	}
	return nil
}

func (r *ProjectRole) setData(event *es_models.Event) error {
	r.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-d9euw").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *Project) appendAddAppEvent(event *es_models.Event) error {
	app := new(Application)
	err := app.setData(event)
	if err != nil {
		return err
	}
	app.ObjectRoot.CreationDate = event.CreationDate
	p.Applications = append(p.Applications, app)
	return nil
}

func (p *Project) appendChangeAppEvent(event *es_models.Event) error {
	app := new(Application)
	err := app.setData(event)
	if err != nil {
		return err
	}
	for i, a := range p.Applications {
		if a.AppID == app.AppID {
			p.Applications[i].setData(event)
		}
	}
	return nil
}

func (p *Project) appendRemoveAppEvent(event *es_models.Event) error {
	app := new(Application)
	err := app.setData(event)
	if err != nil {
		return err
	}
	for i, a := range p.Applications {
		if a.AppID == app.AppID {
			p.Applications[i] = p.Applications[len(p.Applications)-1]
			p.Applications[len(p.Applications)-1] = nil
			p.Applications = p.Applications[:len(p.Applications)-1]
		}
	}
	return nil
}

func (p *Project) appendAppStateEvent(event *es_models.Event, state model.AppState) error {
	app := new(Application)
	err := app.setData(event)
	if err != nil {
		return err
	}
	for i, a := range p.Applications {
		if a.AppID == app.AppID {
			a.State = int32(state)
			p.Applications[i] = a
		}
	}
	return nil
}

func (a *Application) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-8die3").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *Project) appendAddOIDCConfigEvent(event *es_models.Event) error {
	config := new(OIDCConfig)
	err := config.setData(event)
	if err != nil {
		return err
	}
	config.ObjectRoot.CreationDate = event.CreationDate
	for i, a := range p.Applications {
		if a.AppID == config.AppID {
			p.Applications[i].Type = int32(model.APPTYPE_OIDC)
			p.Applications[i].OIDCConfig = config
		}
	}
	return nil
}

func (p *Project) appendChangeOIDCConfigEvent(event *es_models.Event) error {
	config := new(OIDCConfig)
	err := config.setData(event)
	if err != nil {
		return err
	}
	for i, a := range p.Applications {
		if a.AppID == config.AppID {
			p.Applications[i].OIDCConfig.setData(event)
		}
	}
	return nil
}

func (o *OIDCConfig) setData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-d8e3s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *Project) appendAddGrantEvent(event *es_models.Event) error {
	grant := new(ProjectGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	grant.ObjectRoot.CreationDate = event.CreationDate
	p.Grants = append(p.Grants, grant)
	return nil
}

func (p *Project) appendChangeGrantEvent(event *es_models.Event) error {
	grant := new(ProjectGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	for i, g := range p.Grants {
		if g.GrantID == grant.GrantID {
			p.Grants[i].getData(event)
		}
	}
	return nil
}

func (p *Project) appendGrantStateEvent(event *es_models.Event, state model.ProjectGrantState) error {
	grant := new(ProjectGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	for i, g := range p.Grants {
		if g.GrantID == grant.GrantID {
			g.State = int32(state)
			p.Grants[i] = g
		}
	}
	return nil
}

func (p *Project) appendRemoveGrantEvent(event *es_models.Event) error {
	grant := new(ProjectGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	for i, g := range p.Grants {
		if g.GrantID == grant.GrantID {
			p.Grants[i] = p.Grants[len(p.Grants)-1]
			p.Grants[len(p.Grants)-1] = nil
			p.Grants = p.Grants[:len(p.Grants)-1]
		}
	}
	return nil
}

func (g *ProjectGrant) getData(event *es_models.Event) error {
	g.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, g); err != nil {
		logging.Log("EVEN-4h6gd").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *Project) appendAddGrantMemberEvent(event *es_models.Event) error {
	member := &ProjectGrantMember{}
	err := member.getData(event)
	if err != nil {
		return err
	}
	member.ObjectRoot.CreationDate = event.CreationDate
	for _, g := range p.Grants {
		if g.GrantID == member.GrantID {
			g.Members = append(g.Members, member)
		}
	}
	return nil
}

func (p *Project) appendChangeGrantMemberEvent(event *es_models.Event) error {
	member := &ProjectGrantMember{}
	err := member.getData(event)
	if err != nil {
		return err
	}
	for _, g := range p.Grants {
		if g.GrantID == member.GrantID {
			for i, m := range g.Members {
				if m.UserID == member.UserID {
					g.Members[i].getData(event)
				}
			}
		}
	}
	return nil
}

func (p *Project) appendRemoveGrantMemberEvent(event *es_models.Event) error {
	member := &ProjectGrantMember{}
	err := member.getData(event)
	if err != nil {
		return err
	}
	for _, g := range p.Grants {
		if g.GrantID == member.GrantID {
			for i, m := range g.Members {
				if m.UserID == member.UserID {
					g.Members[i] = g.Members[len(g.Members)-1]
					g.Members[len(g.Members)-1] = nil
					g.Members = g.Members[:len(g.Members)-1]
				}
			}
		}
	}
	return nil
}

func (m *ProjectGrantMember) getData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-8die2").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
