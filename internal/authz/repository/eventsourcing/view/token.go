package view

import (
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	tokenTable = "auth.tokens"
)

func (v *View) TokenByID(tokenID string) (*usr_view_model.TokenView, error) {
	return usr_view.TokenByID(v.Db, tokenTable, tokenID)
}
