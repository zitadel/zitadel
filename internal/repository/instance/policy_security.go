package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	securityPolicyPrefix       = "policy.security."
	SecurityPolicySetEventType = instanceEventTypePrefix + securityPolicyPrefix + "set"

	securityPolicyClientIDMetadataDocumentAllowedURLPrefix           = securityPolicyPrefix + "client_id_metadata_document.allowed_url."
	SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEventType   = instanceEventTypePrefix + securityPolicyClientIDMetadataDocumentAllowedURLPrefix + "added"
	SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEventType = instanceEventTypePrefix + securityPolicyClientIDMetadataDocumentAllowedURLPrefix + "removed"
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

	// EnableClientIDMetadataDocument resolves a client_id that is an absolute HTTPS URL as a
	// Client ID Metadata Document instead of looking it up in the database.
	EnableClientIDMetadataDocument *bool `json:"enable_client_id_metadata_document,omitempty"`
	// ClientIDMetadataDocumentAllowedURLs replaces the client_id URLs the instance resolves as
	// Client ID Metadata Documents.
	ClientIDMetadataDocumentAllowedURLs *[]string `json:"client_id_metadata_document_allowed_urls,omitempty"`
	// ClientIDMetadataDocumentAllowAnyURL lets the instance resolve every client_id URL the
	// system allows, regardless of ClientIDMetadataDocumentAllowedURLs.
	ClientIDMetadataDocumentAllowAnyURL *bool `json:"client_id_metadata_document_allow_any_url,omitempty"`
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

func ChangeSecurityPolicyEnableClientIDMetadataDocument(enabled bool) func(event *SecurityPolicySetEvent) {
	return func(e *SecurityPolicySetEvent) {
		e.EnableClientIDMetadataDocument = &enabled
	}
}

func ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs(allowedURLs []string) func(event *SecurityPolicySetEvent) {
	return func(e *SecurityPolicySetEvent) {
		if len(allowedURLs) == 0 {
			allowedURLs = []string{}
		}
		e.ClientIDMetadataDocumentAllowedURLs = &allowedURLs
	}
}

func ChangeSecurityPolicyClientIDMetadataDocumentAllowAnyURL(allow bool) func(event *SecurityPolicySetEvent) {
	return func(e *SecurityPolicySetEvent) {
		e.ClientIDMetadataDocumentAllowAnyURL = &allow
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

// SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent adds a client_id URL the instance
// resolves as a Client ID Metadata Document. It is a separate event rather than a set of the
// whole list, so concurrent additions do not overwrite each other.
type SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	URL string `json:"url"`
}

func (e *SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	url string,
) *SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent {
	return &SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEventType,
		),
		URL: url,
	}
}

func (e *SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent) Payload() interface{} {
	return e
}

func (e *SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent removes a client_id URL the
// instance resolves as a Client ID Metadata Document.
type SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	URL string `json:"url"`
}

func (e *SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewSecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	url string,
) *SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent {
	return &SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEventType,
		),
		URL: url,
	}
}

func (e *SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent) Payload() interface{} {
	return e
}

func (e *SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
