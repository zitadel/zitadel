package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"time"
)

type GrantedProject struct {
	ProjectID     string    `json:"-" gorm:"column:project_id;primary_key"`
	OrgID         string    `json:"-" gorm:"column:org_id;primary_key"`
	Name          string    `json:"name" gorm:"column:project_name"`
	CreationDate  time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate    time.Time `json:"-" gorm:"column:change_date"`
	State         int32     `json:"-" gorm:"column:project_state"`
	Type          int32     `json:"-" gorm:"column:project_type"`
	ResourceOwner string    `json:"-" gorm:"column:resource_owner"`
	OrgName       string    `json:"-" gorm:"column:org_name"`
	OrgDomain     string    `json:"-" gorm:"column:org_domain"`
	Sequence      uint64    `json:"-" gorm:"column:sequence"`
	GrantID       string    `json:"-" gorm:"column:grant_id"`
}

func GrantedProjectFromModel(project *model.GrantedProject) *GrantedProject {
	return &GrantedProject{
		ProjectID:     project.ProjectID,
		OrgID:         project.OrgID,
		Name:          project.Name,
		ChangeDate:    project.ChangeDate,
		CreationDate:  project.CreationDate,
		State:         int32(project.State),
		Type:          int32(project.Type),
		ResourceOwner: project.ResourceOwner,
		OrgName:       project.OrgName,
		GrantID:       project.GrantID,
		Sequence:      project.Sequence,
	}
}

func GrantedProjectToModel(project *GrantedProject) *model.GrantedProject {
	return &model.GrantedProject{
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

func GrantedProjectsToModel(projects []*GrantedProject) []*model.GrantedProject {
	result := make([]*model.GrantedProject, 0)
	for _, p := range projects {
		result = append(result, GrantedProjectToModel(p))
	}
	return result
}

func (p *GrantedProject) AppendEvents(events ...*models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *GrantedProject) AppendEvent(event *models.Event) (err error) {
	return nil
}
