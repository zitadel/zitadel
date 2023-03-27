package management

import "github.com/zitadel/zitadel/internal/api/grpc/server/middleware"

func (r *ListFlowTypesResponse) Localizers() (localizers []middleware.Localizer) {
	if r == nil {
		return nil
	}

	localizers = make([]middleware.Localizer, 0, len(r.Result))
	for _, typ := range r.Result {
		localizers = append(localizers, typ.Localizers()...)
	}

	return localizers
}

func (r *ListFlowTriggerTypesResponse) Localizers() (localizers []middleware.Localizer) {
	if r == nil {
		return nil
	}

	localizers = make([]middleware.Localizer, 0, len(r.Result))
	for _, typ := range r.Result {
		localizers = append(localizers, typ.Localizers()...)
	}

	return localizers
}

func (r *GetFlowResponse) Localizers() []middleware.Localizer {
	if r == nil {
		return nil
	}

	return r.Flow.Localizers()
}
