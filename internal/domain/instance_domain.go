package domain

import (
	"strings"

	"github.com/zitadel/zitadel/internal/crypto"
)

var (
	domainRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type InstanceDomainState int32

const (
	InstanceDomainStateUnspecified InstanceDomainState = iota
	InstanceDomainStateActive
	InstanceDomainStateRemoved

	instanceDomainStateCount
)

func (f InstanceDomainState) Valid() bool {
	return f >= 0 && f < instanceDomainStateCount
}

func (f InstanceDomainState) Exists() bool {
	return f == InstanceDomainStateActive
}

func NewGeneratedInstanceDomain(instanceName, iamDomain string) (string, error) {
	randomString, err := crypto.GenerateRandomString(6, domainRunes)
	if err != nil {
		return "", err
	}
	instanceName = strings.TrimSpace(instanceName)
	return strings.ToLower(strings.ReplaceAll(instanceName, " ", "-") + "-" + randomString + "." + iamDomain), nil
}
