package idp

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type ConfigAggregate struct {
	eventstore.Aggregate
}

type ConfigType int32

const (
	ConfigTypeOIDC ConfigType = iota
	ConfigTypeSAML

	//count is for validation
	configTypeCount
)

func (f ConfigType) Valid() bool {
	return f >= 0 && f < configTypeCount
}

type ConfigState int32

const (
	ConfigStateActive ConfigState = iota
	ConfigStateInactive
	ConfigStateRemoved

	configStateCount
)

func (f ConfigState) Valid() bool {
	return f >= 0 && f < configStateCount
}

type StylingType int32

const (
	StylingTypeGoogle StylingType = iota + 1

	stylingTypeCount
)

func (f StylingType) Valid() bool {
	return f >= 0 && f < stylingTypeCount
}
