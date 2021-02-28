package object

import (
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/pkg/grpc/object"
	object_pb "github.com/caos/zitadel/pkg/grpc/object"
	"github.com/golang/protobuf/ptypes"
)

func ToDetailsPb(
	sequence uint64,
	creationDate,
	changeDate time.Time,
	resourceOwner string,
) *object_pb.ObjectDetails {
	creationDatePb, err := ptypes.TimestampProto(creationDate)
	logging.Log("ADMIN-yzma4").OnError(err).Debug("unable to parse creation date")
	changeDatePb, err := ptypes.TimestampProto(changeDate)
	logging.Log("ADMIN-NTgjY").OnError(err).Debug("unable to parse change date")

	return &object_pb.ObjectDetails{
		Sequence:      sequence,
		CreationDate:  creationDatePb,
		ChangeDate:    changeDatePb,
		ResourceOwner: resourceOwner,
	}
}

func ToListDetails(
	totalResult,
	processedSequence uint64,
	viewTimestamp time.Time,
) *object.ListDetails {
	viewTs, err := ptypes.TimestampProto(viewTimestamp)
	logging.Log("OBJEC-WoeFH").OnError(err).Debug("")
	return &object_pb.ListDetails{
		TotalResult:       totalResult,
		ProcessedSequence: processedSequence,
		ViewTimestamp:     viewTs,
	}
}

func TextMethodToModel(method object_pb.TextQueryMethod) model.SearchMethod {
	switch method {
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		return model.SearchMethodEquals
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return model.SearchMethodEqualsIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return model.SearchMethodStartsWith
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return model.SearchMethodStartsWithIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS:
		return model.SearchMethodContains
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
		return model.SearchMethodContainsIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
		fallthrough
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		fallthrough
	default:
		return -1
	}
}
