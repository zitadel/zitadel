package session

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// TestTerminateEventMapperUnmarshalsTriggerOrigin pins issue #12948: the
// session.terminated event carries the origin that terminated the session
// in its payload, and the notification handler reads it through
// TerminateEvent.TriggerOrigin() to issue the back-channel logout token
// from the terminating origin. A mapper that skips the payload leaves
// TriggerOrigin() empty, so the logout token's iss falls back to the
// instance's primary domain and relying parties that validate iss like an
// ID token's reject every logout token.
func TestTerminateEventMapperUnmarshalsTriggerOrigin(t *testing.T) {
	event := &eventstore.BaseEvent{
		Agg: &eventstore.Aggregate{
			ID:            "session-1",
			Type:          AggregateType,
			ResourceOwner: "org-1",
			InstanceID:    "instance-1",
			Version:       "v1",
		},
		EventType: TerminateType,
		Creation:  time.Unix(1760000000, 0).UTC(),
		Data:      []byte(`{"triggerOrigin":"https://login.example.com"}`),
	}

	got, err := TerminateEventMapper(event)
	if err != nil {
		t.Fatalf("TerminateEventMapper() error: %v", err)
	}
	terminated, ok := got.(*TerminateEvent)
	if !ok {
		t.Fatalf("TerminateEventMapper() returned %T, want *TerminateEvent", got)
	}
	if want := "https://login.example.com"; terminated.TriggerOrigin() != want {
		t.Errorf("TriggerOrigin() = %q, want %q", terminated.TriggerOrigin(), want)
	}
	if terminated.Aggregate().ID != "session-1" {
		t.Errorf("Aggregate().ID = %q, want the event's aggregate id preserved", terminated.Aggregate().ID)
	}
	if terminated.Type() != TerminateType {
		t.Errorf("Type() = %q, want %q", terminated.Type(), TerminateType)
	}
}
