package handler

import (
	"github.com/caos/logging"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type MailTemplate struct {
	handler
}

const (
	mailTemplateTable = "management.mail_templates"
)

func (m *MailTemplate) ViewModel() string {
	return mailTemplateTable
}

func (m *MailTemplate) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestMailTemplateSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate, iam_es_model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *MailTemplate) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processMailTemplate(event)
	}
	return err
}

func (m *MailTemplate) processMailTemplate(event *models.Event) (err error) {
	template := new(iam_model.MailTemplateView)
	switch event.Type {
	case iam_es_model.MailTemplateAdded, model.MailTemplateAdded:
		err = template.AppendEvent(event)
	case iam_es_model.MailTemplateChanged, model.MailTemplateChanged:
		template, err = m.view.MailTemplateByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = template.AppendEvent(event)
	case model.MailTemplateRemoved:
		return m.view.DeleteMailTemplate(event.AggregateID, event.Sequence)
	default:
		return m.view.ProcessedMailTemplateSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutMailTemplate(template, template.Sequence)
}

func (m *MailTemplate) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Djo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in label template handler")
	return spooler.HandleError(event, err, m.view.GetLatestMailTemplateFailedEvent, m.view.ProcessedMailTemplateFailedEvent, m.view.ProcessedMailTemplateSequence, m.errorCountUntilSkip)
}
