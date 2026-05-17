//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
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

// TestServer_PasswordHistoryReuse exercises the password-reuse prevention feature
// that was added in the password-age-new branch: PasswordComplexityPolicy.HistoryCount.
//
// Scenario A: reuse rejection + outside-window permitted.
//   - Set instance complexity policy to history_count=2.
//   - Change pw0→pw1→pw2→pw3 (four sequential self-service changes).
//   - Attempt pw3→pw0: pw0 is 3 generations back (outside the window of 2) → success.
//   - On a fresh user: pw0→pw1→pw2; attempt pw2→pw0 (pw0 is within window) → INVALID_ARGUMENT.
//
// Scenario B: current-hash inclusion.
//   - Set history_count=1.
//   - Request a password-reset code; attempt SetPassword with new password == current → INVALID_ARGUMENT.
func TestServer_PasswordHistoryReuse(t *testing.T) {
	const (
		pw0 = "H1st0ryP@ss0"
		pw1 = "H1st0ryP@ss1"
		pw2 = "H1st0ryP@ss2"
		pw3 = "H1st0ryP@ss3"
	)

	// Helper: change password using current-password verification (self-service ChangePassword path).
	changePassword := func(t *testing.T, userID, current, newPW string) error {
		t.Helper()
		_, err := Client.SetPassword(OrgCTX, &user.SetPasswordRequest{
			UserId: userID,
			NewPassword: &user.Password{
				Password: newPW,
			},
			Verification: &user.SetPasswordRequest_CurrentPassword{
				CurrentPassword: current,
			},
		})
		return err
	}

	// Helper: admin-set initial password (bypasses history check per spec).
	adminSetPassword := func(t *testing.T, userID, pw string) {
		t.Helper()
		_, err := Client.SetPassword(OrgCTX, &user.SetPasswordRequest{
			UserId: userID,
			NewPassword: &user.Password{
				Password: pw,
			},
		})
		require.NoError(t, err)
	}

	t.Run("scenario A: reuse rejection and outside-window permitted", func(t *testing.T) {
		// Set instance complexity policy: history_count=2.
		_, err := Instance.Client.Admin.UpdatePasswordComplexityPolicy(IamCTX, &admin_pb.UpdatePasswordComplexityPolicyRequest{
			MinLength:    8,
			HistoryCount: 2,
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			// Reset history_count to 0 so other tests are unaffected.
			_, _ = Instance.Client.Admin.UpdatePasswordComplexityPolicy(IamCTX, &admin_pb.UpdatePasswordComplexityPolicyRequest{
				MinLength:    8,
				HistoryCount: 0,
			})
		})

		t.Run("outside-window password is permitted", func(t *testing.T) {
			userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
			// Admin-set initial password (history-exempt).
			adminSetPassword(t, userID, pw0)
			// Chain: pw0 → pw1 → pw2 → pw3 (all self-service).
			require.NoError(t, changePassword(t, userID, pw0, pw1))
			require.NoError(t, changePassword(t, userID, pw1, pw2))
			require.NoError(t, changePassword(t, userID, pw2, pw3))
			// pw0 is now 3 generations back; history_count=2 only checks current+1 previous.
			// Attempt pw3 → pw0: should succeed.
			got, err := Client.SetPassword(OrgCTX, &user.SetPasswordRequest{
				UserId: userID,
				NewPassword: &user.Password{
					Password: pw0,
				},
				Verification: &user.SetPasswordRequest_CurrentPassword{
					CurrentPassword: pw3,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, got.GetDetails())
		})

		t.Run("in-window password is rejected", func(t *testing.T) {
			userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
			adminSetPassword(t, userID, pw0)
			// Chain: pw0 → pw1 → pw2 (two self-service changes).
			require.NoError(t, changePassword(t, userID, pw0, pw1))
			require.NoError(t, changePassword(t, userID, pw1, pw2))
			// pw0 is 2 generations back; history_count=2 window includes current+1 previous → pw0 is in window.
			err := changePassword(t, userID, pw2, pw0)
			require.Error(t, err)
			s, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, s.Code())
			assert.Contains(t, s.Message(), "Reused")
		})
	})

	t.Run("scenario B: current-hash inclusion via verify-code path", func(t *testing.T) {
		// Set instance complexity policy: history_count=1.
		_, err := Instance.Client.Admin.UpdatePasswordComplexityPolicy(IamCTX, &admin_pb.UpdatePasswordComplexityPolicyRequest{
			MinLength:    8,
			HistoryCount: 1,
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			_, _ = Instance.Client.Admin.UpdatePasswordComplexityPolicy(IamCTX, &admin_pb.UpdatePasswordComplexityPolicyRequest{
				MinLength:    8,
				HistoryCount: 0,
			})
		})

		userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
		// Admin-set initial password pw0. This is the "current" hash.
		adminSetPassword(t, userID, pw0)

		// Request a reset code.
		resetResp, err := Client.PasswordReset(OrgCTX, &user.PasswordResetRequest{
			UserId: userID,
			Medium: &user.PasswordResetRequest_ReturnCode{
				ReturnCode: &user.ReturnPasswordResetCode{},
			},
		})
		require.NoError(t, err)
		code := resetResp.GetVerificationCode()
		require.NotEmpty(t, code)

		// Attempt SetPassword via verify-code with new password == current (pw0).
		// Spec: current hash IS in the check list → must be rejected.
		_, err = Client.SetPassword(OrgCTX, &user.SetPasswordRequest{
			UserId: userID,
			NewPassword: &user.Password{
				Password: pw0,
			},
			Verification: &user.SetPasswordRequest_VerificationCode{
				VerificationCode: code,
			},
		})
		require.Error(t, err)
		s, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, s.Code())
		assert.Contains(t, s.Message(), "Reused")
	})
}
