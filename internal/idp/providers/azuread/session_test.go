package azuread

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
		name         string
		clientID     string
		clientSecret string
		redirectURI  string
		scopes       []string
		httpMock     func()
		options      []ProviderOptions
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
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				httpMock: func() {
					gock.New("https://graph.microsoft.com").
						Get("/v1.0/me").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
				tokens:  nil,
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
					gock.New("https://graph.microsoft.com").
						Get("/v1.0/me").
						Reply(http.StatusInternalServerError)
				},
				authURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
					IDTokenClaims: oidc.NewIDTokenClaims(
						"https://login.microsoftonline.com/consumers/oauth2/v2.0",
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
					gock.New("https://graph.microsoft.com").
						Get("/v1.0/me").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
					IDTokenClaims: oidc.NewIDTokenClaims(
						"https://login.microsoftonline.com/consumers/oauth2/v2.0",
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
				user: &User{
					ID:                "id",
					BusinessPhones:    []domain.PhoneNumber{"phone1", "phone2"},
					DisplayName:       "firstname lastname",
					FirstName:         "firstname",
					JobTitle:          "title",
					Email:             "email",
					MobilePhone:       "mobile",
					OfficeLocation:    "office",
					PreferredLanguage: "en",
					LastName:          "lastname",
					UserPrincipalName: "username",
					isEmailVerified:   false,
				},
				id:                "id",
				firstName:         "firstname",
				lastName:          "lastname",
				displayName:       "firstname lastname",
				nickName:          "",
				preferredUsername: "username",
				email:             "email",
				isEmailVerified:   false,
				phone:             "",
				isPhoneVerified:   false,
				preferredLanguage: language.English,
				profile:           "",
			},
		},
		{
			name: "successful fetch with email verified",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				options: []ProviderOptions{
					WithEmailVerified(),
				},
				httpMock: func() {
					gock.New("https://graph.microsoft.com").
						Get("/v1.0/me").
						Reply(200).
						JSON(userinfo())
				},
				authURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
				tokens: &oidc.Tokens{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   oidc.BearerToken,
					},
					IDTokenClaims: oidc.NewIDTokenClaims(
						"https://login.microsoftonline.com/consumers/oauth2/v2.0",
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
				user: &User{
					ID:                "id",
					BusinessPhones:    []domain.PhoneNumber{"phone1", "phone2"},
					DisplayName:       "firstname lastname",
					FirstName:         "firstname",
					JobTitle:          "title",
					Email:             "email",
					MobilePhone:       "mobile",
					OfficeLocation:    "office",
					PreferredLanguage: "en",
					LastName:          "lastname",
					UserPrincipalName: "username",
					isEmailVerified:   true,
				},
				id:                "id",
				firstName:         "firstname",
				lastName:          "lastname",
				displayName:       "firstname lastname",
				nickName:          "",
				preferredUsername: "username",
				email:             "email",
				isEmailVerified:   true,
				phone:             "",
				isPhoneVerified:   false,
				preferredLanguage: language.English,
				profile:           "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.fields.httpMock()
			a := assert.New(t)

			provider, err := New(tt.fields.name, tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes, tt.fields.options...)
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
		ID:                "id",
		BusinessPhones:    []domain.PhoneNumber{"phone1", "phone2"},
		DisplayName:       "firstname lastname",
		FirstName:         "firstname",
		JobTitle:          "title",
		Email:             "email",
		MobilePhone:       "mobile",
		OfficeLocation:    "office",
		PreferredLanguage: "en",
		LastName:          "lastname",
		UserPrincipalName: "username",
	}
}
