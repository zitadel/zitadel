package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"golang.org/x/text/language"
)

type User struct {
	es_models.ObjectRoot

	State UserState
	*Password
	*Profile
	*Email
	*Phone
	*Address
}

type Password struct {
	es_models.ObjectRoot

	SecretString   string
	SecretCrypto   *crypto.CryptoValue
	ChangeRequired bool
}

type Profile struct {
	es_models.ObjectRoot

	UserName          string
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	PreferredLanguage language.Tag
	Gender            Gender
}

type Email struct {
	es_models.ObjectRoot

	EmailAddress    string
	IsEmailVerified bool
}

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     string
	IsPhoneVerified bool
}

type Address struct {
	es_models.ObjectRoot

	Country       string
	Locality      string
	PostalCode    string
	Region        string
	StreetAddress string
}

type UserState int32

const (
	USERSTATE_UNSPECIFIED UserState = iota
	USERSTATE_ACTIVE
	USERSTATE_INACTIVE
	USERSTATE_DELETED
	USERSTATE_LOCKED
	USERSTATE_SUSPEND
	USERSTATE_INITIAL
)

type Gender int32

const (
	GENDER_UNDEFINED Gender = 0
	GENDER_FEMALE    Gender = 1
	GENDER_MALE      Gender = 2
	GENDER_DIVERSE   Gender = 3
)

func (u *User) IsValid() bool {
	if u.Profile == nil || u.FirstName == "" || u.LastName == "" || u.UserName == "" || u.Email == nil || u.EmailAddress == "" {
		return false
	}
	return true
}

func (u *User) IsActive() bool {
	return u.State == USERSTATE_ACTIVE
}

func (u *User) IsInitial() bool {
	return u.State == USERSTATE_INITIAL
}

func (u *User) IsInactive() bool {
	return u.State == USERSTATE_INACTIVE
}

func (u *User) IsLocked() bool {
	return u.State == USERSTATE_LOCKED
}
