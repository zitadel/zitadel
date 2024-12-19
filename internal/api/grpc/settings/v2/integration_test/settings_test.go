//go:build integration

package settings_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
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
			ctx:     Instance.WithAuthorization(CTX, integration.UserTypeOrgOwner),
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

func TestServer_SetSecuritySettings(t *testing.T) {
	type args struct {
		ctx context.Context
		req *settings.SetSecuritySettingsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *settings.SetSecuritySettingsResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: Instance.WithAuthorization(CTX, integration.UserTypeOrgOwner),
				req: &settings.SetSecuritySettingsRequest{
					EmbeddedIframe: &settings.EmbeddedIframeSettings{
						Enabled:        true,
						AllowedOrigins: []string{"foo.com", "bar.com"},
					},
					EnableImpersonation: true,
				},
			},
			wantErr: true,
		},
		{
			name: "success allowed origins",
			args: args{
				ctx: AdminCTX,
				req: &settings.SetSecuritySettingsRequest{
					EmbeddedIframe: &settings.EmbeddedIframeSettings{
						AllowedOrigins: []string{"foo.com", "bar.com"},
					},
				},
			},
			want: &settings.SetSecuritySettingsResponse{
				Details: &object_pb.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "success enable iframe embedding",
			args: args{
				ctx: AdminCTX,
				req: &settings.SetSecuritySettingsRequest{
					EmbeddedIframe: &settings.EmbeddedIframeSettings{
						Enabled: true,
					},
				},
			},
			want: &settings.SetSecuritySettingsResponse{
				Details: &object_pb.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "success impersonation",
			args: args{
				ctx: AdminCTX,
				req: &settings.SetSecuritySettingsRequest{
					EnableImpersonation: true,
				},
			},
			want: &settings.SetSecuritySettingsResponse{
				Details: &object_pb.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "success all",
			args: args{
				ctx: AdminCTX,
				req: &settings.SetSecuritySettingsRequest{
					EmbeddedIframe: &settings.EmbeddedIframeSettings{
						Enabled:        true,
						AllowedOrigins: []string{"foo.com", "bar.com"},
					},
					EnableImpersonation: true,
				},
			},
			want: &settings.SetSecuritySettingsResponse{
				Details: &object_pb.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetSecuritySettings(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
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
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	instance.AddGenericOAuthProvider(isolatedIAMOwnerCTX, gofakeit.AppName()) // inactive
	idpActiveName := gofakeit.AppName()
	idpActiveResp := instance.AddGenericOAuthProvider(isolatedIAMOwnerCTX, idpActiveName)
	instance.AddProviderToDefaultLoginPolicy(isolatedIAMOwnerCTX, idpActiveResp.GetId())
	idpActiveResponse := idpResponse(idpActiveResp.GetId(), idpActiveName, true, true, true, true, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	idpLinkingDisallowedName := gofakeit.AppName()
	idpLinkingDisallowedResp := instance.AddGenericOAuthProviderWithOptions(isolatedIAMOwnerCTX, idpLinkingDisallowedName, false, true, true, idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	instance.AddProviderToDefaultLoginPolicy(isolatedIAMOwnerCTX, idpLinkingDisallowedResp.GetId())
	idpLinkingDisallowedResponse := idpResponse(idpLinkingDisallowedResp.GetId(), idpLinkingDisallowedName, false, true, true, true, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	idpCreationDisallowedName := gofakeit.AppName()
	idpCreationDisallowedResp := instance.AddGenericOAuthProviderWithOptions(isolatedIAMOwnerCTX, idpCreationDisallowedName, true, false, true, idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	instance.AddProviderToDefaultLoginPolicy(isolatedIAMOwnerCTX, idpCreationDisallowedResp.GetId())
	idpCreationDisallowedResponse := idpResponse(idpCreationDisallowedResp.GetId(), idpCreationDisallowedName, true, false, true, true, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	idpNoAutoCreationName := gofakeit.AppName()
	idpNoAutoCreationResp := instance.AddGenericOAuthProviderWithOptions(isolatedIAMOwnerCTX, idpNoAutoCreationName, true, true, false, idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	instance.AddProviderToDefaultLoginPolicy(isolatedIAMOwnerCTX, idpNoAutoCreationResp.GetId())
	idpNoAutoCreationResponse := idpResponse(idpNoAutoCreationResp.GetId(), idpNoAutoCreationName, true, true, false, true, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
	idpNoAutoLinkingName := gofakeit.AppName()
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
				ctx: instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
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
				for i, result := range tt.want.GetIdentityProviders() {
					assert.EqualExportedValues(ct, result, got.GetIdentityProviders()[i])
				}
				integration.AssertListDetails(ct, tt.want, got)
			}, retryDuration, tick)
		})
	}
}
