package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	usr_view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	tokenTable = "auth.tokens"
)

func (v *View) TokenByIDs(tokenID, userID, instanceID string) (*usr_view_model.TokenView, error) {
	return usr_view.TokenByIDs(v.Db, tokenTable, tokenID, userID, instanceID)
}

func (v *View) PutToken(token *usr_view_model.TokenView, event *models.Event) error {
	return usr_view.PutToken(v.Db, tokenTable, token)
}

func (v *View) DeleteToken(tokenID, instanceID string, event *models.Event) error {
	err := usr_view.DeleteToken(v.Db, tokenTable, tokenID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteSessionTokens(agentID, userID, instanceID string, event *models.Event) error {
	err := usr_view.DeleteSessionTokens(v.Db, tokenTable, agentID, userID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) GetLatestState(ctx context.Context) (*handler.CurrentState, error) {
	return handler.QueryCurrentState(ctx, v.Db.DB(), tokenTable)
}
