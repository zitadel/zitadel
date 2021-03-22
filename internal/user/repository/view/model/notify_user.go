package model

import (
	"encoding/json"
	"time"

	iam_model "github.com/caos/zitadel/internal/iam/model"

	"github.com/caos/logging"
	"github.com/lib/pq"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	NotifyUserKeyUserID        = "id"
	NotifyUserKeyResourceOwner = "resource_owner"
)

type NotifyUser struct {
	ID                 string         `json:"-" gorm:"column:id;primary_key"`
	CreationDate       time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate         time.Time      `json:"-" gorm:"column:change_date"`
	ResourceOwner      string         `json:"-" gorm:"column:resource_owner"`
	UserName           string         `json:"userName" gorm:"column:user_name"`
	LoginNames         pq.StringArray `json:"-" gorm:"column:login_names"`
	PreferredLoginName string         `json:"-" gorm:"column:preferred_login_name"`
	FirstName          string         `json:"firstName" gorm:"column:first_name"`
	LastName           string         `json:"lastName" gorm:"column:last_name"`
	NickName           string         `json:"nickName" gorm:"column:nick_name"`
	DisplayName        string         `json:"displayName" gorm:"column:display_name"`
	PreferredLanguage  string         `json:"preferredLanguage" gorm:"column:preferred_language"`
	Gender             int32          `json:"gender" gorm:"column:gender"`
	LastEmail          string         `json:"email" gorm:"column:last_email"`
	VerifiedEmail      string         `json:"-" gorm:"column:verified_email"`
	LastPhone          string         `json:"phone" gorm:"column:last_phone"`
	VerifiedPhone      string         `json:"-" gorm:"column:verified_phone"`
	PasswordSet        bool           `json:"-" gorm:"column:password_set"`
	Sequence           uint64         `json:"-" gorm:"column:sequence"`
	State              int32          `json:"-" gorm:"-"`
}

func NotifyUserFromModel(user *model.NotifyUser) *NotifyUser {
	return &NotifyUser{
		ID:                 user.ID,
		ChangeDate:         user.ChangeDate,
		CreationDate:       user.CreationDate,
		ResourceOwner:      user.ResourceOwner,
		UserName:           user.UserName,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		NickName:           user.NickName,
		DisplayName:        user.DisplayName,
		PreferredLanguage:  user.PreferredLanguage,
		Gender:             int32(user.Gender),
		LastEmail:          user.LastEmail,
		VerifiedEmail:      user.VerifiedEmail,
		LastPhone:          user.LastPhone,
		VerifiedPhone:      user.VerifiedPhone,
		PasswordSet:        user.PasswordSet,
		Sequence:           user.Sequence,
	}
}

func NotifyUserToModel(user *NotifyUser) *model.NotifyUser {
	return &model.NotifyUser{
		ID:                 user.ID,
		ChangeDate:         user.ChangeDate,
		CreationDate:       user.CreationDate,
		ResourceOwner:      user.ResourceOwner,
		UserName:           user.UserName,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		NickName:           user.NickName,
		DisplayName:        user.DisplayName,
		PreferredLanguage:  user.PreferredLanguage,
		Gender:             model.Gender(user.Gender),
		LastEmail:          user.LastEmail,
		VerifiedEmail:      user.VerifiedEmail,
		LastPhone:          user.LastPhone,
		VerifiedPhone:      user.VerifiedPhone,
		PasswordSet:        user.PasswordSet,
		Sequence:           user.Sequence,
	}
}

func (u *NotifyUser) GenerateLoginName(domain string, appendDomain bool) string {
	if !appendDomain {
		return u.UserName
	}
	return u.UserName + "@" + domain
}

func (u *NotifyUser) SetLoginNames(policy *iam_model.OrgIAMPolicy, domains []*org_model.OrgDomain) {
	loginNames := make([]string, 0)
	for _, d := range domains {
		if d.Verified {
			loginNames = append(loginNames, u.GenerateLoginName(d.Domain, true))
		}
	}
	if !policy.UserLoginMustBeDomain {
		loginNames = append(loginNames, u.UserName)
	}
	u.LoginNames = loginNames
}

func (u *NotifyUser) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Sequence
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered,
		es_model.HumanRegistered,
		es_model.HumanAdded,
		es_model.MachineAdded:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		err = u.setData(event)
		if err != nil {
			return err
		}
		err = u.setPasswordData(event)
	case es_model.UserProfileChanged,
		es_model.UserEmailChanged,
		es_model.UserPhoneChanged,
		es_model.HumanProfileChanged,
		es_model.HumanEmailChanged,
		es_model.HumanPhoneChanged,
		es_model.UserUserNameChanged:
		err = u.setData(event)
	case es_model.UserEmailVerified,
		es_model.HumanEmailVerified:
		u.VerifiedEmail = u.LastEmail
	case es_model.UserPhoneRemoved,
		es_model.HumanPhoneRemoved:
		u.VerifiedPhone = ""
		u.LastPhone = ""
	case es_model.UserPhoneVerified,
		es_model.HumanPhoneVerified:
		u.VerifiedPhone = u.LastPhone
	case es_model.UserPasswordChanged,
		es_model.HumanPasswordChanged:
		err = u.setPasswordData(event)
	case es_model.UserRemoved:
		u.State = int32(UserStateDeleted)
	}
	return err
}

func (u *NotifyUser) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
}

func (u *NotifyUser) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("MODEL-lso9e").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-8iows", "could not unmarshal data")
	}
	return nil
}

func (u *NotifyUser) setPasswordData(event *models.Event) error {
	password := new(es_model.Password)
	if err := json.Unmarshal(event.Data, password); err != nil {
		logging.Log("MODEL-dfhw6").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-BHFD2", "could not unmarshal data")
	}
	u.PasswordSet = password.Secret != nil
	return nil
}
