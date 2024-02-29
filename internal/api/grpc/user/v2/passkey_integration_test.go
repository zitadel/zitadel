//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func TestServer_RegisterPasskey(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	reg, err := Client.CreatePasskeyRegistrationLink(CTX, &user.CreatePasskeyRegistrationLinkRequest{
		UserId: userID,
		Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
	})
	require.NoError(t, err)

	// We also need a user session
	Tester.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)

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
					ChangeDate:    timestamppb.Now(),
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
		{
			name: "user setting its own passkey",
			args: args{
				ctx: Tester.WithAuthorizationToken(CTX, sessionToken),
				req: &user.RegisterPasskeyRequest{
					UserId: userID,
				},
			},
			want: &user.RegisterPasskeyResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
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
				_, err = Tester.WebAuthN.CreateAttestationResponse(got.GetPublicKeyCredentialCreationOptions())
				require.NoError(t, err)
			}
		})
	}
}

func TestServer_VerifyPasskeyRegistration(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
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

	attestationResponse, err := Tester.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	require.NoError(t, err)

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
				req: &user.VerifyPasskeyRegistrationRequest{
					PasskeyId:           pkr.GetPasskeyId(),
					PublicKeyCredential: attestationResponse,
					PasskeyName:         "nice name",
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: CTX,
				req: &user.VerifyPasskeyRegistrationRequest{
					UserId:              userID,
					PasskeyId:           pkr.GetPasskeyId(),
					PublicKeyCredential: attestationResponse,
					PasskeyName:         "nice name",
				},
			},
			want: &user.VerifyPasskeyRegistrationResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "wrong credential",
			args: args{
				ctx: CTX,
				req: &user.VerifyPasskeyRegistrationRequest{
					UserId:    userID,
					PasskeyId: pkr.GetPasskeyId(),
					PublicKeyCredential: &structpb.Struct{
						Fields: map[string]*structpb.Value{"foo": {Kind: &structpb.Value_StringValue{StringValue: "bar"}}},
					},
					PasskeyName: "nice name",
				},
			},
			wantErr: true,
		},
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
	userID := Tester.CreateHumanUser(CTX).GetUserId()

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
					ChangeDate:    timestamppb.Now(),
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
							UrlTemplate: gu.Ptr("https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}"),
						},
					},
				},
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
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
					ChangeDate:    timestamppb.Now(),
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
