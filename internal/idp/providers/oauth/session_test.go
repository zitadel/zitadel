package oauth

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
)

func TestProvider_FetchUser(t *testing.T) {
	type fields struct {
		config       *oauth2.Config
		name         string
		userEndpoint string
		httpMock     func(issuer string)
		userMapper   func() idp.User
		authURL      string
		code         string
		tokens       *oidc.Tokens
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
			name: "unauthenticated session, error",
			fields: fields{
				config: &oauth2.Config{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://oauth2.com/authorize",
						TokenURL: "https://oauth2.com/token",
					},
					RedirectURL: "redirectURI",
					Scopes:      []string{"user"},
				},
				userEndpoint: "https://oauth2.com/user",
				httpMock:     func(issuer string) {},
				authURL:      "https://oauth2.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
				tokens:       nil,
			},
			want: want{
				err: func(err error) bool {
					return errors.Is(err, ErrCodeMissing)
				},
			},
		},
		{
			name: "user error",
			fields: fields{
				config: &oauth2.Config{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://oauth2.com/authorize",
						TokenURL: "https://oauth2.com/token",
					},
					RedirectURL: "redirectURI",
					Scopes:      []string{"user"},
				},
				userEndpoint: "https://oauth2.com/user",
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get("/user").
						Reply(http.StatusInternalServerError)
				},
				userMapper: func() idp.User {
					return NewUserMapper("userID")
				},
				authURL: "https://oauth2.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
				},
			},
			want: want{
				err: func(err error) bool {
					return err.Error() == "http status not ok: 500 Internal Server Error "
				},
			},
		},
		{
			name: "successful fetch",
			fields: fields{
				config: &oauth2.Config{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://oauth2.com/authorize",
						TokenURL: "https://oauth2.com/token",
					},
					RedirectURL: "redirectURI",
					Scopes:      []string{"user"},
				},
				userEndpoint: "https://oauth2.com/user",
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get("/user").
						Reply(200).
						JSON(map[string]interface{}{
							"userID": "id",
							"custom": "claim",
						})
				},
				userMapper: func() idp.User {
					return NewUserMapper("userID")
				},
				authURL: "https://issuer.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
				},
			},
			want: want{
				user: &UserMapper{
					idAttribute: "userID",
					RawInfo: map[string]interface{}{
						"userID": "id",
						"custom": "claim",
					},
				},
				id:                "id",
				firstName:         "",
				lastName:          "",
				displayName:       "",
				nickName:          "",
				preferredUsername: "",
				email:             "",
				isEmailVerified:   false,
				phone:             "",
				isPhoneVerified:   false,
				preferredLanguage: language.Und,
				avatarURL:         "",
				profile:           "",
			},
		},
		{
			name: "successful fetch with code exchange",
			fields: fields{
				config: &oauth2.Config{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://oauth2.com/authorize",
						TokenURL: "https://oauth2.com/token",
					},
					RedirectURL: "redirectURI",
					Scopes:      []string{"user"},
				},
				userEndpoint: "https://oauth2.com/user",
				httpMock: func(issuer string) {
					gock.New(issuer).
						Post("/token").
						BodyString("client_id=clientID&client_secret=clientSecret&code=code&grant_type=authorization_code&redirect_uri=redirectURI").
						Reply(200).
						JSON(&oidc.AccessTokenResponse{
							AccessToken:  "accessToken",
							TokenType:    oidc.BearerToken,
							RefreshToken: "",
							ExpiresIn:    3600,
							IDToken:      "",
							State:        "testState"})
					gock.New(issuer).
						Get("/user").
						Reply(200).
						JSON(map[string]interface{}{
							"userID": "id",
							"custom": "claim",
						})
				},
				userMapper: func() idp.User {
					return NewUserMapper("userID")
				},
				authURL: "https://issuer.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
				tokens:  nil,
				code:    "code",
			},
			want: want{
				user: &UserMapper{
					idAttribute: "userID",
					RawInfo: map[string]interface{}{
						"userID": "id",
						"custom": "claim",
					},
				},
				id:                "id",
				firstName:         "",
				lastName:          "",
				displayName:       "",
				nickName:          "",
				preferredUsername: "",
				email:             "",
				isEmailVerified:   false,
				phone:             "",
				isPhoneVerified:   false,
				preferredLanguage: language.Und,
				avatarURL:         "",
				profile:           "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.fields.httpMock("https://oauth2.com")
			a := assert.New(t)

			provider, err := New(tt.fields.config, tt.fields.name, tt.fields.userEndpoint, tt.fields.userMapper)
			require.NoError(t, err)

			session := &Session{
				AuthURL:  tt.fields.authURL,
				Code:     tt.fields.code,
				Tokens:   tt.fields.tokens,
				Provider: provider,
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
				a.Equal(domain.EmailAddress(tt.want.email), user.GetEmail())
				a.Equal(tt.want.isEmailVerified, user.IsEmailVerified())
				a.Equal(domain.PhoneNumber(tt.want.phone), user.GetPhone())
				a.Equal(tt.want.isPhoneVerified, user.IsPhoneVerified())
				a.Equal(tt.want.preferredLanguage, user.GetPreferredLanguage())
				a.Equal(tt.want.avatarURL, user.GetAvatarURL())
				a.Equal(tt.want.profile, user.GetProfile())
			}
		})
	}
}
