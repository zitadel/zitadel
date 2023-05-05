package object

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
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

func ResourceOwnerFromReq(ctx context.Context, req *object.RequestContext) string {
	if req.GetInstance() {
		return authz.GetInstance(ctx).InstanceID()
	}
	if req.GetOrgId() != "" {
		return req.GetOrgId()
	}
	return authz.GetCtxData(ctx).OrgID
}
