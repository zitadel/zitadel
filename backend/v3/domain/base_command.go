package domain

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

// BaseCommand offers methods commonly used by all Commands
type BaseCommand struct{}

// TextOperationMapper maps gRPC TextQueryMethod to database.TextOperation
func (b *BaseCommand) TextOperationMapper(queryOperation object.TextQueryMethod) (database.TextOperation, error) {
	switch queryOperation {
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS:
		return database.TextOperationContains, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
		return database.TextOperationContainsIgnoreCase, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
		return database.TextOperationEndsWith, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		return database.TextOperationEndsWithIgnoreCase, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		return database.TextOperationEqual, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return database.TextOperationEqualIgnoreCase, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return database.TextOperationStartsWith, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return database.TextOperationStartsWithIgnoreCase, nil
	default:
		return 0, NewUnexpectedTextQueryOperationError("DOM-iBRBVe", queryOperation)
	}
}
