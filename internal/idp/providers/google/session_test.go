package google

import (
	"errors"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp"
	oidc2 "github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

func TestSession_FetchUser(t *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		redirectURI  string
		httpMock     func()
		authURL      string
		code         string
		tokens       *oidc.Tokens
	}
	type args struct {
		session idp.Session
	}
	type want struct {
		user idp.User
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "unauthenticated session, error",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				httpMock: func() {
					gock.New("https://openidconnect.googleapis.com").
						Get("/v1/userinfo").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://accounts.google.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
				tokens:  nil,
			},
			args: args{
				&oidc2.Session{},
			},
			want: want{
				err: oidc2.ErrCodeMissing,
			},
		},
		{
			name: "userinfo error",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				httpMock: func() {
					gock.New("https://openidconnect.googleapis.com").
						Get("/v1/userinfo").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://accounts.google.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
					IDTokenClaims: oidc.NewIDTokenClaims(
						issuer,
						"sub2",
						[]string{"clientID"},
						time.Now().Add(1*time.Hour),
						time.Now().Add(-1*time.Second),
						"nonce",
						"",
						nil,
						"clientID",
						0,
					),
				},
			},
			args: args{
				&oidc2.Session{},
			},
			want: want{
				err: rp.ErrUserInfoSubNotMatching,
			},
		},
		{
			name: "successful fetch",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				httpMock: func() {
					gock.New("https://openidconnect.googleapis.com").
						Get("/v1/userinfo").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://accounts.google.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
					IDTokenClaims: oidc.NewIDTokenClaims(
						issuer,
						"sub",
						[]string{"clientID"},
						time.Now().Add(1*time.Hour),
						time.Now().Add(-1*time.Second),
						"nonce",
						"",
						nil,
						"clientID",
						0,
					),
				},
			},
			args: args{
				&oidc2.Session{},
			},
			want: want{
				user: idp.User{
					ID:                "sub",
					FirstName:         "firstname",
					LastName:          "lastname",
					DisplayName:       "firstname lastname",
					NickName:          "nickname",
					PreferredUsername: "username",
					Email:             "email",
					IsEmailVerified:   true,
					Phone:             "phone",
					IsPhoneVerified:   true,
					PreferredLanguage: language.English,
					AvatarURL:         "picture",
					Profile:           "profile",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.fields.httpMock()
			a := assert.New(t)

			// call the real discovery endpoint
			gock.New(issuer).Get(oidc.DiscoveryEndpoint).EnableNetworking()
			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI)
			a.NoError(err)

			session := &oidc2.Session{
				Provider: provider.Provider,
				AuthURL:  tt.fields.authURL,
				Code:     tt.fields.code,
				Tokens:   tt.fields.tokens,
			}

			user, err := session.FetchUser()
			if tt.want.err != nil && !errors.Is(err, tt.want.err) {
				a.Fail("invalid error", "expected %v, got %v", tt.want.err, err)
			}
			if tt.want.err == nil {
				a.NoError(err)
				a.Equal(tt.want.user, user)
			}
		})
	}
}

func userinfo() oidc.UserInfoSetter {
	info := oidc.NewUserInfo()
	info.SetSubject("sub")
	info.SetGivenName("firstname")
	info.SetFamilyName("lastname")
	info.SetName("firstname lastname")
	info.SetNickname("nickname")
	info.SetPreferredUsername("username")
	info.SetEmail("email", true)
	info.SetPhone("phone", true)
	info.SetLocale(language.English)
	info.SetPicture("picture")
	info.SetProfile("profile")
	return info
}
