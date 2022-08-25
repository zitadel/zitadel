package model

import (
	"encoding/json"

	"github.com/zitadel/logging"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type SAMLConfig struct {
	es_models.ObjectRoot
	AppID       string `json:"appId"`
	Metadata    []byte `json:"metadata,omitempty"`
	MetadataURL string `json:"metadataUrl,omitempty"`
}

func (o *SAMLConfig) setData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-d8e3s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
