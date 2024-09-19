//go:build integration

package admin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/v2/internal/integration"
	admin_pb "github.com/zitadel/zitadel/v2/pkg/grpc/admin"
	"github.com/zitadel/zitadel/v2/pkg/grpc/object"
	"github.com/zitadel/zitadel/v2/pkg/grpc/settings"
)

func TestServer_GetSecurityPolicy(t *testing.T) {
	t.Parallel()

	instance := integration.NewInstance(CTX)
	adminCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	_, err := instance.Client.Admin.SetSecurityPolicy(adminCtx, &admin_pb.SetSecurityPolicyRequest{
		EnableIframeEmbedding: true,
		AllowedOrigins:        []string{"foo.com", "bar.com"},
		EnableImpersonation:   true,
	})
	require.NoError(t, err)
	tests := []struct {
		name    string
		ctx     context.Context
		want    *admin_pb.GetSecurityPolicyResponse
		wantErr bool
	}{
		{
			name:    "permission error",
			ctx:     instance.WithAuthorization(CTX, integration.UserTypeOrgOwner),
			wantErr: true,
		},
		{
			name: "success",
			ctx:  adminCtx,
			want: &admin_pb.GetSecurityPolicyResponse{
				Policy: &settings.SecurityPolicy{
					EnableIframeEmbedding: true,
					AllowedOrigins:        []string{"foo.com", "bar.com"},
					EnableImpersonation:   true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := instance.Client.Admin.GetSecurityPolicy(tt.ctx, &admin_pb.GetSecurityPolicyRequest{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			got, want := resp.GetPolicy(), tt.want.GetPolicy()
			assert.Equal(t, want.GetEnableIframeEmbedding(), got.GetEnableIframeEmbedding(), "enable iframe embedding")
			assert.Equal(t, want.GetAllowedOrigins(), got.GetAllowedOrigins(), "allowed origins")
			assert.Equal(t, want.GetEnableImpersonation(), got.GetEnableImpersonation(), "enable impersonation")
		})
	}
}

func TestServer_SetSecurityPolicy(t *testing.T) {
	t.Parallel()

	instance := integration.NewInstance(CTX)
	adminCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		req *admin_pb.SetSecurityPolicyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *admin_pb.SetSecurityPolicyResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: instance.WithAuthorization(CTX, integration.UserTypeOrgOwner),
				req: &admin_pb.SetSecurityPolicyRequest{
					EnableIframeEmbedding: true,
					AllowedOrigins:        []string{"foo.com", "bar.com"},
					EnableImpersonation:   true,
				},
			},
			wantErr: true,
		},
		{
			name: "success allowed origins",
			args: args{
				ctx: adminCtx,
				req: &admin_pb.SetSecurityPolicyRequest{
					AllowedOrigins: []string{"foo.com", "bar.com"},
				},
			},
			want: &admin_pb.SetSecurityPolicyResponse{
				Details: &object.ObjectDetails{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: instance.ID(),
				},
			},
		},
		{
			name: "success iframe embedding",
			args: args{
				ctx: adminCtx,
				req: &admin_pb.SetSecurityPolicyRequest{
					EnableIframeEmbedding: true,
				},
			},
			want: &admin_pb.SetSecurityPolicyResponse{
				Details: &object.ObjectDetails{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: instance.ID(),
				},
			},
		},
		{
			name: "success impersonation",
			args: args{
				ctx: adminCtx,
				req: &admin_pb.SetSecurityPolicyRequest{
					EnableImpersonation: true,
				},
			},
			want: &admin_pb.SetSecurityPolicyResponse{
				Details: &object.ObjectDetails{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: instance.ID(),
				},
			},
		},
		{
			name: "success all",
			args: args{
				ctx: adminCtx,
				req: &admin_pb.SetSecurityPolicyRequest{
					EnableIframeEmbedding: true,
					AllowedOrigins:        []string{"foo.com", "bar.com"},
					EnableImpersonation:   true,
				},
			},
			want: &admin_pb.SetSecurityPolicyResponse{
				Details: &object.ObjectDetails{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: instance.ID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Client.Admin.SetSecurityPolicy(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
