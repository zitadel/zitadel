package authz

import (
	"context"
)

var (
	emptyInstance = &instance{}
)

type Instance interface {
	InstanceID() string
	ProjectID() string
	ConsoleClientID() string
	RequestedDomain() string
}

type InstanceVerifier interface {
	InstanceByHost(context.Context, string) (Instance, error)
}

type instance struct {
	ID     string
	Domain string
}

func (i *instance) InstanceID() string {
	return i.ID
}

func (i *instance) ProjectID() string {
	return ""
}

func (i *instance) ConsoleClientID() string {
	return ""
}

func (i *instance) RequestedDomain() string {
	return i.Domain
}

func GetInstance(ctx context.Context) Instance {
	instance, ok := ctx.Value(instanceKey).(Instance)
	if !ok {
		return emptyInstance
	}
	return instance
}

func WithInstance(ctx context.Context, instance Instance) context.Context {
	return context.WithValue(ctx, instanceKey, instance)
}

func WithInstanceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, instanceKey, &instance{ID: id})
}

func WithRequestedDomain(ctx context.Context, domain string) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}

	i.Domain = domain
	return context.WithValue(ctx, instanceKey, i)
}
