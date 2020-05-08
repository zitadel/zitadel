package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"time"
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
)

type UserView struct {
	ID                string    `json:"-" gorm:"column:id;primary_key"`
	CreationDate      time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner     string    `json:"-" gorm:"column:resource_owner"`
	State             int32     `json:"-" gorm:"column:user_state"`
	PasswordChanged   time.Time `json:"-" gorm:"column:password_change"`
	LastLogin         time.Time `json:"-" gorm:"column:last_login"`
	UserName          string    `json:"userName" gorm:"column:user_name"`
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
	Sequence          uint64    `json:"-" gorm:"column:sequence"`
}

func UserFromModel(user *model.UserView) *UserView {
	return &UserView{
		ID:                user.ID,
		ChangeDate:        user.ChangeDate,
		CreationDate:      user.CreationDate,
		ResourceOwner:     user.ResourceOwner,
		State:             int32(user.State),
		PasswordChanged:   user.PasswordChanged,
		LastLogin:         user.LastLogin,
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		NickName:          user.NickName,
		DisplayName:       user.DisplayName,
		PreferredLanguage: user.PreferredLanguage,
		Gender:            int32(user.Gender),
		Email:             user.Email,
		IsEmailVerified:   user.IsEmailVerified,
		Phone:             user.Phone,
		IsPhoneVerified:   user.IsPhoneVerified,
		Country:           user.Country,
		Locality:          user.Locality,
		PostalCode:        user.PostalCode,
		Region:            user.Region,
		StreetAddress:     user.StreetAddress,
		OTPState:          int32(user.OTPState),
		Sequence:          user.Sequence,
	}
}

func UserToModel(user *UserView) *model.UserView {
	return &model.UserView{
		ID:                user.ID,
		ChangeDate:        user.ChangeDate,
		CreationDate:      user.CreationDate,
		ResourceOwner:     user.ResourceOwner,
		State:             model.UserState(user.State),
		PasswordChanged:   user.PasswordChanged,
		LastLogin:         user.LastLogin,
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		NickName:          user.NickName,
		DisplayName:       user.DisplayName,
		PreferredLanguage: user.PreferredLanguage,
		Gender:            model.Gender(user.Gender),
		Email:             user.Email,
		IsEmailVerified:   user.IsEmailVerified,
		Phone:             user.Phone,
		IsPhoneVerified:   user.IsPhoneVerified,
		Country:           user.Country,
		Locality:          user.Locality,
		PostalCode:        user.PostalCode,
		Region:            user.Region,
		StreetAddress:     user.StreetAddress,
		OTPState:          model.MfaState(user.OTPState),
		Sequence:          user.Sequence,
	}
}

func UsersToModel(users []*UserView) []*model.UserView {
	result := make([]*model.UserView, 0)
	for _, p := range users {
		result = append(result, UserToModel(p))
	}
	return result
}

func (p *UserView) AppendEvent(event *models.Event) (err error) {
	p.ChangeDate = event.CreationDate
	p.Sequence = event.Sequence
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered:
		p.CreationDate = event.CreationDate
		p.setRootData(event)
		err = p.setData(event)
	case es_model.UserProfileChanged,
		es_model.UserAddressChanged:
		err = p.setData(event)
	case es_model.UserEmailChanged:
		p.IsEmailVerified = false
		err = p.setData(event)
	case es_model.UserEmailVerified:
		p.IsEmailVerified = true
	case es_model.UserPhoneChanged:
		p.IsPhoneVerified = false
		err = p.setData(event)
	case es_model.UserPhoneVerified:
		p.IsPhoneVerified = true
	case es_model.UserDeactivated:
		p.State = int32(model.USERSTATE_INACTIVE)
	case es_model.UserReactivated,
		es_model.UserUnlocked:
		p.State = int32(model.USERSTATE_ACTIVE)
	case es_model.UserLocked:
		p.State = int32(model.USERSTATE_LOCKED)
	case es_model.MfaOtpAdded:
		p.OTPState = int32(model.MFASTATE_NOTREADY)
	case es_model.MfaOtpVerified:
		p.OTPState = int32(model.MFASTATE_READY)
	case es_model.MfaOtpRemoved:
		p.OTPState = int32(model.MFASTATE_UNSPECIFIED)
	}
	p.ComputeObject()
	return err
}

func (u *UserView) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
}

func (u *UserView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("EVEN-lso9e").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (u *UserView) ComputeObject() {
	if u.State == int32(model.USERSTATE_UNSPECIFIED) || u.State == int32(model.USERSTATE_INITIAL) {
		if u.IsEmailVerified {
			u.State = int32(model.USERSTATE_ACTIVE)
		} else {
			u.State = int32(model.USERSTATE_INITIAL)
		}
	}
}
