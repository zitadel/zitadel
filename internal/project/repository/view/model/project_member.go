package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/lib/pq"
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
	UserID      string         `json:"userId" gorm:"column:user_id;primary_key"`
	ProjectID   string         `json:"-" gorm:"column:project_id;primary_key"`
	UserName    string         `json:"-" gorm:"column:user_name"`
	Email       string         `json:"-" gorm:"column:email_address"`
	FirstName   string         `json:"-" gorm:"column:first_name"`
	LastName    string         `json:"-" gorm:"column:last_name"`
	DisplayName string         `json:"-" gorm:"column:display_name"`
	Roles       pq.StringArray `json:"roles" gorm:"column:roles"`
	Sequence    uint64         `json:"-" gorm:"column:sequence"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
}

func ProjectMemberToModel(member *ProjectMemberView) *model.ProjectMemberView {
	return &model.ProjectMemberView{
		UserID:       member.UserID,
		ProjectID:    member.ProjectID,
		UserName:     member.UserName,
		Email:        member.Email,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		DisplayName:  member.DisplayName,
		Roles:        member.Roles,
		Sequence:     member.Sequence,
		CreationDate: member.CreationDate,
		ChangeDate:   member.ChangeDate,
	}
}

func ProjectMembersToModel(roles []*ProjectMemberView) []*model.ProjectMemberView {
	result := make([]*model.ProjectMemberView, len(roles))
	for i, r := range roles {
		result[i] = ProjectMemberToModel(r)
	}
	return result
}

func (r *ProjectMemberView) AppendEvent(event *models.Event) (err error) {
	r.Sequence = event.Sequence
	r.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.ProjectMemberAdded:
		r.setRootData(event)
		r.CreationDate = event.CreationDate
		err = r.SetData(event)
	case es_model.ProjectMemberChanged:
		err = r.SetData(event)
	}
	return err
}

func (r *ProjectMemberView) setRootData(event *models.Event) {
	r.ProjectID = event.AggregateID
}

func (r *ProjectMemberView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-slo9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
	}
	return nil
}
