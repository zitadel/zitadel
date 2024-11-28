package handlers

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type Commands interface {
	RequestNotification(ctx context.Context, instanceID string, request *command.NotificationRequest) error
	NotificationCanceled(ctx context.Context, tx *sql.Tx, id, resourceOwner string, err error) error
	NotificationRetryRequested(ctx context.Context, tx *sql.Tx, id, resourceOwner string, request *command.NotificationRetryRequest, err error) error
	NotificationSent(ctx context.Context, tx *sql.Tx, id, instanceID string) error
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
