package feature

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/query"
	feature_pb "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

func Test_systemFeaturesToCommand(t *testing.T) {
	arg := &feature_pb.SetSystemFeaturesRequest{
		LoginDefaultOrg:                     gu.Ptr(true),
		OidcTriggerIntrospectionProjections: gu.Ptr(false),
		OidcLegacyIntrospection:             nil,
		UserSchema:                          gu.Ptr(true),
		Actions:                             gu.Ptr(true),
		OidcTokenExchange:                   gu.Ptr(true),
		ImprovedPerformance:                 nil,
	}
	want := &command.SystemFeatures{
		LoginDefaultOrg:                 gu.Ptr(true),
		TriggerIntrospectionProjections: gu.Ptr(false),
		LegacyIntrospection:             nil,
		UserSchema:                      gu.Ptr(true),
		Actions:                         gu.Ptr(true),
		TokenExchange:                   gu.Ptr(true),
		ImprovedPerformance:             nil,
	}
	got := systemFeaturesToCommand(arg)
	assert.Equal(t, want, got)
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
		TriggerIntrospectionProjections: query.FeatureSource[bool]{
			Level: feature.LevelUnspecified,
			Value: false,
		},
		LegacyIntrospection: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: true,
		},
		UserSchema: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: true,
		},
		Actions: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: true,
		},
		TokenExchange: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: false,
		},
		ImprovedPerformance: query.FeatureSource[[]feature.ImprovedPerformanceType]{
			Level: feature.LevelSystem,
			Value: []feature.ImprovedPerformanceType{feature.ImprovedPerformanceTypeOrgByID},
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
		OidcTriggerIntrospectionProjections: &feature_pb.FeatureFlag{
			Enabled: false,
			Source:  feature_pb.Source_SOURCE_UNSPECIFIED,
		},
		OidcLegacyIntrospection: &feature_pb.FeatureFlag{
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
		Actions: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
		ImprovedPerformance: &feature_pb.ImprovedPerformanceFeatureFlag{
			ExecutionPaths: []feature_pb.ImprovedPerformance{feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_BY_ID},
			Source:         feature_pb.Source_SOURCE_SYSTEM,
		},
	}
	got := systemFeaturesToPb(arg)
	assert.Equal(t, want, got)
}

func Test_instanceFeaturesToCommand(t *testing.T) {
	arg := &feature_pb.SetInstanceFeaturesRequest{
		LoginDefaultOrg:                     gu.Ptr(true),
		OidcTriggerIntrospectionProjections: gu.Ptr(false),
		OidcLegacyIntrospection:             nil,
		UserSchema:                          gu.Ptr(true),
		OidcTokenExchange:                   gu.Ptr(true),
		Actions:                             gu.Ptr(true),
		ImprovedPerformance:                 nil,
	}
	want := &command.InstanceFeatures{
		LoginDefaultOrg:                 gu.Ptr(true),
		TriggerIntrospectionProjections: gu.Ptr(false),
		LegacyIntrospection:             nil,
		UserSchema:                      gu.Ptr(true),
		TokenExchange:                   gu.Ptr(true),
		Actions:                         gu.Ptr(true),
		ImprovedPerformance:             nil,
	}
	got := instanceFeaturesToCommand(arg)
	assert.Equal(t, want, got)
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
		TriggerIntrospectionProjections: query.FeatureSource[bool]{
			Level: feature.LevelUnspecified,
			Value: false,
		},
		LegacyIntrospection: query.FeatureSource[bool]{
			Level: feature.LevelInstance,
			Value: true,
		},
		UserSchema: query.FeatureSource[bool]{
			Level: feature.LevelInstance,
			Value: true,
		},
		Actions: query.FeatureSource[bool]{
			Level: feature.LevelInstance,
			Value: true,
		},
		TokenExchange: query.FeatureSource[bool]{
			Level: feature.LevelSystem,
			Value: false,
		},
		ImprovedPerformance: query.FeatureSource[[]feature.ImprovedPerformanceType]{
			Level: feature.LevelSystem,
			Value: []feature.ImprovedPerformanceType{feature.ImprovedPerformanceTypeOrgByID},
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
		OidcTriggerIntrospectionProjections: &feature_pb.FeatureFlag{
			Enabled: false,
			Source:  feature_pb.Source_SOURCE_UNSPECIFIED,
		},
		OidcLegacyIntrospection: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_INSTANCE,
		},
		UserSchema: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_INSTANCE,
		},
		Actions: &feature_pb.FeatureFlag{
			Enabled: true,
			Source:  feature_pb.Source_SOURCE_INSTANCE,
		},
		OidcTokenExchange: &feature_pb.FeatureFlag{
			Enabled: false,
			Source:  feature_pb.Source_SOURCE_SYSTEM,
		},
		ImprovedPerformance: &feature_pb.ImprovedPerformanceFeatureFlag{
			ExecutionPaths: []feature_pb.ImprovedPerformance{feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_BY_ID},
			Source:         feature_pb.Source_SOURCE_SYSTEM,
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
