package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/key/model"
)

const (
	KeyPairVersion = "v1"
)

type KeyPair struct {
	es_models.ObjectRoot

	Usage      int32  `json:"usage"`
	Algorithm  string `json:"algorithm"`
	PrivateKey *Key   `json:"privateKey"`
	PublicKey  *Key   `json:"publicKey"`
}

type Key struct {
	Key    *crypto.CryptoValue `json:"key"`
	Expiry time.Time           `json:"expiry"`
}

func KeyPairFromModel(pair *model.KeyPair) *KeyPair {
	return &KeyPair{
		ObjectRoot: pair.ObjectRoot,
		Usage:      int32(pair.Usage),
		Algorithm:  pair.Algorithm,
		PrivateKey: KeyFromModel(pair.PrivateKey),
		PublicKey:  KeyFromModel(pair.PublicKey),
	}
}

func KeyPairToModel(pair *KeyPair) *model.KeyPair {
	return &model.KeyPair{
		ObjectRoot: pair.ObjectRoot,
		Usage:      model.KeyUsage(pair.Usage),
		Algorithm:  pair.Algorithm,
		PrivateKey: KeyToModel(pair.PrivateKey),
		PublicKey:  KeyToModel(pair.PublicKey),
	}
}

func KeyFromModel(key *model.Key) *Key {
	return &Key{
		Key:    key.Key,
		Expiry: key.Expiry,
	}
}

func KeyToModel(key *Key) *model.Key {
	return &model.Key{
		Key:    key.Key,
		Expiry: key.Expiry,
	}
}

func (k *KeyPair) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := k.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (k *KeyPair) AppendEvent(event *es_models.Event) error {
	k.ObjectRoot.AppendEvent(event)
	switch event.Type {
	case KeyPairAdded:
		return k.AppendAddKeyPair(event)
	}
	return nil
}

func (k *KeyPair) AppendAddKeyPair(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, k); err != nil {
		logging.Log("EVEN-Je92s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
