package object

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails, ownerType object.OwnerType, ownerId string) *resource_object.Details {
	details := &resource_object.Details{
		Id: objectDetail.ID,
		Owner: &object.Owner{
			Type: ownerType,
			Id:   ownerId,
		},
	}
	if !objectDetail.EventDate.IsZero() {
		details.Changed = timestamppb.New(objectDetail.EventDate)
	}
	if !objectDetail.CreationDate.IsZero() {
		details.Created = timestamppb.New(objectDetail.CreationDate)
	}
	return details
}
