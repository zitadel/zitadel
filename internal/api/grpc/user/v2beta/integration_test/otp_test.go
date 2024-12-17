//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func TestServer_AddOTPSMS(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()
	Instance.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTX, userID)

	otherUser := Instance.CreateHumanUser(CTX).GetUserId()
	Instance.RegisterUserPasskey(CTX, otherUser)
	_, sessionTokenOtherUser, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTX, otherUser)

	userVerified := Instance.CreateHumanUser(CTX)
	_, err := Client.VerifyPhone(CTX, &user.VerifyPhoneRequest{
		UserId:           userVerified.GetUserId(),
		VerificationCode: userVerified.GetPhoneCode(),
	})
	require.NoError(t, err)
	Instance.RegisterUserPasskey(CTX, userVerified.GetUserId())
	_, sessionTokenVerified, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTX, userVerified.GetUserId())

	userVerified2 := Instance.CreateHumanUser(CTX)
	_, err = Client.VerifyPhone(CTX, &user.VerifyPhoneRequest{
		UserId:           userVerified2.GetUserId(),
		VerificationCode: userVerified2.GetPhoneCode(),
	})
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.AddOTPSMSRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.AddOTPSMSResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.AddOTPSMSRequest{},
			},
			wantErr: true,
		},
		{
			name: "no permission",
			args: args{
				ctx: integration.WithAuthorizationToken(context.Background(), sessionTokenOtherUser),
				req: &user.AddOTPSMSRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "phone not verified",
			args: args{
				ctx: integration.WithAuthorizationToken(context.Background(), sessionToken),
				req: &user.AddOTPSMSRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "add success",
			args: args{
				ctx: integration.WithAuthorizationToken(context.Background(), sessionTokenVerified),
				req: &user.AddOTPSMSRequest{
					UserId: userVerified.GetUserId(),
				},
			},
			want: &user.AddOTPSMSResponse{
				Details: &object.Details{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "add success, admin",
			args: args{
				ctx: CTX,
				req: &user.AddOTPSMSRequest{
					UserId: userVerified2.GetUserId(),
				},
			},
			want: &user.AddOTPSMSResponse{
				Details: &object.Details{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddOTPSMS(tt.args.ctx, tt.args.req)
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

func TestServer_RemoveOTPSMS(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()
	Instance.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTX, userID)

	userVerified := Instance.CreateHumanUser(CTX)
	Instance.RegisterUserPasskey(CTX, userVerified.GetUserId())
	_, err := Instance.Client.UserV2beta.VerifyPhone(CTX, &user.VerifyPhoneRequest{
		UserId:           userVerified.GetUserId(),
		VerificationCode: userVerified.GetPhoneCode(),
	})
	require.NoError(t, err)
	_, err = Instance.Client.UserV2beta.AddOTPSMS(CTX, &user.AddOTPSMSRequest{UserId: userVerified.GetUserId()})
	require.NoError(t, err)

	userSelf := Instance.CreateHumanUser(CTX)
	Instance.RegisterUserPasskey(CTX, userSelf.GetUserId())
	_, sessionTokenSelf, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTX, userSelf.GetUserId())
	userSelfCtx := integration.WithAuthorizationToken(context.Background(), sessionTokenSelf)
	_, err = Instance.Client.UserV2beta.VerifyPhone(CTX, &user.VerifyPhoneRequest{
		UserId:           userSelf.GetUserId(),
		VerificationCode: userSelf.GetPhoneCode(),
	})
	require.NoError(t, err)
	_, err = Instance.Client.UserV2beta.AddOTPSMS(CTX, &user.AddOTPSMSRequest{UserId: userSelf.GetUserId()})
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.RemoveOTPSMSRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RemoveOTPSMSResponse
		wantErr bool
	}{
		{
			name: "not added",
			args: args{
				ctx: integration.WithAuthorizationToken(context.Background(), sessionToken),
				req: &user.RemoveOTPSMSRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "success, self",
			args: args{
				ctx: userSelfCtx,
				req: &user.RemoveOTPSMSRequest{
					UserId: userSelf.GetUserId(),
				},
			},
			want: &user.RemoveOTPSMSResponse{
				Details: &object.Details{
					ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
				},
			},
		},
		{
			name: "success",
			args: args{
				ctx: CTX,
				req: &user.RemoveOTPSMSRequest{
					UserId: userVerified.GetUserId(),
				},
			},
			want: &user.RemoveOTPSMSResponse{
				Details: &object.Details{
					ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RemoveOTPSMS(tt.args.ctx, tt.args.req)
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

func TestServer_AddOTPEmail(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()
	Instance.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTX, userID)

	otherUser := Instance.CreateHumanUser(CTX).GetUserId()
	Instance.RegisterUserPasskey(CTX, otherUser)
	_, sessionTokenOtherUser, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTX, otherUser)

	userVerified := Instance.CreateHumanUser(CTX)
	_, err := Client.VerifyEmail(CTX, &user.VerifyEmailRequest{
		UserId:           userVerified.GetUserId(),
		VerificationCode: userVerified.GetEmailCode(),
	})
	require.NoError(t, err)
	Instance.RegisterUserPasskey(CTX, userVerified.GetUserId())
	_, sessionTokenVerified, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTX, userVerified.GetUserId())

	userVerified2 := Instance.CreateHumanUser(CTX)
	_, err = Client.VerifyEmail(CTX, &user.VerifyEmailRequest{
		UserId:           userVerified2.GetUserId(),
		VerificationCode: userVerified2.GetEmailCode(),
	})
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.AddOTPEmailRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.AddOTPEmailResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.AddOTPEmailRequest{},
			},
			wantErr: true,
		},
		{
			name: "user mismatch",
			args: args{
				ctx: integration.WithAuthorizationToken(context.Background(), sessionTokenOtherUser),
				req: &user.AddOTPEmailRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "email not verified",
			args: args{
				ctx: integration.WithAuthorizationToken(context.Background(), sessionToken),
				req: &user.AddOTPEmailRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "add success",
			args: args{
				ctx: integration.WithAuthorizationToken(context.Background(), sessionTokenVerified),
				req: &user.AddOTPEmailRequest{
					UserId: userVerified.GetUserId(),
				},
			},
			want: &user.AddOTPEmailResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "add success, admin",
			args: args{
				ctx: CTX,
				req: &user.AddOTPEmailRequest{
					UserId: userVerified2.GetUserId(),
				},
			},
			want: &user.AddOTPEmailResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddOTPEmail(tt.args.ctx, tt.args.req)
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

func TestServer_RemoveOTPEmail(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()
	Instance.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTX, userID)

	userVerified := Instance.CreateHumanUser(CTX)
	Instance.RegisterUserPasskey(CTX, userVerified.GetUserId())
	_, err := Client.VerifyEmail(CTX, &user.VerifyEmailRequest{
		UserId:           userVerified.GetUserId(),
		VerificationCode: userVerified.GetEmailCode(),
	})
	require.NoError(t, err)
	_, err = Client.AddOTPEmail(CTX, &user.AddOTPEmailRequest{UserId: userVerified.GetUserId()})
	require.NoError(t, err)

	userSelf := Instance.CreateHumanUser(CTX)
	Instance.RegisterUserPasskey(CTX, userSelf.GetUserId())
	_, sessionTokenSelf, _, _ := Instance.CreateVerifiedWebAuthNSession(t, IamCTX, userSelf.GetUserId())
	userSelfCtx := integration.WithAuthorizationToken(context.Background(), sessionTokenSelf)
	_, err = Client.VerifyEmail(CTX, &user.VerifyEmailRequest{
		UserId:           userSelf.GetUserId(),
		VerificationCode: userSelf.GetEmailCode(),
	})
	require.NoError(t, err)
	_, err = Client.AddOTPEmail(CTX, &user.AddOTPEmailRequest{UserId: userSelf.GetUserId()})
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.RemoveOTPEmailRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RemoveOTPEmailResponse
		wantErr bool
	}{
		{
			name: "not added",
			args: args{
				ctx: integration.WithAuthorizationToken(context.Background(), sessionToken),
				req: &user.RemoveOTPEmailRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "success, self",
			args: args{
				ctx: userSelfCtx,
				req: &user.RemoveOTPEmailRequest{
					UserId: userSelf.GetUserId(),
				},
			},
			want: &user.RemoveOTPEmailResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
				},
			},
		},
		{
			name: "success",
			args: args{
				ctx: CTX,
				req: &user.RemoveOTPEmailRequest{
					UserId: userVerified.GetUserId(),
				},
			},
			want: &user.RemoveOTPEmailResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RemoveOTPEmail(tt.args.ctx, tt.args.req)
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
