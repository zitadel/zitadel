package object

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
	v2_object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	v2beta_object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

// TODO(IAM-Marco): Remove in V5 (see https://github.com/zitadel/zitadel/issues/10877)
func TextQueryMethodBetaToV2(txtMethod v2beta_object.TextQueryMethod) v2_object.TextQueryMethod {
	switch txtMethod {
	case v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS:
		return v2_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS
	case v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
		return v2_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE
	case v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
		return v2_object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH
	case v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		return v2_object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE
	case v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return v2_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE
	case v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return v2_object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH
	case v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return v2_object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE
	case v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		fallthrough
	default:
		return v2_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS
	}
}

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
