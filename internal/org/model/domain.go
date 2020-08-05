package model

import (
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type OrgDomain struct {
	es_models.ObjectRoot
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

func (t OrgDomainValidationType) CheckType() http_util.CheckType {
	switch t {
	case OrgDomainValidationTypeHTTP:
		return http_util.CheckTypeHTTP
	}
}

func (t OrgDomainValidationType) IsDNS() bool {
	return t == OrgDomainValidationTypeDNS
}

func NewOrgDomain(orgID, domain string) *OrgDomain {
	return &OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: orgID}, Domain: domain}
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
