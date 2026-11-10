package domain

import (
	"slices"
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

type ActionFunction int32

const (
	ActionFunctionUnspecified ActionFunction = iota
	ActionFunctionPreUserinfo
	ActionFunctionPreAccessToken
	ActionFunctionPreSAMLResponse
	actionFunctionCount
)

func (s ActionFunction) Valid() bool {
	return s >= 0 && s < actionFunctionCount
}

func (s ActionFunction) LocalizationKey() string {
	if !s.Valid() {
		return ActionFunctionUnspecified.LocalizationKey()
	}

	switch s {
	case ActionFunctionPreUserinfo:
		return "preuserinfo"
	case ActionFunctionPreAccessToken:
		return "preaccesstoken"
	case ActionFunctionPreSAMLResponse:
		return "presamlresponse"
	case ActionFunctionUnspecified, actionFunctionCount:
		fallthrough
	default:
		return "unspecified"
	}
}

func AllActionFunctions() []string {
	return []string{
		ActionFunctionPreUserinfo.LocalizationKey(),
		ActionFunctionPreAccessToken.LocalizationKey(),
		ActionFunctionPreSAMLResponse.LocalizationKey(),
	}
}

func ActionFunctionExists() func(string) bool {
	functions := AllActionFunctions()
	return func(s string) bool {
		return slices.Contains(functions, s)
	}
}
