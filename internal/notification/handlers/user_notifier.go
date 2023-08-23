package handlers

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	UserNotificationsProjectionTable = "projections.notifications"
)

type userNotifier struct {
	commands     *command.Commands
	queries      *NotificationQueries
	assetsPrefix func(context.Context) string
	metricSuccessfulDeliveriesEmail,
	metricFailedDeliveriesEmail,
	metricSuccessfulDeliveriesSMS,
	metricFailedDeliveriesSMS string
}

func NewUserNotifier(
	ctx context.Context,
	config handler.Config,
	commands *command.Commands,
	queries *NotificationQueries,
	assetsPrefix func(context.Context) string,
	metricSuccessfulDeliveriesEmail,
	metricFailedDeliveriesEmail,
	metricSuccessfulDeliveriesSMS,
	metricFailedDeliveriesSMS string,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &userNotifier{
		commands:                        commands,
		queries:                         queries,
		assetsPrefix:                    assetsPrefix,
		metricSuccessfulDeliveriesEmail: metricSuccessfulDeliveriesEmail,
		metricFailedDeliveriesEmail:     metricFailedDeliveriesEmail,
		metricSuccessfulDeliveriesSMS:   metricSuccessfulDeliveriesSMS,
		metricFailedDeliveriesSMS:       metricFailedDeliveriesSMS,
	})
}

func (u *userNotifier) Name() string {
	return UserNotificationsProjectionTable
}

func (u *userNotifier) Reducers() []handler.AggregateReducer {
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
	}
}

func (u *userNotifier) reduceInitCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanInitialCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-EFe2f", "reduce.wrong.event.type %s", user.HumanInitialCodeAddedType)
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

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.InitCodeMessageType)
		if err != nil {
			return err
		}

		ctx, origin, err := u.queries.Origin(ctx)
		if err != nil {
			return err
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
			return err
		}
		return u.commands.HumanInitCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reduceEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SWf3g", "reduce.wrong.event.type %s", user.HumanEmailCodeAddedType)
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

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifyEmailMessageType)
		if err != nil {
			return err
		}

		ctx, origin, err := u.queries.Origin(ctx)
		if err != nil {
			return err
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
			return err
		}
		return u.commands.HumanEmailVerificationCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reducePasswordCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Eeg3s", "reduce.wrong.event.type %s", user.HumanPasswordCodeAddedType)
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

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordResetMessageType)
		if err != nil {
			return err
		}

		ctx, origin, err := u.queries.Origin(ctx)
		if err != nil {
			return err
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
			return err
		}
		return u.commands.PasswordCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reduceOTPSMSCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ASF3g", "reduce.wrong.event.type %s", user.HumanOTPSMSCodeAddedType)
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.HumanOTPSMSCodeAddedType, user.HumanOTPSMSCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return handler.NewNoOpStatement(e), nil
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
		e,
		u.metricSuccessfulDeliveriesSMS,
		u.metricFailedDeliveriesSMS,
	)
	err = notify.SendOTPSMSCode(authz.GetInstance(ctx).RequestedDomain(), origin, code, e.Expiry)
	if err != nil {
		return nil, err
	}
	err = u.commands.HumanOTPSMSCodeSent(ctx, e.Aggregate().ID, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}
	return handler.NewNoOpStatement(e), nil
}

func (u *userNotifier) reduceOTPEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-JL3hw", "reduce.wrong.event.type %s", user.HumanOTPEmailCodeAddedType)
	}
	ctx := HandlerContext(event.Aggregate())
	alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
		user.HumanOTPEmailCodeAddedType, user.HumanOTPEmailCodeSentType)
	if err != nil {
		return nil, err
	}
	if alreadyHandled {
		return handler.NewNoOpStatement(e), nil
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
	translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifyEmailOTPMessageType)
	if err != nil {
		return nil, err
	}

	ctx, origin, err := u.queries.Origin(ctx)
	if err != nil {
		return nil, err
	}
	var authRequestID string
	if e.AuthRequestInfo != nil {
		authRequestID = e.AuthRequestInfo.ID
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
	err = notify.SendOTPEmailCode(notifyUser, authz.GetInstance(ctx).RequestedDomain(), origin, code, authRequestID, e.Expiry)
	if err != nil {
		return nil, err
	}
	err = u.commands.HumanOTPEmailCodeSent(ctx, e.Aggregate().ID, e.Aggregate().ResourceOwner)
	if err != nil {
		return nil, err
	}
	return handler.NewNoOpStatement(e), nil
}

func (u *userNotifier) reduceDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Drh5w", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
	}
	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, nil, user.AggregateType,
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

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.DomainClaimedMessageType)
		if err != nil {
			return err
		}

		ctx, origin, err := u.queries.Origin(ctx)
		if err != nil {
			return err
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
			return err
		}
		return u.commands.UserDomainClaimedSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reducePasswordlessCodeRequested(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordlessInitCodeRequestedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-EDtjd", "reduce.wrong.event.type %s", user.HumanPasswordlessInitCodeAddedType)
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

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordlessRegistrationMessageType)
		if err != nil {
			return err
		}

		ctx, origin, err := u.queries.Origin(ctx)
		if err != nil {
			return err
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
			return err
		}
		return u.commands.HumanPasswordlessInitCodeSent(ctx, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.ID)
	}), nil
}

func (u *userNotifier) reducePasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Yko2z8", "reduce.wrong.event.type %s", user.HumanPasswordChangedType)
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, nil, user.AggregateType, user.HumanPasswordChangeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}

		notificationPolicy, err := u.queries.NotificationPolicyByOrg(ctx, true, e.Aggregate().ResourceOwner, false)
		if errors.IsNotFound(err) {
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

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.PasswordChangeMessageType)
		if err != nil {
			return err
		}

		ctx, origin, err := u.queries.Origin(ctx)
		if err != nil {
			return err
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
			return err
		}
		return u.commands.PasswordChangeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) reducePhoneCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-He83g", "reduce.wrong.event.type %s", user.HumanPhoneCodeAddedType)
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

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.Aggregate().ID, false)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, notifyUser.ResourceOwner, domain.VerifyPhoneMessageType)
		if err != nil {
			return err
		}

		ctx, origin, err := u.queries.Origin(ctx)
		if err != nil {
			return err
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
			return err
		}
		return u.commands.HumanPhoneVerificationCodeSent(ctx, e.Aggregate().ResourceOwner, e.Aggregate().ID)
	}), nil
}

func (u *userNotifier) checkIfCodeAlreadyHandledOrExpired(ctx context.Context, event eventstore.Event, expiry time.Duration, data map[string]interface{}, eventTypes ...eventstore.EventType) (bool, error) {
	if event.CreatedAt().Add(expiry).Before(time.Now().UTC()) {
		return true, nil
	}
	return u.queries.IsAlreadyHandled(ctx, event, data, user.AggregateType, eventTypes...)
}
