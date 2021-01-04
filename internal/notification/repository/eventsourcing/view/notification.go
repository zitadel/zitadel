package view

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	notificationTable = "notification.notifications"
)

func (v *View) GetLatestNotificationSequence(aggregateType string) (*repository.CurrentSequence, error) {
	return v.latestSequence(notificationTable, aggregateType)
}

func (v *View) ProcessedNotificationSequence(event *models.Event) error {
	return v.saveCurrentSequence(notificationTable, event)
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
