// Package signalprojection provides an event-store handler that converts
// security-relevant domain events into signals for the signal explorer.
package signalprojection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/signals"
)

const ProjectionTable = "projections.event_signals"

// signalEmitter is an interface matching [signals.Emitter.Emit] to avoid
// coupling tightly to the concrete type.
type signalEmitter interface {
	Emit(signal signals.Signal)
}

// eventSignalHandler is an event-store projection that converts
// security-relevant domain events into signals emitted to the
// fire-and-forget [signals.Emitter]. This gives the signal explorer
// complete visibility into authentication and account lifecycle events,
// not just HTTP requests captured by interceptors.
type eventSignalHandler struct {
	emitter signalEmitter
}

// NewHandler creates a handler.Handler that subscribes to session, user,
// and oidc_session aggregate events and emits them as signals. The handler
// must be started with handler.Handler.Start(ctx).
func NewHandler(
	ctx context.Context,
	config handler.Config,
	emitter signalEmitter,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &eventSignalHandler{
		emitter: emitter,
	})
}

func (h *eventSignalHandler) Name() string {
	return ProjectionTable
}

func (h *eventSignalHandler) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: session.AggregateType,
			EventReducers: []handler.EventReducer{
				{Event: session.AddedType, Reduce: h.reduceSessionEvent},
				{Event: session.UserCheckedType, Reduce: h.reduceSessionEvent},
				{Event: session.PasswordCheckedType, Reduce: h.reduceSessionEvent},
				{Event: session.IntentCheckedType, Reduce: h.reduceSessionEvent},
				{Event: session.WebAuthNChallengedType, Reduce: h.reduceSessionEvent},
				{Event: session.WebAuthNCheckedType, Reduce: h.reduceSessionEvent},
				{Event: session.TOTPCheckedType, Reduce: h.reduceSessionEvent},
				{Event: session.OTPSMSChallengedType, Reduce: h.reduceSessionEvent},
				{Event: session.OTPSMSCheckedType, Reduce: h.reduceSessionEvent},
				{Event: session.OTPEmailChallengedType, Reduce: h.reduceSessionEvent},
				{Event: session.OTPEmailCheckedType, Reduce: h.reduceSessionEvent},
				{Event: session.RecoveryCodeCheckedType, Reduce: h.reduceSessionEvent},
				{Event: session.TokenSetType, Reduce: h.reduceSessionEvent},
				{Event: session.TerminateType, Reduce: h.reduceSessionEvent},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{Event: user.HumanRegisteredType, Reduce: h.reduceUserEvent},
				{Event: user.HumanAddedType, Reduce: h.reduceUserEvent},
				{Event: user.HumanPasswordChangedType, Reduce: h.reduceUserEvent},
				{Event: user.HumanPasswordCheckSucceededType, Reduce: h.reduceUserEvent},
				{Event: user.HumanPasswordCheckFailedType, Reduce: h.reduceUserFailure},
				{Event: user.HumanInitializedCheckSucceededType, Reduce: h.reduceUserEvent},
				{Event: user.HumanInitializedCheckFailedType, Reduce: h.reduceUserFailure},
				{Event: user.HumanSignedOutType, Reduce: h.reduceUserEvent},
				{Event: user.UserLockedType, Reduce: h.reduceUserFailure},
				{Event: user.UserDeactivatedType, Reduce: h.reduceUserEvent},
				{Event: user.UserReactivatedType, Reduce: h.reduceUserEvent},
				{Event: user.UserRemovedType, Reduce: h.reduceUserEvent},
				{Event: user.UserIDPLinkAddedType, Reduce: h.reduceUserEvent},
				{Event: user.UserIDPLinkRemovedType, Reduce: h.reduceUserEvent},
				{Event: user.UserImpersonatedType, Reduce: h.reduceUserEvent},
			},
		},
		{
			Aggregate: oidcsession.AggregateType,
			EventReducers: []handler.EventReducer{
				{Event: oidcsession.AddedType, Reduce: h.reduceOIDCSessionEvent},
				{Event: oidcsession.AccessTokenAddedType, Reduce: h.reduceOIDCSessionEvent},
				{Event: oidcsession.AccessTokenRevokedType, Reduce: h.reduceOIDCSessionEvent},
				{Event: oidcsession.RefreshTokenAddedType, Reduce: h.reduceOIDCSessionEvent},
				{Event: oidcsession.RefreshTokenRenewedType, Reduce: h.reduceOIDCSessionEvent},
				{Event: oidcsession.RefreshTokenRevokedType, Reduce: h.reduceOIDCSessionEvent},
			},
		},
	}
}

// reduceSessionEvent emits a signal for session aggregate events.
func (h *eventSignalHandler) reduceSessionEvent(event eventstore.Event) (*handler.Statement, error) {
	h.emitEvent(event, signals.StreamEvents, signals.OutcomeSuccess)
	return handler.NewNoOpStatement(event), nil
}

// reduceUserEvent emits a signal for user aggregate events.
func (h *eventSignalHandler) reduceUserEvent(event eventstore.Event) (*handler.Statement, error) {
	h.emitEvent(event, signals.StreamEvents, signals.OutcomeSuccess)
	return handler.NewNoOpStatement(event), nil
}

// reduceUserFailure emits a signal for user failure events (lock, failed checks).
func (h *eventSignalHandler) reduceUserFailure(event eventstore.Event) (*handler.Statement, error) {
	h.emitEvent(event, signals.StreamEvents, signals.OutcomeFailure)
	return handler.NewNoOpStatement(event), nil
}

// reduceOIDCSessionEvent emits a signal for OIDC session events.
func (h *eventSignalHandler) reduceOIDCSessionEvent(event eventstore.Event) (*handler.Statement, error) {
	h.emitEvent(event, signals.StreamEvents, signals.OutcomeSuccess)
	return handler.NewNoOpStatement(event), nil
}

// emitEvent converts an eventstore.Event into a Signal and emits it.
func (h *eventSignalHandler) emitEvent(event eventstore.Event, stream signals.SignalStream, outcome signals.Outcome) {
	agg := event.Aggregate()

	ts := event.CreatedAt()
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	h.emitter.Emit(signals.Signal{
		InstanceID: agg.InstanceID,
		UserID:     agg.ID,
		CallerID:   event.Creator(),
		SessionID:  agg.ID,
		Operation:  string(event.Type()),
		Stream:     stream,
		Resource:   string(agg.Type),
		Outcome:    outcome,
		Timestamp:  ts,
	})
}
