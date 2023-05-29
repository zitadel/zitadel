package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/project"
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
	UserID             string                     `json:"userId" gorm:"column:user_id;primary_key"`
	GrantID            string                     `json:"grantId" gorm:"column:grant_id;primary_key"`
	ProjectID          string                     `json:"-" gorm:"column:project_id"`
	UserName           string                     `json:"-" gorm:"column:user_name"`
	Email              string                     `json:"-" gorm:"column:email_address"`
	FirstName          string                     `json:"-" gorm:"column:first_name"`
	LastName           string                     `json:"-" gorm:"column:last_name"`
	DisplayName        string                     `json:"-" gorm:"column:display_name"`
	Roles              database.TextArray[string] `json:"roles" gorm:"column:roles"`
	Sequence           uint64                     `json:"-" gorm:"column:sequence"`
	PreferredLoginName string                     `json:"-" gorm:"column:preferred_login_name"`
	AvatarKey          string                     `json:"-" gorm:"column:avatar_key"`
	UserResourceOwner  string                     `json:"-" gorm:"column:user_resource_owner"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
}

func (r *ProjectGrantMemberView) AppendEvent(event *models.Event) (err error) {
	r.Sequence = event.Seq
	r.ChangeDate = event.CreationDate
	switch event.Type() {
	case project.GrantMemberAddedType:
		r.setRootData(event)
		r.CreationDate = event.CreationDate
		err = r.SetData(event)
	case project.GrantMemberChangedType:
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
