package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query"
	user_model "github.com/zitadel/zitadel/internal/user/model"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	refreshTokenTable = "auth.refresh_tokens"
)

func (v *View) RefreshTokenByID(tokenID, instanceID string) (*model.RefreshTokenView, error) {
	return usr_view.RefreshTokenByID(v.Db, refreshTokenTable, tokenID, instanceID)
}

func (v *View) RefreshTokensByUserID(userID, instanceID string) ([]*model.RefreshTokenView, error) {
	return usr_view.RefreshTokensByUserID(v.Db, refreshTokenTable, userID, instanceID)
}

func (v *View) SearchRefreshTokens(request *user_model.RefreshTokenSearchRequest) ([]*model.RefreshTokenView, uint64, error) {
	return usr_view.SearchRefreshTokens(v.Db, refreshTokenTable, request)
}

func (v *View) GetLatestRefreshTokenSequence(ctx context.Context) (_ *query.CurrentState, err error) {
	q := &query.CurrentStateSearchQueries{
		Queries: make([]query.SearchQuery, 2),
	}
	q.Queries[0], err = query.NewCurrentStatesInstanceIDSearchQuery(authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	q.Queries[1], err = query.NewCurrentStatesProjectionSearchQuery(refreshTokenTable)
	if err != nil {
		return nil, err
	}
	states, err := v.query.SearchCurrentStates(ctx, q)
	if err != nil || states.SearchResponse.Count == 0 {
		return nil, err
	}
	return states.CurrentStates[0], nil
}
