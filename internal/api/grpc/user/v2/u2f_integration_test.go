//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"

	"github.com/zitadel/zitadel/internal/integration"
)

func TestServer_RegisterU2F(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	otherUser := Tester.CreateHumanUser(CTX).GetUserId()

	// We also need a user session
	Tester.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)
	Tester.RegisterUserPasskey(CTX, otherUser)
	_, sessionTokenOtherUser, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, otherUser)

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
			name: "admin user",
			args: args{
				ctx: CTX,
				req: &user.RegisterU2FRequest{
					UserId: userID,
				},
			},
			want: &user.RegisterU2FResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "other user, no permission",
			args: args{
				ctx: Tester.WithAuthorizationToken(CTX, sessionTokenOtherUser),
				req: &user.RegisterU2FRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "user setting its own passkey",
			args: args{
				ctx: Tester.WithAuthorizationToken(CTX, sessionToken),
				req: &user.RegisterU2FRequest{
					UserId: userID,
				},
			},
			want: &user.RegisterU2FResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
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
	ctx, userID, pkr := ctxFromNewUserWithRegisteredU2F(t)

	attestationResponse, err := Tester.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	require.NoError(t, err)

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
				ctx: ctx,
				req: &user.VerifyU2FRegistrationRequest{
					U2FId:     "123",
					TokenName: "nice name",
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: &user.VerifyU2FRegistrationRequest{
					UserId:              userID,
					U2FId:               pkr.GetU2FId(),
					PublicKeyCredential: attestationResponse,
					TokenName:           "nice name",
				},
			},
			want: &user.VerifyU2FRegistrationResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "wrong credential",
			args: args{
				ctx: ctx,
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

func ctxFromNewUserWithRegisteredU2F(t *testing.T) (context.Context, string, *user.RegisterU2FResponse) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	Tester.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)
	ctx := Tester.WithAuthorizationToken(CTX, sessionToken)

	pkr, err := Client.RegisterU2F(ctx, &user.RegisterU2FRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, pkr.GetPublicKeyCredentialCreationOptions())
	return ctx, userID, pkr
}

func ctxFromNewUserWithVerifiedU2F(t *testing.T) (context.Context, string, string) {
	ctx, userID, pkr := ctxFromNewUserWithRegisteredU2F(t)

	attestationResponse, err := Tester.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	require.NoError(t, err)

	_, err = Client.VerifyU2FRegistration(ctx, &user.VerifyU2FRegistrationRequest{
		UserId:              userID,
		U2FId:               pkr.GetU2FId(),
		PublicKeyCredential: attestationResponse,
		TokenName:           "nice name",
	})
	require.NoError(t, err)
	return ctx, userID, pkr.GetU2FId()
}

func TestServer_RemoveU2F(t *testing.T) {
	userIDWithout := Tester.CreateHumanUser(CTX).GetUserId()
	ctxRegistered, userIDRegistered, pkrRegistered := ctxFromNewUserWithRegisteredU2F(t)
	_, userIDVerified, u2fVerified := ctxFromNewUserWithVerifiedU2F(t)
	_, userIDVerifiedPermission, u2fVerifiedPermission := ctxFromNewUserWithVerifiedU2F(t)

	type args struct {
		ctx context.Context
		req *user.RemoveU2FRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RemoveU2FResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: ctxRegistered,
				req: &user.RemoveU2FRequest{
					U2FId: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "missing u2f id",
			args: args{
				ctx: ctxRegistered,
				req: &user.RemoveU2FRequest{
					UserId: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "success, registered",
			args: args{
				ctx: ctxRegistered,
				req: &user.RemoveU2FRequest{
					UserId: userIDRegistered,
					U2FId:  pkrRegistered.GetU2FId(),
				},
			},
			want: &user.RemoveU2FResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "no u2f, error",
			args: args{
				ctx: IamCTX,
				req: &user.RemoveU2FRequest{
					UserId: userIDWithout,
					U2FId:  pkrRegistered.GetU2FId(),
				},
			},
			wantErr: true,
		},
		{
			name: "success, IAMOwner permission, verified",
			args: args{
				ctx: IamCTX,
				req: &user.RemoveU2FRequest{
					UserId: userIDVerified,
					U2FId:  u2fVerified,
				},
			},
			want: &user.RemoveU2FResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "verified, permission error",
			args: args{
				ctx: UserCTX,
				req: &user.RemoveU2FRequest{
					UserId: userIDVerifiedPermission,
					U2FId:  u2fVerifiedPermission,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RemoveU2F(tt.args.ctx, tt.args.req)
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
