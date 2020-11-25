package idp

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

type ConfigWriteModel struct {
	eventstore.WriteModel

	ConfigID    string
	Name        string
	StylingType StylingType
	OIDCConfig  *oidc.ConfigWriteModel
}

func (rm *ConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
	rm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {
		case *oidc.ConfigAddedEvent:
			rm.OIDCConfig = &oidc.ConfigWriteModel{}
			rm.OIDCConfig.AppendEvents(event)
		case *oidc.ConfigChangedEvent:
			rm.OIDCConfig.AppendEvents(event)
		}
	}
}

func (rm *ConfigWriteModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *ConfigAddedEvent:
			rm.ConfigID = e.ConfigID
			rm.Name = e.Name
			rm.StylingType = e.StylingType
		case *ConfigChangedEvent:
			if e.Name != "" {
				rm.Name = e.Name
			}
			if e.StylingType.Valid() {
				rm.StylingType = e.StylingType
			}
		}
	}
	if err := rm.OIDCConfig.Reduce(); err != nil {
		return err
	}
	return rm.WriteModel.Reduce()
}
