//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func TestServer_RegisterPasskey(t *testing.T) {
	userID := createHumanUser(t).GetUserId()
	reg, err := Client.CreatePasskeyRegistrationLink(CTX, &user.CreatePasskeyRegistrationLinkRequest{
		UserId: userID,
		Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
	})
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.RegisterPasskeyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RegisterPasskeyResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{},
			},
			wantErr: true,
		},
		{
			name: "register code",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{
					UserId:        userID,
					Code:          reg.GetCode(),
					Authenticator: user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_PLATFORM,
				},
			},
			want: &user.RegisterPasskeyResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "reuse code (not allowed)",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{
					UserId:        userID,
					Code:          reg.GetCode(),
					Authenticator: user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_PLATFORM,
				},
			},
			wantErr: true,
		},
		{
			name: "wrong code",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{
					UserId: userID,
					Code: &user.PasskeyRegistrationCode{
						Id:   reg.GetCode().GetId(),
						Code: "foobar",
					},
					Authenticator: user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_CROSS_PLATFORM,
				},
			},
			wantErr: true,
		},
		{
			name: "user mismatch",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		/* TODO after we are able to obtain a Bearer token for a human user
		{
			name: "human user",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{
					UserId: humanUserID,
				},
			},
			want: &user.RegisterPasskeyResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RegisterPasskey(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			integration.AssertDetails(t, tt.want, got)
			if tt.want != nil {
				assert.NotEmpty(t, got.GetPasskeyId())
				assert.NotEmpty(t, got.GetPublicKeyCredentialCreationOptions())
			}
		})
	}
}

func TestServer_VerifyPasskeyRegistration(t *testing.T) {
	userID := createHumanUser(t).GetUserId()
	reg, err := Client.CreatePasskeyRegistrationLink(CTX, &user.CreatePasskeyRegistrationLinkRequest{
		UserId: userID,
		Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
	})
	require.NoError(t, err)
	pkr, err := Client.RegisterPasskey(CTX, &user.RegisterPasskeyRequest{
		UserId: userID,
		Code:   reg.GetCode(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, pkr.GetPasskeyId())
	require.NotEmpty(t, pkr.GetPublicKeyCredentialCreationOptions())

	/*
		The following is copied from https://github.com/descope/virtualwebauthn,
		but I din't manage to make it work (yet). There are some decoding errors
		between the libraries.

			// The relying party settings should mirror those on the actual WebAuthn server
			rp := virtualwebauthn.RelyingParty{Name: "Example Corp", ID: "example.com", Origin: "https://example.com"}

			// A mock authenticator that represents a security key or biometrics module
			authenticator := virtualwebauthn.NewAuthenticator()

			// Create a new credential that we'll try to register with the relying party
			credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)

			// Parses the attestation options we got from the relying party to ensure they're valid
			parsedAttestationOptions, err := virtualwebauthn.ParseAttestationOptions(string(pkr.GetPublicKeyCredentialCreationOptions()))
			require.NoError(t, err)
			// Creates an attestation response that we can send to the relying party as if it came from
			// an actual browser and authenticator.
			attestationResponse := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsedAttestationOptions)
	*/

	type args struct {
		ctx context.Context
		req *user.VerifyPasskeyRegistrationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.VerifyPasskeyRegistrationResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.VerifyPasskeyRegistrationRequest{},
			},
			wantErr: true,
		},
		/*
			{
				name: "success",
				args: args{
					ctx: CTX,
					req: &user.VerifyPasskeyRegistrationRequest{
						UserId:              userID,
						PasskeyId:           pkr.GetPasskeyId(),
						PublicKeyCredential: []byte(attestationResponse),
						PasskeyName:         "nice name",
					},
				},
				//wantErr: true,
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.VerifyPasskeyRegistration(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_CreatePasskeyRegistrationLink(t *testing.T) {
	userID := createHumanUser(t).GetUserId()

	type args struct {
		ctx context.Context
		req *user.CreatePasskeyRegistrationLinkRequest
	}
	tests := []struct {
		name     string
		args     args
		want     *user.CreatePasskeyRegistrationLinkResponse
		wantCode bool
		wantErr  bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{},
			},
			wantErr: true,
		},
		{
			name: "send default mail",
			args: args{
				ctx: CTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{
					UserId: userID,
				},
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "send custom url",
			args: args{
				ctx: CTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{
					UserId: userID,
					Medium: &user.CreatePasskeyRegistrationLinkRequest_SendLink{
						SendLink: &user.SendPasskeyRegistrationLink{
							UrlTemplate: gu.Ptr("https://example.com/passkey/register?userID={{.UserID}}&orgID={{.ResourceOwner}}&codeID={{.CodeID}}&code={{.Code}}"),
						},
					},
				},
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "return code",
			args: args{
				ctx: CTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{
					UserId: userID,
					Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
				},
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
			wantCode: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.CreatePasskeyRegistrationLink(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			integration.AssertDetails(t, tt.want, got)
			if tt.wantCode {
				assert.NotEmpty(t, got.GetCode().GetId())
				assert.NotEmpty(t, got.GetCode().GetId())
			}
		})
	}
}
