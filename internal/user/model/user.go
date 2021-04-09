package model

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"strings"

	caos_errors "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type User struct {
	es_models.ObjectRoot
	State    UserState
	UserName string

	*Human
	*Machine
}

type UserState int32

const (
	UserStateUnspecified UserState = iota
	UserStateActive
	UserStateInactive
	UserStateDeleted
	UserStateLocked
	UserStateSuspend
	UserStateInitial
)

func (u *User) CheckOrgIAMPolicy(policy *iam_model.OrgIAMPolicyView) error {
	if policy == nil {
		return caos_errors.ThrowPreconditionFailed(nil, "MODEL-zSH73", "Errors.Users.OrgIamPolicyNil")
	}
	if !policy.UserLoginMustBeDomain && u.Profile != nil && u.UserName == "" && u.Email != nil {
		u.UserName = u.EmailAddress
	}
	return nil
}

func (u *User) SetNamesAsDisplayname() {
	if u.Profile != nil && u.DisplayName == "" && u.FirstName != "" && u.LastName != "" {
		u.DisplayName = u.FirstName + " " + u.LastName
	}
}

type UserChanges struct {
	Changes      []*UserChange
	LastSequence uint64
}

type UserChange struct {
	ChangeDate   *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType    string               `json:"eventType,omitempty"`
	Sequence     uint64               `json:"sequence,omitempty"`
	ModifierID   string               `json:"modifierUser,omitempty"`
	ModifierName string               `json:"-"`
	Data         interface{}          `json:"data,omitempty"`
}

func (u *User) IsActive() bool {
	return u.State == UserStateActive
}

func (u *User) IsInitial() bool {
	return u.State == UserStateInitial
}

func (u *User) IsInactive() bool {
	return u.State == UserStateInactive
}

func (u *User) IsLocked() bool {
	return u.State == UserStateLocked
}

func (u *User) IsValid() bool {
	if u.Human == nil && u.Machine == nil || u.UserName == "" {
		return false
	}
	if u.Human != nil {
		return u.Human.IsValid()
	}
	return u.Machine.IsValid()
}

func (u *User) CheckOrgIamPolicy(policy *iam_model.OrgIAMPolicy) error {
	if policy == nil {
		return caos_errors.ThrowPreconditionFailed(nil, "MODEL-zSH7j", "Errors.Users.OrgIamPolicyNil")
	}
	if !policy.UserLoginMustBeDomain && u.Profile != nil && u.UserName == "" && u.Email != nil {
		u.UserName = u.EmailAddress
	}
	return nil
}
