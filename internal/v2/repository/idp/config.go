package idp

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

type ConfigAggregate struct {
	eventstore.Aggregate

	ConfigID    string
	Type        ConfigType
	Name        string
	StylingType StylingType
	State       ConfigState
	// OIDCConfig  *oidc.ConfigReadModel
}

type ConfigReadModel struct {
	eventstore.ReadModel

	ConfigID    string
	Type        ConfigType
	Name        string
	StylingType StylingType
	State       ConfigState
	OIDCConfig  *oidc.ConfigReadModel
}

func (rm *ConfigReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.ReadModel.AppendEvents(events...)
}

func (rm *ConfigReadModel) Reduce() error {
	return nil
}

type ConfigType int32

const (
	ConfigTypeOIDC ConfigType = iota
	ConfigTypeSAML
)

type ConfigState int32

const (
	ConfigStateActive ConfigState = iota
	ConfigStateInactive
	ConfigStateRemoved
)

type StylingType int32

const (
	StylingTypeUnspecified StylingType = iota
	StylingTypeGoogle
)
