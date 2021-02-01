package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/key/model"
	proj_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	user_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	AuthNKeyKeyID      = "key_id"
	AuthNKeyObjectID   = "object_id"
	AuthNKeyObjectType = "object_type"
)

type AuthNKeyView struct {
	ID             string    `json:"keyId" gorm:"column:key_id;primary_key"`
	ObjectID       string    `json:"-" gorm:"column:object_id;primary_key"`
	ObjectType     int32     `json:"-" gorm:"column:object_type;primary_key"`
	AuthIdentifier string    `json:"-" gorm:"column:auth_identifier;primary_key"`
	Type           int32     `json:"type" gorm:"column:key_type"`
	ExpirationDate time.Time `json:"expirationDate" gorm:"column:expiration_date"`
	Sequence       uint64    `json:"-" gorm:"column:sequence"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`

	PublicKey []byte `json:"publicKey" gorm:"column:public_key"`
}

func AuthNKeyViewFromModel(key *model.AuthNKeyView) *AuthNKeyView {
	return &AuthNKeyView{
		ID:             key.ID,
		ObjectID:       key.ObjectID,
		ObjectType:     int32(key.ObjectType),
		Type:           int32(key.Type),
		ExpirationDate: key.ExpirationDate,
		Sequence:       key.Sequence,
		CreationDate:   key.CreationDate,
	}
}

func AuthNKeyToModel(key *AuthNKeyView) *model.AuthNKeyView {
	return &model.AuthNKeyView{
		ID:             key.ID,
		ObjectID:       key.ObjectID,
		ObjectType:     model.ObjectType(key.ObjectType),
		AuthIdentifier: key.AuthIdentifier,
		Type:           model.AuthNKeyType(key.Type),
		ExpirationDate: key.ExpirationDate,
		Sequence:       key.Sequence,
		CreationDate:   key.CreationDate,
		PublicKey:      key.PublicKey,
	}
}

func AuthNKeysToModel(keys []*AuthNKeyView) []*model.AuthNKeyView {
	result := make([]*model.AuthNKeyView, len(keys))
	for i, key := range keys {
		result[i] = AuthNKeyToModel(key)
	}
	return result
}

func (k *AuthNKeyView) AppendEvent(event *models.Event) (err error) {
	k.Sequence = event.Sequence
	switch event.Type {
	case user_model.MachineKeyAdded:
		k.setRootData(event)
		k.CreationDate = event.CreationDate
		err = k.SetUserData(event)
	case proj_model.ClientKeyAdded:
		k.setRootData(event)
		k.CreationDate = event.CreationDate
		err = k.SetClientData(event)
	}
	return err
}

func (k *AuthNKeyView) setRootData(event *models.Event) {
	switch event.AggregateType {
	case user_model.UserAggregate:
		k.ObjectType = int32(model.AuthNKeyObjectTypeUser)
		k.ObjectID = event.AggregateID
		k.AuthIdentifier = event.AggregateID
	case proj_model.ProjectAggregate:
		k.ObjectType = int32(model.AuthNKeyObjectTypeApplication)
	}
}

func (r *AuthNKeyView) SetUserData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-Sj90d").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
	}
	return nil
}

func (r *AuthNKeyView) SetClientData(event *models.Event) error {
	key := new(proj_model.ClientKey)
	if err := json.Unmarshal(event.Data, key); err != nil {
		logging.Log("EVEN-Dgsgg").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-ADbfz", "Could not unmarshal data")
	}
	r.ObjectID = key.ApplicationID
	r.AuthIdentifier = key.ClientID
	r.ID = key.KeyID
	r.ExpirationDate = key.ExpirationDate
	r.PublicKey = key.PublicKey
	r.Type = key.Type
	return nil
}
