package view

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	sequencesTable = "notification.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, event *models.Event) error {
	return repository.SaveCurrentSequence(v.Db, sequencesTable, viewName, event.InstanceID, event.Sequence, event.CreationDate)
}

func (v *View) latestSequence(viewName, instanceID string) (*repository.CurrentSequence, error) {
	return repository.LatestSequence(v.Db, sequencesTable, viewName, instanceID)
}

func (v *View) latestSequences(viewName string) ([]*repository.CurrentSequence, error) {
	return repository.LatestSequences(v.Db, sequencesTable, viewName)
}

func (v *View) updateSpoolerRunSequence(viewName string) error {
	currentSequences, err := repository.LatestSequences(v.Db, sequencesTable, viewName)
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
