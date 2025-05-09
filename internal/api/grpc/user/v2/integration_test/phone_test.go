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

func TestServer_ResendPhoneCode(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()
	verifiedUserID := Instance.CreateHumanUserVerified(CTX, Instance.DefaultOrg.Id, gofakeit.Email(), gofakeit.Phone()).GetUserId()

	tests := []struct {
		name    string
		req     *user.ResendPhoneCodeRequest
		want    *user.ResendPhoneCodeResponse
		wantErr bool
	}{
		{
			name: "user not existing",
			req: &user.ResendPhoneCodeRequest{
				UserId: "xxx",
			},
			wantErr: true,
		},
		{
			name: "user not existing",
			req: &user.ResendPhoneCodeRequest{
				UserId: verifiedUserID,
			},
			wantErr: true,
		},
		{
			name: "resend code",
			req: &user.ResendPhoneCodeRequest{
				UserId: userID,
				Verification: &user.ResendPhoneCodeRequest_SendCode{
					SendCode: &user.SendPhoneVerificationCode{},
				},
			},
			want: &user.ResendPhoneCodeResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "return code",
			req: &user.ResendPhoneCodeRequest{
				UserId: userID,
				Verification: &user.ResendPhoneCodeRequest_ReturnCode{
					ReturnCode: &user.ReturnPhoneVerificationCode{},
				},
			},
			want: &user.ResendPhoneCodeResponse{
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
			got, err := Client.ResendPhoneCode(CTX, tt.req)
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

func TestServer_VerifyPhone(t *testing.T) {
	userResp := Instance.CreateHumanUser(CTX)
	tests := []struct {
		name    string
		req     *user.VerifyPhoneRequest
		want    *user.VerifyPhoneResponse
		wantErr bool
	}{
		{
			name: "wrong code",
			req: &user.VerifyPhoneRequest{
				UserId:           userResp.GetUserId(),
				VerificationCode: "xxx",
			},
			wantErr: true,
		},
		{
			name: "wrong user",
			req: &user.VerifyPhoneRequest{
				UserId:           "xxx",
				VerificationCode: userResp.GetPhoneCode(),
			},
			wantErr: true,
		},
		{
			name: "verify user",
			req: &user.VerifyPhoneRequest{
				UserId:           userResp.GetUserId(),
				VerificationCode: userResp.GetPhoneCode(),
			},
			want: &user.VerifyPhoneResponse{
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
			got, err := Client.VerifyPhone(CTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}
