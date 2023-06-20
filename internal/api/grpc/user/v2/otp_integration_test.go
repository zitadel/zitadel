//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func TestServer_RegisterOTP(t *testing.T) {
	// userID := Tester.CreateHumanUser(CTX).GetUserId()

	type args struct {
		ctx context.Context
		req *user.RegisterOTPRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RegisterOTPResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.RegisterOTPRequest{},
			},
			wantErr: true,
		},
		{
			name: "user mismatch",
			args: args{
				ctx: CTX,
				req: &user.RegisterOTPRequest{
					UserId: "wrong",
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
				req: &user.RegisterOTPRequest{
					UserId: userID,
				},
			},
			want: &user.RegisterOTPResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RegisterOTP(tt.args.ctx, tt.args.req)
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

func TestServer_VerifyOTPRegistration(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()

	/* TODO: after we are able to obtain a Bearer token for a human user
	reg, err := Client.RegisterOTP(CTX, &user.RegisterOTPRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	code, err := totp.GenerateCode(reg.Secret, time.Now())
	require.NoError(t, err)
	*/

	type args struct {
		ctx context.Context
		req *user.VerifyOTPRegistrationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.VerifyOTPRegistrationResponse
		wantErr bool
	}{
		{
			name: "user mismatch",
			args: args{
				ctx: CTX,
				req: &user.VerifyOTPRegistrationRequest{
					UserId: "wrong",
				},
			},
			wantErr: true,
		},
		{
			name: "wrong code",
			args: args{
				ctx: CTX,
				req: &user.VerifyOTPRegistrationRequest{
					UserId: userID,
					Code:   "123",
				},
			},
			wantErr: true,
		},
		/* TODO: after we are able to obtain a Bearer token for a human user
		https://github.com/zitadel/zitadel/issues/6022
		{
			name: "success",
			args: args{
				ctx: CTX,
				req: &user.VerifyOTPRegistrationRequest{
					UserId: userID,
					Code:   code,
				},
			},
			want: &user.VerifyOTPRegistrationResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ResourceOwner,
				},
			},
		},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.VerifyOTPRegistration(tt.args.ctx, tt.args.req)
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
