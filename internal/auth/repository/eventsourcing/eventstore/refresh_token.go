package eventstore

import (
	"context"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

type RefreshTokenRepo struct {
	Eventstore   v1.Eventstore
	View         *view.View
	SearchLimit  uint64
	KeyAlgorithm crypto.EncryptionAlgorithm
}

func (r *RefreshTokenRepo) RefreshTokenByID(ctx context.Context, refreshToken string) (*usr_model.RefreshTokenView, error) {
	userID, tokenID, token, err := domain.FromRefreshToken(refreshToken, r.KeyAlgorithm)
	if err != nil {
		return nil, err
	}
	tokenView, viewErr := r.View.RefreshTokenByID(tokenID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		tokenView = new(model.RefreshTokenView)
		tokenView.ID = tokenID
		tokenView.UserID = userID
	}

	events, esErr := r.getUserEvents(ctx, userID, tokenView.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-BHB52", "Errors.User.RefreshToken.Invalid")
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
	if !tokenView.Expiration.After(time.Now()) || tokenView.Token != token {
		return nil, errors.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.User.RefreshToken.Invalid")
	}
	return model.RefreshTokenViewToModel(tokenView), nil
}

func (r *RefreshTokenRepo) SearchMyRefreshTokens(ctx context.Context, userID string, request *usr_model.RefreshTokenSearchRequest) (*usr_model.RefreshTokenSearchResponse, error) {
	err := request.EnsureLimit(r.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, err := r.View.GetLatestRefreshTokenSequence()
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
		Sequence:    sequence.CurrentSequence,
		Timestamp:   sequence.LastSuccessfulSpoolerRun,
		Result:      model.RefreshTokenViewsToModel(tokens),
	}, nil
}

func (r *RefreshTokenRepo) getUserEvents(ctx context.Context, userID string, sequence uint64) ([]*models.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.FilterEvents(ctx, query)
}
