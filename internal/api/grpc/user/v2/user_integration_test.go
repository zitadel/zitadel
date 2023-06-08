//go:build integration

package user_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/repository/idp"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX    context.Context
	ErrCTX context.Context
	Tester *integration.Tester
	Client user.UserServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(time.Hour)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX, ErrCTX = Tester.WithSystemAuthorization(ctx, integration.OrgOwner), errCtx
		Client = Tester.Client.UserV2
		return m.Run()
	}())
}

func createProvider(t *testing.T) string {
	ctx := authz.WithInstance(context.Background(), Tester.Instance)
	id, _, err := Tester.Commands.AddOrgGenericOAuthProvider(ctx, Tester.Organisation.ID, command.GenericOAuthProvider{
		"idp",
		"clientID",
		"clientSecret",
		"https://example.com/oauth/v2/authorize",
		"https://example.com/oauth/v2/token",
		"https://api.example.com/user",
		[]string{"openid", "profile", "email"},
		"id",
		idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	require.NoError(t, err)
	return id
}

func createIntent(t *testing.T, idpID string) string {
	ctx := authz.WithInstance(context.Background(), Tester.Instance)
	id, _, err := Tester.Commands.CreateIntent(ctx, idpID, "https://example.com/success", "https://example.com/failure", Tester.Organisation.ID)
	require.NoError(t, err)
	return id
}

func createSuccessfulIntent(t *testing.T, idpID string) (string, string, time.Time, uint64) {
	ctx := authz.WithInstance(context.Background(), Tester.Instance)
	intentID := createIntent(t, idpID)
	writeModel, err := Tester.Commands.GetIntentWriteModel(ctx, intentID, Tester.Organisation.ID)
	require.NoError(t, err)
	idpUser := &oauth.UserMapper{
		RawInfo: map[string]interface{}{
			"id": "id",
		},
	}
	idpSession := &oauth.Session{
		Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
			Token: &oauth2.Token{
				AccessToken: "accessToken",
			},
			IDToken: "idToken",
		},
	}
	token, err := Tester.Commands.SucceedIDPIntent(ctx, writeModel, idpUser, idpSession, "")
	require.NoError(t, err)
	return intentID, token, writeModel.ChangeDate, writeModel.ProcessedSequence
}

func TestServer_AddHumanUser(t *testing.T) {
	idpID := createProvider(t)
	type args struct {
		ctx context.Context
		req *user.AddHumanUserRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.AddHumanUserResponse
		wantErr bool
	}{
		{
			name: "default verification",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "return verification code",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_ReturnCode{
							ReturnCode: &user.ReturnEmailVerificationCode{},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
				EmailCode: gu.Ptr("something"),
			},
		},
		{
			name: "custom template",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{
							SendCode: &user.SendEmailVerificationCode{
								UrlTemplate: gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
							},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "custom template error",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{
							SendCode: &user.SendEmailVerificationCode{
								UrlTemplate: gu.Ptr("{{"),
							},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing REQUIRED profile",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_ReturnCode{
							ReturnCode: &user.ReturnEmailVerificationCode{},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing REQUIRED email",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing idp",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Email: "livio@zitadel.com",
						Verification: &user.SetHumanEmail_IsVerified{
							IsVerified: true,
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: false,
						},
					},
					IdpLinks: []*user.IDPLink{
						{
							IdpId:         "idpID",
							IdpExternalId: "externalID",
							DisplayName:   "displayName",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "with idp",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Email: "livio@zitadel.com",
						Verification: &user.SetHumanEmail_IsVerified{
							IsVerified: true,
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: false,
						},
					},
					IdpLinks: []*user.IDPLink{
						{
							IdpId:         idpID,
							IdpExternalId: "externalID",
							DisplayName:   "displayName",
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := fmt.Sprint(time.Now().UnixNano() + int64(i))
			tt.args.req.UserId = &userID
			if email := tt.args.req.GetEmail(); email != nil {
				email.Email = fmt.Sprintf("%s@me.now", userID)
			}

			if tt.want != nil {
				tt.want.UserId = userID
			}

			got, err := Client.AddHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.GetUserId(), got.GetUserId())
			if tt.want.GetEmailCode() != "" {
				assert.NotEmpty(t, got.GetEmailCode())
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_AddIDPLink(t *testing.T) {
	idpID := createProvider(t)
	type args struct {
		ctx context.Context
		req *user.AddIDPLinkRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.AddIDPLinkResponse
		wantErr bool
	}{
		{
			name: "user does not exist",
			args: args{
				CTX,
				&user.AddIDPLinkRequest{
					UserId: "userID",
					IdpLink: &user.IDPLink{
						IdpId:         idpID,
						IdpExternalId: "externalID",
						DisplayName:   "displayName",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "idp does not exist",
			args: args{
				CTX,
				&user.AddIDPLinkRequest{
					UserId: Tester.Users[integration.OrgOwner].ID,
					IdpLink: &user.IDPLink{
						IdpId:         "idpID",
						IdpExternalId: "externalID",
						DisplayName:   "displayName",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "add link",
			args: args{
				CTX,
				&user.AddIDPLinkRequest{
					UserId: Tester.Users[integration.OrgOwner].ID,
					IdpLink: &user.IDPLink{
						IdpId:         idpID,
						IdpExternalId: "externalID",
						DisplayName:   "displayName",
					},
				},
			},
			want: &user.AddIDPLinkResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddIDPLink(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_StartIdentityProviderFlow(t *testing.T) {
	idpID := createProvider(t)
	type args struct {
		ctx context.Context
		req *user.StartIdentityProviderFlowRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.StartIdentityProviderFlowResponse
		wantErr bool
	}{
		{
			name: "missing urls",
			args: args{
				CTX,
				&user.StartIdentityProviderFlowRequest{
					IdpId: idpID,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "next step auth url",
			args: args{
				CTX,
				&user.StartIdentityProviderFlowRequest{
					IdpId:      idpID,
					SuccessUrl: "https://example.com/success",
					FailureUrl: "https://example.com/failure",
				},
			},
			want: &user.StartIdentityProviderFlowResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
				NextStep: &user.StartIdentityProviderFlowResponse_AuthUrl{
					AuthUrl: "https://example.com/oauth/v2/authorize?client_id=clientID&prompt=select_account&redirect_uri=https%3A%2F%2Flocalhost%3A8080%2Fidps%2Fcallback&response_type=code&scope=openid+profile+email&state=",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.StartIdentityProviderFlow(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if nextStep := tt.want.GetNextStep(); nextStep != nil {
				if !strings.HasPrefix(got.GetAuthUrl(), tt.want.GetAuthUrl()) {
					assert.Failf(t, "auth url does not match", "expected: %s, but got: %s", tt.want.GetAuthUrl(), got.GetAuthUrl())
				}
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_RetrieveIdentityProviderInformation(t *testing.T) {
	idpID := createProvider(t)
	intentID := createIntent(t, idpID)
	successfulID, token, changeDate, sequence := createSuccessfulIntent(t, idpID)
	type args struct {
		ctx context.Context
		req *user.RetrieveIdentityProviderInformationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RetrieveIdentityProviderInformationResponse
		wantErr bool
	}{
		{
			name: "failed intent",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderInformationRequest{
					IntentId: intentID,
					Token:    "",
				},
			},
			wantErr: true,
		},
		{
			name: "wrong token",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderInformationRequest{
					IntentId: successfulID,
					Token:    "wrong token",
				},
			},
			wantErr: true,
		},
		{
			name: "retrieve successful intent",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderInformationRequest{
					IntentId: successfulID,
					Token:    token,
				},
			},
			want: &user.RetrieveIdentityProviderInformationResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(changeDate),
					ResourceOwner: Tester.Organisation.ID,
					Sequence:      sequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpInformation: []byte(`{"RawInfo":{"id":"id"}}`),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RetrieveIdentityProviderInformation(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want.GetDetails(), got.GetDetails())
			require.Equal(t, tt.want.GetIdpInformation(), got.GetIdpInformation())
		})
	}
}
