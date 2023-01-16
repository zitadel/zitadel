package github

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
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
		scopes       []string
		options      []oauth.ProviderOpts
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
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				httpMock: func() {
					gock.New("https://api.github.com").
						Get("/user").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://github.com/login/oauth/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&state=testState",
				tokens:  nil,
			},
			args: args{
				&oauth.Session{},
			},
			want: want{
				err: func(err error) bool {
					return errors.Is(err, oauth.ErrCodeMissing)
				},
			},
		},
		{
			name: "user error",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				httpMock: func() {
					gock.New("https://api.github.com").
						Get("/user").
						Reply(http.StatusInternalServerError)
				},
				authURL: "https://github.com/login/oauth/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
				},
			},
			args: args{
				&oauth.Session{},
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
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				httpMock: func() {
					gock.New("https://api.github.com").
						Get("/user").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://github.com/login/oauth/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
				},
			},
			args: args{
				&oauth.Session{},
			},
			want: want{
				user: idp.User{
					ID:                "1",
					DisplayName:       "name",
					NickName:          "login",
					PreferredUsername: "login",
					Email:             "email",
					IsEmailVerified:   true,
					AvatarURL:         "avatarURL",
					Profile:           "htmlURL",
					RawData: &User{
						Login:      "login",
						ID:         1,
						AvatarUrl:  "avatarURL",
						GravatarId: "gravatarID",
						Name:       "name",
						Email:      "email",
						HtmlUrl:    "htmlURL",
						CreatedAt:  time.Date(2023, 01, 10, 11, 10, 35, 0, time.Local),
						UpdatedAt:  time.Date(2023, 01, 10, 11, 10, 35, 0, time.Local),
					},
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
			gock.New("https://api.github.com").Get("/user").EnableNetworking()
			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes, tt.fields.options...)
			a.NoError(err)

			session := &oauth.Session{
				AuthURL:  tt.fields.authURL,
				Code:     tt.fields.code,
				Tokens:   tt.fields.tokens,
				Provider: provider.Provider,
			}

			user, err := session.FetchUser(context.Background())
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

func userinfo() *User {
	return &User{
		Login:      "login",
		ID:         1,
		AvatarUrl:  "avatarURL",
		GravatarId: "gravatarID",
		Name:       "name",
		Email:      "email",
		HtmlUrl:    "htmlURL",
		CreatedAt:  time.Date(2023, 01, 10, 11, 10, 35, 0, time.Local),
		UpdatedAt:  time.Date(2023, 01, 10, 11, 10, 35, 0, time.Local),
	}
}
