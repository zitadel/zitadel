package object

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails) *object.Details {
	details := &object.Details{
		Sequence:      objectDetail.Sequence,
		ResourceOwner: objectDetail.ResourceOwner,
	}
	if !objectDetail.EventDate.IsZero() {
		details.ChangeDate = timestamppb.New(objectDetail.EventDate)
	}
	return details
}

func ToListDetails(response query.SearchResponse) *object.ListDetails {
	details := &object.ListDetails{
		TotalResult:       response.Count,
		ProcessedSequence: response.Sequence,
		Timestamp:         timestamppb.New(response.EventCreatedAt),
	}

	return details
}
func ListQueryToQuery(query *object.ListQuery) (offset, limit uint64, asc bool) {
	if query == nil {
		return 0, 0, false
	}
	return query.Offset, uint64(query.Limit), query.Asc
}

func ResourceOwnerFromReq(ctx context.Context, req *object.RequestContext) string {
	if req.GetInstance() {
		return authz.GetInstance(ctx).InstanceID()
	}
	if req.GetOrgId() != "" {
		return req.GetOrgId()
	}
	return authz.GetCtxData(ctx).OrgID
}

func TextMethodToQuery(method object.TextQueryMethod) query.TextComparison {
	switch method {
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		return query.TextEquals
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return query.TextEqualsIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return query.TextStartsWith
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return query.TextStartsWithIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS:
		return query.TextContains
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
		return query.TextContainsIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
		return query.TextEndsWith
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		return query.TextEndsWithIgnoreCase
	default:
		return -1
	}
}
