package command

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/feature"
)

func Test_reduceInstanceFeature(t *testing.T) {
	type args struct {
		features *InstanceFeatures
		key      feature.Key
		value    any
	}
	tt := []struct {
		name     string
		args     args
		expected *InstanceFeatures
	}{
		{
			name: "key unspecified",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyUnspecified,
				value:    true,
			},
			expected: &InstanceFeatures{},
		},
		{
			name: "login default org",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyLoginDefaultOrg,
				value:    true,
			},
			expected: &InstanceFeatures{LoginDefaultOrg: gu.Ptr(true)},
		},
		{
			name: "token exchange",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyTokenExchange,
				value:    true,
			},
			expected: &InstanceFeatures{TokenExchange: gu.Ptr(true)},
		},
		{
			name: "user schema",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyUserSchema,
				value:    true,
			},
			expected: &InstanceFeatures{UserSchema: gu.Ptr(true)},
		},
		{
			name: "improved performance",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyImprovedPerformance,
				value:    []feature.ImprovedPerformanceType{feature.ImprovedPerformanceTypeProject},
			},
			expected: &InstanceFeatures{ImprovedPerformance: []feature.ImprovedPerformanceType{feature.ImprovedPerformanceTypeProject}},
		},
		{
			name: "debug OIDC parent error",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyDebugOIDCParentError,
				value:    true,
			},
			expected: &InstanceFeatures{DebugOIDCParentError: gu.Ptr(true)},
		},
		{
			name: "OIDC single v1 session termination",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyOIDCSingleV1SessionTermination,
				value:    true,
			},
			expected: &InstanceFeatures{OIDCSingleV1SessionTermination: gu.Ptr(true)},
		},
		{
			name: "enable back channel logout",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyEnableBackChannelLogout,
				value:    true,
			},
			expected: &InstanceFeatures{EnableBackChannelLogout: gu.Ptr(true)},
		},
		{
			name: "login v2",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyLoginV2,
				value:    &feature.LoginV2{},
			},
			expected: &InstanceFeatures{LoginV2: &feature.LoginV2{}},
		},
		{
			name: "permission check v2",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyPermissionCheckV2,
				value:    true,
			},
			expected: &InstanceFeatures{PermissionCheckV2: gu.Ptr(true)},
		},
		{
			name: "console use v2 user api",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyConsoleUseV2UserApi,
				value:    true,
			},
			expected: &InstanceFeatures{ConsoleUseV2UserApi: gu.Ptr(true)},
		},
		{
			name: "enable relational tables",
			args: args{
				features: &InstanceFeatures{},
				key:      feature.KeyEnableRelationalTables,
				value:    true,
			},
			expected: &InstanceFeatures{EnableRelationalTables: gu.Ptr(true)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Test
			reduceInstanceFeature(tc.args.features, tc.args.key, tc.args.value)

			// Verify
			assert.Equal(t, tc.expected, tc.args.features)
		})
	}
}
