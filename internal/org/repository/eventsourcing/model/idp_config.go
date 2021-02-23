package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddIDPConfigEvent(event *es_models.Event) error {
	idp := new(iam_es_model.IDPConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}
	idp.ObjectRoot.CreationDate = event.CreationDate
	o.IDPs = append(o.IDPs, idp)
	return nil
}

func (o *Org) appendChangeIDPConfigEvent(event *es_models.Event) error {
	idp := new(iam_es_model.IDPConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}
	if i, idpConfig := iam_es_model.GetIDPConfig(o.IDPs, idp.IDPConfigID); idpConfig != nil {
		o.IDPs[i].SetData(event)
	}
	return nil
}

func (o *Org) appendRemoveIDPConfigEvent(event *es_models.Event) error {
	idp := new(iam_es_model.IDPConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}
	if i, idpConfig := iam_es_model.GetIDPConfig(o.IDPs, idp.IDPConfigID); idpConfig != nil {
		o.IDPs[i] = o.IDPs[len(o.IDPs)-1]
		o.IDPs[len(o.IDPs)-1] = nil
		o.IDPs = o.IDPs[:len(o.IDPs)-1]
	}
	return nil
}

func (o *Org) appendIDPConfigStateEvent(event *es_models.Event, state model.IDPConfigState) error {
	idp := new(iam_es_model.IDPConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}

	if i, idpConfig := iam_es_model.GetIDPConfig(o.IDPs, idp.IDPConfigID); idpConfig != nil {
		idpConfig.State = int32(state)
		o.IDPs[i] = idpConfig
	}
	return nil
}

func (o *Org) appendAddOIDCIDPConfigEvent(event *es_models.Event) error {
	config := new(iam_es_model.OIDCIDPConfig)
	err := config.SetData(event)
	if err != nil {
		return err
	}
	config.ObjectRoot.CreationDate = event.CreationDate
	if i, idpConfig := iam_es_model.GetIDPConfig(o.IDPs, config.IDPConfigID); idpConfig != nil {
		o.IDPs[i].Type = int32(model.IDPConfigTypeOIDC)
		o.IDPs[i].OIDCIDPConfig = config
	}
	return nil
}

func (o *Org) appendChangeOIDCIDPConfigEvent(event *es_models.Event) error {
	config := new(iam_es_model.OIDCIDPConfig)
	err := config.SetData(event)
	if err != nil {
		return err
	}

	if i, idpConfig := iam_es_model.GetIDPConfig(o.IDPs, config.IDPConfigID); idpConfig != nil {
		o.IDPs[i].OIDCIDPConfig.SetData(event)
	}
	return nil
}
