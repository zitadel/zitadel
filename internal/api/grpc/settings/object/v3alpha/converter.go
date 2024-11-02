package object

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	settings_object "github.com/zitadel/zitadel/pkg/grpc/settings/object/v3alpha"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails, owner *object.Owner) *settings_object.Details {
	details := &settings_object.Details{
		Sequence: objectDetail.Sequence,
		Owner:    owner,
	}
	if !objectDetail.EventDate.IsZero() {
		details.ChangeDate = timestamppb.New(objectDetail.EventDate)
	}
	return details
}
