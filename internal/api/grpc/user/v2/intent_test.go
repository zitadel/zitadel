package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
)

// fakeIDPUser is a minimal idp.User implementation for exercising the brokered
// onboarding fallbacks without depending on a concrete provider mapper.
type fakeIDPUser struct {
	id                string
	firstName         string
	lastName          string
	displayName       string
	nickname          string
	preferredUsername string
	email             domain.EmailAddress
	phone             domain.PhoneNumber
}

func (u fakeIDPUser) GetID() string                      { return u.id }
func (u fakeIDPUser) GetFirstName() string               { return u.firstName }
func (u fakeIDPUser) GetLastName() string                { return u.lastName }
func (u fakeIDPUser) GetDisplayName() string             { return u.displayName }
func (u fakeIDPUser) GetNickname() string                { return u.nickname }
func (u fakeIDPUser) GetPreferredUsername() string       { return u.preferredUsername }
func (u fakeIDPUser) GetEmail() domain.EmailAddress      { return u.email }
func (u fakeIDPUser) IsEmailVerified() bool              { return false }
func (u fakeIDPUser) GetPhone() domain.PhoneNumber       { return u.phone }
func (u fakeIDPUser) IsPhoneVerified() bool              { return false }
func (u fakeIDPUser) GetPreferredLanguage() language.Tag { return language.Und }
func (u fakeIDPUser) GetAvatarURL() string               { return "" }
func (u fakeIDPUser) GetProfile() string                 { return "" }

func Test_idpUserFallbackUsername(t *testing.T) {
	tests := []struct {
		name string
		user fakeIDPUser
		want string
	}{
		{
			name: "preferred username wins",
			user: fakeIDPUser{preferredUsername: "minnie", email: "minnie@example.com", id: "sub-1"},
			want: "minnie",
		},
		{
			name: "falls back to email",
			user: fakeIDPUser{email: "brightsparc@gmail.com", id: "sub-1"},
			want: "brightsparc@gmail.com",
		},
		{
			name: "falls back to external subject",
			user: fakeIDPUser{id: "google-oauth2|118092360802169819584"},
			want: "google-oauth2|118092360802169819584",
		},
		{
			name: "whitespace preferred username is ignored",
			user: fakeIDPUser{preferredUsername: "   ", email: "user@example.com"},
			want: "user@example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, idpUserFallbackUsername(tt.user))
		})
	}
}

func Test_idpUserFallbackNames(t *testing.T) {
	tests := []struct {
		name       string
		user       fakeIDPUser
		wantGiven  string
		wantFamily string
	}{
		{
			name:       "provided names win",
			user:       fakeIDPUser{firstName: "Julian", lastName: "Bright", displayName: "ignored me"},
			wantGiven:  "Julian",
			wantFamily: "Bright",
		},
		{
			name:       "derive both from display name",
			user:       fakeIDPUser{displayName: "Julian Bright", email: "brightsparc@gmail.com"},
			wantGiven:  "Julian",
			wantFamily: "Bright",
		},
		{
			name:       "single-token display name reuses given as family",
			user:       fakeIDPUser{displayName: "Julian"},
			wantGiven:  "Julian",
			wantFamily: "Julian",
		},
		{
			name:       "derive from email local part when no names or display name",
			user:       fakeIDPUser{email: "brightsparc@gmail.com"},
			wantGiven:  "brightsparc",
			wantFamily: "brightsparc",
		},
		{
			name:       "missing family is derived from display name",
			user:       fakeIDPUser{firstName: "Julian", displayName: "Julian Bright"},
			wantGiven:  "Julian",
			wantFamily: "Bright",
		},
		{
			name:       "email-only profile still yields non-empty names",
			user:       fakeIDPUser{id: "sub-only"},
			wantGiven:  "sub-only",
			wantFamily: "sub-only",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			given, family := idpUserFallbackNames(tt.user)
			assert.Equal(t, tt.wantGiven, given)
			assert.Equal(t, tt.wantFamily, family)
			assert.NotEmpty(t, given)
			assert.NotEmpty(t, family)
		})
	}
}

// Test_idpUserToAddHumanUser_PartialProfiles mirrors the reported Auth0 and
// Supabase brokered-onboarding cases: the request handed to AddHumanUser must
// satisfy the non-empty given_name/family_name and IDPLink.user_name proto
// constraints even when the upstream IdP omits those claims.
func Test_idpUserToAddHumanUser_PartialProfiles(t *testing.T) {
	t.Run("auth0 without preferred_username", func(t *testing.T) {
		idpUser := fakeIDPUser{
			id:          "google-oauth2|118092360802169819584",
			email:       "brightsparc@gmail.com",
			displayName: "Julian Bright",
			firstName:   "Julian",
			lastName:    "Bright",
		}

		req := idpUserToAddHumanUser(idpUser, "idpID")

		assert.Equal(t, "Julian", req.GetProfile().GetGivenName())
		assert.Equal(t, "Bright", req.GetProfile().GetFamilyName())
		assert.NotEmpty(t, req.GetIdpLinks()[0].GetUserName())
		assert.Equal(t, "brightsparc@gmail.com", req.GetIdpLinks()[0].GetUserName())
	})

	t.Run("supabase email only", func(t *testing.T) {
		idpUser := fakeIDPUser{
			id:    "sub-123",
			email: "brightsparc@gmail.com",
		}

		req := idpUserToAddHumanUser(idpUser, "idpID")

		assert.NotEmpty(t, req.GetProfile().GetGivenName())
		assert.NotEmpty(t, req.GetProfile().GetFamilyName())
		assert.Equal(t, "brightsparc", req.GetProfile().GetGivenName())
		assert.Equal(t, "brightsparc@gmail.com", req.GetIdpLinks()[0].GetUserName())
	})
}
