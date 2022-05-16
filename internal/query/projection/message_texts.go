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

type messageTextProjection struct {
	crdb.StatementHandler
}

const (
	MessageTextTable = "zitadel.projections.message_texts"

	MessageTextAggregateIDCol  = "aggregate_id"
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
)

func newMessageTextProjection(ctx context.Context, config crdb.StatementHandlerConfig) *messageTextProjection {
	p := &messageTextProjection{}
	config.ProjectionName = MessageTextTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *messageTextProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
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
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.CustomTextSetEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.CustomTextRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  iam.CustomTextTemplateRemovedEventType,
					Reduce: p.reduceTemplateRemoved,
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
	case *iam.CustomTextSetEvent:
		templateEvent = e.CustomTextSetEvent
	default:
		logging.LogWithFields("PROJE-2N9fg", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.CustomTextSetEventType, iam.CustomTextSetEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-2n90r", "reduce.wrong.event.type")
	}
	if !isMessageTemplate(templateEvent.Template) {
		return crdb.NewNoOpStatement(event), nil
	}

	cols := []handler.Column{
		handler.NewCol(MessageTextAggregateIDCol, templateEvent.Aggregate().ID),
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
	return crdb.NewUpsertStatement(
		&templateEvent,
		cols), nil
}

func (p *messageTextProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	var templateEvent policy.CustomTextRemovedEvent
	switch e := event.(type) {
	case *org.CustomTextRemovedEvent:
		templateEvent = e.CustomTextRemovedEvent
	case *iam.CustomTextRemovedEvent:
		templateEvent = e.CustomTextRemovedEvent
	default:
		logging.LogWithFields("PROJE-3m022", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.CustomTextRemovedEventType, iam.CustomTextRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-fm0ge", "reduce.wrong.event.type")
	}
	if !isMessageTemplate(templateEvent.Template) {
		return crdb.NewNoOpStatement(event), nil
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
	return crdb.NewUpdateStatement(
		&templateEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(MessageTextAggregateIDCol, templateEvent.Aggregate().ID),
			handler.NewCond(MessageTextTypeCol, templateEvent.Template),
			handler.NewCond(MessageTextLanguageCol, templateEvent.Language.String()),
		},
	), nil
}

func (p *messageTextProjection) reduceTemplateRemoved(event eventstore.Event) (*handler.Statement, error) {
	var templateEvent policy.CustomTextTemplateRemovedEvent
	switch e := event.(type) {
	case *org.CustomTextTemplateRemovedEvent:
		templateEvent = e.CustomTextTemplateRemovedEvent
	case *iam.CustomTextTemplateRemovedEvent:
		templateEvent = e.CustomTextTemplateRemovedEvent
	default:
		logging.LogWithFields("PROJE-m03ng", "seq", event.Sequence(), "expectedType", org.CustomTextTemplateRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-2n9rs", "reduce.wrong.event.type")
	}
	if !isMessageTemplate(templateEvent.Template) {
		return crdb.NewNoOpStatement(event), nil
	}
	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(MessageTextAggregateIDCol, templateEvent.Aggregate().ID),
			handler.NewCond(MessageTextTypeCol, templateEvent.Template),
			handler.NewCond(MessageTextLanguageCol, templateEvent.Language.String()),
		},
	), nil
}

func isMessageTemplate(template string) bool {
	return template == domain.InitCodeMessageType ||
		template == domain.PasswordResetMessageType ||
		template == domain.VerifyEmailMessageType ||
		template == domain.VerifyPhoneMessageType ||
		template == domain.DomainClaimedMessageType ||
		template == domain.PasswordlessRegistrationMessageType
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
