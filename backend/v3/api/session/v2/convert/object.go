package convert

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	objpb "github.com/zitadel/zitadel/pkg/grpc/object"
)

var grpcTimestampOpToDomain = map[objpb.TimestampQueryMethod]database.NumberOperation{
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_EQUALS:            database.NumberOperationEqual,
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER:           database.NumberOperationGreaterThan,
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS: database.NumberOperationGreaterThanOrEqual,
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS:              database.NumberOperationLessThan,
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS:    database.NumberOperationLessThanOrEqual,
}
