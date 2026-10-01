//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_RequestPasswordReset(t *testing.T) {
	userID := Instance.CreateHumanUser(OrgCTX).GetUserId()

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
					ResourceOwner: Instance.DefaultOrg.Id,
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
					ResourceOwner: Instance.DefaultOrg.Id,
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
					ResourceOwner: Instance.DefaultOrg.Id,
				},
				VerificationCode: gu.Ptr("xxx"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.PasswordReset(OrgCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
			if tt.want.GetVerificationCode() != "" {
				assert.NotEmpty(t, got.GetVerificationCode())
			}
		})
	}
}

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
				ctx: OrgCTX,
				req: &user.SetPasswordRequest{},
			},
			wantErr: true,
		},
		{
			name: "set successful",
			prepare: func(request *user.SetPasswordRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				ctx: OrgCTX,
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
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				_, err := Client.SetPassword(OrgCTX, &user.SetPasswordRequest{
					UserId: userID,
					NewPassword: &user.Password{
						Password: "InitialPassw0rd!",
					},
				})
				return err
			},
			args: args{
				ctx: OrgCTX,
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
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				resp, err := Client.PasswordReset(OrgCTX, &user.PasswordResetRequest{
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
				ctx: OrgCTX,
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

func TestServer_CreateUser_FirebaseScryptPassword(t *testing.T) {
	// Public sample from the firebase/scrypt README.
	const (
		hash     = "$firebasescrypt$ln=14,r=8$42xEC+ixf3L2lw==$lSrfV15cpx95/sZS2W9c9Kp6i/LVgQNDNC/qzrCnh1SAyZvqmZqAjTdn3aoItz+VHjoZilo78198JAdRuid5lQ==$Bw==$jxspr8Ki0RYycVU8zykbdLGjFQ3McFUH0uiiTvC8pVMXAn210wjLNmdZJzxUECKbm0QsEmYUSDzZvpjeJ9WmXA=="
		password = "user1password"
	)
	email := integration.Email()
	createResp, err := Client.CreateUser(OrgCTX, &user.CreateUserRequest{
		OrganizationId: Instance.DefaultOrg.Id,
		Username:       &email,
		UserType: &user.CreateUserRequest_Human_{
			Human: &user.CreateUserRequest_Human{
				Profile: &user.SetHumanProfile{
					GivenName:  "Donald",
					FamilyName: "Duck",
				},
				Email: &user.SetHumanEmail{
					Email: email,
					Verification: &user.SetHumanEmail_IsVerified{
						IsVerified: true,
					},
				},
				PasswordType: &user.CreateUserRequest_Human_HashedPassword{
					HashedPassword: &user.HashedPassword{
						Hash: hash,
					},
				},
			},
		},
	})
	require.NoError(t, err)
	userID := createResp.GetId()

	Instance.TriggerUserByID(OrgCTX, userID)

	_, err = Instance.Client.SessionV2.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: userID},
			},
			Password: &session.CheckPassword{
				Password: "wrong",
			},
		},
	})
	integration.AssertGrpcStatus(t, codes.InvalidArgument, err)

	Instance.CreatePasswordSession(t, LoginCTX, userID, password)
}
