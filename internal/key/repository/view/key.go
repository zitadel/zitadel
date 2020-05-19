package view

import (
	"time"

	"github.com/jinzhu/gorm"

	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view"
)

func KeyByIDAndType(db *gorm.DB, table, keyID string, private bool) (*model.KeyView, error) {
	key := new(model.KeyView)
	query := view.PrepareGetByQuery(table,
		model.KeySearchQuery{Key: key_model.KEYSEARCHKEY_ID, Method: global_model.SEARCHMETHOD_EQUALS, Value: keyID},
		model.KeySearchQuery{Key: key_model.KEYSEARCHKEY_PRIVATE, Method: global_model.SEARCHMETHOD_EQUALS, Value: private},
	)
	err := query(db, key)
	return key, err
}

func GetSigningKey(db *gorm.DB, table string) (*model.KeyView, error) {
	key := new(model.KeyView)
	query := view.PrepareGetByQuery(table,
		model.KeySearchQuery{Key: key_model.KEYSEARCHKEY_PRIVATE, Method: global_model.SEARCHMETHOD_EQUALS, Value: true},
		model.KeySearchQuery{Key: key_model.KEYSEARCHKEY_USAGE, Method: global_model.SEARCHMETHOD_EQUALS, Value: key_model.KeyUsageSigning},
		model.KeySearchQuery{Key: key_model.KEYSEARCHKEY_EXPIRY, Method: global_model.SEARCHMETHOD_GREATER_THAN, Value: time.Now().UTC()},
	)
	err := query(db, key)
	return key, err
}

func GetActivePublicKeys(db *gorm.DB, table string) ([]*model.KeyView, error) {
	keys := make([]*model.KeyView, 0)
	query := view.PrepareSearchQuery(table,
		model.KeySearchRequest{
			Queries: []*key_model.KeySearchQuery{
				{Key: key_model.KEYSEARCHKEY_PRIVATE, Method: global_model.SEARCHMETHOD_EQUALS, Value: false},
				{Key: key_model.KEYSEARCHKEY_USAGE, Method: global_model.SEARCHMETHOD_EQUALS, Value: key_model.KeyUsageSigning},
				{Key: key_model.KEYSEARCHKEY_EXPIRY, Method: global_model.SEARCHMETHOD_GREATER_THAN, Value: time.Now().UTC()},
			},
		},
	)
	_, err := query(db, keys)
	return keys, err
}

func PutKey(db *gorm.DB, table string, key *model.KeyView) error {
	save := view.PrepareSave(table)
	return save(db, key)
}

func DeleteKey(db *gorm.DB, table, keyID string, private bool) error {
	delete := view.PrepareDeleteByKeys(table,
		view.Key{Key: model.KeySearchKey(key_model.KEYSEARCHKEY_ID), Value: keyID},
		view.Key{Key: model.KeySearchKey(key_model.KEYSEARCHKEY_PRIVATE), Value: private},
	)
	return delete(db)
}

func DeleteKeyPair(db *gorm.DB, table, keyID string) error {
	delete := view.PrepareDeleteByKey(table, model.KeySearchKey(key_model.KEYSEARCHKEY_ID), keyID)
	return delete(db)
}
