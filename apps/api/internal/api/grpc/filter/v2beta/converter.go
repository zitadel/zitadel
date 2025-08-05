package filter

import (
	"fmt"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
)

func TextMethodPbToQuery(method filter.TextFilterMethod) query.TextComparison {
	switch method {
	case filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS:
		return query.TextEquals
	case filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS_IGNORE_CASE:
		return query.TextEqualsIgnoreCase
	case filter.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH:
		return query.TextStartsWith
	case filter.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH_IGNORE_CASE:
		return query.TextStartsWithIgnoreCase
	case filter.TextFilterMethod_TEXT_FILTER_METHOD_CONTAINS:
		return query.TextContains
	case filter.TextFilterMethod_TEXT_FILTER_METHOD_CONTAINS_IGNORE_CASE:
		return query.TextContainsIgnoreCase
	case filter.TextFilterMethod_TEXT_FILTER_METHOD_ENDS_WITH:
		return query.TextEndsWith
	case filter.TextFilterMethod_TEXT_FILTER_METHOD_ENDS_WITH_IGNORE_CASE:
		return query.TextEndsWithIgnoreCase
	default:
		return -1
	}
}

func TimestampMethodPbToQuery(method filter.TimestampFilterMethod) query.TimestampComparison {
	switch method {
	case filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_EQUALS:
		return query.TimestampEquals
	case filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_LESS:
		return query.TimestampLess
	case filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_GREATER:
		return query.TimestampGreater
	case filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_LESS_OR_EQUALS:
		return query.TimestampLessOrEquals
	case filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_GREATER_OR_EQUALS:
		return query.TimestampGreaterOrEquals
	default:
		return -1
	}
}

func PaginationPbToQuery(defaults systemdefaults.SystemDefaults, query *filter.PaginationRequest) (offset, limit uint64, asc bool, err error) {
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

func QueryToPaginationPb(request query.SearchRequest, response query.SearchResponse) *filter.PaginationResponse {
	return &filter.PaginationResponse{
		AppliedLimit: request.Limit,
		TotalResult:  response.Count,
	}
}
