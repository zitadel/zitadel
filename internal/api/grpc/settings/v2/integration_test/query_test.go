//go:build integration

package settings_test

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
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp/v2"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func TestServer_GetSecuritySettings(t *testing.T) {
	_, err := Client.SetSecuritySettings(AdminCTX, &settings.SetSecuritySettingsRequest{
		EmbeddedIframe: &settings.EmbeddedIframeSettings{
			Enabled:        true,
			AllowedOrigins: []string{"foo", "bar"},
		},
		EnableImpersonation: true,
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		ctx     context.Context
		want    *settings.GetSecuritySettingsResponse
		wantErr bool
	}{
		{
			name:    "permission error",
			ctx:     Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			wantErr: true,
		},
		{
			name: "success",
			ctx:  AdminCTX,
			want: &settings.GetSecuritySettingsResponse{
				Settings: &settings.SecuritySettings{
					EmbeddedIframe: &settings.EmbeddedIframeSettings{
						Enabled:        true,
						AllowedOrigins: []string{"foo", "bar"},
					},
					EnableImpersonation: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.ctx, time.Minute)
			assert.EventuallyWithT(t, func(ct *assert.CollectT) {
				resp, err := Client.GetSecuritySettings(tt.ctx, &settings.GetSecuritySettingsRequest{})
				if tt.wantErr {
					assert.Error(ct, err)
					return
				}
				if !assert.NoError(ct, err) {
					return
				}
				got, want := resp.GetSettings(), tt.want.GetSettings()
				assert.Equal(ct, want.GetEmbeddedIframe().GetEnabled(), got.GetEmbeddedIframe().GetEnabled(), "enable iframe embedding")
				assert.Equal(ct, want.GetEmbeddedIframe().GetAllowedOrigins(), got.GetEmbeddedIframe().GetAllowedOrigins(), "allowed origins")
				assert.Equal(ct, want.GetEnableImpersonation(), got.GetEnableImpersonation(), "enable impersonation")
			}, retryDuration, tick)
		})
	}
}

func idpResponse(id, name string, linking, creation, autoCreation, autoUpdate bool, autoLinking idp_pb.AutoLinkingOption) *settings.IdentityProvider {
	return &settings.IdentityProvider{
		Id:   id,
		Name: name,
		Type: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_OAUTH,
		Options: &idp_pb.Options{
			IsLinkingAllowed:  linking,
			IsCreationAllowed: creation,
			IsAutoCreation:    autoCreation,
			IsAutoUpdate:      autoUpdate,
			AutoLinking:       autoLinking,
		},
	}
}

func TestServer_GetActiveIdentityProviders(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	instance.AddGenericOAuthProvider(isolatedIAMOwnerCTX, integration.IDPName()) // inactive
	idpActiveName := integration.IDPName()
	idpActiveResp := instance.AddGenericOAuthProvider(isolatedIAMOwnerCTX, idpActiveName)
	instance.AddProviderToDefaultLoginPolicy(isolatedIAMOwnerCTX, idpActiveResp.GetId())
	idpActiveResponse := idpResponse(idpActiveResp.GetId(), idpActiveName, true, true, true, true, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	idpLinkingDisallowedName := integration.IDPName()
	idpLinkingDisallowedResp := instance.AddGenericOAuthProviderWithOptions(isolatedIAMOwnerCTX, idpLinkingDisallowedName, false, true, true, idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	instance.AddProviderToDefaultLoginPolicy(isolatedIAMOwnerCTX, idpLinkingDisallowedResp.GetId())
	idpLinkingDisallowedResponse := idpResponse(idpLinkingDisallowedResp.GetId(), idpLinkingDisallowedName, false, true, true, true, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	idpCreationDisallowedName := integration.IDPName()
	idpCreationDisallowedResp := instance.AddGenericOAuthProviderWithOptions(isolatedIAMOwnerCTX, idpCreationDisallowedName, true, false, true, idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	instance.AddProviderToDefaultLoginPolicy(isolatedIAMOwnerCTX, idpCreationDisallowedResp.GetId())
	idpCreationDisallowedResponse := idpResponse(idpCreationDisallowedResp.GetId(), idpCreationDisallowedName, true, false, true, true, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	idpNoAutoCreationName := integration.IDPName()
	idpNoAutoCreationResp := instance.AddGenericOAuthProviderWithOptions(isolatedIAMOwnerCTX, idpNoAutoCreationName, true, true, false, idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	instance.AddProviderToDefaultLoginPolicy(isolatedIAMOwnerCTX, idpNoAutoCreationResp.GetId())
	idpNoAutoCreationResponse := idpResponse(idpNoAutoCreationResp.GetId(), idpNoAutoCreationName, true, true, false, true, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	idpNoAutoLinkingName := integration.IDPName()
	idpNoAutoLinkingResp := instance.AddGenericOAuthProviderWithOptions(isolatedIAMOwnerCTX, idpNoAutoLinkingName, true, true, true, idp.AutoLinkingOption_AUTO_LINKING_OPTION_UNSPECIFIED)
	instance.AddProviderToDefaultLoginPolicy(isolatedIAMOwnerCTX, idpNoAutoLinkingResp.GetId())
	idpNoAutoLinkingResponse := idpResponse(idpNoAutoLinkingResp.GetId(), idpNoAutoLinkingName, true, true, true, true, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_UNSPECIFIED)

	type args struct {
		ctx context.Context
		req *settings.GetActiveIdentityProvidersRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *settings.GetActiveIdentityProvidersResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &settings.GetActiveIdentityProvidersRequest{},
			},
			wantErr: true,
		},
		{
			name: "success, all",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 5,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpActiveResponse,
					idpLinkingDisallowedResponse,
					idpCreationDisallowedResponse,
					idpNoAutoCreationResponse,
					idpNoAutoLinkingResponse,
				},
			},
		},
		{
			name: "success, exclude linking disallowed",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{
					LinkingAllowed: gu.Ptr(true),
				},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 4,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpActiveResponse,
					idpCreationDisallowedResponse,
					idpNoAutoCreationResponse,
					idpNoAutoLinkingResponse,
				},
			},
		},
		{
			name: "success, only linking disallowed",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{
					LinkingAllowed: gu.Ptr(false),
				},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpLinkingDisallowedResponse,
				},
			},
		},
		{
			name: "success, exclude creation disallowed",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{
					CreationAllowed: gu.Ptr(true),
				},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 4,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpActiveResponse,
					idpLinkingDisallowedResponse,
					idpNoAutoCreationResponse,
					idpNoAutoLinkingResponse,
				},
			},
		},
		{
			name: "success, only creation disallowed",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{
					CreationAllowed: gu.Ptr(false),
				},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpCreationDisallowedResponse,
				},
			},
		},
		{
			name: "success, auto creation",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{
					AutoCreation: gu.Ptr(true),
				},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 4,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpActiveResponse,
					idpLinkingDisallowedResponse,
					idpCreationDisallowedResponse,
					idpNoAutoLinkingResponse,
				},
			},
		},
		{
			name: "success, no auto creation",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{
					AutoCreation: gu.Ptr(false),
				},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpNoAutoCreationResponse,
				},
			},
		},
		{
			name: "success, auto linking",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{
					AutoLinking: gu.Ptr(true),
				},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 4,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpActiveResponse,
					idpLinkingDisallowedResponse,
					idpCreationDisallowedResponse,
					idpNoAutoCreationResponse,
				},
			},
		},
		{
			name: "success, no auto linking",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{
					AutoLinking: gu.Ptr(false),
				},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpNoAutoLinkingResponse,
				},
			},
		},
		{
			name: "success, exclude all",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &settings.GetActiveIdentityProvidersRequest{
					LinkingAllowed:  gu.Ptr(true),
					CreationAllowed: gu.Ptr(true),
					AutoCreation:    gu.Ptr(true),
					AutoLinking:     gu.Ptr(true),
				},
			},
			want: &settings.GetActiveIdentityProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				IdentityProviders: []*settings.IdentityProvider{
					idpActiveResponse,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			assert.EventuallyWithT(t, func(ct *assert.CollectT) {
				got, err := instance.Client.SettingsV2.GetActiveIdentityProviders(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ct, err)
					return
				}
				if !assert.NoError(ct, err) {
					return
				}
				assert.ElementsMatch(ct, tt.want.GetIdentityProviders(), got.GetIdentityProviders())
				integration.AssertListDetails(ct, tt.want, got)
			}, retryDuration, tick)
		})
	}
}

func TestServer_GetHostedLoginTranslation(t *testing.T) {
	// Given
	translations := map[string]any{"loginTitle": integration.Slogan()}

	protoTranslations, err := structpb.NewStruct(translations)
	require.NoError(t, err)

	setupRequest := &settings.SetHostedLoginTranslationRequest{
		Level: &settings.SetHostedLoginTranslationRequest_OrganizationId{
			OrganizationId: Instance.DefaultOrg.GetId(),
		},
		Translations: protoTranslations,
		Locale:       integration.Language(),
	}
	savedTranslation, err := Client.SetHostedLoginTranslation(AdminCTX, setupRequest)
	require.NoError(t, err)

	tt := []struct {
		testName     string
		inputCtx     context.Context
		inputRequest *settings.GetHostedLoginTranslationRequest

		expectedErrorCode codes.Code
		expectedErrorMsg  string
		expectedResponse  *settings.GetHostedLoginTranslationResponse
	}{
		{
			testName:          "when unauthN context should return unauthN error",
			inputCtx:          CTX,
			inputRequest:      &settings.GetHostedLoginTranslationRequest{Locale: "en-US"},
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName:          "when unauthZ context should return unauthZ error",
			inputCtx:          OrgOwnerCtx,
			inputRequest:      &settings.GetHostedLoginTranslationRequest{Locale: "en-US"},
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when authZ request should save to db and return etag",
			inputCtx: AdminCTX,
			inputRequest: &settings.GetHostedLoginTranslationRequest{
				Level: &settings.GetHostedLoginTranslationRequest_OrganizationId{
					OrganizationId: Instance.DefaultOrg.GetId(),
				},
				Locale: setupRequest.GetLocale(),
			},
			expectedResponse: &settings.GetHostedLoginTranslationResponse{
				Etag:         savedTranslation.GetEtag(),
				Translations: protoTranslations,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// When
			res, err := Client.GetHostedLoginTranslation(tc.inputCtx, tc.inputRequest)

			// Then
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				require.NoError(t, err)
				assert.NotEmpty(t, res.GetEtag())
				assert.NotEmpty(t, res.GetTranslations().GetFields())
			}
		})
	}
}

func TestServer_ListOrganizationSettings(t *testing.T) {
	instance := integration.NewInstance(CTX)
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		dep func(*settings.ListOrganizationSettingsRequest, *settings.ListOrganizationSettingsResponse)
		req *settings.ListOrganizationSettingsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *settings.ListOrganizationSettingsResponse
		wantErr bool
	}{
		{
			name: "list by id, unauthenticated",
			args: args{
				ctx: CTX,
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp.GetOrganizationId()},
						},
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp.GetOrganizationId()},
						},
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{},
			},
		},
		{
			name: "list by id, missing permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp.GetOrganizationId()},
						},
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{},
			},
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,

				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{
						Filter: &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
							InOrganizationIdsFilter: &filter.InIDsFilter{
								Ids: []string{"notexisting"},
							},
						},
					}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp.GetOrganizationId()},
						},
					}
					response.OrganizationSettings[0] = &settings.OrganizationSettings{
						OrganizationId:              orgResp.GetOrganizationId(),
						CreationDate:                settingsResp.GetSetDate(),
						ChangeDate:                  settingsResp.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{{}},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp1 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp1 := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp1.GetOrganizationId(), true)
					orgResp2 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp2 := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp2.GetOrganizationId(), true)
					orgResp3 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp3 := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp3.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp1.GetOrganizationId(), orgResp2.GetOrganizationId(), orgResp3.GetOrganizationId()},
						},
					}
					response.OrganizationSettings[2] = &settings.OrganizationSettings{
						OrganizationId:              orgResp1.GetOrganizationId(),
						CreationDate:                settingsResp1.GetSetDate(),
						ChangeDate:                  settingsResp1.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
					response.OrganizationSettings[1] = &settings.OrganizationSettings{
						OrganizationId:              orgResp2.GetOrganizationId(),
						CreationDate:                settingsResp2.GetSetDate(),
						ChangeDate:                  settingsResp2.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
					response.OrganizationSettings[0] = &settings.OrganizationSettings{
						OrganizationId:              orgResp3.GetOrganizationId(),
						CreationDate:                settingsResp3.GetSetDate(),
						ChangeDate:                  settingsResp3.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{{}, {}, {}},
			},
		},
		{
			name: "list multiple id, only org scoped usernames",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp1 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp1.GetOrganizationId(), false)
					orgResp2 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp2 := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp2.GetOrganizationId(), true)
					orgResp3 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp3.GetOrganizationId(), false)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp1.GetOrganizationId(), orgResp2.GetOrganizationId(), orgResp3.GetOrganizationId()},
						},
					}
					request.Filters[1].Filter = &settings.OrganizationSettingsSearchFilter_OrganizationScopedUsernamesFilter{
						OrganizationScopedUsernamesFilter: &settings.OrganizationScopedUsernamesFilter{
							OrganizationScopedUsernames: true,
						},
					}
					response.OrganizationSettings[0] = &settings.OrganizationSettings{
						OrganizationId:              orgResp2.GetOrganizationId(),
						CreationDate:                settingsResp2.GetSetDate(),
						ChangeDate:                  settingsResp2.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}, {}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := instance.Client.SettingsV2.ListOrganizationSettings(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.OrganizationSettings, len(tt.want.OrganizationSettings)) {
					for i := range tt.want.OrganizationSettings {
						assert.EqualExportedValues(ttt, tt.want.OrganizationSettings[i], got.OrganizationSettings[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func assertPaginationResponse(t *assert.CollectT, expected *filter.PaginationResponse, actual *filter.PaginationResponse) {
	assert.Equal(t, expected.AppliedLimit, actual.AppliedLimit)
	assert.Equal(t, expected.TotalResult, actual.TotalResult)
}
