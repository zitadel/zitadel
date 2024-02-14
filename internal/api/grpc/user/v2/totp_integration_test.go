//go:build integration

package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func TestServer_RegisterTOTP(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	Tester.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)
	ctx := Tester.WithAuthorizationToken(CTX, sessionToken)

	type args struct {
		ctx context.Context
		req *user.RegisterTOTPRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RegisterTOTPResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: ctx,
				req: &user.RegisterTOTPRequest{},
			},
			wantErr: true,
		},
		{
			name: "user mismatch",
			args: args{
				ctx: ctx,
				req: &user.RegisterTOTPRequest{
					UserId: "wrong",
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: &user.RegisterTOTPRequest{
					UserId: userID,
				},
			},
			want: &user.RegisterTOTPResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RegisterTOTP(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			integration.AssertDetails(t, tt.want, got)
			assert.NotEmpty(t, got.GetUri())
			assert.NotEmpty(t, got.GetSecret())
		})
	}
}

func TestServer_VerifyTOTPRegistration(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	Tester.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)
	ctx := Tester.WithAuthorizationToken(CTX, sessionToken)

	reg, err := Client.RegisterTOTP(ctx, &user.RegisterTOTPRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	code, err := totp.GenerateCode(reg.Secret, time.Now())
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.VerifyTOTPRegistrationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.VerifyTOTPRegistrationResponse
		wantErr bool
	}{
		{
			name: "user mismatch",
			args: args{
				ctx: ctx,
				req: &user.VerifyTOTPRegistrationRequest{
					UserId: "wrong",
				},
			},
			wantErr: true,
		},
		{
			name: "wrong code",
			args: args{
				ctx: ctx,
				req: &user.VerifyTOTPRegistrationRequest{
					UserId: userID,
					Code:   "123",
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: &user.VerifyTOTPRegistrationRequest{
					UserId: userID,
					Code:   code,
				},
			},
			want: &user.VerifyTOTPRegistrationResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ResourceOwner,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.VerifyTOTPRegistration(tt.args.ctx, tt.args.req)
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
