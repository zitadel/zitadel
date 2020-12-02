package view

import (
	"github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	notificationTable = "notification.notifications"
)

func (v *View) GetLatestNotificationSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(notificationTable)
}

func (v *View) ProcessedNotificationSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(notificationTable, eventSequence, eventTimestamp)
}

func (v *View) UpdateNotificationSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(notificationTable)
}

func (v *View) GetLatestNotificationFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(notificationTable, sequence)
}

func (v *View) ProcessedNotificationFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
