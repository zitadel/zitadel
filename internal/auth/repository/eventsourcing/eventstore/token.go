package eventstore

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	token_model "github.com/caos/zitadel/internal/token/model"
	token_view_model "github.com/caos/zitadel/internal/token/repository/view/model"
)

type TokenRepo struct {
	View *view.View
}

func (repo *TokenRepo) CreateToken(ctx context.Context, agentID, applicationID, userID string, scopes []string, lifetime time.Duration) (*token_model.Token, error) {
	token, err := repo.View.CreateToken(agentID, applicationID, userID, scopes, lifetime)
	if err != nil {
		return nil, err
	}
	return token_view_model.TokenToModel(token), nil
}

func (repo *TokenRepo) IsTokenValid(ctx context.Context, tokenID string) (bool, error) {
	return repo.View.IsTokenValid(tokenID)
}

func (repo *TokenRepo) TokenByID(ctx context.Context, tokenID string) (*token_model.Token, error) {
	token, err := repo.View.TokenByID(tokenID)
	if err != nil {
		return nil, err
	}
	return token_view_model.TokenToModel(token), nil
}
