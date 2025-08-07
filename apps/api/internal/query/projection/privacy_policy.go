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
	PrivacyPolicyTable = "projections.privacy_policies4"

	PrivacyPolicyIDCol             = "id"
	PrivacyPolicyCreationDateCol   = "creation_date"
	PrivacyPolicyChangeDateCol     = "change_date"
	PrivacyPolicySequenceCol       = "sequence"
	PrivacyPolicyStateCol          = "state"
	PrivacyPolicyIsDefaultCol      = "is_default"
	PrivacyPolicyResourceOwnerCol  = "resource_owner"
	PrivacyPolicyInstanceIDCol     = "instance_id"
	PrivacyPolicyPrivacyLinkCol    = "privacy_link"
	PrivacyPolicyTOSLinkCol        = "tos_link"
	PrivacyPolicyHelpLinkCol       = "help_link"
	PrivacyPolicySupportEmailCol   = "support_email"
	PrivacyPolicyDocsLinkCol       = "docs_link"
	PrivacyPolicyCustomLinkCol     = "custom_link"
	PrivacyPolicyCustomLinkTextCol = "custom_link_text"
	PrivacyPolicyOwnerRemovedCol   = "owner_removed"
)

type privacyPolicyProjection struct{}

func newPrivacyPolicyProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(privacyPolicyProjection))
}

func (*privacyPolicyProjection) Name() string {
	return PrivacyPolicyTable
}

func (*privacyPolicyProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(PrivacyPolicyIDCol, handler.ColumnTypeText),
			handler.NewColumn(PrivacyPolicyCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(PrivacyPolicyChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(PrivacyPolicySequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(PrivacyPolicyStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(PrivacyPolicyIsDefaultCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(PrivacyPolicyResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(PrivacyPolicyInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(PrivacyPolicyPrivacyLinkCol, handler.ColumnTypeText),
			handler.NewColumn(PrivacyPolicyTOSLinkCol, handler.ColumnTypeText),
			handler.NewColumn(PrivacyPolicyHelpLinkCol, handler.ColumnTypeText),
			handler.NewColumn(PrivacyPolicySupportEmailCol, handler.ColumnTypeText),
			handler.NewColumn(PrivacyPolicyDocsLinkCol, handler.ColumnTypeText, handler.Default("https://zitadel.com/docs")),
			handler.NewColumn(PrivacyPolicyCustomLinkCol, handler.ColumnTypeText),
			handler.NewColumn(PrivacyPolicyCustomLinkTextCol, handler.ColumnTypeText),
			handler.NewColumn(PrivacyPolicyOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(PrivacyPolicyInstanceIDCol, PrivacyPolicyIDCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{PrivacyPolicyOwnerRemovedCol})),
		),
	)
}

func (p *privacyPolicyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.PrivacyPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.PrivacyPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.PrivacyPolicyRemovedEventType,
					Reduce: p.reduceRemoved,
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
					Event:  instance.PrivacyPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.PrivacyPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(PrivacyPolicyInstanceIDCol),
				},
			},
		},
	}
}

func (p *privacyPolicyProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PrivacyPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.PrivacyPolicyAddedEvent:
		policyEvent = e.PrivacyPolicyAddedEvent
		isDefault = false
	case *instance.PrivacyPolicyAddedEvent:
		policyEvent = e.PrivacyPolicyAddedEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kRNh8", "reduce.wrong.event.type %v", []eventstore.EventType{org.PrivacyPolicyAddedEventType, instance.PrivacyPolicyAddedEventType})
	}
	return handler.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(PrivacyPolicyCreationDateCol, policyEvent.CreationDate()),
			handler.NewCol(PrivacyPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(PrivacyPolicySequenceCol, policyEvent.Sequence()),
			handler.NewCol(PrivacyPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCol(PrivacyPolicyStateCol, domain.PolicyStateActive),
			handler.NewCol(PrivacyPolicyPrivacyLinkCol, policyEvent.PrivacyLink),
			handler.NewCol(PrivacyPolicyTOSLinkCol, policyEvent.TOSLink),
			handler.NewCol(PrivacyPolicyHelpLinkCol, policyEvent.HelpLink),
			handler.NewCol(PrivacyPolicySupportEmailCol, policyEvent.SupportEmail),
			handler.NewCol(PrivacyPolicyDocsLinkCol, policyEvent.DocsLink),
			handler.NewCol(PrivacyPolicyCustomLinkCol, policyEvent.CustomLink),
			handler.NewCol(PrivacyPolicyCustomLinkTextCol, policyEvent.CustomLinkText),
			handler.NewCol(PrivacyPolicyIsDefaultCol, isDefault),
			handler.NewCol(PrivacyPolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(PrivacyPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *privacyPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PrivacyPolicyChangedEvent
	switch e := event.(type) {
	case *org.PrivacyPolicyChangedEvent:
		policyEvent = e.PrivacyPolicyChangedEvent
	case *instance.PrivacyPolicyChangedEvent:
		policyEvent = e.PrivacyPolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-91weZ", "reduce.wrong.event.type %v", []eventstore.EventType{org.PrivacyPolicyChangedEventType, instance.PrivacyPolicyChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(PrivacyPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(PrivacyPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.PrivacyLink != nil {
		cols = append(cols, handler.NewCol(PrivacyPolicyPrivacyLinkCol, *policyEvent.PrivacyLink))
	}
	if policyEvent.TOSLink != nil {
		cols = append(cols, handler.NewCol(PrivacyPolicyTOSLinkCol, *policyEvent.TOSLink))
	}
	if policyEvent.HelpLink != nil {
		cols = append(cols, handler.NewCol(PrivacyPolicyHelpLinkCol, *policyEvent.HelpLink))
	}
	if policyEvent.SupportEmail != nil {
		cols = append(cols, handler.NewCol(PrivacyPolicySupportEmailCol, *policyEvent.SupportEmail))
	}
	if policyEvent.DocsLink != nil {
		cols = append(cols, handler.NewCol(PrivacyPolicyDocsLinkCol, *policyEvent.DocsLink))
	}
	if policyEvent.CustomLink != nil {
		cols = append(cols, handler.NewCol(PrivacyPolicyCustomLinkCol, *policyEvent.CustomLink))
	}
	if policyEvent.CustomLinkText != nil {
		cols = append(cols, handler.NewCol(PrivacyPolicyCustomLinkTextCol, *policyEvent.CustomLinkText))
	}
	return handler.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(PrivacyPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(PrivacyPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *privacyPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PrivacyPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-FvtGO", "reduce.wrong.event.type %s", org.PrivacyPolicyRemovedEventType)
	}
	return handler.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(PrivacyPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(PrivacyPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *privacyPolicyProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-bxJCY", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(PrivacyPolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(PrivacyPolicyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
