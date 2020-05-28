package view

import (
	"github.com/caos/zitadel/internal/view"
)

const (
	sequencesTable = "admin_api.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, sequence uint64) error {
	return view.SaveCurrentSequence(v.Db, sequencesTable, viewName, sequence)
}

func (v *View) latestSequence(viewName string) (uint64, error) {
	return view.LatestSequence(v.Db, sequencesTable, viewName)
}
