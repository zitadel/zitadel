package object

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/pkg/grpc/object"
	object_pb "github.com/caos/zitadel/pkg/grpc/object"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails) *object_pb.ObjectDetails {
	return &object_pb.ObjectDetails{
		Sequence:      objectDetail.Sequence,
		ChangeDate:    timestamppb.New(objectDetail.ChangeDate),
		ResourceOwner: objectDetail.ResourceOwner,
	}
}

func ToViewDetailsPb(
	sequence uint64,
	creationDate,
	changeDate time.Time,
	resourceOwner string,
) *object_pb.ObjectDetails {
	return &object_pb.ObjectDetails{
		Sequence:      sequence,
		CreationDate:  timestamppb.New(creationDate),
		ChangeDate:    timestamppb.New(changeDate),
		ResourceOwner: resourceOwner,
	}
}

func ToDetailsPb(
	sequence uint64,
	changeDate time.Time,
	resourceOwner string,
) *object_pb.ObjectDetails {
	return &object_pb.ObjectDetails{
		Sequence:      sequence,
		ChangeDate:    timestamppb.New(changeDate),
		ResourceOwner: resourceOwner,
	}
}

func ToListDetails(
	totalResult,
	processedSequence uint64,
	viewTimestamp time.Time,
) *object.ListDetails {
	return &object_pb.ListDetails{
		TotalResult:       totalResult,
		ProcessedSequence: processedSequence,
		ViewTimestamp:     timestamppb.New(viewTimestamp),
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
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
		return domain.SearchMethodEndsWith
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		return domain.SearchMethodEndsWithIgnoreCase
	default:
		return -1
	}
}

func ListQueryToModel(query *object_pb.ListQuery) (offset, limit uint64, asc bool) {
	if query == nil {
		return
	}
	return query.Offset, uint64(query.Limit), query.Asc
}
