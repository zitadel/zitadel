package domain

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type OrgIAMPolicy struct {
	models.ObjectRoot

	UserLoginMustBeDomain bool
	Default               bool
}
