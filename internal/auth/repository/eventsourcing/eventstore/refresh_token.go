package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1"
	"time"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

type RefreshTokenRepo struct {
	Eventstore   v1.Eventstore
	View         *view.View
	KeyAlgorithm crypto.EncryptionAlgorithm
}

//func (repo *RefreshTokenRepo) IsTokenValid(ctx context.Context, userID, tokenID string) (bool, error) {
//	token, err := repo.TokenByID(ctx, userID, tokenID)
//	if err == nil {
//		return token.Expiration.After(time.Now().UTC()), nil
//	}
//	if errors.IsNotFound(err) {
//		return false, nil
//	}
//	return false, err
//}

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

func (r *RefreshTokenRepo) getUserEvents(ctx context.Context, userID string, sequence uint64) ([]*models.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.FilterEvents(ctx, query)
}
