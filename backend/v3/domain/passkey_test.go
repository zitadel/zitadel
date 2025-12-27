package domain_test

import (
	"context"
	"testing"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/api/http"
)

func TestPasskeysToCredentials(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx      context.Context
		passkeys []*domain.Passkey
		rpID     string
	}
	tt := []struct {
		name string
		args args
		want []webauthn.Credential
	}{
		{
			name: "matching rpID",
			args: args{
				ctx: context.Background(),
				passkeys: []*domain.Passkey{
					{
						KeyID:                        []byte("key1"),
						PublicKey:                    []byte("publicKey1"),
						AttestationType:              "attestation1",
						AuthenticatorAttestationGUID: []byte("aaguid1"),
						SignCount:                    1,
						RelyingPartyID:               "example.com",
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
			name: "not matching rpID",
			args: args{
				ctx: context.Background(),
				passkeys: []*domain.Passkey{
					{
						KeyID:                        []byte("key1"),
						PublicKey:                    []byte("publicKey1"),
						AttestationType:              "attestation1",
						AuthenticatorAttestationGUID: []byte("aaguid1"),
						SignCount:                    1,
						RelyingPartyID:               "other.com",
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
				passkeys: []*domain.Passkey{
					{
						KeyID:                        []byte("key1"),
						PublicKey:                    []byte("publicKey1"),
						AttestationType:              "attestation1",
						AuthenticatorAttestationGUID: []byte("aaguid1"),
						SignCount:                    1,
						RelyingPartyID:               "",
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
				passkeys: []*domain.Passkey{
					{
						KeyID:                        []byte("key1"),
						PublicKey:                    []byte("publicKey1"),
						AttestationType:              "attestation1",
						AuthenticatorAttestationGUID: []byte("aaguid1"),
						SignCount:                    1,
						RelyingPartyID:               "",
					},
				},
				rpID: "example.com",
			},
			want: []webauthn.Credential{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tc.want, domain.PasskeysToCredentials(tc.args.ctx, tc.args.passkeys, tc.args.rpID), "PasskeysToCredentials(%v, %v, %v)", tc.args.ctx, tc.args.passkeys, tc.args.rpID)
		})
	}
}
