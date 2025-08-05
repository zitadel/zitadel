package ldap

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
)

type User struct {
	ID                string              `json:"id,omitempty"`
	FirstName         string              `json:"firstName,omitempty"`
	LastName          string              `json:"lastName,omitempty"`
	DisplayName       string              `json:"displayName,omitempty"`
	NickName          string              `json:"nickName,omitempty"`
	PreferredUsername string              `json:"preferredUsername,omitempty"`
	Email             domain.EmailAddress `json:"email,omitempty"`
	EmailVerified     bool                `json:"emailVerified,omitempty"`
	Phone             domain.PhoneNumber  `json:"phone,omitempty"`
	PhoneVerified     bool                `json:"phoneVerified,omitempty"`
	PreferredLanguage language.Tag        `json:"preferredLanguage,omitempty"`
	AvatarURL         string              `json:"avatarURL,omitempty"`
	Profile           string              `json:"profile,omitempty"`
}

func NewUser(
	id string,
	firstName string,
	lastName string,
	displayName string,
	nickName string,
	preferredUsername string,
	email domain.EmailAddress,
	emailVerified bool,
	phone domain.PhoneNumber,
	phoneVerified bool,
	preferredLanguage language.Tag,
	avatarURL string,
	profile string,
) *User {
	return &User{
		id,
		firstName,
		lastName,
		displayName,
		nickName,
		preferredUsername,
		email,
		emailVerified,
		phone,
		phoneVerified,
		preferredLanguage,
		avatarURL,
		profile,
	}
}

func (u *User) GetID() string {
	return u.ID
}
func (u *User) GetFirstName() string {
	return u.FirstName
}
func (u *User) GetLastName() string {
	return u.LastName
}
func (u *User) GetDisplayName() string {
	return u.DisplayName
}
func (u *User) GetNickname() string {
	return u.NickName
}
func (u *User) GetPreferredUsername() string {
	return u.PreferredUsername
}
func (u *User) GetEmail() domain.EmailAddress {
	return u.Email
}
func (u *User) IsEmailVerified() bool {
	return u.EmailVerified
}
func (u *User) GetPhone() domain.PhoneNumber {
	return u.Phone
}
func (u *User) IsPhoneVerified() bool {
	return u.PhoneVerified
}
func (u *User) GetPreferredLanguage() language.Tag {
	return u.PreferredLanguage
}
func (u *User) GetAvatarURL() string {
	return u.AvatarURL
}
func (u *User) GetProfile() string {
	return u.Profile
}
