package object

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
	v2_object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func TextQueryMethodToTextOperation(txtMethod v2_object.TextQueryMethod) (database.TextOperation, error) {
	switch txtMethod {
	case v2_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS:
		return database.TextOperationContains, nil
	case v2_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
		return database.TextOperationContainsIgnoreCase, nil
	case v2_object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
		return database.TextOperationEndsWith, nil
	case v2_object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		return database.TextOperationEndsWithIgnoreCase, nil
	case v2_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return database.TextOperationEqualIgnoreCase, nil
	case v2_object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return database.TextOperationStartsWith, nil
	case v2_object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return database.TextOperationStartsWithIgnoreCase, nil
	case v2_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		return database.TextOperationEqual, nil
	default:
		return 0, zerrors.ThrowInvalidArgument(nil, "OBJ-iBRBVe", "invalid text query method")
	}
}
