package object

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resources_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails, owner *object.Owner, id string) *resources_object.Details {
	details := &resources_object.Details{
		Id:       id,
		Sequence: objectDetail.Sequence,
		Owner:    owner,
	}
	if !objectDetail.EventDate.IsZero() {
		details.ChangeDate = timestamppb.New(objectDetail.EventDate)
	}
	return details
}
