package view

import (
	"github.com/caos/zitadel/internal/view"
)

const (
	errTable = "authz.failed_event"
)

func (v *View) saveFailedEvent(failedEvent *view.FailedEvent) error {
	return view.SaveFailedEvent(v.Db, errTable, failedEvent)
}

func (v *View) latestFailedEvent(viewName string, sequence uint64) (*view.FailedEvent, error) {
	return view.LatestFailedEvent(v.Db, errTable, viewName, sequence)
}
