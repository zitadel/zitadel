//go:build integration

package admin_test

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
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
)

func Test_AddZitadelProvider(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.AddZitadelProviderRequest
	}
	tests := []struct {
		name         string
		args         args
		wantErr      error
		wantResponse *admin_pb.AddZitadelProviderResponse
	}{
		{
			name: "no permissions, error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &admin_pb.AddZitadelProviderRequest{
					Name:         "Zitadel Support IdP",
					Issuer:       "zitadel.example.com",
					ClientId:     "test-client",
					ClientSecret: "test-secret",
					Scopes:       []string{"email", "profile"},
				},
			},
			wantErr: status.Error(codes.NotFound, "membership not found (AUTHZ-cdgFk)"),
		},
		{
			name: "insufficient permissions, error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeLogin),
				req: &admin_pb.AddZitadelProviderRequest{
					Name:         "Zitadel Support IdP",
					Issuer:       "zitadel.example.com",
					ClientId:     "test-client",
					ClientSecret: "test-secret",
					Scopes:       []string{"email", "profile"},
				},
			},
			wantErr: status.Error(codes.PermissionDenied, "No matching permissions found (AUTH-5mWD2)"),
		},
		{
			name: "missing required field: name",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.AddZitadelProviderRequest{},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid AddZitadelProviderRequest.Name: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: issuer",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.AddZitadelProviderRequest{
					Name: "Zitadel Support IdP",
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid AddZitadelProviderRequest.Issuer: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: client_id",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.AddZitadelProviderRequest{
					Name:   "Zitadel Support IdP",
					Issuer: "zitadel.example.com",
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid AddZitadelProviderRequest.ClientId: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: client_secret",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.AddZitadelProviderRequest{
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
				ctx: AdminCTX,
				req: &admin_pb.AddZitadelProviderRequest{
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
				ctx: AdminCTX,
				req: &admin_pb.AddZitadelProviderRequest{
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
				ctx: AdminCTX,
				req: &admin_pb.AddZitadelProviderRequest{
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
			wantResponse: &admin_pb.AddZitadelProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.Instance.Id,
				},
			},
		},
		{
			name: "valid request with instance roles info",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.AddZitadelProviderRequest{
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
			wantResponse: &admin_pb.AddZitadelProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.Instance.Id,
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
			t.Cleanup(func() {
				_, err := Client.DeleteProvider(AdminCTX, &admin_pb.DeleteProviderRequest{Id: got.GetId()})
				require.NoError(t, err)
			})
		})
	}
}

func Test_UpdateZitadelProvider(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.UpdateZitadelProviderRequest
	}
	tests := []struct {
		name                string
		args                args
		wantErr             error
		wantResponse        *admin_pb.UpdateZitadelProviderResponse
		wantUpdatedProvider *idp_pb.Provider
	}{
		{
			name: "no permissions, error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &admin_pb.UpdateZitadelProviderRequest{
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
				req: &admin_pb.UpdateZitadelProviderRequest{
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
				ctx: AdminCTX,
				req: &admin_pb.UpdateZitadelProviderRequest{},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid UpdateZitadelProviderRequest.Name: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: issuer",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.UpdateZitadelProviderRequest{
					Name: "Zitadel Support IdP updated",
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid UpdateZitadelProviderRequest.Issuer: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: client_id",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.UpdateZitadelProviderRequest{
					Name:   "Zitadel Support IdP updated",
					Issuer: "acme.example.com",
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid UpdateZitadelProviderRequest.ClientId: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "update, ok",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.UpdateZitadelProviderRequest{
					Name:     "Zitadel Support IdP updated",
					Issuer:   "acme.example.com",
					ClientId: "test-client",
					Scopes:   []string{"email", "profile", "openid", "offline_access"},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
						IsAutoCreation:    false,
					},
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
			wantResponse: &admin_pb.UpdateZitadelProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.Instance.Id,
				},
			},
			wantUpdatedProvider: &idp_pb.Provider{
				State: idp_pb.IDPState_IDP_STATE_ACTIVE,
				Name:  "Zitadel Support IdP updated",
				Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM,
				Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
				Config: &idp_pb.ProviderConfig{
					Options: &idp_pb.Options{
						IsCreationAllowed: true,
					},
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
				ctx: AdminCTX,
				req: &admin_pb.UpdateZitadelProviderRequest{
					Name:     "Zitadel Support IdP updated 1",
					Issuer:   "acme.example.com",
					ClientId: "test-client",
					Scopes:   []string{},
					ProviderOptions: &idp_pb.Options{
						IsCreationAllowed: true,
						IsAutoCreation:    false,
					},
					InstanceRolesInfo: []*idp_pb.InstanceRolesInfo{
						{
							OrganizationId:     "org3",
							OrganizationDomain: "org3.com",
						},
					},
				},
			},
			wantResponse: &admin_pb.UpdateZitadelProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.Instance.Id,
				},
			},
			wantUpdatedProvider: &idp_pb.Provider{
				State: idp_pb.IDPState_IDP_STATE_ACTIVE,
				Name:  "Zitadel Support IdP updated 1",
				Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM,
				Type:  idp_pb.ProviderType_PROVIDER_TYPE_ZITADEL,
				Config: &idp_pb.ProviderConfig{
					Options: &idp_pb.Options{
						IsCreationAllowed: true,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a new provider per subtest
			existingProvider := Instance.AddZitadelProvider(AdminCTX, integration.IDPName())
			t.Cleanup(func() {
				_, err := Client.DeleteProvider(AdminCTX, &admin_pb.DeleteProviderRequest{Id: existingProvider.GetId()})
				require.NoError(t, err)
			})
			// build request using this provider ID
			tt.args.req.Id = existingProvider.GetId()

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
			getResp, err := Client.GetProviderByID(AdminCTX, &admin_pb.GetProviderByIDRequest{Id: existingProvider.GetId()})
			require.NoError(t, err)
			assert.Equal(t, updateResp.GetDetails().GetChangeDate().AsTime(), getResp.GetIdp().GetDetails().GetChangeDate().AsTime())
			tt.wantUpdatedProvider.Id = existingProvider.GetId()
			assertProvider(t, tt.wantUpdatedProvider, getResp.GetIdp())
		})
	}
}

func Test_UpdateZitadelProvider_MissingID(t *testing.T) {
	existingProvider := Instance.AddZitadelProvider(AdminCTX, integration.IDPName())
	t.Cleanup(func() {
		_, err := Client.DeleteProvider(AdminCTX, &admin_pb.DeleteProviderRequest{Id: existingProvider.GetId()})
		require.NoError(t, err)
	})
	// Attempt to update the provider without specifying the ID
	updateResp, err := Client.UpdateZitadelProvider(AdminCTX, &admin_pb.UpdateZitadelProviderRequest{})
	require.Error(t, err)
	grpcStatus, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, grpcStatus.Code())
	assert.Equal(t, "invalid UpdateZitadelProviderRequest.Id: value length must be between 1 and 200 runes, inclusive", grpcStatus.Message())
	require.Nil(t, updateResp)
}

func Test_GetProviderByID(t *testing.T) {
	providerName := integration.IDPName()
	existingProvider := Instance.AddZitadelProvider(AdminCTX, providerName)
	t.Cleanup(func() {
		_, err := Client.DeleteProvider(AdminCTX, &admin_pb.DeleteProviderRequest{Id: existingProvider.GetId()})
		require.NoError(t, err)
	})

	tests := []struct {
		name     string
		ctx      context.Context
		req      *admin_pb.GetProviderByIDRequest
		wantErr  error
		wantResp *admin_pb.GetProviderByIDResponse
	}{
		{
			name: "no permissions",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &admin_pb.GetProviderByIDRequest{
				Id: "idp-id",
			},
			wantErr: status.Error(codes.NotFound, "membership not found (AUTHZ-cdgFk)"),
		},
		{
			name:    "invalid ID",
			ctx:     AdminCTX,
			wantErr: status.Error(codes.InvalidArgument, "invalid GetProviderByIDRequest.Id: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "non-existing ID",
			ctx:  AdminCTX,
			req: &admin_pb.GetProviderByIDRequest{
				Id: "non-existing-id",
			},
			wantErr: status.Error(codes.NotFound, "Identity Provider Configuration doesn't exist (QUERY-SAFrt)"),
		},
		{
			name: "insufficient permissions", // no iam.idp.read permissions set
			ctx:  integration.WithSystemUserWithNoPermissionsAuthorization(CTX),
			req: &admin_pb.GetProviderByIDRequest{
				Id: existingProvider.GetId(),
			},
			wantErr: status.Error(codes.PermissionDenied, "No matching permissions found (AUTH-5mWD2)"),
		},
		{
			name: "found",
			ctx:  AdminCTX,
			req: &admin_pb.GetProviderByIDRequest{
				Id: existingProvider.GetId(),
			},
			wantResp: &admin_pb.GetProviderByIDResponse{
				Idp: &idp_pb.Provider{
					Id:    existingProvider.GetId(),
					State: idp_pb.IDPState_IDP_STATE_ACTIVE,
					Name:  providerName,
					Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM,
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
	provider1Name := integration.IDPName()
	provider1 := Instance.AddZitadelProvider(AdminCTX, provider1Name)

	provider2Name := integration.IDPName()
	provider2 := Instance.AddZitadelProvider(AdminCTX, provider2Name)
	t.Cleanup(func() {
		_, err := Client.DeleteProvider(AdminCTX, &admin_pb.DeleteProviderRequest{Id: provider1.GetId()})
		require.NoError(t, err)
		_, err = Client.DeleteProvider(AdminCTX, &admin_pb.DeleteProviderRequest{Id: provider2.GetId()})
		require.NoError(t, err)
	})

	tests := []struct {
		name     string
		ctx      context.Context
		req      *admin_pb.ListProvidersRequest
		wantResp *admin_pb.ListProvidersResponse
		wantErr  error
	}{
		{
			name:    "no permissions",
			ctx:     Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req:     &admin_pb.ListProvidersRequest{},
			wantErr: status.Error(codes.NotFound, "membership not found (AUTHZ-cdgFk)"),
		},
		{
			name:    "insufficient permissions", // no iam.idp.read permissions set
			ctx:     integration.WithSystemUserWithNoPermissionsAuthorization(CTX),
			req:     &admin_pb.ListProvidersRequest{},
			wantErr: status.Error(codes.PermissionDenied, "No matching permissions found (AUTH-5mWD2)"),
		},
		{
			name: "list all providers",
			ctx:  AdminCTX,
			req: &admin_pb.ListProvidersRequest{
				Query: &object_pb.ListQuery{
					Asc: true,
				},
			},
			wantResp: &admin_pb.ListProvidersResponse{
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
						Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM,
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
						Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM,
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
			name: "list by id",
			ctx:  AdminCTX,
			req: &admin_pb.ListProvidersRequest{
				Query: &object_pb.ListQuery{
					Asc: false,
				},
				Queries: []*admin_pb.ProviderQuery{
					{
						Query: &admin_pb.ProviderQuery_IdpIdQuery{
							IdpIdQuery: &idp_pb.IDPIDQuery{
								Id: provider2.GetId(),
							},
						},
					},
				},
			},
			wantResp: &admin_pb.ListProvidersResponse{
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
						Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM,
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
			ctx:  AdminCTX,
			req: &admin_pb.ListProvidersRequest{
				Query: &object_pb.ListQuery{
					Asc: false,
				},
				Queries: []*admin_pb.ProviderQuery{
					{
						Query: &admin_pb.ProviderQuery_IdpNameQuery{
							IdpNameQuery: &idp_pb.IDPNameQuery{
								Name: provider1Name,
							},
						},
					},
				},
			},
			wantResp: &admin_pb.ListProvidersResponse{
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
						Owner: idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM,
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
				assert.Equal(t, Instance.ID(), got.GetResult()[i].GetDetails().GetResourceOwner())
				assertProvider(t, want, got.GetResult()[i])
			}
		})
	}
}

func Test_DeleteZitadelProvider(t *testing.T) {
	existingProvider := Instance.AddZitadelProvider(AdminCTX, integration.IDPName())
	t.Cleanup(func() {
		_, err := Client.DeleteProvider(AdminCTX, &admin_pb.DeleteProviderRequest{Id: existingProvider.GetId()})
		if err != nil && status.Code(err) != codes.NotFound {
			require.NoError(t, err)
		}
	})

	tests := []struct {
		name         string
		ctx          context.Context
		req          *admin_pb.DeleteProviderRequest
		wantResponse *admin_pb.DeleteProviderResponse
		wantErr      error
	}{
		{
			name:    "no permissions, error",
			ctx:     Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req:     &admin_pb.DeleteProviderRequest{Id: "idp-id"},
			wantErr: status.Error(codes.NotFound, "membership not found (AUTHZ-cdgFk)"),
		},
		{
			name:    "insufficient permissions, error", // no iam.idp.write permission
			ctx:     integration.WithSystemUserWithNoPermissionsAuthorization(CTX),
			req:     &admin_pb.DeleteProviderRequest{Id: "idp-id"},
			wantErr: status.Error(codes.PermissionDenied, "No matching permissions found (AUTH-5mWD2)"),
		},
		{
			name:    "not found, error",
			ctx:     AdminCTX,
			req:     &admin_pb.DeleteProviderRequest{Id: "idp-id"},
			wantErr: status.Error(codes.NotFound, "Identity Provider Configuration doesn't exist (INST-Se3tg)"),
		},
		{
			name: "delete, ok",
			ctx:  AdminCTX,
			req:  &admin_pb.DeleteProviderRequest{Id: existingProvider.GetId()},
			wantResponse: &admin_pb.DeleteProviderResponse{
				Details: &object_pb.ObjectDetails{
					ResourceOwner: Instance.Instance.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.DeleteProvider(tt.ctx, tt.req)
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
			assert.Equal(t, tt.wantResponse.GetDetails().GetResourceOwner(), got.GetDetails().GetResourceOwner())
			assert.WithinRange(t, got.GetDetails().GetChangeDate().AsTime(), existingProvider.GetDetails().GetCreationDate().AsTime(), after)
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
