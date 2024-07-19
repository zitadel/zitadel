package object

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails, owner *object.Owner, id string) *resource_object.Details {
	details := &resource_object.Details{
		Id:       id,
		Sequence: objectDetail.Sequence,
		Owner:    owner,
	}
	if !objectDetail.EventDate.IsZero() {
		details.ChangeDate = timestamppb.New(objectDetail.EventDate)
	}
	return details
}

func ToListDetails(request query.SearchRequest, response query.SearchResponse) *resource_object.ListDetails {
	details := &resource_object.ListDetails{
		AppliedLimit:      uint32(request.Limit),
		EndOfList:         false, // TODO: Implement
		TotalResult:       response.Count,
		ProcessedSequence: response.Sequence,
		Timestamp:         timestamppb.New(response.EventCreatedAt),
	}

	return details
}

func TextMethodToQuery(method resource_object.TextFilterMethod) query.TextComparison {
	switch method {
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS:
		return query.TextEquals
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS_IGNORE_CASE:
		return query.TextEqualsIgnoreCase
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH:
		return query.TextStartsWith
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH_IGNORE_CASE:
		return query.TextStartsWithIgnoreCase
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_CONTAINS:
		return query.TextContains
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_CONTAINS_IGNORE_CASE:
		return query.TextContainsIgnoreCase
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_ENDS_WITH:
		return query.TextEndsWith
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_ENDS_WITH_IGNORE_CASE:
		return query.TextEndsWithIgnoreCase
	default:
		return -1
	}
}

func ListQueryToQuery(query *resource_object.ListQuery) (offset, limit uint64, asc bool) {
	if query == nil {
		return 0, 0, false
	}
	return query.Offset, uint64(query.Limit), query.Asc
}
