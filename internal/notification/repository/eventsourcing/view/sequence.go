package view

import (
	"github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	sequencesTable = "notification.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, sequence uint64) error {
	return repository.SaveCurrentSequence(v.Db, sequencesTable, viewName, sequence)
}

func (v *View) latestSequence(viewName string) (uint64, time.Time, error) {
	return repository.LatestSequence(v.Db, sequencesTable, viewName)
}
