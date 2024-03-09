//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func TestServer_RequestPasswordReset(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()

	tests := []struct {
		name    string
		req     *user.PasswordResetRequest
		want    *user.PasswordResetResponse
		wantErr bool
	}{
		{
			name: "default medium",
			req: &user.PasswordResetRequest{
				UserId: userID,
			},
			want: &user.PasswordResetResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "custom url template",
			req: &user.PasswordResetRequest{
				UserId: userID,
				Medium: &user.PasswordResetRequest_SendLink{
					SendLink: &user.SendPasswordResetLink{
						NotificationType: user.NotificationType_NOTIFICATION_TYPE_Email,
						UrlTemplate:      gu.Ptr("https://example.com/password/change?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
					},
				},
			},
			want: &user.PasswordResetResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "template error",
			req: &user.PasswordResetRequest{
				UserId: userID,
				Medium: &user.PasswordResetRequest_SendLink{
					SendLink: &user.SendPasswordResetLink{
						UrlTemplate: gu.Ptr("{{"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "return code",
			req: &user.PasswordResetRequest{
				UserId: userID,
				Medium: &user.PasswordResetRequest_ReturnCode{
					ReturnCode: &user.ReturnPasswordResetCode{},
				},
			},
			want: &user.PasswordResetResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
				VerificationCode: gu.Ptr("xxx"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.PasswordReset(CTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			integration.AssertDetails(t, tt.want, got)
			if tt.want.GetVerificationCode() != "" {
				assert.NotEmpty(t, got.GetVerificationCode())
			}
		})
	}
}

func TestServer_SetPassword(t *testing.T) {
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
				userID := Tester.CreateHumanUser(CTX).GetUserId()
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
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "change successful",
			prepare: func(request *user.SetPasswordRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
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
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "set with code successful",
			prepare: func(request *user.SetPasswordRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
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
					ResourceOwner: Tester.Organisation.ID,
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
