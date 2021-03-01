package management

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func searchMethodToModel(method management.SearchMethod) domain.SearchMethod {
	switch method {
	case management.SearchMethod_SEARCHMETHOD_EQUALS:
		return domain.SearchMethodEquals
	case management.SearchMethod_SEARCHMETHOD_CONTAINS:
		return domain.SearchMethodContains
	case management.SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return domain.SearchMethodStartsWith
	case management.SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return domain.SearchMethodEqualsIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return domain.SearchMethodContainsIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return domain.SearchMethodStartsWithIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_NOT_EQUALS:
		return domain.SearchMethodNotEquals
	case management.SearchMethod_SEARCHMETHOD_IS_ONE_OF:
		return domain.SearchMethodIsOneOf
	case management.SearchMethod_SEARCHMETHOD_LIST_CONTAINS:
		return domain.SearchMethodListContains
	default:
		return domain.SearchMethodEquals
	}
}
