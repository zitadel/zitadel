package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/lib/pq"
	"time"
)

const (
	GrantedProjectKeyProjectID     = "project_id"
	GrantedProjectKeyGrantID       = "grant_id"
	GrantedProjectKeyOrgID         = "org_id"
	GrantedProjectKeyResourceOwner = "resource_owner"
	GrantedProjectKeyName          = "name"
)

type GrantedProjectView struct {
	ProjectID       string         `json:"-" gorm:"column:project_id;primary_key"`
	OrgID           string         `json:"-" gorm:"column:org_id;primary_key"`
	Name            string         `json:"name" gorm:"column:project_name"`
	CreationDate    time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate      time.Time      `json:"-" gorm:"column:change_date"`
	State           int32          `json:"-" gorm:"column:project_state"`
	Type            int32          `json:"-" gorm:"column:project_type"`
	ResourceOwner   string         `json:"-" gorm:"column:resource_owner"`
	OrgName         string         `json:"-" gorm:"column:org_name"`
	OrgDomain       string         `json:"-" gorm:"column:org_domain"`
	Sequence        uint64         `json:"-" gorm:"column:sequence"`
	GrantID         string         `json:"-" gorm:"column:grant_id"`
	GrantedRoleKeys pq.StringArray `json:"-" gorm:"column:granted_role_keys"`
}

type ProjectGrant struct {
	GrantID      string   `json:"grantId"`
	GrantedOrgID string   `json:"grantedOrgId"`
	RoleKeys     []string `json:"roleKeys"`
}

func GrantedProjectFromModel(project *model.GrantedProjectView) *GrantedProjectView {
	return &GrantedProjectView{
		ProjectID:       project.ProjectID,
		OrgID:           project.OrgID,
		Name:            project.Name,
		ChangeDate:      project.ChangeDate,
		CreationDate:    project.CreationDate,
		State:           int32(project.State),
		Type:            int32(project.Type),
		ResourceOwner:   project.ResourceOwner,
		OrgName:         project.OrgName,
		GrantID:         project.GrantID,
		GrantedRoleKeys: project.GrantedRoleKeys,
		Sequence:        project.Sequence,
	}
}

func GrantedProjectToModel(project *GrantedProjectView) *model.GrantedProjectView {
	return &model.GrantedProjectView{
		ProjectID:     project.ProjectID,
		OrgID:         project.OrgID,
		Name:          project.Name,
		ChangeDate:    project.ChangeDate,
		CreationDate:  project.CreationDate,
		State:         model.ProjectState(project.State),
		Type:          model.ProjectType(project.Type),
		ResourceOwner: project.ResourceOwner,
		OrgName:       project.OrgName,
		GrantID:       project.GrantID,
		Sequence:      project.Sequence,
	}
}

func GrantedProjectsToModel(projects []*GrantedProjectView) []*model.GrantedProjectView {
	result := make([]*model.GrantedProjectView, 0)
	for _, p := range projects {
		result = append(result, GrantedProjectToModel(p))
	}
	return result
}

func (p *GrantedProjectView) AppendEvent(event *models.Event) error {
	p.ChangeDate = event.CreationDate
	p.Sequence = event.Sequence
	switch event.Type {
	case es_model.ProjectAdded:
		p.setRootData(event)
		p.setData(event)
		p.State = int32(model.PROJECTSTATE_ACTIVE)
		p.CreationDate = event.CreationDate
	case es_model.ProjectChanged:
		p.setData(event)
	case es_model.ProjectDeactivated:
		p.State = int32(model.PROJECTSTATE_INACTIVE)
	case es_model.ProjectReactivated:
		p.State = int32(model.PROJECTSTATE_ACTIVE)
	case es_model.ProjectGrantAdded:
		p.setRootData(event)
		p.setProjectGrantData(event)
		p.State = int32(model.PROJECTSTATE_ACTIVE)
		p.CreationDate = event.CreationDate
	case es_model.ProjectGrantDeactivated:
		p.State = int32(model.PROJECTSTATE_INACTIVE)
	case es_model.ProjectGrantReactivated:
		p.State = int32(model.PROJECTSTATE_ACTIVE)
	}
	return nil
}

func (p *GrantedProjectView) setRootData(event *models.Event) {
	p.ProjectID = event.AggregateID
	p.OrgID = event.ResourceOwner
	p.ResourceOwner = event.ResourceOwner
}

func (p *GrantedProjectView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-dlo92").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *GrantedProjectView) setProjectGrantData(event *models.Event) error {
	grant := new(ProjectGrant)
	err := grant.SetData(event)
	if err != nil {
		return err
	}
	p.OrgID = grant.GrantedOrgID
	p.GrantID = grant.GrantID
	p.GrantedRoleKeys = grant.RoleKeys
	return nil
}

func (p *ProjectGrant) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-dlo92").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
