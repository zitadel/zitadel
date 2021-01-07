package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type OrgDomain struct {
	models.ObjectRoot

	Domain         string
	Primary        bool
	Verified       bool
	ValidationType OrgDomainValidationType
	ValidationCode *crypto.CryptoValue
}

type OrgDomainValidationType int32

const (
	OrgDomainValidationTypeUnspecified OrgDomainValidationType = iota
	OrgDomainValidationTypeHTTP
	OrgDomainValidationTypeDNS
)

type OrgDomainState int32

const (
	OrgDomainStateUnspecified OrgDomainState = iota
	OrgDomainStateActive
	OrgDomainStateRemoved

	orgDomainStateCount
)

func (f OrgDomainState) Valid() bool {
	return f >= 0 && f < orgDomainStateCount
}
