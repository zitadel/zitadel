package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	sequencesTable = "auth.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, event *models.Event) error {
	return repository.SaveCurrentSequence(v.Db, sequencesTable, viewName, event.InstanceID, event.Sequence, event.CreationDate)
}

func (v *View) latestSequence(ctx context.Context, viewName, instanceID string) (*repository.CurrentSequence, error) {
	return repository.LatestSequence(v.Db, v.TimeTravel(ctx, sequencesTable), viewName, instanceID)
}
