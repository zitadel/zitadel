package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

const (
	LabelPolicyTable = "projections.label_policies2"

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

type labelPolicyProjection struct {
	crdb.StatementHandler
}

func newLabelPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *labelPolicyProjection {
	p := new(labelPolicyProjection)
	config.ProjectionName = LabelPolicyTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(LabelPolicyIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LabelPolicyCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(LabelPolicyChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(LabelPolicySequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(LabelPolicyStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(LabelPolicyIsDefaultCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(LabelPolicyResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(LabelPolicyInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LabelPolicyHideLoginNameSuffixCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(LabelPolicyWatermarkDisabledCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(LabelPolicyShouldErrorPopupCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(LabelPolicyFontURLCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyLightPrimaryColorCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyLightWarnColorCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyLightBackgroundColorCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyLightFontColorCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyLightLogoURLCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyLightIconURLCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyDarkPrimaryColorCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyDarkWarnColorCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyDarkBackgroundColorCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyDarkFontColorCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyDarkLogoURLCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyDarkIconURLCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LabelPolicyOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(LabelPolicyInstanceIDCol, LabelPolicyIDCol, LabelPolicyStateCol),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{LabelPolicyOwnerRemovedCol})),
		),
	)
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
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-CSE7A", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyAddedEventType, instance.LabelPolicyAddedEventType})
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyChangedEventType, instance.LabelPolicyChangedEventType})
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
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.LabelPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-ATMBz", "reduce.wrong.event.type %s", org.LabelPolicyRemovedEventType)
	}
	return crdb.NewDeleteStatement(
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-dldEU", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyActivatedEventType, instance.LabelPolicyActivatedEventType})
	}
	return crdb.NewCopyStatement(
		event,
		[]handler.Column{
			handler.NewCol(LabelPolicyInstanceIDCol, nil),
			handler.NewCol(LabelPolicyIDCol, nil),
			handler.NewCol(LabelPolicyStateCol, nil),
		},
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, event.CreationDate()),
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-4wbOI", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoAddedEventType, instance.LabelPolicyLogoAddedEventType, org.LabelPolicyLogoDarkAddedEventType, instance.LabelPolicyLogoDarkAddedEventType})
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoRemovedEventType, instance.LabelPolicyLogoRemovedEventType, org.LabelPolicyLogoDarkRemovedEventType, instance.LabelPolicyLogoDarkRemovedEventType})
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-e2JFz", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconAddedEventType, instance.LabelPolicyIconAddedEventType, org.LabelPolicyIconDarkAddedEventType, instance.LabelPolicyIconDarkAddedEventType})
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-gfgbY", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconRemovedEventType, instance.LabelPolicyIconRemovedEventType, org.LabelPolicyIconDarkRemovedEventType, instance.LabelPolicyIconDarkRemovedEventType})
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-65i9W", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontAddedEventType, instance.LabelPolicyFontAddedEventType})
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-xf32J", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontRemovedEventType, instance.LabelPolicyFontRemovedEventType})
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
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceAssetsRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *org.LabelPolicyAssetsRemovedEvent, *instance.LabelPolicyAssetsRemovedEvent:
		//ok
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-qi39A", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyAssetsRemovedEventType, instance.LabelPolicyAssetsRemovedEventType})
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
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *labelPolicyProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Su6pX", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(LabelPolicyChangeDateCol, e.CreationDate()),
			handler.NewCol(LabelPolicySequenceCol, e.Sequence()),
			handler.NewCol(LabelPolicyOwnerRemovedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(LabelPolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(LabelPolicyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
