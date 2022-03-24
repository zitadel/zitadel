package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type DomainPolicy struct {
	models.ObjectRoot

	UserLoginMustBeDomain bool
	Default               bool
}
