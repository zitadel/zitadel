package app

import (
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
)

func (a *App) Localizers() []middleware.Localizer {
	if a == nil {
		return nil
	}

	switch configType := a.Config.(type) {
	case *App_OidcConfig:
		return configType.ComplianceLocalizers()
	}
	return nil
}

func (o *App_OidcConfig) ComplianceLocalizers() []middleware.Localizer {
	if o.OidcConfig == nil {
		return nil
	}

	if !o.OidcConfig.NoneCompliant {
		return nil
	}
	localizers := make([]middleware.Localizer, len(o.OidcConfig.ComplianceProblems))
	for i, problem := range o.OidcConfig.ComplianceProblems {
		localizers[i] = problem
	}
	return localizers
}

type AppConfig = isApp_Config
