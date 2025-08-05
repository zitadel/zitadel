package api

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
)

func V2BetaPaginationToDomain(pagination *filter.PaginationRequest) domain.Pagination {
	return domain.Pagination{
		Limit:     pagination.Limit,
		Offset:    uint32(pagination.Offset),
		Ascending: pagination.Asc,
	}
}
