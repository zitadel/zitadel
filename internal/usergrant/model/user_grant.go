package model

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type UserGrant struct {
	es_models.ObjectRoot

	State     UserGrantState
	UserID    string
	ProjectID string
	RoleKeys  []string
}

type UserGrantState int32

const (
	USERGRANTSTATE_ACTIVE UserGrantState = iota
	USERGRANTSTATE_INACTIVE
	USERGRANTSTATE_REMOVED
)

func (u *UserGrant) IsValid() bool {
	return u.ProjectID != "" && u.UserID != ""
}

func (u *UserGrant) IsActive() bool {
	return u.State == USERGRANTSTATE_ACTIVE
}

func (u *UserGrant) IsInactive() bool {
	return u.State == USERGRANTSTATE_INACTIVE
}

func (u *UserGrant) RemoveRoleKeyIfExisting(key string) bool {
	for i, role := range u.RoleKeys {
		if role == key {
			u.RoleKeys[i] = u.RoleKeys[len(u.RoleKeys)-1]
			u.RoleKeys[len(u.RoleKeys)-1] = ""
			u.RoleKeys = u.RoleKeys[:len(u.RoleKeys)-1]
			return true
		}
	}
	return false
}

func (u *UserGrant) RemoveRoleKeysIfExisting(keys []string) bool {
	exists := false
	for _, key := range keys {
		keyExists := u.RemoveRoleKeyIfExisting(key)
		if keyExists {
			exists = true
		}
	}
	return exists
}
