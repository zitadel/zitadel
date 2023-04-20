package object

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"

	"github.com/zitadel/zitadel/internal/domain"
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
