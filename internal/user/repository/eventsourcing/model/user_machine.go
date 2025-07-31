package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Machine struct {
	user *User `json:"-"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (sa *Machine) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := sa.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (sa *Machine) AppendEvent(event *es_models.Event) (err error) {
	switch event.Type() {
	case user_repo.MachineAddedEventType, user_repo.MachineChangedEventType:
		err = sa.setData(event)
	}

	return err
}

func (sa *Machine) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, sa); err != nil {
		logging.Log("EVEN-8ujgd").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-GwjY9", "could not unmarshal event")
	}
	return nil
}

type MachineKey struct {
	es_models.ObjectRoot `json:"-"`
	KeyID                string    `json:"keyId,omitempty"`
	Type                 int32     `json:"type,omitempty"`
	ExpirationDate       time.Time `json:"expirationDate"`
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

func (key *MachineKey) AppendEvent(event *es_models.Event) (err error) {
	key.ObjectRoot.AppendEvent(event)
	switch event.Type() {
	case user_repo.MachineKeyAddedEventType:
		err = json.Unmarshal(event.Data, key)
		if err != nil {
			return zerrors.ThrowInternal(err, "MODEL-SjI4S", "Errors.Internal")
		}
	case user_repo.MachineKeyRemovedEventType:
		key.ExpirationDate = event.CreationDate
	}
	return err
}
