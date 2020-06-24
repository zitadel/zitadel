package view

import (
	"github.com/caos/zitadel/internal/view/repository"
	"time"

	"github.com/jinzhu/gorm"

	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
)

func KeyByIDAndType(db *gorm.DB, table, keyID string, private bool) (*model.KeyView, error) {
	key := new(model.KeyView)
	query := repository.PrepareGetByQuery(table,
		model.KeySearchQuery{Key: key_model.KeySearchKeyID, Method: global_model.SearchMethodEquals, Value: keyID},
		model.KeySearchQuery{Key: key_model.KeySearchKeyPrivate, Method: global_model.SearchMethodEquals, Value: private},
	)
	err := query(db, key)
	return key, err
}

func GetSigningKey(db *gorm.DB, table string) (*model.KeyView, error) {
	key := new(model.KeyView)
	query := repository.PrepareGetByQuery(table,
		model.KeySearchQuery{Key: key_model.KeySearchKeyPrivate, Method: global_model.SearchMethodEquals, Value: true},
		model.KeySearchQuery{Key: key_model.KeySearchKeyUsage, Method: global_model.SearchMethodEquals, Value: key_model.KeyUsageSigning},
		model.KeySearchQuery{Key: key_model.KeySearchKeyExpiry, Method: global_model.SearchMethodGreaterThan, Value: time.Now().UTC()},
	)
	err := query(db, key)
	return key, err
}

func GetActivePublicKeys(db *gorm.DB, table string) ([]*model.KeyView, error) {
	keys := make([]*model.KeyView, 0)
	query := repository.PrepareSearchQuery(table,
		model.KeySearchRequest{
			Queries: []*key_model.KeySearchQuery{
				{Key: key_model.KeySearchKeyPrivate, Method: global_model.SearchMethodEquals, Value: false},
				{Key: key_model.KeySearchKeyUsage, Method: global_model.SearchMethodEquals, Value: key_model.KeyUsageSigning},
				{Key: key_model.KeySearchKeyExpiry, Method: global_model.SearchMethodGreaterThan, Value: time.Now().UTC()},
			},
		},
	)
	_, err := query(db, &keys)
	return keys, err
}

func PutKeys(db *gorm.DB, table string, privateKey, publicKey *model.KeyView) error {
	save := repository.PrepareSave(table)
	err := save(db, privateKey)
	if err != nil {
		return err
	}
	return save(db, publicKey)
}

func DeleteKey(db *gorm.DB, table, keyID string, private bool) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.KeySearchKey(key_model.KeySearchKeyID), Value: keyID},
		repository.Key{Key: model.KeySearchKey(key_model.KeySearchKeyPrivate), Value: private},
	)
	return delete(db)
}

func DeleteKeyPair(db *gorm.DB, table, keyID string) error {
	delete := repository.PrepareDeleteByKey(table, model.KeySearchKey(key_model.KeySearchKeyID), keyID)
	return delete(db)
}
