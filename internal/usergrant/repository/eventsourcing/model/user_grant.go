package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type UserGrant struct {
	es_models.ObjectRoot

	State     int32    `json:"-"`
	UserID    string   `json:"userId,omitempty"`
	ProjectID string   `json:"projectId,omitempty"`
	GrantID   string   `json:"grantId,omitempty"`
	RoleKeys  []string `json:"roleKeys,omitempty"`
}
