//go:build integration

package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"

	"github.com/zitadel/zitadel/internal/integration"
)

func TestServer_RegisterPasskey(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	reg, err := Client.CreatePasskeyRegistrationLink(CTX, &user.CreatePasskeyRegistrationLinkRequest{
		UserId: userID,
		Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
	})
	require.NoError(t, err)

	// We also need a user session
	Tester.RegisterUserPasskey(CTX, userID)
	_, sessionToken, _, _ := Tester.CreateVerifiedWebAuthNSession(t, CTX, userID)

	type args struct {
		ctx context.Context
		req *user.RegisterPasskeyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RegisterPasskeyResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{},
			},
			wantErr: true,
		},
		{
			name: "register code",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{
					UserId:        userID,
					Code:          reg.GetCode(),
					Authenticator: user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_PLATFORM,
				},
			},
			want: &user.RegisterPasskeyResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "reuse code (not allowed)",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{
					UserId:        userID,
					Code:          reg.GetCode(),
					Authenticator: user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_PLATFORM,
				},
			},
			wantErr: true,
		},
		{
			name: "wrong code",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{
					UserId: userID,
					Code: &user.PasskeyRegistrationCode{
						Id:   reg.GetCode().GetId(),
						Code: "foobar",
					},
					Authenticator: user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_CROSS_PLATFORM,
				},
			},
			wantErr: true,
		},
		{
			name: "user mismatch",
			args: args{
				ctx: CTX,
				req: &user.RegisterPasskeyRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "user setting its own passkey",
			args: args{
				ctx: Tester.WithAuthorizationToken(CTX, sessionToken),
				req: &user.RegisterPasskeyRequest{
					UserId: userID,
				},
			},
			want: &user.RegisterPasskeyResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RegisterPasskey(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			integration.AssertDetails(t, tt.want, got)
			if tt.want != nil {
				assert.NotEmpty(t, got.GetPasskeyId())
				assert.NotEmpty(t, got.GetPublicKeyCredentialCreationOptions())
				_, err = Tester.WebAuthN.CreateAttestationResponse(got.GetPublicKeyCredentialCreationOptions())
				require.NoError(t, err)
			}
		})
	}
}

func TestServer_VerifyPasskeyRegistration(t *testing.T) {
	userID, pkr := userWithPasskeyRegistered(t)

	attestationResponse, err := Tester.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.VerifyPasskeyRegistrationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.VerifyPasskeyRegistrationResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.VerifyPasskeyRegistrationRequest{
					PasskeyId:           pkr.GetPasskeyId(),
					PublicKeyCredential: attestationResponse,
					PasskeyName:         "nice name",
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: CTX,
				req: &user.VerifyPasskeyRegistrationRequest{
					UserId:              userID,
					PasskeyId:           pkr.GetPasskeyId(),
					PublicKeyCredential: attestationResponse,
					PasskeyName:         "nice name",
				},
			},
			want: &user.VerifyPasskeyRegistrationResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "wrong credential",
			args: args{
				ctx: CTX,
				req: &user.VerifyPasskeyRegistrationRequest{
					UserId:    userID,
					PasskeyId: pkr.GetPasskeyId(),
					PublicKeyCredential: &structpb.Struct{
						Fields: map[string]*structpb.Value{"foo": {Kind: &structpb.Value_StringValue{StringValue: "bar"}}},
					},
					PasskeyName: "nice name",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.VerifyPasskeyRegistration(tt.args.ctx, tt.args.req)
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

func TestServer_CreatePasskeyRegistrationLink(t *testing.T) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()

	type args struct {
		ctx context.Context
		req *user.CreatePasskeyRegistrationLinkRequest
	}
	tests := []struct {
		name     string
		args     args
		want     *user.CreatePasskeyRegistrationLinkResponse
		wantCode bool
		wantErr  bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: CTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{},
			},
			wantErr: true,
		},
		{
			name: "send default mail",
			args: args{
				ctx: CTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{
					UserId: userID,
				},
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "send custom url",
			args: args{
				ctx: CTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{
					UserId: userID,
					Medium: &user.CreatePasskeyRegistrationLinkRequest_SendLink{
						SendLink: &user.SendPasskeyRegistrationLink{
							UrlTemplate: gu.Ptr("https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}"),
						},
					},
				},
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "return code",
			args: args{
				ctx: CTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{
					UserId: userID,
					Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
				},
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
			wantCode: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.CreatePasskeyRegistrationLink(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			integration.AssertDetails(t, tt.want, got)
			if tt.wantCode {
				assert.NotEmpty(t, got.GetCode().GetId())
				assert.NotEmpty(t, got.GetCode().GetId())
			}
		})
	}
}

func userWithPasskeyRegistered(t *testing.T) (string, *user.RegisterPasskeyResponse) {
	userID := Tester.CreateHumanUser(CTX).GetUserId()
	return userID, passkeyRegister(t, userID)
}

func userWithPasskeyVerified(t *testing.T) (string, string) {
	userID, pkr := userWithPasskeyRegistered(t)
	return userID, passkeyVerify(t, userID, pkr)
}

func passkeyRegister(t *testing.T, userID string) *user.RegisterPasskeyResponse {
	reg, err := Client.CreatePasskeyRegistrationLink(CTX, &user.CreatePasskeyRegistrationLinkRequest{
		UserId: userID,
		Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
	})
	require.NoError(t, err)
	pkr, err := Client.RegisterPasskey(CTX, &user.RegisterPasskeyRequest{
		UserId: userID,
		Code:   reg.GetCode(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, pkr.GetPasskeyId())
	require.NotEmpty(t, pkr.GetPublicKeyCredentialCreationOptions())
	return pkr
}

func passkeyVerify(t *testing.T, userID string, pkr *user.RegisterPasskeyResponse) string {
	attestationResponse, err := Tester.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	require.NoError(t, err)

	_, err = Client.VerifyPasskeyRegistration(CTX, &user.VerifyPasskeyRegistrationRequest{
		UserId:              userID,
		PasskeyId:           pkr.GetPasskeyId(),
		PublicKeyCredential: attestationResponse,
		PasskeyName:         "nice name",
	})
	require.NoError(t, err)
	return pkr.GetPasskeyId()
}

func TestServer_RemovePasskey(t *testing.T) {
	userIDWithout := Tester.CreateHumanUser(CTX).GetUserId()
	userIDRegistered, pkrRegistered := userWithPasskeyRegistered(t)
	userIDVerified, passkeyIDVerified := userWithPasskeyVerified(t)
	userIDVerifiedPermission, passkeyIDVerifiedPermission := userWithPasskeyVerified(t)

	type args struct {
		ctx context.Context
		req *user.RemovePasskeyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RemovePasskeyResponse
		wantErr bool
	}{
		{
			name: "missing user id",
			args: args{
				ctx: IamCTX,
				req: &user.RemovePasskeyRequest{
					PasskeyId: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "missing passkey id",
			args: args{
				ctx: IamCTX,
				req: &user.RemovePasskeyRequest{
					UserId: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "success, registered",
			args: args{
				ctx: IamCTX,
				req: &user.RemovePasskeyRequest{
					UserId:    userIDRegistered,
					PasskeyId: pkrRegistered.GetPasskeyId(),
				},
			},
			want: &user.RemovePasskeyResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "no passkey, error",
			args: args{
				ctx: IamCTX,
				req: &user.RemovePasskeyRequest{
					UserId:    userIDWithout,
					PasskeyId: pkrRegistered.GetPasskeyId(),
				},
			},
			wantErr: true,
		},
		{
			name: "success, verified",
			args: args{
				ctx: IamCTX,
				req: &user.RemovePasskeyRequest{
					UserId:    userIDVerified,
					PasskeyId: passkeyIDVerified,
				},
			},
			want: &user.RemovePasskeyResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "verified, permission error",
			args: args{
				ctx: UserCTX,
				req: &user.RemovePasskeyRequest{
					UserId:    userIDVerifiedPermission,
					PasskeyId: passkeyIDVerifiedPermission,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RemovePasskey(tt.args.ctx, tt.args.req)
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

func TestServer_ListPasskeys(t *testing.T) {
	userIDWithout := Tester.CreateHumanUser(CTX).GetUserId()
	userIDRegistered, _ := userWithPasskeyRegistered(t)
	userIDVerified, passkeyIDVerified := userWithPasskeyVerified(t)

	userIDMulti, passkeyIDMulti1 := userWithPasskeyVerified(t)
	passkeyIDMulti2 := passkeyVerify(t, userIDMulti, passkeyRegister(t, userIDMulti))

	type args struct {
		ctx context.Context
		req *user.ListPasskeysRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.ListPasskeysResponse
		wantErr bool
	}{
		{
			name: "list passkeys, no permission",
			args: args{
				UserCTX,
				&user.ListPasskeysRequest{
					UserId: userIDVerified,
				},
			},
			want: &user.ListPasskeysResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				Result: []*user.Passkey{},
			},
		},
		{
			name: "list passkeys, none",
			args: args{
				UserCTX,
				&user.ListPasskeysRequest{
					UserId: userIDWithout,
				},
			},
			want: &user.ListPasskeysResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				Result: []*user.Passkey{},
			},
		},
		{
			name: "list passkeys, registered",
			args: args{
				UserCTX,
				&user.ListPasskeysRequest{
					UserId: userIDRegistered,
				},
			},
			want: &user.ListPasskeysResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				Result: []*user.Passkey{},
			},
		},
		{
			name: "list passkeys, ok",
			args: args{
				IamCTX,
				&user.ListPasskeysRequest{
					UserId: userIDVerified,
				},
			},
			want: &user.ListPasskeysResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				Result: []*user.Passkey{
					{
						Id:    passkeyIDVerified,
						State: user.AuthFactorState_AUTH_FACTOR_STATE_READY,
						Name:  "nice name",
					},
				},
			},
		},
		{
			name: "list idp links, multi, ok",
			args: args{
				IamCTX,
				&user.ListPasskeysRequest{
					UserId: userIDMulti,
				},
			},
			want: &user.ListPasskeysResponse{
				Details: &object.ListDetails{
					TotalResult: 2,
					Timestamp:   timestamppb.Now(),
				},
				Result: []*user.Passkey{
					{
						Id:    passkeyIDMulti1,
						State: user.AuthFactorState_AUTH_FACTOR_STATE_READY,
						Name:  "nice name",
					},
					{
						Id:    passkeyIDMulti2,
						State: user.AuthFactorState_AUTH_FACTOR_STATE_READY,
						Name:  "nice name",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration := time.Minute
			if ctxDeadline, ok := CTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := Client.ListPasskeys(tt.args.ctx, tt.args.req)
				assertErr := assert.NoError
				if tt.wantErr {
					assertErr = assert.Error
				}
				assertErr(ttt, listErr)
				if listErr != nil {
					return
				}
				// always first check length, otherwise its failed anyway
				assert.Len(ttt, got.Result, len(tt.want.Result))
				for i := range tt.want.Result {
					assert.Contains(ttt, got.Result, tt.want.Result[i])
				}
				integration.AssertListDetails(t, tt.want, got)
			}, retryDuration, time.Millisecond*100, "timeout waiting for expected idplinks result")
		})
	}
}
