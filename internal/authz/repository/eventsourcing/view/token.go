package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	usr_view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	tokenTable = "auth.tokens"
)

func (v *View) TokenByIDs(tokenID, userID, instanceID string) (*usr_view_model.TokenView, error) {
	return usr_view.TokenByIDs(v.Db, tokenTable, tokenID, userID, instanceID)
}

func (v *View) GetLatestState(ctx context.Context) (_ *query.CurrentState, err error) {
	q := &query.CurrentStateSearchQueries{
		Queries: make([]query.SearchQuery, 2),
	}
	q.Queries[0], err = query.NewCurrentStatesInstanceIDSearchQuery(authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	q.Queries[1], err = query.NewCurrentStatesProjectionSearchQuery(tokenTable)
	if err != nil {
		return nil, err
	}
	states, err := v.Query.SearchCurrentStates(ctx, q)
	if err != nil || states.SearchResponse.Count == 0 {
		return nil, err
	}
	return states.CurrentStates[0], nil
}
