package model

import (
	"time"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type User struct {
	es_models.ObjectRoot

	State UserState
	*Password
	*Profile
	*Email
	*Phone
	*Address
	InitCode       *InitUserCode
	EmailCode      *EmailCode
	PhoneCode      *PhoneCode
	PasswordCode   *PasswordCode
	OTP            *OTP
	SkippedMfaInit time.Time
}

type InitUserCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
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
	GENDER_UNDEFINED Gender = iota
	GENDER_FEMALE
	GENDER_MALE
	GENDER_DIVERSE
)

func (u *User) IsValid() bool {
	return u.Profile != nil && u.FirstName != "" && u.LastName != "" && u.UserName != "" && u.Email != nil && u.EmailAddress != ""
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

func (u *User) PasswordVerified(userAgentID string) (bool, uint16) {
	return true, 0 //TODO: ???
}

func (u *User) MfaVerified(userAgentID string) (bool, uint16) {
	return true, 0 //TODO: ???
}

func (u *User) MfaTypesReadyAndSufficient(level MfaLevel) []MfaType {
	types := make([]MfaType, 0, 1)
	if u.OTP != nil && MfaIsReady(u.OTP) && MfaLevelSufficient(u.OTP, level) {
		types = append(types, u.OTP)
	}
	return types
}

func (u *User) MfaTypesPossible(level MfaLevel) []MfaType {
	types := make([]MfaType, 0, 1)
	if u.OTP == nil || !MfaIsReady(u.OTP) && MfaLevelSufficient(u.OTP, level) {
		types = append(types, u.OTP)
	}
	return types
}
