package view

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	sequencesTable = "auth.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, event *models.Event) error {
	return repository.SaveCurrentSequence(v.Db, sequencesTable, viewName, event.Sequence, event.CreationDate)
}

func (v *View) latestSequence(viewName string) (*repository.CurrentSequence, error) {
	return repository.LatestSequence(v.Db, sequencesTable, viewName)
}

func (v *View) updateSpoolerRunSequence(viewName string) error {
	currentSequence, err := repository.LatestSequence(v.Db, sequencesTable, viewName)
	if err != nil {
		return err
	}
	if currentSequence.ViewName == "" {
		currentSequence.ViewName = viewName
	}
	currentSequence.LastSuccessfulSpoolerRun = time.Now()
	//update all aggregate types
	//TODO: not sure if all scenarios work as expected
	currentSequence.AggregateType = ""
	return repository.UpdateCurrentSequence(v.Db, sequencesTable, currentSequence)
}
