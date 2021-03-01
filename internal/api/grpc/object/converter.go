package object

import (
	"github.com/caos/zitadel/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/pkg/grpc/object"
	object_pb "github.com/caos/zitadel/pkg/grpc/object"
	"github.com/golang/protobuf/ptypes"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails) *object_pb.ObjectDetails {
	return &object_pb.ObjectDetails{
		Sequence:      objectDetail.Sequence,
		ChangeDate:    timestamppb.New(objectDetail.ChangeDate),
		ResourceOwner: objectDetail.ResourceOwner,
	}
}

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

func TextMethodToModel(method object_pb.TextQueryMethod) domain.SearchMethod {
	switch method {
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		return domain.SearchMethodEquals
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return domain.SearchMethodEqualsIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return domain.SearchMethodStartsWith
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return domain.SearchMethodStartsWithIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS:
		return domain.SearchMethodContains
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
		return domain.SearchMethodContainsIgnoreCase
		//TODO: uncomment when added in proto
	//case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
	//	fallthrough
	//case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
	//	fallthrough
	default:
		return -1
	}
}
