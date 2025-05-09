//go:build integration

package user_test

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_RequestPasswordReset(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()

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
			got, err := Client.PasswordReset(CTX, tt.req)
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
