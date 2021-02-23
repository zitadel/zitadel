package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	tokenTable = "auth.tokens"
)

func (v *View) TokenByID(tokenID string) (*model.TokenView, error) {
	return usr_view.TokenByID(v.Db, tokenTable, tokenID)
}

func (v *View) TokensByUserID(userID string) ([]*model.TokenView, error) {
	return usr_view.TokensByUserID(v.Db, tokenTable, userID)
}

func (v *View) PutToken(token *model.TokenView, event *models.Event) error {
	err := usr_view.PutToken(v.Db, tokenTable, token)
	if err != nil {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) PutTokens(token []*model.TokenView, event *models.Event) error {
	err := usr_view.PutTokens(v.Db, tokenTable, token...)
	if err != nil {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) DeleteToken(tokenID string, event *models.Event) error {
	err := usr_view.DeleteToken(v.Db, tokenTable, tokenID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) DeleteSessionTokens(agentID, userID string, event *models.Event) error {
	err := usr_view.DeleteSessionTokens(v.Db, tokenTable, agentID, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) DeleteUserTokens(userID string, event *models.Event) error {
	err := usr_view.DeleteUserTokens(v.Db, tokenTable, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) DeleteApplicationTokens(event *models.Event, ids ...string) error {
	err := usr_view.DeleteApplicationTokens(v.Db, tokenTable, ids)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) GetLatestTokenSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(tokenTable)
}

func (v *View) ProcessedTokenSequence(event *models.Event) error {
	return v.saveCurrentSequence(tokenTable, event)
}

func (v *View) UpdateTokenSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(tokenTable)
}

func (v *View) GetLatestTokenFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(tokenTable, sequence)
}

func (v *View) ProcessedTokenFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
