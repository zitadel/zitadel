package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	UserNotificationsProjectionTable = "projections.notifications"
)

type userNotifier struct {
	crdb.StatementHandler
	commands     *command.Commands
	queries      *NotificationQueries
	assetsPrefix func(context.Context) string
	otpEmailTmpl string
	metricSuccessfulDeliveriesEmail,
	metricFailedDeliveriesEmail,
	metricSuccessfulDeliveriesSMS,
	metricFailedDeliveriesSMS string
}

func NewUserNotifier(
	ctx context.Context,
	config crdb.StatementHandlerConfig,
	commands *command.Commands,
	queries *NotificationQueries,
	assetsPrefix func(context.Context) string,
	otpEmailTmpl string,
	metricSuccessfulDeliveriesEmail,
	metricFailedDeliveriesEmail,
	metricSuccessfulDeliveriesSMS,
	metricFailedDeliveriesSMS string,
) *userNotifier {
	p := new(userNotifier)
	config.ProjectionName = UserNotificationsProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	p.commands = commands
	p.queries = queries
	p.assetsPrefix = assetsPrefix
	p.otpEmailTmpl = otpEmailTmpl
	p.metricSuccessfulDeliveriesEmail = metricSuccessfulDeliveriesEmail
	p.metricFailedDeliveriesEmail = metricFailedDeliveriesEmail
	p.metricSuccessfulDeliveriesSMS = metricSuccessfulDeliveriesSMS
	p.metricFailedDeliveriesSMS = metricFailedDeliveriesSMS
	projection.NotificationsProjection = p
	return p
}

func (u *userNotifier) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserV1InitialCodeAddedType,
					Reduce: u.reduceInitCodeAdded,
				},
				{
					Event:  user.HumanInitialCodeAddedType,
					Reduce: u.reduceInitCodeAdded,
				},
				{
					Event:  user.UserV1EmailCodeAddedType,
					Reduce: u.reduceEmailCodeAdded,
				},
				{
					Event:  user.HumanEmailCodeAddedType,
					Reduce: u.reduceEmailCodeAdded,
				},
				{
					Event:  user.UserV1PasswordCodeAddedType,
					Reduce: u.reducePasswordCodeAdded,
				},
				{
					Event:  user.HumanPasswordCodeAddedType,
					Reduce: u.reducePasswordCodeAdded,
				},
				{
					Event:  user.UserDomainClaimedType,
					Reduce: u.reduceDomainClaimed,
				},
				{
					Event:  user.HumanPasswordlessInitCodeRequestedType,
					Reduce: u.reducePasswordlessCodeRequested,
				},
				{
					Event:  user.UserV1PhoneCodeAddedType,
					Reduce: u.reducePhoneCodeAdded,
				},
				{
					Event:  user.HumanPhoneCodeAddedType,
					Reduce: u.reducePhoneCodeAdded,
				},
				{
					Event:  user.HumanPasswordChangedType,
					Reduce: u.reducePasswordChanged,
				},
				{
					Event:  user.HumanOTPSMSCodeAddedType,
					Reduce: u.reduceOTPSMSCodeAdded,
				},
				{
					Event:  user.HumanOTPEmailCodeAddedType,
					Reduce: u.reduceOTPEmailCodeAdded,
				},
			},
		},
		{
			Aggregate: session.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  session.OTPSMSChallengedType,
					Reduce: u.reduceSessionOTPSMSChallenged,
				},
				{
					Event:  session.OTPEmailChallengedType,
					Reduce: u.reduceSessionOTPEmailChallenged,
				},
			},
		},
	}
}

func (u *userNotifier) reduceInitCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanInitialCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-EFe2f", "reduce.wrong.event.type %s", user.HumanInitialCodeAddedType)
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.UserV1InitialCodeAddedType, user.UserV1InitialCodeSentType,
		user.HumanInitialCodeAddedType, user.HumanInitialCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.InitCodeMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := u.queries.Origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		u.queries.GetSMTPConfig,
		u.queries.GetFileSystemProvider,
		u.queries.GetLogProvider,
		colors,
		u.assetsPrefix(ctx),
		e,
		u.metricSuccessfulDeliveriesEmail,
		u.metricFailedDeliveriesEmail,
	).SendUserInitCode(notifyUser, origin, code)
	if err != nil {
		return nil, err
	}
	err = u.commands.HumanInitCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (u *userNotifier) reduceEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SWf3g", "reduce.wrong.event.type %s", user.HumanEmailCodeAddedType)
	}

	if e.CodeReturned {
		return crdb.NewNoOpStatement(e), nil
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.UserV1EmailCodeAddedType, user.UserV1EmailCodeSentType,
		user.HumanEmailCodeAddedType, user.HumanEmailCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifyEmailMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := u.queries.Origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		u.queries.GetSMTPConfig,
		u.queries.GetFileSystemProvider,
		u.queries.GetLogProvider,
		colors,
		u.assetsPrefix(ctx),
		e,
		u.metricSuccessfulDeliveriesEmail,
		u.metricFailedDeliveriesEmail,
	).SendEmailVerificationCode(notifyUser, origin, code, e.URLTemplate)
	if err != nil {
		return nil, err
	}
	err = u.commands.HumanEmailVerificationCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (u *userNotifier) reducePasswordCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Eeg3s", "reduce.wrong.event.type %s", user.HumanPasswordCodeAddedType)
	}
	if e.CodeReturned {
		return crdb.NewNoOpStatement(e), nil
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.UserV1PasswordCodeAddedType, user.UserV1PasswordCodeSentType,
		user.HumanPasswordCodeAddedType, user.HumanPasswordCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordResetMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := u.queries.Origin(ctx)
	if err != nil {
		return nil, err
	}
	notify := types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		u.queries.GetSMTPConfig,
		u.queries.GetFileSystemProvider,
		u.queries.GetLogProvider,
		colors,
		u.assetsPrefix(ctx),
		e,
		u.metricSuccessfulDeliveriesEmail,
		u.metricFailedDeliveriesEmail,
	)
	if e.NotificationType == domain.NotificationTypeSms {
		notify = types.SendSMSTwilio(
			ctx,
			translator,
			notifyUser,
			u.queries.GetTwilioConfig,
			u.queries.GetFileSystemProvider,
			u.queries.GetLogProvider,
			colors,
			u.assetsPrefix(ctx),
			e,
			u.metricSuccessfulDeliveriesSMS,
			u.metricFailedDeliveriesSMS,
		)
	}
	err = notify.SendPasswordCode(notifyUser, origin, code, e.URLTemplate)
	if err != nil {
		return nil, err
	}
	err = u.commands.PasswordCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (u *userNotifier) reduceOTPSMSCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ASF3g", "reduce.wrong.event.type %s", user.HumanOTPSMSCodeAddedType)
	}
	return u.reduceOTPSMS(
		e,
		e.Code,
		e.Expiry,
		e.Aggregate().ID,
		e.Aggregate().ResourceOwner,
		u.commands.HumanOTPSMSCodeSent,
		user.HumanOTPSMSCodeAddedType,
		user.HumanOTPSMSCodeSentType,
	)
}

func (u *userNotifier) reduceSessionOTPSMSChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.OTPSMSChallengedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Sk32L", "reduce.wrong.event.type %s", session.OTPSMSChallengedType)
	}
	if e.CodeReturned {
		return crdb.NewNoOpStatement(e), nil
	}
	ctx := HandlerContext(event.Aggregate())
	s, err := u.queries.SessionByID(ctx, true, e.Aggregate().ID, "")
	if err != nil {
		return nil, err
	}
	return u.reduceOTPSMS(
		e,
		e.Code,
		e.Expiry,
		s.UserFactor.UserID,
		s.UserFactor.ResourceOwner,
		u.commands.OTPSMSSent,
		session.OTPSMSChallengedType,
		session.OTPSMSSentType,
	)
}

func (u *userNotifier) reduceOTPSMS(
	event eventstore.Event,
	code *crypto.CryptoValue,
	expiry time.Duration,
	userID,
	resourceOwner string,
	sentCommand func(ctx context.Context, userID string, resourceOwner string) (err error),
	eventTypes ...eventstore.EventType,
) (*handler.Statement, error) {
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, expiry, nil, eventTypes...)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(event), nil
	}
	plainCode, err := crypto.DecryptString(code, u.queries.UserDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, resourceOwner, false)
	if err != nil {
		return nil, err
	}

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, userID, false)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifySMSOTPMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := u.queries.Origin(ctx)
	if err != nil {
		return nil, err
	}
	notify := types.SendSMSTwilio(
		ctx,
		translator,
		notifyUser,
		u.queries.GetTwilioConfig,
		u.queries.GetFileSystemProvider,
		u.queries.GetLogProvider,
		colors,
		u.assetsPrefix(ctx),
		event,
		u.metricSuccessfulDeliveriesSMS,
		u.metricFailedDeliveriesSMS,
	)
	err = notify.SendOTPSMSCode(authz.GetInstance(ctx).RequestedDomain(), origin, plainCode, expiry)
	if err != nil {
		return nil, err
	}
	err = sentCommand(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(event), nil
}

func (u *userNotifier) reduceOTPEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-JL3hw", "reduce.wrong.event.type %s", user.HumanOTPEmailCodeAddedType)
	}
	var authRequestID string
	if e.AuthRequestInfo != nil {
		authRequestID = e.AuthRequestInfo.ID
	}
	url := func(code, origin string, _ *query.NotifyUser) (string, error) {
		return login.OTPLink(origin, authRequestID, code, domain.MFATypeOTPEmail), nil
	}
	return u.reduceOTPEmail(
		e,
		e.Code,
		e.Expiry,
		e.Aggregate().ID,
		e.Aggregate().ResourceOwner,
		url,
		u.commands.HumanOTPEmailCodeSent,
		user.HumanOTPEmailCodeAddedType,
		user.HumanOTPEmailCodeSentType,
	)
}

func (u *userNotifier) reduceSessionOTPEmailChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.OTPEmailChallengedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-zbsgt", "reduce.wrong.event.type %s", session.OTPEmailChallengedType)
	}
	if e.ReturnCode {
		return crdb.NewNoOpStatement(e), nil
	}
	ctx := HandlerContext(event.Aggregate())
	s, err := u.queries.SessionByID(ctx, true, e.Aggregate().ID, "")
	if err != nil {
		return nil, err
	}
	url := func(code, origin string, user *query.NotifyUser) (string, error) {
		var buf strings.Builder
		urlTmpl := origin + u.otpEmailTmpl
		if e.URLTmpl != "" {
			urlTmpl = e.URLTmpl
		}
		if err := domain.RenderOTPEmailURLTemplate(&buf, urlTmpl, code, user.ID, user.PreferredLoginName, user.DisplayName, user.PreferredLanguage); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return u.reduceOTPEmail(
		e,
		e.Code,
		e.Expiry,
		s.UserFactor.UserID,
		s.UserFactor.ResourceOwner,
		url,
		u.commands.OTPEmailSent,
		user.HumanOTPEmailCodeAddedType,
		user.HumanOTPEmailCodeSentType,
	)
}

func (u *userNotifier) reduceOTPEmail(
	event eventstore.Event,
	code *crypto.CryptoValue,
	expiry time.Duration,
	userID,
	resourceOwner string,
	urlTmpl func(code, origin string, user *query.NotifyUser) (string, error),
	sentCommand func(ctx context.Context, userID string, resourceOwner string) (err error),
	eventTypes ...eventstore.EventType,
) (*handler.Statement, error) {
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, expiry, nil, eventTypes...)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(event), nil
	}
	plainCode, err := crypto.DecryptString(code, u.queries.UserDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, resourceOwner, false)
	if err != nil {
		return nil, err
	}

	template, err := u.queries.MailTemplateByOrg(ctx, resourceOwner, false)
	if err != nil {
		return nil, err
	}

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, userID, false)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, resourceOwner, domain.VerifyEmailOTPMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := u.queries.Origin(ctx)
	if err != nil {
		return nil, err
	}
	url, err := urlTmpl(plainCode, origin, notifyUser)
	if err != nil {
		return nil, err
	}
	notify := types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		u.queries.GetSMTPConfig,
		u.queries.GetFileSystemProvider,
		u.queries.GetLogProvider,
		colors,
		u.assetsPrefix(ctx),
		event,
		u.metricSuccessfulDeliveriesEmail,
		u.metricFailedDeliveriesEmail,
	)
	err = notify.SendOTPEmailCode(notifyUser, url, authz.GetInstance(ctx).RequestedDomain(), origin, plainCode, expiry)
	if err != nil {
		return nil, err
	}
	err = sentCommand(ctx, event.Aggregate().ID, event.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(event), nil
}

func (u *userNotifier) reduceDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Drh5w", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, nil, user.AggregateType,
		user.UserDomainClaimedType, user.UserDomainClaimedSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.DomainClaimedMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := u.queries.Origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		u.queries.GetSMTPConfig,
		u.queries.GetFileSystemProvider,
		u.queries.GetLogProvider,
		colors,
		u.assetsPrefix(ctx),
		e,
		u.metricSuccessfulDeliveriesEmail,
		u.metricFailedDeliveriesEmail,
	).SendDomainClaimed(notifyUser, origin, e.UserName)
	if err != nil {
		return nil, err
	}
	err = u.commands.UserDomainClaimedSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (u *userNotifier) reducePasswordlessCodeRequested(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordlessInitCodeRequestedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-EDtjd", "reduce.wrong.event.type %s", user.HumanPasswordlessInitCodeAddedType)
	}
	if e.CodeReturned {
		return crdb.NewNoOpStatement(e), nil
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, map[string]interface{}{"id": e.ID}, user.HumanPasswordlessInitCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordlessRegistrationMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := u.queries.Origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendEmail(
		ctx,
		string(template.Template),
		translator,
		notifyUser,
		u.queries.GetSMTPConfig,
		u.queries.GetFileSystemProvider,
		u.queries.GetLogProvider,
		colors,
		u.assetsPrefix(ctx),
		e,
		u.metricSuccessfulDeliveriesEmail,
		u.metricFailedDeliveriesEmail,
	).SendPasswordlessRegistrationLink(notifyUser, origin, code, e.ID, e.URLTemplate)
	if err != nil {
		return nil, err
	}
	err = u.commands.HumanPasswordlessInitCodeSent(ctx, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (u *userNotifier) reducePasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Yko2z8", "reduce.wrong.event.type %s", user.HumanPasswordChangedType)
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, nil, user.AggregateType, user.HumanPasswordChangeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}

	notificationPolicy, err := u.queries.NotificationPolicyByOrg(ctx, true, e.Aggregate().ResourceOwner, false)
	if errors.IsNotFound(err) {
		return crdb.NewNoOpStatement(e), nil
	}
	if err != nil {
		return nil, err
	}

	if notificationPolicy.PasswordChange {
		colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return nil, err
		}

		template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return nil, err
		}

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
		if err != nil {
			return nil, err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordChangeMessageType)
		if err != nil {
			return nil, err
		}

		ctx, origin, err := u.queries.Origin(ctx)
		if err != nil {
			return nil, err
		}
		err = types.SendEmail(
			ctx,
			string(template.Template),
			translator,
			notifyUser,
			u.queries.GetSMTPConfig,
			u.queries.GetFileSystemProvider,
			u.queries.GetLogProvider,
			colors,
			u.assetsPrefix(ctx),
			e,
			u.metricSuccessfulDeliveriesEmail,
			u.metricFailedDeliveriesEmail,
		).SendPasswordChange(notifyUser, origin)
		if err != nil {
			return nil, err
		}
		err = u.commands.PasswordChangeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
		if err != nil {
			return nil, err
		}
	}
	return crdb.NewNoOpStatement(e), nil
}

func (u *userNotifier) reducePhoneCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-He83g", "reduce.wrong.event.type %s", user.HumanPhoneCodeAddedType)
	}
	if e.CodeReturned {
		return crdb.NewNoOpStatement(e), nil
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.UserV1PhoneCodeAddedType, user.UserV1PhoneCodeSentType,
		user.HumanPhoneCodeAddedType, user.HumanPhoneCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return crdb.NewNoOpStatement(e), nil
	}
	code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
	if err != nil {
		return nil, err
	}

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifyPhoneMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := u.queries.Origin(ctx)
	if err != nil {
		return nil, err
	}
	err = types.SendSMSTwilio(
		ctx,
		translator,
		notifyUser,
		u.queries.GetTwilioConfig,
		u.queries.GetFileSystemProvider,
		u.queries.GetLogProvider,
		colors,
		u.assetsPrefix(ctx),
		e,
		u.metricSuccessfulDeliveriesSMS,
		u.metricFailedDeliveriesSMS,
	).SendPhoneVerificationCode(notifyUser, origin, code, authz.GetInstance(ctx).RequestedDomain())
	if err != nil {
		return nil, err
	}
	err = u.commands.HumanPhoneVerificationCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	if err != nil {
		return nil, err
	}
	return crdb.NewNoOpStatement(e), nil
}

func (u *userNotifier) checkIfCodeAlreadyHandledOrExpired(ctx context.Context, event eventstore.Event, expiry time.Duration, data map[string]interface{}, eventTypes ...eventstore.EventType) (bool, error) {
	if event.CreationDate().Add(expiry).Before(time.Now().UTC()) {
		return true, nil
	}
	return u.queries.IsAlreadyHandled(ctx, event, data, user.AggregateType, eventTypes...)
}
