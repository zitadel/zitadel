package jwt

import (
	"context"
	"errors"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

func TestSession_FetchUser(t *testing.T) {
	type fields struct {
		issuer       string
		jwtEndpoint  string
		keysEndpoint string
		headerName   string
		authURL      string
		tokens       *oidc.Tokens
	}
	type want struct {
		user idp.User
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
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
				authURL:      "https://auth.com/jwt?authRequestID=testState",
				tokens: &oidc.Tokens{
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

			//provider, err := New(tt.fields.issuer, tt.fields.jwtEndpoint, tt.fields.keysEndpoint, tt.fields.headerName)
			//a.NoError(err)

			session := &Session{
				AuthURL: tt.fields.authURL,
				Tokens:  tt.fields.tokens,
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
