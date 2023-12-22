package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/zitadel/zitadel/internal/iam/repository/eventsourcing/model"
	org_es_model "github.com/zitadel/zitadel/internal/org/repository/eventsourcing/model"
	proj_es_model "github.com/zitadel/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UserMembershipKeyUserID        = "user_id"
	UserMembershipKeyAggregateID   = "aggregate_id"
	UserMembershipKeyObjectID      = "object_id"
	UserMembershipKeyResourceOwner = "resource_owner"
	UserMembershipKeyMemberType    = "member_type"
	UserMembershipKeyInstanceID    = "instance_id"
)

type UserMembershipView struct {
	UserID      string `json:"-" gorm:"column:user_id;primary_key"`
	MemberType  int32  `json:"-" gorm:"column:member_type;primary_key"`
	AggregateID string `json:"-" gorm:"column:aggregate_id;primary_key"`
	ObjectID    string `json:"-" gorm:"column:object_id;primary_key"`

	Roles             database.TextArray[string] `json:"-" gorm:"column:roles"`
	DisplayName       string                     `json:"-" gorm:"column:display_name"`
	CreationDate      time.Time                  `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time                  `json:"-" gorm:"column:change_date"`
	ResourceOwner     string                     `json:"-" gorm:"column:resource_owner"`
	ResourceOwnerName string                     `json:"-" gorm:"column:resource_owner_name"`
	Sequence          uint64                     `json:"-" gorm:"column:sequence"`
	InstanceID        string                     `json:"instanceID" gorm:"column:instance_id;primary_key"`
}

func (u *UserMembershipView) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Seq

	switch event.Type() {
	case instance.MemberAddedEventType:
		u.setRootData(event, model.MemberTypeIam)
		err = u.setIamMemberData(event)
	case instance.MemberChangedEventType,
		instance.MemberRemovedEventType,
		instance.MemberCascadeRemovedEventType:
		err = u.setIamMemberData(event)
	case org.MemberAddedEventType:
		u.setRootData(event, model.MemberTypeOrganisation)
		err = u.setOrgMemberData(event)
	case org.MemberChangedEventType,
		org.MemberRemovedEventType,
		org.MemberCascadeRemovedEventType:
		err = u.setOrgMemberData(event)
	case project.MemberAddedType:
		u.setRootData(event, model.MemberTypeProject)
		err = u.setProjectMemberData(event)
	case project.MemberChangedType,
		project.MemberRemovedType,
		project.MemberCascadeRemovedType:
		err = u.setProjectMemberData(event)
	case project.GrantMemberAddedType:
		u.setRootData(event, model.MemberTypeProjectGrant)
		err = u.setProjectGrantMemberData(event)
	case project.GrantMemberChangedType,
		project.GrantMemberRemovedType,
		project.GrantMemberCascadeRemovedType:
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
	u.InstanceID = event.InstanceID
}

func (u *UserMembershipView) setIamMemberData(event *models.Event) error {
	member := new(iam_es_model.IAMMember)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.New().WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.UserID = member.UserID
	u.Roles = member.Roles
	return nil
}

func (u *UserMembershipView) setOrgMemberData(event *models.Event) error {
	member := new(org_es_model.OrgMember)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.New().WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.UserID = member.UserID
	u.Roles = member.Roles
	return nil
}

func (u *UserMembershipView) setProjectMemberData(event *models.Event) error {
	member := new(proj_es_model.ProjectMember)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.New().WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.UserID = member.UserID
	u.Roles = member.Roles
	return nil
}

func (u *UserMembershipView) setProjectGrantMemberData(event *models.Event) error {
	member := new(proj_es_model.ProjectGrantMember)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.New().WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.UserID = member.UserID
	u.ObjectID = member.GrantID
	u.Roles = member.Roles
	return nil
}
