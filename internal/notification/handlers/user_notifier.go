package handlers

import (
	"context"
	"strings"
	"time"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UserNotificationsProjectionTable = "projections.notifications"
)

type userNotifier struct {
	commands     Commands
	queries      *NotificationQueries
	channels     types.ChannelChains
	otpEmailTmpl string
}

func NewUserNotifier(
	ctx context.Context,
	config handler.Config,
	commands Commands,
	queries *NotificationQueries,
	channels types.ChannelChains,
	otpEmailTmpl string,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &userNotifier{
		commands:     commands,
		queries:      queries,
		otpEmailTmpl: otpEmailTmpl,
		channels:     channels,
	})
}

func (u *userNotifier) Name() string {
	return UserNotificationsProjectionTable
}

func (u *userNotifier) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-EFe2f", "reduce.wrong.event.type %s", user.HumanInitialCodeAddedType)
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.UserV1InitialCodeAddedType, user.UserV1InitialCodeSentType,
			user.HumanInitialCodeAddedType, user.HumanInitialCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
		if err != nil {
			return err
		}
		colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.InitCodeMessageType)
		if err != nil {
			return err
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		err = types.SendEmail(ctx, u.channels, string(template.Template), translator, notifyUser, colors, e).
			SendUserInitCode(ctx, notifyUser, code, e.AuthRequestID)
		if err != nil {
			return err
		}
		return u.commands.HumanInitCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reduceEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SWf3g", "reduce.wrong.event.type %s", user.HumanEmailCodeAddedType)
	}

	if e.CodeReturned {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.UserV1EmailCodeAddedType, user.UserV1EmailCodeSentType,
			user.HumanEmailCodeAddedType, user.HumanEmailCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
		if err != nil {
			return err
		}
		colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifyEmailMessageType)
		if err != nil {
			return err
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		err = types.SendEmail(ctx, u.channels, string(template.Template), translator, notifyUser, colors, e).
			SendEmailVerificationCode(ctx, notifyUser, code, e.URLTemplate, e.AuthRequestID)
		if err != nil {
			return err
		}
		return u.commands.HumanEmailVerificationCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reducePasswordCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Eeg3s", "reduce.wrong.event.type %s", user.HumanPasswordCodeAddedType)
	}
	if e.CodeReturned {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.UserV1PasswordCodeAddedType, user.UserV1PasswordCodeSentType,
			user.HumanPasswordCodeAddedType, user.HumanPasswordCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
		if err != nil {
			return err
		}
		colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordResetMessageType)
		if err != nil {
			return err
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		notify := types.SendEmail(ctx, u.channels, string(template.Template), translator, notifyUser, colors, e)

		if e.NotificationType == domain.NotificationTypeSms {
			twilioVerificationEnabled, err := u.isTwilioVerificationAPIEnabled(ctx)
			if err != nil {
				return err
			}

			if twilioVerificationEnabled {
				notify = types.SendSMSTwilioVerifyRequest(ctx, u.channels, notifyUser, e)
			} else {
				notify = types.SendSMSTwilio(ctx, u.channels, translator, notifyUser, colors, e)
			}
		}

		err = notify.SendPasswordCode(ctx, notifyUser, code, e.URLTemplate, e.AuthRequestID)
		if err != nil {
			return err
		}
		return u.commands.PasswordCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reduceOTPSMSCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASF3g", "reduce.wrong.event.type %s", user.HumanOTPSMSCodeAddedType)
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Sk32L", "reduce.wrong.event.type %s", session.OTPSMSChallengedType)
	}
	if e.CodeReturned {
		return handler.NewNoOpStatement(e), nil
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
		return handler.NewNoOpStatement(event), nil
	}
	plainCode, err := crypto.DecryptString(code, u.queries.UserDataCrypto)
	if err != nil {
		return nil, err
	}
	colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, resourceOwner, false)
	if err != nil {
		return nil, err
	}

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, userID)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifySMSOTPMessageType)
	if err != nil {
		return nil, err
	}
	ctx, err = u.queries.Origin(ctx, event)
	if err != nil {
		return nil, err
	}

	twilioVerificationEnabled, err := u.isTwilioVerificationAPIEnabled(ctx)
	if err != nil {
		return nil, err
	}

	var notify types.Notify
	if twilioVerificationEnabled {
		notify = types.SendSMSTwilioVerifyRequest(ctx, u.channels, notifyUser, event)
	} else {
		notify = types.SendSMSTwilio(ctx, u.channels, translator, notifyUser, colors, event)
	}

	err = notify.SendOTPSMSCode(ctx, plainCode, expiry)
	if err != nil {
		return nil, err
	}

	err = sentCommand(ctx, event.Aggregate().ID, event.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}
	return handler.NewNoOpStatement(event), nil
}

func (u *userNotifier) reduceOTPEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JL3hw", "reduce.wrong.event.type %s", user.HumanOTPEmailCodeAddedType)
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-zbsgt", "reduce.wrong.event.type %s", session.OTPEmailChallengedType)
	}
	if e.ReturnCode {
		return handler.NewNoOpStatement(e), nil
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
		return handler.NewNoOpStatement(event), nil
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

	notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, userID)
	if err != nil {
		return nil, err
	}
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, resourceOwner, domain.VerifyEmailOTPMessageType)
	if err != nil {
		return nil, err
	}
	ctx, err = u.queries.Origin(ctx, event)
	if err != nil {
		return nil, err
	}
	url, err := urlTmpl(plainCode, http_util.ComposedOrigin(ctx), notifyUser)
	if err != nil {
		return nil, err
	}
	notify := types.SendEmail(ctx, u.channels, string(template.Template), translator, notifyUser, colors, event)
	err = notify.SendOTPEmailCode(ctx, url, plainCode, expiry)
	if err != nil {
		return nil, err
	}
	err = sentCommand(ctx, event.Aggregate().ID, event.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}
	return handler.NewNoOpStatement(event), nil
}

func (u *userNotifier) reduceDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Drh5w", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
	}
	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, nil,
			user.UserDomainClaimedType, user.UserDomainClaimedSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.DomainClaimedMessageType)
		if err != nil {
			return err
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		err = types.SendEmail(ctx, u.channels, string(template.Template), translator, notifyUser, colors, e).
			SendDomainClaimed(ctx, notifyUser, e.UserName)
		if err != nil {
			return err
		}
		return u.commands.UserDomainClaimedSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reducePasswordlessCodeRequested(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordlessInitCodeRequestedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-EDtjd", "reduce.wrong.event.type %s", user.HumanPasswordlessInitCodeAddedType)
	}
	if e.CodeReturned {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, map[string]interface{}{"id": e.ID}, user.HumanPasswordlessInitCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
		if err != nil {
			return err
		}
		colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordlessRegistrationMessageType)
		if err != nil {
			return err
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		err = types.SendEmail(ctx, u.channels, string(template.Template), translator, notifyUser, colors, e).
			SendPasswordlessRegistrationLink(ctx, notifyUser, code, e.ID, e.URLTemplate)
		if err != nil {
			return err
		}
		return u.commands.HumanPasswordlessInitCodeSent(ctx, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.ID)
	}), nil
}

func (u *userNotifier) reducePasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yko2z8", "reduce.wrong.event.type %s", user.HumanPasswordChangedType)
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, nil, user.HumanPasswordChangeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}

		notificationPolicy, err := u.queries.NotificationPolicyByOrg(ctx, true, e.Aggregate().ResourceOwner, false)
		if zerrors.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}

		if !notificationPolicy.PasswordChange {
			return nil
		}

		colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		template, err := u.queries.MailTemplateByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordChangeMessageType)
		if err != nil {
			return err
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		err = types.SendEmail(ctx, u.channels, string(template.Template), translator, notifyUser, colors, e).
			SendPasswordChange(ctx, notifyUser)
		if err != nil {
			return err
		}
		return u.commands.PasswordChangeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reducePhoneCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-He83g", "reduce.wrong.event.type %s", user.HumanPhoneCodeAddedType)
	}
	if e.CodeReturned {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.UserV1PhoneCodeAddedType, user.UserV1PhoneCodeSentType,
			user.HumanPhoneCodeAddedType, user.HumanPhoneCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		code, err := crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
		if err != nil {
			return err
		}
		colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.Aggregate().ResourceOwner, false)
		if err != nil {
			return err
		}

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifyPhoneMessageType)
		if err != nil {
			return err
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}

		twilioVerificationEnabled, err := u.isTwilioVerificationAPIEnabled(ctx)
		if err != nil {
			return err
		}

		var notify types.Notify
		if twilioVerificationEnabled {
			notify = types.SendSMSTwilioVerifyRequest(ctx, u.channels, notifyUser, e)
		} else {
			notify = types.SendSMSTwilio(ctx, u.channels, translator, notifyUser, colors, e)
		}

		err = notify.SendPhoneVerificationCode(ctx, code)
		if err != nil {
			return err
		}

		return u.commands.HumanPhoneVerificationCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) checkIfCodeAlreadyHandledOrExpired(ctx context.Context, event eventstore.Event, expiry time.Duration, data map[string]interface{}, eventTypes ...eventstore.EventType) (bool, error) {
	if event.CreatedAt().Add(expiry).Before(time.Now().UTC()) {
		return true, nil
	}
	return u.queries.IsAlreadyHandled(ctx, event, data, eventTypes...)
}

func (u *userNotifier) isTwilioVerificationAPIEnabled(ctx context.Context) (bool, error) {
	_, twilioConfig, err := u.channels.SMS(ctx)
	if err != nil {
		return false, err
	}
	return twilioConfig.VerifyServiceSID != "", nil
}
