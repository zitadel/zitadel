package ldap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_Options(t *testing.T) {
	type fields struct {
		name                string
		host                string
		baseDN              string
		userObjectClass     string
		userUniqueAttribute string
		admin               string
		password            string
		loginUrl            string
		opts                []ProviderOpts
	}
	type want struct {
		name                       string
		port                       string
		tls                        bool
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
				name:                "ldap",
				host:                "host",
				baseDN:              "base",
				userObjectClass:     "class",
				userUniqueAttribute: "attr",
				admin:               "admin",
				password:            "password",
				loginUrl:            "url",
				opts:                nil,
			},
			want: want{
				name:            "ldap",
				port:            DefaultPort,
				tls:             true,
				linkingAllowed:  false,
				creationAllowed: false,
				autoCreation:    false,
				autoUpdate:      false,
				idAttribute:     "attr",
			},
		},
		{
			name: "all true",
			fields: fields{
				name:                "ldap",
				host:                "host",
				baseDN:              "base",
				userObjectClass:     "class",
				userUniqueAttribute: "attr",
				admin:               "admin",
				password:            "password",
				loginUrl:            "url",
				opts: []ProviderOpts{
					WithLinkingAllowed(),
					WithCreationAllowed(),
					WithAutoCreation(),
					WithAutoUpdate(),
				},
			},
			want: want{
				name:            "ldap",
				port:            DefaultPort,
				tls:             true,
				linkingAllowed:  true,
				creationAllowed: true,
				autoCreation:    true,
				autoUpdate:      true,
				idAttribute:     "attr",
			},
		}, {
			name: "all true, attributes set",
			fields: fields{
				name:                "ldap",
				host:                "host",
				baseDN:              "base",
				userObjectClass:     "class",
				userUniqueAttribute: "attr",
				admin:               "admin",
				password:            "password",
				loginUrl:            "url",
				opts: []ProviderOpts{
					Insecure(),
					WithCustomPort("port"),
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
				port:                       "port",
				tls:                        false,
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
			provider := New(tt.fields.name, tt.fields.host, tt.fields.baseDN, tt.fields.userObjectClass, tt.fields.userUniqueAttribute, tt.fields.admin, tt.fields.password, tt.fields.loginUrl, tt.fields.opts...)

			a.Equal(tt.want.name, provider.Name())
			a.Equal(tt.want.port, provider.port)
			a.Equal(tt.want.tls, provider.tls)
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
