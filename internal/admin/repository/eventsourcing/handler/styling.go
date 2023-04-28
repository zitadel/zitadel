package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/zitadel/zitadel/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/static"
)

const (
	stylingTable = "adminapi.styling2"
)

type Styling struct {
	handler
	static       static.Storage
	subscription *v1.Subscription
}

func newStyling(ctx context.Context, handler handler, static static.Storage) *Styling {
	h := &Styling{
		handler: handler,
		static:  static,
	}
	h.subscribe(ctx)

	return h
}

func (m *Styling) subscribe(ctx context.Context) {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(ctx, m, event)
		}
	}()
}

func (m *Styling) ViewModel() string {
	return stylingTable
}

func (m *Styling) Subscription() *v1.Subscription {
	return m.subscription
}

func (_ *Styling) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{org.AggregateType, instance.AggregateType}
}

func (m *Styling) CurrentSequence(ctx context.Context, instanceID string) (uint64, error) {
	sequence, err := m.view.GetLatestStylingSequence(ctx, instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *Styling) EventQuery(ctx context.Context, instanceIDs []string) (*models.SearchQuery, error) {
	sequences, err := m.view.GetLatestStylingSequences(ctx, instanceIDs)
	if err != nil {
		return nil, err
	}
	searchQuery := models.NewSearchQuery()
	for _, instanceID := range instanceIDs {
		var seq uint64
		for _, sequence := range sequences {
			if sequence.InstanceID == instanceID {
				seq = sequence.CurrentSequence
				break
			}
		}
		searchQuery.AddQuery().
			AggregateTypeFilter(m.AggregateTypes()...).
			LatestSequenceFilter(seq).
			InstanceIDFilter(instanceID)
	}
	return searchQuery, nil
}

func (m *Styling) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case org.AggregateType, instance.AggregateType:
		err = m.processLabelPolicy(event)
	}
	return err
}

func (m *Styling) processLabelPolicy(event *models.Event) (err error) {
	policy := new(iam_model.LabelPolicyView)
	switch eventstore.EventType(event.Type) {
	case instance.LabelPolicyAddedEventType,
		org.LabelPolicyAddedEventType:
		err = policy.AppendEvent(event)
	case instance.LabelPolicyChangedEventType,
		org.LabelPolicyChangedEventType,
		instance.LabelPolicyLogoAddedEventType,
		org.LabelPolicyLogoAddedEventType,
		instance.LabelPolicyLogoRemovedEventType,
		org.LabelPolicyLogoRemovedEventType,
		instance.LabelPolicyIconAddedEventType,
		org.LabelPolicyIconAddedEventType,
		instance.LabelPolicyIconRemovedEventType,
		org.LabelPolicyIconRemovedEventType,
		instance.LabelPolicyLogoDarkAddedEventType,
		org.LabelPolicyLogoDarkAddedEventType,
		instance.LabelPolicyLogoDarkRemovedEventType,
		org.LabelPolicyLogoDarkRemovedEventType,
		instance.LabelPolicyIconDarkAddedEventType,
		org.LabelPolicyIconDarkAddedEventType,
		instance.LabelPolicyIconDarkRemovedEventType,
		org.LabelPolicyIconDarkRemovedEventType,
		instance.LabelPolicyFontAddedEventType,
		org.LabelPolicyFontAddedEventType,
		instance.LabelPolicyFontRemovedEventType,
		org.LabelPolicyFontRemovedEventType,
		instance.LabelPolicyAssetsRemovedEventType,
		org.LabelPolicyAssetsRemovedEventType:
		policy, err = m.view.StylingByAggregateIDAndState(event.AggregateID, event.InstanceID, int32(domain.LabelPolicyStatePreview))
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)

	case instance.LabelPolicyActivatedEventType,
		org.LabelPolicyActivatedEventType:
		policy, err = m.view.StylingByAggregateIDAndState(event.AggregateID, event.InstanceID, int32(domain.LabelPolicyStatePreview))
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
		if err != nil {
			return err
		}
		err = m.generateStylingFile(policy)
	case instance.InstanceRemovedEventType:
		err = m.deleteInstanceFilesFromStorage(event.InstanceID)
		if err != nil {
			return err
		}
		return m.view.DeleteInstanceStyling(event)
	case org.OrgRemovedEventType:
		return m.view.UpdateOrgOwnerRemovedStyling(event)
	default:
		return m.view.ProcessedStylingSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutStyling(policy, event)
}

func (m *Styling) OnError(event *models.Event, err error) error {
	logging.WithFields("id", event.AggregateID).WithError(err).Warn("something went wrong in label policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestStylingFailedEvent, m.view.ProcessedStylingFailedEvent, m.view.ProcessedStylingSequence, m.errorCountUntilSkip)
}

func (m *Styling) OnSuccess(instanceIDs []string) error {
	return spooler.HandleSuccess(m.view.UpdateStylingSpoolerRunTimestamp, instanceIDs)
}

func (m *Styling) generateStylingFile(policy *iam_model.LabelPolicyView) error {
	reader, size, err := m.writeFile(policy)
	if err != nil {
		return err
	}
	return m.uploadFilesToStorage(policy.InstanceID, policy.AggregateID, "text/css", reader, size)
}

func (m *Styling) writeFile(policy *iam_model.LabelPolicyView) (io.Reader, int64, error) {
	cssContent := ""
	cssContent += ":root {"
	if policy.PrimaryColor != "" {
		palette := m.generateColorPaletteRGBA255(policy.PrimaryColor)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-primary-%v: %s;", i, color)
		}
	}

	if policy.BackgroundColor != "" {
		palette := m.generateColorPaletteRGBA255(policy.BackgroundColor)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-background-%v: %s;", i, color)
		}
	}
	if policy.WarnColor != "" {
		palette := m.generateColorPaletteRGBA255(policy.WarnColor)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-warn-%v: %s;", i, color)
		}
	}
	if policy.FontColor != "" {
		palette := m.generateColorPaletteRGBA255(policy.FontColor)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-text-%v: %s;", i, color)
		}
	}
	var fontname string
	if policy.FontURL != "" {
		split := strings.Split(policy.FontURL, "/")
		fontname = split[len(split)-1]
		cssContent += fmt.Sprintf("--zitadel-font-family: %s;", fontname)
	}
	cssContent += "}"
	if policy.FontURL != "" {
		cssContent += fmt.Sprintf(fontFaceTemplate, fontname, login.HandlerPrefix+login.EndpointDynamicResources, policy.AggregateID, policy.FontURL)
	}
	cssContent += ".lgn-dark-theme {"
	if policy.PrimaryColorDark != "" {
		palette := m.generateColorPaletteRGBA255(policy.PrimaryColorDark)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-primary-%v: %s;", i, color)
		}
	}
	if policy.BackgroundColorDark != "" {
		palette := m.generateColorPaletteRGBA255(policy.BackgroundColorDark)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-background-%v: %s;", i, color)
		}
	}
	if policy.WarnColorDark != "" {
		palette := m.generateColorPaletteRGBA255(policy.WarnColorDark)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-warn-%v: %s;", i, color)
		}
	}
	if policy.FontColorDark != "" {
		palette := m.generateColorPaletteRGBA255(policy.FontColorDark)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-text-%v: %s;", i, color)
		}
	}
	cssContent += "}"

	data := []byte(cssContent)
	buffer := bytes.NewBuffer(data)
	return buffer, int64(buffer.Len()), nil
}

const fontFaceTemplate = `
@font-face {
	font-family: '%s';
	font-style: normal;
	font-display: swap;
	src: url(%s?orgId=%s&filename=%s);
}
`

func (m *Styling) uploadFilesToStorage(instanceID, aggregateID, contentType string, reader io.Reader, size int64) error {
	fileName := domain.CssPath + "/" + domain.CssVariablesFileName
	//TODO: handle location as soon as possible
	_, err := m.static.PutObject(context.Background(), instanceID, "", aggregateID, fileName, contentType, static.ObjectTypeStyling, reader, size)
	return err
}

func (m *Styling) deleteInstanceFilesFromStorage(instanceID string) error {
	return m.static.RemoveInstanceObjects(context.Background(), instanceID)
}

func (m *Styling) generateColorPaletteRGBA255(hex string) map[string]string {
	palette := make(map[string]string)
	defaultColor := gamut.Hex(hex)

	color50, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 0.52))
	if ok {
		palette["50"] = cssRGB(color50.RGB255())
	}

	color100, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 0.37))
	if ok {
		palette["100"] = cssRGB(color100.RGB255())
	}

	color200, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 0.26))
	if ok {
		palette["200"] = cssRGB(color200.RGB255())
	}

	color300, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 0.12))
	if ok {
		palette["300"] = cssRGB(color300.RGB255())
	}

	color400, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 0.06))
	if ok {
		palette["400"] = cssRGB(color400.RGB255())
	}

	color500, ok := colorful.MakeColor(defaultColor)
	if ok {
		palette["500"] = cssRGB(color500.RGB255())
	}

	color600, ok := colorful.MakeColor(gamut.Darker(defaultColor, 0.06))
	if ok {
		palette["600"] = cssRGB(color600.RGB255())
	}

	color700, ok := colorful.MakeColor(gamut.Darker(defaultColor, 0.12))
	if ok {
		palette["700"] = cssRGB(color700.RGB255())
	}

	color800, ok := colorful.MakeColor(gamut.Darker(defaultColor, 0.18))
	if ok {
		palette["800"] = cssRGB(color800.RGB255())
	}

	color900, ok := colorful.MakeColor(gamut.Darker(defaultColor, 0.24))
	if ok {
		palette["900"] = cssRGB(color900.RGB255())
	}

	colorContrast, ok := colorful.MakeColor(gamut.Contrast(defaultColor))
	if ok {
		palette["contrast"] = cssRGB(colorContrast.RGB255())
	}

	return palette
}

func cssRGB(r, g, b uint8) string {
	return fmt.Sprintf("rgb(%v, %v, %v)", r, g, b)
}
