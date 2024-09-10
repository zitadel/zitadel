package readmodel

import (
	"context"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
	"github.com/zitadel/zitadel/internal/v2/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	_ object         = (*AuthZInstance)(nil)
	_ authz.Instance = (*AuthZInstance)(nil)
)

type AuthZInstance struct {
	readModel

	*projection.AuthZInstance
	*InstanceFeatures
}

func (i *AuthZInstance) Reducers() map[string]map[string]v2_es.ReduceEvent {
	if i.reducers != nil {
		return i.reducers
	}

	i.reducers = mergeReducers(i.AuthZInstance.Reducers(), i.InstanceFeatures.Reducers())

	return i.reducers
}

func (i *AuthZInstance) InstanceID() string {
	return i.ID
}

func (i *AuthZInstance) ProjectID() string {
	return i.AuthZInstance.ProjectID
}

func (i *AuthZInstance) ConsoleClientID() string {
	return i.AuthZInstance.ConsoleClientID
}

func (i *AuthZInstance) ConsoleApplicationID() string {
	return i.AuthZInstance.ConsoleAppID
}

func (i *AuthZInstance) DefaultLanguage() language.Tag {
	return i.AuthZInstance.DefaultLanguage
}

func (i *AuthZInstance) DefaultOrganisationID() string {
	return i.DefaultOrgID
}

func (i *AuthZInstance) SecurityPolicyAllowedOrigins() []string {
	return nil
}

func (i *AuthZInstance) EnableImpersonation() bool {
	return false
}

func (i *AuthZInstance) Block() *bool {
	return nil
}

func (i *AuthZInstance) AuditLogRetention() *time.Duration {
	return nil
}

func (i *AuthZInstance) Features() feature.Features {
	return feature.Features{
		LoginDefaultOrg:                 i.LoginDefaultOrg.Value(),
		TriggerIntrospectionProjections: i.TriggerIntrospectionProjections.Value(),
		LegacyIntrospection:             i.LegacyIntrospection.Value(),
		UserSchema:                      i.UserSchema.Value(),
		TokenExchange:                   i.TokenExchange.Value(),
		Actions:                         i.Actions.Value(),
		ImprovedPerformance:             i.ImprovedPerformance.Value(),
		WebKey:                          i.WebKey.Value(),
		DebugOIDCParentError:            i.DebugOIDCParentError.Value(),
		OIDCSingleV1SessionTermination:  i.OIDCSingleV1SessionTermination.Value(),
		InMemoryProjections:             i.InMemoryProjections.Value(),
	}
}

var _ listManager = (*AuthZInstances)(nil)

// AuthZInstances is the manager for the instance list read model.
type AuthZInstances struct {
	readModel

	idCache     Cache[string, *AuthZInstance]
	domainCache Cache[string, string]
	object      *objectReadModel

	systemFeatures *SystemFeatures
	query          authzInstanceQueries
}

type authzInstanceQueries interface {
	InstanceByHost(ctx context.Context, instanceHost, publicHost string) (authz.Instance, error)
}

func NewAuthZInstances(ctx context.Context, es *eventstore.Eventstore, systemFeatures *SystemFeatures, query authzInstanceQueries) *AuthZInstances {
	instances := &AuthZInstances{
		idCache:        NewMapCache[string, *AuthZInstance](),
		domainCache:    NewMapCache[string, string](),
		systemFeatures: systemFeatures,
		query:          query,
	}
	instances.object = newObjectReadModel(ctx, instances, es)
	return instances
}

func (instances *AuthZInstances) InstanceByHost(ctx context.Context, instanceHost, publicHost string) (authz.Instance, error) {
	instanceDomain := strings.Split(instanceHost, ":")[0] // remove possible port

	instanceID, ok := instances.domainCache.Get(instanceDomain)
	if !ok {
		publicDomain := strings.Split(publicHost, ":")[0] // remove possible port
		instanceID, ok = instances.domainCache.Get(publicDomain)
		if !ok {
			queriedInstance, err := instances.query.InstanceByHost(ctx, instanceHost, publicHost)
			if err != nil {
				return nil, err
			}
			if err = instances.loadInstance(ctx, queriedInstance.InstanceID()); err != nil {
				return nil, err
			}
			instanceID = queriedInstance.InstanceID()
		}
	}

	return instances.InstanceByID(ctx, instanceID)
}

func (instances *AuthZInstances) InstanceByID(ctx context.Context, id string) (authz.Instance, error) {
	instance, ok := instances.idCache.Get(id)
	if !ok {
		// TODO: query eventstore?
		return nil, zerrors.ThrowNotFound(nil, "READM-m9NSw", "Errors.Instance.NotFound")
	}

	return instance, nil
}

// Name implements manager.
func (i *AuthZInstances) Name() string {
	return "authz_instances"
}

// EventstoreV3Query implements manager.
func (i *AuthZInstances) EventstoreV3Query(position decimal.Decimal) *eventstore.SearchQueryBuilder {
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderAsc()

	for aggregateType, eventReduces := range i.Reducers() {
		eventTypes := make([]eventstore.EventType, 0, len(eventReduces))
		for eventType := range eventReduces {
			eventTypes = append(eventTypes, eventstore.EventType(eventType))
		}
		builder = builder.AddQuery().AggregateTypes(eventstore.AggregateType(aggregateType)).EventTypes(eventTypes...).Builder()
	}

	return builder
}

// Reducers implements manager.
func (i *AuthZInstances) Reducers() map[string]map[string]v2_es.ReduceEvent {
	if i.reducers != nil {
		return i.reducers
	}

	i.reducers = map[string]map[string]v2_es.ReduceEvent{
		instance.AggregateType: {
			instance.AddedType:              i.reduceAdded,
			instance.ChangedType:            i.reduce,
			instance.DefaultOrgSetType:      i.reduce,
			instance.ProjectSetType:         i.reduce,
			instance.ConsoleSetType:         i.reduce,
			instance.DefaultLanguageSetType: i.reduce,
			instance.RemovedType:            i.reduceRemoved,

			instance.DomainAddedType:      i.reduceDomainAdded,
			instance.DomainVerifiedType:   i.reduce,
			instance.DomainPrimarySetType: i.reduce,
			instance.DomainRemovedType:    i.reduceDomainRemoved,
		},
	}

	return i.reducers
}

func (i *AuthZInstances) reduceAdded(event *v2_es.StorageEvent) error {
	instance, ok := i.idCache.Get(event.Aggregate.ID)
	if !ok {
		instance = &AuthZInstance{
			InstanceFeatures: NewInstanceFeatures(context.TODO(), i.object.es, i.systemFeatures, event.Aggregate.ID),
			AuthZInstance:    projection.NewAuthZInstanceFromEvent(event),
		}
	}
	err := instance.Reducers()[event.Aggregate.Type][event.Type](event)
	if err != nil {
		return err
	}
	return i.idCache.Set(instance.ID, instance)
}

func (i *AuthZInstances) reduceRemoved(event *v2_es.StorageEvent) error {
	instance, ok := i.idCache.Get(event.Aggregate.ID)
	if !ok {
		return nil
	}
	err := i.idCache.Remove(instance.ID)
	if err != nil {
		return err
	}
	for _, domain := range instance.Domains {
		err = i.domainCache.Remove(domain.Name)
		i.object.logEvent(event).OnError(err).WithField("domain", domain.Domain).Warn("could not remove domain from cache")
	}
	return nil
}

func (i *AuthZInstances) reduce(event *v2_es.StorageEvent) error {
	instance, ok := i.idCache.Get(event.Aggregate.ID)
	if !ok {
		return nil
	}
	err := instance.Reducers()[event.Aggregate.Type][event.Type](event)
	if err != nil {
		return err
	}
	return i.idCache.Set(instance.ID, instance)
}

func (i *AuthZInstances) reduceDomainAdded(event *v2_es.StorageEvent) error {
	if err := i.reduce(event); err != nil {
		return err
	}
	e, err := instance.DomainAddedEventFromStorage(event)
	if err != nil {
		return err
	}
	return i.domainCache.Set(e.Payload.Name, e.Aggregate.ID)
}

func (i *AuthZInstances) reduceDomainRemoved(event *v2_es.StorageEvent) error {
	if err := i.reduce(event); err != nil {
		return err
	}
	e, err := instance.DomainRemovedEventFromStorage(event)
	if err != nil {
		return err
	}
	return i.domainCache.Remove(e.Payload.Name)
}

func (i *AuthZInstances) loadInstance(ctx context.Context, id string) error {
	return i.object.es.FilterToReducer(ctx, i.EventstoreV3Query(decimal.Zero).InstanceID(id), i.object)
}
