//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_Deprecated_SetPassword(t *testing.T) {
	type args struct {
		ctx context.Context
		req *user.SetPasswordRequest
	}
	tests := []struct {
		name    string
		prepare func(request *user.SetPasswordRequest) error
		args    args
		want    *user.SetPasswordResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			prepare: func(request *user.SetPasswordRequest) error {
				return nil
			},
			args: args{
				ctx: CTX,
				req: &user.SetPasswordRequest{},
			},
			wantErr: true,
		},
		{
			name: "set successful",
			prepare: func(request *user.SetPasswordRequest) error {
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &user.SetPasswordRequest{
					NewPassword: &user.Password{
						Password: "Secr3tP4ssw0rd!",
					},
				},
			},
			want: &user.SetPasswordResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "change successful",
			prepare: func(request *user.SetPasswordRequest) error {
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				_, err := Client.SetPassword(CTX, &user.SetPasswordRequest{
					UserId: userID,
					NewPassword: &user.Password{
						Password: "InitialPassw0rd!",
					},
				})
				return err
			},
			args: args{
				ctx: CTX,
				req: &user.SetPasswordRequest{
					NewPassword: &user.Password{
						Password: "Secr3tP4ssw0rd!",
					},
					Verification: &user.SetPasswordRequest_CurrentPassword{
						CurrentPassword: "InitialPassw0rd!",
					},
				},
			},
			want: &user.SetPasswordResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "set with code successful",
			prepare: func(request *user.SetPasswordRequest) error {
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				resp, err := Client.PasswordReset(CTX, &user.PasswordResetRequest{
					UserId: userID,
					Medium: &user.PasswordResetRequest_ReturnCode{
						ReturnCode: &user.ReturnPasswordResetCode{},
					},
				})
				if err != nil {
					return err
				}
				request.Verification = &user.SetPasswordRequest_VerificationCode{
					VerificationCode: resp.GetVerificationCode(),
				}
				return nil
			},
			args: args{
				ctx: CTX,
				req: &user.SetPasswordRequest{
					NewPassword: &user.Password{
						Password: "Secr3tP4ssw0rd!",
					},
				},
			},
			want: &user.SetPasswordResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.SetPassword(tt.args.ctx, tt.args.req)
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
