package model

import (
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Group struct {
	es_models.ObjectRoot

	Name        string
	Description string
	State       GroupState
	Members     []*GroupMember
}

type GroupState int32

const (
	GroupStateActive GroupState = iota
	GroupStateInactive
	GroupStateRemoved
)
