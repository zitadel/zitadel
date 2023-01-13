package oauth

import (
	"encoding/json"

	"golang.org/x/text/language"
)

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
	GetPreferredLanguange() language.Tag
	GetProfile() string
	RawData() any
}

var _ UserInfoMapper = (*UserMapper)(nil)

type UserMapper struct {
	ID          string
	DisplayName string
	info        map[string]interface{}
}

func (u *UserMapper) UnmarshalJSON(data []byte) error {
	if u.info == nil {
		u.info = make(map[string]interface{})
	}
	return json.Unmarshal(data, &u.info)
}

func (u *UserMapper) GetID() string {
	id, _ := u.info[u.ID].(string)
	return id
}

func (u *UserMapper) GetFirstName() string {
	return ""
}

func (u *UserMapper) GetLastName() string {
	return ""
}

func (u *UserMapper) GetDisplayName() string {
	displayName, _ := u.info[u.DisplayName].(string)
	return displayName
}

func (u *UserMapper) GetNickName() string {
	return ""
}

func (u *UserMapper) GetPreferredUsername() string {
	return ""
}

func (u *UserMapper) GetEmail() string {
	return ""
}

func (u *UserMapper) IsEmailVerified() bool {
	return false
}

func (u *UserMapper) GetPhone() string {
	return ""
}

func (u *UserMapper) IsPhoneVerified() bool {
	return false
}

func (u *UserMapper) GetPreferredLanguange() language.Tag {
	return language.Und
}

func (u *UserMapper) GetAvatarURL() string {
	return ""
}

func (u *UserMapper) GetProfile() string {
	return ""
}

func (u *UserMapper) RawData() any {
	return u.info
}
