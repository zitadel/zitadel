package view

import (
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	notificationTable = "notification.notifications"
)

func (v *View) GetLatestNotificationSequence() (uint64, error) {
	return v.latestSequence(notificationTable)
}

func (v *View) ProcessedNotificationSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(notificationTable, eventSequence)
}

func (v *View) GetLatestNotificationFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(notificationTable, sequence)
}

func (v *View) ProcessedNotificationFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
