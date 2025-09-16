package domain

import (
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

const (
	maxGroupNameLen = 200
)

// Group represents a user group in an organization
type Group struct {
	models.ObjectRoot

	Name           string
	Description    string
	OrganizationID string
}

func (g *Group) IsValid() bool {
	groupName := strings.TrimSpace(g.Name)
	return groupName != "" && len(groupName) <= maxGroupNameLen
}

type GroupState int32

const (
	GroupStateUnspecified GroupState = iota
	GroupStateActive
	GroupStateInactive
	GroupStateRemoved
)

func (g GroupState) Exists() bool {
	return g != GroupStateRemoved && g != GroupStateUnspecified
}
