//go:build integration

package user_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/sink"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServer_StartIdentityProviderIntent(t *testing.T) {
	idpResp := Instance.AddGenericOAuthProvider(IamCTX, Instance.DefaultOrg.Id)
	orgIdpResp := Instance.AddOrgGenericOAuthProvider(OrgCTX, Instance.DefaultOrg.Id)
	orgResp := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	notDefaultOrgIdpResp := Instance.AddOrgGenericOAuthProvider(IamCTX, orgResp.OrganizationId)
	samlIdpID := Instance.AddSAMLProvider(IamCTX)
	samlRedirectIdpID := Instance.AddSAMLRedirectProvider(IamCTX, "")
	samlPostIdpID := Instance.AddSAMLPostProvider(IamCTX)
	jwtIdPID := Instance.AddJWTProvider(IamCTX)
	type args struct {
		ctx context.Context
		req *user.StartIdentityProviderIntentRequest
	}
	type want struct {
		details            *object.Details
		url                string
		parametersExisting []string
		parametersEqual    map[string]string
		postForm           bool
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "missing urls",
			args: args{
				OrgCTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: idpResp.Id,
				},
			},
			wantErr: true,
		},
		{
			name: "next step oauth auth url",
			args: args{
				OrgCTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: idpResp.Id,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url: "https://example.com/oauth/v2/authorize",
				parametersEqual: map[string]string{
					"client_id":     "clientID",
					"prompt":        "select_account",
					"redirect_uri":  "http://" + Instance.Domain + ":8082/idps/callback",
					"response_type": "code",
					"scope":         "openid profile email",
				},
				parametersExisting: []string{"state"},
			},
			wantErr: false,
		},
		{
			name: "next step oauth auth url, default org",
			args: args{
				OrgCTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: orgIdpResp.Id,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url: "https://example.com/oauth/v2/authorize",
				parametersEqual: map[string]string{
					"client_id":     "clientID",
					"prompt":        "select_account",
					"redirect_uri":  "http://" + Instance.Domain + ":8082/idps/callback",
					"response_type": "code",
					"scope":         "openid profile email",
				},
				parametersExisting: []string{"state"},
			},
			wantErr: false,
		},
		{
			name: "next step oauth auth url, default org",
			args: args{
				OrgCTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: notDefaultOrgIdpResp.Id,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url: "https://example.com/oauth/v2/authorize",
				parametersEqual: map[string]string{
					"client_id":     "clientID",
					"prompt":        "select_account",
					"redirect_uri":  "http://" + Instance.Domain + ":8082/idps/callback",
					"response_type": "code",
					"scope":         "openid profile email",
				},
				parametersExisting: []string{"state"},
			},
			wantErr: false,
		},
		{
			name: "next step oauth auth url org",
			args: args{
				OrgCTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: orgIdpResp.Id,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url: "https://example.com/oauth/v2/authorize",
				parametersEqual: map[string]string{
					"client_id":     "clientID",
					"prompt":        "select_account",
					"redirect_uri":  "http://" + Instance.Domain + ":8082/idps/callback",
					"response_type": "code",
					"scope":         "openid profile email",
				},
				parametersExisting: []string{"state"},
			},
			wantErr: false,
		},
		{
			name: "next step saml default",
			args: args{
				OrgCTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: samlIdpID,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url:                "http://localhost:8000/sso",
				parametersExisting: []string{"RelayState", "SAMLRequest"},
			},
			wantErr: false,
		},
		{
			name: "next step saml auth url",
			args: args{
				OrgCTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: samlRedirectIdpID,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url:                "http://localhost:8000/sso",
				parametersExisting: []string{"RelayState", "SAMLRequest"},
			},
			wantErr: false,
		},
		{
			name: "next step saml form",
			args: args{
				OrgCTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: samlPostIdpID,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url:                "http://localhost:8000/sso",
				parametersExisting: []string{"RelayState", "SAMLRequest"},
				postForm:           true,
			},
			wantErr: false,
		},
		{
			name: "next step jwt idp",
			args: args{
				OrgCTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: jwtIdPID,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url:                "https://example.com/jwt",
				parametersExisting: []string{"authRequestID", "userAgentID"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.StartIdentityProviderIntent(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.want.url != "" && !tt.want.postForm {
				authUrl, err := url.Parse(got.GetAuthUrl())
				require.NoError(t, err)

				assert.Equal(t, tt.want.url, authUrl.Scheme+"://"+authUrl.Host+authUrl.Path)
				require.Len(t, authUrl.Query(), len(tt.want.parametersEqual)+len(tt.want.parametersExisting))

				for _, existing := range tt.want.parametersExisting {
					assert.True(t, authUrl.Query().Has(existing))
				}
				for key, equal := range tt.want.parametersEqual {
					assert.Equal(t, equal, authUrl.Query().Get(key))
				}
			}
			if tt.want.postForm {
				assert.Equal(t, tt.want.url, got.GetFormData().GetUrl())

				require.Len(t, got.GetFormData().GetFields(), len(tt.want.parametersEqual)+len(tt.want.parametersExisting))
				for _, existing := range tt.want.parametersExisting {
					assert.Contains(t, got.GetFormData().GetFields(), existing)
				}
				for key, equal := range tt.want.parametersEqual {
					assert.Equal(t, got.GetFormData().GetFields()[key], equal)
				}
			}
			integration.AssertDetails(t, &user.StartIdentityProviderIntentResponse{
				Details: tt.want.details,
			}, got)
		})
	}
}

func TestServer_RetrieveIdentityProviderIntent(t *testing.T) {
	oauthIdpID := Instance.AddGenericOAuthProvider(IamCTX, integration.IDPName()).GetId()
	azureIdpID := Instance.AddAzureADProvider(IamCTX, integration.IDPName()).GetId()
	oidcIdpID := Instance.AddGenericOIDCProvider(IamCTX, integration.IDPName()).GetId()
	samlIdpID := Instance.AddSAMLPostProvider(IamCTX)
	ldapIdpID := Instance.AddLDAPProvider(IamCTX)
	jwtIdPID := Instance.AddJWTProvider(IamCTX)
	authURL, err := url.Parse(Instance.CreateIntent(OrgCTX, oauthIdpID).GetAuthUrl())
	require.NoError(t, err)
	intentID := authURL.Query().Get("state")
	expiry := time.Now().Add(1 * time.Hour)
	expiryFormatted := expiry.Round(time.Millisecond).UTC().Format("2006-01-02T15:04:05.999Z07:00")

	intentUser := Instance.CreateHumanUser(IamCTX)
	_, err = Instance.CreateUserIDPlink(IamCTX, intentUser.GetUserId(), "idpUserID", oauthIdpID, "username")
	require.NoError(t, err)

	successfulID, token, changeDate, sequence, err := sink.SuccessfulOAuthIntent(Instance.ID(), oauthIdpID, "id", "", expiry)
	require.NoError(t, err)
	successfulWithUserID, withUsertoken, withUserchangeDate, withUsersequence, err := sink.SuccessfulOAuthIntent(Instance.ID(), oauthIdpID, "id", "user", expiry)
	require.NoError(t, err)
	successfulExpiredID, expiredToken, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), oauthIdpID, "id", "user", time.Now().Add(time.Second))
	require.NoError(t, err)
	// make sure the intent is expired
	time.Sleep(2 * time.Second)
	successfulConsumedID, consumedToken, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), oauthIdpID, "idpUserID", intentUser.GetUserId(), expiry)
	require.NoError(t, err)
	// make sure the intent is consumed
	Instance.CreateIntentSession(t, IamCTX, intentUser.GetUserId(), successfulConsumedID, consumedToken)

	azureADSuccessful, azureADToken, azureADChangeDate, azureADSequence, err := sink.SuccessfulAzureADIntent(Instance.ID(), azureIdpID, "id", "", expiry)
	require.NoError(t, err)
	azureADSuccessfulWithUserID, azureADWithUserIDToken, azureADWithUserIDChangeDate, azureADWithUserIDSequence, err := sink.SuccessfulAzureADIntent(Instance.ID(), azureIdpID, "id", "user", expiry)
	require.NoError(t, err)

	oidcSuccessful, oidcToken, oidcChangeDate, oidcSequence, err := sink.SuccessfulOIDCIntent(Instance.ID(), oidcIdpID, "id", "", expiry)
	require.NoError(t, err)
	oidcSuccessfulWithUserID, oidcWithUserIDToken, oidcWithUserIDChangeDate, oidcWithUserIDSequence, err := sink.SuccessfulOIDCIntent(Instance.ID(), oidcIdpID, "id", "user", expiry)
	require.NoError(t, err)

	ldapSuccessfulID, ldapToken, ldapChangeDate, ldapSequence, err := sink.SuccessfulLDAPIntent(Instance.ID(), ldapIdpID, "id", "")
	require.NoError(t, err)
	ldapSuccessfulWithUserID, ldapWithUserToken, ldapWithUserChangeDate, ldapWithUserSequence, err := sink.SuccessfulLDAPIntent(Instance.ID(), ldapIdpID, "id", "user")
	require.NoError(t, err)

	samlSuccessfulID, samlToken, samlChangeDate, samlSequence, err := sink.SuccessfulSAMLIntent(Instance.ID(), samlIdpID, "id", "", expiry)
	require.NoError(t, err)
	samlSuccessfulWithUserID, samlWithUserToken, samlWithUserChangeDate, samlWithUserSequence, err := sink.SuccessfulSAMLIntent(Instance.ID(), samlIdpID, "id", "user", expiry)
	require.NoError(t, err)

	jwtSuccessfulID, jwtToken, jwtChangeDate, jwtSequence, err := sink.SuccessfulJWTIntent(Instance.ID(), jwtIdPID, "id", "", expiry)
	require.NoError(t, err)
	jwtSuccessfulWithUserID, jwtWithUserToken, jwtWithUserChangeDate, jwtWithUserSequence, err := sink.SuccessfulJWTIntent(Instance.ID(), jwtIdPID, "id", "user", expiry)
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.RetrieveIdentityProviderIntentRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RetrieveIdentityProviderIntentResponse
		wantErr bool
	}{
		{
			name: "failed intent",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    intentID,
					IdpIntentToken: "",
				},
			},
			wantErr: true,
		},
		{
			name: "wrong token",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulID,
					IdpIntentToken: "wrong token",
				},
			},
			wantErr: true,
		},
		{
			name: "retrieve successful oauth intent",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulID,
					IdpIntentToken: token,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(changeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      sequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    oauthIdpID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"RawInfo": map[string]interface{}{
								"id":                 "id",
								"preferred_username": "username",
							},
						})
						require.NoError(t, err)
						return s
					}(),
				},
				AddHumanUser: &user.AddHumanUserRequest{
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
					},
					IdpLinks: []*user.IDPLink{
						{IdpId: oauthIdpID, UserId: "id"},
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{SendCode: &user.SendEmailVerificationCode{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful intent with linked user",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulWithUserID,
					IdpIntentToken: withUsertoken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(withUserchangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      withUsersequence,
				},
				UserId: "user",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    oauthIdpID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"RawInfo": map[string]interface{}{
								"id":                 "id",
								"preferred_username": "username",
							},
						})
						require.NoError(t, err)
						return s
					}(),
				},
				UpdateHumanUser: &user.UpdateHumanUserRequest{
					UserId: "user",
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful expired intent",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulExpiredID,
					IdpIntentToken: expiredToken,
				},
			},
			wantErr: true,
		},
		{
			name: "retrieve successful consumed intent",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulConsumedID,
					IdpIntentToken: consumedToken,
				},
			},
			wantErr: true,
		},
		{
			name: "retrieve successful azure AD intent",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    azureADSuccessful,
					IdpIntentToken: azureADToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(azureADChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      azureADSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    azureIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id":                "id",
							"userPrincipalName": "username",
							"displayName":       "displayname",
							"givenName":         "firstname",
							"surname":           "lastname",
							"mail":              "email@email.com",
						})
						require.NoError(t, err)
						return s
					}(),
				},
				AddHumanUser: &user.AddHumanUserRequest{
					Username: gu.Ptr("username"),
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
						GivenName:         "firstname",
						FamilyName:        "lastname",
						DisplayName:       gu.Ptr("displayname"),
					},
					IdpLinks: []*user.IDPLink{
						{IdpId: azureIdpID, UserId: "id", UserName: "username"},
					},
					Email: &user.SetHumanEmail{
						Email:        "email@email.com",
						Verification: &user.SetHumanEmail_SendCode{SendCode: &user.SendEmailVerificationCode{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful azure AD intent with user ID",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    azureADSuccessfulWithUserID,
					IdpIntentToken: azureADWithUserIDToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(azureADWithUserIDChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      azureADWithUserIDSequence,
				},
				UserId: "user",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    azureIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id":                "id",
							"userPrincipalName": "username",
							"displayName":       "displayname",
							"givenName":         "firstname",
							"surname":           "lastname",
							"mail":              "email@email.com",
						})
						require.NoError(t, err)
						return s
					}(),
				},
				UpdateHumanUser: &user.UpdateHumanUserRequest{
					Username: gu.Ptr("username"),
					UserId:   "user",
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
						GivenName:         "firstname",
						FamilyName:        "lastname",
						DisplayName:       gu.Ptr("displayname"),
					},
					Email: &user.SetHumanEmail{
						Email: "email@email.com",
						Verification: &user.SetHumanEmail_SendCode{SendCode: &user.SendEmailVerificationCode{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful oidc intent",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    oidcSuccessful,
					IdpIntentToken: oidcToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(oidcChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      oidcSequence,
				},
				UserId: "",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    oidcIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"sub":                "id",
							"preferred_username": "username",
						})
						require.NoError(t, err)
						return s
					}(),
				},
				AddHumanUser: &user.AddHumanUserRequest{
					Username: gu.Ptr("username"),
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
					},
					IdpLinks: []*user.IDPLink{
						{IdpId: oidcIdpID, UserId: "id", UserName: "username"},
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{SendCode: &user.SendEmailVerificationCode{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful oidc intent with linked user",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    oidcSuccessfulWithUserID,
					IdpIntentToken: oidcWithUserIDToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(oidcWithUserIDChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      oidcWithUserIDSequence,
				},
				UserId: "user",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    oidcIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"sub":                "id",
							"preferred_username": "username",
						})
						require.NoError(t, err)
						return s
					}(),
				},
				UpdateHumanUser: &user.UpdateHumanUserRequest{
					Username: gu.Ptr("username"),
					UserId:   "user",
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful ldap intent",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    ldapSuccessfulID,
					IdpIntentToken: ldapToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(ldapChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      ldapSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Ldap{
						Ldap: &user.IDPLDAPAccessInformation{
							Attributes: func() *structpb.Struct {
								s, err := structpb.NewStruct(map[string]interface{}{
									"id":       []interface{}{"id"},
									"username": []interface{}{"username"},
									"language": []interface{}{"en"},
								})
								require.NoError(t, err)
								return s
							}(),
						},
					},
					IdpId:    ldapIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id":                "id",
							"preferredUsername": "username",
							"preferredLanguage": "en",
						})
						require.NoError(t, err)
						return s
					}(),
				},
				AddHumanUser: &user.AddHumanUserRequest{
					Username: gu.Ptr("username"),
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("en"),
					},
					IdpLinks: []*user.IDPLink{
						{IdpId: ldapIdpID, UserId: "id", UserName: "username"},
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{SendCode: &user.SendEmailVerificationCode{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful ldap intent with linked user",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    ldapSuccessfulWithUserID,
					IdpIntentToken: ldapWithUserToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(ldapWithUserChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      ldapWithUserSequence,
				},
				UserId: "user",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Ldap{
						Ldap: &user.IDPLDAPAccessInformation{
							Attributes: func() *structpb.Struct {
								s, err := structpb.NewStruct(map[string]interface{}{
									"id":       []interface{}{"id"},
									"username": []interface{}{"username"},
									"language": []interface{}{"en"},
								})
								require.NoError(t, err)
								return s
							}(),
						},
					},
					IdpId:    ldapIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id":                "id",
							"preferredUsername": "username",
							"preferredLanguage": "en",
						})
						require.NoError(t, err)
						return s
					}(),
				},
				UpdateHumanUser: &user.UpdateHumanUserRequest{
					Username: gu.Ptr("username"),
					UserId:   "user",
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("en"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful saml intent",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    samlSuccessfulID,
					IdpIntentToken: samlToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(samlChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      samlSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Saml{
						Saml: &user.IDPSAMLAccessInformation{
							Assertion: []byte(fmt.Sprintf(`<Assertion xmlns="urn:oasis:names:tc:SAML:2.0:assertion" ID="id" IssueInstant="0001-01-01T00:00:00Z" Version=""><Issuer xmlns="urn:oasis:names:tc:SAML:2.0:assertion" NameQualifier="" SPNameQualifier="" Format="" SPProvidedID=""></Issuer><Conditions NotBefore="0001-01-01T00:00:00Z" NotOnOrAfter="%s"></Conditions></Assertion>`, expiryFormatted)),
						},
					},
					IdpId:    samlIdpID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id": "id",
							"attributes": map[string]interface{}{
								"attribute1": []interface{}{"value1"},
							},
						})
						require.NoError(t, err)
						return s
					}(),
				},
				AddHumanUser: &user.AddHumanUserRequest{
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
					},
					IdpLinks: []*user.IDPLink{
						{IdpId: samlIdpID, UserId: "id"},
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{SendCode: &user.SendEmailVerificationCode{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful saml intent with linked user",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    samlSuccessfulWithUserID,
					IdpIntentToken: samlWithUserToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(samlWithUserChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      samlWithUserSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Saml{
						Saml: &user.IDPSAMLAccessInformation{
							Assertion: []byte(fmt.Sprintf(`<Assertion xmlns="urn:oasis:names:tc:SAML:2.0:assertion" ID="id" IssueInstant="0001-01-01T00:00:00Z" Version=""><Issuer xmlns="urn:oasis:names:tc:SAML:2.0:assertion" NameQualifier="" SPNameQualifier="" Format="" SPProvidedID=""></Issuer><Conditions NotBefore="0001-01-01T00:00:00Z" NotOnOrAfter="%s"></Conditions></Assertion>`, expiryFormatted)),
						},
					},
					IdpId:    samlIdpID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id": "id",
							"attributes": map[string]interface{}{
								"attribute1": []interface{}{"value1"},
							},
						})
						require.NoError(t, err)
						return s
					}(),
				},
				UserId: "user",
				UpdateHumanUser: &user.UpdateHumanUserRequest{
					UserId: "user",
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful jwt intent",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    jwtSuccessfulID,
					IdpIntentToken: jwtToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(jwtChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      jwtSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							IdToken: gu.Ptr("idToken"),
						},
					},
					IdpId:    jwtIdPID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"sub": "id",
						})
						require.NoError(t, err)
						return s
					}(),
				},
				AddHumanUser: &user.AddHumanUserRequest{
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
					},
					IdpLinks: []*user.IDPLink{
						{IdpId: jwtIdPID, UserId: "id"},
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{SendCode: &user.SendEmailVerificationCode{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful jwt intent with linked user",
			args: args{
				OrgCTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    jwtSuccessfulWithUserID,
					IdpIntentToken: jwtWithUserToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(jwtWithUserChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      jwtWithUserSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							IdToken: gu.Ptr("idToken"),
						},
					},
					IdpId:    jwtIdPID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"sub": "id",
						})
						require.NoError(t, err)
						return s
					}(),
				},
				UserId: "user",
				UpdateHumanUser: &user.UpdateHumanUserRequest{
					UserId: "user",
					Profile: &user.SetHumanProfile{
						PreferredLanguage: gu.Ptr("und"),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RetrieveIdentityProviderIntent(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}