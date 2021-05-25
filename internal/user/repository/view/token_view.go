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

func TokenByID(db *gorm.DB, table, tokenID string) (*usr_model.TokenView, error) {
	token := new(usr_model.TokenView)
	query := repository.PrepareGetByKey(table, usr_model.TokenSearchKey(model.TokenSearchKeyTokenID), tokenID)
	err := query(db, token)
	if errors.IsNotFound(err) {
		return nil, errors.ThrowNotFound(nil, "VIEW-6ub3p", "Errors.Token.NotFound")
	}
	return token, err
}

func TokensByUserID(db *gorm.DB, table, userID string) ([]*usr_model.TokenView, error) {
	tokens := make([]*usr_model.TokenView, 0)
	userIDQuery := &model.TokenSearchQuery{
		Key:    model.TokenSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  userID,
	}
	query := repository.PrepareSearchQuery(table, usr_model.TokenSearchRequest{
		Queries: []*model.TokenSearchQuery{userIDQuery},
	})
	_, err := query(db, &tokens)
	return tokens, err
}

func PutToken(db *gorm.DB, table string, token *usr_model.TokenView) error {
	save := repository.PrepareSaveOnConflict(table,
		[]string{"user_id", "user_agent_id", "application_id"},
		[]string{"id", "creation_date", "change_date", "expiration", "sequence", "scopes", "audience", "preferred_language"},
	)
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

func DeleteToken(db *gorm.DB, table, tokenID string) error {
	delete := repository.PrepareDeleteByKey(table, usr_model.TokenSearchKey(model.TokenSearchKeyTokenID), tokenID)
	return delete(db)
}

func DeleteSessionTokens(db *gorm.DB, table, agentID, userID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: usr_model.TokenSearchKey(model.TokenSearchKeyUserAgentID), Value: agentID},
		repository.Key{Key: usr_model.TokenSearchKey(model.TokenSearchKeyUserID), Value: userID},
	)
	return delete(db)
}

func DeleteUserTokens(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, usr_model.TokenSearchKey(model.TokenSearchKeyUserID), userID)
	return delete(db)
}

func DeleteApplicationTokens(db *gorm.DB, table string, appIDs []string) error {
	delete := repository.PrepareDeleteByKey(table, usr_model.TokenSearchKey(model.TokenSearchKeyApplicationID), pq.StringArray(appIDs))
	return delete(db)
}
