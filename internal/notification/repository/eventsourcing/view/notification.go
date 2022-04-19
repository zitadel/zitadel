package view

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	notificationTable = "notification.notifications"
)

func (v *View) GetLatestNotificationSequence(instanceID string) (*repository.CurrentSequence, error) {
	return v.latestSequence(notificationTable, instanceID)
}

func (v *View) GetLatestNotificationSequences() ([]*repository.CurrentSequence, error) {
	return v.latestSequences(notificationTable)
}

func (v *View) ProcessedNotificationSequence(event *models.Event) error {
	return v.saveCurrentSequence(notificationTable, event)
}

func (v *View) UpdateNotificationSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(notificationTable)
}

func (v *View) GetLatestNotificationFailedEvent(sequence uint64, instanceID string) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(notificationTable, instanceID, sequence)
}

func (v *View) ProcessedNotificationFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
