package oauth

import (
	"encoding/json"
	"fmt"
	"strconv"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.User = (*UserMapper)(nil)

// UserMapper is an implementation of [idp.User].
// It can be used in ZITADEL actions to map the `RawInfo`
type UserMapper struct {
	idAttribute string
	RawInfo     map[string]interface{}
}

func NewUserMapper(idAttribute string) *UserMapper {
	return &UserMapper{
		idAttribute: idAttribute,
		RawInfo:     make(map[string]interface{}),
	}
}

func (u *UserMapper) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &u.RawInfo)
}

// GetID is an implementation of the [idp.User] interface.
func (u *UserMapper) GetID() string {
	id, ok := u.RawInfo[u.idAttribute]
	if !ok {
		return ""
	}
	switch i := id.(type) {
	case string:
		return i
	case int:
		return strconv.Itoa(i)
	case float64:
		return strconv.FormatFloat(i, 'f', -1, 64)
	default:
		return fmt.Sprint(i)
	}
}

// GetFirstName is an implementation of the [idp.User] interface.
func (u *UserMapper) GetFirstName() string {
	return ""
}

// GetLastName is an implementation of the [idp.User] interface.
func (u *UserMapper) GetLastName() string {
	return ""
}

// GetDisplayName is an implementation of the [idp.User] interface.
func (u *UserMapper) GetDisplayName() string {
	return ""
}

// GetNickname is an implementation of the [idp.User] interface.
func (u *UserMapper) GetNickname() string {
	return ""
}

// GetPreferredUsername is an implementation of the [idp.User] interface.
func (u *UserMapper) GetPreferredUsername() string {
	return ""
}

// GetEmail is an implementation of the [idp.User] interface.
func (u *UserMapper) GetEmail() string {
	return ""
}

// IsEmailVerified is an implementation of the [idp.User] interface.
func (u *UserMapper) IsEmailVerified() bool {
	return false
}

// GetPhone is an implementation of the [idp.User] interface.
func (u *UserMapper) GetPhone() string {
	return ""
}

// IsPhoneVerified is an implementation of the [idp.User] interface.
func (u *UserMapper) IsPhoneVerified() bool {
	return false
}

// GetPreferredLanguage is an implementation of the [idp.User] interface.
func (u *UserMapper) GetPreferredLanguage() language.Tag {
	return language.Und
}

// GetAvatarURL is an implementation of the [idp.User] interface.
func (u *UserMapper) GetAvatarURL() string {
	return ""
}

// GetProfile is an implementation of the [idp.User] interface.
func (u *UserMapper) GetProfile() string {
	return ""
}
