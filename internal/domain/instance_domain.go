package domain

import (
	"strings"
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

func NewGeneratedInstanceDomain(instanceName, iamDomain string) string {
	//TODO: Add random number/string to be unique
	instanceName = strings.TrimSpace(instanceName)
	return strings.ToLower(strings.ReplaceAll(instanceName, " ", "-") + "." + iamDomain)
}
