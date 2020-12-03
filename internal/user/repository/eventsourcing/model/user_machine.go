package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

type Machine struct {
	user *User `json:"-"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (sa *Machine) AppendEvents(events ...*models.Event) error {
	for _, event := range events {
		if err := sa.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (sa *Machine) AppendEvent(event *models.Event) (err error) {
	switch event.Type {
	case MachineAdded, MachineChanged:
		err = sa.setData(event)
	}

	return err
}

func (sa *Machine) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, sa); err != nil {
		logging.Log("EVEN-8ujgd").WithError(err).Error("could not unmarshal event data")
		return errors.ThrowInternal(err, "MODEL-GwjY9", "could not unmarshal event")
	}
	return nil
}

func (sa *Machine) Changes(updatedAccount *Machine) map[string]interface{} {
	changes := make(map[string]interface{})
	if updatedAccount.Description != "" && updatedAccount.Description != sa.Description {
		changes["description"] = updatedAccount.Description
	}
	return changes
}

func MachineFromModel(machine *model.Machine) *Machine {
	return &Machine{
		Description: machine.Description,
		Name:        machine.Name,
	}
}

func MachineToModel(machine *Machine) *model.Machine {
	return &model.Machine{
		Description: machine.Description,
		Name:        machine.Name,
	}
}

type MachineKey struct {
	es_models.ObjectRoot `json:"-"`
	KeyID                string    `json:"keyId,omitempty"`
	Type                 int32     `json:"type,omitempty"`
	ExpirationDate       time.Time `json:"expirationDate,omitempty"`
	PublicKey            []byte    `json:"publicKey,omitempty"`
	privateKey           []byte
}

func (key *MachineKey) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		err := key.AppendEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (key *MachineKey) AppendEvent(event *es_models.Event) error {
	switch event.Type {
	case MachineKeyAdded:
	}
	return nil
}

func MachineKeyFromModel(machine *model.MachineKey) *MachineKey {
	return &MachineKey{
		ObjectRoot:     machine.ObjectRoot,
		ExpirationDate: machine.ExpirationDate,
		KeyID:          machine.KeyID,
		Type:           int32(machine.Type),
	}
}

func MachineKeyToModel(machine *MachineKey) *model.MachineKey {
	return &model.MachineKey{
		ObjectRoot:     machine.ObjectRoot,
		ExpirationDate: machine.ExpirationDate,
		KeyID:          machine.KeyID,
		PrivateKey:     machine.privateKey,
		Type:           model.MachineKeyType(machine.Type),
	}
}

func (key *MachineKey) GenerateMachineKeyPair(keySize int, alg crypto.EncryptionAlgorithm) error {
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
