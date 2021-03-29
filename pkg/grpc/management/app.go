package management

import (
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
)

func (a *ListAppsResponse) Localizers() []middleware.Localizer {
	if a == nil {
		return nil
	}

	localizers := make([]middleware.Localizer, 0)
	for _, a := range a.Result {
		localizers = append(localizers, a.Localizers()...)
	}
	return localizers
}

func (a *GetAppByIDResponse) Localizers() []middleware.Localizer {
	if a == nil || (a != nil && a.App == nil) {
		return nil
	}

	return a.App.Localizers()
}

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
