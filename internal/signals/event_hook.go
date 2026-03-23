package signals

// PREVIEW: Identity Signals is a preview feature. APIs, storage format,
// and configuration may change between releases without notice.

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

// NewEventSignalHook returns a hook function that converts every pushed
// event into a Signal on the "events" stream. The conversion runs in a
// background goroutine so the eventstore push path is never delayed.
// Events are emitted fire-and-forget through the given Emitter.
func NewEventSignalHook(emitter *Emitter) func(ctx context.Context, events []eventstore.Event) {
	return func(ctx context.Context, events []eventstore.Event) {
		traceID := tracing.TraceIDFromCtx(ctx)
		spanID := spanIDFromCtx(ctx)

		// Snapshot event data before spawning goroutine — Event
		// interface values may not be safe to read after Push returns.
		type snap struct {
			instanceID    string
			aggID         string
			aggType       string
			resourceOwner string
			creator       string
			eventType     string
			ts            time.Time
			payload       string
		}
		snaps := make([]snap, len(events))
		for i, e := range events {
			agg := e.Aggregate()
			ts := e.CreatedAt()
			if ts.IsZero() {
				ts = time.Now().UTC()
			}
			var payload string
			if b := e.DataAsBytes(); len(b) > 0 {
				payload = string(b)
			}
			snaps[i] = snap{
				instanceID:    agg.InstanceID,
				aggID:         agg.ID,
				aggType:       string(agg.Type),
				resourceOwner: agg.ResourceOwner,
				creator:       e.Creator(),
				eventType:     string(e.Type()),
				ts:            ts,
				payload:        payload,
			}
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.ErrorContext(ctx, "identity_signals.event_hook_panic",
						slog.Any("panic", r),
						slog.Int("batch_size", len(snaps)),
					)
				}
			}()
			for _, s := range snaps {
				ids := extractIDs(s.aggType, s.aggID, s.payload)

				// Fallback: use the event creator (acting user) when no
				// userID could be extracted from aggregate or payload.
				// Skip system-level identifiers like "SYSTEM" / "system"
				// which are set when the OIDC/SAML server acts on behalf
				// of a user (not a real user ID).
				userID := ids.userID
				if userID == "" && isRealUser(s.creator) {
					userID = s.creator
				}

				emitter.Emit(Signal{
					InstanceID: s.instanceID,
					UserID:     userID,
					CallerID:   s.creator,
					SessionID:  ids.sessionID,
					ClientID:   ids.clientID,
					OrgID:      s.resourceOwner,
					Operation:  s.eventType,
					Stream:     StreamEvents,
					Resource:   s.aggType + "/" + s.aggID,
					Outcome:    outcomeFromEventType(s.eventType),
					Timestamp:  s.ts,
					Payload:    s.payload,
					TraceID:    traceID,
					SpanID:     spanID,
				})
			}
		}()
	}
}

// outcomeFromEventType derives the outcome from an event type name.
// Events ending in ".failed" are classified as failures; all others
// are treated as successes.
func outcomeFromEventType(eventType string) Outcome {
	if strings.HasSuffix(eventType, ".failed") {
		return OutcomeFailure
	}
	return OutcomeSuccess
}

// extractedIDs holds the identity fields extracted from an event.
type extractedIDs struct {
	userID    string
	sessionID string
	clientID  string
}

// extractIDs pulls user/session/client IDs from the aggregate type + ID and
// the JSON payload. Different event types use inconsistent JSON field names
// (camelCase vs snake_case), so we parse into a generic map and check all
// known variants.
func extractIDs(aggType, aggID, payload string) extractedIDs {
	var ids extractedIDs

	// Aggregate-level: the aggregate ID is the entity itself for user/session types
	switch {
	case strings.HasPrefix(aggType, "user"):
		ids.userID = aggID
	case strings.HasPrefix(aggType, "session"):
		ids.sessionID = aggID
	}

	// Parse the JSON payload once into a generic map to handle all field
	// name variants used across the codebase.
	// Cap at 1MB to prevent OOM on extremely large payloads.
	const maxPayloadSize = 1 << 20 // 1MB
	if payload == "" || len(payload) > maxPayloadSize {
		return ids
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(payload), &m); err != nil {
		return ids
	}

	// User ID — try all variants: "userID", "user_id", "userId", "hint_user_id"
	if ids.userID == "" {
		ids.userID = firstStringField(m, "userID", "user_id", "userId", "hint_user_id")
	}

	// Session ID — try: "sessionID", "session_id"
	if ids.sessionID == "" {
		ids.sessionID = firstStringField(m, "sessionID", "session_id")
	}

	// Client ID — try: "clientID", "clientId", "client_id"
	ids.clientID = firstStringField(m, "clientID", "clientId", "client_id")

	return ids
}

// firstStringField returns the value of the first key found in the map that
// unmarshals to a non-empty string.
func firstStringField(m map[string]json.RawMessage, keys ...string) string {
	for _, k := range keys {
		raw, ok := m[k]
		if !ok {
			continue
		}
		var v string
		if err := json.Unmarshal(raw, &v); err == nil && v != "" {
			return v
		}
	}
	return ""
}

// isRealUser returns true when creator looks like a real user ID rather
// than a well-known system identifier. The OIDC/SAML server and migration
// code use "SYSTEM" / "system" as creator; these are not real users.
func isRealUser(creator string) bool {
	if creator == "" {
		return false
	}
	up := strings.ToUpper(creator)
	return up != "SYSTEM"
}
