package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type Commands interface {
	HumanInitCodeSent(ctx context.Context, orgID, userID string) error
	HumanEmailVerificationCodeSent(ctx context.Context, orgID, userID string) error
	PasswordCodeSent(ctx context.Context, orgID, userID string, generatorInfo *senders.CodeGeneratorInfo) error
	HumanOTPSMSCodeSent(ctx context.Context, userID, resourceOwner string, generatorInfo *senders.CodeGeneratorInfo) error
	HumanOTPEmailCodeSent(ctx context.Context, userID, resourceOwner string) error
	OTPSMSSent(ctx context.Context, sessionID, resourceOwner string, generatorInfo *senders.CodeGeneratorInfo) error
	OTPEmailSent(ctx context.Context, sessionID, resourceOwner string) error
	UserDomainClaimedSent(ctx context.Context, orgID, userID string) error
	HumanPasswordlessInitCodeSent(ctx context.Context, userID, resourceOwner, codeID string) error
	PasswordChangeSent(ctx context.Context, orgID, userID string) error
	HumanPhoneVerificationCodeSent(ctx context.Context, orgID, userID string, generatorInfo *senders.CodeGeneratorInfo) error
	InviteCodeSent(ctx context.Context, orgID, userID string) error
	UsageNotificationSent(ctx context.Context, dueEvent *quota.NotificationDueEvent) error
	MilestonePushed(ctx context.Context, instanceID string, msType milestone.Type, endpoints []string) error
}
