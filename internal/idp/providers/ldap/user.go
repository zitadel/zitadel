package ldap

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
)

type User struct {
	id                string
	firstName         string
	lastName          string
	displayName       string
	nickName          string
	preferredUsername string
	email             domain.EmailAddress
	emailVerified     bool
	phone             domain.PhoneNumber
	phoneVerified     bool
	preferredLanguage language.Tag
	avatarURL         string
	profile           string
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
	return u.id
}
func (u *User) GetFirstName() string {
	return u.firstName
}
func (u *User) GetLastName() string {
	return u.lastName
}
func (u *User) GetDisplayName() string {
	return u.displayName
}
func (u *User) GetNickname() string {
	return u.nickName
}
func (u *User) GetPreferredUsername() string {
	return u.preferredUsername
}
func (u *User) GetEmail() domain.EmailAddress {
	return u.email
}
func (u *User) IsEmailVerified() bool {
	return u.emailVerified
}
func (u *User) GetPhone() domain.PhoneNumber {
	return u.phone
}
func (u *User) IsPhoneVerified() bool {
	return u.phoneVerified
}
func (u *User) GetPreferredLanguage() language.Tag {
	return u.preferredLanguage
}
func (u *User) GetAvatarURL() string {
	return u.avatarURL
}
func (u *User) GetProfile() string {
	return u.profile
}
