package domain

import (
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type DomainPolicy struct {
	models.ObjectRoot

	UserLoginMustBeDomain                  bool
	ValidateOrgDomains                     bool
	SMTPSenderAddressMatchesInstanceDomain bool
	Default                                bool
}
