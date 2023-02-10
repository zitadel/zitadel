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
	DefaultOrganisationID() string
	SecurityPolicyAllowedOrigins() []string
}

type InstanceVerifier interface {
	InstanceByHost(ctx context.Context, host string) (Instance, error)
	InstanceByID(ctx context.Context) (Instance, error)
}

type instance struct {
	id        string
	domain    string
	projectID string
	appID     string
	clientID  string
	orgID     string
}

func (i *instance) InstanceID() string {
	return i.id
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
	return i.domain
}

func (i *instance) RequestedHost() string {
	return i.domain
}

func (i *instance) DefaultLanguage() language.Tag {
	return language.Und
}

func (i *instance) DefaultOrganisationID() string {
	return i.orgID
}

func (i *instance) SecurityPolicyAllowedOrigins() []string {
	return nil
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
	return context.WithValue(ctx, instanceKey, &instance{id: id})
}

func WithRequestedDomain(ctx context.Context, domain string) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}

	i.domain = domain
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
