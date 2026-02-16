package authz

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/feature"
)

var (
	emptyInstance          = &instance{}
	_             Instance = (*instance)(nil)
)

type Instance interface {
	InstanceID() string
	ProjectID() string
	ManagementConsoleClientID() string
	ManagementConsoleApplicationID() string
	DefaultLanguage() language.Tag
	AllowedLanguages() []language.Tag
	DefaultOrganisationID() string
	SecurityPolicyAllowedOrigins() []string
	EnableImpersonation() bool
	Block() *bool
	AuditLogRetention() *time.Duration
	Features() feature.Features
	ExecutionRouter() target.Router
}

type InstanceVerifier interface {
	// InstanceByHost returns the instance for the given instanceDomain or publicDomain.
	// Previously it used the host (hostname[:port]) to find the instance, but is now using the domain (hostname) only.
	// For preventing issues in backports, the name of the method is not changed.
	InstanceByHost(ctx context.Context, instanceDomain, publicDomain string) (Instance, error)
	InstanceByID(ctx context.Context, id string) (Instance, error)
}

type instance struct {
	id               string
	projectID        string
	appID            string
	clientID         string
	orgID            string
	defaultLanguage  language.Tag
	allowedLanguages []language.Tag
	features         feature.Features
	executionTargets target.Router
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

func (i *instance) ManagementConsoleClientID() string {
	return i.clientID
}

func (i *instance) ManagementConsoleApplicationID() string {
	return i.appID
}

func (i *instance) DefaultLanguage() language.Tag {
	return i.defaultLanguage
}

func (i *instance) AllowedLanguages() []language.Tag {
	return i.allowedLanguages
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

func (i *instance) ExecutionRouter() target.Router {
	return i.executionTargets
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
	ctx = instrumentation.SetInstance(ctx, instance)
	return context.WithValue(ctx, instanceKey, instance)
}

func WithInstanceID(ctx context.Context, id string) context.Context {
	return WithInstance(ctx, &instance{id: id})
}

func WithDefaultLanguage(ctx context.Context, defaultLanguage language.Tag) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}

	i.defaultLanguage = defaultLanguage
	return context.WithValue(ctx, instanceKey, i)
}

func WithManagementConsole(ctx context.Context, projectID, appID string) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}

	i.projectID = projectID
	i.appID = appID
	return context.WithValue(ctx, instanceKey, i)
}

func WithManagementConsoleClientID(ctx context.Context, clientID string) context.Context {
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

func WithExecutionRouter(ctx context.Context, router target.Router) context.Context {
	i, ok := ctx.Value(instanceKey).(*instance)
	if !ok {
		i = new(instance)
	}
	i.executionTargets = router
	return context.WithValue(ctx, instanceKey, i)
}
