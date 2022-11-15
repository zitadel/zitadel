package quotas

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/repository/quota"

	"github.com/zitadel/zitadel/internal/repository/instance"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	NotificationsProjectionTable = "projections.quotas"
)

func Start(
	ctx context.Context,
	customConfig projection.CustomConfig,
	commands *command.Commands,
	queries *query.Queries,
	es *eventstore.Eventstore,
) {
	projection.QuotasProjection = newQuotasProjection(ctx, projection.ApplyCustomConfig(customConfig), commands, queries, es)
}

type notificationsProjection struct {
	crdb.StatementHandler
	commands *command.Commands
	queries  *query.Queries
	es       *eventstore.Eventstore
}

func newQuotasProjection(
	ctx context.Context,
	config crdb.StatementHandlerConfig,
	commands *command.Commands,
	queries *query.Queries,
	es *eventstore.Eventstore,
) *notificationsProjection {
	p := new(notificationsProjection)
	config.ProjectionName = NotificationsProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	p.commands = commands
	p.queries = queries
	p.es = es
	return p
}

func (p *notificationsProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  quota.AddedEventType,
					Reduce: p.reduceQuotaAdded,
				},
				{
					Event:  quota.RemovedEventType,
					Reduce: p.reduceInitCodeAdded,
				},
			},
		},
	}
}

func (p *notificationsProjection) reduceQuotaAdded(event eventstore.Event) (*handler.Statement, error) {
	e := event.(*quota.AddedEvent)

	ctx := setQuotaContext(event.Aggregate())
	alreadyHandled, err := p.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.UserV1InitialCodeAddedType, user.UserV1InitialCodeSentType,
		user.HumanInitialCodeAddedType, user.HumanInitialCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, p.userDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := p.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	template, err := p.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	notifyUser, err := p.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	translator, err := p.getTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.InitCodeMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := p.origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		p.getSMTPConfig,
		p.getFileSystemProvider,
		p.getLogProvider,
		colors,
		p.assetsPrefix(ctx),
	).SendUserInitCode(notifyUser, origin, code)
	if err != nil {
		return nil, err
	}
	err = p.commands.HumanInitCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (p *notificationsProjection) reduceEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SWf3g", "reduce.wrong.event.type %s", user.HumanEmailCodeAddedType)
	}
	ctx := setQuotaContext(event.Aggregate())
	alreadyHandled, err := p.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.UserV1EmailCodeAddedType, user.UserV1EmailCodeSentType,
		user.HumanEmailCodeAddedType, user.HumanEmailCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, p.userDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := p.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	template, err := p.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	notifyUser, err := p.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	translator, err := p.getTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifyEmailMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := p.origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		p.getSMTPConfig,
		p.getFileSystemProvider,
		p.getLogProvider,
		colors,
		p.assetsPrefix(ctx),
	).SendEmailVerificationCode(notifyUser, origin, code)
	if err != nil {
		return nil, err
	}
	err = p.commands.HumanEmailVerificationCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (p *notificationsProjection) reducePasswordCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Eeg3s", "reduce.wrong.event.type %s", user.HumanPasswordCodeAddedType)
	}
	ctx := setQuotaContext(event.Aggregate())
	alreadyHandled, err := p.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.UserV1PasswordCodeAddedType, user.UserV1PasswordCodeSentType,
		user.HumanPasswordCodeAddedType, user.HumanPasswordCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, p.userDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := p.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	template, err := p.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	notifyUser, err := p.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	translator, err := p.getTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordResetMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := p.origin(ctx)
	if err != nil {
		return nil, err
	}
	notify := types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		p.getSMTPConfig,
		p.getFileSystemProvider,
		p.getLogProvider,
		colors,
		p.assetsPrefix(ctx),
	)
	if e.NotificationType == domain.NotificationTypeSms {
		notify = types.SendSMSTwilio(
			ctx,
			translator,
			notifyUser,
			p.getTwilioConfig,
			p.getFileSystemProvider,
			p.getLogProvider,
			colors,
			p.assetsPrefix(ctx),
		)
	}
	err = notify.SendPasswordCode(notifyUser, origin, code)
	if err != nil {
		return nil, err
	}
	err = p.commands.PasswordCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (p *notificationsProjection) reduceDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Drh5w", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
	}
	ctx := setQuotaContext(event.Aggregate())
	alreadyHandled, err := p.checkIfAlreadyHandled(ctx, event, nil,
		user.UserDomainClaimedType, user.UserDomainClaimedSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	colors, err := p.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	template, err := p.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	notifyUser, err := p.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	translator, err := p.getTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.DomainClaimedMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := p.origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		p.getSMTPConfig,
		p.getFileSystemProvider,
		p.getLogProvider,
		colors,
		p.assetsPrefix(ctx),
	).SendDomainClaimed(notifyUser, origin, e.UserName)
	if err != nil {
		return nil, err
	}
	err = p.commands.UserDomainClaimedSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (p *notificationsProjection) reducePasswordlessCodeRequested(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordlessInitCodeRequestedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-EDtjd", "reduce.wrong.event.type %s", user.HumanPasswordlessInitCodeAddedType)
	}
	ctx := setQuotaContext(event.Aggregate())
	alreadyHandled, err := p.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, map[string]interface{}{"id": e.ID}, user.HumanPasswordlessInitCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, p.userDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := p.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	template, err := p.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	notifyUser, err := p.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	translator, err := p.getTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordlessRegistrationMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := p.origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		p.getSMTPConfig,
		p.getFileSystemProvider,
		p.getLogProvider,
		colors,
		p.assetsPrefix(ctx),
	).SendPasswordlessRegistrationLink(notifyUser, origin, code, e.ID)
	if err != nil {
		return nil, err
	}
	err = p.commands.HumanPasswordlessInitCodeSent(ctx, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (p *notificationsProjection) reducePhoneCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-He83g", "reduce.wrong.event.type %s", user.HumanPhoneCodeAddedType)
	}
	ctx := setQuotaContext(event.Aggregate())
	alreadyHandled, err := p.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.UserV1PhoneCodeAddedType, user.UserV1PhoneCodeSentType,
		user.HumanPhoneCodeAddedType, user.HumanPhoneCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, p.userDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := p.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}

	notifyUser, err := p.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	translator, err := p.getTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifyPhoneMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := p.origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendSMSTwilio(
		ctx,
		translator,
		notifyUser,
		p.getTwilioConfig,
		p.getFileSystemProvider,
		p.getLogProvider,
		colors,
		p.assetsPrefix(ctx),
	).SendPhoneVerificationCode(notifyUser, origin, code)
	if err != nil {
		return nil, err
	}
	err = p.commands.HumanPhoneVerificationCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (p *notificationsProjection) checkIfCodeAlreadyHandledOrExpired(ctx context.Context, event eventstore.Event, expiry time.Duration, data map[string]interface{}, eventTypes ...eventstore.EventType) (bool, error) {
	if event.CreationDate().Add(expiry).Before(time.Now().UTC()) {
		return true, nil
	}
	return p.checkIfAlreadyHandled(ctx, event, data, eventTypes...)
}

func (p *notificationsProjection) checkIfAlreadyHandled(ctx context.Context, event eventstore.Event, data map[string]interface{}, eventTypes ...eventstore.EventType) (bool, error) {
	events, err := p.es.Filter(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(event.Aggregate().InstanceID).
			AddQuery().
			AggregateTypes(user.AggregateType).
			AggregateIDs(event.Aggregate().ID).
			SequenceGreater(event.Sequence()).
			EventTypes(eventTypes...).
			EventData(data).
			Builder(),
	)
	if err != nil {
		return false, err
	}
	return len(events) > 0, nil
}
func (p *notificationsProjection) getSMTPConfig(ctx context.Context) (*smtp.EmailConfig, error) {
	config, err := p.queries.SMTPConfigByAggregateID(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	password, err := crypto.DecryptString(config.Password, p.smtpPasswordCrypto)
	if err != nil {
		return nil, err
	}
	return &smtp.EmailConfig{
		From:     config.SenderAddress,
		FromName: config.SenderName,
		Tls:      config.TLS,
		SMTP: smtp.SMTP{
			Host:     config.Host,
			User:     config.User,
			Password: password,
		},
	}, nil
}

// Read iam twilio config
func (p *notificationsProjection) getTwilioConfig(ctx context.Context) (*twilio.TwilioConfig, error) {
	active, err := query.NewSMSProviderStateQuery(domain.SMSConfigStateActive)
	if err != nil {
		return nil, err
	}
	config, err := p.queries.SMSProviderConfig(ctx, active)
	if err != nil {
		return nil, err
	}
	if config.TwilioConfig == nil {
		return nil, errors.ThrowNotFound(nil, "HANDLER-8nfow", "Errors.SMS.Twilio.NotFound")
	}
	token, err := crypto.DecryptString(config.TwilioConfig.Token, p.smsTokenCrypto)
	if err != nil {
		return nil, err
	}
	return &twilio.TwilioConfig{
		SID:          config.TwilioConfig.SID,
		Token:        token,
		SenderNumber: config.TwilioConfig.SenderNumber,
	}, nil
}

// Read iam filesystem provider config
func (p *notificationsProjection) getFileSystemProvider(ctx context.Context) (*fs.FSConfig, error) {
	config, err := p.queries.NotificationProviderByIDAndType(ctx, authz.GetInstance(ctx).InstanceID(), domain.NotificationProviderTypeFile)
	if err != nil {
		return nil, err
	}
	return &fs.FSConfig{
		Compact: config.Compact,
		Path:    p.fileSystemPath,
	}, nil
}

// Read iam log provider config
func (p *notificationsProjection) getLogProvider(ctx context.Context) (*log.LogConfig, error) {
	config, err := p.queries.NotificationProviderByIDAndType(ctx, authz.GetInstance(ctx).InstanceID(), domain.NotificationProviderTypeLog)
	if err != nil {
		return nil, err
	}
	return &log.LogConfig{
		Compact: config.Compact,
	}, nil
}

func (p *notificationsProjection) getTranslatorWithOrgTexts(ctx context.Context, orgID, textType string) (*i18n.Translator, error) {
	translator, err := i18n.NewTranslator(p.statikDir, p.queries.GetDefaultLanguage(ctx), "")
	if err != nil {
		return nil, err
	}

	allCustomTexts, err := p.queries.CustomTextListByTemplate(ctx, authz.GetInstance(ctx).InstanceID(), textType)
	if err != nil {
		return translator, nil
	}
	customTexts, err := p.queries.CustomTextListByTemplate(ctx, orgID, textType)
	if err != nil {
		return translator, nil
	}
	allCustomTexts.CustomTexts = append(allCustomTexts.CustomTexts, customTexts.CustomTexts...)

	for _, text := range allCustomTexts.CustomTexts {
		msg := i18n.Message{
			ID:   text.Template + "." + text.Key,
			Text: text.Text,
		}
		err = translator.AddMessages(text.Language, msg)
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "orgID", orgID, "messageType", textType, "messageID", msg.ID).
			OnError(err).
			Warn("could not add translation message")
	}
	return translator, nil
}

func (p *notificationsProjection) origin(ctx context.Context) (context.Context, string, error) {
	primary, err := query.NewInstanceDomainPrimarySearchQuery(true)
	if err != nil {
		return ctx, "", err
	}
	domains, err := p.queries.SearchInstanceDomains(ctx, &query.InstanceDomainSearchQueries{
		Queries: []query.SearchQuery{primary},
	})
	if err != nil {
		return ctx, "", err
	}
	if len(domains.Domains) < 1 {
		return ctx, "", errors.ThrowInternal(nil, "NOTIF-Ef3r1", "Errors.Notification.NoDomain")
	}
	ctx = authz.WithRequestedDomain(ctx, domains.Domains[0].Domain)
	return ctx, http_utils.BuildHTTP(domains.Domains[0].Domain, p.externalPort, p.externalSecure), nil
}

func setQuotaContext(event eventstore.Aggregate) context.Context {
	return authz.WithInstanceID(context.Background(), event.InstanceID)
}
