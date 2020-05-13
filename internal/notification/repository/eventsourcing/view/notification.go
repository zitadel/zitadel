package view

import (
	global_view "github.com/caos/zitadel/internal/view"
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

func (v *View) GetLatestNotificationFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(notificationTable, sequence)
}

func (v *View) ProcessedNotificationFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
