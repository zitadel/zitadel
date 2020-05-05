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
	ProjectGrantMemberKeyUserID    = "user_id"
	ProjectGrantMemberKeyGrantID   = "grant_id"
	ProjectGrantMemberKeyUserName  = "user_name"
	ProjectGrantMemberKeyEmail     = "email"
	ProjectGrantMemberKeyFirstName = "first_name"
	ProjectGrantMemberKeyLastName  = "last_name"
)

type ProjectGrantMemberView struct {
	UserID    string         `json:"userId" gorm:"column:user_id;primary_key"`
	GrantID   string         `json:"grantId" gorm:"column:grant_id;primary_key"`
	ProjectID string         `json:"-" gorm:"column:project_id"`
	UserName  string         `json:"-" gorm:"column:user_name"`
	Email     string         `json:"-" gorm:"column:email_address"`
	FirstName string         `json:"-" gorm:"column:first_name"`
	LastName  string         `json:"-" gorm:"column:last_name"`
	Roles     pq.StringArray `json:"roles" gorm:"column:roles"`
	Sequence  uint64         `json:"-" gorm:"column:sequence"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
}

func ProjectGrantMemberViewFromModel(member *model.ProjectGrantMemberView) *ProjectGrantMemberView {
	return &ProjectGrantMemberView{
		UserID:       member.UserID,
		GrantID:      member.GrantID,
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

func ProjectGrantMemberToModel(member *ProjectGrantMemberView) *model.ProjectGrantMemberView {
	return &model.ProjectGrantMemberView{
		UserID:       member.UserID,
		GrantID:      member.GrantID,
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

func ProjectGrantMembersToModel(roles []*ProjectGrantMemberView) []*model.ProjectGrantMemberView {
	result := make([]*model.ProjectGrantMemberView, 0)
	for _, r := range roles {
		result = append(result, ProjectGrantMemberToModel(r))
	}
	return result
}

func (r *ProjectGrantMemberView) AppendEvent(event *models.Event) error {
	r.Sequence = event.Sequence
	switch event.Type {
	case es_model.ProjectGrantMemberAdded:
		r.setRootData(event)
		r.SetData(event)
		r.CreationDate = event.CreationDate
	case es_model.ProjectGrantMemberChanged:
		r.SetData(event)
	}
	return nil
}

func (r *ProjectGrantMemberView) setRootData(event *models.Event) {
	r.ProjectID = event.AggregateID
	r.ChangeDate = event.CreationDate
}

func (r *ProjectGrantMemberView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-slo9s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
