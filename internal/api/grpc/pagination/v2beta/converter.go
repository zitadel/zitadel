package pagination

import (
	"fmt"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	pagination "github.com/zitadel/zitadel/pkg/grpc/pagination/v2beta"
)

func SearchQueryPbToQuery(defaults systemdefaults.SystemDefaults, query *pagination.ListQuery) (offset, limit uint64, asc bool, err error) {
	limit = defaults.DefaultQueryLimit
	if query == nil {
		return 0, limit, asc, nil
	}
	offset = query.Offset
	asc = query.Asc
	if defaults.MaxQueryLimit > 0 && uint64(query.Limit) > defaults.MaxQueryLimit {
		return 0, 0, false, zerrors.ThrowInvalidArgumentf(fmt.Errorf("given: %d, allowed: %d", query.Limit, defaults.MaxQueryLimit), "QUERY-4M0fs", "Errors.Query.LimitExceeded")
	}
	if query.Limit > 0 {
		limit = uint64(query.Limit)
	}
	return offset, limit, asc, nil
}

func ToSearchDetailsPb(request query.SearchRequest, response query.SearchResponse) *pagination.ListDetails {
	return &pagination.ListDetails{
		AppliedLimit: request.Limit,
		TotalResult:  response.Count,
	}
}
