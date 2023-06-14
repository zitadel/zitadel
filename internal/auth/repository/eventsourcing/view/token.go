package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	tokenTable = "auth.tokens"
)

func (v *View) TokenByIDs(tokenID, userID, instanceID string) (*model.TokenView, error) {
	return usr_view.TokenByIDs(v.Db, tokenTable, tokenID, userID, instanceID)
}

func (v *View) TokensByUserID(userID, instanceID string) ([]*model.TokenView, error) {
	return usr_view.TokensByUserID(v.Db, tokenTable, userID, instanceID)
}

func (v *View) PutToken(token *model.TokenView) error {
	return usr_view.PutToken(v.Db, tokenTable, token)
}

func (v *View) PutTokens(token []*model.TokenView) error {
	return usr_view.PutTokens(v.Db, tokenTable, token...)
}

func (v *View) DeleteToken(tokenID, instanceID string) error {
	err := usr_view.DeleteToken(v.Db, tokenTable, tokenID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteSessionTokens(agentID string, event eventstore.Event) error {
	err := usr_view.DeleteSessionTokens(v.Db, tokenTable, agentID, event.Aggregate().ID, event.Aggregate().InstanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteUserTokens(event eventstore.Event) error {
	err := usr_view.DeleteUserTokens(v.Db, tokenTable, event.Aggregate().ID, event.Aggregate().InstanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteApplicationTokens(event eventstore.Event, ids ...string) error {
	err := usr_view.DeleteApplicationTokens(v.Db, tokenTable, event.Aggregate().InstanceID, ids)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteTokensFromRefreshToken(refreshTokenID, instanceID string) error {
	err := usr_view.DeleteTokensFromRefreshToken(v.Db, tokenTable, refreshTokenID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteInstanceTokens(event eventstore.Event) error {
	err := usr_view.DeleteInstanceTokens(v.Db, tokenTable, event.Aggregate().InstanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteOrgTokens(event eventstore.Event) error {
	err := usr_view.DeleteOrgTokens(v.Db, tokenTable, event.Aggregate().InstanceID, event.Aggregate().ResourceOwner)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}
