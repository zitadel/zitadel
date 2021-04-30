package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

type IDPConfigsReadModel struct {
	eventstore.ReadModel

	Configs []*IDPConfigReadModel
}

func (rm *IDPConfigsReadModel) ConfigByID(id string) (idx int, config *IDPConfigReadModel) {
	for idx, config = range rm.Configs {
		if config.ConfigID == id {
			return idx, config
		}
	}
	return -1, nil
}

func (rm *IDPConfigsReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *idpconfig.IDPConfigAddedEvent:
			config := NewIDPConfigReadModel(e.ConfigID)
			rm.Configs = append(rm.Configs, config)
			config.AppendEvents(event)
		case *idpconfig.IDPConfigChangedEvent:
			_, config := rm.ConfigByID(e.ConfigID)
			config.AppendEvents(e)
		case *idpconfig.IDPConfigDeactivatedEvent:
			_, config := rm.ConfigByID(e.ConfigID)
			config.AppendEvents(e)
		case *idpconfig.IDPConfigReactivatedEvent:
			_, config := rm.ConfigByID(e.ConfigID)
			config.AppendEvents(e)
		case *idpconfig.OIDCConfigAddedEvent:
			_, config := rm.ConfigByID(e.IDPConfigID)
			config.AppendEvents(e)
		case *idpconfig.OIDCConfigChangedEvent:
			_, config := rm.ConfigByID(e.IDPConfigID)
			config.AppendEvents(e)
		case *idpconfig.IDPConfigRemovedEvent:
			idx, _ := rm.ConfigByID(e.ConfigID)
			if idx < 0 {
				continue
			}
			copy(rm.Configs[idx:], rm.Configs[idx+1:])
			rm.Configs[len(rm.Configs)-1] = nil
			rm.Configs = rm.Configs[:len(rm.Configs)-1]
		}
	}
}

func (rm *IDPConfigsReadModel) Reduce() error {
	for _, config := range rm.Configs {
		if err := config.Reduce(); err != nil {
			return err
		}
	}
	return nil
}
