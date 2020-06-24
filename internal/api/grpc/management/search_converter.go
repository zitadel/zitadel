package management

import (
	"github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/pkg/management/grpc"
)

func searchMethodToModel(method grpc.SearchMethod) model.SearchMethod {
	switch method {
	case grpc.SearchMethod_SEARCHMETHOD_EQUALS:
		return model.SearchMethodEquals
	case grpc.SearchMethod_SEARCHMETHOD_CONTAINS:
		return model.SearchMethodContains
	case grpc.SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return model.SearchMethodStartsWith
	case grpc.SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return model.SearchMethodEqualsIgnoreCase
	case grpc.SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return model.SearchMethodContainsIgnoreCase
	case grpc.SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return model.SearchMethodStartsWithIgnoreCase
	case grpc.SearchMethod_SEARCHMETHOD_NOT_EQUALS:
		return model.SearchMethodNotEquals
	case grpc.SearchMethod_SEARCHMETHOD_IS_ONE_OF:
		return model.SearchMethodIsOneOf
	case grpc.SearchMethod_SEARCHMETHOD_LIST_CONTAINS:
		return model.SearchMethodListContains
	default:
		return model.SearchMethodEquals
	}
}
