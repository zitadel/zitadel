package feature

import (
	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/query"
	feature_pb "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
)

func systemFeaturesToCommand(req *feature_pb.SetSystemFeaturesRequest) *command.SystemFeatures {
	return &command.SystemFeatures{
		LoginDefaultOrg:                 req.LoginDefaultOrg,
		TriggerIntrospectionProjections: req.OidcTriggerIntrospectionProjections,
		LegacyIntrospection:             req.OidcLegacyIntrospection,
		UserSchema:                      req.UserSchema,
		Actions:                         req.Actions,
		TokenExchange:                   req.OidcTokenExchange,
		ImprovedPerformance:             improvedPerformanceListToDomain(req.ImprovedPerformance),
	}
}

func systemFeaturesToPb(f *query.SystemFeatures) *feature_pb.GetSystemFeaturesResponse {
	return &feature_pb.GetSystemFeaturesResponse{
		Details:                             object.DomainToDetailsPb(f.Details),
		LoginDefaultOrg:                     featureSourceToFlagPb(&f.LoginDefaultOrg),
		OidcTriggerIntrospectionProjections: featureSourceToFlagPb(&f.TriggerIntrospectionProjections),
		OidcLegacyIntrospection:             featureSourceToFlagPb(&f.LegacyIntrospection),
		UserSchema:                          featureSourceToFlagPb(&f.UserSchema),
		OidcTokenExchange:                   featureSourceToFlagPb(&f.TokenExchange),
		Actions:                             featureSourceToFlagPb(&f.Actions),
		ImprovedPerformance:                 featureSourceToImprovedPerformanceFlagPb(&f.ImprovedPerformance),
	}
}

func instanceFeaturesToCommand(req *feature_pb.SetInstanceFeaturesRequest) *command.InstanceFeatures {
	return &command.InstanceFeatures{
		LoginDefaultOrg:                 req.LoginDefaultOrg,
		TriggerIntrospectionProjections: req.OidcTriggerIntrospectionProjections,
		LegacyIntrospection:             req.OidcLegacyIntrospection,
		UserSchema:                      req.UserSchema,
		TokenExchange:                   req.OidcTokenExchange,
		Actions:                         req.Actions,
		ImprovedPerformance:             improvedPerformanceListToDomain(req.ImprovedPerformance),
	}
}

func instanceFeaturesToPb(f *query.InstanceFeatures) *feature_pb.GetInstanceFeaturesResponse {
	return &feature_pb.GetInstanceFeaturesResponse{
		Details:                             object.DomainToDetailsPb(f.Details),
		LoginDefaultOrg:                     featureSourceToFlagPb(&f.LoginDefaultOrg),
		OidcTriggerIntrospectionProjections: featureSourceToFlagPb(&f.TriggerIntrospectionProjections),
		OidcLegacyIntrospection:             featureSourceToFlagPb(&f.LegacyIntrospection),
		UserSchema:                          featureSourceToFlagPb(&f.UserSchema),
		OidcTokenExchange:                   featureSourceToFlagPb(&f.TokenExchange),
		Actions:                             featureSourceToFlagPb(&f.Actions),
		ImprovedPerformance:                 featureSourceToImprovedPerformanceFlagPb(&f.ImprovedPerformance),
	}
}

func featureSourceToImprovedPerformanceFlagPb(fs *query.FeatureSource[[]feature.ImprovedPerformanceType]) *feature_pb.ImprovedPerformanceFeatureFlag {
	return &feature_pb.ImprovedPerformanceFeatureFlag{
		ExecutionPaths: improvedPerformanceTypesToPb(fs.Value),
		Source:         featureLevelToSourcePb(fs.Level),
	}
}

func featureSourceToFlagPb(fs *query.FeatureSource[bool]) *feature_pb.FeatureFlag {
	return &feature_pb.FeatureFlag{
		Enabled: fs.Value,
		Source:  featureLevelToSourcePb(fs.Level),
	}
}

func featureLevelToSourcePb(level feature.Level) feature_pb.Source {
	switch level {
	case feature.LevelUnspecified:
		return feature_pb.Source_SOURCE_UNSPECIFIED
	case feature.LevelSystem:
		return feature_pb.Source_SOURCE_SYSTEM
	case feature.LevelInstance:
		return feature_pb.Source_SOURCE_INSTANCE
	case feature.LevelOrg:
		return feature_pb.Source_SOURCE_ORGANIZATION
	case feature.LevelProject:
		return feature_pb.Source_SOURCE_PROJECT
	case feature.LevelApp:
		return feature_pb.Source_SOURCE_APP
	case feature.LevelUser:
		return feature_pb.Source_SOURCE_USER
	default:
		return feature_pb.Source(level)
	}
}

func improvedPerformanceTypesToPb(types []feature.ImprovedPerformanceType) []feature_pb.ImprovedPerformance {
	res := make([]feature_pb.ImprovedPerformance, len(types))

	for i, typ := range types {
		res[i] = improvedPerformanceTypeToPb(typ)
	}

	return res
}

func improvedPerformanceTypeToPb(typ feature.ImprovedPerformanceType) feature_pb.ImprovedPerformance {
	switch typ {
	case feature.ImprovedPerformanceTypeUnknown:
		return feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_UNSPECIFIED
	case feature.ImprovedPerformanceTypeOrgByID:
		return feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_BY_ID
	case feature.ImprovedPerformanceTypeProjectGrant:
		return feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_PROJECT_GRANT
	case feature.ImprovedPerformanceTypeProject:
		return feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_PROJECT
	case feature.ImprovedPerformanceTypeUserGrant:
		return feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_USER_GRANT
	case feature.ImprovedPerformanceTypeOrgDomainVerified:
		return feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED
	default:
		return feature_pb.ImprovedPerformance(typ)
	}
}

func improvedPerformanceListToDomain(list []feature_pb.ImprovedPerformance) []feature.ImprovedPerformanceType {
	if list == nil {
		return nil
	}
	res := make([]feature.ImprovedPerformanceType, len(list))

	for i, typ := range list {
		res[i] = improvedPerformanceToDomain(typ)
	}

	return res
}

func improvedPerformanceToDomain(typ feature_pb.ImprovedPerformance) feature.ImprovedPerformanceType {
	switch typ {
	case feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_UNSPECIFIED:
		return feature.ImprovedPerformanceTypeUnknown
	case feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_BY_ID:
		return feature.ImprovedPerformanceTypeOrgByID
	case feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_PROJECT_GRANT:
		return feature.ImprovedPerformanceTypeProjectGrant
	case feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_PROJECT:
		return feature.ImprovedPerformanceTypeProject
	case feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_USER_GRANT:
		return feature.ImprovedPerformanceTypeUserGrant
	case feature_pb.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED:
		return feature.ImprovedPerformanceTypeOrgDomainVerified
	default:
		return feature.ImprovedPerformanceTypeUnknown
	}
}
