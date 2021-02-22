package domain

import (
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type ApplicationKey struct {
	models.ObjectRoot

	ApplicationID  string
	ClientID       string
	KeyID          string
	Type           AuthNKeyType
	ExpirationDate time.Time
	PrivateKey     []byte
	PublicKey      []byte
}

func (k *ApplicationKey) setPublicKey(publicKey []byte) {
	k.PublicKey = publicKey
}

func (k *ApplicationKey) setPrivateKey(privateKey []byte) {
	k.PrivateKey = privateKey
}

func (k *ApplicationKey) expirationDate() time.Time {
	return k.ExpirationDate
}

func (k *ApplicationKey) setExpirationDate(expiration time.Time) {
	k.ExpirationDate = expiration
}

func (k *ApplicationKey) Detail() ([]byte, error) {
	if k.Type == AuthNKeyTypeJSON {
		return k.MarshalJSON()
	}
	return nil, errors.ThrowPreconditionFailed(nil, "KEY-dsg52", "Errors.Internal")
}

func (k *ApplicationKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string `json:"type"`
		KeyID    string `json:"keyId"`
		Key      string `json:"key"`
		AppID    string `json:"appId"`
		ClientID string `json:"clientID"`
	}{
		Type:     "application",
		KeyID:    k.KeyID,
		Key:      string(k.PrivateKey),
		AppID:    k.ApplicationID,
		ClientID: k.ClientID,
	})
}
