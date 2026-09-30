package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

//go:generate mockgen -typed -package mock -destination ./mock/commands.mock.go . Commands
type Commands interface {
	HumanInitCodeSent(ctx context.Context, orgID, userID string, deliverySuppressed bool) error
	HumanEmailVerificationCodeSent(ctx context.Context, orgID, userID string, deliverySuppressed bool) error
	PasswordCodeSent(ctx context.Context, orgID, userID string, generatorInfo *senders.CodeGeneratorInfo, deliverySuppressed bool) error
	HumanOTPSMSCodeSent(ctx context.Context, userID, resourceOwner string, generatorInfo *senders.CodeGeneratorInfo) error
	HumanOTPEmailCodeSent(ctx context.Context, userID, resourceOwner string, deliverySuppressed bool) error
	OTPSMSSent(ctx context.Context, sessionID, resourceOwner string, generatorInfo *senders.CodeGeneratorInfo) error
	OTPEmailSent(ctx context.Context, sessionID, resourceOwner string, deliverySuppressed bool) error
	UserDomainClaimedSent(ctx context.Context, orgID, userID string, deliverySuppressed bool) error
	HumanPasswordlessInitCodeSent(ctx context.Context, userID, resourceOwner, codeID string, deliverySuppressed bool) error
	PasswordChangeSent(ctx context.Context, orgID, userID string, deliverySuppressed bool) error
	HumanPhoneVerificationCodeSent(ctx context.Context, orgID, userID string, generatorInfo *senders.CodeGeneratorInfo) error
	InviteCodeSent(ctx context.Context, userID, orgID string, deliverySuppressed bool) error
	UsageNotificationSent(ctx context.Context, dueEvent *quota.NotificationDueEvent) error
	MilestonePushed(ctx context.Context, instanceID string, msType milestone.Type, endpoints []string) error
	BackChannelLogoutSent(ctx context.Context, id, oidcSessionID, instanceID string) (err error)
}
