package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	securityPolicyPrefix       = "policy.security."
	SecurityPolicySetEventType = instanceEventTypePrefix + securityPolicyPrefix + "set"
)

type SecurityPolicySetEvent struct {
	eventstore.BaseEvent `json:"-"`

	// Enabled is a legacy field which was used before for Iframe Embedding.
	// It is kept so older events can still be reduced.
	Enabled               *bool     `json:"enabled,omitempty"`
	EnableIframeEmbedding *bool     `json:"enable_iframe_embedding,omitempty"`
	AllowedOrigins        *[]string `json:"allowedOrigins,omitempty"`
	EnableImpersonation   *bool     `json:"enable_impersonation,omitempty"`

	// EnableDynamicClientRegistration serves and advertises the OAuth 2.0 Dynamic Client
	// Registration endpoint (RFC 7591).
	EnableDynamicClientRegistration *bool `json:"enable_dynamic_client_registration,omitempty"`
	// AllowUnauthenticatedDynamicClientRegistration additionally allows registration
	// without an access token. It only has an effect if EnableDynamicClientRegistration.
	AllowUnauthenticatedDynamicClientRegistration *bool `json:"allow_unauthenticated_dynamic_client_registration,omitempty"`
}

func NewSecurityPolicySetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []SecurityPolicyChanges,
) (*SecurityPolicySetEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "POLICY-EWsf3", "Errors.NoChangesFound")
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

func ChangeSecurityPolicyEnableIframeEmbedding(enabled bool) func(event *SecurityPolicySetEvent) {
	return func(e *SecurityPolicySetEvent) {
		e.EnableIframeEmbedding = &enabled
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

func ChangeSecurityPolicyEnableImpersonation(enabled bool) func(event *SecurityPolicySetEvent) {
	return func(e *SecurityPolicySetEvent) {
		e.EnableImpersonation = &enabled
	}
}

func ChangeSecurityPolicyEnableDynamicClientRegistration(enabled bool) func(event *SecurityPolicySetEvent) {
	return func(e *SecurityPolicySetEvent) {
		e.EnableDynamicClientRegistration = &enabled
	}
}

func ChangeSecurityPolicyAllowUnauthenticatedDynamicClientRegistration(allow bool) func(event *SecurityPolicySetEvent) {
	return func(e *SecurityPolicySetEvent) {
		e.AllowUnauthenticatedDynamicClientRegistration = &allow
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
		return nil, zerrors.ThrowInternal(err, "INST-soiwj", "unable to unmarshal oidc config added")
	}

	return securityPolicyAdded, nil
}
