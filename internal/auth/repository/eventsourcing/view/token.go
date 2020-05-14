package view

import (
	"time"

	"github.com/caos/zitadel/internal/token/repository/view"
	"github.com/caos/zitadel/internal/token/repository/view/model"
)

const (
	tokenTable = "auth.tokens"
)

func (v *View) TokenByID(tokenID string) (*model.Token, error) {
	return view.TokenByID(v.Db, tokenTable, tokenID)
}

func (v *View) IsTokenValid(tokenID string) (bool, error) {
	return view.IsTokenValid(v.Db, tokenTable, tokenID)
}

func (v *View) CreateToken(agentID, applicationID, userID string, lifetime time.Duration) (*model.Token, error) {
	now := time.Now().UTC()
	token := &model.Token{
		CreationDate:  now,
		UserID:        userID,
		ApplicationID: applicationID,
		UserAgentID:   agentID,
		Expiration:    now.Add(lifetime),
	}
	err := view.PutToken(v.Db, tokenTable, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
