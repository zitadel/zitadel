package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

const (
	MessageTextTable = "projections.message_texts2"

	MessageTextAggregateIDCol  = "aggregate_id"
	MessageTextInstanceIDCol   = "instance_id"
	MessageTextCreationDateCol = "creation_date"
	MessageTextChangeDateCol   = "change_date"
	MessageTextSequenceCol     = "sequence"
	MessageTextStateCol        = "state"
	MessageTextTypeCol         = "type"
	MessageTextLanguageCol     = "language"
	MessageTextTitleCol        = "title"
	MessageTextPreHeaderCol    = "pre_header"
	MessageTextSubjectCol      = "subject"
	MessageTextGreetingCol     = "greeting"
	MessageTextTextCol         = "text"
	MessageTextButtonTextCol   = "button_text"
	MessageTextFooterCol       = "footer_text"
	MessageTextOwnerRemovedCol = "owner_removed"
)

type messageTextProjection struct{}

func newMessageTextProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(messageTextProjection))
}

func (*messageTextProjection) Name() string {
	return MessageTextTable
}

func (*messageTextProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(MessageTextAggregateIDCol, handler.ColumnTypeText),
			handler.NewColumn(MessageTextInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(MessageTextCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(MessageTextChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(MessageTextSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(MessageTextStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(MessageTextTypeCol, handler.ColumnTypeText),
			handler.NewColumn(MessageTextLanguageCol, handler.ColumnTypeText),
			handler.NewColumn(MessageTextTitleCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MessageTextPreHeaderCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MessageTextSubjectCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MessageTextGreetingCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MessageTextTextCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MessageTextButtonTextCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MessageTextFooterCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MessageTextOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(MessageTextInstanceIDCol, MessageTextAggregateIDCol, MessageTextTypeCol, MessageTextLanguageCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{MessageTextOwnerRemovedCol})),
		),
	)
}

func (p *messageTextProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.CustomTextSetEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.CustomTextRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  org.CustomTextTemplateRemovedEventType,
					Reduce: p.reduceTemplateRemoved,
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
					Event:  instance.CustomTextSetEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.CustomTextRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  instance.CustomTextTemplateRemovedEventType,
					Reduce: p.reduceTemplateRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(MessageTextInstanceIDCol),
				},
			},
		},
	}
}

func (p *messageTextProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var templateEvent policy.CustomTextSetEvent
	switch e := event.(type) {
	case *org.CustomTextSetEvent:
		templateEvent = e.CustomTextSetEvent
	case *instance.CustomTextSetEvent:
		templateEvent = e.CustomTextSetEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-2n90r", "reduce.wrong.event.type %v", []eventstore.EventType{org.CustomTextSetEventType, instance.CustomTextSetEventType})
	}
	if !isMessageTemplate(templateEvent.Template) {
		return handler.NewNoOpStatement(event), nil
	}

	cols := []handler.Column{
		handler.NewCol(MessageTextAggregateIDCol, templateEvent.Aggregate().ID),
		handler.NewCol(MessageTextInstanceIDCol, templateEvent.Aggregate().InstanceID),
		handler.NewCol(MessageTextCreationDateCol, templateEvent.CreationDate()),
		handler.NewCol(MessageTextChangeDateCol, templateEvent.CreationDate()),
		handler.NewCol(MessageTextSequenceCol, templateEvent.Sequence()),
		handler.NewCol(MessageTextStateCol, domain.PolicyStateActive),
		handler.NewCol(MessageTextTypeCol, templateEvent.Template),
		handler.NewCol(MessageTextLanguageCol, templateEvent.Language.String()),
	}
	if isTitle(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextTitleCol, templateEvent.Text))
	}
	if isPreHeader(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextPreHeaderCol, templateEvent.Text))
	}
	if isSubject(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextSubjectCol, templateEvent.Text))
	}
	if isGreeting(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextGreetingCol, templateEvent.Text))
	}
	if isText(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextTextCol, templateEvent.Text))
	}
	if isButtonText(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextButtonTextCol, templateEvent.Text))
	}
	if isFooterText(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextFooterCol, templateEvent.Text))
	}
	return handler.NewUpsertStatement(
		&templateEvent,
		[]handler.Column{
			handler.NewCol(MessageTextInstanceIDCol, nil),
			handler.NewCol(MessageTextAggregateIDCol, nil),
			handler.NewCol(MessageTextTypeCol, nil),
			handler.NewCol(MessageTextLanguageCol, nil),
		},
		cols,
	), nil
}

func (p *messageTextProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	var templateEvent policy.CustomTextRemovedEvent
	switch e := event.(type) {
	case *org.CustomTextRemovedEvent:
		templateEvent = e.CustomTextRemovedEvent
	case *instance.CustomTextRemovedEvent:
		templateEvent = e.CustomTextRemovedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-fm0ge", "reduce.wrong.event.type %v", []eventstore.EventType{org.CustomTextRemovedEventType, instance.CustomTextRemovedEventType})
	}
	if !isMessageTemplate(templateEvent.Template) {
		return handler.NewNoOpStatement(event), nil
	}
	cols := []handler.Column{
		handler.NewCol(MessageTextChangeDateCol, templateEvent.CreationDate()),
		handler.NewCol(MessageTextSequenceCol, templateEvent.Sequence()),
	}
	if isTitle(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextTitleCol, ""))
	}
	if isPreHeader(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextPreHeaderCol, ""))
	}
	if isSubject(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextSubjectCol, ""))
	}
	if isGreeting(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextGreetingCol, ""))
	}
	if isText(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextTextCol, ""))
	}
	if isButtonText(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextButtonTextCol, ""))
	}
	if isFooterText(templateEvent.Key) {
		cols = append(cols, handler.NewCol(MessageTextFooterCol, ""))
	}
	return handler.NewUpdateStatement(
		&templateEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(MessageTextAggregateIDCol, templateEvent.Aggregate().ID),
			handler.NewCond(MessageTextTypeCol, templateEvent.Template),
			handler.NewCond(MessageTextLanguageCol, templateEvent.Language.String()),
			handler.NewCond(MessageTextInstanceIDCol, templateEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *messageTextProjection) reduceTemplateRemoved(event eventstore.Event) (*handler.Statement, error) {
	var templateEvent policy.CustomTextTemplateRemovedEvent
	switch e := event.(type) {
	case *org.CustomTextTemplateRemovedEvent:
		templateEvent = e.CustomTextTemplateRemovedEvent
	case *instance.CustomTextTemplateRemovedEvent:
		templateEvent = e.CustomTextTemplateRemovedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-2n9rs", "reduce.wrong.event.type %s", org.CustomTextTemplateRemovedEventType)
	}
	if !isMessageTemplate(templateEvent.Template) {
		return handler.NewNoOpStatement(event), nil
	}
	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(MessageTextAggregateIDCol, templateEvent.Aggregate().ID),
			handler.NewCond(MessageTextTypeCol, templateEvent.Template),
			handler.NewCond(MessageTextLanguageCol, templateEvent.Language.String()),
			handler.NewCond(MessageTextInstanceIDCol, templateEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *messageTextProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-mLsQw", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(MessageTextInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(MessageTextAggregateIDCol, e.Aggregate().ID),
		},
	), nil
}

func isMessageTemplate(template string) bool {
	return template == domain.InitCodeMessageType ||
		template == domain.PasswordResetMessageType ||
		template == domain.VerifyEmailMessageType ||
		template == domain.VerifyPhoneMessageType ||
		template == domain.VerifySMSOTPMessageType ||
		template == domain.VerifyEmailOTPMessageType ||
		template == domain.DomainClaimedMessageType ||
		template == domain.PasswordlessRegistrationMessageType ||
		template == domain.PasswordChangeMessageType
}
func isTitle(key string) bool {
	return key == domain.MessageTitle
}
func isPreHeader(key string) bool {
	return key == domain.MessagePreHeader
}
func isSubject(key string) bool {
	return key == domain.MessageSubject
}
func isGreeting(key string) bool {
	return key == domain.MessageGreeting
}
func isText(key string) bool {
	return key == domain.MessageText
}
func isButtonText(key string) bool {
	return key == domain.MessageButtonText
}
func isFooterText(key string) bool {
	return key == domain.MessageFooterText
}
