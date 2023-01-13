package jwt

import (
	"errors"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		issuer       string
		jwtEndpoint  string
		keysEndpoint string
		headerName   string
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "successful auth",
			fields: fields{
				issuer:       "https://jwt.com",
				jwtEndpoint:  "https://auth.com/jwt",
				keysEndpoint: "https://jwt.com/keys",
				headerName:   "jwt-header",
			},
			want: &Session{AuthURL: "https://auth.com/jwt?authRequestID=testState"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			a := assert.New(t)

			provider, err := New(tt.fields.issuer, tt.fields.jwtEndpoint, tt.fields.keysEndpoint, tt.fields.headerName)
			a.NoError(err)

			session, err := provider.BeginAuth("testState")
			a.NoError(err)

			a.Equal(tt.want.GetAuthURL(), session.GetAuthURL())
		})
	}
}

func TestProvider_FetchUser(t *testing.T) {
	type fields struct {
		issuer       string
		jwtEndpoint  string
		keysEndpoint string
		headerName   string
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
			name: "no tokens",
			fields: fields{
				issuer:       "https://jwt.com",
				jwtEndpoint:  "https://auth.com/jwt",
				keysEndpoint: "https://jwt.com/keys",
				headerName:   "jwt-header",
			},
			args: args{
				&Session{AuthURL: "https://auth.com/jwt?authRequestID=testState"},
			},
			want: want{
				err: func(err error) bool {
					return errors.Is(err, ErrNoTokens)
				},
			},
		},
		{
			name: "successful fetch",
			fields: fields{
				issuer:       "https://jwt.com",
				jwtEndpoint:  "https://auth.com/jwt",
				keysEndpoint: "https://jwt.com/keys",
				headerName:   "jwt-header",
			},
			args: args{
				&Session{
					AuthURL: "https://auth.com/jwt?authRequestID=testState",
					Tokens: &oidc.Tokens{
						Token: &oauth2.Token{},
						IDTokenClaims: func() oidc.IDTokenClaims {
							claims := oidc.EmptyIDTokenClaims()
							userinfo := oidc.NewUserInfo()
							userinfo.SetSubject("sub")
							userinfo.SetPicture("picture")
							userinfo.SetName("firstname lastname")
							userinfo.SetEmail("email", true)
							userinfo.SetGivenName("firstname")
							userinfo.SetFamilyName("lastname")
							userinfo.SetNickname("nickname")
							claims.SetUserinfo(userinfo)
							return claims
						}(),
					},
				},
			},
			want: want{
				user: idp.User{
					ID:          "sub",
					DisplayName: "firstname lastname",
					NickName:    "nickname",
					Email:       "email",
					AvatarURL:   "picture",
					FirstName:   "firstname",
					LastName:    "lastname",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			a := assert.New(t)

			provider, err := New(tt.fields.issuer, tt.fields.jwtEndpoint, tt.fields.keysEndpoint, tt.fields.headerName)
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
