package view

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	token_model "github.com/caos/zitadel/internal/token/model"
	"github.com/caos/zitadel/internal/token/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func TokenByID(db *gorm.DB, table, tokenID string) (*model.Token, error) {
	token := new(model.Token)
	query := repository.PrepareGetByKey(table, model.TokenSearchKey(token_model.TokenSearchKeyTokenID), tokenID)
	err := query(db, token)
	if errors.IsNotFound(err) {
		return nil, errors.ThrowNotFound(nil, "VIEW-6ub3p", "Errors.Token.NotFound")
	}
	return token, err
}

func ValidTokenByID(db *gorm.DB, table, tokenID string) (*model.Token, error) {
	token := new(model.Token)
	query := repository.PrepareGetByQuery(table,
		model.TokenSearchQuery{
			Key:    token_model.TokenSearchKeyTokenID,
			Method: global_model.SearchMethodEquals,
			Value:  tokenID,
		},
		model.TokenSearchQuery{
			Key:    token_model.TokenSearchKeyExpiration,
			Method: global_model.SearchMethodGreaterThan,
			Value:  time.Now().UTC(),
		},
	)
	err := query(db, token)
	if errors.IsNotFound(err) {
		return nil, errors.ThrowNotFound(nil, "VIEW-6ub3p", "Errors.Token.NotFound")
	}
	return token, err
}

func TokensByUserID(db *gorm.DB, table, userID string) ([]*model.Token, error) {
	tokens := make([]*model.Token, 0)
	userIDQuery := &token_model.TokenSearchQuery{
		Key:    token_model.TokenSearchKeyUserID,
		Method: global_model.SearchMethodEquals,
		Value:  userID,
	}
	query := repository.PrepareSearchQuery(table, model.TokenSearchRequest{
		Queries: []*token_model.TokenSearchQuery{userIDQuery},
	})
	_, err := query(db, &tokens)
	return tokens, err
}

func PutToken(db *gorm.DB, table string, token *model.Token) error {
	save := repository.PrepareSave(table)
	return save(db, token)
}

func PutTokens(db *gorm.DB, table string, tokens ...*model.Token) error {
	save := repository.PrepareBulkSave(table)
	t := make([]interface{}, len(tokens))
	for i, token := range tokens {
		t[i] = token
	}
	return save(db, t...)
}

func DeleteToken(db *gorm.DB, table, tokenID string) error {
	delete := repository.PrepareDeleteByKey(table, model.TokenSearchKey(token_model.TokenSearchKeyTokenID), tokenID)
	return delete(db)
}

func DeleteSessionTokens(db *gorm.DB, table, agentID, userID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.TokenSearchKey(token_model.TokenSearchKeyUserAgentID), Value: agentID},
		repository.Key{Key: model.TokenSearchKey(token_model.TokenSearchKeyUserID), Value: userID},
	)
	return delete(db)
}

func DeleteUserTokens(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.TokenSearchKey(token_model.TokenSearchKeyUserID), userID)
	return delete(db)
}

func DeleteApplicationTokens(db *gorm.DB, table string, appIDs []string) error {
	delete := repository.PrepareDeleteByKey(table, model.TokenSearchKey(token_model.TokenSearchKeyApplicationID), pq.StringArray(appIDs))
	return delete(db)
}
