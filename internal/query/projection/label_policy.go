package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	LabelPolicyTable = "projections.label_policies3"

	LabelPolicyIDCol                  = "id"
	LabelPolicyCreationDateCol        = "creation_date"
	LabelPolicyChangeDateCol          = "change_date"
	LabelPolicySequenceCol            = "sequence"
	LabelPolicyStateCol               = "state"
	LabelPolicyIsDefaultCol           = "is_default"
	LabelPolicyResourceOwnerCol       = "resource_owner"
	LabelPolicyInstanceIDCol          = "instance_id"
	LabelPolicyHideLoginNameSuffixCol = "hide_login_name_suffix"
	LabelPolicyWatermarkDisabledCol   = "watermark_disabled"
	LabelPolicyShouldErrorPopupCol    = "should_error_popup"
	LabelPolicyFontURLCol             = "font_url"
	LabelPolicyOwnerRemovedCol        = "owner_removed"
	LabelPolicyThemeModeCol           = "theme_mode"

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

type labelPolicyProjection struct{}

func newLabelPolicyProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(labelPolicyProjection))
}

func (*labelPolicyProjection) Name() string {
	return LabelPolicyTable
}

func (*labelPolicyProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(LabelPolicyIDCol, handler.ColumnTypeText),
			handler.NewColumn(LabelPolicyCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(LabelPolicyChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(LabelPolicySequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(LabelPolicyStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(LabelPolicyIsDefaultCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(LabelPolicyResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(LabelPolicyInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(LabelPolicyHideLoginNameSuffixCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(LabelPolicyWatermarkDisabledCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(LabelPolicyShouldErrorPopupCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(LabelPolicyFontURLCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyLightPrimaryColorCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyLightWarnColorCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyLightBackgroundColorCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyLightFontColorCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyLightLogoURLCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyLightIconURLCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyDarkPrimaryColorCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyDarkWarnColorCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyDarkBackgroundColorCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyDarkFontColorCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyDarkLogoURLCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyDarkIconURLCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LabelPolicyOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(LabelPolicyThemeModeCol, handler.ColumnTypeEnum, handler.Default(0)),
		},
			handler.NewPrimaryKey(LabelPolicyInstanceIDCol, LabelPolicyIDCol, LabelPolicyStateCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{LabelPolicyOwnerRemovedCol})),
		),
	)
}

func (p *labelPolicyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
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
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.LabelPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.LabelPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  instance.LabelPolicyActivatedEventType,
					Reduce: p.reduceActivated,
				},
				{
					Event:  instance.LabelPolicyLogoAddedEventType,
					Reduce: p.reduceLogoAdded,
				},
				{
					Event:  instance.LabelPolicyLogoRemovedEventType,
					Reduce: p.reduceLogoRemoved,
				},
				{
					Event:  instance.LabelPolicyIconAddedEventType,
					Reduce: p.reduceIconAdded,
				},
				{
					Event:  instance.LabelPolicyIconRemovedEventType,
					Reduce: p.reduceIconRemoved,
				},
				{
					Event:  instance.LabelPolicyLogoDarkAddedEventType,
					Reduce: p.reduceLogoAdded,
				},
				{
					Event:  instance.LabelPolicyLogoDarkRemovedEventType,
					Reduce: p.reduceLogoRemoved,
				},
				{
					Event:  instance.LabelPolicyIconDarkAddedEventType,
					Reduce: p.reduceIconAdded,
				},
				{
					Event:  instance.LabelPolicyIconDarkRemovedEventType,
					Reduce: p.reduceIconRemoved,
				},
				{
					Event:  instance.LabelPolicyFontAddedEventType,
					Reduce: p.reduceFontAdded,
				},
				{
					Event:  instance.LabelPolicyFontRemovedEventType,
					Reduce: p.reduceFontRemoved,
				},
				{
					Event:  instance.LabelPolicyAssetsRemovedEventType,
					Reduce: p.reduceAssetsRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(LabelPolicyInstanceIDCol),
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
	case *instance.LabelPolicyAddedEvent:
		policyEvent = e.LabelPolicyAddedEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-CSE7A", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyAddedEventType, instance.LabelPolicyAddedEventType})
	}
	return handler.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LabelPolicyCreationDateCol, policyEvent.CreatedAt()),
			handler.NewCol(LabelPolicyChangeDateCol, policyEvent.CreatedAt()),
			handler.NewCol(LabelPolicySequenceCol, policyEvent.Sequence()),
			handler.NewCol(LabelPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCol(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewCol(LabelPolicyIsDefaultCol, isDefault),
			handler.NewCol(LabelPolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(LabelPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
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
			handler.NewCol(LabelPolicyThemeModeCol, policyEvent.ThemeMode),
		}), nil
}

func (p *labelPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.LabelPolicyChangedEvent
	switch e := event.(type) {
	case *org.LabelPolicyChangedEvent:
		policyEvent = e.LabelPolicyChangedEvent
	case *instance.LabelPolicyChangedEvent:
		policyEvent = e.LabelPolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyChangedEventType, instance.LabelPolicyChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(LabelPolicyChangeDateCol, policyEvent.CreatedAt()),
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
	if policyEvent.ThemeMode != nil {
		cols = append(cols, handler.NewCol(LabelPolicyThemeModeCol, *policyEvent.ThemeMode))
	}
	return handler.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.LabelPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-ATMBz", "reduce.wrong.event.type %s", org.LabelPolicyRemovedEventType)
	}
	return handler.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceActivated(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *org.LabelPolicyActivatedEvent, *instance.LabelPolicyActivatedEvent:
		// everything ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-dldEU", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyActivatedEventType, instance.LabelPolicyActivatedEventType})
	}
	return handler.NewCopyStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyInstanceIDCol, nil),
			handler.NewCol(LabelPolicyIDCol, nil),
			handler.NewCol(LabelPolicyStateCol, nil),
		},
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			handler.NewCol(LabelPolicyStateCol, domain.LabelPolicyStateActive),
			handler.NewCol(LabelPolicyCreationDateCol, nil),
			handler.NewCol(LabelPolicyResourceOwnerCol, nil),
			handler.NewCol(LabelPolicyInstanceIDCol, nil),
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
			handler.NewCol(LabelPolicyThemeModeCol, nil),
		},
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, nil),
			handler.NewCol(LabelPolicySequenceCol, nil),
			handler.NewCol(LabelPolicyStateCol, nil),
			handler.NewCol(LabelPolicyCreationDateCol, nil),
			handler.NewCol(LabelPolicyResourceOwnerCol, nil),
			handler.NewCol(LabelPolicyInstanceIDCol, nil),
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
			handler.NewCol(LabelPolicyThemeModeCol, nil),
		},
		[]handler.NamespacedCondition{
			handler.NewNamespacedCondition(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewNamespacedCondition(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewNamespacedCondition(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceLogoAdded(event eventstore.Event) (*handler.Statement, error) {
	var storeKey handler.Column
	switch e := event.(type) {
	case *org.LabelPolicyLogoAddedEvent:
		storeKey = handler.NewCol(LabelPolicyLightLogoURLCol, e.StoreKey)
	case *instance.LabelPolicyLogoAddedEvent:
		storeKey = handler.NewCol(LabelPolicyLightLogoURLCol, e.StoreKey)
	case *org.LabelPolicyLogoDarkAddedEvent:
		storeKey = handler.NewCol(LabelPolicyDarkLogoURLCol, e.StoreKey)
	case *instance.LabelPolicyLogoDarkAddedEvent:
		storeKey = handler.NewCol(LabelPolicyDarkLogoURLCol, e.StoreKey)
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-4wbOI", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoAddedEventType, instance.LabelPolicyLogoAddedEventType, org.LabelPolicyLogoDarkAddedEventType, instance.LabelPolicyLogoDarkAddedEventType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			storeKey,
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceLogoRemoved(event eventstore.Event) (*handler.Statement, error) {
	var col string
	switch event.(type) {
	case *org.LabelPolicyLogoRemovedEvent:
		col = LabelPolicyLightLogoURLCol
	case *instance.LabelPolicyLogoRemovedEvent:
		col = LabelPolicyLightLogoURLCol
	case *org.LabelPolicyLogoDarkRemovedEvent:
		col = LabelPolicyDarkLogoURLCol
	case *instance.LabelPolicyLogoDarkRemovedEvent:
		col = LabelPolicyDarkLogoURLCol
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoRemovedEventType, instance.LabelPolicyLogoRemovedEventType, org.LabelPolicyLogoDarkRemovedEventType, instance.LabelPolicyLogoDarkRemovedEventType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			handler.NewCol(col, nil),
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceIconAdded(event eventstore.Event) (*handler.Statement, error) {
	var storeKey handler.Column
	switch e := event.(type) {
	case *org.LabelPolicyIconAddedEvent:
		storeKey = handler.NewCol(LabelPolicyLightIconURLCol, e.StoreKey)
	case *instance.LabelPolicyIconAddedEvent:
		storeKey = handler.NewCol(LabelPolicyLightIconURLCol, e.StoreKey)
	case *org.LabelPolicyIconDarkAddedEvent:
		storeKey = handler.NewCol(LabelPolicyDarkIconURLCol, e.StoreKey)
	case *instance.LabelPolicyIconDarkAddedEvent:
		storeKey = handler.NewCol(LabelPolicyDarkIconURLCol, e.StoreKey)
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-e2JFz", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconAddedEventType, instance.LabelPolicyIconAddedEventType, org.LabelPolicyIconDarkAddedEventType, instance.LabelPolicyIconDarkAddedEventType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			storeKey,
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceIconRemoved(event eventstore.Event) (*handler.Statement, error) {
	var col string
	switch event.(type) {
	case *org.LabelPolicyIconRemovedEvent:
		col = LabelPolicyLightIconURLCol
	case *instance.LabelPolicyIconRemovedEvent:
		col = LabelPolicyLightIconURLCol
	case *org.LabelPolicyIconDarkRemovedEvent:
		col = LabelPolicyDarkIconURLCol
	case *instance.LabelPolicyIconDarkRemovedEvent:
		col = LabelPolicyDarkIconURLCol
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gfgbY", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconRemovedEventType, instance.LabelPolicyIconRemovedEventType, org.LabelPolicyIconDarkRemovedEventType, instance.LabelPolicyIconDarkRemovedEventType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			handler.NewCol(col, nil),
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceFontAdded(event eventstore.Event) (*handler.Statement, error) {
	var storeKey handler.Column
	switch e := event.(type) {
	case *org.LabelPolicyFontAddedEvent:
		storeKey = handler.NewCol(LabelPolicyFontURLCol, e.StoreKey)
	case *instance.LabelPolicyFontAddedEvent:
		storeKey = handler.NewCol(LabelPolicyFontURLCol, e.StoreKey)
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-65i9W", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontAddedEventType, instance.LabelPolicyFontAddedEventType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			storeKey,
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceFontRemoved(event eventstore.Event) (*handler.Statement, error) {
	var col string
	switch event.(type) {
	case *org.LabelPolicyFontRemovedEvent:
		col = LabelPolicyFontURLCol
	case *instance.LabelPolicyFontRemovedEvent:
		col = LabelPolicyFontURLCol
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-xf32J", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontRemovedEventType, instance.LabelPolicyFontRemovedEventType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
			handler.NewCol(col, nil),
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceAssetsRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *org.LabelPolicyAssetsRemovedEvent, *instance.LabelPolicyAssetsRemovedEvent:
		//ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qi39A", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyAssetsRemovedEventType, instance.LabelPolicyAssetsRemovedEventType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
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
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Su6pX", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(LabelPolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(LabelPolicyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
