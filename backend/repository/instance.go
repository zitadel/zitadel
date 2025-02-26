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

type InstanceSetuper interface {
	SetUp(ctx context.Context, instance *Instance) error
}

type instanceByIDQuerier interface {
	ByID(ctx context.Context, id string) (*Instance, error)
}

type instanceByDomainQuerier interface {
	ByDomain(ctx context.Context, domain string) (*Instance, error)
}

type InstanceLister interface {
	List(ctx context.Context) ([]*Instance, error)
}
