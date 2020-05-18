package view

import (
	"time"

	"github.com/jinzhu/gorm"

	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view"
)

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
