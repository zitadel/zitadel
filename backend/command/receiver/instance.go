package receiver

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/receiver/cache"
)

type InstanceState uint8

const (
	InstanceStateActive InstanceState = iota
	InstanceStateDeleted
)

type Instance struct {
	ID      string
	Name    string
	State   InstanceState
	Domains []*Domain
}

type InstanceIndex uint8

var InstanceIndices = []InstanceIndex{
	InstanceByID,
	InstanceByDomain,
}

const (
	InstanceByID InstanceIndex = iota
	InstanceByDomain
)

var _ cache.Entry[InstanceIndex, string] = (*Instance)(nil)

// Keys implements [cache.Entry].
func (i *Instance) Keys(index InstanceIndex) (key []string) {
	switch index {
	case InstanceByID:
		return []string{i.ID}
	case InstanceByDomain:
		return []string{i.Name}
	}
	return nil
}

type InstanceManipulator interface {
	Create(ctx context.Context, instance *Instance) error
	Delete(ctx context.Context, instance *Instance) error
	AddDomain(ctx context.Context, instance *Instance, domain *Domain) error
	SetPrimaryDomain(ctx context.Context, instance *Instance, domain *Domain) error
}

type InstanceReader interface {
	ByID(ctx context.Context, id string) (*Instance, error)
}
