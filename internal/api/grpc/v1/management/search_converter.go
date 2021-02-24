package management

import (
	"github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func searchMethodToModel(method management.SearchMethod) model.SearchMethod {
	switch method {
	case management.SearchMethod_SEARCHMETHOD_EQUALS:
		return model.SearchMethodEquals
	case management.SearchMethod_SEARCHMETHOD_CONTAINS:
		return model.SearchMethodContains
	case management.SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return model.SearchMethodStartsWith
	case management.SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return model.SearchMethodEqualsIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return model.SearchMethodContainsIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return model.SearchMethodStartsWithIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_NOT_EQUALS:
		return model.SearchMethodNotEquals
	case management.SearchMethod_SEARCHMETHOD_IS_ONE_OF:
		return model.SearchMethodIsOneOf
	case management.SearchMethod_SEARCHMETHOD_LIST_CONTAINS:
		return model.SearchMethodListContains
	default:
		return model.SearchMethodEquals
	}
}
