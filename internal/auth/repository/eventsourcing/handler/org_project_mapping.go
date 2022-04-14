package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/repository/project"
)

const (
	orgProjectMappingTable = "auth.org_project_mapping"
)

type OrgProjectMapping struct {
	handler
	subscription *v1.Subscription
}

func newOrgProjectMapping(
	handler handler,
) *OrgProjectMapping {
	h := &OrgProjectMapping{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (k *OrgProjectMapping) subscribe() {
	k.subscription = k.es.Subscribe(k.AggregateTypes()...)
	go func() {
		for event := range k.subscription.Events {
			query.ReduceEvent(k, event)
		}
	}()
}

func (p *OrgProjectMapping) ViewModel() string {
	return orgProjectMappingTable
}

func (p *OrgProjectMapping) Subscription() *v1.Subscription {
	return p.subscription
}

func (_ *OrgProjectMapping) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{project.AggregateType}
}

func (p *OrgProjectMapping) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestOrgProjectMappingSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *OrgProjectMapping) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestOrgProjectMappingSequence()
	if err != nil {
		return nil, err
	}
	return proj_view.ProjectQuery(sequence.CurrentSequence), nil
}

func (p *OrgProjectMapping) Reduce(event *es_models.Event) (err error) {
	mapping := new(view_model.OrgProjectMapping)
	switch eventstore.EventType(event.Type) {
	case project.ProjectAddedType:
		mapping.OrgID = event.ResourceOwner
		mapping.ProjectID = event.AggregateID
		mapping.InstanceID = event.InstanceID
	case project.ProjectRemovedType:
		err := p.view.DeleteOrgProjectMappingsByProjectID(event.AggregateID)
		if err == nil {
			return p.view.ProcessedOrgProjectMappingSequence(event)
		}
	case project.GrantAddedType:
		projectGrant := new(view_model.ProjectGrant)
		projectGrant.SetData(event)
		mapping.OrgID = projectGrant.GrantedOrgID
		mapping.ProjectID = event.AggregateID
		mapping.ProjectGrantID = projectGrant.GrantID
		mapping.InstanceID = projectGrant.InstanceID
	case project.GrantRemovedType:
		projectGrant := new(view_model.ProjectGrant)
		projectGrant.SetData(event)
		err := p.view.DeleteOrgProjectMappingsByProjectGrantID(event.AggregateID)
		if err == nil {
			return p.view.ProcessedOrgProjectMappingSequence(event)
		}
	default:
		return p.view.ProcessedOrgProjectMappingSequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutOrgProjectMapping(mapping, event)
}

func (p *OrgProjectMapping) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-2k0fS", "id", event.AggregateID).WithError(err).Warn("something went wrong in org project mapping handler")
	return spooler.HandleError(event, err, p.view.GetLatestOrgProjectMappingFailedEvent, p.view.ProcessedOrgProjectMappingFailedEvent, p.view.ProcessedOrgProjectMappingSequence, p.errorCountUntilSkip)
}

func (p *OrgProjectMapping) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateOrgProjectMappingSpoolerRunTimestamp)
}
