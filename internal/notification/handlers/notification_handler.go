package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/repository/notification"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type notificationHandler struct {
	commands Commands
	queries  *NotificationQueries
	channels types.ChannelChains
	senders  map[eventstore.EventType]Sent
}

type Sent func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error

func NewNotificationHandler(
	ctx context.Context,
	config handler.Config,
	commands Commands,
	queries *NotificationQueries,
	channels types.ChannelChains,
) *handler.Handler {
	x := &notificationHandler{
		commands: commands,
		queries:  queries,
		channels: channels,
		senders:  make(map[eventstore.EventType]Sent),
	}
	test(x)
	return handler.NewHandler(ctx, &config, x)
}

func (u *notificationHandler) Name() string {
	return UserNotificationsProjectionTable + "_worker" //TODO: remove
}

func (u *notificationHandler) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: notification.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  notification.RequestedType,
					Reduce: u.reduceNotificationAdded,
				},
			},
		},
	}
}

func (u *notificationHandler) RegisterSender(eventType eventstore.EventType, sent Sent) {
	u.senders[eventType] = sent
}

func (u *notificationHandler) reduceNotificationAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*notification.RequestedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SFA3gs", "reduce.wrong.event.type %s", notification.RequestedType)
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) (err error) {
		ctx := HandlerContext(event.Aggregate())
		var code string
		if e.Code != nil {
			code, err = crypto.DecryptString(e.Code, u.queries.UserDataCrypto)
			if err != nil {
				return err
			}
		}
		colors, err := u.queries.ActiveLabelPolicyByOrg(ctx, e.UserResourceOwner, false)
		if err != nil {
			return err
		}

		notifyUser, err := u.queries.GetNotifyUserByID(ctx, true, e.UserID)
		if err != nil {
			return err
		}
		translator, err := u.queries.GetTranslatorWithOrgTexts(ctx, e.UserResourceOwner, e.MessageType)
		if err != nil {
			return err
		}

		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}

		args := e.Args
		if len(args) == 0 {
			args = make(map[string]any)
		}
		if code != "" {
			args["Code"] = code
			// existing notifications use `OTP` as argument for the code
			if e.IsOTP {
				args["OTP"] = code
			}
		}

		generatorInfo := new(senders.CodeGeneratorInfo)
		var notify types.Notify
		switch e.NotificationType {
		case domain.NotificationTypeEmail:
			template, err := u.queries.MailTemplateByOrg(ctx, notifyUser.ResourceOwner, false)
			if err != nil {
				return err
			}

			notify = types.SendEmail(ctx, u.channels, string(template.Template), translator, notifyUser, colors, e)
		case domain.NotificationTypeSms:
			notify = types.SendSMS(ctx, u.channels, translator, notifyUser, colors, e, generatorInfo)
		}
		if err := notify(e.URLTemplate, args, e.MessageType, e.UnverifiedNotificationChannel); err != nil {
			return u.commands.NotificationFailed(ctx, e.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
		}
		sender, ok := u.senders[e.EventType]
		if !ok {
			return nil
		}
		//if err := sender.Send(ctx, notify, notifyUser, code); err != nil {
		//	return err
		//}
		return sender(ctx, u.commands, e.NotificationAggregateID(), e.NotificationAggregateResourceOwner(), generatorInfo, args)
	}), nil
}

func test(x *notificationHandler) {
	x.RegisterSender(user.HumanInitialCodeAddedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code, urlTmpl string) error {
		//	return notify.SendUserInitCode(ctx, notifyUser, code, e.AuthRequestID)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, _ *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanInitCodeSent(ctx, orgID, id)
		},
	)
	x.RegisterSender(user.HumanEmailCodeAddedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendEmailVerificationCode(ctx, notifyUser, code, e.URLTemplate, e.AuthRequestID)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanEmailVerificationCodeSent(ctx, orgID, id)
		},
	)
	x.RegisterSender(user.HumanPasswordCodeAddedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendPasswordCode(ctx, notifyUser, code, e.URLTemplate, e.AuthRequestID)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.PasswordCodeSent(ctx, orgID, id, generatorInfo)
		},
	)
	x.RegisterSender(user.HumanOTPSMSCodeAddedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendOTPSMSCode(ctx, code, e.CodeExpiry)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanOTPSMSCodeSent(ctx, orgID, id, generatorInfo)
		},
	)
	x.RegisterSender(session.OTPSMSChallengedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendOTPSMSCode(ctx, code, e.CodeExpiry)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.OTPSMSSent(ctx, orgID, id, generatorInfo)
		},
	)
	x.RegisterSender(user.HumanOTPEmailCodeAddedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendOTPEmailCode(ctx, url, code, e.CodeExpiry)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanOTPEmailCodeSent(ctx, orgID, id)
		},
	)
	x.RegisterSender(session.OTPEmailChallengedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendOTPEmailCode(ctx, url, code, e.CodeExpiry)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.OTPEmailSent(ctx, orgID, id)
		},
	)
	x.RegisterSender(user.UserDomainClaimedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendDomainClaimed(ctx, notifyUser, e.UserName)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.UserDomainClaimedSent(ctx, orgID, id)
		},
	)
	x.RegisterSender(user.HumanPasswordlessInitCodeAddedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendPasswordlessRegistrationLink(ctx, notifyUser, e.CodeID, e.URLTemplate)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanPasswordlessInitCodeSent(ctx, id, orgID, args["CodeID"].(string))
		},
	)
	x.RegisterSender(user.HumanPasswordChangedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendPasswordChange(ctx, notifyUser)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.PasswordChangeSent(ctx, orgID, id)
		},
	)
	x.RegisterSender(user.HumanPhoneCodeAddedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendPhoneVerificationCode(ctx, code)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanPhoneVerificationCodeSent(ctx, orgID, id, generatorInfo)
		},
	)
	x.RegisterSender(user.HumanInviteCodeAddedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendInviteCode(ctx, code, e.ApplicationName, e.URLTemplate, e.AuthRequestID)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, _ *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.InviteCodeSent(ctx, orgID, id)
		},
	)
	x.RegisterSender(user.HumanPhoneCodeAddedType,
		//func(ctx context.Context, notify types.Notify, notifyUser *query.NotifyUser, code string) error {
		//	return notify.SendPhoneVerificationCode(ctx, code)
		//},
		func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error {
			return commands.HumanPhoneVerificationCodeSent(ctx, orgID, id, generatorInfo)
		},
	)
}
