package object

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
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
