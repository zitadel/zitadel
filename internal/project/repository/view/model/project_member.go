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
	ProjectMemberKeyUserID    = "user_id"
	ProjectMemberKeyProjectID = "project_id"
	ProjectMemberKeyUserName  = "user_name"
	ProjectMemberKeyEmail     = "email"
	ProjectMemberKeyFirstName = "first_name"
	ProjectMemberKeyLastName  = "last_name"
)

type ProjectMemberView struct {
	UserID    string         `json:"userId" gorm:"column:user_id;primary_key"`
	ProjectID string         `json:"-" gorm:"column:project_id;primary_key"`
	UserName  string         `json:"userName,omitempty" gorm:"column:user_name"`
	Email     string         `json:"email,omitempty" gorm:"column:email_address"`
	FirstName string         `json:"firstName,omitempty" gorm:"column:first_name"`
	LastName  string         `json:"lastName,omitempty" gorm:"column:last_name"`
	Roles     pq.StringArray `json:"roles,omitempty" gorm:"column:roles"`
	Sequence  uint64         `json:"-" gorm:"column:sequence"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
}

func ProjectMemberViewFromModel(member *model.ProjectMemberView) *ProjectMemberView {
	return &ProjectMemberView{
		UserID:       member.UserID,
		ProjectID:    member.ProjectID,
		UserName:     member.UserName,
		Email:        member.Email,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		Roles:        member.Roles,
		Sequence:     member.Sequence,
		CreationDate: member.CreationDate,
		ChangeDate:   member.ChangeDate,
	}
}

func ProjectMemberToModel(member *ProjectMemberView) *model.ProjectMemberView {
	return &model.ProjectMemberView{
		UserID:       member.UserID,
		ProjectID:    member.ProjectID,
		UserName:     member.UserName,
		Email:        member.Email,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		Roles:        member.Roles,
		Sequence:     member.Sequence,
		CreationDate: member.CreationDate,
		ChangeDate:   member.ChangeDate,
	}
}

func ProjectMembersToModel(roles []*ProjectMemberView) []*model.ProjectMemberView {
	result := make([]*model.ProjectMemberView, 0)
	for _, r := range roles {
		result = append(result, ProjectMemberToModel(r))
	}
	return result
}

func (r *ProjectMemberView) AppendEvent(event *models.Event) error {
	r.Sequence = event.Sequence
	switch event.Type {
	case es_model.ProjectMemberAdded:
		r.setRootData(event)
		r.SetData(event)
		r.CreationDate = event.CreationDate
	case es_model.ProjectMemberChanged:
		r.SetData(event)
	}
	return nil
}

func (r *ProjectMemberView) setRootData(event *models.Event) {
	r.ProjectID = event.AggregateID
	r.ChangeDate = event.CreationDate
}

func (r *ProjectMemberView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-slo9s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
