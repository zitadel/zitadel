package webauthn

import (
	"context"
	"testing"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestWebAuthNsToCredentials(t *testing.T) {
	type args struct {
		ctx       context.Context
		webAuthNs []*domain.WebAuthNToken
		rpID      string
	}
	tests := []struct {
		name string
		args args
		want []webauthn.Credential
	}{
		{
			name: "unready credential",
			args: args{
				ctx: context.Background(),
				webAuthNs: []*domain.WebAuthNToken{
					{
						KeyID:           []byte("key1"),
						PublicKey:       []byte("publicKey1"),
						AttestationType: "attestation1",
						AAGUID:          []byte("aaguid1"),
						SignCount:       1,
						State:           domain.MFAStateNotReady,
					},
				},
				rpID: "example.com",
			},
			want: []webauthn.Credential{},
		},
		{
			name: "not matching rpID",
			args: args{
				ctx: context.Background(),
				webAuthNs: []*domain.WebAuthNToken{
					{
						KeyID:           []byte("key1"),
						PublicKey:       []byte("publicKey1"),
						AttestationType: "attestation1",
						AAGUID:          []byte("aaguid1"),
						SignCount:       1,
						State:           domain.MFAStateReady,
						RPID:            "other.com",
					},
				},
				rpID: "example.com",
			},
			want: []webauthn.Credential{},
		},
		{
			name: "matching rpID",
			args: args{
				ctx: context.Background(),
				webAuthNs: []*domain.WebAuthNToken{
					{
						KeyID:           []byte("key1"),
						PublicKey:       []byte("publicKey1"),
						AttestationType: "attestation1",
						AAGUID:          []byte("aaguid1"),
						SignCount:       1,
						State:           domain.MFAStateReady,
						RPID:            "example.com",
					},
				},
				rpID: "example.com",
			},
			want: []webauthn.Credential{
				{
					ID:              []byte("key1"),
					PublicKey:       []byte("publicKey1"),
					AttestationType: "attestation1",
					Authenticator: webauthn.Authenticator{
						AAGUID:    []byte("aaguid1"),
						SignCount: 1,
					},
				},
			},
		},
		{
			name: "no rpID, different host",
			args: args{
				ctx: http.WithDomainContext(context.Background(), &http.DomainCtx{
					InstanceHost: "other.com:443",
					PublicHost:   "other.com:443",
					Protocol:     "https",
				}),
				webAuthNs: []*domain.WebAuthNToken{
					{
						KeyID:           []byte("key1"),
						PublicKey:       []byte("publicKey1"),
						AttestationType: "attestation1",
						AAGUID:          []byte("aaguid1"),
						SignCount:       1,
						State:           domain.MFAStateReady,
						RPID:            "",
					},
				},
				rpID: "example.com",
			},
			want: []webauthn.Credential{},
		},
		{
			name: "no rpID, same host",
			args: args{
				ctx: http.WithDomainContext(context.Background(), &http.DomainCtx{
					InstanceHost: "example.com:443",
					PublicHost:   "example.com:443",
					Protocol:     "https",
				}),
				webAuthNs: []*domain.WebAuthNToken{
					{
						KeyID:           []byte("key1"),
						PublicKey:       []byte("publicKey1"),
						AttestationType: "attestation1",
						AAGUID:          []byte("aaguid1"),
						SignCount:       1,
						State:           domain.MFAStateReady,
						RPID:            "",
					},
				},
				rpID: "example.com",
			},
			want: []webauthn.Credential{
				{
					ID:              []byte("key1"),
					PublicKey:       []byte("publicKey1"),
					AttestationType: "attestation1",
					Authenticator: webauthn.Authenticator{
						AAGUID:    []byte("aaguid1"),
						SignCount: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, WebAuthNsToCredentials(tt.args.ctx, tt.args.webAuthNs, tt.args.rpID), "WebAuthNsToCredentials(%v, %v, %v)", tt.args.ctx, tt.args.webAuthNs, tt.args.rpID)
		})
	}
}
