package api

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func V2BetaOrgStateToDomain(state org.OrgState) domain.OrgState {
	switch state {
	case org.OrgState_ORG_STATE_ACTIVE:
		return domain.OrgStateActive
	case org.OrgState_ORG_STATE_INACTIVE:
		return domain.OrgStateInactive
	default:
		// TODO: removed is not supported in the domain
		panic("unknown org state: " + state.String())
	}
}

func V2BetaPaginationToDomain(pagination *filter.PaginationRequest) domain.Pagination {
	return domain.Pagination{
		Limit:     pagination.Limit,
		Offset:    uint32(pagination.Offset),
		Ascending: pagination.Asc,
	}
}
