package oauth

import (
	"encoding/json"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.User = (*UserMapper)(nil)

// UserMapper is an implementation of [idp.User].
// It can be used in ZITADEL actions to map the raw `info`
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

// GetID is an implementation of the [idp.User] interface.
func (u *UserMapper) GetID() string {
	return u.ID
}

// GetFirstName is an implementation of the [idp.User] interface.
func (u *UserMapper) GetFirstName() string {
	return u.FirstName
}

// GetLastName is an implementation of the [idp.User] interface.
func (u *UserMapper) GetLastName() string {
	return u.LastName
}

// GetDisplayName is an implementation of the [idp.User] interface.
func (u *UserMapper) GetDisplayName() string {
	return u.DisplayName
}

// GetNickname is an implementation of the [idp.User] interface.
func (u *UserMapper) GetNickname() string {
	return u.NickName
}

// GetPreferredUsername is an implementation of the [idp.User] interface.
func (u *UserMapper) GetPreferredUsername() string {
	return u.PreferredUsername
}

// GetEmail is an implementation of the [idp.User] interface.
func (u *UserMapper) GetEmail() string {
	return u.Email
}

// IsEmailVerified is an implementation of the [idp.User] interface.
func (u *UserMapper) IsEmailVerified() bool {
	return u.EmailVerified
}

// GetPhone is an implementation of the [idp.User] interface.
func (u *UserMapper) GetPhone() string {
	return u.Phone
}

// IsPhoneVerified is an implementation of the [idp.User] interface.
func (u *UserMapper) IsPhoneVerified() bool {
	return u.PhoneVerified
}

// GetPreferredLanguage is an implementation of the [idp.User] interface.
func (u *UserMapper) GetPreferredLanguage() language.Tag {
	return language.Make(u.PreferredLanguage)
}

// GetAvatarURL is an implementation of the [idp.User] interface.
func (u *UserMapper) GetAvatarURL() string {
	return u.AvatarURL
}

// GetProfile is an implementation of the [idp.User] interface.
func (u *UserMapper) GetProfile() string {
	return u.Profile
}
