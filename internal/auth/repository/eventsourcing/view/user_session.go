package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
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

func (v *View) ActiveUserSessionsCount() (uint64, error) {
	return view.ActiveUserSessions(v.Db, userSessionTable)
}

func (v *View) PutUserSession(userSession *model.UserSessionView, event *models.Event) error {
	err := view.PutUserSession(v.Db, userSessionTable, userSession)
	if err != nil {
		return err
	}
	return v.ProcessedUserSessionSequence(event)
}

func (v *View) PutUserSessions(userSession []*model.UserSessionView, event *models.Event) error {
	err := view.PutUserSessions(v.Db, userSessionTable, userSession...)
	if err != nil {
		return err
	}
	return v.ProcessedUserSessionSequence(event)
}

func (v *View) DeleteUserSessions(userID, instanceID string, event *models.Event) error {
	err := view.DeleteUserSessions(v.Db, userSessionTable, userID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserSessionSequence(event)
}

func (v *View) GetLatestUserSessionSequence(instanceID string) (*repository.CurrentSequence, error) {
	return v.latestSequence(userSessionTable, instanceID)
}

func (v *View) GetLatestUserSessionSequences() ([]*repository.CurrentSequence, error) {
	return v.latestSequences(userSessionTable)
}

func (v *View) ProcessedUserSessionSequence(event *models.Event) error {
	return v.saveCurrentSequence(userSessionTable, event)
}

func (v *View) UpdateUserSessionSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(userSessionTable)
}

func (v *View) GetLatestUserSessionFailedEvent(sequence uint64, instanceID string) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userSessionTable, instanceID, sequence)
}

func (v *View) ProcessedUserSessionFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
