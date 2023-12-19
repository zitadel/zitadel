package eventstore

import (
	"context"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type TokenRepo struct {
	Eventstore *eventstore.Eventstore
	View       *view.View
}

func (repo *TokenRepo) TokenByIDs(ctx context.Context, userID, tokenID string) (*usr_model.TokenView, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()

	// always load the latest sequence first, so in case the token was not found by id,
	// the sequence will be equal or lower than the actual projection and no events are lost
	sequence, err := repo.View.GetLatestTokenSequence(ctx, instanceID)
	logging.WithFields("instanceID", instanceID, "userID", userID, "tokenID", tokenID).
		OnError(err).
		Errorf("could not get current sequence for TokenByIDs")

	token, viewErr := repo.View.TokenByIDs(tokenID, userID, instanceID)
	if viewErr != nil && !zerrors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if zerrors.IsNotFound(viewErr) {

		token = new(model.TokenView)
		token.ID = tokenID
		token.UserID = userID
		token.InstanceID = instanceID
		if sequence != nil {
			token.ChangeDate = sequence.EventCreatedAt
		}
	}

	events, esErr := repo.getUserEvents(ctx, userID, token.InstanceID, token.ChangeDate, token.GetRelevantEventTypes())
	if zerrors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-4T90g", "Errors.Token.NotFound")
	}

	if esErr != nil {
		logging.Log("EVENT-5Nm9s").WithError(viewErr).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error retrieving new events")
		return model.TokenViewToModel(token), nil
	}
	viewToken := *token
	for _, event := range events {
		err := token.AppendEventIfMyToken(event)
		if err != nil {
			return model.TokenViewToModel(&viewToken), nil
		}
	}
	if !token.Expiration.After(time.Now().UTC()) || token.Deactivated {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.Token.NotFound")
	}
	return model.TokenViewToModel(token), nil
}

func (r *TokenRepo) getUserEvents(ctx context.Context, userID, instanceID string, changeDate time.Time, eventTypes []eventstore.EventType) ([]eventstore.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, instanceID, changeDate, eventTypes)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.Filter(ctx, query)
}
