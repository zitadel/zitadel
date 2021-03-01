package view

import (
	"github.com/caos/zitadel/internal/domain"
	"time"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/view/repository"

	"github.com/jinzhu/gorm"

	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/view/model"
)

func KeyByIDAndType(db *gorm.DB, table, keyID string, private bool) (*model.KeyView, error) {
	key := new(model.KeyView)
	query := repository.PrepareGetByQuery(table,
		model.KeySearchQuery{Key: key_model.KeySearchKeyID, Method: domain.SearchMethodEquals, Value: keyID},
		model.KeySearchQuery{Key: key_model.KeySearchKeyPrivate, Method: domain.SearchMethodEquals, Value: private},
	)
	err := query(db, key)
	return key, err
}

func GetSigningKey(db *gorm.DB, table string, expiry time.Time) (*model.KeyView, error) {
	if expiry.IsZero() {
		expiry = time.Now().UTC()
	}
	keys := make([]*model.KeyView, 0)
	query := repository.PrepareSearchQuery(table,
		model.KeySearchRequest{
			Queries: []*key_model.KeySearchQuery{
				{Key: key_model.KeySearchKeyPrivate, Method: domain.SearchMethodEquals, Value: true},
				{Key: key_model.KeySearchKeyUsage, Method: domain.SearchMethodEquals, Value: key_model.KeyUsageSigning},
				{Key: key_model.KeySearchKeyExpiry, Method: domain.SearchMethodGreaterThan, Value: time.Now().UTC()},
			},
			SortingColumn: key_model.KeySearchKeyExpiry,
			Limit:         1,
		},
	)
	_, err := query(db, &keys)
	if err != nil {
		return nil, err
	}
	if len(keys) != 1 {
		return nil, caos_errs.ThrowNotFound(err, "VIEW-BGD41", "key not found")
	}
	return keys[0], nil
}

func GetActivePublicKeys(db *gorm.DB, table string) ([]*model.KeyView, error) {
	keys := make([]*model.KeyView, 0)
	query := repository.PrepareSearchQuery(table,
		model.KeySearchRequest{
			Queries: []*key_model.KeySearchQuery{
				{Key: key_model.KeySearchKeyPrivate, Method: domain.SearchMethodEquals, Value: false},
				{Key: key_model.KeySearchKeyUsage, Method: domain.SearchMethodEquals, Value: key_model.KeyUsageSigning},
				{Key: key_model.KeySearchKeyExpiry, Method: domain.SearchMethodGreaterThan, Value: time.Now().UTC()},
			},
		},
	)
	_, err := query(db, &keys)
	return keys, err
}

func PutKeys(db *gorm.DB, table string, privateKey, publicKey *model.KeyView) error {
	save := repository.PrepareBulkSave(table)
	return save(db, privateKey, publicKey)
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
