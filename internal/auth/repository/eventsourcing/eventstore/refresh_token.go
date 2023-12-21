package eventstore

import (
	"context"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type RefreshTokenRepo struct {
	Eventstore   *eventstore.Eventstore
	View         *view.View
	SearchLimit  uint64
	KeyAlgorithm crypto.EncryptionAlgorithm
}

func (r *RefreshTokenRepo) RefreshTokenByToken(ctx context.Context, refreshToken string) (*usr_model.RefreshTokenView, error) {
	userID, tokenID, token, err := domain.FromRefreshToken(refreshToken, r.KeyAlgorithm)
	if err != nil {
		return nil, err
	}
	tokenView, err := r.RefreshTokenByID(ctx, tokenID, userID)
	if err != nil {
		return nil, err
	}
	if tokenView.Token != token {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.User.RefreshToken.Invalid")
	}
	return tokenView, nil
}

func (r *RefreshTokenRepo) RefreshTokenByID(ctx context.Context, tokenID, userID string) (*usr_model.RefreshTokenView, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()

	// always load the latest sequence first, so in case the token was not found by id,
	// the sequence will be equal or lower than the actual projection and no events are lost
	sequence, err := r.View.GetLatestRefreshTokenSequence(ctx)
	logging.WithFields("instanceID", instanceID, "userID", userID, "tokenID", tokenID).
		OnError(err).
		Errorf("could not get current sequence for RefreshTokenByID")

	tokenView, viewErr := r.View.RefreshTokenByID(tokenID, instanceID)
	if viewErr != nil && !zerrors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if zerrors.IsNotFound(viewErr) {
		tokenView = new(model.RefreshTokenView)
		tokenView.ID = tokenID
		tokenView.UserID = userID
		tokenView.InstanceID = instanceID
		if sequence != nil {
			tokenView.ChangeDate = sequence.EventCreatedAt
		}
	}

	events, esErr := r.getUserEvents(ctx, userID, tokenView.InstanceID, tokenView.ChangeDate, tokenView.GetRelevantEventTypes())
	if zerrors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-BHB52", "Errors.User.RefreshToken.Invalid")
	}

	if esErr != nil {
		logging.Log("EVENT-AE462").WithError(viewErr).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error retrieving new events")
		return model.RefreshTokenViewToModel(tokenView), nil
	}
	viewToken := *tokenView
	for _, event := range events {
		err := tokenView.AppendEventIfMyRefreshToken(event)
		if err != nil {
			return model.RefreshTokenViewToModel(&viewToken), nil
		}
	}
	if !tokenView.Expiration.After(time.Now()) {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.User.RefreshToken.Invalid")
	}
	return model.RefreshTokenViewToModel(tokenView), nil
}

func (r *RefreshTokenRepo) SearchMyRefreshTokens(ctx context.Context, userID string, request *usr_model.RefreshTokenSearchRequest) (*usr_model.RefreshTokenSearchResponse, error) {
	err := request.EnsureLimit(r.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, err := r.View.GetLatestRefreshTokenSequence(ctx)
	logging.Log("EVENT-GBdn4").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest refresh token sequence")
	request.Queries = append(request.Queries, &usr_model.RefreshTokenSearchQuery{Key: usr_model.RefreshTokenSearchKeyUserID, Method: domain.SearchMethodEquals, Value: userID})
	tokens, count, err := r.View.SearchRefreshTokens(request)
	if err != nil {
		return nil, err
	}
	return &usr_model.RefreshTokenSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Sequence:    sequence.Sequence,
		Timestamp:   sequence.LastRun,
		Result:      model.RefreshTokenViewsToModel(tokens),
	}, nil
}

func (r *RefreshTokenRepo) getUserEvents(ctx context.Context, userID, instanceID string, changeDate time.Time, eventTypes []eventstore.EventType) ([]eventstore.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, instanceID, changeDate, eventTypes)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.Filter(ctx, query)
}
