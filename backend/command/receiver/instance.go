package receiver

import "context"

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

type InstanceManipulator interface {
	Create(ctx context.Context, instance *Instance) error
	Delete(ctx context.Context, instance *Instance) error
	AddDomain(ctx context.Context, instance *Instance, domain *Domain) error
	SetPrimaryDomain(ctx context.Context, instance *Instance, domain *Domain) error
}

type InstanceReader interface {
	ByID(ctx context.Context, id string) (*Instance, error)
}
