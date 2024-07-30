//go:build integration

package feature_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

var (
	SystemCTX context.Context
	IamCTX    context.Context
	OrgCTX    context.Context
	Tester    *integration.Tester
	Client    feature.FeatureServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()
		Tester = integration.NewTester(ctx)
		SystemCTX = Tester.WithAuthorization(ctx, integration.SystemUser)
		IamCTX = Tester.WithAuthorization(ctx, integration.IAMOwner)
		OrgCTX = Tester.WithAuthorization(ctx, integration.OrgOwner)

		defer Tester.Done()
		Client = Tester.Client.FeatureV2

		return m.Run()
	}())
}

func TestServer_SetSystemFeatures(t *testing.T) {
	type args struct {
		ctx context.Context
		req *feature.SetSystemFeaturesRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *feature.SetSystemFeaturesResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: IamCTX,
				req: &feature.SetSystemFeaturesRequest{
					OidcTriggerIntrospectionProjections: gu.Ptr(true),
				},
			},
			wantErr: true,
		},
		{
			name: "no changes error",
			args: args{
				ctx: SystemCTX,
				req: &feature.SetSystemFeaturesRequest{},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: SystemCTX,
				req: &feature.SetSystemFeaturesRequest{
					OidcTriggerIntrospectionProjections: gu.Ptr(true),
				},
			},
			want: &feature.SetSystemFeaturesResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: "SYSTEM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				// make sure we have a clean state after each test
				_, err := Client.ResetSystemFeatures(SystemCTX, &feature.ResetSystemFeaturesRequest{})
				require.NoError(t, err)
			})
			got, err := Client.SetSystemFeatures(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_ResetSystemFeatures(t *testing.T) {
	_, err := Client.SetSystemFeatures(SystemCTX, &feature.SetSystemFeaturesRequest{
		LoginDefaultOrg: gu.Ptr(true),
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		ctx     context.Context
		want    *feature.ResetSystemFeaturesResponse
		wantErr bool
	}{
		{
			name:    "permission error",
			ctx:     IamCTX,
			wantErr: true,
		},
		{
			name: "success",
			ctx:  SystemCTX,
			want: &feature.ResetSystemFeaturesResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: "SYSTEM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.ResetSystemFeatures(tt.ctx, &feature.ResetSystemFeaturesRequest{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_GetSystemFeatures(t *testing.T) {
	type args struct {
		ctx context.Context
		req *feature.GetSystemFeaturesRequest
	}
	tests := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		want    *feature.GetSystemFeaturesResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: IamCTX,
				req: &feature.GetSystemFeaturesRequest{},
			},
			wantErr: true,
		},
		{
			name: "nothing set",
			args: args{
				ctx: SystemCTX,
				req: &feature.GetSystemFeaturesRequest{},
			},
			want: &feature.GetSystemFeaturesResponse{},
		},
		{
			name: "some features",
			prepare: func(t *testing.T) {
				_, err := Client.SetSystemFeatures(SystemCTX, &feature.SetSystemFeaturesRequest{
					LoginDefaultOrg:                     gu.Ptr(true),
					OidcTriggerIntrospectionProjections: gu.Ptr(false),
				})
				require.NoError(t, err)
			},
			args: args{
				ctx: SystemCTX,
				req: &feature.GetSystemFeaturesRequest{},
			},
			want: &feature.GetSystemFeaturesResponse{
				LoginDefaultOrg: &feature.FeatureFlag{
					Enabled: true,
					Source:  feature.Source_SOURCE_SYSTEM,
				},
				OidcTriggerIntrospectionProjections: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_SYSTEM,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				// make sure we have a clean state after each test
				_, err := Client.ResetSystemFeatures(SystemCTX, &feature.ResetSystemFeaturesRequest{})
				require.NoError(t, err)
			})
			if tt.prepare != nil {
				tt.prepare(t)
			}
			got, err := Client.GetSystemFeatures(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assertFeatureFlag(t, tt.want.LoginDefaultOrg, got.LoginDefaultOrg)
			assertFeatureFlag(t, tt.want.OidcTriggerIntrospectionProjections, got.OidcTriggerIntrospectionProjections)
			assertFeatureFlag(t, tt.want.OidcLegacyIntrospection, got.OidcLegacyIntrospection)
			assertFeatureFlag(t, tt.want.UserSchema, got.UserSchema)
			assertFeatureFlag(t, tt.want.Actions, got.Actions)
		})
	}
}

func TestServer_SetInstanceFeatures(t *testing.T) {
	type args struct {
		ctx context.Context
		req *feature.SetInstanceFeaturesRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *feature.SetInstanceFeaturesResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: OrgCTX,
				req: &feature.SetInstanceFeaturesRequest{
					OidcTriggerIntrospectionProjections: gu.Ptr(true),
				},
			},
			wantErr: true,
		},
		{
			name: "no changes error",
			args: args{
				ctx: IamCTX,
				req: &feature.SetInstanceFeaturesRequest{},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: IamCTX,
				req: &feature.SetInstanceFeaturesRequest{
					OidcTriggerIntrospectionProjections: gu.Ptr(true),
				},
			},
			want: &feature.SetInstanceFeaturesResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				// make sure we have a clean state after each test
				_, err := Client.ResetInstanceFeatures(IamCTX, &feature.ResetInstanceFeaturesRequest{})
				require.NoError(t, err)
			})
			got, err := Client.SetInstanceFeatures(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_ResetInstanceFeatures(t *testing.T) {
	_, err := Client.SetInstanceFeatures(IamCTX, &feature.SetInstanceFeaturesRequest{
		LoginDefaultOrg: gu.Ptr(true),
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		ctx     context.Context
		want    *feature.ResetInstanceFeaturesResponse
		wantErr bool
	}{
		{
			name:    "permission error",
			ctx:     OrgCTX,
			wantErr: true,
		},
		{
			name: "success",
			ctx:  IamCTX,
			want: &feature.ResetInstanceFeaturesResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.ResetInstanceFeatures(tt.ctx, &feature.ResetInstanceFeaturesRequest{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_GetInstanceFeatures(t *testing.T) {
	_, err := Client.SetSystemFeatures(SystemCTX, &feature.SetSystemFeaturesRequest{
		OidcLegacyIntrospection: gu.Ptr(true),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := Client.ResetSystemFeatures(SystemCTX, &feature.ResetSystemFeaturesRequest{})
		require.NoError(t, err)
	})

	type args struct {
		ctx context.Context
		req *feature.GetInstanceFeaturesRequest
	}
	tests := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		want    *feature.GetInstanceFeaturesResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: OrgCTX,
				req: &feature.GetInstanceFeaturesRequest{},
			},
			wantErr: true,
		},
		{
			name: "defaults, no inheritance",
			args: args{
				ctx: IamCTX,
				req: &feature.GetInstanceFeaturesRequest{},
			},
			want: &feature.GetInstanceFeaturesResponse{},
		},
		{
			name: "defaults, inheritance",
			args: args{
				ctx: IamCTX,
				req: &feature.GetInstanceFeaturesRequest{
					Inheritance: true,
				},
			},
			want: &feature.GetInstanceFeaturesResponse{
				LoginDefaultOrg: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_UNSPECIFIED,
				},
				OidcTriggerIntrospectionProjections: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_UNSPECIFIED,
				},
				OidcLegacyIntrospection: &feature.FeatureFlag{
					Enabled: true,
					Source:  feature.Source_SOURCE_SYSTEM,
				},
				UserSchema: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_UNSPECIFIED,
				},
				Actions: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_UNSPECIFIED,
				},
			},
		},
		{
			name: "some features, no inheritance",
			prepare: func(t *testing.T) {
				_, err := Client.SetInstanceFeatures(IamCTX, &feature.SetInstanceFeaturesRequest{
					LoginDefaultOrg:                     gu.Ptr(true),
					OidcTriggerIntrospectionProjections: gu.Ptr(false),
					UserSchema:                          gu.Ptr(true),
					Actions:                             gu.Ptr(true),
				})
				require.NoError(t, err)
			},
			args: args{
				ctx: IamCTX,
				req: &feature.GetInstanceFeaturesRequest{},
			},
			want: &feature.GetInstanceFeaturesResponse{
				LoginDefaultOrg: &feature.FeatureFlag{
					Enabled: true,
					Source:  feature.Source_SOURCE_INSTANCE,
				},
				OidcTriggerIntrospectionProjections: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_INSTANCE,
				},
				UserSchema: &feature.FeatureFlag{
					Enabled: true,
					Source:  feature.Source_SOURCE_INSTANCE,
				},
				Actions: &feature.FeatureFlag{
					Enabled: true,
					Source:  feature.Source_SOURCE_INSTANCE,
				},
			},
		},
		{
			name: "one feature, inheritance",
			prepare: func(t *testing.T) {
				_, err := Client.SetInstanceFeatures(IamCTX, &feature.SetInstanceFeaturesRequest{
					LoginDefaultOrg: gu.Ptr(true),
				})
				require.NoError(t, err)
			},
			args: args{
				ctx: IamCTX,
				req: &feature.GetInstanceFeaturesRequest{
					Inheritance: true,
				},
			},
			want: &feature.GetInstanceFeaturesResponse{
				LoginDefaultOrg: &feature.FeatureFlag{
					Enabled: true,
					Source:  feature.Source_SOURCE_INSTANCE,
				},
				OidcTriggerIntrospectionProjections: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_UNSPECIFIED,
				},
				OidcLegacyIntrospection: &feature.FeatureFlag{
					Enabled: true,
					Source:  feature.Source_SOURCE_SYSTEM,
				},
				UserSchema: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_UNSPECIFIED,
				},
				Actions: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_UNSPECIFIED,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				// make sure we have a clean state after each test
				_, err := Client.ResetInstanceFeatures(IamCTX, &feature.ResetInstanceFeaturesRequest{})
				require.NoError(t, err)
			})
			if tt.prepare != nil {
				tt.prepare(t)
			}
			got, err := Client.GetInstanceFeatures(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assertFeatureFlag(t, tt.want.LoginDefaultOrg, got.LoginDefaultOrg)
			assertFeatureFlag(t, tt.want.OidcTriggerIntrospectionProjections, got.OidcTriggerIntrospectionProjections)
			assertFeatureFlag(t, tt.want.OidcLegacyIntrospection, got.OidcLegacyIntrospection)
			assertFeatureFlag(t, tt.want.UserSchema, got.UserSchema)
		})
	}
}

func assertFeatureFlag(t *testing.T, expected, actual *feature.FeatureFlag) {
	t.Helper()
	assert.Equal(t, expected.GetEnabled(), actual.GetEnabled(), "enabled")
	assert.Equal(t, expected.GetSource(), actual.GetSource(), "source")
}
