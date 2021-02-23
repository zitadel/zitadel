package view

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	sequencesTable = "adminapi.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, event *models.Event) error {
	return repository.SaveCurrentSequence(v.Db, sequencesTable, viewName, event.Sequence, event.CreationDate)
}

func (v *View) latestSequence(viewName string) (*repository.CurrentSequence, error) {
	return repository.LatestSequence(v.Db, sequencesTable, viewName)
}

func (v *View) AllCurrentSequences(db string) ([]*repository.CurrentSequence, error) {
	return repository.AllCurrentSequences(v.Db, db+".current_sequences")
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

func (v *View) GetCurrentSequence(db, viewName string) (*repository.CurrentSequence, error) {
	sequenceTable := db + ".current_sequences"
	fullView := db + "." + viewName
	return repository.LatestSequence(v.Db, sequenceTable, fullView)
}

func (v *View) ClearView(db, viewName string) error {
	truncateView := db + "." + viewName
	sequenceTable := db + ".current_sequences"
	return repository.ClearView(v.Db, truncateView, sequenceTable)
}
