package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

func writeModelToObjectDetails(writeModel *eventstore.WriteModel) *domain.ObjectDetails {
	return &domain.ObjectDetails{
		Sequence:      writeModel.ProcessedSequence,
		ResourceOwner: writeModel.ResourceOwner,
		EventDate:     writeModel.ChangeDate,
	}
}
