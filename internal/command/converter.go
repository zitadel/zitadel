package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func writeModelToObjectDetails(writeModel *eventstore.WriteModel) *domain.ObjectDetails {
	return &domain.ObjectDetails{
		Sequence:      writeModel.ProcessedSequence,
		ResourceOwner: writeModel.ResourceOwner,
		EventDate:     writeModel.ChangeDate,
	}
}
