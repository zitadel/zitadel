package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	userSessionTable = "auth.user_sessions"
)

func (v *View) UserSessionByIDs(agentID, userID, instanceID string) (*model.UserSessionView, error) {
	return view.UserSessionByIDs(v.Db, userSessionTable, agentID, userID, instanceID)
}

func (v *View) UserSessionsByUserID(userID, instanceID string) ([]*model.UserSessionView, error) {
	return view.UserSessionsByUserID(v.Db, userSessionTable, userID, instanceID)
}

func (v *View) UserSessionsByAgentID(agentID, instanceID string) ([]*model.UserSessionView, error) {
	return view.UserSessionsByAgentID(v.Db, userSessionTable, agentID, instanceID)
}

func (v *View) UserSessionsByOrgID(orgID, instanceID string) ([]*model.UserSessionView, error) {
	return view.UserSessionsByOrgID(v.Db, userSessionTable, orgID, instanceID)
}

func (v *View) ActiveUserSessionsCount() (uint64, error) {
	return view.ActiveUserSessions(v.Db, userSessionTable)
}

func (v *View) PutUserSession(userSession *model.UserSessionView) error {
	return view.PutUserSession(v.Db, userSessionTable, userSession)
}

func (v *View) PutUserSessions(userSession []*model.UserSessionView) error {
	return view.PutUserSessions(v.Db, userSessionTable, userSession...)
}

func (v *View) DeleteUserSessions(userID, instanceID string) error {
	err := view.DeleteUserSessions(v.Db, userSessionTable, userID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteInstanceUserSessions(instanceID string) error {
	err := view.DeleteInstanceUserSessions(v.Db, userSessionTable, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteOrgUserSessions(event eventstore.Event) error {
	err := view.DeleteOrgUserSessions(v.Db, userSessionTable, event.Aggregate().InstanceID, event.Aggregate().ResourceOwner)
	if err != nil && !errors.IsNotFound(err) {
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
