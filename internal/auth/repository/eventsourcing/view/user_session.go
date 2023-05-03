package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
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

func (v *View) DeleteInstanceUserSessions(event *models.Event) error {
	err := view.DeleteInstanceUserSessions(v.Db, userSessionTable, event.InstanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserSessionSequence(event)
}

func (v *View) DeleteOrgUserSessions(event *models.Event) error {
	err := view.DeleteOrgUserSessions(v.Db, userSessionTable, event.InstanceID, event.ResourceOwner)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserSessionSequence(event)
}

func (v *View) GetLatestUserSessionSequence(ctx context.Context, instanceID string) (*repository.CurrentSequence, error) {
	return v.latestSequence(ctx, userSessionTable, instanceID)
}

func (v *View) GetLatestUserSessionSequences(ctx context.Context, instanceIDs []string) ([]*repository.CurrentSequence, error) {
	return v.latestSequences(ctx, userSessionTable, instanceIDs)
}

func (v *View) ProcessedUserSessionSequence(event *models.Event) error {
	return v.saveCurrentSequence(userSessionTable, event)
}

func (v *View) UpdateUserSessionSpoolerRunTimestamp(instanceIDs []string) error {
	return v.updateSpoolerRunSequence(userSessionTable, instanceIDs)
}

func (v *View) GetLatestUserSessionFailedEvent(sequence uint64, instanceID string) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userSessionTable, instanceID, sequence)
}

func (v *View) ProcessedUserSessionFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
