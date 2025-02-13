package repository

import "context"

type InstanceRepository interface {
	SetUp(ctx context.Context, instance *Instance) error
	ByID(ctx context.Context, id string) (*Instance, error)
	ByDomain(ctx context.Context, domain string) (*Instance, error)
}

type Instance struct {
	ID   string
	Name string
}
