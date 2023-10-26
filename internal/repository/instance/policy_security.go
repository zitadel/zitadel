package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	securityPolicyPrefix       = "policy.security."
	SecurityPolicySetEventType = instanceEventTypePrefix + securityPolicyPrefix + "set"
)

type SecurityPolicySetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Enabled        *bool     `json:"enabled,omitempty"`
	AllowedOrigins *[]string `json:"allowedOrigins,omitempty"`
}

func NewSecurityPolicySetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []SecurityPolicyChanges,
) (*SecurityPolicySetEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-EWsf3", "Errors.NoChangesFound")
	}
	event := &SecurityPolicySetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SecurityPolicySetEventType,
		),
	}
	for _, change := range changes {
		change(event)
	}
	return event, nil
}

type SecurityPolicyChanges func(event *SecurityPolicySetEvent)

func ChangeSecurityPolicyEnabled(enabled bool) func(event *SecurityPolicySetEvent) {
	return func(e *SecurityPolicySetEvent) {
		e.Enabled = &enabled
	}
}

func ChangeSecurityPolicyAllowedOrigins(allowedOrigins []string) func(event *SecurityPolicySetEvent) {
	return func(e *SecurityPolicySetEvent) {
		if len(allowedOrigins) == 0 {
			allowedOrigins = []string{}
		}
		e.AllowedOrigins = &allowedOrigins
	}
}

func (e *SecurityPolicySetEvent) Payload() interface{} {
	return e
}

func (e *SecurityPolicySetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SecurityPolicySetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	securityPolicyAdded := &SecurityPolicySetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(securityPolicyAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-soiwj", "unable to unmarshal oidc config added")
	}

	return securityPolicyAdded, nil
}
