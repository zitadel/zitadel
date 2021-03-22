package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/lib/pq"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
)

const (
	ProjectGrantMemberKeyUserID    = "user_id"
	ProjectGrantMemberKeyGrantID   = "grant_id"
	ProjectGrantMemberKeyProjectID = "project_id"
	ProjectGrantMemberKeyUserName  = "user_name"
	ProjectGrantMemberKeyEmail     = "email"
	ProjectGrantMemberKeyFirstName = "first_name"
	ProjectGrantMemberKeyLastName  = "last_name"
)

type ProjectGrantMemberView struct {
	UserID             string         `json:"userId" gorm:"column:user_id;primary_key"`
	GrantID            string         `json:"grantId" gorm:"column:grant_id;primary_key"`
	ProjectID          string         `json:"-" gorm:"column:project_id"`
	UserName           string         `json:"-" gorm:"column:user_name"`
	Email              string         `json:"-" gorm:"column:email_address"`
	FirstName          string         `json:"-" gorm:"column:first_name"`
	LastName           string         `json:"-" gorm:"column:last_name"`
	DisplayName        string         `json:"-" gorm:"column:display_name"`
	Roles              pq.StringArray `json:"roles" gorm:"column:roles"`
	Sequence           uint64         `json:"-" gorm:"column:sequence"`
	PreferredLoginName string         `json:"-" gorm:"column:preferred_login_name"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
}

func ProjectGrantMemberToModel(member *ProjectGrantMemberView) *model.ProjectGrantMemberView {
	return &model.ProjectGrantMemberView{
		UserID:             member.UserID,
		GrantID:            member.GrantID,
		ProjectID:          member.ProjectID,
		UserName:           member.UserName,
		Email:              member.Email,
		FirstName:          member.FirstName,
		LastName:           member.LastName,
		DisplayName:        member.DisplayName,
		PreferredLoginName: member.PreferredLoginName,
		Roles:              member.Roles,
		Sequence:           member.Sequence,
		CreationDate:       member.CreationDate,
		ChangeDate:         member.ChangeDate,
	}
}

func ProjectGrantMembersToModel(roles []*ProjectGrantMemberView) []*model.ProjectGrantMemberView {
	result := make([]*model.ProjectGrantMemberView, len(roles))
	for i, r := range roles {
		result[i] = ProjectGrantMemberToModel(r)
	}
	return result
}

func (r *ProjectGrantMemberView) AppendEvent(event *models.Event) (err error) {
	r.Sequence = event.Sequence
	r.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.ProjectGrantMemberAdded:
		r.setRootData(event)
		r.CreationDate = event.CreationDate
		err = r.SetData(event)
	case es_model.ProjectGrantMemberChanged:
		err = r.SetData(event)
	}
	return err
}

func (r *ProjectGrantMemberView) setRootData(event *models.Event) {
	r.ProjectID = event.AggregateID
}

func (r *ProjectGrantMemberView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-slo9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-0plew", "Could not unmarshal data")
	}
	return nil
}
