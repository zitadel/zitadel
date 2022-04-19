package view

import (
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	errTable  = "adminapi.failed_events"
	errColumn = "failed_events"
)

func (v *View) saveFailedEvent(failedEvent *repository.FailedEvent) error {
	return repository.SaveFailedEvent(v.Db, errTable, failedEvent)
}

func (v *View) RemoveFailedEvent(database string, failedEvent *repository.FailedEvent) error {
	return repository.RemoveFailedEvent(v.Db, database+"."+errColumn, failedEvent)
}

func (v *View) latestFailedEvent(viewName, instanceID string, sequence uint64) (*repository.FailedEvent, error) {
	return repository.LatestFailedEvent(v.Db, errTable, viewName, instanceID, sequence)
}

func (v *View) AllFailedEvents(db string) ([]*repository.FailedEvent, error) {
	return repository.AllFailedEvents(v.Db, db+"."+errColumn)
}
