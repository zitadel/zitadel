package apple

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	openid "github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

func TestSession_FetchUser(t *testing.T) {
	type fields struct {
		clientID      string
		teamID        string
		keyID         string
		privateKey    []byte
		redirectURI   string
		scopes        []string
		httpMock      func()
		authURL       string
		code          string
		tokens        *openid.Tokens[*openid.IDTokenClaims]
		userFormValue string
	}
	type want struct {
		err               error
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
		nonceSupported    bool
		isPrivateEmail    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "unauthenticated session, error",
			fields: fields{
				clientID:    "clientID",
				teamID:      "teamID",
				keyID:       "keyID",
				privateKey:  []byte(privateKey),
				redirectURI: "redirectURI",
				scopes:      []string{"openid"},
				httpMock:    func() {},
				authURL:     "https://appleid.apple.com/auth/authorize?client_id=clientID&redirect_uri=redirectURI&response_mode=form_post&response_type=code&scope=openid&state=testState",
				tokens:      nil,
			},
			want: want{
				err: oidc.ErrCodeMissing,
			},
		},
		{
			name: "no user param",
			fields: fields{
				clientID:    "clientID",
				teamID:      "teamID",
				keyID:       "keyID",
				privateKey:  []byte(privateKey),
				redirectURI: "redirectURI",
				scopes:      []string{"openid"},
				httpMock:    func() {},
				authURL:     "https://appleid.apple.com/auth/authorize?client_id=clientID&redirect_uri=redirectURI&response_mode=form_post&response_type=code&scope=openid&state=testState",
				tokens: &openid.Tokens[*openid.IDTokenClaims]{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   openid.BearerToken,
					},
					IDTokenClaims: id_token(),
				},
				userFormValue: "",
			},
			want: want{
				id:                "sub",
				firstName:         "",
				lastName:          "",
				displayName:       "",
				nickName:          "",
				preferredUsername: "email",
				email:             "email",
				isEmailVerified:   true,
				phone:             "",
				isPhoneVerified:   false,
				preferredLanguage: language.Und,
				avatarURL:         "",
				profile:           "",
				nonceSupported:    true,
				isPrivateEmail:    true,
			},
		},
		{
			name: "with user param",
			fields: fields{
				clientID:    "clientID",
				teamID:      "teamID",
				keyID:       "keyID",
				privateKey:  []byte(privateKey),
				redirectURI: "redirectURI",
				scopes:      []string{"openid"},
				httpMock:    func() {},
				authURL:     "https://appleid.apple.com/auth/authorize?client_id=clientID&redirect_uri=redirectURI&response_mode=form_post&response_type=code&scope=openid&state=testState",
				tokens: &openid.Tokens[*openid.IDTokenClaims]{
					Token: &oauth2.Token{
						AccessToken: "accessToken",
						TokenType:   openid.BearerToken,
					},
					IDTokenClaims: id_token(),
				},
				userFormValue: `{"name": {"firstName": "firstName", "lastName": "lastName"}}`,
			},
			want: want{
				id:                "sub",
				firstName:         "firstName",
				lastName:          "lastName",
				displayName:       "",
				nickName:          "",
				preferredUsername: "email",
				email:             "email",
				isEmailVerified:   true,
				phone:             "",
				isPhoneVerified:   false,
				preferredLanguage: language.Und,
				avatarURL:         "",
				profile:           "",
				nonceSupported:    true,
				isPrivateEmail:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.fields.httpMock()
			a := assert.New(t)

			// call the real discovery endpoint
			gock.New(issuer).Get(openid.DiscoveryEndpoint).EnableNetworking()
			provider, err := New(tt.fields.clientID, tt.fields.teamID, tt.fields.keyID, tt.fields.redirectURI, tt.fields.privateKey, tt.fields.scopes)
			require.NoError(t, err)

			session := &Session{
				Session: &oidc.Session{
					Provider: provider.Provider,
					AuthURL:  tt.fields.authURL,
					Code:     tt.fields.code,
					Tokens:   tt.fields.tokens,
				},
				UserFormValue: tt.fields.userFormValue,
			}

			user, err := session.FetchUser(context.Background())
			if tt.want.err != nil && !errors.Is(err, tt.want.err) {
				a.Fail("invalid error", "expected %v, got %v", tt.want.err, err)
			}
			if tt.want.err == nil {
				a.NoError(err)
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

func id_token() *openid.IDTokenClaims {
	return &openid.IDTokenClaims{
		TokenClaims: openid.TokenClaims{
			Issuer:     issuer,
			Subject:    "sub",
			Audience:   []string{"clientID"},
			Expiration: openid.FromTime(time.Now().Add(1 * time.Hour)),
			IssuedAt:   openid.FromTime(time.Now().Add(-1 * time.Second)),
			AuthTime:   openid.FromTime(time.Now().Add(-1 * time.Second)),
			Nonce:      "nonce",
			ClientID:   "clientID",
		},
		UserInfoEmail: openid.UserInfoEmail{
			Email:         "email",
			EmailVerified: true,
		},
		Claims: map[string]any{
			"nonce_supported":  true,
			"is_private_email": true,
		},
	}
}
