package channels

import (
	"errors"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type CancelError struct {
	Err error
}

func (e *CancelError) Error() string {
	return e.Err.Error()
}

func NewCancelError(err error) error {
	return &CancelError{
		Err: err,
	}
}

func (e *CancelError) Is(target error) bool {
	return errors.As(target, &e)
}

func (e *CancelError) Unwrap() error {
	return e.Err
}

// RejectionReason states why a notification was rejected.
type RejectionReason string

const (
	// RejectionReasonLimitExceeded states that the limit of the provider was exceeded.
	RejectionReasonLimitExceeded RejectionReason = "limit_exceeded"
)

// Rejection describes a notification that was intentionally not sent.
// It is the payload of the rejected events of the notification.
type Rejection struct {
	// TriggeringEventType is the type of the event that requested the notification.
	TriggeringEventType eventstore.EventType `json:"triggeringEventType"`
	Reason              RejectionReason      `json:"reason"`
	// ProviderID is the ID of the provider that would have sent the notification.
	ProviderID string `json:"providerId,omitzero"`
}

// RejectionsOf returns the filter for the payload of rejected events
// of notifications requested by events of the given type.
func RejectionsOf(triggeringEventType eventstore.EventType) map[string]any {
	return map[string]any{"triggeringEventType": triggeringEventType}
}

// RejectedError is returned if a notification is intentionally not sent,
// e.g. because a limit of the provider is exceeded.
// Like a [CancelError] the notification is not retried,
// but the rejection is stated as event on the aggregate, so it is visible why nothing was sent.
type RejectedError struct {
	// Rejection is stated on the aggregate.
	// The triggering event type is set by the handler of the notification.
	Rejection Rejection
}

func NewRejectedError(reason RejectionReason, providerID string) error {
	return &RejectedError{
		Rejection: Rejection{
			Reason:     reason,
			ProviderID: providerID,
		},
	}
}

func (e *RejectedError) Error() string {
	return fmt.Sprintf("notification rejected: %s", e.Rejection.Reason)
}
