package grpc

import "github.com/caos/zitadel/internal/model"

func searchMethodToModel(method SearchMethod) model.SearchMethod {
	switch method {
	case SearchMethod_SEARCHMETHOD_EQUALS:
		return model.SEARCHMETHOD_EQUALS
	case SearchMethod_SEARCHMETHOD_CONTAINS:
		return model.SEARCHMETHOD_CONTAINS
	case SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return model.SEARCHMETHOD_STARTS_WITH
	case SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return model.SEARCHMETHOD_EQUALS_IGNORE_CASE
	case SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return model.SEARCHMETHOD_CONTAINS_IGNORE_CASE
	case SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return model.SEARCHMETHOD_STARTS_WITH_IGNORE_CASE
	case SearchMethod_SEARCHMETHOD_NOT_EQUALS:
		return model.SEARCHMETHOD_NOT_EQUALS
	case SearchMethod_SEARCHMETHOD_IS_ONE_OF:
		return model.SEARCHMETHOD_IS_ONE_OF
	case SearchMethod_SEARCHMETHOD_LIST_CONTAINS:
		return model.SEARCHMETHOD_LIST_CONTAINS
	default:
		return model.SEARCHMETHOD_EQUALS
	}
}
