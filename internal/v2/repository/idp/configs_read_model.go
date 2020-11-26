package idp

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

type ConfigsReadModel struct {
	eventstore.ReadModel

	Configs []*ConfigReadModel
}

func (rm *ConfigsReadModel) ConfigByID(id string) (idx int, config *ConfigReadModel) {
	for idx, config = range rm.Configs {
		if config.ConfigID == id {
			return idx, config
		}
	}
	return -1, nil
}

func (rm *ConfigsReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *ConfigAddedEvent:
			config := NewConfigReadModel(e.ConfigID)
			rm.Configs = append(rm.Configs, config)
			config.AppendEvents(event)
		case *ConfigChangedEvent:
			_, config := rm.ConfigByID(e.ConfigID)
			config.AppendEvents(e)
		case *ConfigDeactivatedEvent:
			_, config := rm.ConfigByID(e.ConfigID)
			config.AppendEvents(e)
		case *ConfigReactivatedEvent:
			_, config := rm.ConfigByID(e.ConfigID)
			config.AppendEvents(e)
		case *oidc.ConfigAddedEvent:
			_, config := rm.ConfigByID(e.IDPConfigID)
			config.AppendEvents(e)
		case *oidc.ConfigChangedEvent:
			_, config := rm.ConfigByID(e.IDPConfigID)
			config.AppendEvents(e)
		case *ConfigRemovedEvent:
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

func (rm *ConfigsReadModel) Reduce() error {
	for _, config := range rm.Configs {
		if err := config.Reduce(); err != nil {
			return err
		}
	}
	return nil
}
