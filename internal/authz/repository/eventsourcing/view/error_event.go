package view

import (
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	errTable = "authz.failed_event"
)

func (v *View) saveFailedEvent(failedEvent *repository.FailedEvent) error {
	return repository.SaveFailedEvent(v.Db, errTable, failedEvent)
}

func (v *View) latestFailedEvent(viewName string, sequence uint64) (*repository.FailedEvent, error) {
	return repository.LatestFailedEvent(v.Db, errTable, viewName, sequence)
}
