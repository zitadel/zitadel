//go:build integration

package settings_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	settings "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
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
			ctx:     Tester.WithAuthorization(CTX, integration.OrgOwner),
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
			resp, err := Client.GetSecuritySettings(tt.ctx, &settings.GetSecuritySettingsRequest{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			got, want := resp.GetSettings(), tt.want.GetSettings()
			assert.Equal(t, want.GetEmbeddedIframe().GetEnabled(), got.GetEmbeddedIframe().GetEnabled(), "enable iframe embedding")
			assert.Equal(t, want.GetEmbeddedIframe().GetAllowedOrigins(), got.GetEmbeddedIframe().GetAllowedOrigins(), "allowed origins")
			assert.Equal(t, want.GetEnableImpersonation(), got.GetEnableImpersonation(), "enable impersonation")
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
				ctx: Tester.WithAuthorization(CTX, integration.OrgOwner),
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
					ResourceOwner: Tester.Instance.InstanceID(),
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
					ResourceOwner: Tester.Instance.InstanceID(),
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
					ResourceOwner: Tester.Instance.InstanceID(),
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
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetSecuritySettings(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
