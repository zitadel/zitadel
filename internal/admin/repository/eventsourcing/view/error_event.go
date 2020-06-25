package view

import (
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	errTable = "admin_api.failed_event"
)

func (v *View) saveFailedEvent(failedEvent *repository.FailedEvent) error {
	return repository.SaveFailedEvent(v.Db, errTable, failedEvent)
}

func (v *View) RemoveFailedEvent(database string, failedEvent *repository.FailedEvent) error {
	return repository.RemoveFailedEvent(v.Db, database+".failed_event", failedEvent)
}

func (v *View) latestFailedEvent(viewName string, sequence uint64) (*repository.FailedEvent, error) {
	return repository.LatestFailedEvent(v.Db, errTable, viewName, sequence)
}

func (v *View) AllFailedEvents(db string) ([]*repository.FailedEvent, error) {
	return repository.AllFailedEvents(v.Db, db+".failed_event")
}
