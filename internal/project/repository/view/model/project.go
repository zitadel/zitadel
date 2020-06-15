package model

import (
	"encoding/json"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"time"
)

const (
	ProjectKeyProjectID     = "project_id"
	ProjectKeyResourceOwner = "resource_owner"
	ProjectKeyName          = "project_name"
)

type ProjectView struct {
	ProjectID     string    `json:"-" gorm:"column:project_id;primary_key"`
	Name          string    `json:"name" gorm:"column:project_name"`
	CreationDate  time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate    time.Time `json:"-" gorm:"column:change_date"`
	State         int32     `json:"-" gorm:"column:project_state"`
	ResourceOwner string    `json:"-" gorm:"column:resource_owner"`
	Sequence      uint64    `json:"-" gorm:"column:sequence"`
}

func ProjectFromModel(project *model.ProjectView) *ProjectView {
	return &ProjectView{
		ProjectID:     project.ProjectID,
		Name:          project.Name,
		ChangeDate:    project.ChangeDate,
		CreationDate:  project.CreationDate,
		State:         int32(project.State),
		ResourceOwner: project.ResourceOwner,
		Sequence:      project.Sequence,
	}
}

func ProjectToModel(project *ProjectView) *model.ProjectView {
	return &model.ProjectView{
		ProjectID:     project.ProjectID,
		Name:          project.Name,
		ChangeDate:    project.ChangeDate,
		CreationDate:  project.CreationDate,
		State:         model.ProjectState(project.State),
		ResourceOwner: project.ResourceOwner,
		Sequence:      project.Sequence,
	}
}

func ProjectsToModel(projects []*ProjectView) []*model.ProjectView {
	result := make([]*model.ProjectView, len(projects))
	for i, p := range projects {
		result[i] = ProjectToModel(p)
	}
	return result
}

func (p *ProjectView) AppendEvent(event *models.Event) (err error) {
	p.ChangeDate = event.CreationDate
	p.Sequence = event.Sequence
	switch event.Type {
	case es_model.ProjectAdded:
		p.State = int32(model.PROJECTSTATE_ACTIVE)
		p.CreationDate = event.CreationDate
		p.setRootData(event)
		err = p.setData(event)
	case es_model.ProjectChanged:
		err = p.setData(event)
	case es_model.ProjectDeactivated:
		p.State = int32(model.PROJECTSTATE_INACTIVE)
	case es_model.ProjectReactivated:
		p.State = int32(model.PROJECTSTATE_ACTIVE)
	}
	return err
}

func (p *ProjectView) setRootData(event *models.Event) {
	p.ProjectID = event.AggregateID
	p.ResourceOwner = event.ResourceOwner
}

func (p *ProjectView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-dlo92").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *ProjectView) setProjectData(event *models.Event) error {
	project := new(ProjectView)
	err := project.SetData(event)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-sk9Sj").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-s9ols", "Could not unmarshal data")
	}
	return nil
}
