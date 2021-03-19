package management

import "github.com/caos/zitadel/internal/api/grpc/server/middleware"

func (a *AddOIDCAppResponse) Localizers() []middleware.Localizer {
	if a == nil {
		return nil
	}

	if !a.NoneCompliant {
		return nil
	}
	localizers := make([]middleware.Localizer, len(a.ComplianceProblems))
	for i, problem := range a.ComplianceProblems {
		localizers[i] = problem
	}
	return localizers

}
