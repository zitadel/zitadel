package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	userSessionTable = "auth.user_sessions"
)

func (v *View) UserSessionByIDs(agentID, userID, instanceID string) (*model.UserSessionView, error) {
	return view.UserSessionByIDs(v.client, agentID, userID, instanceID)
}

func (v *View) UserSessionsByAgentID(agentID, instanceID string) ([]*model.UserSessionView, error) {
	return view.UserSessionsByAgentID(v.client, agentID, instanceID)
}

func (v *View) PutUserSession(userSession *model.UserSessionView) error {
	return view.PutUserSession(v.Db, userSessionTable, userSession)
}

func (v *View) DeleteUserSessions(userID, instanceID string) error {
	err := view.DeleteUserSessions(v.Db, userSessionTable, userID, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteInstanceUserSessions(instanceID string) error {
	err := view.DeleteInstanceUserSessions(v.Db, userSessionTable, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteOrgUserSessions(event eventstore.Event) error {
	err := view.DeleteOrgUserSessions(v.Db, userSessionTable, event.Aggregate().InstanceID, event.Aggregate().ResourceOwner)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
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
