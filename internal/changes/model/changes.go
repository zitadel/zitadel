package model

import (
	es_type "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type ChangeObjectType int32

const (
	ChangeObjectTypeUser ChangeObjectType = iota
	ChangeObjectTypeApp
	ChangeObjectTypeOrg
	ChangeObjectTypeProject

	Project      es_type.AggregateType = "project"
	Application  es_type.AggregateType = "application"
	ProjectGrant es_type.AggregateType = "project_grant"

	Org             es_type.AggregateType = "org"
	UniqueOrgName   es_type.AggregateType = "unique.org.name"
	UniqueOrgDomain es_type.AggregateType = "unique.org.domain"

	User           es_type.AggregateType = "user"
	UniqueUsername es_type.AggregateType = "unique.user.username"
	UniqueEmail    es_type.AggregateType = "unique.user.email"
	UserGrant      es_type.AggregateType = "user_grant"
)

type Changes struct {
	Changes      []*Change
	LastSequence uint64
}

type Change struct {
	ChangeDate *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType  string               `json:"eventType,omitempty"`
	Sequence   uint64               `json:"sequence,omitempty"`
	Modifier   string               `json:"modifierUser,omitempty"`
	Data       interface{}          `json:"data,omitempty"`
}
