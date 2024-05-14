package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	user_model "github.com/zitadel/zitadel/internal/user/model"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	refreshTokenTable = "auth.refresh_tokens"
)

func (v *View) RefreshTokenByID(tokenID, instanceID string) (*model.RefreshTokenView, error) {
	return usr_view.RefreshTokenByID(v.Db, refreshTokenTable, tokenID, instanceID)
}

func (v *View) RefreshTokensByUserID(userID, instanceID string) ([]*model.RefreshTokenView, error) {
	return usr_view.RefreshTokensByUserID(v.Db, refreshTokenTable, userID, instanceID)
}

func (v *View) SearchRefreshTokens(request *user_model.RefreshTokenSearchRequest) ([]*model.RefreshTokenView, uint64, error) {
	return usr_view.SearchRefreshTokens(v.Db, refreshTokenTable, request)
}

func (v *View) PutRefreshToken(token *model.RefreshTokenView) error {
	return usr_view.PutRefreshToken(v.Db, refreshTokenTable, token)
}

func (v *View) PutRefreshTokens(token []*model.RefreshTokenView) error {
	return usr_view.PutRefreshTokens(v.Db, refreshTokenTable, token...)
}

func (v *View) DeleteRefreshToken(tokenID, instanceID string) error {
	err := usr_view.DeleteRefreshToken(v.Db, refreshTokenTable, tokenID, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteUserRefreshTokens(userID, instanceID string) error {
	err := usr_view.DeleteUserRefreshTokens(v.Db, refreshTokenTable, userID, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteApplicationRefreshTokens(event *models.Event, ids ...string) error {
	err := usr_view.DeleteApplicationTokens(v.Db, refreshTokenTable, event.InstanceID, ids)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteInstanceRefreshTokens(instanceID string) error {
	err := usr_view.DeleteInstanceRefreshTokens(v.Db, refreshTokenTable, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteOrgRefreshTokens(event eventstore.Event) error {
	err := usr_view.DeleteOrgRefreshTokens(v.Db, refreshTokenTable, event.Aggregate().InstanceID, event.Aggregate().ResourceOwner)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) GetLatestRefreshTokenSequence(ctx context.Context) (_ *query.CurrentState, err error) {
	q := &query.CurrentStateSearchQueries{
		Queries: make([]query.SearchQuery, 2),
	}
	q.Queries[0], err = query.NewCurrentStatesInstanceIDSearchQuery(authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	q.Queries[1], err = query.NewCurrentStatesProjectionSearchQuery(refreshTokenTable)
	if err != nil {
		return nil, err
	}
	states, err := v.query.SearchCurrentStates(ctx, q)
	if err != nil || states.SearchResponse.Count == 0 {
		return nil, err
	}
	return states.CurrentStates[0], nil
}
