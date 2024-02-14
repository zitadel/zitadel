//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func TestServer_AddOTPSMS(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	Tester.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)

	// TODO: add when phone can be added to user
	/*
		userIDPhone := Tester.CreateHumanUser(CTX).GetUserId()
		Tester.RegisterUserPasskey(CTX, userIDPhone)
		_, sessionTokenPhone, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userIDPhone)
	*/
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
			name: "user mismatch",
			args: args{
				ctx: CTX,
				req: &user.AddOTPSMSRequest{
					UserId: "wrong",
				},
			},
			wantErr: true,
		},
		{
			name: "phone not verified",
			args: args{
				ctx: Tester.WithAuthorizationToken(context.Background(), sessionToken),
				req: &user.AddOTPSMSRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		// TODO: add when phone can be added to user
		/*
			{
				name: "add success",
				args: args{
					ctx: Tester.WithAuthorizationToken(context.Background(), sessionTokenPhone),
					req: &user.AddOTPSMSRequest{
						UserId: userID,
					},
				},
				want: &user.AddOTPSMSResponse{
					Details: &object.Details{
						ResourceOwner: Tester.Organisation.ID,
					},
				},
			},
		*/
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
	// TODO: add when phone can be added to user
	/*
		userID := Tester.CreateHumanUser(CTX).GetUserId()
		Tester.RegisterUserPasskey(CTX, userID)
		_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)
	*/

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
				ctx: CTX,
				req: &user.RemoveOTPSMSRequest{
					UserId: "wrong",
				},
			},
			wantErr: true,
		},
		// TODO: add when phone can be added to user
		/*
			{
				name: "success",
				args: args{
					ctx: Tester.WithAuthorizationToken(context.Background(), sessionToken),
					req: &user.RemoveOTPSMSRequest{
						UserId: userID,
					},
				},
				want: &user.RemoveOTPSMSResponse{
					Details: &object.Details{
						ResourceOwner: Tester.Organisation.ResourceOwner,
					},
				},
			},
		*/
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
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	Tester.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)

	userVerified := Tester.CreateHumanUser(CTX)
	_, err := Tester.Client.UserV2.VerifyEmail(CTX, &user.VerifyEmailRequest{
		UserId:           userVerified.GetUserId(),
		VerificationCode: userVerified.GetEmailCode(),
	})
	require.NoError(t, err)
	Tester.RegisterUserPasskey(CTX, userVerified.GetUserId())
	_, sessionTokenVerified, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userVerified.GetUserId())

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
				ctx: CTX,
				req: &user.AddOTPEmailRequest{
					UserId: "wrong",
				},
			},
			wantErr: true,
		},
		{
			name: "email not verified",
			args: args{
				ctx: Tester.WithAuthorizationToken(context.Background(), sessionToken),
				req: &user.AddOTPEmailRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "add success",
			args: args{
				ctx: Tester.WithAuthorizationToken(context.Background(), sessionTokenVerified),
				req: &user.AddOTPEmailRequest{
					UserId: userVerified.GetUserId(),
				},
			},
			want: &user.AddOTPEmailResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
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
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	Tester.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)

	userVerified := Tester.CreateHumanUser(CTX)
	Tester.RegisterUserPasskey(CTX, userVerified.GetUserId())
	_, sessionTokenVerified, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userVerified.GetUserId())
	userVerifiedCtx := Tester.WithAuthorizationToken(context.Background(), sessionTokenVerified)
	_, err := Tester.Client.UserV2.VerifyEmail(userVerifiedCtx, &user.VerifyEmailRequest{
		UserId:           userVerified.GetUserId(),
		VerificationCode: userVerified.GetEmailCode(),
	})
	require.NoError(t, err)
	_, err = Tester.Client.UserV2.AddOTPEmail(userVerifiedCtx, &user.AddOTPEmailRequest{UserId: userVerified.GetUserId()})
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
				ctx: Tester.WithAuthorizationToken(context.Background(), sessionToken),
				req: &user.RemoveOTPEmailRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: userVerifiedCtx,
				req: &user.RemoveOTPEmailRequest{
					UserId: userVerified.GetUserId(),
				},
			},
			want: &user.RemoveOTPEmailResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ResourceOwner,
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
