package view

import (
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	errTable = "auth.failed_events"
)

func (v *View) saveFailedEvent(failedEvent *repository.FailedEvent) error {
	return repository.SaveFailedEvent(v.Db, errTable, failedEvent)
}

func (v *View) latestFailedEvent(viewName, instanceID, eventID string) (*repository.FailedEvent, error) {
	return repository.LatestFailedEvent(v.Db, errTable, viewName, instanceID, eventID)
}
