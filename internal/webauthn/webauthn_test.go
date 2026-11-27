package webauthn

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestConfig_serverFromContext(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     string
		origin string
	}
	tests := []struct {
		name    string
		args    args
		want    *webauthn.WebAuthn
		wantErr error
	}{
		{
			name:    "webauthn error",
			args:    args{context.Background(), "", ""},
			wantErr: zerrors.ThrowInternal(nil, "WEBAU-UX9ta", "Errors.User.WebAuthN.ServerConfig"),
		},
		{
			name: "success from ctx",
			args: args{
				ctx:    http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "example.com", Protocol: "https"}),
				id:     "",
				origin: "",
			},
			want: &webauthn.WebAuthn{
				Config: &webauthn.Config{
					RPDisplayName: "DisplayName",
					RPID:          "example.com",
					RPOrigins:     []string{"https://example.com"},
				},
			},
		},
		{
			name: "success from id",
			args: args{
				ctx:    http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "example.com", Protocol: "https"}),
				id:     "external.com",
				origin: "https://external.com",
			},
			want: &webauthn.WebAuthn{
				Config: &webauthn.Config{
					RPDisplayName: "DisplayName",
					RPID:          "external.com",
					RPOrigins:     []string{"https://external.com"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Config{
				DisplayName:    "DisplayName",
				ExternalSecure: true,
			}
			got, err := w.serverFromContext(tt.args.ctx, tt.args.id, tt.args.origin)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.want != nil {
				require.NotNil(t, got)
				assert.Equal(t, tt.want.Config.RPDisplayName, got.Config.RPDisplayName)
				assert.Equal(t, tt.want.Config.RPID, got.Config.RPID)
				assert.Equal(t, tt.want.Config.RPOrigins, got.Config.RPOrigins)
			}
		})
	}
}

func TestConfig_BeginRegistration(t *testing.T) {
	type args struct {
		authType     domain.AuthenticatorAttachment
		verification domain.UserVerificationRequirement
		residentKey  domain.ResidentKeyRequirement
	}
	type selection struct {
		ResidentKey        string `json:"residentKey"`
		RequireResidentKey *bool  `json:"requireResidentKey"`
	}
	tests := []struct {
		name     string
		args     args
		wantSel  selection
		wantAbst bool
	}{
		{
			// Regression test for #12942: a passkey is a discoverable credential,
			// so its registration must request residentKey "required" instead of
			// leaving the requirement absent (WebAuthn default: "discouraged").
			name: "passkey requires a discoverable credential",
			args: args{
				authType:     domain.AuthenticatorAttachmentUnspecified,
				verification: domain.UserVerificationRequirementRequired,
				residentKey:  domain.ResidentKeyRequirementRequired,
			},
			wantSel: selection{ResidentKey: "required", RequireResidentKey: boolPtr(true)},
		},
		{
			// U2F (second factor) credentials are non-discoverable by design:
			// both members must stay absent so the options stay byte-identical
			// to previous releases.
			name: "u2f keeps the resident key members absent",
			args: args{
				authType:     domain.AuthenticatorAttachmentUnspecified,
				verification: domain.UserVerificationRequirementDiscouraged,
				residentKey:  domain.ResidentKeyRequirementUnspecified,
			},
			wantAbst: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "example.com", Protocol: "https"})
			w := &Config{DisplayName: "ZITADEL", ExternalSecure: true}
			user := &domain.Human{}
			user.AggregateID = "user1"
			user.Username = "test-user"
			user.Profile = &domain.Profile{}

			tok, err := w.BeginRegistration(ctx, user, "test-user", tt.args.authType, tt.args.verification, tt.args.residentKey, "")
			require.NoError(t, err)
			require.NotEmpty(t, tok.CredentialCreationData)

			var creation struct {
				PublicKey struct {
					AuthenticatorSelection selection `json:"authenticatorSelection"`
				} `json:"publicKey"`
			}
			require.NoError(t, json.Unmarshal(tok.CredentialCreationData, &creation))
			got := creation.PublicKey.AuthenticatorSelection

			if tt.wantAbst {
				assert.Empty(t, got.ResidentKey)
				assert.Nil(t, got.RequireResidentKey)
				return
			}
			assert.Equal(t, tt.wantSel, got)
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
