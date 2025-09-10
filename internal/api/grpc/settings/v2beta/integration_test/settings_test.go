//go:build integration

package settings_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
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
					ResourceOwner: Instance.Instance.Id,
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
					ResourceOwner: Instance.Instance.Id,
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
					ResourceOwner: Instance.Instance.Id,
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
					ResourceOwner: Instance.Instance.Id,
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

func TestServer_SetOrganizationSettings(t *testing.T) {
	instance := integration.NewInstance(CTX)
	iamOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		req *settings.SetOrganizationSettingsRequest
	}
	type want struct {
		set     bool
		setDate bool
	}
	tests := []struct {
		name    string
		prepare func(req *settings.SetOrganizationSettingsRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &settings.SetOrganizationSettingsRequest{
					OrganizationId:              Instance.DefaultOrg.GetId(),
					OrganizationScopedUsernames: gu.Ptr(true),
				},
			},
			wantErr: true,
		},
		{
			name: "org not provided",
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.SetOrganizationSettingsRequest{
					OrganizationId:              "",
					OrganizationScopedUsernames: gu.Ptr(true),
				},
			},
			wantErr: true,
		},
		{
			name: "org not existing",
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.SetOrganizationSettingsRequest{
					OrganizationId:              "notexisting",
					OrganizationScopedUsernames: gu.Ptr(true),
				},
			},
			wantErr: true,
		},
		{
			name: "success no changes",
			prepare: func(req *settings.SetOrganizationSettingsRequest) {
				orgResp := instance.CreateOrganization(iamOwnerCTX, integration.OrganizationName(), integration.Email())
				req.OrganizationId = orgResp.GetOrganizationId()
			},
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.SetOrganizationSettingsRequest{},
			},
			want: want{
				set:     false,
				setDate: true,
			},
		},
		{
			name: "success user uniqueness",
			prepare: func(req *settings.SetOrganizationSettingsRequest) {
				orgResp := instance.CreateOrganization(iamOwnerCTX, integration.OrganizationName(), integration.Email())
				req.OrganizationId = orgResp.GetOrganizationId()
			},
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.SetOrganizationSettingsRequest{
					OrganizationScopedUsernames: gu.Ptr(true),
				},
			},
			want: want{
				set:     true,
				setDate: true,
			},
		},
		{
			name: "success no change",
			prepare: func(req *settings.SetOrganizationSettingsRequest) {
				orgResp := instance.CreateOrganization(iamOwnerCTX, integration.OrganizationName(), integration.Email())
				req.OrganizationId = orgResp.GetOrganizationId()
			},
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.SetOrganizationSettingsRequest{
					OrganizationScopedUsernames: gu.Ptr(false),
				},
			},
			want: want{
				set:     false,
				setDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			if tt.prepare != nil {
				tt.prepare(tt.args.req)
			}

			got, err := instance.Client.SettingsV2beta.SetOrganizationSettings(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			setDate := time.Time{}
			if tt.want.set {
				setDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertOrganizationSettingsResponse(t, creationDate, setDate, tt.want.setDate, got)
		})
	}
}

func assertOrganizationSettingsResponse(t *testing.T, creationDate, setDate time.Time, expectedSetDate bool, actualResp *settings.SetOrganizationSettingsResponse) {
	if expectedSetDate {
		if !setDate.IsZero() {
			assert.WithinRange(t, actualResp.GetSetDate().AsTime(), creationDate, setDate)
		} else {
			assert.WithinRange(t, actualResp.GetSetDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.SetDate)
	}
}

func TestServer_DeleteOrganizationSettings(t *testing.T) {
	instance := integration.NewInstance(CTX)
	iamOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		req *settings.DeleteOrganizationSettingsRequest
	}
	type want struct {
		deletion     bool
		deletionDate bool
	}
	tests := []struct {
		name    string
		prepare func(t *testing.T, req *settings.DeleteOrganizationSettingsRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "permission error",
			prepare: func(t *testing.T, req *settings.DeleteOrganizationSettingsRequest) {
				orgResp := instance.CreateOrganization(iamOwnerCTX, integration.OrganizationName(), integration.Email())
				req.OrganizationId = orgResp.GetOrganizationId()
				instance.SetOrganizationSettings(iamOwnerCTX, t, orgResp.GetOrganizationId(), true)
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &settings.DeleteOrganizationSettingsRequest{
					OrganizationId: Instance.DefaultOrg.GetId(),
				},
			},
			wantErr: true,
		},
		{
			name: "org not provided",
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.DeleteOrganizationSettingsRequest{
					OrganizationId: "",
				},
			},
			wantErr: true,
		},
		{
			name: "org not existing",
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.DeleteOrganizationSettingsRequest{
					OrganizationId: "notexisting",
				},
			},
			want: want{
				deletion:     false,
				deletionDate: false,
			},
		},
		{
			name: "success user uniqueness",
			prepare: func(t *testing.T, req *settings.DeleteOrganizationSettingsRequest) {
				orgResp := instance.CreateOrganization(iamOwnerCTX, integration.OrganizationName(), integration.Email())
				req.OrganizationId = orgResp.GetOrganizationId()
				instance.SetOrganizationSettings(iamOwnerCTX, t, orgResp.GetOrganizationId(), true)
			},
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.DeleteOrganizationSettingsRequest{},
			},
			want: want{
				deletion:     true,
				deletionDate: true,
			},
		},
		{
			name: "success no existing",
			prepare: func(t *testing.T, req *settings.DeleteOrganizationSettingsRequest) {
				orgResp := instance.CreateOrganization(iamOwnerCTX, integration.OrganizationName(), integration.Email())
				req.OrganizationId = orgResp.GetOrganizationId()
			},
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.DeleteOrganizationSettingsRequest{},
			},
			want: want{
				deletion:     false,
				deletionDate: true,
			},
		},
		{
			name: "success already deleted",
			prepare: func(t *testing.T, req *settings.DeleteOrganizationSettingsRequest) {
				orgResp := instance.CreateOrganization(iamOwnerCTX, integration.OrganizationName(), integration.Email())
				req.OrganizationId = orgResp.GetOrganizationId()
				instance.SetOrganizationSettings(iamOwnerCTX, t, orgResp.GetOrganizationId(), true)
				instance.DeleteOrganizationSettings(iamOwnerCTX, t, orgResp.GetOrganizationId())
			},
			args: args{
				ctx: iamOwnerCTX,
				req: &settings.DeleteOrganizationSettingsRequest{},
			},
			want: want{
				deletion:     false,
				deletionDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			if tt.prepare != nil {
				tt.prepare(t, tt.args.req)
			}

			got, err := instance.Client.SettingsV2beta.DeleteOrganizationSettings(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			deletionDate := time.Time{}
			if tt.want.deletion {
				deletionDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertDeleteOrganizationSettingsResponse(t, creationDate, deletionDate, tt.want.deletionDate, got)
		})
	}
}

func assertDeleteOrganizationSettingsResponse(t *testing.T, creationDate, deletionDate time.Time, expectedDeletionDate bool, actualResp *settings.DeleteOrganizationSettingsResponse) {
	if expectedDeletionDate {
		if !deletionDate.IsZero() {
			assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), creationDate, deletionDate)
		} else {
			assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.DeletionDate)
	}
}
