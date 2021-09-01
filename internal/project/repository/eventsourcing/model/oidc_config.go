package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/project/model"
)

type OIDCConfig struct {
	es_models.ObjectRoot
	Version                  int32               `json:"oidcVersion,omitempty"`
	AppID                    string              `json:"appId"`
	ClientID                 string              `json:"clientId,omitempty"`
	ClientSecret             *crypto.CryptoValue `json:"clientSecret,omitempty"`
	RedirectUris             []string            `json:"redirectUris,omitempty"`
	ResponseTypes            []int32             `json:"responseTypes,omitempty"`
	GrantTypes               []int32             `json:"grantTypes,omitempty"`
	ApplicationType          int32               `json:"applicationType,omitempty"`
	AuthMethodType           int32               `json:"authMethodType,omitempty"`
	PostLogoutRedirectUris   []string            `json:"postLogoutRedirectUris,omitempty"`
	DevMode                  bool                `json:"devMode,omitempty"`
	AccessTokenType          int32               `json:"accessTokenType,omitempty"`
	AccessTokenRoleAssertion bool                `json:"accessTokenRoleAssertion,omitempty"`
	IDTokenRoleAssertion     bool                `json:"idTokenRoleAssertion,omitempty"`
	IDTokenUserinfoAssertion bool                `json:"idTokenUserinfoAssertion,omitempty"`
	ClockSkew                time.Duration       `json:"clockSkew,omitempty"`
	ClientKeys               []*ClientKey        `json:"-"`
}

func (o *OIDCConfig) setData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-d8e3s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

type ClientKey struct {
	es_models.ObjectRoot `json:"-"`
	ApplicationID        string    `json:"applicationId,omitempty"`
	ClientID             string    `json:"clientId,omitempty"`
	KeyID                string    `json:"keyId,omitempty"`
	Type                 int32     `json:"type,omitempty"`
	ExpirationDate       time.Time `json:"expirationDate,omitempty"`
	PublicKey            []byte    `json:"publicKey,omitempty"`
	privateKey           []byte
}

func (key *ClientKey) SetData(event *es_models.Event) error {
	key.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, key); err != nil {
		logging.Log("EVEN-SADdg").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (key *ClientKey) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		err := key.AppendEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (key *ClientKey) AppendEvent(event *es_models.Event) (err error) {
	key.ObjectRoot.AppendEvent(event)
	switch event.Type {
	case ClientKeyAdded:
		err = json.Unmarshal(event.Data, key)
		if err != nil {
			return errors.ThrowInternal(err, "MODEL-Fetg3", "Errors.Internal")
		}
	case ClientKeyRemoved:
		key.ExpirationDate = event.CreationDate
	}
	return err
}

func ClientKeysToModel(keys []*ClientKey) []*model.ClientKey {
	clientKeys := make([]*model.ClientKey, len(keys))
	for i, key := range keys {
		clientKeys[i] = ClientKeyToModel(key)
	}
	return clientKeys
}

func ClientKeyToModel(key *ClientKey) *model.ClientKey {
	return &model.ClientKey{
		ObjectRoot:     key.ObjectRoot,
		ExpirationDate: key.ExpirationDate,
		ApplicationID:  key.ApplicationID,
		ClientID:       key.ClientID,
		KeyID:          key.KeyID,
		PrivateKey:     key.privateKey,
		Type:           key_model.AuthNKeyType(key.Type),
	}
}

func (key *ClientKey) GenerateClientKeyPair(keySize int) error {
	privateKey, publicKey, err := crypto.GenerateKeyPair(keySize)
	if err != nil {
		return err
	}
	key.PublicKey, err = crypto.PublicKeyToBytes(publicKey)
	if err != nil {
		return err
	}
	key.privateKey = crypto.PrivateKeyToBytes(privateKey)
	return nil
}
