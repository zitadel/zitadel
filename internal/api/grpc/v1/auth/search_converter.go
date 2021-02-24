package auth

import (
	"github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func searchMethodToModel(method auth.SearchMethod) model.SearchMethod {
	switch method {
	case auth.SearchMethod_SEARCHMETHOD_EQUALS:
		return model.SearchMethodEquals
	case auth.SearchMethod_SEARCHMETHOD_CONTAINS:
		return model.SearchMethodContains
	case auth.SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return model.SearchMethodStartsWith
	case auth.SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return model.SearchMethodEqualsIgnoreCase
	case auth.SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return model.SearchMethodContainsIgnoreCase
	case auth.SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return model.SearchMethodStartsWithIgnoreCase
	default:
		return model.SearchMethodEquals
	}
}
