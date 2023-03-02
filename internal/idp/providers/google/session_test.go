package google

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	openid "github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

func TestSession_FetchUser(t *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		redirectURI  string
		scopes       []string
		httpMock     func()
		authURL      string
		code         string
		tokens       *openid.Tokens
	}
	type want struct {
		err               error
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
		hostedDomain      string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "unauthenticated session, error",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				httpMock: func() {
					gock.New("https://openidconnect.googleapis.com").
						Get("/v1/userinfo").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://accounts.google.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
				tokens:  nil,
			},
			want: want{
				err: oidc.ErrCodeMissing,
			},
		},
		{
			name: "userinfo error",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				httpMock: func() {
					gock.New("https://openidconnect.googleapis.com").
						Get("/v1/userinfo").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://accounts.google.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
				tokens: &openid.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   openid.BearerToken,
					},
					IDTokenClaims: openid.NewIDTokenClaims(
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
				scopes:       []string{"openid"},
				httpMock: func() {
					gock.New("https://openidconnect.googleapis.com").
						Get("/v1/userinfo").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://accounts.google.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
				tokens: &openid.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   openid.BearerToken,
					},
					IDTokenClaims: openid.NewIDTokenClaims(
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
			want: want{
				id:                "sub",
				firstName:         "firstname",
				lastName:          "lastname",
				displayName:       "firstname lastname",
				nickName:          "",
				preferredUsername: "email",
				email:             "email",
				isEmailVerified:   true,
				phone:             "",
				isPhoneVerified:   false,
				preferredLanguage: language.English,
				avatarURL:         "picture",
				profile:           "",
				hostedDomain:      "hosted domain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.fields.httpMock()
			a := assert.New(t)

			// call the real discovery endpoint
			gock.New(issuer).Get(openid.DiscoveryEndpoint).EnableNetworking()
			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes)
			require.NoError(t, err)

			session := &oidc.Session{
				Provider: provider.Provider,
				AuthURL:  tt.fields.authURL,
				Code:     tt.fields.code,
				Tokens:   tt.fields.tokens,
			}

			user, err := session.FetchUser(context.Background())
			if tt.want.err != nil && !errors.Is(err, tt.want.err) {
				a.Fail("invalid error", "expected %v, got %v", tt.want.err, err)
			}
			if tt.want.err == nil {
				a.NoError(err)
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

func userinfo() openid.UserInfoSetter {
	info := openid.NewUserInfo()
	info.SetSubject("sub")
	info.SetGivenName("firstname")
	info.SetFamilyName("lastname")
	info.SetName("firstname lastname")
	info.SetEmail("email", true)
	info.SetLocale(language.English)
	info.SetPicture("picture")
	info.AppendClaims("hd", "hosted domain")
	return info
}
