package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddIdpConfigEvent(event *es_models.Event) error {
	idp := new(iam_es_model.IdpConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}
	idp.ObjectRoot.CreationDate = event.CreationDate
	o.IDPs = append(o.IDPs, idp)
	return nil
}

func (o *Org) appendChangeIdpConfigEvent(event *es_models.Event) error {
	idp := new(iam_es_model.IdpConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}
	if i, a := iam_es_model.GetIdpConfig(o.IDPs, idp.IDPConfigID); a != nil {
		o.IDPs[i].SetData(event)
	}
	return nil
}

func (o *Org) appendRemoveIdpConfigEvent(event *es_models.Event) error {
	idp := new(iam_es_model.IdpConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}
	if i, a := iam_es_model.GetIdpConfig(o.IDPs, idp.IDPConfigID); a != nil {
		o.IDPs[i] = o.IDPs[len(o.IDPs)-1]
		o.IDPs[len(o.IDPs)-1] = nil
		o.IDPs = o.IDPs[:len(o.IDPs)-1]
	}
	return nil
}

func (o *Org) appendIdpConfigStateEvent(event *es_models.Event, state model.IdpConfigState) error {
	idp := new(iam_es_model.IdpConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}

	if i, a := iam_es_model.GetIdpConfig(o.IDPs, idp.IDPConfigID); a != nil {
		a.State = int32(state)
		o.IDPs[i] = a
	}
	return nil
}

func (o *Org) appendAddOidcIdpConfigEvent(event *es_models.Event) error {
	config := new(iam_es_model.OidcIdpConfig)
	err := config.SetData(event)
	if err != nil {
		return err
	}
	config.ObjectRoot.CreationDate = event.CreationDate
	if i, a := iam_es_model.GetIdpConfig(o.IDPs, config.IdpConfigID); a != nil {
		o.IDPs[i].Type = int32(model.IDPConfigTypeOIDC)
		o.IDPs[i].OIDCIDPConfig = config
	}
	return nil
}

func (o *Org) appendChangeOidcIdpConfigEvent(event *es_models.Event) error {
	config := new(iam_es_model.OidcIdpConfig)
	err := config.SetData(event)
	if err != nil {
		return err
	}

	if i, a := iam_es_model.GetIdpConfig(o.IDPs, config.IdpConfigID); a != nil {
		o.IDPs[i].OIDCIDPConfig.SetData(event)
	}
	return nil
}
