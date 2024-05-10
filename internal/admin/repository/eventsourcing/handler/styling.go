package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"

	admin_view "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/static"
)

const (
	stylingTable = "adminapi.styling2"
)

var _ handler.Projection = (*Styling)(nil)

type Styling struct {
	static static.Storage
	view   *admin_view.View
}

func newStyling(ctx context.Context, config handler.Config, static static.Storage, view *admin_view.View) *handler.Handler {
	return handler.NewHandler(
		ctx,
		&config,
		&Styling{
			static: static,
			view:   view,
		},
	)
}

// Name implements [handler.Projection]
func (*Styling) Name() string {
	return stylingTable
}

// Reducers implements [handler.Projection]
func (s *Styling) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.LabelPolicyAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyChangedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyLogoAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyLogoRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyIconAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyIconRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyLogoDarkAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyLogoDarkRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyIconDarkAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyIconDarkRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyFontAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyFontRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyAssetsRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.LabelPolicyActivatedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.LabelPolicyAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyChangedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyLogoAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyLogoRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyIconAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyIconRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyLogoDarkAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyLogoDarkRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyIconDarkAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyIconDarkRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyFontAddedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyFontRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyAssetsRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.LabelPolicyActivatedEventType,
					Reduce: s.processLabelPolicy,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: s.processLabelPolicy,
				},
			},
		},
	}
}

func (m *Styling) processLabelPolicy(event eventstore.Event) (_ *handler.Statement, err error) {
	policy := &model.LabelPolicyView{
		AggregateID: event.Aggregate().ID,
		InstanceID:  event.Aggregate().InstanceID,
		State:       int32(domain.LabelPolicyStatePreview),
	}
	switch event.Type() {
	case instance.LabelPolicyAddedEventType,
		org.LabelPolicyAddedEventType,
		instance.LabelPolicyChangedEventType,
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

		err = policy.AppendEvent(event)
		if err != nil {
			return nil, err
		}
		changes := policy.Changes()
		if len(changes) == 0 {
			return handler.NewNoOpStatement(event), nil
		}
		return handler.NewUpsertStatement(
			event,
			policy.PKColumns(),
			append(policy.PKColumns(), changes...),
		), nil
	case instance.LabelPolicyActivatedEventType,
		org.LabelPolicyActivatedEventType:

		err = policy.AppendEvent(event)
		if err != nil {
			return nil, err
		}
		changes := policy.Changes()
		if len(changes) == 0 {
			return handler.NewNoOpStatement(event), nil
		}
		return handler.NewMultiStatement(
			event,
			handler.AddStatement(func(ex handler.Executer, projectionName string) error {
				return m.generateStylingFile(policy)
			}),
			handler.AddCopyStatement(
				[]handler.Column{
					handler.NewCol(model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyInstanceID).ToColumnName(), nil),
					handler.NewCol(model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyAggregateID).ToColumnName(), nil),
					handler.NewCol(model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyState).ToColumnName(), nil),
				},
				[]handler.Column{
					handler.NewCol(model.LabelPolicyKeyChangeDate, event.CreatedAt()),
					handler.NewCol(model.LabelPolicyKeySequence, event.Sequence()),
					handler.NewCol(model.LabelPolicyKeyState, domain.LabelPolicyStateActive),
					handler.NewCol(model.LabelPolicyKeyCreationDate, nil),
					handler.NewCol(model.LabelPolicyKeyInstanceID, nil),
					handler.NewCol(model.LabelPolicyKeyAggregateID, nil),
					handler.NewCol(model.LabelPolicyKeyHideLoginNameSuffix, nil),
					handler.NewCol(model.LabelPolicyKeyFontURL, nil),
					handler.NewCol(model.LabelPolicyKeyDisableWatermark, nil),
					handler.NewCol(model.LabelPolicyKeyErrorMsgPopup, nil),
					handler.NewCol(model.LabelPolicyKeyPrimaryColor, nil),
					handler.NewCol(model.LabelPolicyKeyWarnColor, nil),
					handler.NewCol(model.LabelPolicyKeyBackgroundColor, nil),
					handler.NewCol(model.LabelPolicyKeyFontColor, nil),
					handler.NewCol(model.LabelPolicyKeyLogoURL, nil),
					handler.NewCol(model.LabelPolicyKeyIconURL, nil),
					handler.NewCol(model.LabelPolicyKeyPrimaryColorDark, nil),
					handler.NewCol(model.LabelPolicyKeyWarnColorDark, nil),
					handler.NewCol(model.LabelPolicyKeyBackgroundColorDark, nil),
					handler.NewCol(model.LabelPolicyKeyFontColorDark, nil),
					handler.NewCol(model.LabelPolicyKeyLogoDarkURL, nil),
					handler.NewCol(model.LabelPolicyKeyIconDarkURL, nil),
				},
				[]handler.Column{
					handler.NewCol(model.LabelPolicyKeyChangeDate, nil),
					handler.NewCol(model.LabelPolicyKeySequence, nil),
					handler.NewCol(model.LabelPolicyKeyState, nil),
					handler.NewCol(model.LabelPolicyKeyCreationDate, nil),
					handler.NewCol(model.LabelPolicyKeyInstanceID, nil),
					handler.NewCol(model.LabelPolicyKeyAggregateID, nil),
					handler.NewCol(model.LabelPolicyKeyHideLoginNameSuffix, nil),
					handler.NewCol(model.LabelPolicyKeyFontURL, nil),
					handler.NewCol(model.LabelPolicyKeyDisableWatermark, nil),
					handler.NewCol(model.LabelPolicyKeyErrorMsgPopup, nil),
					handler.NewCol(model.LabelPolicyKeyPrimaryColor, nil),
					handler.NewCol(model.LabelPolicyKeyWarnColor, nil),
					handler.NewCol(model.LabelPolicyKeyBackgroundColor, nil),
					handler.NewCol(model.LabelPolicyKeyFontColor, nil),
					handler.NewCol(model.LabelPolicyKeyLogoURL, nil),
					handler.NewCol(model.LabelPolicyKeyIconURL, nil),
					handler.NewCol(model.LabelPolicyKeyPrimaryColorDark, nil),
					handler.NewCol(model.LabelPolicyKeyWarnColorDark, nil),
					handler.NewCol(model.LabelPolicyKeyBackgroundColorDark, nil),
					handler.NewCol(model.LabelPolicyKeyFontColorDark, nil),
					handler.NewCol(model.LabelPolicyKeyLogoDarkURL, nil),
					handler.NewCol(model.LabelPolicyKeyIconDarkURL, nil),
				},
				[]handler.NamespacedCondition{
					handler.NewNamespacedCondition(model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyAggregateID).ToColumnName(), event.Aggregate().ID),
					handler.NewNamespacedCondition(model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyState).ToColumnName(), domain.LabelPolicyStatePreview),
					handler.NewNamespacedCondition(model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyInstanceID).ToColumnName(), event.Aggregate().InstanceID),
				},
			),
			handler.AddUpsertStatement(
				policy.PKColumns(),
				append(policy.PKColumns(), changes...),
			),
		), nil
	case instance.InstanceRemovedEventType:
		return handler.NewMultiStatement(
			event,
			handler.AddStatement(func(ex handler.Executer, projectionName string) error {
				return m.deleteInstanceFilesFromStorage(event.Aggregate().InstanceID)
			}),
			handler.AddDeleteStatement(
				[]handler.Condition{
					handler.NewCond(
						model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyInstanceID).ToColumnName(),
						event.Aggregate().InstanceID,
					),
				},
			),
		), nil
	case org.OrgRemovedEventType:
		return handler.NewUpdateStatement(
			event,
			[]handler.Column{
				handler.NewCol(
					model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyOwnerRemoved).ToColumnName(),
					true,
				),
			},
			[]handler.Condition{
				handler.NewCond(
					model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyInstanceID).ToColumnName(),
					event.Aggregate().InstanceID,
				),
				handler.NewCond(
					model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyAggregateID).ToColumnName(),
					event.Aggregate().ID,
				),
			},
		), nil
	default:
		return handler.NewNoOpStatement(event), nil
	}
}

func (m *Styling) generateStylingFile(policy *model.LabelPolicyView) error {
	reader, size, err := m.writeFile(policy)
	if err != nil {
		return err
	}
	return m.uploadFilesToStorage(policy.InstanceID, policy.AggregateID, "text/css", reader, size)
}

func (m *Styling) writeFile(policy *model.LabelPolicyView) (io.Reader, int64, error) {
	cssContent := ""
	cssContent += ":root {"
	if policy.PrimaryColor.Value() != "" {
		palette := m.generateColorPaletteRGBA255(policy.PrimaryColor.Value())
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-primary-%v: %s;", i, color)
		}
	}

	if policy.BackgroundColor.Value() != "" {
		palette := m.generateColorPaletteRGBA255(policy.BackgroundColor.Value())
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-background-%v: %s;", i, color)
		}
	}
	if policy.WarnColor.Value() != "" {
		palette := m.generateColorPaletteRGBA255(policy.WarnColor.Value())
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-warn-%v: %s;", i, color)
		}
	}
	if policy.FontColor.Value() != "" {
		cssContent += fmt.Sprintf("--zitadel-color-label: %s;", policy.FontColor.Value())
		palette := m.generateColorPaletteRGBA255(policy.FontColor.Value())
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-text-%v: %s;", i, color)
		}
	}
	var fontName string
	if policy.FontURL.Value() != "" {
		split := strings.Split(policy.FontURL.Value(), "/")
		fontName = split[len(split)-1]
		cssContent += fmt.Sprintf("--zitadel-font-family: %s;", fontName)
	}
	cssContent += "}"
	if policy.FontURL.Value() != "" {
		cssContent += fmt.Sprintf(fontFaceTemplate, fontName, login.HandlerPrefix+login.EndpointDynamicResources, policy.AggregateID, policy.FontURL.Value())
	}
	cssContent += ".lgn-dark-theme {"
	if policy.PrimaryColorDark.Value() != "" {
		palette := m.generateColorPaletteRGBA255(policy.PrimaryColorDark.Value())
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-primary-%v: %s;", i, color)
		}
	}
	if policy.BackgroundColorDark.Value() != "" {
		palette := m.generateColorPaletteRGBA255(policy.BackgroundColorDark.Value())
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-background-%v: %s;", i, color)
		}
	}
	if policy.WarnColorDark.Value() != "" {
		palette := m.generateColorPaletteRGBA255(policy.WarnColorDark.Value())
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-warn-%v: %s;", i, color)
		}
	}
	if policy.FontColorDark.Value() != "" {
		cssContent += fmt.Sprintf("--zitadel-color-label: %s;", policy.FontColorDark.Value())
		palette := m.generateColorPaletteRGBA255(policy.FontColorDark.Value())
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
