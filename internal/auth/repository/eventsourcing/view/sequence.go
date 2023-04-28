package view

import (
	"context"
	"time"

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

func (v *View) latestSequences(ctx context.Context, viewName string, instanceIDs []string) ([]*repository.CurrentSequence, error) {
	return repository.LatestSequences(v.Db, v.TimeTravel(ctx, sequencesTable), viewName, instanceIDs)
}

func (v *View) updateSpoolerRunSequence(viewName string, instanceIDs []string) error {
	currentSequences, err := repository.LatestSequences(v.Db, sequencesTable, viewName, instanceIDs)
	if err != nil {
		return err
	}
	for _, currentSequence := range currentSequences {
		if currentSequence.ViewName == "" {
			currentSequence.ViewName = viewName
		}
		currentSequence.LastSuccessfulSpoolerRun = time.Now()
	}
	return repository.UpdateCurrentSequences(v.Db, sequencesTable, currentSequences)
}
