package model

import (
	"encoding/json"

	"github.com/zitadel/logging"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ExternalIDP struct {
	es_models.ObjectRoot

	IDPConfigID string `json:"idpConfigId,omitempty"`
	UserID      string `json:"userId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func GetExternalIDP(idps []*ExternalIDP, id string) (int, *ExternalIDP) {
	for i, idp := range idps {
		if idp.UserID == id {
			return i, idp
		}
	}
	return -1, nil
}
func (u *Human) appendExternalIDPAddedEvent(event *es_models.Event) error {
	idp := new(ExternalIDP)
	err := idp.setData(event)
	if err != nil {
		return err
	}
	idp.ObjectRoot.CreationDate = event.CreationDate
	u.ExternalIDPs = append(u.ExternalIDPs, idp)
	return nil
}

func (u *Human) appendExternalIDPRemovedEvent(event *es_models.Event) error {
	idp := new(ExternalIDP)
	err := idp.setData(event)
	if err != nil {
		return err
	}
	if i, externalIdp := GetExternalIDP(u.ExternalIDPs, idp.UserID); externalIdp != nil {
		u.ExternalIDPs[i] = u.ExternalIDPs[len(u.ExternalIDPs)-1]
		u.ExternalIDPs[len(u.ExternalIDPs)-1] = nil
		u.ExternalIDPs = u.ExternalIDPs[:len(u.ExternalIDPs)-1]
	}
	return nil
}

func (pw *ExternalIDP) setData(event *es_models.Event) error {
	pw.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, pw); err != nil {
		logging.Log("EVEN-Msi9d").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-A9osf", "could not unmarshal event")
	}
	return nil
}
