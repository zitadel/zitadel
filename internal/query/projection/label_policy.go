package projection

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type labelPolicyProjection struct {
	crdb.StatementHandler
}

const (
	LabelPolicyTable = "zitadel.projections.label_policies"
)

func newLabelPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *labelPolicyProjection {
	p := &labelPolicyProjection{}
	config.ProjectionName = LabelPolicyTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *labelPolicyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.LabelPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.LabelPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.LabelPolicyRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  org.LabelPolicyActivatedEventType,
					Reduce: p.reduceActivated,
				},
				{
					Event:  org.LabelPolicyLogoAddedEventType,
					Reduce: p.reduceLogoAdded,
				},
				{
					Event:  org.LabelPolicyLogoRemovedEventType,
					Reduce: p.reduceLogoRemoved,
				},
				{
					Event:  org.LabelPolicyIconAddedEventType,
					Reduce: p.reduceIconAdded,
				},
				{
					Event:  org.LabelPolicyIconRemovedEventType,
					Reduce: p.reduceIconRemoved,
				},
				{
					Event:  org.LabelPolicyLogoDarkAddedEventType,
					Reduce: p.reduceLogoAdded,
				},
				{
					Event:  org.LabelPolicyLogoDarkRemovedEventType,
					Reduce: p.reduceLogoRemoved,
				},
				{
					Event:  org.LabelPolicyIconDarkAddedEventType,
					Reduce: p.reduceIconAdded,
				},
				{
					Event:  org.LabelPolicyIconDarkRemovedEventType,
					Reduce: p.reduceIconRemoved,
				},
				{
					Event:  org.LabelPolicyFontAddedEventType,
					Reduce: p.reduceFontAdded,
				},
				{
					Event:  org.LabelPolicyFontRemovedEventType,
					Reduce: p.reduceFontRemoved,
				},
				{
					Event:  org.LabelPolicyAssetsRemovedEventType,
					Reduce: p.reduceAssetsRemoved,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.LabelPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.LabelPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  iam.LabelPolicyActivatedEventType,
					Reduce: p.reduceActivated,
				},
				{
					Event:  iam.LabelPolicyLogoAddedEventType,
					Reduce: p.reduceLogoAdded,
				},
				{
					Event:  iam.LabelPolicyLogoRemovedEventType,
					Reduce: p.reduceLogoRemoved,
				},
				{
					Event:  iam.LabelPolicyIconAddedEventType,
					Reduce: p.reduceIconAdded,
				},
				{
					Event:  iam.LabelPolicyIconRemovedEventType,
					Reduce: p.reduceIconRemoved,
				},
				{
					Event:  iam.LabelPolicyLogoDarkAddedEventType,
					Reduce: p.reduceLogoAdded,
				},
				{
					Event:  iam.LabelPolicyLogoDarkRemovedEventType,
					Reduce: p.reduceLogoRemoved,
				},
				{
					Event:  iam.LabelPolicyIconDarkAddedEventType,
					Reduce: p.reduceIconAdded,
				},
				{
					Event:  iam.LabelPolicyIconDarkRemovedEventType,
					Reduce: p.reduceIconRemoved,
				},
				{
					Event:  iam.LabelPolicyFontAddedEventType,
					Reduce: p.reduceFontAdded,
				},
				{
					Event:  iam.LabelPolicyFontRemovedEventType,
					Reduce: p.reduceFontRemoved,
				},
				{
					Event:  iam.LabelPolicyAssetsRemovedEventType,
					Reduce: p.reduceAssetsRemoved,
				},
			},
		},
	}
}

func (p *labelPolicyProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.LabelPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.LabelPolicyAddedEvent:
		policyEvent = e.LabelPolicyAddedEvent
		isDefault = false
	case *iam.LabelPolicyAddedEvent:
		policyEvent = e.LabelPolicyAddedEvent
		isDefault = true
	default:
		logging.LogWithFields("PROJE-zR6h0", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyAddedEventType, iam.LabelPolicyAddedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-CSE7A", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LabelPolicyCreationDateCol, policyEvent.CreationDate()),
			handler.NewCol(LabelPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, policyEvent.Sequence()),
			handler.NewCol(LabelPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCol(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewCol(LabelPolicyIsDefaultCol, isDefault),
			handler.NewCol(LabelPolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(LabelPolicyLightPrimaryColorCol, policyEvent.PrimaryColor),
			handler.NewCol(LabelPolicyLightBackgroundColorCol, policyEvent.BackgroundColor),
			handler.NewCol(LabelPolicyLightWarnColorCol, policyEvent.WarnColor),
			handler.NewCol(LabelPolicyLightFontColorCol, policyEvent.FontColor),
			handler.NewCol(LabelPolicyDarkPrimaryColorCol, policyEvent.PrimaryColorDark),
			handler.NewCol(LabelPolicyDarkBackgroundColorCol, policyEvent.BackgroundColorDark),
			handler.NewCol(LabelPolicyDarkWarnColorCol, policyEvent.WarnColorDark),
			handler.NewCol(LabelPolicyDarkFontColorCol, policyEvent.FontColorDark),
			handler.NewCol(LabelPolicyHideLoginNameSuffixCol, policyEvent.HideLoginNameSuffix),
			handler.NewCol(LabelPolicyShouldErrorPopupCol, policyEvent.ErrorMsgPopup),
			handler.NewCol(LabelPolicyWatermarkDisabledCol, policyEvent.DisableWatermark),
		}), nil
}

func (p *labelPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.LabelPolicyChangedEvent
	switch e := event.(type) {
	case *org.LabelPolicyChangedEvent:
		policyEvent = e.LabelPolicyChangedEvent
	case *iam.LabelPolicyChangedEvent:
		policyEvent = e.LabelPolicyChangedEvent
	default:
		logging.LogWithFields("PROJE-2VrlG", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyChangedEventType, iam.LabelPolicyChangedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-qgVug", "reduce.wrong.event.type")
	}
	cols := []handler.Column{
		handler.NewCol(LabelPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(LabelPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.PrimaryColor != nil {
		cols = append(cols, handler.NewCol(LabelPolicyLightPrimaryColorCol, *policyEvent.PrimaryColor))
	}
	if policyEvent.BackgroundColor != nil {
		cols = append(cols, handler.NewCol(LabelPolicyLightBackgroundColorCol, *policyEvent.BackgroundColor))
	}
	if policyEvent.WarnColor != nil {
		cols = append(cols, handler.NewCol(LabelPolicyLightWarnColorCol, *policyEvent.WarnColor))
	}
	if policyEvent.FontColor != nil {
		cols = append(cols, handler.NewCol(LabelPolicyLightFontColorCol, *policyEvent.FontColor))
	}
	if policyEvent.PrimaryColorDark != nil {
		cols = append(cols, handler.NewCol(LabelPolicyDarkPrimaryColorCol, *policyEvent.PrimaryColorDark))
	}
	if policyEvent.BackgroundColorDark != nil {
		cols = append(cols, handler.NewCol(LabelPolicyDarkBackgroundColorCol, *policyEvent.BackgroundColorDark))
	}
	if policyEvent.WarnColorDark != nil {
		cols = append(cols, handler.NewCol(LabelPolicyDarkWarnColorCol, *policyEvent.WarnColorDark))
	}
	if policyEvent.FontColorDark != nil {
		cols = append(cols, handler.NewCol(LabelPolicyDarkFontColorCol, *policyEvent.FontColorDark))
	}
	if policyEvent.HideLoginNameSuffix != nil {
		cols = append(cols, handler.NewCol(LabelPolicyHideLoginNameSuffixCol, *policyEvent.HideLoginNameSuffix))
	}
	if policyEvent.ErrorMsgPopup != nil {
		cols = append(cols, handler.NewCol(LabelPolicyShouldErrorPopupCol, *policyEvent.ErrorMsgPopup))
	}
	if policyEvent.DisableWatermark != nil {
		cols = append(cols, handler.NewCol(LabelPolicyWatermarkDisabledCol, *policyEvent.DisableWatermark))
	}
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
		}), nil
}

func (p *labelPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.LabelPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-izDbs", "seq", event.Sequence(), "expectedType", org.LabelPolicyRemovedEventType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-ATMBz", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}

func (p *labelPolicyProjection) reduceActivated(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *org.LabelPolicyActivatedEvent, *iam.LabelPolicyActivatedEvent:
		// everything ok
	default:
		logging.LogWithFields("PROJE-ZQO7J", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyActivatedEventType, iam.LabelPolicyActivatedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-dldEU", "reduce.wrong.event.type")
	}
	return crdb.NewCopyStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			handler.NewCol(LabelPolicyStateCol, domain.LabelPolicyStateActive),
			handler.NewCol(LabelPolicyCreationDateCol, nil),
			handler.NewCol(LabelPolicyResourceOwnerCol, nil),
			handler.NewCol(LabelPolicyIDCol, nil),
			handler.NewCol(LabelPolicyIsDefaultCol, nil),
			handler.NewCol(LabelPolicyHideLoginNameSuffixCol, nil),
			handler.NewCol(LabelPolicyFontURLCol, nil),
			handler.NewCol(LabelPolicyWatermarkDisabledCol, nil),
			handler.NewCol(LabelPolicyShouldErrorPopupCol, nil),
			handler.NewCol(LabelPolicyLightPrimaryColorCol, nil),
			handler.NewCol(LabelPolicyLightWarnColorCol, nil),
			handler.NewCol(LabelPolicyLightBackgroundColorCol, nil),
			handler.NewCol(LabelPolicyLightFontColorCol, nil),
			handler.NewCol(LabelPolicyLightLogoURLCol, nil),
			handler.NewCol(LabelPolicyLightIconURLCol, nil),
			handler.NewCol(LabelPolicyDarkPrimaryColorCol, nil),
			handler.NewCol(LabelPolicyDarkWarnColorCol, nil),
			handler.NewCol(LabelPolicyDarkBackgroundColorCol, nil),
			handler.NewCol(LabelPolicyDarkFontColorCol, nil),
			handler.NewCol(LabelPolicyDarkLogoURLCol, nil),
			handler.NewCol(LabelPolicyDarkIconURLCol, nil),
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
		}), nil
}

func (p *labelPolicyProjection) reduceLogoAdded(event eventstore.Event) (*handler.Statement, error) {
	var storeKey handler.Column
	switch e := event.(type) {
	case *org.LabelPolicyLogoAddedEvent:
		storeKey = handler.NewCol(LabelPolicyLightLogoURLCol, e.StoreKey)
	case *iam.LabelPolicyLogoAddedEvent:
		storeKey = handler.NewCol(LabelPolicyLightLogoURLCol, e.StoreKey)
	case *org.LabelPolicyLogoDarkAddedEvent:
		storeKey = handler.NewCol(LabelPolicyDarkLogoURLCol, e.StoreKey)
	case *iam.LabelPolicyLogoDarkAddedEvent:
		storeKey = handler.NewCol(LabelPolicyDarkLogoURLCol, e.StoreKey)
	default:
		logging.LogWithFields("PROJE-NHrbi", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyLogoAddedEventType, iam.LabelPolicyLogoAddedEventType, org.LabelPolicyLogoDarkAddedEventType, iam.LabelPolicyLogoDarkAddedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-4wbOI", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			storeKey,
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
		}), nil
}

func (p *labelPolicyProjection) reduceLogoRemoved(event eventstore.Event) (*handler.Statement, error) {
	var col string
	switch event.(type) {
	case *org.LabelPolicyLogoRemovedEvent:
		col = LabelPolicyLightLogoURLCol
	case *iam.LabelPolicyLogoRemovedEvent:
		col = LabelPolicyLightLogoURLCol
	case *org.LabelPolicyLogoDarkRemovedEvent:
		col = LabelPolicyDarkLogoURLCol
	case *iam.LabelPolicyLogoDarkRemovedEvent:
		col = LabelPolicyDarkLogoURLCol
	default:
		logging.LogWithFields("PROJE-oUmnS", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyLogoRemovedEventType, iam.LabelPolicyLogoRemovedEventType, org.LabelPolicyLogoDarkRemovedEventType, iam.LabelPolicyLogoDarkRemovedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-kg8H4", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			handler.NewCol(col, nil),
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
		}), nil
}

func (p *labelPolicyProjection) reduceIconAdded(event eventstore.Event) (*handler.Statement, error) {
	var storeKey handler.Column
	switch e := event.(type) {
	case *org.LabelPolicyIconAddedEvent:
		storeKey = handler.NewCol(LabelPolicyLightIconURLCol, e.StoreKey)
	case *iam.LabelPolicyIconAddedEvent:
		storeKey = handler.NewCol(LabelPolicyLightIconURLCol, e.StoreKey)
	case *org.LabelPolicyIconDarkAddedEvent:
		storeKey = handler.NewCol(LabelPolicyDarkIconURLCol, e.StoreKey)
	case *iam.LabelPolicyIconDarkAddedEvent:
		storeKey = handler.NewCol(LabelPolicyDarkIconURLCol, e.StoreKey)
	default:
		logging.LogWithFields("PROJE-6efFw", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyIconAddedEventType, iam.LabelPolicyIconAddedEventType, org.LabelPolicyIconDarkAddedEventType, iam.LabelPolicyIconDarkAddedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-e2JFz", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			storeKey,
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
		}), nil
}

func (p *labelPolicyProjection) reduceIconRemoved(event eventstore.Event) (*handler.Statement, error) {
	var col string
	switch event.(type) {
	case *org.LabelPolicyIconRemovedEvent:
		col = LabelPolicyLightIconURLCol
	case *iam.LabelPolicyIconRemovedEvent:
		col = LabelPolicyLightIconURLCol
	case *org.LabelPolicyIconDarkRemovedEvent:
		col = LabelPolicyDarkIconURLCol
	case *iam.LabelPolicyIconDarkRemovedEvent:
		col = LabelPolicyDarkIconURLCol
	default:
		logging.LogWithFields("PROJE-0BiAZ", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyIconRemovedEventType, iam.LabelPolicyIconRemovedEventType, org.LabelPolicyIconDarkRemovedEventType, iam.LabelPolicyIconDarkRemovedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-gfgbY", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			handler.NewCol(col, nil),
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
		}), nil
}

func (p *labelPolicyProjection) reduceFontAdded(event eventstore.Event) (*handler.Statement, error) {
	var storeKey handler.Column
	switch e := event.(type) {
	case *org.LabelPolicyFontAddedEvent:
		storeKey = handler.NewCol(LabelPolicyFontURLCol, e.StoreKey)
	case *iam.LabelPolicyFontAddedEvent:
		storeKey = handler.NewCol(LabelPolicyFontURLCol, e.StoreKey)
	default:
		logging.LogWithFields("PROJE-DCzfX", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyFontAddedEventType, iam.LabelPolicyFontAddedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-65i9W", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			storeKey,
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
		}), nil
}

func (p *labelPolicyProjection) reduceFontRemoved(event eventstore.Event) (*handler.Statement, error) {
	var col string
	switch event.(type) {
	case *org.LabelPolicyFontRemovedEvent:
		col = LabelPolicyFontURLCol
	case *iam.LabelPolicyFontRemovedEvent:
		col = LabelPolicyFontURLCol
	default:
		logging.LogWithFields("PROJE-YKwG4", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyFontRemovedEventType, iam.LabelPolicyFontRemovedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-xf32J", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			handler.NewCol(col, nil),
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
		}), nil
}

func (p *labelPolicyProjection) reduceAssetsRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *org.LabelPolicyAssetsRemovedEvent, *iam.LabelPolicyAssetsRemovedEvent:
		//ok
	default:
		logging.LogWithFields("PROJE-YKwG4", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LabelPolicyAssetsRemovedEventType, iam.LabelPolicyAssetsRemovedEventType}).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-qi39A", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			handler.NewCol(LabelPolicyLightLogoURLCol, nil),
			handler.NewCol(LabelPolicyLightIconURLCol, nil),
			handler.NewCol(LabelPolicyDarkLogoURLCol, nil),
			handler.NewCol(LabelPolicyDarkIconURLCol, nil),
			handler.NewCol(LabelPolicyFontURLCol, nil),
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
		}), nil
}

const (
	LabelPolicyCreationDateCol        = "creation_date"
	LabelPolicyChangeDateCol          = "change_date"
	LabelPolicySequenceCol            = "sequence"
	LabelPolicyIDCol                  = "id"
	LabelPolicyStateCol               = "state"
	LabelPolicyIsDefaultCol           = "is_default"
	LabelPolicyResourceOwnerCol       = "resource_owner"
	LabelPolicyHideLoginNameSuffixCol = "hide_login_name_suffix"
	LabelPolicyFontURLCol             = "font_url"
	LabelPolicyWatermarkDisabledCol   = "watermark_disabled"
	LabelPolicyShouldErrorPopupCol    = "should_error_popup"

	LabelPolicyLightPrimaryColorCol    = "light_primary_color"
	LabelPolicyLightWarnColorCol       = "light_warn_color"
	LabelPolicyLightBackgroundColorCol = "light_background_color"
	LabelPolicyLightFontColorCol       = "light_font_color"
	LabelPolicyLightLogoURLCol         = "light_logo_url"
	LabelPolicyLightIconURLCol         = "light_icon_url"

	LabelPolicyDarkPrimaryColorCol    = "dark_primary_color"
	LabelPolicyDarkWarnColorCol       = "dark_warn_color"
	LabelPolicyDarkBackgroundColorCol = "dark_background_color"
	LabelPolicyDarkFontColorCol       = "dark_font_color"
	LabelPolicyDarkLogoURLCol         = "dark_logo_url"
	LabelPolicyDarkIconURLCol         = "dark_icon_url"
)
