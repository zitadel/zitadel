package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/user/model"
	usr_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func RefreshTokenByID(db *gorm.DB, table, tokenID, instanceID string) (*usr_model.RefreshTokenView, error) {
	token := new(usr_model.RefreshTokenView)
	query := repository.PrepareGetByQuery(table,
		&usr_model.RefreshTokenSearchQuery{Key: model.RefreshTokenSearchKeyRefreshTokenID, Method: domain.SearchMethodEquals, Value: tokenID},
		&usr_model.RefreshTokenSearchQuery{Key: model.RefreshTokenSearchKeyInstanceID, Method: domain.SearchMethodEquals, Value: instanceID},
	)
	err := query(db, token)
	if zerrors.IsNotFound(err) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-6ub3p", "Errors.RefreshToken.NotFound")
	}
	return token, err
}

func RefreshTokensByUserID(db *gorm.DB, table, userID, instanceID string) ([]*usr_model.RefreshTokenView, error) {
	tokens := make([]*usr_model.RefreshTokenView, 0)
	userIDQuery := &model.RefreshTokenSearchQuery{
		Key:    model.RefreshTokenSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  userID,
	}
	instanceIDQuery := &model.RefreshTokenSearchQuery{
		Key:    model.RefreshTokenSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	query := repository.PrepareSearchQuery(table, usr_model.RefreshTokenSearchRequest{
		Queries: []*model.RefreshTokenSearchQuery{userIDQuery, instanceIDQuery},
	})
	_, err := query(db, &tokens)
	return tokens, err
}

func PutRefreshToken(db *gorm.DB, table string, token *usr_model.RefreshTokenView) error {
	save := repository.PrepareSaveOnConflict(table,
		[]string{"client_id", "user_agent_id", "user_id"},
		[]string{"id", "creation_date", "change_date", "token", "auth_time", "idle_expiration", "expiration", "sequence", "scopes", "audience", "amr"},
	)
	return save(db, token)
}

func PutRefreshTokens(db *gorm.DB, table string, tokens ...*usr_model.RefreshTokenView) error {
	save := repository.PrepareBulkSave(table)
	t := make([]interface{}, len(tokens))
	for i, token := range tokens {
		t[i] = token
	}
	return save(db, t...)
}

func SearchRefreshTokens(db *gorm.DB, table string, req *model.RefreshTokenSearchRequest) ([]*usr_model.RefreshTokenView, uint64, error) {
	tokens := make([]*usr_model.RefreshTokenView, 0)
	query := repository.PrepareSearchQuery(table, usr_model.RefreshTokenSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &tokens)
	return tokens, count, err
}

func DeleteRefreshToken(db *gorm.DB, table, tokenID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyRefreshTokenID), tokenID},
		repository.Key{usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}

func DeleteSessionRefreshTokens(db *gorm.DB, table, agentID, userID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyUserAgentID), Value: agentID},
		repository.Key{Key: usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyUserID), Value: userID},
	)
	return delete(db)
}

func DeleteUserRefreshTokens(db *gorm.DB, table, userID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyUserID), userID},
		repository.Key{usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}

func DeleteApplicationRefreshTokens(db *gorm.DB, table string, instanceID string, appIDs []string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyInstanceID), Value: instanceID},
		repository.Key{Key: usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyApplicationID), Value: appIDs},
	)
	return delete(db)
}

func DeleteOrgRefreshTokens(db *gorm.DB, table string, instanceID, orgID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyInstanceID), Value: instanceID},
		repository.Key{Key: usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyResourceOwner), Value: orgID},
	)
	return delete(db)
}

func DeleteInstanceRefreshTokens(db *gorm.DB, table string, instanceID string) error {
	delete := repository.PrepareDeleteByKey(table,
		usr_model.RefreshTokenSearchKey(model.RefreshTokenSearchKeyInstanceID),
		instanceID,
	)
	return delete(db)
}
