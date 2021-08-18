package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Project struct {
	models.ObjectRoot

	State                ProjectState
	Name                 string
	ProjectRoleAssertion bool
	ProjectRoleCheck     bool
	HasProjectCheck      bool
	PrivateLabeling      PrivateLabelingSetting
}

type ProjectState int32

const (
	ProjectStateUnspecified ProjectState = iota
	ProjectStateActive
	ProjectStateInactive
	ProjectStateRemoved
)

type PrivateLabelingSetting int32

const (
	PrivateLabelingSettingUnspecified PrivateLabelingSetting = iota
	PrivateLabelingSettingProjectResourceOwnerSetting
	PrivateLabelingSettingLoginUserResourceOwnerSetting
)

func (o *Project) IsValid() bool {
	return o.Name != ""
}
