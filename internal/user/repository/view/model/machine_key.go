package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	MachineKeyKeyID     = "id"
	MachineKeyKeyUserID = "user_id"
)

type MachineKeyView struct {
	ID             string    `json:"keyId" gorm:"column:id;primary_key"`
	UserID         string    `json:"-" gorm:"column:user_id;primary_key"`
	Type           int32     `json:"type" gorm:"column:machine_type"`
	ExpirationDate time.Time `json:"expirationDate" gorm:"column:expiration_date"`
	Sequence       uint64    `json:"-" gorm:"column:sequence"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
}

func MachineKeyViewFromModel(key *model.MachineKeyView) *MachineKeyView {
	return &MachineKeyView{
		ID:             key.ID,
		UserID:         key.UserID,
		Type:           int32(key.Type),
		ExpirationDate: key.ExpirationDate,
		Sequence:       key.Sequence,
		CreationDate:   key.CreationDate,
	}
}

func MachineKeyToModel(key *MachineKeyView) *model.MachineKeyView {
	return &model.MachineKeyView{
		ID:             key.ID,
		UserID:         key.UserID,
		Type:           model.MachineKeyType(key.Type),
		ExpirationDate: key.ExpirationDate,
		Sequence:       key.Sequence,
		CreationDate:   key.CreationDate,
	}
}

func MachineKeysToModel(keys []*MachineKeyView) []*model.MachineKeyView {
	result := make([]*model.MachineKeyView, len(keys))
	for i, key := range keys {
		result[i] = MachineKeyToModel(key)
	}
	return result
}

func (k *MachineKeyView) AppendEvent(event *models.Event) (err error) {
	k.Sequence = event.Sequence
	switch event.Type {
	case es_model.MachineKeyAdded:
		k.setRootData(event)
		k.CreationDate = event.CreationDate
		err = k.SetData(event)
	}
	return err
}

func (k *MachineKeyView) setRootData(event *models.Event) {
	k.UserID = event.AggregateID
}

func (r *MachineKeyView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-Sj90d").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
	}
	return nil
}
