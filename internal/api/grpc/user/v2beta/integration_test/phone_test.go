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

func TestServer_SetPhone(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()

	tests := []struct {
		name    string
		req     *user.SetPhoneRequest
		want    *user.SetPhoneResponse
		wantErr bool
	}{
		{
			name: "default verification",
			req: &user.SetPhoneRequest{
				UserId: userID,
				Phone:  "+41791234568",
			},
			want: &user.SetPhoneResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "send verification",
			req: &user.SetPhoneRequest{
				UserId: userID,
				Phone:  "+41791234569",
				Verification: &user.SetPhoneRequest_SendCode{
					SendCode: &user.SendPhoneVerificationCode{},
				},
			},
			want: &user.SetPhoneResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "return code",
			req: &user.SetPhoneRequest{
				UserId: userID,
				Phone:  "+41791234566",
				Verification: &user.SetPhoneRequest_ReturnCode{
					ReturnCode: &user.ReturnPhoneVerificationCode{},
				},
			},
			want: &user.SetPhoneResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
				VerificationCode: gu.Ptr("xxx"),
			},
		},
		{
			name: "is verified true",
			req: &user.SetPhoneRequest{
				UserId: userID,
				Phone:  "+41791234565",
				Verification: &user.SetPhoneRequest_IsVerified{
					IsVerified: true,
				},
			},
			want: &user.SetPhoneResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "is verified false",
			req: &user.SetPhoneRequest{
				UserId: userID,
				Phone:  "+41791234564",
				Verification: &user.SetPhoneRequest_IsVerified{
					IsVerified: false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetPhone(CTX, tt.req)
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

func TestServer_ResendPhoneCode(t *testing.T) {
	userID := Instance.CreateHumanUser(CTX).GetUserId()
	verifiedUserID := Instance.CreateHumanUserVerified(CTX, Instance.DefaultOrg.Id, integration.Email(), integration.Phone()).GetUserId()

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
			} else {
				assert.Empty(t, got.GetVerificationCode())
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

func TestServer_RemovePhone(t *testing.T) {
	userResp := Instance.CreateHumanUser(CTX)
	failResp := Instance.CreateHumanUserNoPhone(CTX)
	otherUser := Instance.CreateHumanUser(CTX).GetUserId()
	doubleRemoveUser := Instance.CreateHumanUser(CTX)

	Instance.RegisterUserPasskey(CTX, otherUser)
	_, sessionTokenOtherUser, _, _ := Instance.CreateVerifiedWebAuthNSession(t, LoginCTX, otherUser)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *user.RemovePhoneRequest
		want    *user.RemovePhoneResponse
		wantErr bool
		dep     func(ctx context.Context, userID string) (*user.RemovePhoneResponse, error)
	}{
		{
			name: "remove phone",
			ctx:  CTX,
			req: &user.RemovePhoneRequest{
				UserId: userResp.GetUserId(),
			},
			want: &user.RemovePhoneResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
			dep: func(ctx context.Context, userID string) (*user.RemovePhoneResponse, error) {
				return nil, nil
			},
		},
		{
			name: "user without phone",
			ctx:  CTX,
			req: &user.RemovePhoneRequest{
				UserId: failResp.GetUserId(),
			},
			wantErr: true,
			dep: func(ctx context.Context, userID string) (*user.RemovePhoneResponse, error) {
				return nil, nil
			},
		},
		{
			name: "remove previously deleted phone",
			ctx:  CTX,
			req: &user.RemovePhoneRequest{
				UserId: doubleRemoveUser.GetUserId(),
			},
			wantErr: true,
			dep: func(ctx context.Context, userID string) (*user.RemovePhoneResponse, error) {
				return Client.RemovePhone(ctx, &user.RemovePhoneRequest{
					UserId: doubleRemoveUser.GetUserId(),
				})
			},
		},
		{
			name:    "no user id",
			ctx:     CTX,
			req:     &user.RemovePhoneRequest{},
			wantErr: true,
			dep: func(ctx context.Context, userID string) (*user.RemovePhoneResponse, error) {
				return nil, nil
			},
		},
		{
			name: "other user, no permission",
			ctx:  integration.WithAuthorizationToken(CTX, sessionTokenOtherUser),
			req: &user.RemovePhoneRequest{
				UserId: userResp.GetUserId(),
			},
			wantErr: true,
			dep: func(ctx context.Context, userID string) (*user.RemovePhoneResponse, error) {
				return nil, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, depErr := tt.dep(tt.ctx, tt.req.UserId)
			require.NoError(t, depErr)

			got, err := Client.RemovePhone(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}
