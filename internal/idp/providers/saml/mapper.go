package saml

import (
	"github.com/crewjam/saml"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.User = (*UserMapper)(nil)

// UserMapper is an implementation of [idp.User].
type UserMapper struct {
	ID         string              `json:"id,omitempty"`
	Attributes map[string][]string `json:"attributes,omitempty"`
}

func NewUser() *UserMapper {
	return &UserMapper{Attributes: map[string][]string{}}
}

func (u *UserMapper) SetID(id *saml.NameID) {
	u.ID = id.Value
}

// GetID is an implementation of the [idp.User] interface.
func (u *UserMapper) GetID() string {
	return u.ID
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
func (u *UserMapper) GetEmail() domain.EmailAddress {
	return ""
}

// IsEmailVerified is an implementation of the [idp.User] interface.
func (u *UserMapper) IsEmailVerified() bool {
	return false
}

// GetPhone is an implementation of the [idp.User] interface.
func (u *UserMapper) GetPhone() domain.PhoneNumber {
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
