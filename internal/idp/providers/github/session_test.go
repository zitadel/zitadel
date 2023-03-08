package github

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
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
				user: &User{
					Login:      "login",
					ID:         1,
					AvatarUrl:  "avatarURL",
					GravatarId: "gravatarID",
					Name:       "name",
					Email:      "email",
					HtmlUrl:    "htmlURL",
					CreatedAt:  time.Date(2023, 01, 10, 11, 10, 35, 0, time.UTC),
					UpdatedAt:  time.Date(2023, 01, 10, 11, 10, 35, 0, time.UTC),
				},
				id:                "1",
				firstName:         "",
				lastName:          "",
				displayName:       "name",
				nickName:          "login",
				preferredUsername: "login",
				email:             "email",
				isEmailVerified:   true,
				phone:             "",
				isPhoneVerified:   false,
				preferredLanguage: language.Und,
				avatarURL:         "avatarURL",
				profile:           "htmlURL",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.fields.httpMock()
			a := assert.New(t)

			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes, tt.fields.options...)
			require.NoError(t, err)

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

func userinfo() *User {
	return &User{
		Login:      "login",
		ID:         1,
		AvatarUrl:  "avatarURL",
		GravatarId: "gravatarID",
		Name:       "name",
		Email:      "email",
		HtmlUrl:    "htmlURL",
		CreatedAt:  time.Date(2023, 01, 10, 11, 10, 35, 0, time.UTC),
		UpdatedAt:  time.Date(2023, 01, 10, 11, 10, 35, 0, time.UTC),
	}
}
