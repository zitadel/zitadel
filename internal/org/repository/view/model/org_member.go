package model

import (
	"encoding/json"
	"time"

	es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/org/model"
	"github.com/lib/pq"
)

const (
	OrgMemberKeyUserID    = "user_id"
	OrgMemberKeyOrgID     = "org_id"
	OrgMemberKeyUserName  = "user_name"
	OrgMemberKeyEmail     = "email"
	OrgMemberKeyFirstName = "first_name"
	OrgMemberKeyLastName  = "last_name"
)

type OrgMemberView struct {
	UserID             string         `json:"userId" gorm:"column:user_id;primary_key"`
	OrgID              string         `json:"-" gorm:"column:org_id;primary_key"`
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

func OrgMemberToModel(member *OrgMemberView) *model.OrgMemberView {
	return &model.OrgMemberView{
		UserID:             member.UserID,
		OrgID:              member.OrgID,
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

func OrgMembersToModel(roles []*OrgMemberView) []*model.OrgMemberView {
	result := make([]*model.OrgMemberView, len(roles))
	for i, r := range roles {
		result[i] = OrgMemberToModel(r)
	}
	return result
}

func (r *OrgMemberView) AppendEvent(event *models.Event) (err error) {
	r.Sequence = event.Sequence
	r.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.OrgMemberAdded:
		r.setRootData(event)
		r.CreationDate = event.CreationDate
		err = r.SetData(event)
	case es_model.OrgMemberChanged:
		err = r.SetData(event)
	}
	return err
}

func (r *OrgMemberView) setRootData(event *models.Event) {
	r.OrgID = event.AggregateID
}

func (r *OrgMemberView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-slo9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
	}
	return nil
}
