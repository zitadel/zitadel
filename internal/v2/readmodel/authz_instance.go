package readmodel

// import (
// 	"context"
// 	"strings"
// 	"time"

// 	"github.com/zitadel/logging"
// 	"golang.org/x/text/language"

// 	"github.com/zitadel/zitadel/internal/api/authz"
// 	"github.com/zitadel/zitadel/internal/eventstore"
// 	"github.com/zitadel/zitadel/internal/feature"
// 	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
// 	"github.com/zitadel/zitadel/internal/v2/instance"
// 	"github.com/zitadel/zitadel/internal/v2/projection"
// 	"github.com/zitadel/zitadel/internal/zerrors"
// )

// type AuthZInstances struct {
// 	cache   *CachedReadModel[*AuthZInstances, *AuthZInstance]
// 	queries authzInstanceQueries
// }

// type authzInstanceQueries interface {
// 	InstanceByHost(ctx context.Context, instanceHost, publicHost string) (authz.Instance, error)
// }

// func NewAuthZInstances(ctx context.Context, eventStore *eventstore.Eventstore, query authzInstanceQueries) *AuthZInstances {
// 	instances := &AuthZInstances{
// 		queries: query,
// 	}
// 	instances.cache = NewCachedReadModel[*AuthZInstances, *AuthZInstance](ctx, instances, eventStore)
// 	return instances
// }

// func (instances *AuthZInstances) InstanceByHost(ctx context.Context, instanceHost, publicHost string) (authz.Instance, error) {
// 	return instances.ByHost(ctx, instanceHost, publicHost)
// }

// func (instances *AuthZInstances) InstanceByID(ctx context.Context, id string) (authz.Instance, error) {
// 	return nil, zerrors.ThrowUnimplemented(nil, "READM-m9NSw", "method not implemented")
// }

// func (instances *AuthZInstances) ByHost(ctx context.Context, instanceHost, publicHost string) (*AuthZInstance, error) {
// 	instanceDomain := strings.Split(instanceHost, ":")[0] // remove possible port
// 	publicDomain := strings.Split(publicHost, ":")[0]     // remove possible port

// 	if instance, ok := instances.cache.get(instanceDomain); ok {
// 		return instance, nil
// 	}
// 	if instance, ok := instances.cache.get(publicDomain); ok {
// 		return instance, nil
// 	}
// 	queriedInstance, err := instances.queries.InstanceByHost(ctx, instanceHost, publicHost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	instance := AuthZInstanceFromAuthZ(queriedInstance)
// 	if instanceDomain != "" {
// 		err = instances.cache.set(instanceDomain, instance)
// 		logging.OnError(err).Debug("failed to cache instance")
// 	}
// 	if publicDomain != "" {
// 		err = instances.cache.set(publicDomain, instance)
// 		logging.OnError(err).Debug("failed to cache instance")
// 	}
// 	return instance, nil
// }

// func (instances *AuthZInstances) Reducers() map[string]map[string]v2_es.ReduceEvent {
// 	return map[string]map[string]v2_es.ReduceEvent{
// 		instance.AggregateType: {
// 			instance.AddedType:              instances.reduceEvent,
// 			instance.ChangedType:            instances.reduceEvent,
// 			instance.DefaultOrgSetType:      instances.reduceEvent,
// 			instance.ProjectSetType:         instances.reduceEvent,
// 			instance.ConsoleSetType:         instances.reduceEvent,
// 			instance.DefaultLanguageSetType: instances.reduceEvent,
// 			instance.RemovedType:            instances.reduceEvent,
// 		},
// 	}
// }

// // Reduce implements [v2_es.Reducer].
// func (i *AuthZInstances) Reduce(events ...*v2_es.StorageEvent) error {
// 	for _, event := range events {
// 		if err := i.reduceEvent(event); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (i *AuthZInstances) reduceEvent(event *v2_es.StorageEvent) error {
// 	if event.Type == instance.AddedType {
// 		instance := NewAuthZInstanceFromEvent(event)
// 		if err := i.cache.set(instance.ID, instance); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	instance, ok := i.cache.get(event.Aggregate.ID)
// 	if !ok {
// 		return nil
// 	}
// 	if err := instance.Reduce(event); err != nil {
// 		return err
// 	}
// 	if err := i.cache.set(instance.ID, instance); err != nil {
// 		return err
// 	}
// 	return nil
// }

// var (
// 	_ authz.Instance = (*AuthZInstance)(nil)
// 	// _ cacheModel     = (*AuthZInstance)(nil)
// )

// type AuthZInstance struct {
// 	projection.AuthZInstance
// 	InstanceFeatures
// 	SystemFeatures
// }

// func NewAuthZInstanceFromEvent(event *v2_es.StorageEvent) *AuthZInstance {
// 	return &AuthZInstance{
// 		AuthZInstance:  *projection.NewAuthZInstanceFromEvent(event),
// 		SystemFeatures: *NewSystemFeatures(),
// 	}
// }

// func AuthZInstanceFromAuthZ(authZInstance authz.Instance) *AuthZInstance {
// 	state := projection.NewInstanceStateProjection(authZInstance.InstanceID())
// 	state.State = instance.ActiveState

// 	return &AuthZInstance{
// 		AuthZInstance: projection.AuthZInstance{
// 			ID:              authZInstance.InstanceID(),
// 			ProjectID:       authZInstance.ProjectID(),
// 			ConsoleAppID:    authZInstance.ConsoleApplicationID(),
// 			ConsoleClientID: authZInstance.ConsoleClientID(),
// 			DefaultOrgID:    authZInstance.DefaultOrganisationID(),
// 			DefaultLanguage: authZInstance.DefaultLanguage(),
// 			State:           state,
// 		},
// 	}
// }

// // Reduce implements model.
// // Subtle: this method shadows the method (Instance).Reduce of AuthZInstance.Instance.
// func (i *AuthZInstance) Reduce(events ...*v2_es.StorageEvent) error {
// 	return i.AuthZInstance.Reduce(events...)
// }

// // AuditLogRetention implements [authz.Instance].
// func (i *AuthZInstance) AuditLogRetention() *time.Duration {
// 	// TODO: implement
// 	return nil
// }

// // Block implements [authz.Instance].
// func (i *AuthZInstance) Block() *bool {
// 	// TODO: implement
// 	return nil
// }

// // ConsoleApplicationID implements [authz.Instance].
// func (i *AuthZInstance) ConsoleApplicationID() string {
// 	return i.ConsoleAppID
// }

// // ConsoleClientID implements [authz.Instance].
// func (i *AuthZInstance) ConsoleClientID() string {
// 	return i.AuthZInstance.ConsoleClientID
// }

// // DefaultLanguage implements [authz.Instance].
// func (i *AuthZInstance) DefaultLanguage() language.Tag {
// 	return i.AuthZInstance.DefaultLanguage
// }

// // DefaultOrganisationID implements [authz.Instance].
// func (i *AuthZInstance) DefaultOrganisationID() string {
// 	return i.DefaultOrgID
// }

// // EnableImpersonation implements [authz.Instance].
// func (i *AuthZInstance) EnableImpersonation() bool {
// 	// TODO: implement
// 	return true
// }

// // Features implements [authz.Instance].
// func (i *AuthZInstance) Features() feature.Features {
// 	// TODO: implement
// 	return feature.Features{
// 		LoginDefaultOrg:                 true,
// 		TriggerIntrospectionProjections: false,
// 		LegacyIntrospection:             false,
// 		UserSchema:                      true,
// 		TokenExchange:                   false,
// 		Actions:                         true,
// 		ImprovedPerformance: []feature.ImprovedPerformanceType{
// 			feature.ImprovedPerformanceTypeOrgByID,
// 			feature.ImprovedPerformanceTypeProjectGrant,
// 			feature.ImprovedPerformanceTypeProject,
// 			feature.ImprovedPerformanceTypeUserGrant,
// 			feature.ImprovedPerformanceTypeOrgDomainVerified,
// 		},
// 		WebKey:                         false,
// 		DebugOIDCParentError:           true,
// 		OIDCSingleV1SessionTermination: false,
// 		InMemoryProjections:            true,
// 	}
// }

// // InstanceID implements [authz.Instance].
// func (i *AuthZInstance) InstanceID() string {
// 	return i.ID
// }

// // ProjectID implements [authz.Instance].
// func (i *AuthZInstance) ProjectID() string {
// 	return i.AuthZInstance.ProjectID
// }

// // SecurityPolicyAllowedOrigins implements [authz.Instance].
// func (i *AuthZInstance) SecurityPolicyAllowedOrigins() []string {
// 	// TODO: implement
// 	return nil
// }
