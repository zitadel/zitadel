package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     string
	IsPhoneVerified bool
}

func (p *Phone) IsValid() bool {
	return p.PhoneNumber != ""
}

func (u *User) appendUserPhoneChangedEvent(event *es_models.Event) error {
	u.Phone = new(Phone)
	u.Phone.setData(event)
	u.IsPhoneVerified = false
	return nil
}

func (u *User) appendUserPhoneVerifiedEvent() error {
	u.IsPhoneVerified = true
	return nil
}

func (a *Phone) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-dlo9s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
