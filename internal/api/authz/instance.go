package authz

import (
	"context"

	"golang.org/x/text/language"
)

var (
	emptyInstance = &instance{}
)

type Instance interface {
	InstanceID() string
	ProjectID() string
	ConsoleClientID() string
	ConsoleApplicationID() string
	RequestedDomain() string
	RequestedHost() string
	DefaultLanguage() language.Tag
}

type InstanceVerifier interface {
	InstanceByHost(context.Context, string) (Instance, error)
}

type instance struct {
	ID        string
	Domain    string
	projectID string
	appID     string
	clientID  string
}

func (i *instance) InstanceID() string {
	return i.ID
}

func (i *instance) ProjectID() string {
	return i.projectID
}

func (i *instance) ConsoleClientID() string {
	return i.clientID
}

func (i *instance) ConsoleApplicationID() string {
	return i.appID
}

func (i *instance) RequestedDomain() string {
	return i.Domain
}

func (i *instance) RequestedHost() string {
	return i.Domain
}

func (i *instance) DefaultLanguage() language.Tag {
	return language.Und
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

func WithConsole(ctx context.Context, projectID, appID string) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}

	i.projectID = projectID
	i.appID = appID
	//i.clientID = clientID
	return context.WithValue(ctx, instanceKey, i)
}
