package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func MachineKeyByIDs(db *gorm.DB, table, userID, keyID string) (*model.MachineKeyView, error) {
	key := new(model.MachineKeyView)
	query := repository.PrepareGetByQuery(table,
		model.MachineKeySearchQuery{Key: usr_model.MachineKeyKeyUserID, Method: global_model.SearchMethodEquals, Value: userID},
		model.MachineKeySearchQuery{Key: usr_model.MachineKeyKeyID, Method: global_model.SearchMethodEquals, Value: keyID},
	)
	err := query(db, key)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-3Dk9s", "Errors.User.KeyNotFound")
	}
	return key, err
}

func SearchMachineKeys(db *gorm.DB, table string, req *usr_model.MachineKeySearchRequest) ([]*model.MachineKeyView, uint64, error) {
	members := make([]*model.MachineKeyView, 0)
	query := repository.PrepareSearchQuery(table, model.MachineKeySearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &members)
	if err != nil {
		return nil, 0, err
	}
	return members, count, nil
}

func MachineKeysByUserID(db *gorm.DB, table string, userID string) ([]*model.MachineKeyView, error) {
	keys := make([]*model.MachineKeyView, 0)
	queries := []*usr_model.MachineKeySearchQuery{
		{
			Key:    usr_model.MachineKeyKeyUserID,
			Value:  userID,
			Method: global_model.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.MachineKeySearchRequest{Queries: queries})
	_, err := query(db, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func PutMachineKey(db *gorm.DB, table string, role *model.MachineKeyView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func DeleteMachineKey(db *gorm.DB, table, keyID string) error {
	delete := repository.PrepareDeleteByKey(table, model.MachineKeySearchKey(usr_model.MachineKeyKeyID), keyID)
	return delete(db)
}

func DeleteMachineKeysByUserID(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.MachineKeySearchKey(usr_model.MachineKeyKeyUserID), userID)
	return delete(db)
}
