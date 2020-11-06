package iam

import "github.com/caos/zitadel/internal/eventstore/v2"

type Aggregate struct {
	eventstore.Aggregate
}

type Step int8

type SetupStepEvent struct {
	eventstore.BaseEvent `json:"-"`

	Step Step
	//Done if the setup is started earlier
	Done bool `json:"-"`
}

func (e *SetupStepEvent) CheckPrevious() bool {
	return e.Type() == "iam.setup.started"
}

//Type implements event
func (e *SetupStepEvent) Type() eventstore.EventType {
	if e.Done {
		return "iam.setup.done"
	}
	return "iam.setup.started"
}

func (e *SetupStepEvent) Data() interface{} {
	return e
}

type MemberAddedEvent struct {
}
