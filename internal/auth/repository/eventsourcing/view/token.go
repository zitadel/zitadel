package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteSessionTokens(agentID string, event eventstore.Event) error {
	err := usr_view.DeleteSessionTokens(v.Db, tokenTable, agentID, event.Aggregate().ID, event.Aggregate().InstanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteUserTokens(event eventstore.Event) error {
	err := usr_view.DeleteUserTokens(v.Db, tokenTable, event.Aggregate().ID, event.Aggregate().InstanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteApplicationTokens(event eventstore.Event, ids ...string) error {
	err := usr_view.DeleteApplicationTokens(v.Db, tokenTable, event.Aggregate().InstanceID, ids)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteTokensFromRefreshToken(refreshTokenID, instanceID string) error {
	err := usr_view.DeleteTokensFromRefreshToken(v.Db, tokenTable, refreshTokenID, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteInstanceTokens(event eventstore.Event) error {
	err := usr_view.DeleteInstanceTokens(v.Db, tokenTable, event.Aggregate().InstanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteOrgTokens(event eventstore.Event) error {
	err := usr_view.DeleteOrgTokens(v.Db, tokenTable, event.Aggregate().InstanceID, event.Aggregate().ResourceOwner)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) GetLatestTokenSequence(ctx context.Context, instanceID string) (_ *query.CurrentState, err error) {
	q := &query.CurrentStateSearchQueries{
		Queries: make([]query.SearchQuery, 2),
	}
	q.Queries[0], err = query.NewCurrentStatesInstanceIDSearchQuery(instanceID)
	if err != nil {
		return nil, err
	}
	q.Queries[1], err = query.NewCurrentStatesProjectionSearchQuery(tokenTable)
	if err != nil {
		return nil, err
	}
	states, err := v.query.SearchCurrentStates(ctx, q)
	if err != nil || states.SearchResponse.Count == 0 {
		return nil, err
	}
	return states.CurrentStates[0], nil
}
