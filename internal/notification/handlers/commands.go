package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

//go:generate mockgen -typed -package mock -destination ./mock/commands.mock.go . Commands
type Commands interface {
	HumanInitCodeSent(ctx context.Context, orgID, userID string, deliveryInfo senders.DeliveryInfo) error
	HumanEmailVerificationCodeSent(ctx context.Context, orgID, userID string, deliveryInfo senders.DeliveryInfo) error
	PasswordCodeSent(ctx context.Context, orgID, userID string, generatorInfo *senders.CodeGeneratorInfo, deliveryInfo senders.DeliveryInfo) error
	HumanOTPSMSCodeSent(ctx context.Context, userID, resourceOwner string, generatorInfo *senders.CodeGeneratorInfo) error
	HumanOTPEmailCodeSent(ctx context.Context, userID, resourceOwner string, deliveryInfo senders.DeliveryInfo) error
	OTPSMSSent(ctx context.Context, sessionID, resourceOwner string, generatorInfo *senders.CodeGeneratorInfo) error
	OTPEmailSent(ctx context.Context, sessionID, resourceOwner string, deliveryInfo senders.DeliveryInfo) error
	UserDomainClaimedSent(ctx context.Context, orgID, userID string, deliveryInfo senders.DeliveryInfo) error
	HumanPasswordlessInitCodeSent(ctx context.Context, userID, resourceOwner, codeID string, deliveryInfo senders.DeliveryInfo) error
	PasswordChangeSent(ctx context.Context, orgID, userID string, deliveryInfo senders.DeliveryInfo) error
	HumanPhoneVerificationCodeSent(ctx context.Context, orgID, userID string, generatorInfo *senders.CodeGeneratorInfo) error
	InviteCodeSent(ctx context.Context, userID, orgID string, deliveryInfo senders.DeliveryInfo) error
	NotificationRejected(ctx context.Context, aggregate *eventstore.Aggregate, rejection channels.Rejection) error
	UsageNotificationSent(ctx context.Context, dueEvent *quota.NotificationDueEvent) error
	MilestonePushed(ctx context.Context, instanceID string, msType milestone.Type, endpoints []string) error
	BackChannelLogoutSent(ctx context.Context, id, oidcSessionID, instanceID string) (err error)
}
