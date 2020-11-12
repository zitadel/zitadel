package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

type MailTemplate struct {
	handler
}

const (
	mailTemplateTable = "adminapi.mail_templates"
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
		AggregateTypeFilter(model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *MailTemplate) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = m.processMailTemplate(event)
	}
	return err
}

func (m *MailTemplate) processMailTemplate(event *models.Event) (err error) {
	template := new(iam_model.MailTemplateView)
	switch event.Type {
	case model.MailTemplateAdded:
		err = template.AppendEvent(event)
	case model.MailTemplateChanged:
		template, err = m.view.MailTemplateByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = template.AppendEvent(event)
	default:
		return m.view.ProcessedMailTemplateSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutMailTemplate(template, template.Sequence)
}

func (m *MailTemplate) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Wj8sf", "id", event.AggregateID).WithError(err).Warn("something went wrong in label template handler")
	return spooler.HandleError(event, err, m.view.GetLatestMailTemplateFailedEvent, m.view.ProcessedMailTemplateFailedEvent, m.view.ProcessedMailTemplateSequence, m.errorCountUntilSkip)
}
