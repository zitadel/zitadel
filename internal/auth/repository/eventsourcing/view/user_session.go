package view

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/query"
	"github.com/zitadel/zitadel/v2/internal/user/repository/view"
	"github.com/zitadel/zitadel/v2/internal/user/repository/view/model"
)

const (
	userSessionTable = "auth.user_sessions"
)

func (v *View) UserSessionByIDs(ctx context.Context, agentID, userID, instanceID string) (*model.UserSessionView, error) {
	return view.UserSessionByIDs(ctx, v.client, agentID, userID, instanceID)
}

func (v *View) UserSessionsByAgentID(ctx context.Context, agentID, instanceID string) ([]*model.UserSessionView, error) {
	return view.UserSessionsByAgentID(ctx, v.client, agentID, instanceID)
}

func (v *View) UserAgentIDBySessionID(ctx context.Context, sessionID, instanceID string) (string, error) {
	return view.UserAgentIDBySessionID(ctx, v.client, sessionID, instanceID)
}

func (v *View) ActiveUserIDsBySessionID(ctx context.Context, sessionID, instanceID string) (userAgentID string, userIDs []string, err error) {
	return view.ActiveUserIDsBySessionID(ctx, v.client, sessionID, instanceID)
}

func (v *View) GetLatestUserSessionSequence(ctx context.Context, instanceID string) (_ *query.CurrentState, err error) {
	q := &query.CurrentStateSearchQueries{
		Queries: make([]query.SearchQuery, 2),
	}
	q.Queries[0], err = query.NewCurrentStatesInstanceIDSearchQuery(instanceID)
	if err != nil {
		return nil, err
	}
	q.Queries[1], err = query.NewCurrentStatesProjectionSearchQuery(userSessionTable)
	if err != nil {
		return nil, err
	}
	states, err := v.query.SearchCurrentStates(ctx, q)
	if err != nil || states.SearchResponse.Count == 0 {
		return nil, err
	}
	return states.CurrentStates[0], nil
}
