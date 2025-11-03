package handlers

import (
	"context"
	"time"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/queue"
	"github.com/zitadel/zitadel/internal/repository/notification"
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
			return commands.HumanOTPSMSCodeSent(ctx, id, orgID, generatorInfo)
		},
	)
	RegisterSentHandler(session.OTPSMSChallengedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.OTPSMSSent(ctx, id, orgID, generatorInfo)
		},
	)
	RegisterSentHandler(user.HumanOTPEmailCodeAddedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanOTPEmailCodeSent(ctx, id, orgID)
		},
	)
	RegisterSentHandler(session.OTPEmailChallengedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.OTPEmailSent(ctx, id, orgID)
		},
	)
	RegisterSentHandler(user.UserDomainClaimedType,
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.UserDomainClaimedSent(ctx, orgID, id)
		},
	)
	RegisterSentHandler(user.HumanPasswordlessInitCodeRequestedType,
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
}

const (
	UserNotificationsProjectionTable = "projections.notifications"
)

type userNotifier struct {
	queries      *NotificationQueries
	otpEmailTmpl string

	queue       Queue
	maxAttempts uint8
}

func NewUserNotifier(
	ctx context.Context,
	config handler.Config,
	commands Commands,
	queries *NotificationQueries,
	channels types.ChannelChains,
	otpEmailTmpl string,
	workerConfig WorkerConfig,
	queue Queue,
) *handler.Handler {
	if workerConfig.LegacyEnabled {
		return NewUserNotifierLegacy(ctx, config, commands, queries, channels, otpEmailTmpl)
	}
	return handler.NewHandler(ctx, &config, &userNotifier{
		queries:      queries,
		otpEmailTmpl: otpEmailTmpl,
		queue:        queue,
		maxAttempts:  workerConfig.MaxAttempts,
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
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.UserV1InitialCodeAddedType, user.UserV1InitialCodeSentType,
			user.HumanInitialCodeAddedType, user.HumanInitialCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		origin := http_util.DomainContext(ctx).Origin()
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:                     e.Aggregate(),
				UserID:                        e.Aggregate().ID,
				UserResourceOwner:             e.Aggregate().ResourceOwner,
				TriggeredAtOrigin:             origin,
				EventType:                     e.EventType,
				NotificationType:              domain.NotificationTypeEmail,
				MessageType:                   domain.InitCodeMessageType,
				Code:                          e.Code,
				CodeExpiry:                    e.Expiry,
				IsOTP:                         false,
				UnverifiedNotificationChannel: true,
				URLTemplate:                   login.InitUserLinkTemplate(origin, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.AuthRequestID),
				Args: &domain.NotificationArguments{
					AuthRequestID: e.AuthRequestID,
				},
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.UserV1EmailCodeAddedType, user.UserV1EmailCodeSentType,
			user.HumanEmailCodeAddedType, user.HumanEmailCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		origin := http_util.DomainContext(ctx).Origin()
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:                     e.Aggregate(),
				UserID:                        e.Aggregate().ID,
				UserResourceOwner:             e.Aggregate().ResourceOwner,
				TriggeredAtOrigin:             origin,
				EventType:                     e.EventType,
				NotificationType:              domain.NotificationTypeEmail,
				MessageType:                   domain.VerifyEmailMessageType,
				Code:                          e.Code,
				CodeExpiry:                    e.Expiry,
				IsOTP:                         false,
				UnverifiedNotificationChannel: true,
				URLTemplate:                   u.emailCodeTemplate(origin, e),
				Args: &domain.NotificationArguments{
					AuthRequestID: e.AuthRequestID,
				},
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
		)
	}), nil
}

func (u *userNotifier) emailCodeTemplate(origin string, e *user.HumanEmailCodeAddedEvent) string {
	if e.URLTemplate != "" {
		return e.URLTemplate
	}
	return login.MailVerificationLinkTemplate(origin, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.AuthRequestID)
}

func (u *userNotifier) reducePasswordCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Eeg3s", "reduce.wrong.event.type %s", user.HumanPasswordCodeAddedType)
	}
	if e.CodeReturned {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.UserV1PasswordCodeAddedType, user.UserV1PasswordCodeSentType,
			user.HumanPasswordCodeAddedType, user.HumanPasswordCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		origin := http_util.DomainContext(ctx).Origin()
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:                     e.Aggregate(),
				UserID:                        e.Aggregate().ID,
				UserResourceOwner:             e.Aggregate().ResourceOwner,
				TriggeredAtOrigin:             origin,
				EventType:                     e.EventType,
				NotificationType:              e.NotificationType,
				MessageType:                   domain.PasswordResetMessageType,
				Code:                          e.Code,
				CodeExpiry:                    e.Expiry,
				IsOTP:                         false,
				UnverifiedNotificationChannel: true,
				URLTemplate:                   u.passwordCodeTemplate(origin, e),
				Args: &domain.NotificationArguments{
					AuthRequestID: e.AuthRequestID,
				},
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
		)
	}), nil
}

func (u *userNotifier) passwordCodeTemplate(origin string, e *user.HumanPasswordCodeAddedEvent) string {
	if e.URLTemplate != "" {
		return e.URLTemplate
	}
	return login.InitPasswordLinkTemplate(origin, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.AuthRequestID)
}

func (u *userNotifier) reduceOTPSMSCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASF3g", "reduce.wrong.event.type %s", user.HumanOTPSMSCodeAddedType)
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.HumanOTPSMSCodeAddedType,
			user.HumanOTPSMSCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:         e.Aggregate(),
				UserID:            e.Aggregate().ID,
				UserResourceOwner: e.Aggregate().ResourceOwner,
				TriggeredAtOrigin: http_util.DomainContext(ctx).Origin(),
				EventType:         e.EventType,
				NotificationType:  domain.NotificationTypeSms,
				MessageType:       domain.VerifySMSOTPMessageType,
				Code:              e.Code,
				CodeExpiry:        e.Expiry,
				IsOTP:             true,
				Args:              otpArgs(ctx, e.Expiry),
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			session.OTPSMSChallengedType,
			session.OTPSMSSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}

		sessionWriteModel := command.NewSessionWriteModel(e.Aggregate().ID, e.Aggregate().InstanceID)
		err = u.queries.es.FilterToQueryReducer(ctx, sessionWriteModel)
		if err != nil {
			return err
		}

		args := otpArgs(ctx, e.Expiry)
		args.SessionID = e.Aggregate().ID
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:         e.Aggregate(),
				UserID:            sessionWriteModel.UserID,
				UserResourceOwner: sessionWriteModel.UserResourceOwner,
				TriggeredAtOrigin: http_util.DomainContext(ctx).Origin(),
				EventType:         e.EventType,
				NotificationType:  domain.NotificationTypeSms,
				MessageType:       domain.VerifySMSOTPMessageType,
				Code:              e.Code,
				CodeExpiry:        e.Expiry,
				IsOTP:             true,
				Args:              args,
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
		)
	}), nil
}

func (u *userNotifier) reduceOTPEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JL3hw", "reduce.wrong.event.type %s", user.HumanOTPEmailCodeAddedType)
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.HumanOTPEmailCodeAddedType,
			user.HumanOTPEmailCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		origin := http_util.DomainContext(ctx).Origin()
		var authRequestID string
		if e.AuthRequestInfo != nil {
			authRequestID = e.AuthRequestInfo.ID
		}
		args := otpArgs(ctx, e.Expiry)
		args.AuthRequestID = authRequestID
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:         e.Aggregate(),
				UserID:            e.Aggregate().ID,
				UserResourceOwner: e.Aggregate().ResourceOwner,
				TriggeredAtOrigin: origin,
				EventType:         e.EventType,
				NotificationType:  domain.NotificationTypeEmail,
				MessageType:       domain.VerifyEmailOTPMessageType,
				Code:              e.Code,
				CodeExpiry:        e.Expiry,
				IsOTP:             true,
				URLTemplate:       login.OTPLinkTemplate(origin, authRequestID, domain.MFATypeOTPEmail),
				Args:              args,
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
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
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			session.OTPEmailChallengedType,
			session.OTPEmailSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		s, err := u.queries.SessionByID(ctx, true, e.Aggregate().ID, "", nil)
		if err != nil {
			return err
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		origin := http_util.DomainContext(ctx).Origin()

		args := otpArgs(ctx, e.Expiry)
		args.SessionID = e.Aggregate().ID
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:         e.Aggregate(),
				UserID:            s.UserFactor.UserID,
				UserResourceOwner: s.UserFactor.ResourceOwner,
				TriggeredAtOrigin: origin,
				EventType:         e.EventType,
				NotificationType:  domain.NotificationTypeEmail,
				MessageType:       domain.VerifyEmailOTPMessageType,
				Code:              e.Code,
				CodeExpiry:        e.Expiry,
				IsOTP:             true,
				URLTemplate:       u.otpEmailTemplate(origin, e),
				Args:              args,
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
		)
	}), nil
}

func (u *userNotifier) otpEmailTemplate(origin string, e *session.OTPEmailChallengedEvent) string {
	if e.URLTmpl != "" {
		return e.URLTmpl
	}
	return origin + u.otpEmailTmpl
}

func otpArgs(ctx context.Context, expiry time.Duration) *domain.NotificationArguments {
	domainCtx := http_util.DomainContext(ctx)
	return &domain.NotificationArguments{
		Origin: domainCtx.Origin(),
		Domain: domainCtx.RequestedDomain(),
		Expiry: expiry,
	}
}

func (u *userNotifier) reduceDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Drh5w", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
	}
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, nil,
			user.UserDomainClaimedType, user.UserDomainClaimedSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		origin := http_util.DomainContext(ctx).Origin()
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:                     e.Aggregate(),
				UserID:                        e.Aggregate().ID,
				UserResourceOwner:             e.Aggregate().ResourceOwner,
				TriggeredAtOrigin:             origin,
				EventType:                     e.EventType,
				NotificationType:              domain.NotificationTypeEmail,
				MessageType:                   domain.DomainClaimedMessageType,
				URLTemplate:                   login.LoginLink(origin, e.Aggregate().ResourceOwner),
				UnverifiedNotificationChannel: true,
				Args: &domain.NotificationArguments{
					TempUsername: e.UserName,
				},
				RequiresPreviousDomain: true,
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, map[string]interface{}{"id": e.ID}, user.HumanPasswordlessInitCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		origin := http_util.DomainContext(ctx).Origin()
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:         e.Aggregate(),
				UserID:            e.Aggregate().ID,
				UserResourceOwner: e.Aggregate().ResourceOwner,
				TriggeredAtOrigin: origin,
				EventType:         e.EventType,
				NotificationType:  domain.NotificationTypeEmail,
				MessageType:       domain.PasswordlessRegistrationMessageType,
				URLTemplate:       u.passwordlessCodeTemplate(origin, e),
				Args: &domain.NotificationArguments{
					CodeID: e.ID,
				},
				CodeExpiry: e.Expiry,
				Code:       e.Code,
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
		)
	}), nil
}

func (u *userNotifier) passwordlessCodeTemplate(origin string, e *user.HumanPasswordlessInitCodeRequestedEvent) string {
	if e.URLTemplate != "" {
		return e.URLTemplate
	}
	return domain.PasswordlessInitCodeLinkTemplate(origin+login.HandlerPrefix+login.EndpointPasswordlessRegistration, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.ID)
}

func (u *userNotifier) reducePasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yko2z8", "reduce.wrong.event.type %s", user.HumanPasswordChangedType)
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.queries.IsAlreadyHandled(ctx, event, nil, user.HumanPasswordChangeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}

		notificationPolicy, err := u.queries.NotificationPolicyByOrg(ctx, true, e.Aggregate().ResourceOwner, false)
		if err != nil && !zerrors.IsNotFound(err) {
			return err
		}

		if !notificationPolicy.PasswordChange {
			return nil
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		origin := http_util.DomainContext(ctx).Origin()

		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:                     e.Aggregate(),
				UserID:                        e.Aggregate().ID,
				UserResourceOwner:             e.Aggregate().ResourceOwner,
				TriggeredAtOrigin:             origin,
				EventType:                     e.EventType,
				NotificationType:              domain.NotificationTypeEmail,
				MessageType:                   domain.PasswordChangeMessageType,
				URLTemplate:                   console.LoginHintLink(origin, "{{.PreferredLoginName}}"),
				UnverifiedNotificationChannel: true,
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.UserV1PhoneCodeAddedType, user.UserV1PhoneCodeSentType,
			user.HumanPhoneCodeAddedType, user.HumanPhoneCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}

		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:                     e.Aggregate(),
				UserID:                        e.Aggregate().ID,
				UserResourceOwner:             e.Aggregate().ResourceOwner,
				TriggeredAtOrigin:             http_util.DomainContext(ctx).Origin(),
				EventType:                     e.EventType,
				NotificationType:              domain.NotificationTypeSms,
				MessageType:                   domain.VerifyPhoneMessageType,
				CodeExpiry:                    e.Expiry,
				Code:                          e.Code,
				UnverifiedNotificationChannel: true,
				Args: &domain.NotificationArguments{
					Domain: http_util.DomainContext(ctx).RequestedDomain(),
				},
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx = HandlerContext(ctx, event.Aggregate())
		alreadyHandled, err := u.checkIfCodeAlreadyHandledOrExpired(ctx, event, e.Expiry, nil,
			user.HumanInviteCodeAddedType, user.HumanInviteCodeSentType)
		if err != nil {
			return err
		}
		if alreadyHandled {
			return nil
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		origin := http_util.DomainContext(ctx).Origin()

		applicationName := e.ApplicationName
		if applicationName == "" {
			applicationName = "ZITADEL"
		}
		return u.queue.Insert(ctx,
			&notification.Request{
				Aggregate:                     e.Aggregate(),
				UserID:                        e.Aggregate().ID,
				UserResourceOwner:             e.Aggregate().ResourceOwner,
				TriggeredAtOrigin:             origin,
				EventType:                     e.EventType,
				NotificationType:              domain.NotificationTypeEmail,
				MessageType:                   domain.InviteUserMessageType,
				CodeExpiry:                    e.Expiry,
				Code:                          e.Code,
				UnverifiedNotificationChannel: true,
				URLTemplate:                   u.inviteCodeTemplate(origin, e),
				Args: &domain.NotificationArguments{
					AuthRequestID:   e.AuthRequestID,
					ApplicationName: applicationName,
				},
			},
			queue.WithQueueName(notification.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
		)
	}), nil
}

func (u *userNotifier) inviteCodeTemplate(origin string, e *user.HumanInviteCodeAddedEvent) string {
	if e.URLTemplate != "" {
		return e.URLTemplate
	}
	return login.InviteUserLinkTemplate(origin, e.Aggregate().ID, e.Aggregate().ResourceOwner, e.AuthRequestID)
}

func (u *userNotifier) checkIfCodeAlreadyHandledOrExpired(ctx context.Context, event eventstore.Event, expiry time.Duration, data map[string]interface{}, eventTypes ...eventstore.EventType) (bool, error) {
	if expiry > 0 && event.CreatedAt().Add(expiry).Before(time.Now().UTC()) {
		return true, nil
	}
	return u.queries.IsAlreadyHandled(ctx, event, data, eventTypes...)
}
