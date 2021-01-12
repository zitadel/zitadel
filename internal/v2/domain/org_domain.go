package domain

import (
	http_util "github.com/caos/zitadel/internal/api/http"
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

func (domain *OrgDomain) IsValid() bool {
	return domain.AggregateID != "" && domain.Domain != ""
}

func (domain *OrgDomain) GenerateVerificationCode(codeGenerator crypto.Generator) (string, error) {
	validationCodeCrypto, validationCode, err := crypto.NewCode(codeGenerator)
	if err != nil {
		return "", err
	}
	domain.ValidationCode = validationCodeCrypto
	return validationCode, nil
}

type OrgDomainValidationType int32

const (
	OrgDomainValidationTypeUnspecified OrgDomainValidationType = iota
	OrgDomainValidationTypeHTTP
	OrgDomainValidationTypeDNS
)

func (t OrgDomainValidationType) CheckType() (http_util.CheckType, bool) {
	switch t {
	case OrgDomainValidationTypeHTTP:
		return http_util.CheckTypeHTTP, true
	case OrgDomainValidationTypeDNS:
		return http_util.CheckTypeDNS, true
	default:
		return -1, false
	}
}

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
