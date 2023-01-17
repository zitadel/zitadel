package oauth

import (
	"encoding/json"

	"golang.org/x/text/language"
)

// UserInfoMapper needs to be implemented for an oauth Provider
// to map the info returned by the userEndpoint to an idp.User
type UserInfoMapper interface {
	GetID() string
	GetDisplayName() string
	GetNickName() string
	GetEmail() string
	IsEmailVerified() bool
	GetAvatarURL() string
	GetFirstName() string
	GetLastName() string
	GetPreferredUsername() string
	GetPhone() string
	IsPhoneVerified() bool
	GetPreferredLanguage() language.Tag
	GetProfile() string
	RawData() any
}

var _ UserInfoMapper = (*UserMapper)(nil)

// UserMapper is an implementation of UserInfoMapper
// it can be used in ZITADEL actions to map the raw info
type UserMapper struct {
	ID                string
	FirstName         string
	LastName          string
	DisplayName       string
	NickName          string
	PreferredUsername string
	Email             string
	EmailVerified     bool
	Phone             string
	PhoneVerified     bool
	PreferredLanguage string
	AvatarURL         string
	Profile           string
	info              map[string]interface{}
}

func (u *UserMapper) UnmarshalJSON(data []byte) error {
	if u.info == nil {
		u.info = make(map[string]interface{})
	}
	return json.Unmarshal(data, &u.info)
}

func (u *UserMapper) GetID() string {
	return u.ID
}

func (u *UserMapper) GetFirstName() string {
	return u.FirstName
}

func (u *UserMapper) GetLastName() string {
	return u.LastName
}

func (u *UserMapper) GetDisplayName() string {
	return u.DisplayName
}

func (u *UserMapper) GetNickName() string {
	return u.NickName
}

func (u *UserMapper) GetPreferredUsername() string {
	return u.PreferredUsername
}

func (u *UserMapper) GetEmail() string {
	return u.Email
}

func (u *UserMapper) IsEmailVerified() bool {
	return u.EmailVerified
}

func (u *UserMapper) GetPhone() string {
	return u.Phone
}

func (u *UserMapper) IsPhoneVerified() bool {
	return u.PhoneVerified
}

func (u *UserMapper) GetPreferredLanguage() language.Tag {
	return language.Make(u.PreferredLanguage)
}

func (u *UserMapper) GetAvatarURL() string {
	return u.AvatarURL
}

func (u *UserMapper) GetProfile() string {
	return u.Profile
}

func (u *UserMapper) RawData() any {
	return u.info
}
