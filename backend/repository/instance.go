package repository

import "context"

type InstanceRepository interface {
	InstanceSetuper
	instanceByIDQuerier
	instanceByDomainQuerier
}

type Instance struct {
	ID   string
	Name string
}

type SetUpInstance func(ctx context.Context, instance *Instance) error

type InstanceSetuper interface {
	SetUp(ctx context.Context, instance *Instance) error
}

type InstanceByID func(ctx context.Context, id string) (*Instance, error)

type instanceByIDQuerier interface {
	ByID(ctx context.Context, id string) (*Instance, error)
}

type InstanceByDomain func(ctx context.Context, domain string) (*Instance, error)

type instanceByDomainQuerier interface {
	ByDomain(ctx context.Context, domain string) (*Instance, error)
}

type ListInstances func(ctx context.Context) ([]*Instance, error)

type InstanceLister interface {
	List(ctx context.Context) ([]*Instance, error)
}
