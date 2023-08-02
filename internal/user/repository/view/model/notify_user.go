package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	org_model "github.com/zitadel/zitadel/internal/org/model"
	"github.com/zitadel/zitadel/internal/repository/user"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	NotifyUserKeyUserID        = "id"
	NotifyUserKeyResourceOwner = "resource_owner"
	NotifyUserKeyInstanceID    = "instance_id"
)

type NotifyUser struct {
	ID                 string               `json:"-" gorm:"column:id;primary_key"`
	CreationDate       time.Time            `json:"-" gorm:"column:creation_date"`
	ChangeDate         time.Time            `json:"-" gorm:"column:change_date"`
	ResourceOwner      string               `json:"-" gorm:"column:resource_owner"`
	UserName           string               `json:"userName" gorm:"column:user_name"`
	LoginNames         database.StringArray `json:"-" gorm:"column:login_names"`
	PreferredLoginName string               `json:"-" gorm:"column:preferred_login_name"`
	FirstName          string               `json:"firstName" gorm:"column:first_name"`
	LastName           string               `json:"lastName" gorm:"column:last_name"`
	NickName           string               `json:"nickName" gorm:"column:nick_name"`
	DisplayName        string               `json:"displayName" gorm:"column:display_name"`
	PreferredLanguage  string               `json:"preferredLanguage" gorm:"column:preferred_language"`
	Gender             int32                `json:"gender" gorm:"column:gender"`
	LastEmail          string               `json:"email" gorm:"column:last_email"`
	VerifiedEmail      string               `json:"-" gorm:"column:verified_email"`
	LastPhone          string               `json:"phone" gorm:"column:last_phone"`
	VerifiedPhone      string               `json:"-" gorm:"column:verified_phone"`
	PasswordSet        bool                 `json:"-" gorm:"column:password_set"`
	Sequence           uint64               `json:"-" gorm:"column:sequence"`
	State              int32                `json:"-" gorm:"-"`
	InstanceID         string               `json:"instanceID" gorm:"column:instance_id;primary_key"`
}

func (u *NotifyUser) GenerateLoginName(domain string, appendDomain bool) string {
	if !appendDomain {
		return u.UserName
	}
	return u.UserName + "@" + domain
}

func (u *NotifyUser) SetLoginNames(userLoginMustBeDomain bool, domains []*org_model.OrgDomain) {
	loginNames := make([]string, 0)
	for _, d := range domains {
		if d.Verified {
			loginNames = append(loginNames, u.GenerateLoginName(d.Domain, true))
		}
	}
	if !userLoginMustBeDomain {
		loginNames = append(loginNames, u.UserName)
	}
	u.LoginNames = loginNames
}

func (u *NotifyUser) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Sequence
	switch eventstore.EventType(event.Type) {
	case user.UserV1AddedType,
		user.UserV1RegisteredType,
		user.HumanRegisteredType,
		user.HumanAddedType,
		user.MachineAddedEventType:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		err = u.setData(event)
		if err != nil {
			return err
		}
		err = u.setPasswordData(event)
	case user.UserV1ProfileChangedType,
		user.UserV1EmailChangedType,
		user.UserV1PhoneChangedType,
		user.HumanProfileChangedType,
		user.HumanEmailChangedType,
		user.HumanPhoneChangedType,
		user.UserUserNameChangedType:
		err = u.setData(event)
	case user.UserV1EmailVerifiedType,
		user.HumanEmailVerifiedType:
		u.VerifiedEmail = u.LastEmail
	case user.UserV1PhoneRemovedType,
		user.HumanPhoneRemovedType:
		u.VerifiedPhone = ""
		u.LastPhone = ""
	case user.UserV1PhoneVerifiedType,
		user.HumanPhoneVerifiedType:
		u.VerifiedPhone = u.LastPhone
	case user.UserV1PasswordChangedType,
		user.HumanPasswordChangedType:
		err = u.setPasswordData(event)
	case user.UserRemovedType:
		u.State = int32(UserStateDeleted)
	}
	return err
}

func (u *NotifyUser) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
	u.InstanceID = event.InstanceID
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
	u.PasswordSet = password.Secret != nil || password.EncodedHash != ""
	return nil
}
