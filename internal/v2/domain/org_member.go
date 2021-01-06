package domain

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type OrgMember struct {
	models.ObjectRoot

	UserID string
	Roles  []string
}
