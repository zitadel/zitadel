package domain

import (
	"strings"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type InstanceDomain struct {
	models.ObjectRoot

	Domain    string
	Generated bool
}

func (i *InstanceDomain) IsValid() bool {
	return i.AggregateID != "" && i.Domain != ""
}

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
