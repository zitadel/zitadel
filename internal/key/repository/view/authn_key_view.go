package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func AuthNKeyByIDs(db *gorm.DB, table, objectID, keyID string) (*model.AuthNKeyView, error) {
	key := new(model.AuthNKeyView)
	query := repository.PrepareGetByQuery(table,
		model.AuthNKeySearchQuery{Key: key_model.AuthNKeyObjectID, Method: global_model.SearchMethodEquals, Value: objectID},
		model.AuthNKeySearchQuery{Key: key_model.AuthNKeyKeyID, Method: global_model.SearchMethodEquals, Value: keyID},
	)
	err := query(db, key)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-3Dk9s", "Errors.User.KeyNotFound")
	}
	return key, err
}

func SearchAuthNKeys(db *gorm.DB, table string, req *key_model.AuthNKeySearchRequest) ([]*model.AuthNKeyView, uint64, error) {
	members := make([]*model.AuthNKeyView, 0)
	query := repository.PrepareSearchQuery(table, model.AuthNKeySearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &members)
	if err != nil {
		return nil, 0, err
	}
	return members, count, nil
}

func AuthNKeysByObjectID(db *gorm.DB, table string, objectID string) ([]*model.AuthNKeyView, error) {
	keys := make([]*model.AuthNKeyView, 0)
	queries := []*key_model.AuthNKeySearchQuery{
		{
			Key:    key_model.AuthNKeyObjectID,
			Value:  objectID,
			Method: global_model.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.AuthNKeySearchRequest{Queries: queries})
	_, err := query(db, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func AuthNKeyByID(db *gorm.DB, table string, keyID string) (*model.AuthNKeyView, error) {
	key := new(model.AuthNKeyView)
	query := repository.PrepareGetByQuery(table,
		model.AuthNKeySearchQuery{Key: key_model.AuthNKeyKeyID, Method: global_model.SearchMethodEquals, Value: keyID},
	)
	err := query(db, key)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-BjN6x", "Errors.User.KeyNotFound")
	}
	return key, err
}

func PutAuthNKey(db *gorm.DB, table string, role *model.AuthNKeyView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func DeleteAuthNKey(db *gorm.DB, table, keyID string) error {
	delete := repository.PrepareDeleteByKey(table, model.AuthNKeySearchKey(key_model.AuthNKeyKeyID), keyID)
	return delete(db)
}

func DeleteAuthNKeysByUserID(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.AuthNKeySearchKey(key_model.AuthNKeyObjectID), userID)
	return delete(db)
}
