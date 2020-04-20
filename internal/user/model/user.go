package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"golang.org/x/text/language"
)

type User struct {
	es_models.ObjectRoot

	State UserState
	*Profile
	*Email
	*Phone
	*Address
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

	Email           string
	IsEmailVerified bool
}

type Phone struct {
	es_models.ObjectRoot

	Phone           string
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
	USERSTATE_ACTIVE UserState = iota
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
