package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/repository/view/model"
)

type Org struct {
	handler
}

const (
	orgTable = "management.orgs"
)

func (o *Org) ViewModel() string {
	return orgTable
}

func (_ *Org) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate}
}

func (o *Org) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := o.view.GetLatestOrgSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (o *Org) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := o.view.GetLatestOrgSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.OrgQuery(sequence.CurrentSequence), nil
}

func (o *Org) Reduce(event *es_models.Event) (err error) {
	org := new(org_model.OrgView)

	switch event.Type {
	case model.OrgAdded:
		err = org.AppendEvent(event)
	case model.OrgChanged:
		err = org.SetData(event)
		if err != nil {
			return err
		}
		org, err = o.view.OrgByID(org.ID)
		if err != nil {
			return err
		}
		err = org.AppendEvent(event)
	default:
		return o.view.ProcessedOrgSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return o.view.PutOrg(org, event.CreationDate)
}

func (o *Org) OnError(event *es_models.Event, spoolerErr error) error {
	logging.LogWithFields("SPOOL-ls9ew", "id", event.AggregateID).WithError(spoolerErr).Warn("something went wrong in project app handler")
	return spooler.HandleError(event, spoolerErr, o.view.GetLatestOrgFailedEvent, o.view.ProcessedOrgFailedEvent, o.view.ProcessedOrgSequence, o.errorCountUntilSkip)
}

func (o *Org) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateOrgSpoolerRunTimestamp)
}
