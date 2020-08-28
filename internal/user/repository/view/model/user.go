package model

import (
	"encoding/json"
	"time"

	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/lib/pq"

	"github.com/caos/logging"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	UserKeyUserID        = "id"
	UserKeyUserName      = "user_name"
	UserKeyFirstName     = "first_name"
	UserKeyLastName      = "last_name"
	UserKeyNickName      = "nick_name"
	UserKeyDisplayName   = "display_name"
	UserKeyEmail         = "email"
	UserKeyState         = "user_state"
	UserKeyResourceOwner = "resource_owner"
	UserKeyLoginNames    = "login_names"
	UserKeyType          = "user_type"
)

type userType string

const (
	userTypeHuman   = "human"
	userTypeMachine = "machine"
)

type UserView struct {
	ID                 string         `json:"-" gorm:"column:id;primary_key"`
	CreationDate       time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate         time.Time      `json:"-" gorm:"column:change_date"`
	ResourceOwner      string         `json:"-" gorm:"column:resource_owner"`
	State              int32          `json:"-" gorm:"column:user_state"`
	LastLogin          time.Time      `json:"-" gorm:"column:last_login"`
	LoginNames         pq.StringArray `json:"-" gorm:"column:login_names"`
	PreferredLoginName string         `json:"-" gorm:"column:preferred_login_name"`
	Sequence           uint64         `json:"-" gorm:"column:sequence"`
	Type               userType       `json:"-" gorm:"column:user_type"`
	UserName           string         `json:"userName" gorm:"column:user_name"`
	*MachineView
	*HumanView
}

type HumanView struct {
	FirstName         string    `json:"firstName" gorm:"column:first_name"`
	LastName          string    `json:"lastName" gorm:"column:last_name"`
	NickName          string    `json:"nickName" gorm:"column:nick_name"`
	DisplayName       string    `json:"displayName" gorm:"column:display_name"`
	PreferredLanguage string    `json:"preferredLanguage" gorm:"column:preferred_language"`
	Gender            int32     `json:"gender" gorm:"column:gender"`
	Email             string    `json:"email" gorm:"column:email"`
	IsEmailVerified   bool      `json:"-" gorm:"column:is_email_verified"`
	Phone             string    `json:"phone" gorm:"column:phone"`
	IsPhoneVerified   bool      `json:"-" gorm:"column:is_phone_verified"`
	Country           string    `json:"country" gorm:"column:country"`
	Locality          string    `json:"locality" gorm:"column:locality"`
	PostalCode        string    `json:"postalCode" gorm:"column:postal_code"`
	Region            string    `json:"region" gorm:"column:region"`
	StreetAddress     string    `json:"streetAddress" gorm:"column:street_address"`
	OTPState          int32     `json:"-" gorm:"column:otp_state"`
	MfaMaxSetUp       int32     `json:"-" gorm:"column:mfa_max_set_up"`
	MfaInitSkipped    time.Time `json:"-" gorm:"column:mfa_init_skipped"`
	InitRequired      bool      `json:"-" gorm:"column:init_required"`

	PasswordSet            bool      `json:"-" gorm:"column:password_set"`
	PasswordChangeRequired bool      `json:"-" gorm:"column:password_change_required"`
	PasswordChanged        time.Time `json:"-" gorm:"column:password_change"`
}

func (h *HumanView) IsZero() bool {
	return h == nil || h.FirstName == ""
}

type MachineView struct {
	Name        string `json:"name" gorm:"column:machine_name"`
	Description string `json:"description" gorm:"column:machine_description"`
}

func (m *MachineView) IsZero() bool {
	return m == nil || m.Name == ""
}

func UserToModel(user *UserView) *model.UserView {
	userView := &model.UserView{
		ID:                 user.ID,
		UserName:           user.UserName,
		ChangeDate:         user.ChangeDate,
		CreationDate:       user.CreationDate,
		ResourceOwner:      user.ResourceOwner,
		State:              model.UserState(user.State),
		LastLogin:          user.LastLogin,
		PreferredLoginName: user.PreferredLoginName,
		LoginNames:         user.LoginNames,
		Sequence:           user.Sequence,
	}
	if !user.HumanView.IsZero() {
		userView.HumanView = &model.HumanView{
			PasswordSet:            user.PasswordSet,
			PasswordChangeRequired: user.PasswordChangeRequired,
			PasswordChanged:        user.PasswordChanged,
			FirstName:              user.FirstName,
			LastName:               user.LastName,
			NickName:               user.NickName,
			DisplayName:            user.DisplayName,
			PreferredLanguage:      user.PreferredLanguage,
			Gender:                 model.Gender(user.Gender),
			Email:                  user.Email,
			IsEmailVerified:        user.IsEmailVerified,
			Phone:                  user.Phone,
			IsPhoneVerified:        user.IsPhoneVerified,
			Country:                user.Country,
			Locality:               user.Locality,
			PostalCode:             user.PostalCode,
			Region:                 user.Region,
			StreetAddress:          user.StreetAddress,
			OTPState:               model.MfaState(user.OTPState),
			MfaMaxSetUp:            req_model.MfaLevel(user.MfaMaxSetUp),
			MfaInitSkipped:         user.MfaInitSkipped,
			InitRequired:           user.InitRequired,
		}
	}

	if !user.MachineView.IsZero() {
		userView.MachineView = &model.MachineView{
			Description: user.MachineView.Description,
			Name:        user.MachineView.Name,
		}
	}
	return userView
}

func UsersToModel(users []*UserView) []*model.UserView {
	result := make([]*model.UserView, len(users))
	for i, p := range users {
		result[i] = UserToModel(p)
	}
	return result
}

func (u *UserView) GenerateLoginName(domain string, appendDomain bool) string {
	var name string
	if u.MachineView != nil {
		name = u.MachineView.Name
	} else {
		name = u.UserName
	}
	if !appendDomain {
		return name
	}
	return name + "@" + domain
}

func (u *UserView) SetLoginNames(policy *org_model.OrgIamPolicy, domains []*org_model.OrgDomain) {
	loginNames := make([]string, 0)
	for _, d := range domains {
		if d.Verified {
			loginNames = append(loginNames, u.GenerateLoginName(d.Domain, true))
		}
	}
	if !policy.UserLoginMustBeDomain {
		if u.MachineView != nil {
			loginNames = append(loginNames, u.MachineView.Name)
		} else {
			loginNames = append(loginNames, u.UserName)
		}
		loginNames = append(loginNames, u.UserName)
	}
	u.LoginNames = loginNames
}

func (u *UserView) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Sequence
	switch event.Type {
	case es_model.MachineAdded:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		u.Type = userTypeMachine
		err = u.setData(event)
		if err != nil {
			return err
		}
	case es_model.UserAdded,
		es_model.UserRegistered,
		es_model.HumanRegistered,
		es_model.HumanAdded:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		u.Type = userTypeHuman
		err = u.setData(event)
		if err != nil {
			return err
		}
		err = u.setPasswordData(event)
	case es_model.UserPasswordChanged,
		es_model.HumanPasswordChanged:
		err = u.setPasswordData(event)
	case es_model.UserProfileChanged,
		es_model.UserAddressChanged,
		es_model.DomainClaimed,
		es_model.HumanProfileChanged,
		es_model.HumanAddressChanged,
		es_model.MachineChanged:
		err = u.setData(event)
	case es_model.UserEmailChanged,
		es_model.HumanEmailChanged:
		u.IsEmailVerified = false
		err = u.setData(event)
	case es_model.UserEmailVerified,
		es_model.HumanEmailVerified:
		u.IsEmailVerified = true
	case es_model.UserPhoneChanged,
		es_model.HumanPhoneChanged:
		u.IsPhoneVerified = false
		err = u.setData(event)
	case es_model.UserPhoneVerified,
		es_model.HumanPhoneVerified:
		u.IsPhoneVerified = true
	case es_model.UserPhoneRemoved,
		es_model.HumanPhoneRemoved:
		u.Phone = ""
		u.IsPhoneVerified = false
	case es_model.UserDeactivated:
		u.State = int32(model.UserStateInactive)
	case es_model.UserReactivated,
		es_model.UserUnlocked:
		u.State = int32(model.UserStateActive)
	case es_model.UserLocked:
		u.State = int32(model.UserStateLocked)
	case es_model.MfaOtpAdded,
		es_model.HumanMfaOtpAdded:
		u.OTPState = int32(model.MfaStateNotReady)
	case es_model.MfaOtpVerified,
		es_model.HumanMfaOtpVerified:
		u.OTPState = int32(model.MfaStateReady)
		u.MfaInitSkipped = time.Time{}
	case es_model.MfaOtpRemoved,
		es_model.HumanMfaOtpRemoved:
		u.OTPState = int32(model.MfaStateUnspecified)
	case es_model.MfaInitSkipped,
		es_model.HumanMfaInitSkipped:
		u.MfaInitSkipped = event.CreationDate
	case es_model.InitializedUserCodeAdded,
		es_model.InitializedHumanCodeAdded:
		u.InitRequired = true
	case es_model.InitializedUserCheckSucceeded,
		es_model.InitializedHumanCheckSucceeded:
		u.InitRequired = false
	}
	u.ComputeObject()
	return err
}

func (u *UserView) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
}

func (u *UserView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("MODEL-lso9e").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-8iows", "could not unmarshal data")
	}
	return nil
}

func (u *UserView) setPasswordData(event *models.Event) error {
	password := new(es_model.Password)
	if err := json.Unmarshal(event.Data, password); err != nil {
		logging.Log("MODEL-sdw4r").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.PasswordSet = password.Secret != nil
	u.PasswordChangeRequired = password.ChangeRequired
	u.PasswordChanged = event.CreationDate
	return nil
}

func (u *UserView) ComputeObject() {
	if u.MachineView != nil {
		if u.State == int32(model.UserStateUnspecified) {
			u.State = int32(model.UserStateActive)
		}
		return
	}
	if u.State == int32(model.UserStateUnspecified) || u.State == int32(model.UserStateInitial) {
		if u.IsEmailVerified {
			u.State = int32(model.UserStateActive)
		} else {
			u.State = int32(model.UserStateInitial)
		}
	}
	if u.OTPState != int32(model.MfaStateReady) {
		u.MfaMaxSetUp = int32(req_model.MfaLevelNotSetUp)
	}
	if u.OTPState == int32(model.MfaStateReady) {
		u.MfaMaxSetUp = int32(req_model.MfaLevelSoftware)
	}
}
