package view

import (
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	userSessionTable = "auth.user_sessions"
)

func (v *View) UserSessionByIDs(agentID, userID string) (*model.UserSessionView, error) {
	return view.UserSessionByIDs(v.Db, userSessionTable, agentID, userID)
}

func (v *View) UserSessionsByUserID(userID string) ([]*model.UserSessionView, error) {
	return view.UserSessionsByUserID(v.Db, userSessionTable, userID)
}

func (v *View) UserSessionsByAgentID(agentID string) ([]*model.UserSessionView, error) {
	return view.UserSessionsByAgentID(v.Db, userSessionTable, agentID)
}

func (v *View) ActiveUserSessions() ([]*model.UserSessionView, error) {
	return view.ActiveUserSessions(v.Db, userSessionTable)
}

func (v *View) PutUserSession(userSession *model.UserSessionView) error {
	err := view.PutUserSession(v.Db, userSessionTable, userSession)
	if err != nil {
		return err
	}
	return v.ProcessedUserSessionSequence(userSession.Sequence)
}

func (v *View) PutUserSessions(userSession []*model.UserSessionView, sequence uint64) error {
	err := view.PutUserSessions(v.Db, userSessionTable, userSession...)
	if err != nil {
		return err
	}
	return v.ProcessedUserSessionSequence(sequence)
}

func (v *View) DeleteUserSessions(userID string, eventSequence uint64) error {
	err := view.DeleteUserSessions(v.Db, userSessionTable, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedUserSessionSequence(eventSequence)
}

func (v *View) GetLatestUserSessionSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(userSessionTable)
}

func (v *View) ProcessedUserSessionSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(userSessionTable, eventSequence)
}

func (v *View) GetLatestUserSessionFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userSessionTable, sequence)
}

func (v *View) ProcessedUserSessionFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
