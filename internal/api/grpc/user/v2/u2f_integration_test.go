//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/integration"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func TestServer_RegisterU2F(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()

	type args struct {
		ctx context.Context
		req *user.RegisterU2FRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RegisterU2FResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.RegisterU2FRequest{},
			},
			wantErr: true,
		},
		{
			name: "user mismatch",
			args: args{
				ctx: CTX,
				req: &user.RegisterU2FRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		/* TODO: after we are able to obtain a Bearer token for a human user
		https://github.com/zitadel/zitadel/issues/6022
		{
			name: "human user",
			args: args{
				ctx: CTX,
				req: &user.RegisterU2FRequest{
					UserId: userID,
				},
			},
			want: &user.RegisterU2FResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RegisterU2F(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			integration.AssertDetails(t, tt.want, got)
			if tt.want != nil {
				assert.NotEmpty(t, got.GetU2FId())
				assert.NotEmpty(t, got.GetPublicKeyCredentialCreationOptions())
				_, err = Tester.WebAuthN.CreateAttestationResponse(got.GetPublicKeyCredentialCreationOptions())
				require.NoError(t, err)
			}
		})
	}
}

func TestServer_VerifyU2FRegistration(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	/* TODO after we are able to obtain a Bearer token for a human user
	pkr, err := Client.RegisterU2F(CTX, &user.RegisterU2FRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, pkr.GetPublicKeyCredentialCreationOptions())

	attestationResponse, err := Tester.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	require.NoError(t, err)
	*/

	type args struct {
		ctx context.Context
		req *user.VerifyU2FRegistrationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.VerifyU2FRegistrationResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.VerifyU2FRegistrationRequest{
					U2FId:     "123",
					TokenName: "nice name",
				},
			},
			wantErr: true,
		},
		/* TODO after we are able to obtain a Bearer token for a human user
		{
			name: "success",
			args: args{
				ctx: CTX,
				req: &user.VerifyU2FRegistrationRequest{
					UserId:              userID,
					U2FId:               pkr.GetU2FId(),
					PublicKeyCredential: attestationResponse,
					TokenName:           "nice name",
				},
			},
			want: &user.VerifyU2FRegistrationResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		*/
		{
			name: "wrong credential",
			args: args{
				ctx: CTX,
				req: &user.VerifyU2FRegistrationRequest{
					UserId: userID,
					U2FId:  "123",
					PublicKeyCredential: &structpb.Struct{
						Fields: map[string]*structpb.Value{"foo": {Kind: &structpb.Value_StringValue{StringValue: "bar"}}},
					},
					TokenName: "nice name",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.VerifyU2FRegistration(tt.args.ctx, tt.args.req)
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
