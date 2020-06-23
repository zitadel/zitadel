package grpc

import "github.com/caos/zitadel/internal/model"

func searchMethodToModel(method SearchMethod) model.SearchMethod {
	switch method {
	case SearchMethod_SEARCHMETHOD_EQUALS:
		return model.SearchMethodEquals
	case SearchMethod_SEARCHMETHOD_CONTAINS:
		return model.SearchMethodContains
	case SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return model.SearchMethodStartsWith
	case SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return model.SearchMethodEqualsIgnoreCase
	case SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return model.SearchMethodContainsIgnoreCase
	case SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return model.SearchMethodStartsWithIgnoreCase
	case SearchMethod_SEARCHMETHOD_NOT_EQUALS:
		return model.SearchMethodNotEquals
	case SearchMethod_SEARCHMETHOD_IS_ONE_OF:
		return model.SearchMethodIsOneOf
	case SearchMethod_SEARCHMETHOD_LIST_CONTAINS:
		return model.SearchMethodListContains
	default:
		return model.SearchMethodEquals
	}
}
