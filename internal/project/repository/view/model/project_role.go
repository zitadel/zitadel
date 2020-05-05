package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"time"
)

const (
	ProjectRoleKeyNameKey           = "role_key"
	ProjectRoleKeyNameOrgID         = "org_id"
	ProjectRoleKeyNameProjectID     = "project_id"
	ProjectRoleKeyNameResourceOwner = "resource_owner"
)

type ProjectRoleView struct {
	OrgID       string `json:"-" gorm:"column:org_id;primary_key"`
	ProjectID   string `json:"projectId,omitempty" gorm:"column:project_id;primary_key"`
	Key         string `json:"key" gorm:"column:role_key;primary_key"`
	DisplayName string `json:"displayName" gorm:"column:display_name"`
	Group       string `json:"group" gorm:"column:group_name"`
	Sequence    uint64 `json:"-" gorm:"column:sequence"`

	ResourceOwner string    `json:"-" gorm:"resource_owner"`
	CreationDate  time.Time `json:"-" gorm:"column:creation_date"`
}

func ProjectRoleViewFromModel(role *model.ProjectRoleView) *ProjectRoleView {
	return &ProjectRoleView{
		ResourceOwner: role.ResourceOwner,
		OrgID:         role.OrgID,
		ProjectID:     role.ProjectID,
		Key:           role.Key,
		DisplayName:   role.DisplayName,
		Group:         role.Group,
		Sequence:      role.Sequence,
		CreationDate:  role.CreationDate,
	}
}

func ProjectRoleToModel(role *ProjectRoleView) *model.ProjectRoleView {
	return &model.ProjectRoleView{
		ResourceOwner: role.ResourceOwner,
		OrgID:         role.OrgID,
		ProjectID:     role.ProjectID,
		Key:           role.Key,
		DisplayName:   role.DisplayName,
		Group:         role.Group,
		Sequence:      role.Sequence,
		CreationDate:  role.CreationDate,
	}
}

func ProjectRolesToModel(roles []*ProjectRoleView) []*model.ProjectRoleView {
	result := make([]*model.ProjectRoleView, 0)
	for _, r := range roles {
		result = append(result, ProjectRoleToModel(r))
	}
	return result
}

func (r *ProjectRoleView) AppendEvent(event *models.Event) error {
	r.Sequence = event.Sequence
	switch event.Type {
	case es_model.ProjectRoleAdded:
		r.setRootData(event)
		r.SetData(event)
		r.CreationDate = event.CreationDate
	case es_model.ProjectRoleChanged:
		r.SetData(event)
	}
	return nil
}

func (r *ProjectRoleView) setRootData(event *models.Event) {
	r.ProjectID = event.AggregateID
	r.OrgID = event.ResourceOwner
	r.ResourceOwner = event.ResourceOwner
}

func (r *ProjectRoleView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-slo9s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
