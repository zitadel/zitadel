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

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		config       *oauth2.Config
		userEndpoint string
		userMapper   func() UserInfoMapper
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "successful auth",
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
			},
			want: &Session{AuthURL: "https://oauth2.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			a := assert.New(t)

			provider, err := New(tt.fields.config, tt.fields.userEndpoint, tt.fields.userMapper)
			a.NoError(err)

			session, err := provider.BeginAuth("testState")
			a.NoError(err)

			a.Equal(tt.want.GetAuthURL(), session.GetAuthURL())
		})
	}
}

func TestProvider_FetchUser(t *testing.T) {
	type fields struct {
		config       *oauth2.Config
		userEndpoint string
		httpMock     func(issuer string)
		userMapper   func() UserInfoMapper
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
			},
			args: args{
				&Session{
					AuthURL: "https://oauth2.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
					Tokens:  nil,
				},
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
			},
			args: args{
				&Session{
					AuthURL: "https://oauth2.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
					Tokens: &oidc.Tokens{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
							TokenType:   oidc.BearerToken,
						},
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
				userMapper: func() UserInfoMapper {
					return &UserMapper{
						ID: "userID",
					}
				},
			},
			args: args{
				&Session{
					AuthURL: "https://issuer.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
					Tokens: &oidc.Tokens{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
							TokenType:   oidc.BearerToken,
						},
					},
				},
			},
			want: want{
				user: idp.User{
					ID: "id",
					//Name:      "firstname lastname",
					//NickName:  "nickname",
					//Email:     "email",
					//AvatarURL: "picture",
					//FirstName: "firstname",
					//LastName:  "lastname",
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
					return &UserMapper{
						ID: "userID",
					}
				},
			},
			args: args{
				&Session{
					AuthURL: "https://issuer.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState",
					Tokens:  nil,
					Code:    "code",
				},
			},
			want: want{
				user: idp.User{
					ID: "id",
					//Name:      "firstname lastname",
					//NickName:  "nickname",
					//Email:     "email",
					//AvatarURL: "picture",
					//FirstName: "firstname",
					//LastName:  "lastname",
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

			provider, err := New(tt.fields.config, tt.fields.userEndpoint, tt.fields.userMapper)
			a.NoError(err)

			user, err := provider.FetchUser(tt.args.session)
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
