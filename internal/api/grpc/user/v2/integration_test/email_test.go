//go:build integration

package user_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_ResendEmailCode(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()
	verifiedUserID := Instance.CreateHumanUserVerified(CTX, Instance.DefaultOrg.Id, gofakeit.Email(), gofakeit.Phone()).GetUserId()

	tests := []struct {
		name    string
		req     *user.ResendEmailCodeRequest
		want    *user.ResendEmailCodeResponse
		wantErr bool
	}{
		{
			name: "user not existing",
			req: &user.ResendEmailCodeRequest{
				UserId: "xxx",
			},
			wantErr: true,
		},
		{
			name: "user no code",
			req: &user.ResendEmailCodeRequest{
				UserId: verifiedUserID,
			},
			wantErr: true,
		},
		{
			name: "resend",
			req: &user.ResendEmailCodeRequest{
				UserId: userID,
			},
			want: &user.ResendEmailCodeResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "custom url template",
			req: &user.ResendEmailCodeRequest{
				UserId: userID,
				Verification: &user.ResendEmailCodeRequest_SendCode{
					SendCode: &user.SendEmailVerificationCode{
						UrlTemplate: gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
					},
				},
			},
			want: &user.ResendEmailCodeResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "template error",
			req: &user.ResendEmailCodeRequest{
				UserId: userID,
				Verification: &user.ResendEmailCodeRequest_SendCode{
					SendCode: &user.SendEmailVerificationCode{
						UrlTemplate: gu.Ptr("{{"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "return code",
			req: &user.ResendEmailCodeRequest{
				UserId: userID,
				Verification: &user.ResendEmailCodeRequest_ReturnCode{
					ReturnCode: &user.ReturnEmailVerificationCode{},
				},
			},
			want: &user.ResendEmailCodeResponse{
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
			got, err := Client.ResendEmailCode(CTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
			if tt.want.GetVerificationCode() != "" {
				assert.NotEmpty(t, got.GetVerificationCode())
			} else {
				assert.Empty(t, got.GetVerificationCode())
			}
		})
	}
}

func TestServer_SendEmailCode(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()
	verifiedUserID := Instance.CreateHumanUserVerified(CTX, Instance.DefaultOrg.Id, gofakeit.Email(), gofakeit.Phone()).GetUserId()

	tests := []struct {
		name    string
		req     *user.SendEmailCodeRequest
		want    *user.SendEmailCodeResponse
		wantErr bool
	}{
		{
			name: "user not existing",
			req: &user.SendEmailCodeRequest{
				UserId: "xxx",
			},
			wantErr: true,
		},
		{
			name: "user no code",
			req: &user.SendEmailCodeRequest{
				UserId: verifiedUserID,
			},
			want: &user.SendEmailCodeResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "resend",
			req: &user.SendEmailCodeRequest{
				UserId: userID,
			},
			want: &user.SendEmailCodeResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "custom url template",
			req: &user.SendEmailCodeRequest{
				UserId: userID,
				Verification: &user.SendEmailCodeRequest_SendCode{
					SendCode: &user.SendEmailVerificationCode{
						UrlTemplate: gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
					},
				},
			},
			want: &user.SendEmailCodeResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "template error",
			req: &user.SendEmailCodeRequest{
				UserId: userID,
				Verification: &user.SendEmailCodeRequest_SendCode{
					SendCode: &user.SendEmailVerificationCode{
						UrlTemplate: gu.Ptr("{{"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "return code",
			req: &user.SendEmailCodeRequest{
				UserId: userID,
				Verification: &user.SendEmailCodeRequest_ReturnCode{
					ReturnCode: &user.ReturnEmailVerificationCode{},
				},
			},
			want: &user.SendEmailCodeResponse{
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
			got, err := Client.SendEmailCode(CTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
			if tt.want.GetVerificationCode() != "" {
				assert.NotEmpty(t, got.GetVerificationCode())
			} else {
				assert.Empty(t, got.GetVerificationCode())
			}
		})
	}
}

func TestServer_VerifyEmail(t *testing.T) {
	userResp := Instance.CreateHumanUser(CTX)
	tests := []struct {
		name    string
		req     *user.VerifyEmailRequest
		want    *user.VerifyEmailResponse
		wantErr bool
	}{
		{
			name: "wrong code",
			req: &user.VerifyEmailRequest{
				UserId:           userResp.GetUserId(),
				VerificationCode: "xxx",
			},
			wantErr: true,
		},
		{
			name: "wrong user",
			req: &user.VerifyEmailRequest{
				UserId:           "xxx",
				VerificationCode: userResp.GetEmailCode(),
			},
			wantErr: true,
		},
		{
			name: "verify user",
			req: &user.VerifyEmailRequest{
				UserId:           userResp.GetUserId(),
				VerificationCode: userResp.GetEmailCode(),
			},
			want: &user.VerifyEmailResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.VerifyEmail(CTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
