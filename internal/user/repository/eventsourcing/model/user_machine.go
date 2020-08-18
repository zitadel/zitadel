package model

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

type Machine struct {
	objectRoot models.ObjectRoot
	state      int32 `json:"-"`

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
		sa.state = int32(model.UserStateActive)
		err = sa.setData(event)
	case KeyAdded:
		fallthrough
	case KeyRemoved:
		logging.Log("MODEL-iBgOc").Panic("key unimplemented")
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

func ServiceAccountFromModel(serviceAccount *model.Machine) *Machine {
	return &Machine{
		Description: serviceAccount.Description,
		Name:        serviceAccount.Name,
	}
}

func ServiceAccountToModel(serviceAccount *Machine) *model.Machine {
	return &model.Machine{
		Description: serviceAccount.Description,
		Name:        serviceAccount.Name,
	}
}
