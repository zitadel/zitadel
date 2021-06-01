package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	user_model "github.com/caos/zitadel/internal/user/model"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	refreshTokenTable = "auth.refresh_tokens"
)

func (v *View) RefreshTokenByID(tokenID string) (*model.RefreshTokenView, error) {
	return usr_view.RefreshTokenByID(v.Db, refreshTokenTable, tokenID)
}

func (v *View) RefreshTokensByUserID(userID string) ([]*model.RefreshTokenView, error) {
	return usr_view.RefreshTokensByUserID(v.Db, refreshTokenTable, userID)
}

func (v *View) SearchRefreshTokens(request *user_model.RefreshTokenSearchRequest) ([]*model.RefreshTokenView, uint64, error) {
	return usr_view.SearchRefreshTokens(v.Db, refreshTokenTable, request)
}

func (v *View) PutRefreshToken(token *model.RefreshTokenView, event *models.Event) error {
	err := usr_view.PutRefreshToken(v.Db, refreshTokenTable, token)
	if err != nil {
		return err
	}
	return v.ProcessedTokenSequence(event)
}

func (v *View) PutRefreshTokens(token []*model.RefreshTokenView, event *models.Event) error {
	err := usr_view.PutRefreshTokens(v.Db, refreshTokenTable, token...)
	if err != nil {
		return err
	}
	return v.ProcessedRefreshTokenSequence(event)
}

func (v *View) DeleteRefreshToken(tokenID string, event *models.Event) error {
	err := usr_view.DeleteRefreshToken(v.Db, refreshTokenTable, tokenID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedRefreshTokenSequence(event)
}

func (v *View) DeleteUserRefreshTokens(userID string, event *models.Event) error {
	err := usr_view.DeleteUserRefreshTokens(v.Db, refreshTokenTable, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedRefreshTokenSequence(event)
}

func (v *View) DeleteApplicationRefreshTokens(event *models.Event, ids ...string) error {
	err := usr_view.DeleteApplicationTokens(v.Db, refreshTokenTable, ids)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedRefreshTokenSequence(event)
}

func (v *View) GetLatestRefreshTokenSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(refreshTokenTable)
}

func (v *View) ProcessedRefreshTokenSequence(event *models.Event) error {
	return v.saveCurrentSequence(refreshTokenTable, event)
}

func (v *View) UpdateRefreshTokenSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(refreshTokenTable)
}

func (v *View) GetLatestRefreshTokenFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(refreshTokenTable, sequence)
}

func (v *View) ProcessedRefreshTokenFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
