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

func NewGeneratedInstanceDomain(instanceName, iamDomain string) string {
	return strings.ToLower(strings.ReplaceAll(instanceName, " ", "-") + "." + iamDomain)
}
