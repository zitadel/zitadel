package domain

import (
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
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

func (k *ApplicationKey) SetPublicKey(publicKey []byte) {
	k.PublicKey = publicKey
}

func (k *ApplicationKey) SetPrivateKey(privateKey []byte) {
	k.PrivateKey = privateKey
}

func (k *ApplicationKey) GetExpirationDate() time.Time {
	return k.ExpirationDate
}

func (k *ApplicationKey) SetExpirationDate(expiration time.Time) {
	k.ExpirationDate = expiration
}

func (k *ApplicationKey) Detail() ([]byte, error) {
	if k.Type == AuthNKeyTypeJSON {
		return k.MarshalJSON()
	}
	return nil, zerrors.ThrowPreconditionFailed(nil, "KEY-dsg52", "Errors.Internal")
}

func (k *ApplicationKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string `json:"type"`
		KeyID    string `json:"keyId"`
		Key      string `json:"key"`
		AppID    string `json:"appId"`
		ClientID string `json:"clientId"`
	}{
		Type:     "application",
		KeyID:    k.KeyID,
		Key:      string(k.PrivateKey),
		AppID:    k.ApplicationID,
		ClientID: k.ClientID,
	})
}
