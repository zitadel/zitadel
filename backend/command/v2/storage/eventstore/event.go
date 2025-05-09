package eventstore

import "github.com/zitadel/zitadel/backend/command/v2/pattern"

type Event struct {
	AggregateType string `json:"aggregateType"`
	AggregateID   string `json:"aggregateId"`
}

type EventCommander interface {
	pattern.Command
	Event() *Event
}
