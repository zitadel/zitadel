package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/key/model"
	es_model "github.com/caos/zitadel/internal/key/repository/eventsourcing/model"
)

const (
	KeyKeyID     = "id"
	KeyPrivate   = "private"
	KeyUsage     = "usage"
	KeyAlgorithm = "algorithm"
	KeyExpiry    = "expiry"
)

type KeyView struct {
	ID        string              `json:"-" gorm:"column:id;primary_key"`
	Private   sql.NullBool        `json:"-" gorm:"column:private;primary_key"`
	Expiry    time.Time           `json:"-" gorm:"column:expiry"`
	Algorithm string              `json:"-" gorm:"column:algorithm"`
	Usage     int32               `json:"-" gorm:"column:usage"`
	Key       *crypto.CryptoValue `json:"-" gorm:"column:key"`
	Sequence  uint64              `json:"-" gorm:"column:sequence"`
}

func KeysFromPairEvent(event *models.Event) (*KeyView, *KeyView, error) {
	pair := new(es_model.KeyPair)
	if err := json.Unmarshal(event.Data, pair); err != nil {
		logging.Log("MODEL-s3Ga1").WithError(err).Error("could not unmarshal event data")
		return nil, nil, caos_errs.ThrowInternal(nil, "MODEL-G3haa", "could not unmarshal data")
	}
	privateKey := &KeyView{
		ID:        event.AggregateID,
		Private:   sql.NullBool{Bool: true, Valid: true},
		Expiry:    pair.PrivateKey.Expiry,
		Algorithm: pair.Algorithm,
		Usage:     pair.Usage,
		Key:       pair.PrivateKey.Key,
		Sequence:  event.Sequence,
	}
	publicKey := &KeyView{
		ID:        event.AggregateID,
		Private:   sql.NullBool{Bool: false, Valid: true},
		Expiry:    pair.PublicKey.Expiry,
		Algorithm: pair.Algorithm,
		Usage:     pair.Usage,
		Key:       pair.PublicKey.Key,
		Sequence:  event.Sequence,
	}
	return privateKey, publicKey, nil
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

func (k *KeyView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, k); err != nil {
		logging.Log("MODEL-4ag41").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-GFQ31", "could not unmarshal data")
	}
	return nil
}
