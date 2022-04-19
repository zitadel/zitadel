package view

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/user/model"
	usr_model "github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func TokenByID(db *gorm.DB, table, tokenID, instanceID string) (*usr_model.TokenView, error) {
	token := new(usr_model.TokenView)
	query := repository.PrepareGetByQuery(table,
		&usr_model.TokenSearchQuery{Key: model.TokenSearchKeyTokenID, Method: domain.SearchMethodEquals, Value: tokenID},
		&usr_model.TokenSearchQuery{Key: model.TokenSearchKeyInstanceID, Method: domain.SearchMethodEquals, Value: instanceID},
	)
	err := query(db, token)
	if errors.IsNotFound(err) {
		return nil, errors.ThrowNotFound(nil, "VIEW-6ub3p", "Errors.Token.NotFound")
	}
	return token, err
}

func TokensByUserID(db *gorm.DB, table, userID, instanceID string) ([]*usr_model.TokenView, error) {
	tokens := make([]*usr_model.TokenView, 0)
	userIDQuery := &model.TokenSearchQuery{
		Key:    model.TokenSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  userID,
	}
	instanceIDQuery := &model.TokenSearchQuery{
		Key:    model.TokenSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	query := repository.PrepareSearchQuery(table, usr_model.TokenSearchRequest{
		Queries: []*model.TokenSearchQuery{userIDQuery, instanceIDQuery},
	})
	_, err := query(db, &tokens)
	return tokens, err
}

func PutToken(db *gorm.DB, table string, token *usr_model.TokenView) error {
	save := repository.PrepareSave(table)
	return save(db, token)
}

func PutTokens(db *gorm.DB, table string, tokens ...*usr_model.TokenView) error {
	save := repository.PrepareBulkSave(table)
	t := make([]interface{}, len(tokens))
	for i, token := range tokens {
		t[i] = token
	}
	return save(db, t...)
}

func DeleteToken(db *gorm.DB, table, tokenID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{usr_model.TokenSearchKey(model.TokenSearchKeyTokenID), tokenID},
		repository.Key{usr_model.TokenSearchKey(model.TokenSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}

func DeleteSessionTokens(db *gorm.DB, table, agentID, userID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: usr_model.TokenSearchKey(model.TokenSearchKeyUserAgentID), Value: agentID},
		repository.Key{Key: usr_model.TokenSearchKey(model.TokenSearchKeyUserID), Value: userID},
		repository.Key{Key: usr_model.TokenSearchKey(model.TokenSearchKeyInstanceID), Value: instanceID},
	)
	return delete(db)
}

func DeleteUserTokens(db *gorm.DB, table, userID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{usr_model.TokenSearchKey(model.TokenSearchKeyUserID), userID},
		repository.Key{usr_model.TokenSearchKey(model.TokenSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}

func DeleteTokensFromRefreshToken(db *gorm.DB, table, refreshTokenID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{usr_model.TokenSearchKey(model.TokenSearchKeyRefreshTokenID), refreshTokenID},
		repository.Key{usr_model.TokenSearchKey(model.TokenSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}

func DeleteApplicationTokens(db *gorm.DB, table, instanceID string, appIDs []string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{usr_model.TokenSearchKey(model.TokenSearchKeyApplicationID), pq.StringArray(appIDs)},
		repository.Key{usr_model.TokenSearchKey(model.TokenSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}
