package view

import (
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	sequencesTable = "admin_api.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, sequence uint64) error {
	return repository.SaveCurrentSequence(v.Db, sequencesTable, viewName, sequence)
}

func (v *View) latestSequence(viewName string) (uint64, error) {
	return repository.LatestSequence(v.Db, sequencesTable, viewName)
}

func (v *View) AllCurrentSequences(db string) ([]*repository.CurrentSequence, error) {
	return repository.AllCurrentSequences(v.Db, db+".current_sequences")
}

func (v *View) ClearView(db, viewName string) error {
	truncateView := db + "." + viewName
	sequenceTable := db + ".current_sequences"
	return repository.ClearView(v.Db, truncateView, sequenceTable)
}
