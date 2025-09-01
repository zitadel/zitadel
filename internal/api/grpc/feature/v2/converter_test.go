package feature

import (
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/query"
	feature_pb "github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func Test_systemFeaturesToCommand(t *testing.T) {
	arg := &feature_pb.SetSystemFeaturesRequest{
		LoginDefaultOrg:                gu.Ptr(true),
		UserSchema:                     gu.Ptr(true),
		OidcTokenExchange:              gu.Ptr(true),
		ImprovedPerformance:            nil,
		OidcSingleV1SessionTermination: gu.Ptr(true),
		LoginV2: &feature_pb.LoginV2{
			Required: true,
			BaseUri:  gu.Ptr("https://login.com"),
		},
	}
	want := &command.SystemFeatures{
		LoginDefaultOrg:                gu.Ptr(true),
		UserSchema:                     gu.Ptr(true),
		TokenExchange:                  gu.Ptr(true),
		ImprovedPerformance:            nil,
		OIDCSingleV1SessionTermination: gu.Ptr(true),
		LoginV2: &feature.LoginV2{
			Required: true,
			BaseURI:  &url.URL{Scheme: "https", Host: "login.com"},
		},
	}
	got, err := systemFeaturesToCommand(arg)
	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func Test_systemFeaturesToPb(t *testing.T) {
	arg := &query.SystemFeatures{
		Details: &domain.ObjectDetails{
			Sequence:      22,
			EventDate:     time.Unix(123, 0),
			ResourceOwner: "SYSTEM",
		},
		LoginDefaultOrg: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: true,
		},
		UserSchema: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: true,
		},
		TokenExchange: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: false,
		},
		ImprovedPerformance: query.FeatureSource[[]feature.ImprovedPerformanceType]{
			Level: feature.LevelSystem,
			Value: []feature.ImprovedPerformanceType{feature.ImprovedPerformanceTypeOrgDomainVerified},
		},
		OIDCSingleV1SessionTermination: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: true,
		},
		EnableBackChannelLogout: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: true,
		},
		LoginV2: query.FeatureSource[*feature.LoginV2]{
			Level: feature.LevelSystem,
			Value: &feature.LoginV2{
				Required: true,
				BaseURI:  &url.URL{Scheme: "https", Host: "login.com"},
			},
		},
		PermissionCheckV2: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: true,
		},
	}
	want := &feature_pb.GetSystemFeaturesResponse{
		Details: &object.Details{
			Sequence:      22,
			ChangeDate:    &timestamppb.Timestamp{Seconds: 123},
			ResourceOwner: "SYSTEM",
		},
		LoginDefaultOrg: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
		UserSchema: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
		OidcTokenExchange: &feature_pb.FeatureFlag{
			Enabled: false,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
		ImprovedPerformance: &feature_pb.ImprovedPerformanceFeatureFlag{
			ExecutionPaths: []feature_pb.ImprovedPerformance{feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED},
			Source:         feature_pb.Source_SOURCE_SYSTEM,
		},
		OidcSingleV1SessionTermination: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
		EnableBackChannelLogout: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
		LoginV2: &feature_pb.LoginV2FeatureFlag{
			Required: true,
			BaseUri:  gu.Ptr("https://login.com"),
			Source:   feature_pb.Source_SOURCE_SYSTEM,
		},
		PermissionCheckV2: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
	}
	got := systemFeaturesToPb(arg)
	assert.Equal(t, want, got)
}

func Test_instanceFeaturesToCommand(t *testing.T) {
	arg := &feature_pb.SetInstanceFeaturesRequest{
		LoginDefaultOrg:                gu.Ptr(true),
		UserSchema:                     gu.Ptr(true),
		OidcTokenExchange:              gu.Ptr(true),
		ImprovedPerformance:            nil,
		DebugOidcParentError:           gu.Ptr(true),
		OidcSingleV1SessionTermination: gu.Ptr(true),
		EnableBackChannelLogout:        gu.Ptr(true),
		LoginV2: &feature_pb.LoginV2{
			Required: true,
			BaseUri:  gu.Ptr("https://login.com"),
		},
		ConsoleUseV2UserApi: gu.Ptr(true),
	}
	want := &command.InstanceFeatures{
		LoginDefaultOrg:                gu.Ptr(true),
		UserSchema:                     gu.Ptr(true),
		TokenExchange:                  gu.Ptr(true),
		ImprovedPerformance:            nil,
		DebugOIDCParentError:           gu.Ptr(true),
		OIDCSingleV1SessionTermination: gu.Ptr(true),
		EnableBackChannelLogout:        gu.Ptr(true),
		LoginV2: &feature.LoginV2{
			Required: true,
			BaseURI:  &url.URL{Scheme: "https", Host: "login.com"},
		},
		ConsoleUseV2UserApi: gu.Ptr(true),
	}
	got, err := instanceFeaturesToCommand(arg)
	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func Test_instanceFeaturesToPb(t *testing.T) {
	arg := &query.InstanceFeatures{
		Details: &domain.ObjectDetails{
			Sequence:      22,
			EventDate:     time.Unix(123, 0),
			ResourceOwner: "instance1",
		},
		LoginDefaultOrg: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: true,
		},
		UserSchema: query.FeatureSource[bool]{
			Level: feature.LevelInstance,
			Value: true,
		},
		TokenExchange: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: false,
		},
		ImprovedPerformance: query.FeatureSource[[]feature.ImprovedPerformanceType]{
			Level: feature.LevelSystem,
			Value: []feature.ImprovedPerformanceType{feature.ImprovedPerformanceTypeOrgDomainVerified},
		},
		OIDCSingleV1SessionTermination: query.FeatureSource[bool]{
			Level: feature.LevelInstance,
			Value: true,
		},
		EnableBackChannelLogout: query.FeatureSource[bool]{
			Level: feature.LevelInstance,
			Value: true,
		},
		LoginV2: query.FeatureSource[*feature.LoginV2]{
			Level: feature.LevelInstance,
			Value: &feature.LoginV2{
				Required: true,
				BaseURI:  &url.URL{Scheme: "https", Host: "login.com"},
			},
		},
		PermissionCheckV2: query.FeatureSource[bool]{
			Level: feature.LevelInstance,
			Value: true,
		},
		ConsoleUseV2UserApi: query.FeatureSource[bool]{
			Level: feature.LevelInstance,
			Value: true,
		},
	}
	want := &feature_pb.GetInstanceFeaturesResponse{
		Details: &object.Details{
			Sequence:      22,
			ChangeDate:    &timestamppb.Timestamp{Seconds: 123},
			ResourceOwner: "instance1",
		},
		LoginDefaultOrg: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
		UserSchema: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_INSTANCE,
		},
		OidcTokenExchange: &feature_pb.FeatureFlag{
			Enabled: false,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
		ImprovedPerformance: &feature_pb.ImprovedPerformanceFeatureFlag{
			ExecutionPaths: []feature_pb.ImprovedPerformance{feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED},
			Source:         feature_pb.Source_SOURCE_SYSTEM,
		},
		DebugOidcParentError: &feature_pb.FeatureFlag{
			Enabled: false,
			Source:  feature_pb.Source_SOURCE_UNSPECIFIED,
		},
		OidcSingleV1SessionTermination: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_INSTANCE,
		},
		EnableBackChannelLogout: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_INSTANCE,
		},
		LoginV2: &feature_pb.LoginV2FeatureFlag{
			Required: true,
			BaseUri:  gu.Ptr("https://login.com"),
			Source:   feature_pb.Source_SOURCE_INSTANCE,
		},
		PermissionCheckV2: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_INSTANCE,
		},
		ConsoleUseV2UserApi: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_INSTANCE,
		},
	}
	got := instanceFeaturesToPb(arg)
	assert.Equal(t, want, got)
}

func Test_featureLevelToSourcePb(t *testing.T) {
	tests := []struct {
		name  string
		level feature.Level
		want  feature_pb.Source
	}{
		{
			name:  "unspecified",
			level: feature.LevelUnspecified,
			want:  feature_pb.Source_SOURCE_UNSPECIFIED,
		},
		{
			name:  "system",
			level: feature.LevelSystem,
			want:  feature_pb.Source_SOURCE_SYSTEM,
		},
		{
			name:  "instance",
			level: feature.LevelInstance,
			want:  feature_pb.Source_SOURCE_INSTANCE,
		},
		{
			name:  "org",
			level: feature.LevelOrg,
			want:  feature_pb.Source_SOURCE_ORGANIZATION,
		},
		{
			name:  "project",
			level: feature.LevelProject,
			want:  feature_pb.Source_SOURCE_PROJECT,
		},
		{
			name:  "app",
			level: feature.LevelApp,
			want:  feature_pb.Source_SOURCE_APP,
		},
		{
			name:  "user",
			level: feature.LevelUser,
			want:  feature_pb.Source_SOURCE_USER,
		},
		{
			name:  "unknown",
			level: 99,
			want:  99,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := featureLevelToSourcePb(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}
