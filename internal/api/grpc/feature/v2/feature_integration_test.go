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
	"github.com/zitadel/zitadel/internal/integration"
	feature "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client feature.FeatureServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(5 * time.Minute)
		CTX = ctx
		defer cancel()

		Tester = integration.NewTester(ctx)
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
				ctx: Tester.WithAuthorization(CTX, integration.IAMOwner),
				req: &feature.SetSystemFeaturesRequest{
					OidcTriggerIntrospectionProjections: gu.Ptr(true),
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: Tester.WithAuthorization(CTX, integration.SystemUser),
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
				_, err := Client.ResetSystemFeatures(Tester.WithAuthorization(CTX, integration.SystemUser), &feature.ResetSystemFeaturesRequest{})
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
	_, err := Client.SetSystemFeatures(Tester.WithAuthorization(CTX, integration.SystemUser), &feature.SetSystemFeaturesRequest{
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
			ctx:     Tester.WithAuthorization(CTX, integration.IAMOwner),
			wantErr: true,
		},
		{
			name: "success",
			ctx:  Tester.WithAuthorization(CTX, integration.SystemUser),
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
				ctx: Tester.WithAuthorization(CTX, integration.IAMOwner),
				req: &feature.GetSystemFeaturesRequest{},
			},
			wantErr: true,
		},
		{
			name: "defaults, no inheritance",
			args: args{
				ctx: Tester.WithAuthorization(CTX, integration.SystemUser),
				req: &feature.GetSystemFeaturesRequest{},
			},
			want: &feature.GetSystemFeaturesResponse{},
		},
		{
			name: "defaults, inheritance",
			args: args{
				ctx: Tester.WithAuthorization(CTX, integration.SystemUser),
				req: &feature.GetSystemFeaturesRequest{
					Inheritance: true,
				},
			},
			want: &feature.GetSystemFeaturesResponse{
				LoginDefaultOrg: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_DEFAULT,
				},
				OidcTriggerIntrospectionProjections: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_DEFAULT,
				},
				OidcLegacyIntrospection: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_DEFAULT,
				},
			},
		},
		{
			name: "some features, no inheritance",
			prepare: func(t *testing.T) {
				_, err := Client.SetSystemFeatures(Tester.WithAuthorization(CTX, integration.SystemUser), &feature.SetSystemFeaturesRequest{
					LoginDefaultOrg:                     gu.Ptr(true),
					OidcTriggerIntrospectionProjections: gu.Ptr(false),
				})
				require.NoError(t, err)
			},
			args: args{
				ctx: Tester.WithAuthorization(CTX, integration.SystemUser),
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
		{
			name: "some features, inheritance",
			prepare: func(t *testing.T) {
				_, err := Client.SetSystemFeatures(Tester.WithAuthorization(CTX, integration.SystemUser), &feature.SetSystemFeaturesRequest{
					LoginDefaultOrg:                     gu.Ptr(true),
					OidcTriggerIntrospectionProjections: gu.Ptr(false),
				})
				require.NoError(t, err)
			},
			args: args{
				ctx: Tester.WithAuthorization(CTX, integration.SystemUser),
				req: &feature.GetSystemFeaturesRequest{
					Inheritance: true,
				},
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
				OidcLegacyIntrospection: &feature.FeatureFlag{
					Enabled: false,
					Source:  feature.Source_SOURCE_DEFAULT,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				// make sure we have a clean state after each test
				_, err := Client.ResetSystemFeatures(Tester.WithAuthorization(CTX, integration.SystemUser), &feature.ResetSystemFeaturesRequest{})
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
		})
	}
}

func assertFeatureFlag(t *testing.T, expected, actual *feature.FeatureFlag) {
	t.Helper()
	assert.Equal(t, expected.GetEnabled(), actual.GetEnabled(), "enabled")
	assert.Equal(t, expected.GetSource(), actual.GetSource(), "source")
}
