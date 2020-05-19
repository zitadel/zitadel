package model

import (
	"database/sql"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/key/model"
)

const (
	KeyKeyID     = "id"
	KeyPrivate   = "private"
	KeyUsage     = "usage"
	KeyAlgorithm = "algorithm"
	KeyExpiry    = "expiry"
)

type KeyView struct {
	ID        string              `json:"-" gorm:"column:key_id;primary_key"`
	Private   sql.NullBool        `json:"-" gorm:"column:private;primary_key"`
	Expiry    time.Time           `json:"-" gorm:"column:expiry"`
	Algorithm string              `json:"-" gorm:"column:algorithm"`
	Usage     int32               `json:"-" gorm:"column:usage"`
	Key       *crypto.CryptoValue `json:"key" gorm:"column:key"`
	Sequence  uint64              `json:"-" gorm:"column:sequence"`
}

func KeyViewsToModel(keys []*KeyView) []*model.KeyView {
	converted := make([]*model.KeyView, len(keys))
	for i, key := range keys {
		converted[i] = KeyViewToModel(key)
	}
	return converted
}

func KeyViewToModel(key *KeyView) *model.KeyView {
	return &model.KeyView{
		ID:        key.ID,
		Private:   key.Private.Bool,
		Expiry:    key.Expiry,
		Algorithm: key.Algorithm,
		Usage:     model.KeyUsage(key.Usage),
		Key:       key.Key,
		Sequence:  key.Sequence,
	}
}
