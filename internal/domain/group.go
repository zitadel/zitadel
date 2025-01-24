package domain

import (
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Group struct {
	models.ObjectRoot

	State       GroupState
	Name        string
	Description string
}

type GroupState int32

const (
	GroupStateUnspecified GroupState = iota
	GroupStateActive
	GroupStateInactive
	GroupStateRemoved

	groupStateMax
)

func (s GroupState) Valid() bool {
	return s > GroupStateUnspecified && s < groupStateMax
}

func (o *Group) IsValid() bool {
	return o.Name != ""
}
