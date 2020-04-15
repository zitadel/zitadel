package eventsourcing

import (
	"encoding/json"
	"github.com/caos/logging"
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
	OIDCConfig *OIDCConfig `json:"-"`
}

type ApplicationID struct {
	es_models.ObjectRoot
	AppID string `json:"appId"`
}

type OIDCConfig struct {
	es_models.ObjectRoot
	AppID                  string   `json:"appId"`
	ClientID               string   `json:"clientId,omitempty"`
	ClientSecret           []byte   `json:"clientSecret,omitempty"`
	RedirectUris           []string `json:"redirectUris,omitempty"`
	ResponseTypes          []int32  `json:"responseTypes,omitempty"`
	GrantTypes             []int32  `json:"grantTypes,omitempty"`
	ApplicationType        int32    `json:"applicationType,omitempty"`
	AuthMethodType         int32    `json:"authMethodType,omitempty"`
	PostLogoutRedirectUris []string `json:"postLogoutRedirectUris,omitempty"`
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
	if changed.Name != "" && a.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func (c *OIDCConfig) Changes(changed *OIDCConfig) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if !reflect.DeepEqual(c.RedirectUris, changed.RedirectUris) {
		changes["redirectUris"] = c.RedirectUris
	}
	if !reflect.DeepEqual(c.ResponseTypes, changed.ResponseTypes) {
		changes["responseTypes"] = c.ResponseTypes
	}
	if !reflect.DeepEqual(c.GrantTypes, changed.GrantTypes) {
		changes["grantTypes"] = c.GrantTypes
	}
	if c.ApplicationType != changed.ApplicationType {
		changes["applicationType"] = c.ApplicationType
	}
	if c.AuthMethodType != changed.AuthMethodType {
		changes["authMethodType"] = c.AuthMethodType
	}
	if !reflect.DeepEqual(c.PostLogoutRedirectUris, changed.PostLogoutRedirectUris) {
		changes["postLogoutRedirectUris"] = c.PostLogoutRedirectUris
	}
	return changes
}

func ProjectFromModel(project *model.Project) *Project {
	members := ProjectMembersFromModel(project.Members)
	roles := ProjectRolesFromModel(project.Roles)
	return &Project{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ObjectRoot.ID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		Name:    project.Name,
		State:   model.ProjectStateToInt(project.State),
		Members: members,
		Roles:   roles,
	}
}

func ProjectToModel(project *Project) *model.Project {
	members := ProjectMembersToModel(project.Members)
	roles := ProjectRolesToModel(project.Roles)
	return &model.Project{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ID,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
			Sequence:     project.Sequence,
		},
		Name:    project.Name,
		State:   model.ProjectStateFromInt(project.State),
		Members: members,
		Roles:   roles,
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
			ID:           member.ObjectRoot.ID,
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
			ID:           member.ID,
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
			ID:           role.ObjectRoot.ID,
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
			ID:           role.ID,
			ChangeDate:   role.ChangeDate,
			CreationDate: role.CreationDate,
			Sequence:     role.Sequence,
		},
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
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
		p.State = model.ProjectStateToInt(model.PROJECTSTATE_ACTIVE)
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
		return p.appendAddAppEvent(event)
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
	}
	return nil
}

func (p *Project) appendDeactivatedEvent() error {
	p.State = model.ProjectStateToInt(model.PROJECTSTATE_INACTIVE)
	return nil
}

func (p *Project) appendReactivatedEvent() error {
	p.State = model.ProjectStateToInt(model.PROJECTSTATE_ACTIVE)
	return nil
}

func (p *Project) appendAddMemberEvent(event *es_models.Event) error {
	member := &ProjectMember{}
	err := member.getData(event)
	if err != nil {
		return err
	}
	member.ObjectRoot.CreationDate = event.CreationDate
	p.Members = append(p.Members, member)
	return nil
}

func (p *Project) appendChangeMemberEvent(event *es_models.Event) error {
	member := &ProjectMember{}
	err := member.getData(event)
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
	err := member.getData(event)
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

func (m *ProjectMember) getData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-e4dkp").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *Project) appendAddRoleEvent(event *es_models.Event) error {
	role := new(ProjectRole)
	err := role.getData(event)
	if err != nil {
		return nil
	}
	role.ObjectRoot.CreationDate = event.CreationDate
	p.Roles = append(p.Roles, role)
	return nil
}

func (p *Project) appendChangeRoleEvent(event *es_models.Event) error {
	role := new(ProjectRole)
	err := role.getData(event)
	if err != nil {
		return nil
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
	err := role.getData(event)
	if err != nil {
		return nil
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

func (r *ProjectRole) getData(event *es_models.Event) error {
	r.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-d9euw").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *Project) appendAddAppEvent(event *es_models.Event) error {
	app := new(Application)
	err := app.getData(event)
	if err != nil {
		return nil
	}
	app.ObjectRoot.CreationDate = event.CreationDate
	p.Applications = append(p.Applications, app)
	return nil
}

func (p *Project) appendChangeAppEvent(event *es_models.Event) error {
	app := new(Application)
	err := app.getData(event)
	if err != nil {
		return nil
	}
	for i, a := range p.Applications {
		if a.AppID == app.AppID {
			p.Applications[i].getData(event)
		}
	}
	return nil
}

func (p *Project) appendRemoveAppEvent(event *es_models.Event) error {
	app := new(Application)
	err := app.getData(event)
	if err != nil {
		return nil
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
	err := app.getData(event)
	if err != nil {
		return nil
	}
	for i, a := range p.Applications {
		if a.AppID == app.AppID {
			a.State = model.AppStateToInt(state)
			p.Applications[i] = a
		}
	}
	return nil
}

func (a *Application) getData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-8die3").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *Project) appendAddOIDCConfigEvent(event *es_models.Event) error {
	config := new(OIDCConfig)
	err := config.getData(event)
	if err != nil {
		return nil
	}
	config.ObjectRoot.CreationDate = event.CreationDate
	for i, a := range p.Applications {
		if a.AppID == config.AppID {
			p.Applications[i].OIDCConfig = config
		}
	}
	return nil
}

func (p *Project) appendChangeOIDCConfigEvent(event *es_models.Event) error {
	config := new(OIDCConfig)
	err := config.getData(event)
	if err != nil {
		return nil
	}
	for i, a := range p.Applications {
		if a.AppID == config.AppID {
			p.Applications[i].OIDCConfig.getData(event)
		}
	}
	return nil
}

func (o *OIDCConfig) getData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-d8e3s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
