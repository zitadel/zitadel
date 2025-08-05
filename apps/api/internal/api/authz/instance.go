package authz

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/feature"
)

var emptyInstance = &instance{}

type Instance interface {
	InstanceID() string
	ProjectID() string
	ConsoleClientID() string
	ConsoleApplicationID() string
	DefaultLanguage() language.Tag
	DefaultOrganisationID() string
	SecurityPolicyAllowedOrigins() []string
	EnableImpersonation() bool
	Block() *bool
	AuditLogRetention() *time.Duration
	Features() feature.Features
}

type InstanceVerifier interface {
	InstanceByHost(ctx context.Context, host, publicDomain string) (Instance, error)
	InstanceByID(ctx context.Context, id string) (Instance, error)
}

type instance struct {
	id              string
	projectID       string
	appID           string
	clientID        string
	orgID           string
	defaultLanguage language.Tag
	features        feature.Features
}

func (i *instance) Block() *bool {
	return nil
}

func (i *instance) AuditLogRetention() *time.Duration {
	return nil
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

func (i *instance) DefaultLanguage() language.Tag {
	return i.defaultLanguage
}

func (i *instance) DefaultOrganisationID() string {
	return i.orgID
}

func (i *instance) SecurityPolicyAllowedOrigins() []string {
	return nil
}

func (i *instance) EnableImpersonation() bool {
	return false
}

func (i *instance) Features() feature.Features {
	return i.features
}

func GetInstance(ctx context.Context) Instance {
	instance, ok := ctx.Value(instanceKey).(Instance)
	if !ok {
		return emptyInstance
	}
	return instance
}

func GetFeatures(ctx context.Context) feature.Features {
	return GetInstance(ctx).Features()
}

func WithInstance(ctx context.Context, instance Instance) context.Context {
	return context.WithValue(ctx, instanceKey, instance)
}

func WithInstanceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, instanceKey, &instance{id: id})
}

func WithDefaultLanguage(ctx context.Context, defaultLanguage language.Tag) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}

	i.defaultLanguage = defaultLanguage
	return context.WithValue(ctx, instanceKey, i)
}

func WithConsole(ctx context.Context, projectID, appID string) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}

	i.projectID = projectID
	i.appID = appID
	return context.WithValue(ctx, instanceKey, i)
}

func WithConsoleClientID(ctx context.Context, clientID string) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}
	i.clientID = clientID
	return context.WithValue(ctx, instanceKey, i)
}

func WithFeatures(ctx context.Context, f feature.Features) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}
	i.features = f
	return context.WithValue(ctx, instanceKey, i)
}
