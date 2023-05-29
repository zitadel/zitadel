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
	ProjectMemberKeyUserID    = "user_id"
	ProjectMemberKeyProjectID = "project_id"
	ProjectMemberKeyUserName  = "user_name"
	ProjectMemberKeyEmail     = "email"
	ProjectMemberKeyFirstName = "first_name"
	ProjectMemberKeyLastName  = "last_name"
)

type ProjectMemberView struct {
	UserID             string                     `json:"userId" gorm:"column:user_id;primary_key"`
	ProjectID          string                     `json:"-" gorm:"column:project_id;primary_key"`
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

func (r *ProjectMemberView) AppendEvent(event *models.Event) (err error) {
	r.Sequence = event.Seq
	r.ChangeDate = event.CreationDate
	switch event.Type() {
	case project.MemberAddedType:
		r.setRootData(event)
		r.CreationDate = event.CreationDate
		err = r.SetData(event)
	case project.MemberChangedType:
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
