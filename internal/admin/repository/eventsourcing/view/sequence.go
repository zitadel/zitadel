package view

import (
	"github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	sequencesTable = "adminapi.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, sequence uint64, eventTimeStamp time.Time) error {
	return repository.SaveCurrentSequence(v.Db, sequencesTable, viewName, sequence, eventTimeStamp)
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
	currentSequence.LastSuccessfulSpoolerRun = time.Now()
	return repository.UpdateCurrentSequence(v.Db, sequencesTable, currentSequence)
}

func (v *View) ClearView(db, viewName string) error {
	truncateView := db + "." + viewName
	sequenceTable := db + ".current_sequences"
	return repository.ClearView(v.Db, truncateView, sequenceTable)
}
