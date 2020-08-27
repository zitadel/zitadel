package model

import (
	"encoding/json"
	"time"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/lib/pq"
)

const (
	IamMemberKeyUserID    = "user_id"
	IamMemberKeyIamID     = "iam_id"
	IamMemberKeyUserName  = "user_name"
	IamMemberKeyEmail     = "email"
	IamMemberKeyFirstName = "first_name"
	IamMemberKeyLastName  = "last_name"
)

type IamMemberView struct {
	UserID      string         `json:"userId" gorm:"column:user_id;primary_key"`
	IamID       string         `json:"-" gorm:"column:iam_id"`
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

func IamMemberToModel(member *IamMemberView) *model.IamMemberView {
	return &model.IamMemberView{
		UserID:       member.UserID,
		IamID:        member.IamID,
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

func IamMembersToModel(roles []*IamMemberView) []*model.IamMemberView {
	result := make([]*model.IamMemberView, len(roles))
	for i, r := range roles {
		result[i] = IamMemberToModel(r)
	}
	return result
}

func (r *IamMemberView) AppendEvent(event *models.Event) (err error) {
	r.Sequence = event.Sequence
	r.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.IamMemberAdded:
		r.setRootData(event)
		r.CreationDate = event.CreationDate
		err = r.SetData(event)
	case es_model.IamMemberChanged:
		err = r.SetData(event)
	}
	return err
}

func (r *IamMemberView) setRootData(event *models.Event) {
	r.IamID = event.AggregateID
}

func (r *IamMemberView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-Psl89").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
	}
	return nil
}
