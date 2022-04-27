package actions

import (
	"github.com/zitadel/zitadel/internal/domain"
	"golang.org/x/text/language"
)

type API map[string]interface{}

func (a API) set(name string, value interface{}) {
	map[string]interface{}(a)[name] = value
}

func (a *API) SetHuman(human *domain.Human) *API {
	a.set("setFirstName", func(firstName string) {
		human.FirstName = firstName
	})
	a.set("setLastName", func(lastName string) {
		human.LastName = lastName
	})
	a.set("setNickName", func(nickName string) {
		human.NickName = nickName
	})
	a.set("setDisplayName", func(displayName string) {
		human.DisplayName = displayName
	})
	a.set("setPreferredLanguage", func(preferredLanguage string) {
		human.PreferredLanguage = language.Make(preferredLanguage)
	})
	a.set("setGender", func(gender domain.Gender) {
		human.Gender = gender
	})
	a.set("setUsername", func(username string) {
		human.Username = username
	})
	a.set("setEmail", func(email string) {
		if human.Email == nil {
			human.Email = &domain.Email{}
		}
		human.Email.EmailAddress = email
	})
	a.set("setEmailVerified", func(verified bool) {
		if human.Email == nil {
			return
		}
		human.Email.IsEmailVerified = verified
	})
	a.set("setPhone", func(email string) {
		if human.Phone == nil {
			human.Phone = &domain.Phone{}
		}
		human.Phone.PhoneNumber = email
	})
	a.set("setPhoneVerified", func(verified bool) {
		if human.Phone == nil {
			return
		}
		human.Phone.IsPhoneVerified = verified
	})
	return a
}

func (a *API) SetExternalUser(user *domain.ExternalUser) *API {
	a.set("setFirstName", func(firstName string) {
		user.FirstName = firstName
	})
	a.set("setLastName", func(lastName string) {
		user.LastName = lastName
	})
	a.set("setNickName", func(nickName string) {
		user.NickName = nickName
	})
	a.set("setDisplayName", func(displayName string) {
		user.DisplayName = displayName
	})
	a.set("setPreferredLanguage", func(preferredLanguage string) {
		user.PreferredLanguage = language.Make(preferredLanguage)
	})
	a.set("setPreferredUsername", func(username string) {
		user.PreferredUsername = username
	})
	a.set("setEmail", func(email string) {
		user.Email = email
	})
	a.set("setEmailVerified", func(verified bool) {
		user.IsEmailVerified = verified
	})
	a.set("setPhone", func(phone string) {
		user.Phone = phone
	})
	a.set("setPhoneVerified", func(verified bool) {
		user.IsPhoneVerified = verified
	})
	return a
}

func (a *API) SetMetadata(metadata *[]*domain.Metadata) *API {
	a.set("metadata", metadata)
	return a
}

func (a *API) SetUserGrants(usergrants *[]UserGrant) *API {
	a.set("userGrants", usergrants)
	return a
}
