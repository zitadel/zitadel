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

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		redirectURI  string
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "successful auth",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
			},
			want: &oidc2.Session{
				AuthURL: "https://accounts.google.com/o/oauth2/v2/auth?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI)
			a.NoError(err)

			session, err := provider.BeginAuth("testState")
			a.NoError(err)

			//authUrl, err := url.Parse(session.GetAuthURL())
			//a.NoError(err)
			//
			a.Equal(tt.want.GetAuthURL(), session.GetAuthURL())
			//a.Equal("/authorize", authUrl.Path)
			//a.Equal("clientID", authUrl.Query().Get("client_id"))
			//a.Equal("testState", authUrl.Query().Get("state"))
			//a.Equal("redirectURI", authUrl.Query().Get("redirect_uri"))
			//a.Equal("openid", authUrl.Query().Get("scope"))
			//
			//if !tt.wantErr(t, err, fmt.Sprintf("BeginAuth(%v)", tt.fields.state)) {
			//	return
			//}
			//assert.Equalf(t, tt.want, got, "BeginAuth(%v)", tt.args.state)
		})
	}
}

func TestProvider_FetchUser(t *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		redirectURI  string
		httpMock     func()
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
			},
			args: args{
				&oidc2.Session{
					AuthURL: "https://accounts.google.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
					Tokens:  nil,
				},
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
			},
			args: args{
				&oidc2.Session{
					AuthURL: "https://accounts.google.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
					Tokens: &oidc.Tokens{
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
			},
			args: args{
				&oidc2.Session{
					AuthURL: "https://accounts.google.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
					Tokens: &oidc.Tokens{
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

			user, err := provider.FetchUser(tt.args.session)
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
