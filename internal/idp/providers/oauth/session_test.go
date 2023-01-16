package oauth

import (
	"errors"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

func TestProvider_FetchUser(t *testing.T) {
	type fields struct {
		config       *oauth2.Config
		name         string
		userEndpoint string
		httpMock     func(issuer string)
		userMapper   func() UserInfoMapper
		authURL      string
		code         string
		tokens       *oidc.Tokens
	}
	type args struct {
		session idp.Session
	}
	type want struct {
		user idp.User
		err  func(error) bool
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
			args: args{
				&Session{},
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
				userMapper: func() UserInfoMapper {
					return &UserMapper{
						ID: "userID",
					}
				},
				authURL: "https://oauth2.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
				},
			},
			args: args{
				&Session{},
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
				userMapper: func() UserInfoMapper {
					return &UserMapper{}
				},
				authURL: "https://issuer.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
				},
			},
			args: args{
				&Session{},
			},
			want: want{
				user: idp.User{
					RawData: map[string]interface{}{
						"userID": "id",
						"custom": "claim",
					},
				},
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
				userMapper: func() UserInfoMapper {
					return &UserMapper{}
				},
				authURL: "https://issuer.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
				tokens:  nil,
				code:    "code",
			},
			args: args{
				&Session{},
			},
			want: want{
				user: idp.User{
					RawData: map[string]interface{}{
						"userID": "id",
						"custom": "claim",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.fields.httpMock("https://oauth2.com")
			a := assert.New(t)

			provider, err := New(tt.fields.config, tt.fields.name, tt.fields.userEndpoint, tt.fields.userMapper)
			a.NoError(err)

			session := &Session{
				AuthURL:  tt.fields.authURL,
				Code:     tt.fields.code,
				Tokens:   tt.fields.tokens,
				Provider: provider,
			}

			user, err := session.FetchUser()
			if tt.want.err != nil && !tt.want.err(err) {
				a.Fail("invalid error", err)
			}
			if tt.want.err == nil {
				a.NoError(err)
				a.Equal(tt.want.user, user)
			}
		})
	}
}
