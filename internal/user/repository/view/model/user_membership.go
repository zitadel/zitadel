package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/lib/pq"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/user/model"
)

const (
	UserMembershipKeyUserID        = "user_id"
	UserMembershipKeyAggregateID   = "aggregate_id"
	UserMembershipKeyObjectID      = "object_id"
	UserMembershipKeyResourceOwner = "resource_owner"
	UserMembershipKeyMemberType    = "member_type"
)

type UserMembershipView struct {
	UserID      string `json:"-" gorm:"column:user_id;primary_key"`
	MemberType  int32  `json:"-" gorm:"column:member_type;primary_key"`
	AggregateID string `json:"-" gorm:"column:aggregate_id;primary_key"`
	ObjectID    string `json:"-" gorm:"column:object_id;primary_key"`

	Roles             pq.StringArray `json:"-" gorm:"column:roles"`
	DisplayName       string         `json:"-" gorm:"column:display_name"`
	CreationDate      time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time      `json:"-" gorm:"column:change_date"`
	ResourceOwner     string         `json:"-" gorm:"column:resource_owner"`
	ResourceOwnerName string         `json:"-" gorm:"column:resource_owner_name"`
	Sequence          uint64         `json:"-" gorm:"column:sequence"`
}

func UserMembershipToModel(membership *UserMembershipView) *model.UserMembershipView {
	return &model.UserMembershipView{
		UserID:            membership.UserID,
		MemberType:        model.MemberType(membership.MemberType),
		AggregateID:       membership.AggregateID,
		ObjectID:          membership.ObjectID,
		Roles:             membership.Roles,
		DisplayName:       membership.DisplayName,
		ChangeDate:        membership.ChangeDate,
		CreationDate:      membership.CreationDate,
		ResourceOwner:     membership.ResourceOwner,
		ResourceOwnerName: membership.ResourceOwnerName,
		Sequence:          membership.Sequence,
	}
}

func UserMembershipsToModel(memberships []*UserMembershipView) []*model.UserMembershipView {
	result := make([]*model.UserMembershipView, len(memberships))
	for i, m := range memberships {
		result[i] = UserMembershipToModel(m)
	}
	return result
}

func (u *UserMembershipView) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Sequence

	switch event.Type {
	case iam_es_model.IAMMemberAdded:
		u.setRootData(event, model.MemberTypeIam)
		err = u.setIamMemberData(event)
	case iam_es_model.IAMMemberChanged,
		iam_es_model.IAMMemberRemoved:
		err = u.setIamMemberData(event)
	case org_es_model.OrgMemberAdded:
		u.setRootData(event, model.MemberTypeOrganisation)
		err = u.setOrgMemberData(event)
	case org_es_model.OrgMemberChanged,
		org_es_model.OrgMemberRemoved:
		err = u.setOrgMemberData(event)
	case proj_es_model.ProjectMemberAdded:
		u.setRootData(event, model.MemberTypeProject)
		err = u.setProjectMemberData(event)
	case proj_es_model.ProjectMemberChanged,
		proj_es_model.ProjectMemberRemoved:
		err = u.setProjectMemberData(event)
	case proj_es_model.ProjectGrantMemberAdded:
		u.setRootData(event, model.MemberTypeProjectGrant)
		err = u.setProjectGrantMemberData(event)
	case proj_es_model.ProjectGrantMemberChanged,
		proj_es_model.ProjectGrantMemberRemoved:
		err = u.setProjectGrantMemberData(event)
	}
	return err
}

func (u *UserMembershipView) setRootData(event *models.Event, memberType model.MemberType) {
	u.CreationDate = event.CreationDate
	u.AggregateID = event.AggregateID
	u.ObjectID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
	u.MemberType = int32(memberType)
}

func (u *UserMembershipView) setIamMemberData(event *models.Event) error {
	member := new(iam_es_model.IAMMember)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.Log("MODEL-Ec9sf").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.UserID = member.UserID
	u.Roles = member.Roles
	return nil
}

func (u *UserMembershipView) setOrgMemberData(event *models.Event) error {
	member := new(org_es_model.OrgMember)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.Log("MODEL-Lps0e").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.UserID = member.UserID
	u.Roles = member.Roles
	return nil
}

func (u *UserMembershipView) setProjectMemberData(event *models.Event) error {
	member := new(proj_es_model.ProjectMember)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.Log("MODEL-Esu8k").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.UserID = member.UserID
	u.Roles = member.Roles
	return nil
}

func (u *UserMembershipView) setProjectGrantMemberData(event *models.Event) error {
	member := new(proj_es_model.ProjectGrantMember)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.Log("MODEL-MCn8s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.UserID = member.UserID
	u.ObjectID = member.GrantID
	u.Roles = member.Roles
	return nil
}
