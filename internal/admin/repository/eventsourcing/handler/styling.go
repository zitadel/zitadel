package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"
	"github.com/wellington/go-libsass"

	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/static"
)

const (
	stylingTable = "adminapi.styling"
)

type Styling struct {
	handler
	static       static.Storage
	subscription *v1.Subscription
	dir          http.FileSystem
}

func newStyling(handler handler, static static.Storage) *Styling {
	statikFS, err := fs.NewWithNamespace("login")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start listener")

	h := &Styling{
		handler: handler,
		static:  static,
		dir:     statikFS,
	}

	h.subscribe()

	return h
}

func (m *Styling) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (m *Styling) ViewModel() string {
	return stylingTable
}

func (_ *Styling) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (m *Styling) CurrentSequence() (uint64, error) {
	sequence, err := m.view.GetLatestStylingSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *Styling) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestStylingSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *Styling) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processLabelPolicy(event)
	}
	return err
}

func (m *Styling) processLabelPolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.LabelPolicyView)
	switch event.Type {
	case iam_es_model.LabelPolicyAdded, model.LabelPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.LabelPolicyChanged, model.LabelPolicyChanged,
		iam_es_model.LabelPolicyLogoAdded, model.LabelPolicyLogoAdded,
		iam_es_model.LabelPolicyLogoRemoved, model.LabelPolicyLogoRemoved,
		iam_es_model.LabelPolicyIconAdded, model.LabelPolicyIconAdded,
		iam_es_model.LabelPolicyIconRemoved, model.LabelPolicyIconRemoved,
		iam_es_model.LabelPolicyLogoDarkAdded, model.LabelPolicyLogoDarkAdded,
		iam_es_model.LabelPolicyLogoDarkRemoved, model.LabelPolicyLogoDarkRemoved,
		iam_es_model.LabelPolicyIconDarkAdded, model.LabelPolicyIconDarkAdded,
		iam_es_model.LabelPolicyIconDarkRemoved, model.LabelPolicyIconDarkRemoved:
		policy, err = m.view.StylingByAggregateIDAndState(event.AggregateID, int32(domain.LabelPolicyStatePreview))
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)

	case iam_es_model.LabelPolicyActivated, model.LabelPolicyActivated:
		policy, err = m.view.StylingByAggregateIDAndState(event.AggregateID, int32(domain.LabelPolicyStatePreview))
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
		if err != nil {
			return err
		}
		err = m.generateStylingFile(policy)
	default:
		return m.view.ProcessedStylingSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutStyling(policy, event)
}

func (m *Styling) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-2m9fs", "id", event.AggregateID).WithError(err).Warn("something went wrong in label policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestLabelPolicyFailedEvent, m.view.ProcessedLabelPolicyFailedEvent, m.view.ProcessedLabelPolicySequence, m.errorCountUntilSkip)
}

func (m *Styling) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateLabelPolicySpoolerRunTimestamp)
}

func (m *Styling) generateStylingFile(policy *iam_model.LabelPolicyView) error {
	reader, size, err := m.writeFile(policy)
	if err != nil {
		return err
	}
	return m.uploadFilesToBucket(policy.AggregateID, "text/css", reader, size)
}

func (m *Styling) writeFile(policy *iam_model.LabelPolicyView) (io.Reader, int64, error) {
	r, err := m.dir.Open("/resources/themes/scss/zitadel-alternative.scss")
	if err != nil {
		logging.LogWithFields("SPOOL-2n8fs", "policy", policy.AggregateID).WithError(err).Warn("something went wrong create file")
		return nil, 0, caos_errors.ThrowInternal(err, "SPOOL-f83nf", "error create file")
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, 0, caos_errors.ThrowInternal(err, "SPOOL-2n8fs", "error create file")
	}

	fmt.Println(string(contents))

	output := bytes.NewBuffer(nil)
	comp, err := libsass.New(output, r)
	if err != nil {
		logging.LogWithFields("SPOOL-2n8fs", "policy", policy.AggregateID).WithError(err).Warn("something went wrong create file")
		return nil, 0, caos_errors.ThrowInternal(err, "SPOOL-f83nf", "error create file")
	}

	if err := comp.Run(); err != nil {
		logging.LogWithFields("SPOOL-2n8fs", "policy", policy.AggregateID).WithError(err).Warn("something went wrong create file")
		return nil, 0, caos_errors.ThrowInternal(err, "SPOOL-f83nf", "error create file")
	}
	return output, int64(output.Len()), nil
}

func (m *Styling) uploadFilesToBucket(aggregateID, contentType string, reader io.Reader, size int64) error {
	fileName := domain.OrgCssPath + "/" + "main"
	if aggregateID == domain.IAMID {
		fileName = domain.IAMCssPath + "/" + "main"
	}
	_, err := m.static.PutObject(context.Background(), aggregateID, fileName, contentType, reader, size, true)
	return err
}
