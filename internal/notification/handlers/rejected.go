package handlers

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// notSent reports whether a notification was not sent, so no sent event must be created for it.
// The returned error is nil if the notification must not be retried,
// because it was canceled or rejected. A rejection is stated on the aggregate of the event.
func notSent(ctx context.Context, commands Commands, event eventstore.Event, err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	err = stateRejection(ctx, commands, event.Aggregate(), event.Type(), err)
	if errors.Is(err, &channels.CancelError{}) {
		return true, nil
	}
	return true, err
}

// stateRejection states a rejected notification on the aggregate that triggered it.
// The notification is then canceled, so it is not retried. Any other error is returned unchanged.
// If the rejection cannot be stated, the error of the command is returned, so the notification is retried,
// unless the aggregate does not exist anymore.
func stateRejection(ctx context.Context, commands Commands, aggregate *eventstore.Aggregate, triggeringEventType eventstore.EventType, err error) error {
	rejected := new(channels.RejectedError)
	if !errors.As(err, &rejected) {
		return err
	}
	rejection := rejected.Rejection
	rejection.TriggeringEventType = triggeringEventType

	commandErr := commands.NotificationRejected(ctx, aggregate, rejection)
	if zerrors.IsPreconditionFailed(commandErr) {
		// the aggregate was removed in the meantime, there is nothing to state the rejection on
		logging.Debug(ctx, "rejection of notification not stated", "err", commandErr, "eventType", triggeringEventType, "reason", rejection.Reason)
		return channels.NewCancelError(err)
	}
	if commandErr != nil {
		return commandErr
	}
	// the counter states how many notifications are rejected, the event which ones
	logging.Debug(ctx, "notification rejected", "eventType", triggeringEventType, "reason", rejection.Reason, "provider", rejection.ProviderID)
	countErr := metrics.AddCount(ctx, types.RejectedNotificationsCounter, 1, map[string]attribute.Value{
		"reason":                attribute.StringValue(string(rejection.Reason)),
		"triggering_event_type": attribute.StringValue(string(triggeringEventType)),
	})
	logging.OnError(ctx, countErr).Warn("incrementing counter metric failed", "name", types.RejectedNotificationsCounter)
	return channels.NewCancelError(err)
}
