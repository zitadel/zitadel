package filter

import (
	"github.com/zitadel/zitadel/internal/query"
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
