package feature

import (
	"github.com/zitadel/zitadel/internal/command"
	feature "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
)

func systemFeaturesToCommand(req *feature.SetSystemFeaturesRequest) *command.SystemFeatures {
	return &command.SystemFeatures{
		LoginDefaultOrg:                 req.LoginDefaultOrg,
		TriggerIntrospectionProjections: req.OidcTriggerIntrospectionProjections,
		LegacyIntrospection:             req.OidcLegacyIntrospection,
	}
}

func instanceFeaturesToCommand(req *feature.SetInstanceFeaturesRequest) *command.InstanceFeatures {
	return &command.InstanceFeatures{
		LoginDefaultOrg:                 req.LoginDefaultOrg,
		TriggerIntrospectionProjections: req.OidcTriggerIntrospectionProjections,
		LegacyIntrospection:             req.OidcLegacyIntrospection,
	}
}
