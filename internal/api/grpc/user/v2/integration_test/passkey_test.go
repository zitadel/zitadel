//go:build integration

package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_RegisterPasskey(t *testing.T) {
	userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
	reg, err := Client.CreatePasskeyRegistrationLink(OrgCTX, &user.CreatePasskeyRegistrationLinkRequest{
		UserId: userID,
		Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
	})
	require.NoError(t, err)

	// We also need a user session
	Instance.RegisterUserPasskey(OrgCTX, userID)
	_, sessionToken, _, _ := Instance.CreateVerifiedWebAuthNSession(t, LoginCTX, userID)

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
				ctx: OrgCTX,
				req: &user.RegisterPasskeyRequest{},
			},
			wantErr: true,
		},
		{
			name: "register code",
			args: args{
				ctx: OrgCTX,
				req: &user.RegisterPasskeyRequest{
					UserId:        userID,
					Code:          reg.GetCode(),
					Authenticator: user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_PLATFORM,
				},
			},
			want: &user.RegisterPasskeyResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "reuse code (not allowed)",
			args: args{
				ctx: OrgCTX,
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
				ctx: OrgCTX,
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
			name: "user no permission",
			args: args{
				ctx: UserCTX,
				req: &user.RegisterPasskeyRequest{
					UserId: userID,
				},
			},
			wantErr: true,
		},
		{
			name: "user permission",
			args: args{
				ctx: IamCTX,
				req: &user.RegisterPasskeyRequest{
					UserId: userID,
				},
			},
			want: &user.RegisterPasskeyResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "user setting its own passkey",
			args: args{
				ctx: integration.WithAuthorizationToken(OrgCTX, sessionToken),
				req: &user.RegisterPasskeyRequest{
					UserId: userID,
				},
			},
			want: &user.RegisterPasskeyResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
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
				_, err = Instance.WebAuthN.CreateAttestationResponse(got.GetPublicKeyCredentialCreationOptions())
				require.NoError(t, err)
			}
		})
	}
}

func TestServer_VerifyPasskeyRegistration(t *testing.T) {
	userID, pkr := userWithPasskeyRegistered(t)

	attestationResponse, err := Instance.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
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
				ctx: OrgCTX,
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
				ctx: OrgCTX,
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
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "wrong credential",
			args: args{
				ctx: OrgCTX,
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
	userID := Instance.CreateHumanUser(OrgCTX).GetUserId()

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
				ctx: OrgCTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{},
			},
			wantErr: true,
		},
		{
			name: "send default mail",
			args: args{
				ctx: OrgCTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{
					UserId: userID,
				},
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "send custom url",
			args: args{
				ctx: OrgCTX,
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
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "return code",
			args: args{
				ctx: OrgCTX,
				req: &user.CreatePasskeyRegistrationLinkRequest{
					UserId: userID,
					Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
				},
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
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

// TestServer_CreatePasskeyRegistrationLink_CrossOrg simulates the cross-organization passkey
// enrollment code issuance: the RPC's auth annotation carries user.passkey.write with no
// org_field, so the interceptor only checks the caller's permission in the request-header org.
// An org owner of one organization must not be able to mint an enrollment code - and with it a
// takeover of the account - for a user in another organization.
func TestServer_CreatePasskeyRegistrationLink_CrossOrg(t *testing.T) {
	// The victim lives in another organization than the caller.
	orgB := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	victimID := Instance.CreateHumanUserVerified(
		IamCTX, orgB.GetOrganizationId(), integration.Email(), integration.Phone(),
	).GetUserId()

	t.Run("org owner of another org is denied", func(t *testing.T) {
		// The attacker is ORG_OWNER of the default org only, and pins the request-header org
		// to it - exactly the scope the interceptor verifies.
		attackerCTX := integration.SetOrgID(OrgCTX, Instance.DefaultOrg.GetId())

		got, err := Client.CreatePasskeyRegistrationLink(attackerCTX, &user.CreatePasskeyRegistrationLinkRequest{
			UserId: victimID,
			Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
		})
		require.Error(t, err)
		// The default permission path reports "membership not found" (NotFound) because the
		// caller has no membership in the victim's org at all; the v2 permission check reports
		// PermissionDenied. Either is a correct denial.
		assert.Contains(t, []codes.Code{codes.NotFound, codes.PermissionDenied}, status.Code(err))
		assert.Empty(t, got.GetCode().GetCode())
	})

	t.Run("instance wide user manager is allowed", func(t *testing.T) {
		// IAM_USER_MANAGER holds user.passkey.write instance wide, so it must keep working
		// across organizations - the check narrows the org, it must not narrow the role set.
		_, pat, err := Instance.CreateMachineUserPATWithMembership(IamCTX, "IAM_USER_MANAGER")
		require.NoError(t, err)
		managerCTX := integration.WithAuthorizationToken(CTX, pat)

		got, err := Client.CreatePasskeyRegistrationLink(managerCTX, &user.CreatePasskeyRegistrationLinkRequest{
			UserId: victimID,
			Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
		})
		require.NoError(t, err)
		assert.NotEmpty(t, got.GetCode().GetId())
	})
}

func userWithPasskeyRegistered(t *testing.T) (string, *user.RegisterPasskeyResponse) {
	userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
	return userID, passkeyRegister(t, userID)
}

func userWithPasskeyVerified(t *testing.T) (string, string) {
	userID, pkr := userWithPasskeyRegistered(t)
	return userID, passkeyVerify(t, userID, pkr)
}

func passkeyRegister(t *testing.T, userID string) *user.RegisterPasskeyResponse {
	reg, err := Client.CreatePasskeyRegistrationLink(OrgCTX, &user.CreatePasskeyRegistrationLinkRequest{
		UserId: userID,
		Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
	})
	require.NoError(t, err)
	pkr, err := Client.RegisterPasskey(OrgCTX, &user.RegisterPasskeyRequest{
		UserId: userID,
		Code:   reg.GetCode(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, pkr.GetPasskeyId())
	require.NotEmpty(t, pkr.GetPublicKeyCredentialCreationOptions())
	return pkr
}

func passkeyVerify(t *testing.T, userID string, pkr *user.RegisterPasskeyResponse) string {
	attestationResponse, err := Instance.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	require.NoError(t, err)

	_, err = Client.VerifyPasskeyRegistration(OrgCTX, &user.VerifyPasskeyRegistrationRequest{
		UserId:              userID,
		PasskeyId:           pkr.GetPasskeyId(),
		PublicKeyCredential: attestationResponse,
		PasskeyName:         "nice name",
	})
	require.NoError(t, err)
	return pkr.GetPasskeyId()
}

func TestServer_RemovePasskey(t *testing.T) {
	userIDWithout := Instance.CreateHumanUser(OrgCTX).GetUserId()
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
					ResourceOwner: Instance.DefaultOrg.Id,
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
					ResourceOwner: Instance.DefaultOrg.Id,
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
	userIDWithout := Instance.CreateHumanUser(OrgCTX).GetUserId()
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
			name: "list passkeys, no userID",
			args: args{
				IamCTX,
				&user.ListPasskeysRequest{
					UserId: "",
				},
			},
			wantErr: true,
		},
		{
			name: "list passkeys, no permission",
			args: args{
				UserCTX,
				&user.ListPasskeysRequest{
					UserId: userIDVerified,
				},
			},
			wantErr: true,
		},
		{
			name: "list passkeys, none",
			args: args{
				IamCTX,
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
				IamCTX,
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
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListPasskeys(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)
				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Result, len(tt.want.Result)) {
					for i := range tt.want.Result {
						assert.Contains(ttt, got.Result, tt.want.Result[i])
					}
				}
				integration.AssertListDetails(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected idplinks result")
		})
	}
}
