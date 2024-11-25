package handlers

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func init() {
	RegisterSentHandler(user.HumanInitialCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, _ *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanInitCodeSent(ctx, orgID, id)
		},
	)
	RegisterSentHandler(user.HumanEmailCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanEmailVerificationCodeSent(ctx, orgID, id)
		},
	)
	RegisterSentHandler(user.HumanPasswordCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.PasswordCodeSent(ctx, orgID, id, generatorInfo)
		},
	)
	RegisterSentHandler(user.HumanOTPSMSCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanOTPSMSCodeSent(ctx, orgID, id, generatorInfo)
		},
	)
	RegisterSentHandler(session.OTPSMSChallengedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.OTPSMSSent(ctx, orgID, id, generatorInfo)
		},
	)
	RegisterSentHandler(user.HumanOTPEmailCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanOTPEmailCodeSent(ctx, orgID, id)
		},
	)
	RegisterSentHandler(session.OTPEmailChallengedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.OTPEmailSent(ctx, orgID, id)
		},
	)
	RegisterSentHandler(user.UserDomainClaimedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.UserDomainClaimedSent(ctx, orgID, id)
		},
	)
	RegisterSentHandler(user.HumanPasswordlessInitCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanPasswordlessInitCodeSent(ctx, id, orgID, args["CodeID"].(string))
		},
	)
	RegisterSentHandler(user.HumanPasswordChangedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.PasswordChangeSent(ctx, orgID, id)
		},
	)
	RegisterSentHandler(user.HumanPhoneCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanPhoneVerificationCodeSent(ctx, orgID, id, generatorInfo)
		},
	)
	RegisterSentHandler(user.HumanInviteCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, _ *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.InviteCodeSent(ctx, orgID, id)
		},
	)
	RegisterSentHandler(user.HumanPhoneCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanPhoneVerificationCodeSent(ctx, orgID, id, generatorInfo)
		},
	)
}

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
	otpEmailTmpl string,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &userNotifier{
		commands:     commands,
		queries:      queries,
		otpEmailTmpl: otpEmailTmpl,
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
				{
					Event:  user.HumanInviteCodeAddedType,
					Reduce: u.reduceInviteCodeAdded,
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
		return u.commands.RequestNotification(
			ctx,
			authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeEmail,
				domain.InitCodeMessageType,
			).
				WithURLTemplate(login.InitUserLinkTemplate(e.TriggeredAtOrigin, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.AuthRequestID)).
				WithCode(e.Code, e.Expiry).
				WithUnverifiedChannel(),
		)
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
		urlTmpl := e.URLTemplate
		if urlTmpl == "" {
			urlTmpl = login.MailVerificationLinkTemplate(e.TriggeredAtOrigin, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.AuthRequestID)
		}
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeEmail,
				domain.VerifyEmailMessageType,
			).
				WithURLTemplate(urlTmpl).
				WithCode(e.Code, e.Expiry).
				WithUnverifiedChannel(),
		)
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
		urlTmpl := e.URLTemplate
		if urlTmpl == "" {
			urlTmpl = login.InitPasswordLinkTemplate(e.TriggeredAtOrigin, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.AuthRequestID)
		}
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				e.NotificationType,
				domain.PasswordResetMessageType,
			).
				WithURLTemplate(urlTmpl).
				WithCode(e.Code, e.Expiry).
				WithUnverifiedChannel(),
		)
	}), nil
}

func (u *userNotifier) reduceOTPSMSCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASF3g", "reduce.wrong.event.type %s", user.HumanOTPSMSCodeAddedType)
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.HumanOTPSMSCodeAddedType,
			user.HumanOTPSMSCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeSms,
				domain.VerifySMSOTPMessageType,
			).
				WithCode(e.Code, e.Expiry).
				WithArgs(otpArgs(ctx, e.Expiry)).
				WithOTP(),
		)
	}), nil
}

func (u *userNotifier) reduceSessionOTPSMSChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.OTPSMSChallengedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Sk32L", "reduce.wrong.event.type %s", session.OTPSMSChallengedType)
	}
	if e.CodeReturned {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			session.OTPSMSChallengedType,
			session.OTPSMSSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		s, err := u.queries.SessionByID(ctx, true, e.Aggregate().ID, "")
		if err != nil {
			return err
		}
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				s.UserFactor.UserID,
				s.UserFactor.ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeSms,
				domain.VerifySMSOTPMessageType,
			).
				WithAggregate(e.Aggregate().ID, e.Aggregate().ResourceOwner).
				WithCode(e.Code, e.Expiry).
				WithArgs(otpArgs(ctx, e.Expiry)),
		)
	}), nil
}

func (u *userNotifier) reduceOTPEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JL3hw", "reduce.wrong.event.type %s", user.HumanOTPEmailCodeAddedType)
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.HumanOTPEmailCodeAddedType,
			user.HumanOTPEmailCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		var authRequestID string
		if e.AuthRequestInfo != nil {
			authRequestID = e.AuthRequestInfo.ID
		}
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeEmail,
				domain.VerifyEmailOTPMessageType,
			).
				WithURLTemplate(login.OTPLinkTemplate(e.TriggeredAtOrigin, authRequestID, domain.MFATypeOTPEmail)).
				WithCode(e.Code, e.Expiry).
				WithArgs(otpArgs(ctx, e.Expiry)),
		)
	}), nil
}

func (u *userNotifier) reduceSessionOTPEmailChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.OTPEmailChallengedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-zbsgt", "reduce.wrong.event.type %s", session.OTPEmailChallengedType)
	}
	if e.ReturnCode {
		return handler.NewNoOpStatement(e), nil
	}
	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			session.OTPEmailChallengedType,
			session.OTPEmailSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		s, err := u.queries.SessionByID(ctx, true, e.Aggregate().ID, "")
		if err != nil {
			return err
		}
		urlTmpl := e.TriggeredAtOrigin + u.otpEmailTmpl
		if e.URLTmpl != "" {
			urlTmpl = e.URLTmpl
		}
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				s.UserFactor.UserID,
				s.UserFactor.ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeEmail,
				domain.VerifyEmailOTPMessageType,
			).
				WithAggregate(e.Aggregate().ID, e.Aggregate().ResourceOwner).
				WithURLTemplate(urlTmpl).
				WithCode(e.Code, e.Expiry).
				WithArgs(otpArgs(ctx, e.Expiry)),
		)
	}), nil
}

func otpArgs(ctx context.Context, expiry time.Duration) map[string]interface{} {
	domainCtx := http_util.DomainContext(ctx)
	args := make(map[string]interface{})
	args["Origin"] = domainCtx.Origin()
	args["Domain"] = domainCtx.RequestedDomain()
	args["Expiry"] = expiry
	return args
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
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeEmail,
				domain.DomainClaimedMessageType,
			).
				WithURLTemplate(login.LoginLink(e.TriggeredAtOrigin, e.Aggregate().ResourceOwner)).
				WithUnverifiedChannel().
				WithArgs(map[string]any{
					"TempUsername": e.UserName,
				}),
		)
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
		urlTmpl := e.URLTemplate
		if urlTmpl == "" {
			urlTmpl = domain.PasswordlessInitCodeLinkTemplate(e.TriggeredAtOrigin+login.HandlerPrefix+login.EndpointPasswordlessRegistration, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.ID)
		}
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeEmail,
				domain.PasswordlessRegistrationMessageType,
			).
				WithURLTemplate(urlTmpl).
				WithCode(e.Code, e.Expiry),
		)
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

		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeEmail,
				domain.PasswordChangeMessageType,
			).
				WithURLTemplate(console.LoginHintLink(e.TriggeredAtOrigin, "{{.PreferredLoginName}}")).
				WithUnverifiedChannel(),
		)
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
		ctx, err = enrichCtx(ctx, e.TriggeredAtOrigin)
		if err != nil {
			return err
		}
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeSms,
				domain.VerifyPhoneMessageType,
			).
				WithCode(e.Code, e.Expiry).
				WithUnverifiedChannel().
				WithArgs(map[string]any{
					"Domain": http_util.DomainContext(ctx).RequestedDomain(),
				}),
		)
	}), nil
}

func (u *userNotifier) reduceInviteCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanInviteCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Eeg3s", "reduce.wrong.event.type %s", user.HumanInviteCodeAddedType)
	}
	if e.CodeReturned {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.HumanInviteCodeAddedType, user.HumanInviteCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		applicationName := e.ApplicationName
		if applicationName == "" {
			applicationName = "ZITADEL"
		}
		urlTmpl := e.URLTemplate
		if urlTmpl == "" {
			urlTmpl = login.InviteUserLinkTemplate(e.TriggeredAtOrigin, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.AuthRequestID)
		}
		return u.commands.RequestNotification(ctx, authz.GetInstance(ctx).InstanceID(),
			command.NewNotificationRequest(
				e.Aggregate().ID,
				e.Aggregate().ResourceOwner,
				e.TriggeredAtOrigin,
				e.EventType,
				domain.NotificationTypeEmail,
				domain.InviteUserMessageType,
			).
				WithURLTemplate(urlTmpl).
				WithCode(e.Code, e.Expiry).
				WithUnverifiedChannel().
				WithArgs(map[string]any{
					"ApplicationName": applicationName,
				}),
		)
	}), nil
}

func (u *userNotifier) checkIfCodeAlreadyHandledOrExpired(ctx context.Context, event eventstore.Event, expiry time.Duration, data map[string]interface{}, eventTypes ...eventstore.EventType) (bool, error) {
	if expiry > 0 && event.CreatedAt().Add(expiry).Before(time.Now().UTC()) {
		return true, nil
	}
	return u.queries.IsAlreadyHandled(ctx, event, data, eventTypes...)
}
