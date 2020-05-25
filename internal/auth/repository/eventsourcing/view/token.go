package view

import (
	"time"

	"github.com/caos/zitadel/internal/token/repository/view"
	"github.com/caos/zitadel/internal/token/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
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

func (v *View) CreateToken(agentID, applicationID, userID string, scopes []string, lifetime time.Duration) (*model.Token, error) {
	id, err := v.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	token := &model.Token{
		ID:            id,
		CreationDate:  now,
		UserID:        userID,
		ApplicationID: applicationID,
		UserAgentID:   agentID,
		Scopes:        scopes,
		Expiration:    now.Add(lifetime),
	}
	if err := view.PutToken(v.Db, tokenTable, token); err != nil {
		return nil, err
	}
	return token, nil
}

func (v *View) PutToken(token *model.Token) error {
	err := view.PutToken(v.Db, tokenTable, token)
	if err != nil {
		return err
	}
	return v.ProcessedTokenSequence(token.Sequence)
}

func (v *View) DeleteToken(tokenID string, eventSequence uint64) error {
	err := view.DeleteToken(v.Db, tokenTable, tokenID)
	if err != nil {
		return nil
	}
	return v.ProcessedTokenSequence(eventSequence)
}

func (v *View) DeleteSessionTokens(agentID, userID string, eventSequence uint64) error {
	err := view.DeleteTokens(v.Db, tokenTable, agentID, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedTokenSequence(eventSequence)
}

func (v *View) GetLatestTokenSequence() (uint64, error) {
	return v.latestSequence(tokenTable)
}

func (v *View) ProcessedTokenSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(tokenTable, eventSequence)
}

func (v *View) GetLatestTokenFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(tokenTable, sequence)
}

func (v *View) ProcessedTokenFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
