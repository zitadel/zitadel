package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/static"
)

var (
	includePaths = []string{
		"/resources/themes/scss/main.scss",
		"/resources/themes/scss/bundle.scss",
		"/resources/themes/scss/styles/a/a.scss",
		"/resources/themes/scss/styles/a/a_theme.scss",
		"/resources/themes/scss/styles/account_selection/account_selection.scss",
		"/resources/themes/scss/styles/account_selection/account_selection_theme.scss",
		"/resources/themes/scss/styles/avatar/avatar.scss",
		"/resources/themes/scss/styles/avatar/avatar_theme.scss",
		"/resources/themes/scss/styles/button/button.scss",
		"/resources/themes/scss/styles/button/button_base.scss",
		"/resources/themes/scss/styles/button/button_theme.scss",
		"/resources/themes/scss/styles/checkbox/checkbox.scss",
		"/resources/themes/scss/styles/checkbox/checkbox_base.scss",
		"/resources/themes/scss/styles/checkbox/checkbox_theme.scss",
		"/resources/themes/scss/styles/color/all_color.scss",
		"/resources/themes/scss/styles/container/container.scss",
		"/resources/themes/scss/styles/container/container_theme.scss",
		"/resources/themes/scss/styles/core/core.scss",
		"/resources/themes/scss/styles/elevation/elevation.scss",
		"/resources/themes/scss/styles/error/error.scss",
		"/resources/themes/scss/styles/error/error_theme.scss",
		"/resources/themes/scss/styles/footer/footer.scss",
		"/resources/themes/scss/styles/footer/footer_theme.scss",
		"/resources/themes/scss/styles/header/header.scss",
		"/resources/themes/scss/styles/header/header_theme.scss",
		"/resources/themes/scss/styles/identity_provider/identity_provider.scss",
		"/resources/themes/scss/styles/identity_provider/identity_provider_base.scss",
		"/resources/themes/scss/styles/identity_provider/identity_provider_theme.scss",
		"/resources/themes/scss/styles/input/input.scss",
		"/resources/themes/scss/styles/input/input_base.scss",
		"/resources/themes/scss/styles/input/input_theme.scss",
		"/resources/themes/scss/styles/label/label.scss",
		"/resources/themes/scss/styles/label/label_base.scss",
		"/resources/themes/scss/styles/label/label_theme.scss",
		"/resources/themes/scss/styles/list/list.scss",
		"/resources/themes/scss/styles/list/list_base.scss",
		"/resources/themes/scss/styles/list/list_theme.scss",
		"/resources/themes/scss/styles/progress_bar/progress_bar.scss",
		"/resources/themes/scss/styles/progress_bar/progress_bar_base.scss",
		"/resources/themes/scss/styles/progress_bar/progress_bar_theme.scss",
		"/resources/themes/scss/styles/qrcode/qrcode.scss",
		"/resources/themes/scss/styles/qrcode/qrcode_theme.scss",
		"/resources/themes/scss/styles/radio/radio.scss",
		"/resources/themes/scss/styles/radio/radio_base.scss",
		"/resources/themes/scss/styles/radio/radio_theme.scss",
		"/resources/themes/scss/styles/register/register.scss",
		"/resources/themes/scss/styles/select/select.scss",
		"/resources/themes/scss/styles/select/select_base.scss",
		"/resources/themes/scss/styles/select/select_theme.scss",
		"/resources/themes/scss/styles/success_label/success_label.scss",
		"/resources/themes/scss/styles/success_label/success_label_base.scss",
		"/resources/themes/scss/styles/success_label/success_label_theme.scss",
		"/resources/themes/scss/styles/theming/all.scss",
		"/resources/themes/scss/styles/theming/palette.scss",
		"/resources/themes/scss/styles/theming/theming.scss",
	}
)

const (
	stylingTable = "adminapi.styling"
)

type Styling struct {
	handler
	static       static.Storage
	subscription *v1.Subscription
}

func newStyling(handler handler, static static.Storage) *Styling {
	h := &Styling{
		handler: handler,
		static:  static,
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
	cssContent := ""
	cssContent += fmt.Sprint(":root {")
	if policy.PrimaryColor != "" {
		cssContent += fmt.Sprintf("--primary-color: %s;", policy.PrimaryColor)
	}
	if policy.SecondaryColor != "" {
		cssContent += fmt.Sprintf("--secondary-color: %s;", policy.SecondaryColor)
	}
	if policy.PrimaryColor != "" {
		cssContent += fmt.Sprintf("--warn-color: %s;", policy.SecondaryColor)
	}
	if policy.SecondaryColorDark != "" {
		cssContent += fmt.Sprintf("--primary-color-dark: %s;", policy.PrimaryColorDark)
	}
	if policy.SecondaryColorDark != "" {
		cssContent += fmt.Sprintf("--secondary-color-dark: %s;", policy.SecondaryColorDark)
	}
	if policy.WarnColorDark != "" {
		cssContent += fmt.Sprintf("--warn-color-dark: %s;", policy.WarnColorDark)
	}
	if policy.LogoURL != "" {
		cssContent += fmt.Sprintf("--logo-url: %s;", policy.LogoURL)
	}
	if policy.LogoDarkURL != "" {
		cssContent += fmt.Sprintf("--logo-url-dark: %s;", policy.LogoDarkURL)
	}
	if policy.IconURL != "" {
		cssContent += fmt.Sprintf("--icon-url: %s;", policy.IconURL)
	}
	if policy.IconDarkURL != "" {
		cssContent += fmt.Sprintf("--icon-url-dark: %s;", policy.IconDarkURL)
	}
	if policy.FontURL != "" {
		cssContent += fmt.Sprintf("--font-url: %s;", policy.FontURL)
	}
	cssContent += fmt.Sprint("}")

	data := []byte(cssContent)
	buffer := bytes.NewBuffer(data)
	return buffer, int64(buffer.Len()), nil
}

func (m *Styling) uploadFilesToBucket(aggregateID, contentType string, reader io.Reader, size int64) error {
	fileName := domain.OrgCssPath + "/" + domain.CssVariablesFileName
	if aggregateID == domain.IAMID {
		fileName = domain.IAMCssPath + "/" + domain.CssVariablesFileName
	}
	_, err := m.static.PutObject(context.Background(), aggregateID, fileName, contentType, reader, size, true)
	return err
}
