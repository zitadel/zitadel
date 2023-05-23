package ldap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProvider_Options(t *testing.T) {
	type fields struct {
		name              string
		servers           []string
		baseDN            string
		bindDN            string
		bindPassword      string
		userBase          string
		userObjectClasses []string
		userFilters       []string
		timeout           time.Duration
		loginUrl          string
		opts              []ProviderOpts
	}
	type want struct {
		name                       string
		startTls                   bool
		linkingAllowed             bool
		creationAllowed            bool
		autoCreation               bool
		autoUpdate                 bool
		idAttribute                string
		firstNameAttribute         string
		lastNameAttribute          string
		displayNameAttribute       string
		nickNameAttribute          string
		preferredUsernameAttribute string
		emailAttribute             string
		emailVerifiedAttribute     string
		phoneAttribute             string
		phoneVerifiedAttribute     string
		preferredLanguageAttribute string
		avatarURLAttribute         string
		profileAttribute           string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "default",
			fields: fields{
				name:              "ldap",
				servers:           []string{"server"},
				baseDN:            "base",
				bindDN:            "binddn",
				bindPassword:      "password",
				userBase:          "user",
				userObjectClasses: []string{"object"},
				userFilters:       []string{"filter"},
				timeout:           30 * time.Second,
				loginUrl:          "url",
				opts:              nil,
			},
			want: want{
				name:            "ldap",
				startTls:        true,
				linkingAllowed:  false,
				creationAllowed: false,
				autoCreation:    false,
				autoUpdate:      false,
				idAttribute:     "",
			},
		},
		{
			name: "all true",
			fields: fields{
				name:              "ldap",
				servers:           []string{"server"},
				baseDN:            "base",
				bindDN:            "binddn",
				bindPassword:      "password",
				userBase:          "user",
				userObjectClasses: []string{"object"},
				userFilters:       []string{"filter"},
				timeout:           30 * time.Second,
				loginUrl:          "url",
				opts: []ProviderOpts{
					WithoutStartTLS(),
					WithLinkingAllowed(),
					WithCreationAllowed(),
					WithAutoCreation(),
					WithAutoUpdate(),
				},
			},
			want: want{
				name:            "ldap",
				startTls:        false,
				linkingAllowed:  true,
				creationAllowed: true,
				autoCreation:    true,
				autoUpdate:      true,
				idAttribute:     "",
			},
		}, {
			name: "all true, attributes set",
			fields: fields{
				name:              "ldap",
				servers:           []string{"server"},
				baseDN:            "base",
				bindDN:            "binddn",
				bindPassword:      "password",
				userBase:          "user",
				userObjectClasses: []string{"object"},
				userFilters:       []string{"filter"},
				timeout:           30 * time.Second,
				loginUrl:          "url",
				opts: []ProviderOpts{
					WithoutStartTLS(),
					WithLinkingAllowed(),
					WithCreationAllowed(),
					WithAutoCreation(),
					WithAutoUpdate(),
					WithCustomIDAttribute("id"),
					WithFirstNameAttribute("first"),
					WithLastNameAttribute("last"),
					WithDisplayNameAttribute("display"),
					WithNickNameAttribute("nick"),
					WithPreferredUsernameAttribute("prefUser"),
					WithEmailAttribute("email"),
					WithEmailVerifiedAttribute("emailVerified"),
					WithPhoneAttribute("phone"),
					WithPhoneVerifiedAttribute("phoneVerified"),
					WithPreferredLanguageAttribute("prefLang"),
					WithAvatarURLAttribute("avatar"),
					WithProfileAttribute("profile"),
				},
			},
			want: want{
				name:                       "ldap",
				startTls:                   false,
				linkingAllowed:             true,
				creationAllowed:            true,
				autoCreation:               true,
				autoUpdate:                 true,
				idAttribute:                "id",
				firstNameAttribute:         "first",
				lastNameAttribute:          "last",
				displayNameAttribute:       "display",
				nickNameAttribute:          "nick",
				preferredUsernameAttribute: "prefUser",
				emailAttribute:             "email",
				emailVerifiedAttribute:     "emailVerified",
				phoneAttribute:             "phone",
				phoneVerifiedAttribute:     "phoneVerified",
				preferredLanguageAttribute: "prefLang",
				avatarURLAttribute:         "avatar",
				profileAttribute:           "profile",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			provider := New(
				tt.fields.name,
				tt.fields.servers,
				tt.fields.baseDN,
				tt.fields.bindDN,
				tt.fields.bindPassword,
				tt.fields.userBase,
				tt.fields.userObjectClasses,
				tt.fields.userFilters,
				tt.fields.timeout,
				tt.fields.loginUrl,
				tt.fields.opts...,
			)

			a.Equal(tt.want.name, provider.Name())
			a.Equal(tt.want.startTls, provider.startTLS)
			a.Equal(tt.want.linkingAllowed, provider.IsLinkingAllowed())
			a.Equal(tt.want.creationAllowed, provider.IsCreationAllowed())
			a.Equal(tt.want.autoCreation, provider.IsAutoCreation())
			a.Equal(tt.want.autoUpdate, provider.IsAutoUpdate())

			a.Equal(tt.want.idAttribute, provider.idAttribute)
			a.Equal(tt.want.firstNameAttribute, provider.firstNameAttribute)
			a.Equal(tt.want.lastNameAttribute, provider.lastNameAttribute)
			a.Equal(tt.want.displayNameAttribute, provider.displayNameAttribute)
			a.Equal(tt.want.nickNameAttribute, provider.nickNameAttribute)
			a.Equal(tt.want.preferredUsernameAttribute, provider.preferredUsernameAttribute)
			a.Equal(tt.want.emailAttribute, provider.emailAttribute)
			a.Equal(tt.want.emailVerifiedAttribute, provider.emailVerifiedAttribute)
			a.Equal(tt.want.phoneAttribute, provider.phoneAttribute)
			a.Equal(tt.want.phoneVerifiedAttribute, provider.phoneVerifiedAttribute)
			a.Equal(tt.want.preferredLanguageAttribute, provider.preferredLanguageAttribute)
			a.Equal(tt.want.avatarURLAttribute, provider.avatarURLAttribute)
			a.Equal(tt.want.profileAttribute, provider.profileAttribute)
		})
	}
}
