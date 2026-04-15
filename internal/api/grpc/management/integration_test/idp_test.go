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
