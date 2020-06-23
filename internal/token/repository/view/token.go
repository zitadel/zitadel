package view

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/errors"
	token_model "github.com/caos/zitadel/internal/token/model"
	"github.com/caos/zitadel/internal/token/repository/view/model"
	"github.com/caos/zitadel/internal/view"
)

func TokenByID(db *gorm.DB, table, tokenID string) (*model.Token, error) {
	token := new(model.Token)
	query := view.PrepareGetByKey(table, model.TokenSearchKey(token_model.TokenSearchKeyTokenID), tokenID)
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

func PutToken(db *gorm.DB, table string, token *model.Token) error {
	save := view.PrepareSave(table)
	return save(db, token)
}

func DeleteToken(db *gorm.DB, table, tokenID string) error {
	delete := view.PrepareDeleteByKey(table, model.TokenSearchKey(token_model.TokenSearchKeyTokenID), tokenID)
	return delete(db)
}

func DeleteSessionTokens(db *gorm.DB, table, agentID, userID string) error {
	delete := view.PrepareDeleteByKeys(table,
		view.Key{Key: model.TokenSearchKey(token_model.TokenSearchKeyUserAgentID), Value: agentID},
		view.Key{Key: model.TokenSearchKey(token_model.TokenSearchKeyUserID), Value: userID},
	)
	return delete(db)
}

func DeleteUserTokens(db *gorm.DB, table, userID string) error {
	delete := view.PrepareDeleteByKey(table, model.TokenSearchKey(token_model.TokenSearchKeyUserID), userID)
	return delete(db)
}

func DeleteApplicationTokens(db *gorm.DB, table string, appIDs []string) error {
	delete := view.PrepareDeleteByKey(table, model.TokenSearchKey(token_model.TokenSearchKeyApplicationID), pq.StringArray(appIDs))
	return delete(db)
}
