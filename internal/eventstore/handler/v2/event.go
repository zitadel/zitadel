package handler

import "github.com/zitadel/zitadel/internal/eventstore"

type ProjectionSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	Name                 string `json:"name"`
}

func (p *ProjectionSucceededEvent) Data() interface{} {
	return p
}

func (p *ProjectionSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}
