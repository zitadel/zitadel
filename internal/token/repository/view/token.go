package view

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/errors"
	token_model "github.com/caos/zitadel/internal/token/model"
	"github.com/caos/zitadel/internal/token/repository/view/model"
	"github.com/caos/zitadel/internal/view"
)

func TokenByID(db *gorm.DB, table, tokenID string) (*model.Token, error) {
	token := new(model.Token)
	query := view.PrepareGetByKey(table, model.TokenSearchKey(token_model.TOKENSEARCHKEY_TOKEN_ID), tokenID)
	err := query(db, token)
	return token, err
}

func IsTokenValid(db *gorm.DB, table, tokenID string) (bool, error) {
	token, err := TokenByID(db, table, tokenID)
	if err == nil {
		return token.Expiration.After(time.Now().UTC()), nil
	}
	if errors.IsNotFound(err) {
		return false, nil
	}
	return false, err
}

//
//func TokenByIDs(db *gorm.DB, table, agentID, userID string) (*model.TokenView, error) {
//	userSession := new(model.TokenView)
//	userAgentQuery := model.TokenSearchQuery{
//		Key:    token_model.USERSESSIONSEARCHKEY_USER_AGENT_ID,
//		Method: global_model.SEARCHMETHOD_EQUALS,
//		Value:  agentID,
//	}
//	userQuery := model.TokenSearchQuery{
//		Key:    token_model.USERSESSIONSEARCHKEY_USER_ID,
//		Method: global_model.SEARCHMETHOD_EQUALS,
//		Value:  userID,
//	}
//	query := view.PrepareGetByQuery(table, userAgentQuery, userQuery)
//	err := query(db, userSession)
//	return userSession, err
//}
//
//func TokensByAgentID(db *gorm.DB, table, agentID string) ([]*model.TokenView, error) {
//	userSessions := make([]*model.TokenView, 0)
//	userAgentQuery := &token_model.TokenSearchQuery{
//		Key:    token_model.USERSESSIONSEARCHKEY_USER_AGENT_ID,
//		Method: global_model.SEARCHMETHOD_EQUALS,
//		Value:  agentID,
//	}
//	query := view.PrepareSearchQuery(table, model.TokenSearchRequest{
//		Queries: []*token_model.TokenSearchQuery{userAgentQuery},
//	})
//	_, err := query(db, userSessions)
//	return userSessions, err
//}

func PutToken(db *gorm.DB, table string, token *model.Token) error {
	save := view.PrepareSave(table)
	return save(db, token)
}

func DeleteToken(db *gorm.DB, table, tokenID string) error {
	delete := view.PrepareDeleteByKey(table, model.TokenSearchKey(token_model.TOKENSEARCHKEY_TOKEN_ID), tokenID)
	return delete(db)
}
