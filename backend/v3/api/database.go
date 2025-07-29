package api

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v2beta "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

func V2BetaTextFilterToDatabase(filter v2beta.TextQueryMethod) database.TextOperation {
	switch filter {
	case v2beta.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		return database.TextOperationEqual
	case v2beta.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return database.TextOperationEqualIgnoreCase
	case v2beta.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return database.TextOperationStartsWith
	case v2beta.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return database.TextOperationStartsWithIgnoreCase
	case v2beta.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS, v2beta.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE, v2beta.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH, v2beta.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		panic("unimplemented text query method: " + filter.String())
	default:
		panic("unknown text query method: " + filter.String())
	}
}
