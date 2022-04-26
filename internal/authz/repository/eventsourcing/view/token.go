package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	usr_view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	tokenTable = "auth.tokens"
)

func (v *View) TokenByID(tokenID, instanceID string) (*usr_view_model.TokenView, error) {
	return usr_view.TokenByID(v.Db, tokenTable, tokenID, instanceID)
}

func (v *View) PutToken(token *usr_view_model.TokenView, event *models.Event) error {
	err := usr_view.PutToken(v.Db, tokenTable, token)
	if err != nil {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) DeleteToken(tokenID, instanceID string, event *models.Event) error {
	err := usr_view.DeleteToken(v.Db, tokenTable, tokenID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) DeleteSessionTokens(agentID, userID, instanceID string, event *models.Event) error {
	err := usr_view.DeleteSessionTokens(v.Db, tokenTable, agentID, userID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) GetLatestTokenSequence(instanceID string) (*repository.CurrentSequence, error) {
	return v.latestSequence(tokenTable, instanceID)
}

func (v *View) ProcessedTokenSequence(event *models.Event) error {
	return v.saveCurrentSequence(tokenTable, event)
}

func (v *View) UpdateTokenSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(tokenTable)
}

func (v *View) GetLatestTokenFailedEvent(sequence uint64, instanceID string) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(tokenTable, instanceID, sequence)
}

func (v *View) ProcessedTokenFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
