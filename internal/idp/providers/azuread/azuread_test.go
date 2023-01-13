package azuread

import (
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
	oidc2 "github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		redirectURI  string
		opts         ProviderOpts
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "default common tenant",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
			},
			want: &oidc2.Session{
				AuthURL: "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
			},
		},
		{
			name: "tenant",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				opts: ProviderOpts{
					Tenant: ConsumersTenant,
				},
			},
			want: &oidc2.Session{
				AuthURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.opts)
			a.NoError(err)

			session, err := provider.BeginAuth("testState")
			a.NoError(err)

			a.Equal(tt.want.GetAuthURL(), session.GetAuthURL())
		})
	}
}

func TestProvider_FetchUser(t *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		redirectURI  string
		httpMock     func()
		opts         ProviderOpts
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
					gock.New("https://graph.microsoft.com").
						Get("/oidc/userinfo").
						Reply(200).
						JSON(userinfo())
				},
			},
			args: args{
				&oauth.Session{
					AuthURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
					Tokens:  nil,
				},
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
						Get("/oidc/userinfo").
						Reply(http.StatusInternalServerError)
				},
			},
			args: args{
				&oauth.Session{
					AuthURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
					Tokens: &oidc.Tokens{
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
						Get("/oidc/userinfo").
						Reply(200).
						JSON(userinfo())
				},
			},
			args: args{
				&oauth.Session{
					AuthURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
					Tokens: &oidc.Tokens{
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
			},
			want: want{
				user: idp.User{
					ID:                "sub",
					DisplayName:       "firstname lastname",
					PreferredUsername: "username",
					Email:             "email",
					AvatarURL:         "picture",
					FirstName:         "firstname",
					LastName:          "lastname",
					RawData: &User{
						Sub:               "sub",
						FamilyName:        "lastname",
						GivenName:         "firstname",
						Name:              "firstname lastname",
						PreferredUsername: "username",
						Email:             "email",
						Picture:           "picture",
						isEmailVerified:   false,
					},
				},
			},
		},
		{
			name: "successful fetch with email verified",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				opts: ProviderOpts{
					EmailVerified: true,
				},
				httpMock: func() {
					gock.New("https://graph.microsoft.com").
						Get("/oidc/userinfo").
						Reply(200).
						JSON(userinfo())
				},
			},
			args: args{
				&oauth.Session{
					AuthURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
					Tokens: &oidc.Tokens{
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
			},
			want: want{
				user: idp.User{
					ID:                "sub",
					DisplayName:       "firstname lastname",
					PreferredUsername: "username",
					Email:             "email",
					IsEmailVerified:   true,
					AvatarURL:         "picture",
					FirstName:         "firstname",
					LastName:          "lastname",
					RawData: &User{
						Sub:               "sub",
						FamilyName:        "lastname",
						GivenName:         "firstname",
						Name:              "firstname lastname",
						PreferredUsername: "username",
						Email:             "email",
						Picture:           "picture",
						isEmailVerified:   true,
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
			gock.New("https://login.microsoftonline.com/consumers/oauth2/v2.0").Get(oidc.DiscoveryEndpoint).EnableNetworking()
			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.opts)
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

func userinfo() oidc.UserInfoSetter {
	userinfo := oidc.NewUserInfo()
	userinfo.SetSubject("sub")
	userinfo.SetName("firstname lastname")
	userinfo.SetPreferredUsername("username")
	userinfo.SetNickname("nickname")
	userinfo.SetEmail("email", false) // azure add does not send the email_verified claim
	userinfo.SetPicture("picture")
	userinfo.SetGivenName("firstname")
	userinfo.SetFamilyName("lastname")
	return userinfo
}
