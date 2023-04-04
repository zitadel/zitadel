package ldap

import (
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestProvider_objectClassesToSearchQuery(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		want   string
	}{
		{
			name:   "zero",
			fields: []string{},
			want:   "",
		},
		{
			name:   "one",
			fields: []string{"test"},
			want:   "(objectClass=test)",
		},
		{
			name:   "three",
			fields: []string{"test1", "test2", "test3"},
			want:   "(objectClass=test1)(objectClass=test2)(objectClass=test3)",
		},
		{
			name:   "five",
			fields: []string{"test1", "test2", "test3", "test4", "test5"},
			want:   "(objectClass=test1)(objectClass=test2)(objectClass=test3)(objectClass=test4)(objectClass=test5)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			a.Equal(tt.want, objectClassesToSearchQuery(tt.fields))
		})
	}
}

func TestProvider_userFiltersToSearchQuery(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		username string
		want     string
	}{
		{
			name:     "zero",
			fields:   []string{},
			username: "user",
			want:     "",
		},
		{
			name:     "one",
			fields:   []string{"test"},
			username: "user",
			want:     "(test=user)",
		},
		{
			name:     "three",
			fields:   []string{"test1", "test2", "test3"},
			username: "user",
			want:     "(test1=user)(test2=user)(test3=user)",
		},
		{
			name:     "five",
			fields:   []string{"test1", "test2", "test3", "test4", "test5"},
			username: "user",
			want:     "(test1=user)(test2=user)(test3=user)(test4=user)(test5=user)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			a.Equal(tt.want, userFiltersToSearchQuery(tt.fields, tt.username))
		})
	}
}

func TestProvider_queriesAndToSearchQuery(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		want   string
	}{
		{
			name:   "zero",
			fields: []string{},
			want:   "",
		},
		{
			name:   "one",
			fields: []string{"(test)"},
			want:   "(test)",
		},
		{
			name:   "three",
			fields: []string{"(test1)", "(test2)", "(test3)"},
			want:   "(&(test1)(test2)(test3))",
		},
		{
			name:   "five",
			fields: []string{"(test1)", "(test2)", "(test3)", "(test4)", "(test5)"},
			want:   "(&(test1)(test2)(test3)(test4)(test5))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			a.Equal(tt.want, queriesAndToSearchQuery(tt.fields...))
		})
	}
}

func TestProvider_queriesOrToSearchQuery(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		want   string
	}{
		{
			name:   "zero",
			fields: []string{},
			want:   "",
		},
		{
			name:   "one",
			fields: []string{"(test)"},
			want:   "(test)",
		},
		{
			name:   "three",
			fields: []string{"(test1)", "(test2)", "(test3)"},
			want:   "(|(test1)(test2)(test3))",
		},
		{
			name:   "five",
			fields: []string{"(test1)", "(test2)", "(test3)", "(test4)", "(test5)"},
			want:   "(|(test1)(test2)(test3)(test4)(test5))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			a.Equal(tt.want, queriesOrToSearchQuery(tt.fields...))
		})
	}
}

func TestProvider_mapLDAPEntryToUser(t *testing.T) {
	type fields struct {
		user                       *ldap.Entry
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
	type want struct {
		user *User
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
			fields: fields{
				user: &ldap.Entry{
					Attributes: []*ldap.EntryAttribute{
						{Name: "id", Values: []string{"id"}},
						{Name: "first", Values: []string{"first"}},
						{Name: "last", Values: []string{"last"}},
						{Name: "display", Values: []string{"display"}},
						{Name: "nick", Values: []string{"nick"}},
						{Name: "preferred", Values: []string{"preferred"}},
						{Name: "email", Values: []string{"email"}},
						{Name: "emailVerified", Values: []string{"false"}},
						{Name: "phone", Values: []string{"phone"}},
						{Name: "phoneVerified", Values: []string{"false"}},
						{Name: "lang", Values: []string{"und"}},
						{Name: "avatar", Values: []string{"avatar"}},
						{Name: "profile", Values: []string{"profile"}},
					},
				},
				idAttribute:                "",
				firstNameAttribute:         "",
				lastNameAttribute:          "",
				displayNameAttribute:       "",
				nickNameAttribute:          "",
				preferredUsernameAttribute: "",
				emailAttribute:             "",
				emailVerifiedAttribute:     "",
				phoneAttribute:             "",
				phoneVerifiedAttribute:     "",
				preferredLanguageAttribute: "",
				avatarURLAttribute:         "",
				profileAttribute:           "",
			},
			want: want{
				user: &User{
					id:                "",
					firstName:         "",
					lastName:          "",
					displayName:       "",
					nickName:          "",
					preferredUsername: "",
					email:             "",
					emailVerified:     false,
					phone:             "",
					phoneVerified:     false,
					preferredLanguage: language.Tag{},
					avatarURL:         "",
					profile:           "",
				},
			},
		},
		{
			name: "failed parse emailVerified",
			fields: fields{
				user: &ldap.Entry{
					Attributes: []*ldap.EntryAttribute{
						{Name: "id", Values: []string{"id"}},
						{Name: "first", Values: []string{"first"}},
						{Name: "last", Values: []string{"last"}},
						{Name: "display", Values: []string{"display"}},
						{Name: "nick", Values: []string{"nick"}},
						{Name: "preferred", Values: []string{"preferred"}},
						{Name: "email", Values: []string{"email"}},
						{Name: "emailVerified", Values: []string{"failure"}},
						{Name: "phone", Values: []string{"phone"}},
						{Name: "phoneVerified", Values: []string{"false"}},
						{Name: "lang", Values: []string{"und"}},
						{Name: "avatar", Values: []string{"avatar"}},
						{Name: "profile", Values: []string{"profile"}},
					},
				},
				idAttribute:                "id",
				firstNameAttribute:         "first",
				lastNameAttribute:          "last",
				displayNameAttribute:       "display",
				nickNameAttribute:          "nick",
				preferredUsernameAttribute: "preferred",
				emailAttribute:             "email",
				emailVerifiedAttribute:     "emailVerified",
				phoneAttribute:             "phone",
				phoneVerifiedAttribute:     "phoneVerified",
				preferredLanguageAttribute: "lang",
				avatarURLAttribute:         "avatar",
				profileAttribute:           "profile",
			},
			want: want{
				err: func(err error) bool {
					return err != nil
				},
			},
		},
		{
			name: "failed parse phoneVerified",
			fields: fields{
				user: &ldap.Entry{
					Attributes: []*ldap.EntryAttribute{
						{Name: "id", Values: []string{"id"}},
						{Name: "first", Values: []string{"first"}},
						{Name: "last", Values: []string{"last"}},
						{Name: "display", Values: []string{"display"}},
						{Name: "nick", Values: []string{"nick"}},
						{Name: "preferred", Values: []string{"preferred"}},
						{Name: "email", Values: []string{"email"}},
						{Name: "emailVerified", Values: []string{"false"}},
						{Name: "phone", Values: []string{"phone"}},
						{Name: "phoneVerified", Values: []string{"failure"}},
						{Name: "lang", Values: []string{"und"}},
						{Name: "avatar", Values: []string{"avatar"}},
						{Name: "profile", Values: []string{"profile"}},
					},
				},
				idAttribute:                "id",
				firstNameAttribute:         "first",
				lastNameAttribute:          "last",
				displayNameAttribute:       "display",
				nickNameAttribute:          "nick",
				preferredUsernameAttribute: "preferred",
				emailAttribute:             "email",
				emailVerifiedAttribute:     "emailVerified",
				phoneAttribute:             "phone",
				phoneVerifiedAttribute:     "phoneVerified",
				preferredLanguageAttribute: "lang",
				avatarURLAttribute:         "avatar",
				profileAttribute:           "profile",
			},
			want: want{
				err: func(err error) bool {
					return err != nil
				},
			},
		},
		{
			name: "full user",
			fields: fields{
				user: &ldap.Entry{
					Attributes: []*ldap.EntryAttribute{
						{Name: "id", Values: []string{"id"}},
						{Name: "first", Values: []string{"first"}},
						{Name: "last", Values: []string{"last"}},
						{Name: "display", Values: []string{"display"}},
						{Name: "nick", Values: []string{"nick"}},
						{Name: "preferred", Values: []string{"preferred"}},
						{Name: "email", Values: []string{"email"}},
						{Name: "emailVerified", Values: []string{"false"}},
						{Name: "phone", Values: []string{"phone"}},
						{Name: "phoneVerified", Values: []string{"false"}},
						{Name: "lang", Values: []string{"und"}},
						{Name: "avatar", Values: []string{"avatar"}},
						{Name: "profile", Values: []string{"profile"}},
					},
				},
				idAttribute:                "id",
				firstNameAttribute:         "first",
				lastNameAttribute:          "last",
				displayNameAttribute:       "display",
				nickNameAttribute:          "nick",
				preferredUsernameAttribute: "preferred",
				emailAttribute:             "email",
				emailVerifiedAttribute:     "emailVerified",
				phoneAttribute:             "phone",
				phoneVerifiedAttribute:     "phoneVerified",
				preferredLanguageAttribute: "lang",
				avatarURLAttribute:         "avatar",
				profileAttribute:           "profile",
			},
			want: want{
				user: &User{
					id:                "id",
					firstName:         "first",
					lastName:          "last",
					displayName:       "display",
					nickName:          "nick",
					preferredUsername: "preferred",
					email:             "email",
					emailVerified:     false,
					phone:             "phone",
					phoneVerified:     false,
					preferredLanguage: language.Make("und"),
					avatarURL:         "avatar",
					profile:           "profile",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapLDAPEntryToUser(
				tt.fields.user,
				tt.fields.idAttribute,
				tt.fields.firstNameAttribute,
				tt.fields.lastNameAttribute,
				tt.fields.displayNameAttribute,
				tt.fields.nickNameAttribute,
				tt.fields.preferredUsernameAttribute,
				tt.fields.emailAttribute,
				tt.fields.emailVerifiedAttribute,
				tt.fields.phoneAttribute,
				tt.fields.phoneVerifiedAttribute,
				tt.fields.preferredLanguageAttribute,
				tt.fields.avatarURLAttribute,
				tt.fields.profileAttribute,
			)
			if tt.want.err == nil {
				assert.NoError(t, err)
			}
			if tt.want.err != nil && !tt.want.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.want.err == nil {
				assert.Equal(t, tt.want.user, got)
			}
		})
	}
}
