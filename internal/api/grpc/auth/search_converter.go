package auth

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func searchMethodToModel(method auth.SearchMethod) domain.SearchMethod {
	switch method {
	case auth.SearchMethod_SEARCHMETHOD_EQUALS:
		return domain.SearchMethodEquals
	case auth.SearchMethod_SEARCHMETHOD_CONTAINS:
		return domain.SearchMethodContains
	case auth.SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return domain.SearchMethodStartsWith
	case auth.SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return domain.SearchMethodEqualsIgnoreCase
	case auth.SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return domain.SearchMethodContainsIgnoreCase
	case auth.SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return domain.SearchMethodStartsWithIgnoreCase
	default:
		return domain.SearchMethodEquals
	}
}
