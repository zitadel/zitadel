package domain

import (
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	maxGroupNameLen = 200
)

// Group represents a user group in an organization
type Group struct {
	models.ObjectRoot

	Name        string
	Description string
}

func (g *Group) IsValid() error {
	groupName := strings.TrimSpace(g.Name)
	if groupName == "" || len(groupName) > maxGroupNameLen {
		return zerrors.ThrowInvalidArgument(nil, "GROUP-m177lN", "Errors.Group.InvalidName")
	}
	return nil
}

type GroupState int32

const (
	GroupStateUnspecified GroupState = iota
	GroupStateActive
	GroupStateRemoved
)

func (g GroupState) Exists() bool {
	return g != GroupStateRemoved && g != GroupStateUnspecified
}
