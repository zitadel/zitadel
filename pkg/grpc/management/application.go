package management

import (
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
)

func (a *ApplicationView) Localizers() []middleware.Localizer {
	if a == nil {
		return nil
	}

	switch configType := a.AppConfig.(type) {
	case *ApplicationView_OidcConfig:
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
