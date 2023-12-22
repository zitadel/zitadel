package jwt

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/idp"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		name          string
		issuer        string
		jwtEndpoint   string
		keysEndpoint  string
		headerName    string
		encryptionAlg func(t *testing.T) crypto.EncryptionAlgorithm
	}
	type args struct {
		params []idp.Parameter
	}
	type want struct {
		session idp.Session
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "missing userAgentID error",
			fields: fields{
				issuer:       "https://jwt.com",
				jwtEndpoint:  "https://auth.com/jwt",
				keysEndpoint: "https://jwt.com/keys",
				headerName:   "jwt-header",
				encryptionAlg: func(t *testing.T) crypto.EncryptionAlgorithm {
					return crypto.CreateMockEncryptionAlg(gomock.NewController(t))
				},
			},
			args: args{
				params: nil,
			},
			want: want{
				err: func(err error) bool {
					return errors.Is(err, ErrMissingUserAgentID)
				},
			},
		},
		{
			name: "successful auth",
			fields: fields{
				issuer:       "https://jwt.com",
				jwtEndpoint:  "https://auth.com/jwt",
				keysEndpoint: "https://jwt.com/keys",
				headerName:   "jwt-header",
				encryptionAlg: func(t *testing.T) crypto.EncryptionAlgorithm {
					return crypto.CreateMockEncryptionAlg(gomock.NewController(t))
				},
			},
			args: args{
				params: []idp.Parameter{
					idp.UserAgentID("agent"),
				},
			},
			want: want{
				session: &Session{AuthURL: "https://auth.com/jwt?authRequestID=testState&userAgentID=YWdlbnQ"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			provider, err := New(
				tt.fields.name,
				tt.fields.issuer,
				tt.fields.jwtEndpoint,
				tt.fields.keysEndpoint,
				tt.fields.headerName,
				tt.fields.encryptionAlg(t),
			)
			require.NoError(t, err)

			ctx := context.Background()
			session, err := provider.BeginAuth(ctx, "testState", tt.args.params...)
			if tt.want.err != nil && !tt.want.err(err) {
				a.Fail("invalid error", err)
			}
			if tt.want.err == nil {
				a.NoError(err)
				wantHeaders, wantContent := tt.want.session.GetAuth(ctx)
				gotHeaders, gotContent := session.GetAuth(ctx)
				a.Equal(wantHeaders, gotHeaders)
				a.Equal(wantContent, gotContent)
			}
		})
	}
}

func TestProvider_Options(t *testing.T) {
	type fields struct {
		name          string
		issuer        string
		jwtEndpoint   string
		keysEndpoint  string
		headerName    string
		encryptionAlg func(t *testing.T) crypto.EncryptionAlgorithm
		opts          []ProviderOpts
	}
	type want struct {
		name            string
		linkingAllowed  bool
		creationAllowed bool
		autoCreation    bool
		autoUpdate      bool
		pkce            bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "default",
			fields: fields{
				name:         "jwt",
				issuer:       "https://jwt.com",
				jwtEndpoint:  "https://auth.com/jwt",
				keysEndpoint: "https://jwt.com/keys",
				headerName:   "jwt-header",
				encryptionAlg: func(t *testing.T) crypto.EncryptionAlgorithm {
					return crypto.CreateMockEncryptionAlg(gomock.NewController(t))
				},
				opts: nil,
			},
			want: want{
				name:            "jwt",
				linkingAllowed:  false,
				creationAllowed: false,
				autoCreation:    false,
				autoUpdate:      false,
				pkce:            false,
			},
		},
		{
			name: "all true",
			fields: fields{
				name:         "jwt",
				issuer:       "https://jwt.com",
				jwtEndpoint:  "https://auth.com/jwt",
				keysEndpoint: "https://jwt.com/keys",
				headerName:   "jwt-header",
				encryptionAlg: func(t *testing.T) crypto.EncryptionAlgorithm {
					return crypto.CreateMockEncryptionAlg(gomock.NewController(t))
				},
				opts: []ProviderOpts{
					WithLinkingAllowed(),
					WithCreationAllowed(),
					WithAutoCreation(),
					WithAutoUpdate(),
				},
			},
			want: want{
				name:            "jwt",
				linkingAllowed:  true,
				creationAllowed: true,
				autoCreation:    true,
				autoUpdate:      true,
				pkce:            true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			provider, err := New(
				tt.fields.name,
				tt.fields.issuer,
				tt.fields.jwtEndpoint,
				tt.fields.keysEndpoint,
				tt.fields.headerName,
				tt.fields.encryptionAlg(t),
				tt.fields.opts...,
			)
			require.NoError(t, err)

			a.Equal(tt.want.name, provider.Name())
			a.Equal(tt.want.linkingAllowed, provider.IsLinkingAllowed())
			a.Equal(tt.want.creationAllowed, provider.IsCreationAllowed())
			a.Equal(tt.want.autoCreation, provider.IsAutoCreation())
			a.Equal(tt.want.autoUpdate, provider.IsAutoUpdate())
		})
	}
}
