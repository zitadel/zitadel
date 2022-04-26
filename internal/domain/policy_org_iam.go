package domain

import (
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type OrgIAMPolicy struct {
	models.ObjectRoot

	UserLoginMustBeDomain bool
	Default               bool
}
