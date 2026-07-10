//go:build integration

package management_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/zitadel/zitadel/internal/integration"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
)

func Test_AddZitadelProvider(t *testing.T) {
	type args struct {
		ctx context.Context
		req *mgmt_pb.AddZitadelProviderRequest
	}
	tests := []struct {
		name         string
		args         args
		wantErr      error
		wantResponse *mgmt_pb.AddZitadelProviderResponse
	}{
		{
			name: "no permissions, error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &mgmt_pb.AddZitadelProviderRequest{
					Name:         "Zitadel Support IdP",
					Issuer:       "zitadel.example.com",
					ClientId:     "test-client",
					ClientSecret: "test-secret",
					Scopes:       []string{"email", "profile"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
					},
				},
			},
			wantErr: status.Error(codes.NotFound, "membership not found (AUTHZ-cdgFk)"),
		},
		{
			name: "insufficient permissions, error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeLogin),
				req: &mgmt_pb.AddZitadelProviderRequest{
					Name:         "Zitadel Support IdP",
					Issuer:       "zitadel.example.com",
					ClientId:     "test-client",
					ClientSecret: "test-secret",
					Scopes:       []string{"email", "profile"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
					},
				},
			},
			wantErr: status.Error(codes.PermissionDenied, "No matching permissions found (AUTH-5mWD2)"),
		},
		{
			name: "missing required field: name",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddZitadelProviderRequest{},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid AddZitadelProviderRequest.Name: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: issuer",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddZitadelProviderRequest{
					Name: "Zitadel Support IdP",
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid AddZitadelProviderRequest.Issuer: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: client_id",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddZitadelProviderRequest{
					Name:   "Zitadel Support IdP",
					Issuer: "zitadel.example.com",
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid AddZitadelProviderRequest.ClientId: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: client_secret",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddZitadelProviderRequest{
					Name:     "Zitadel Support IdP",
					Issuer:   "zitadel.example.com",
					ClientId: "test-client",
				},
			},
			wantErr: status.Error(
				codes.InvalidArgument,
				"invalid AddZitadelProviderRequest.ClientSecret: value length must be between 1 and 1000 runes, inclusive",
			),
		},
		{
			name: "missing org ID in instance roles info",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddZitadelProviderRequest{
					Name:         "Zitadel Support IdP",
					Issuer:       "zitadel.example.com",
					ClientId:     "test-client",
					ClientSecret: "test-secret",
					Scopes:       []string{"email", "profile"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
					},
					InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
						{
							OrganizationId: "",
						},
					},
				},
			},
			wantErr: status.Error(
				codes.InvalidArgument,
				"invalid AddZitadelProviderRequest.InstanceRolesInfo[0]: embedded message failed validation | caused by: invalid InstanceRolesInfo.OrganizationId: value length must be between 1 and 200 runes, inclusive",
			),
		},
		{
			name: "missing org domain in instance roles info",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddZitadelProviderRequest{
					Name:         "Zitadel Support IdP",
					Issuer:       "zitadel.example.com",
					ClientId:     "test-client",
					ClientSecret: "test-secret",
					Scopes:       []string{"email", "profile"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
					},
					InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
						{
							OrganizationId:     "org1",
							OrganizationDomain: "org1.com",
						},
						{
							OrganizationId: "org2",
						},
					},
				},
			},
			wantErr: status.Error(
				codes.InvalidArgument,
				"invalid AddZitadelProviderRequest.InstanceRolesInfo[1]: embedded message failed validation | caused by: invalid InstanceRolesInfo.OrganizationDomain: value length must be between 1 and 200 runes, inclusive",
			),
		},
		{
			name: "valid request without instance roles info",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddZitadelProviderRequest{
					Name:         "Zitadel Support IdP",
					Issuer:       "zitadel.example.com",
					ClientId:     "test-client",
					ClientSecret: "test-secret",
					Scopes:       []string{"email", "profile"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
					},
				},
			},
			wantResponse: &mgmt_pb.AddZitadelProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "valid request with instance roles info",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddZitadelProviderRequest{
					Name:         "Zitadel Support IdP",
					Issuer:       "zitadel.example.com",
					ClientId:     "test-client",
					ClientSecret: "test-secret",
					Scopes:       []string{"email", "profile"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
					},
					InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
						{
							OrganizationId:     "org1",
							OrganizationDomain: "org1.com",
						},
					},
				},
			},
			wantResponse: &mgmt_pb.AddZitadelProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			got, err := Client.AddZitadelProvider(tt.args.ctx, tt.args.req)
			after := time.Now()
			if tt.wantErr != nil {
				require.Error(t, err)
				grpcStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, status.Code(tt.wantErr), grpcStatus.Code())
				assert.Equal(t, status.Convert(tt.wantErr).Message(), grpcStatus.Message())
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.NotEmpty(t, got.GetId())
			assert.WithinRange(t, got.GetDetails().GetCreationDate().AsTime(), before, after)
			assert.Equal(t, tt.wantResponse.GetDetails().GetResourceOwner(), got.GetDetails().GetResourceOwner())
		})
	}
}

func Test_UpdateZitadelProvider(t *testing.T) {
	type args struct {
		ctx context.Context
		req *mgmt_pb.UpdateZitadelProviderRequest
	}
	tests := []struct {
		name                string
		args                args
		wantErr             error
		wantResponse        *mgmt_pb.UpdateZitadelProviderResponse
		wantUpdatedProvider *idp_pb.Provider
	}{
		{
			name: "no permissions, error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &mgmt_pb.UpdateZitadelProviderRequest{
					Issuer:       "zitadel.example.com",
					ClientId:     "test-client",
					ClientSecret: "test-secret",
					Scopes:       []string{"email", "profile"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
					},
				},
			},
			wantErr: status.Error(codes.NotFound, "membership not found (AUTHZ-cdgFk)"),
		},
		{
			name: "insufficient permissions, error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeLogin),
				req: &mgmt_pb.UpdateZitadelProviderRequest{
					Issuer: "acme.example.com",
					Scopes: []string{"email", "profile"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
						IsAutoCreation:    false,
					},
				},
			},
			wantErr: status.Error(codes.PermissionDenied, "No matching permissions found (AUTH-5mWD2)"),
		},
		{
			name: "missing required field: name",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.UpdateZitadelProviderRequest{},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid UpdateZitadelProviderRequest.Name: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: issuer",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.UpdateZitadelProviderRequest{
					Name: "Zitadel Support IdP updated",
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid UpdateZitadelProviderRequest.Issuer: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: client_id",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.UpdateZitadelProviderRequest{
					Name:   "Zitadel Support IdP updated",
					Issuer: "acme.example.com",
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid UpdateZitadelProviderRequest.ClientId: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "update, ok",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.UpdateZitadelProviderRequest{
					Name:     "Zitadel Support IdP updated",
					Issuer:   "acme.example.com",
					ClientId: "test-client",
					Scopes:   []string{"email", "profile", "openid", "offline_access"},
					InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
						{
							OrganizationId:     "org1",
							OrganizationDomain: "org1.com",
						},
						{
							OrganizationId:     "org2",
							OrganizationDomain: "org2.com",
						},
					},
				},
			},
			wantResponse: &mgmt_pb.UpdateZitadelProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
			wantUpdatedProvider: &idp_pb.Provider{
				State: idp_pb.IDPState_IDP_STATE_ACTIVE,
				Name:  "Zitadel Support IdP updated",
				Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG,
				Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
				Config: &idp_pb.ProviderConfig{
					Options: &idp_pb.Options{},
					Config: &idp_pb.ProviderConfig_Zitadel{
						Zitadel: &idp_pb.ZitadelConfig{
							Issuer:   "acme.example.com",
							ClientId: "test-client",
							Scopes:   []string{"email", "profile", "openid", "offline_access"},
							InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
								{
									OrganizationId:     "org1",
									OrganizationDomain: "org1.com",
								},
								{
									OrganizationId:     "org2",
									OrganizationDomain: "org2.com",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "update with scopes unset, ok",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.UpdateZitadelProviderRequest{
					Name:     "Zitadel Support IdP updated 1",
					Issuer:   "acme.example.com",
					ClientId: "test-client",
					Scopes:   []string{}, // scopes unset -> will be updated by the API
					ProviderOptions: &idp_pb.Options{
						IsAutoCreation: true,
					},
					InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
						{
							OrganizationId:     "org3",
							OrganizationDomain: "org3.com",
						},
					},
				},
			},
			wantResponse: &mgmt_pb.UpdateZitadelProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
			wantUpdatedProvider: &idp_pb.Provider{
				State: idp_pb.IDPState_IDP_STATE_ACTIVE,
				Name:  "Zitadel Support IdP updated 1",
				Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG,
				Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
				Config: &idp_pb.ProviderConfig{
					Options: &idp_pb.Options{
						IsAutoCreation: true,
					},
					Config: &idp_pb.ProviderConfig_Zitadel{
						Zitadel: &idp_pb.ZitadelConfig{
							Issuer:   "acme.example.com",
							ClientId: "test-client",
							Scopes:   nil, // unset scopes
							InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
								{
									OrganizationId:     "org3",
									OrganizationDomain: "org3.com",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "update with instance roles info not set, ok",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.UpdateZitadelProviderRequest{
					Name:     "Zitadel Support IdP updated 2",
					Issuer:   "acme.example.com",
					ClientId: "test-client",
					Scopes:   []string{"email", "openid"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
					},
					InstanceRolesInfo: nil, // instance roles info isn't set -> will be unset by the API
				},
			},
			wantResponse: &mgmt_pb.UpdateZitadelProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
			wantUpdatedProvider: &idp_pb.Provider{
				State: idp_pb.IDPState_IDP_STATE_ACTIVE,
				Name:  "Zitadel Support IdP updated 2",
				Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG,
				Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
				Config: &idp_pb.ProviderConfig{
					Options: &idp_pb.Options{
						IsCreationAllowed: true,
					},
					Config: &idp_pb.ProviderConfig_Zitadel{
						Zitadel: &idp_pb.ZitadelConfig{
							Issuer:            "acme.example.com",
							ClientId:          "test-client",
							Scopes:            []string{"email", "openid"},
							InstanceRolesInfo: nil, // instance roles info isn't set
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a new provider per subtest and set the ID in the request
			zitadelProvider := Instance.AddOrgZitadelProvider(OrgCTX, integration.IDPName())
			tt.args.req.Id = zitadelProvider.Id

			before := time.Now()
			updateResp, err := Client.UpdateZitadelProvider(tt.args.ctx, tt.args.req)
			after := time.Now()
			if tt.wantErr != nil {
				require.Error(t, err)
				grpcStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, status.Code(tt.wantErr), grpcStatus.Code())
				assert.Equal(t, status.Convert(tt.wantErr).Message(), grpcStatus.Message())
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, updateResp)
			assert.WithinRange(t, updateResp.GetDetails().GetChangeDate().AsTime(), before, after)
			assert.Equal(t, tt.wantResponse.GetDetails().GetResourceOwner(), updateResp.GetDetails().GetResourceOwner())

			// get provider by ID and assert that the updated provider matches the expected values
			getResp, err := Client.GetProviderByID(OrgCTX, &mgmt_pb.GetProviderByIDRequest{Id: zitadelProvider.Id})
			require.NoError(t, err)
			assert.Equal(t, updateResp.GetDetails().GetChangeDate().AsTime(), getResp.GetIdp().GetDetails().GetChangeDate().AsTime())
			tt.wantUpdatedProvider.Id = zitadelProvider.Id
			assertProvider(t, tt.wantUpdatedProvider, getResp.GetIdp())
		})
	}
}

func Test_UpdateZitadelProvider_MissingID(t *testing.T) {
	_ = Instance.AddOrgZitadelProvider(OrgCTX, integration.IDPName())
	// Attempt to update the provider without specifying the ID
	updateResp, err := Client.UpdateZitadelProvider(OrgCTX, &mgmt_pb.UpdateZitadelProviderRequest{})
	require.Error(t, err)
	grpcStatus, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, grpcStatus.Code())
	assert.Equal(t, "invalid UpdateZitadelProviderRequest.Id: value length must be between 1 and 200 runes, inclusive", grpcStatus.Message())
	require.Nil(t, updateResp)
}

func Test_GetProviderByID(t *testing.T) {
	name := integration.IDPName()
	existingProvider := Instance.AddOrgZitadelProvider(OrgCTX, name)

	tests := []struct {
		name     string
		ctx      context.Context
		req      *mgmt_pb.GetProviderByIDRequest
		wantErr  error
		wantResp *mgmt_pb.GetProviderByIDResponse
	}{
		{
			name: "no permissions",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &mgmt_pb.GetProviderByIDRequest{
				Id: "idp-id",
			},
			wantErr: status.Error(codes.NotFound, "membership not found (AUTHZ-cdgFk)"),
		},
		{
			name:    "invalid ID",
			ctx:     OrgCTX,
			wantErr: status.Error(codes.InvalidArgument, "invalid GetProviderByIDRequest.Id: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "non-existing ID",
			ctx:  OrgCTX,
			req: &mgmt_pb.GetProviderByIDRequest{
				Id: "non-existing-id",
			},
			wantErr: status.Error(codes.NotFound, "Identity Provider Configuration doesn't exist (QUERY-SAFrt)"),
		},
		{
			name: "insufficient permissions", // no org.idp.read permissions set
			ctx:  integration.WithSystemUserWithNoPermissionsAuthorization(CTX),
			req: &mgmt_pb.GetProviderByIDRequest{
				Id: existingProvider.GetId(),
			},
			wantErr: status.Error(codes.PermissionDenied, "No matching permissions found (AUTH-5mWD2)"),
		},
		{
			name: "found",
			ctx:  OrgCTX,
			req: &mgmt_pb.GetProviderByIDRequest{
				Id: existingProvider.GetId(),
			},
			wantResp: &mgmt_pb.GetProviderByIDResponse{
				Idp: &idp_pb.Provider{
					Id:    existingProvider.GetId(),
					State: idp_pb.IDPState_IDP_STATE_ACTIVE,
					Name:  name,
					Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG,
					Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
					Config: &idp_pb.ProviderConfig{
						Options: &idp_pb.Options{
							IsCreationAllowed: true,
						},
						Config: &idp_pb.ProviderConfig_Zitadel{
							Zitadel: &idp_pb.ZitadelConfig{
								Issuer:   "zitadel.example.com",
								ClientId: "test-client",
								Scopes:   []string{"email", "profile"},
								InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
									{
										OrganizationId:     "org1",
										OrganizationDomain: "org1.com",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.GetProviderByID(tt.ctx, tt.req)
			if tt.wantErr != nil {
				require.Error(t, err)
				grpcStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, status.Code(tt.wantErr), grpcStatus.Code())
				assert.Equal(t, status.Convert(tt.wantErr).Message(), grpcStatus.Message())
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, existingProvider.GetDetails().GetCreationDate().AsTime(), got.GetIdp().GetDetails().GetCreationDate().AsTime())
			assert.Equal(t, existingProvider.GetDetails().GetResourceOwner(), got.GetIdp().GetDetails().GetResourceOwner())
			assertProvider(t, tt.wantResp.GetIdp(), got.GetIdp())
		})
	}
}

func Test_ListProviders(t *testing.T) {
	org := Instance.CreateOrganization(IAMOwnerCTX, integration.OrganizationName(), integration.Email())
	require.NotEmpty(t, org.GetOrganizationId())
	orgCtx := integration.SetOrgID(IAMOwnerCTX, org.GetOrganizationId())

	provider1Name := integration.IDPName()
	provider1 := Instance.AddOrgZitadelProvider(orgCtx, provider1Name)
	provider2Name := integration.IDPName()
	provider2 := Instance.AddOrgZitadelProvider(orgCtx, provider2Name)

	tests := []struct {
		name     string
		ctx      context.Context
		req      *mgmt_pb.ListProvidersRequest
		wantResp *mgmt_pb.ListProvidersResponse
		wantErr  error
	}{
		{
			name:    "no permissions",
			ctx:     Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req:     &mgmt_pb.ListProvidersRequest{},
			wantErr: status.Error(codes.NotFound, "membership not found (AUTHZ-cdgFk)"),
		},
		{
			name:    "insufficient permissions", // no iam.idp.read permissions set
			ctx:     integration.WithSystemUserWithNoPermissionsAuthorization(CTX),
			req:     &mgmt_pb.ListProvidersRequest{},
			wantErr: status.Error(codes.PermissionDenied, "No matching permissions found (AUTH-5mWD2)"),
		},
		{
			name: "list by id",
			ctx:  orgCtx,
			req: &mgmt_pb.ListProvidersRequest{
				Query: &object_pb.ListQuery{
					Asc: false,
				},
				Queries: []*mgmt_pb.ProviderQuery{
					{
						Query: &mgmt_pb.ProviderQuery_IdpIdQuery{
							IdpIdQuery: &idp_pb.IDPIDQuery{
								Id: provider2.GetId(),
							},
						},
					},
				},
			},
			wantResp: &mgmt_pb.ListProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 1,
				},
				Result: []*idp_pb.Provider{
					{
						Id: provider2.GetId(),
						Details: &object_pb.ObjectDetails{
							CreationDate:  provider2.GetDetails().GetCreationDate(),
							ChangeDate:    provider2.GetDetails().GetChangeDate(),
							ResourceOwner: provider2.GetDetails().GetResourceOwner(),
						},
						State: idp_pb.IDPState_IDP_STATE_ACTIVE,
						Name:  provider2Name,
						Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG,
						Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
						Config: &idp_pb.ProviderConfig{
							Options: &idp_pb.Options{
								IsCreationAllowed: true,
							},
							Config: &idp_pb.ProviderConfig_Zitadel{
								Zitadel: &idp_pb.ZitadelConfig{
									Issuer:   "zitadel.example.com",
									ClientId: "test-client",
									Scopes:   []string{"email", "profile"},
									InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
										{
											OrganizationId:     "org1",
											OrganizationDomain: "org1.com",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list by name",
			ctx:  orgCtx,
			req: &mgmt_pb.ListProvidersRequest{
				Query: &object_pb.ListQuery{
					Asc: false,
				},
				Queries: []*mgmt_pb.ProviderQuery{
					{
						Query: &mgmt_pb.ProviderQuery_IdpNameQuery{
							IdpNameQuery: &idp_pb.IDPNameQuery{
								Name: provider1Name,
							},
						},
					},
				},
			},
			wantResp: &mgmt_pb.ListProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 1,
				},
				Result: []*idp_pb.Provider{
					{
						Id: provider1.GetId(),
						Details: &object_pb.ObjectDetails{
							CreationDate:  provider1.GetDetails().GetCreationDate(),
							ChangeDate:    provider1.GetDetails().GetChangeDate(),
							ResourceOwner: provider1.GetDetails().GetResourceOwner(),
						},
						State: idp_pb.IDPState_IDP_STATE_ACTIVE,
						Name:  provider1Name,
						Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG,
						Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
						Config: &idp_pb.ProviderConfig{
							Options: &idp_pb.Options{
								IsCreationAllowed: true,
							},
							Config: &idp_pb.ProviderConfig_Zitadel{
								Zitadel: &idp_pb.ZitadelConfig{
									Issuer:   "zitadel.example.com",
									ClientId: "test-client",
									Scopes:   []string{"email", "profile"},
									InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
										{
											OrganizationId:     "org1",
											OrganizationDomain: "org1.com",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list all providers",
			ctx:  orgCtx,
			req: &mgmt_pb.ListProvidersRequest{
				Query: &object_pb.ListQuery{
					Asc: true,
				},
			},
			wantResp: &mgmt_pb.ListProvidersResponse{
				Details: &object_pb.ListDetails{
					TotalResult: 2,
				},
				Result: []*idp_pb.Provider{
					{
						Id: provider1.GetId(),
						Details: &object_pb.ObjectDetails{
							CreationDate:  provider1.GetDetails().GetCreationDate(),
							ChangeDate:    provider1.GetDetails().GetChangeDate(),
							ResourceOwner: provider1.GetDetails().GetResourceOwner(),
						},
						State: idp_pb.IDPState_IDP_STATE_ACTIVE,
						Name:  provider1Name,
						Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG,
						Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
						Config: &idp_pb.ProviderConfig{
							Options: &idp_pb.Options{
								IsCreationAllowed: true,
							},
							Config: &idp_pb.ProviderConfig_Zitadel{
								Zitadel: &idp_pb.ZitadelConfig{
									Issuer:   "zitadel.example.com",
									ClientId: "test-client",
									Scopes:   []string{"email", "profile"},
									InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
										{
											OrganizationId:     "org1",
											OrganizationDomain: "org1.com",
										},
									},
								},
							},
						},
					},
					{
						Id: provider2.GetId(),
						Details: &object_pb.ObjectDetails{
							CreationDate:  provider2.GetDetails().GetCreationDate(),
							ChangeDate:    provider2.GetDetails().GetChangeDate(),
							ResourceOwner: provider2.GetDetails().GetResourceOwner(),
						},
						State: idp_pb.IDPState_IDP_STATE_ACTIVE,
						Name:  provider2Name,
						Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG,
						Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
						Config: &idp_pb.ProviderConfig{
							Options: &idp_pb.Options{
								IsCreationAllowed: true,
							},
							Config: &idp_pb.ProviderConfig_Zitadel{
								Zitadel: &idp_pb.ZitadelConfig{
									Issuer:   "zitadel.example.com",
									ClientId: "test-client",
									Scopes:   []string{"email", "profile"},
									InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
										{
											OrganizationId:     "org1",
											OrganizationDomain: "org1.com",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.ListProviders(tt.ctx, tt.req)
			if tt.wantErr != nil {
				require.Error(t, err)
				grpcStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, status.Code(tt.wantErr), grpcStatus.Code())
				assert.Equal(t, status.Convert(tt.wantErr).Message(), grpcStatus.Message())
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantResp.GetDetails().GetTotalResult(), got.GetDetails().GetTotalResult())
			for i, want := range tt.wantResp.GetResult() {
				assert.Equal(t, want.GetDetails().GetCreationDate().AsTime(), got.GetResult()[i].GetDetails().GetCreationDate().AsTime())
				assert.Equal(t, org.GetOrganizationId(), got.GetResult()[i].GetDetails().GetResourceOwner())
				assertProvider(t, want, got.GetResult()[i])
			}
		})
	}
}

func assertProvider(t *testing.T, expected, actual *idp_pb.Provider) {
	assert.Equal(t, expected.GetId(), actual.GetId())
	assert.Equal(t, expected.GetState(), actual.GetState())
	assert.Equal(t, expected.GetName(), actual.GetName())
	assert.Equal(t, expected.GetOwner(), actual.GetOwner())
	assert.Equal(t, expected.GetType(), actual.GetType())
	assert.True(t, proto.Equal(expected.GetConfig(), actual.GetConfig()), "expected: %v, actual: %v", expected.GetConfig(), actual.GetConfig())
}
