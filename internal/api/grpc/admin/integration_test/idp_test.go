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

	"github.com/zitadel/zitadel/internal/integration"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
)

func Test_AddZitadelProvider(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
	existingProvider, err := Instance.Client.Admin.AddZitadelProvider(AdminCTX, &admin_pb.AddZitadelProviderRequest{
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
	})
	require.NoError(t, err)
	require.NotEmpty(t, existingProvider.GetId())

	type args struct {
		ctx context.Context
		req *admin_pb.UpdateZitadelProviderRequest
	}
	tests := []struct {
		name         string
		args         args
		wantErr      error
		wantResponse *admin_pb.UpdateZitadelProviderResponse
	}{
		{
			name: "no permissions, error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &admin_pb.UpdateZitadelProviderRequest{
					Id:           existingProvider.GetId(),
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
					Id:     existingProvider.GetId(),
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
			name: "missing required field: id",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.UpdateZitadelProviderRequest{},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid UpdateZitadelProviderRequest.Id: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: name",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.UpdateZitadelProviderRequest{
					Id: existingProvider.GetId(),
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid UpdateZitadelProviderRequest.Name: value length must be between 1 and 200 runes, inclusive"),
		},
		{
			name: "missing required field: issuer",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.UpdateZitadelProviderRequest{
					Id:   existingProvider.GetId(),
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
					Id:     existingProvider.GetId(),
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
					Id:       existingProvider.GetId(),
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			before := time.Now()
			got, err := Client.UpdateZitadelProvider(tt.args.ctx, tt.args.req)
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
			assert.WithinRange(t, got.GetDetails().GetChangeDate().AsTime(), before, after)
			assert.Equal(t, tt.wantResponse.GetDetails().GetResourceOwner(), got.GetDetails().GetResourceOwner())
		})
	}
}
