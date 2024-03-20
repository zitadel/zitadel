//go:build integration

package admin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/settings"
)

func TestServer_GetSecurityPolicy(t *testing.T) {
	_, err := Client.SetSecurityPolicy(AdminCTX, &admin_pb.SetSecurityPolicyRequest{
		EnableIframeEmbedding: true,
		AllowedOrigins:        []string{"foo.com", "bar.com"},
		EnableImpersonation:   true,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := Client.SetSecurityPolicy(AdminCTX, &admin_pb.SetSecurityPolicyRequest{
			EnableIframeEmbedding: false,
			AllowedOrigins:        []string{},
			EnableImpersonation:   false,
		})
		require.NoError(t, err)
	})

	tests := []struct {
		name    string
		ctx     context.Context
		want    *admin_pb.GetSecurityPolicyResponse
		wantErr bool
	}{
		{
			name:    "permission error",
			ctx:     Tester.WithAuthorization(CTX, integration.OrgOwner),
			wantErr: true,
		},
		{
			name: "success",
			ctx:  AdminCTX,
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
			resp, err := Client.GetSecurityPolicy(tt.ctx, &admin_pb.GetSecurityPolicyRequest{})
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
				ctx: Tester.WithAuthorization(CTX, integration.OrgOwner),
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
				ctx: AdminCTX,
				req: &admin_pb.SetSecurityPolicyRequest{
					AllowedOrigins: []string{"foo.com", "bar.com"},
				},
			},
			want: &admin_pb.SetSecurityPolicyResponse{
				Details: &object.ObjectDetails{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "success iframe embedding",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.SetSecurityPolicyRequest{
					EnableIframeEmbedding: true,
				},
			},
			want: &admin_pb.SetSecurityPolicyResponse{
				Details: &object.ObjectDetails{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "success impersonation",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.SetSecurityPolicyRequest{
					EnableImpersonation: true,
				},
			},
			want: &admin_pb.SetSecurityPolicyResponse{
				Details: &object.ObjectDetails{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "success all",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.SetSecurityPolicyRequest{
					EnableIframeEmbedding: true,
					AllowedOrigins:        []string{"foo.com", "bar.com"},
					EnableImpersonation:   true,
				},
			},
			want: &admin_pb.SetSecurityPolicyResponse{
				Details: &object.ObjectDetails{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetSecurityPolicy(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
