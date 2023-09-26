package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type Commands interface {
	HumanInitCodeSent(ctx context.Context, orgID, userID string) (err error)
	HumanEmailVerificationCodeSent(ctx context.Context, orgID, userID string) (err error)
	PasswordCodeSent(ctx context.Context, orgID, userID string) (err error)
	HumanOTPSMSCodeSent(ctx context.Context, userID, resourceOwner string) (err error)
	HumanOTPEmailCodeSent(ctx context.Context, userID, resourceOwner string) (err error)
	OTPSMSSent(ctx context.Context, sessionID, resourceOwner string) error
	OTPEmailSent(ctx context.Context, sessionID, resourceOwner string) error
	UserDomainClaimedSent(ctx context.Context, orgID, userID string) (err error)
	HumanPasswordlessInitCodeSent(ctx context.Context, userID, resourceOwner, codeID string) error
	PasswordChangeSent(ctx context.Context, orgID, userID string) (err error)
	HumanPhoneVerificationCodeSent(ctx context.Context, orgID, userID string) (err error)
	UsageNotificationSent(ctx context.Context, dueEvent *quota.NotificationDueEvent) error
	MilestonePushed(ctx context.Context, msType milestone.Type, endpoints []string, primaryDomain string) error
}
