package eventstore

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"time"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
)

type TokenRepo struct {
	UserEvents *user_event.UserEventstore
	View       *view.View
}

func (repo *TokenRepo) IsTokenValid(ctx context.Context, userID, tokenID string) (bool, error) {
	token, err := repo.TokenByID(ctx, userID, tokenID)
	if err == nil {
		return token.Expiration.After(time.Now().UTC()), nil
	}
	if errors.IsNotFound(err) {
		return false, nil
	}
	return false, err
}

func (repo *TokenRepo) TokenByID(ctx context.Context, userID, tokenID string) (*usr_model.TokenView, error) {
	token, viewErr := repo.View.TokenByID(tokenID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		token = new(model.TokenView)
		token.ID = tokenID
		token.UserID = userID
	}

	events, esErr := repo.UserEvents.UserEventsByID(ctx, userID, token.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-4T90g", "Errors.Token.NotFound")
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
		return nil, errors.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.Token.NotFound")
	}
	return model.TokenViewToModel(token), nil
}

func AppendAudIfNotExisting(aud string, existingAud []string) []string {
	for _, a := range existingAud {
		if a == aud {
			return existingAud
		}
	}
	return append(existingAud, aud)
}
