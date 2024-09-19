package query

import (
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func readModelToObjectDetails(model *eventstore.ReadModel) *domain.ObjectDetails {
	return &domain.ObjectDetails{
		Sequence:      model.ProcessedSequence,
		ResourceOwner: model.ResourceOwner,
		EventDate:     model.ChangeDate,
	}
}
