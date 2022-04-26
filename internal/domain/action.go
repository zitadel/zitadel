package domain

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Action struct {
	models.ObjectRoot

	Name          string
	Script        string
	Timeout       time.Duration
	AllowedToFail bool
	State         ActionState
}

func (a *Action) IsValid() bool {
	return a.Name != ""
}

type ActionState int32

const (
	ActionStateUnspecified ActionState = iota
	ActionStateActive
	ActionStateInactive
	ActionStateRemoved
	actionStateCount
)

func (s ActionState) Valid() bool {
	return s >= 0 && s < actionStateCount
}

func (s ActionState) Exists() bool {
	return s != ActionStateUnspecified && s != ActionStateRemoved
}

type ActionsAllowed int32

const (
	ActionsNotAllowed ActionsAllowed = iota
	ActionsMaxAllowed
	ActionsAllowedUnlimited
)
