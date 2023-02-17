package jwt

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp"
)

func TestSession_FetchUser(t *testing.T) {
	type fields struct {
		authURL string
		tokens  *oidc.Tokens
	}
	type want struct {
		err               func(error) bool
		user              idp.User
		id                string
		firstName         string
		lastName          string
		displayName       string
		nickName          string
		preferredUsername string
		email             string
		isEmailVerified   bool
		phone             string
		isPhoneVerified   bool
		preferredLanguage language.Tag
		avatarURL         string
		profile           string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name:   "no tokens",
			fields: fields{},
			want: want{
				err: func(err error) bool {
					return errors.Is(err, ErrNoTokens)
				},
			},
		},
		{
			name: "successful fetch",
			fields: fields{
				authURL: "https://auth.com/jwt?authRequestID=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{},
					IDTokenClaims: func() oidc.IDTokenClaims {
						claims := oidc.EmptyIDTokenClaims()
						userinfo := oidc.NewUserInfo()
						userinfo.SetSubject("sub")
						userinfo.SetPicture("picture")
						userinfo.SetName("firstname lastname")
						userinfo.SetEmail("email", true)
						userinfo.SetGivenName("firstname")
						userinfo.SetFamilyName("lastname")
						userinfo.SetNickname("nickname")
						userinfo.SetPreferredUsername("username")
						userinfo.SetProfile("profile")
						userinfo.SetPhone("phone", true)
						userinfo.SetLocale(language.English)
						claims.SetUserinfo(userinfo)
						return claims
					}(),
				},
			},
			want: want{
				user: &User{
					IDTokenClaims: func() oidc.IDTokenClaims {
						claims := oidc.EmptyIDTokenClaims()
						userinfo := oidc.NewUserInfo()
						userinfo.SetSubject("sub")
						userinfo.SetPicture("picture")
						userinfo.SetName("firstname lastname")
						userinfo.SetEmail("email", true)
						userinfo.SetGivenName("firstname")
						userinfo.SetFamilyName("lastname")
						userinfo.SetNickname("nickname")
						userinfo.SetPreferredUsername("username")
						userinfo.SetProfile("profile")
						userinfo.SetPhone("phone", true)
						userinfo.SetLocale(language.English)
						claims.SetUserinfo(userinfo)
						return claims
					}(),
				},
				id:                "sub",
				firstName:         "firstname",
				lastName:          "lastname",
				displayName:       "firstname lastname",
				nickName:          "nickname",
				preferredUsername: "username",
				email:             "email",
				isEmailVerified:   true,
				phone:             "phone",
				isPhoneVerified:   true,
				preferredLanguage: language.English,
				avatarURL:         "picture",
				profile:           "profile",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			session := &Session{
				AuthURL: tt.fields.authURL,
				Tokens:  tt.fields.tokens,
			}

			user, err := session.FetchUser(context.Background())
			if tt.want.err != nil && !tt.want.err(err) {
				a.Fail("invalid error", err)
			}
			if tt.want.err == nil {
				a.NoError(err)
				a.Equal(tt.want.user, user)
				a.Equal(tt.want.id, user.GetID())
				a.Equal(tt.want.firstName, user.GetFirstName())
				a.Equal(tt.want.lastName, user.GetLastName())
				a.Equal(tt.want.displayName, user.GetDisplayName())
				a.Equal(tt.want.nickName, user.GetNickname())
				a.Equal(tt.want.preferredUsername, user.GetPreferredUsername())
				a.Equal(tt.want.email, user.GetEmail())
				a.Equal(tt.want.isEmailVerified, user.IsEmailVerified())
				a.Equal(tt.want.phone, user.GetPhone())
				a.Equal(tt.want.isPhoneVerified, user.IsPhoneVerified())
				a.Equal(tt.want.preferredLanguage, user.GetPreferredLanguage())
				a.Equal(tt.want.avatarURL, user.GetAvatarURL())
				a.Equal(tt.want.profile, user.GetProfile())
			}
		})
	}
}
