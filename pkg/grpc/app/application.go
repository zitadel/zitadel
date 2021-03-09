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
		if !configType.OidcConfig.NoneCompliant {
			return nil
		}
		localizers := make([]middleware.Localizer, len(configType.OidcConfig.ComplianceProblems))
		for i, problem := range configType.OidcConfig.ComplianceProblems {
			localizers[i] = problem
		}
		return localizers
	}
	return nil
}

type AppConfig = isApp_Config
