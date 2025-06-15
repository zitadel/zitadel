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
		TokenExchange:                   req.OidcTokenExchange,
		OIDCSingleV1SessionTermination:  req.OidcSingleV1SessionTermination,
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
		OidcSingleV1SessionTermination:      featureSourceToFlagPb(&f.OIDCSingleV1SessionTermination),
	}
}

func instanceFeaturesToCommand(req *feature_pb.SetInstanceFeaturesRequest) *command.InstanceFeatures {
	return &command.InstanceFeatures{
		LoginDefaultOrg:                 req.LoginDefaultOrg,
		TriggerIntrospectionProjections: req.OidcTriggerIntrospectionProjections,
		LegacyIntrospection:             req.OidcLegacyIntrospection,
		UserSchema:                      req.UserSchema,
		TokenExchange:                   req.OidcTokenExchange,
		WebKey:                          req.WebKey,
		DebugOIDCParentError:            req.DebugOidcParentError,
		OIDCSingleV1SessionTermination:  req.OidcSingleV1SessionTermination,
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
		WebKey:                              featureSourceToFlagPb(&f.WebKey),
		DebugOidcParentError:                featureSourceToFlagPb(&f.DebugOIDCParentError),
		OidcSingleV1SessionTermination:      featureSourceToFlagPb(&f.OIDCSingleV1SessionTermination),
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
