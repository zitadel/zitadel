//go:build integration

package settings_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/v2/internal/integration"
	object_pb "github.com/zitadel/zitadel/v2/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/v2/pkg/grpc/settings/v2"
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
			assert.EventuallyWithT(t, func(ct *assert.CollectT) {
				resp, err := Client.GetSecuritySettings(tt.ctx, &settings.GetSecuritySettingsRequest{})
				if tt.wantErr {
					assert.Error(ct, err)
					return
				}
				require.NoError(ct, err)
				got, want := resp.GetSettings(), tt.want.GetSettings()
				assert.Equal(ct, want.GetEmbeddedIframe().GetEnabled(), got.GetEmbeddedIframe().GetEnabled(), "enable iframe embedding")
				assert.Equal(ct, want.GetEmbeddedIframe().GetAllowedOrigins(), got.GetEmbeddedIframe().GetAllowedOrigins(), "allowed origins")
				assert.Equal(ct, want.GetEnableImpersonation(), got.GetEnableImpersonation(), "enable impersonation")
			}, time.Minute, time.Second/10)
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
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
