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
}

type InstanceVerifier interface {
	InstanceByHost(context.Context, string) (Instance, error)
	//CheckInstance(ctx context.Context, host string) (*Instance, error)
	//instanceRepo
}

type instance struct {
	ID string
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

type instanceRepo interface {
}

//
//func (v InstanceVerifier) CheckInstance(ctx context.Context, host string) (Instance, error) {
//	instance, err := v.instanceRepo.InstanceByHost(ctx, host)
//	if err != nil {
//
//	}
//	return instance, nil
//}
