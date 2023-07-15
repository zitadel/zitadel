package domain

import (
	"regexp"
	"strings"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
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
	return domain.Domain != ""
}

func (domain *OrgDomain) GenerateVerificationCode(codeGenerator crypto.Generator) (string, error) {
	validationCodeCrypto, validationCode, err := crypto.NewCode(codeGenerator)
	if err != nil {
		return "", err
	}
	domain.ValidationCode = validationCodeCrypto
	return validationCode, nil
}

func NewIAMDomainName(orgName, iamDomain string) string {
	// Reference: label domain requirements https://www.nic.ad.jp/timeline/en/20th/appendix1.html

	// Replaces spaces in org name with hyphens
	label := strings.ReplaceAll(orgName, " ", "-")

	// The label must only contains alphanumeric characters and hyphens
	// Invalid characters are replaced with and empty space
	label = string(regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAll([]byte(label), []byte("")))

	// The label cannot exceed 63 characters
	if len(label) > 63 {
		label = label[:63]
	}

	// The total length of the resulting domain can't exceed 253 characters
	domain := label + "." + iamDomain
	if len(domain) > 253 {
		truncateNChars := len(domain) - 253
		label = label[:len(label)-truncateNChars]
	}

	// Label (maybe truncated) can't start with a hyphen
	if len(label) > 0 && label[0:1] == "-" {
		label = label[1:]
	}

	// Label (maybe truncated) can't end with a hyphen
	if len(label) > 0 && label[len(label)-1:] == "-" {
		label = label[:len(label)-1]
	}

	return strings.ToLower(label + "." + iamDomain)
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
