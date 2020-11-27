package view

import (
	"github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	sequencesTable = "authz.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, sequence uint64, eventTimestamp time.Time) error {
	return repository.SaveCurrentSequence(v.Db, sequencesTable, viewName, sequence, eventTimestamp)
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
	return repository.UpdateCurrentSequence(v.Db, sequencesTable, currentSequence)
}
