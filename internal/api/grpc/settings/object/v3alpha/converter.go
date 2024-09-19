package object

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/v2/internal/domain"
	object "github.com/zitadel/zitadel/v2/pkg/grpc/object/v3alpha"
	settings_object "github.com/zitadel/zitadel/v2/pkg/grpc/settings/object/v3alpha"
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
