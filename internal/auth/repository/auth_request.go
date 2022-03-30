package repository

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
)

type AuthRequestRepository interface {
	CreateAuthRequest(ctx context.Context, request *domain.AuthRequest) (*domain.AuthRequest, error)
	AuthRequestByID(ctx context.Context, id, userAgentID, instanceID string) (*domain.AuthRequest, error)
	AuthRequestByIDCheckLoggedIn(ctx context.Context, id, userAgentID, instanceID string) (*domain.AuthRequest, error)
	AuthRequestByCode(ctx context.Context, code, instanceID string) (*domain.AuthRequest, error)
	SaveAuthCode(ctx context.Context, id, code, userAgentID, instanceID string) error
	DeleteAuthRequest(ctx context.Context, id, instanceID string) error

	CheckLoginName(ctx context.Context, id, loginName, userAgentID, instanceID string) error
	CheckExternalUserLogin(ctx context.Context, authReqID, userAgentID, instanceID string, user *domain.ExternalUser, info *domain.BrowserInfo) error
	SetExternalUserLogin(ctx context.Context, authReqID, userAgentID, instanceID string, user *domain.ExternalUser) error
	SelectUser(ctx context.Context, id, userID, userAgentID, instanceID string) error
	SelectExternalIDP(ctx context.Context, authReqID, idpConfigID, userAgentID, instanceID string) error
	VerifyPassword(ctx context.Context, id, userID, resourceOwner, password, userAgentID, instanceID string, info *domain.BrowserInfo) error

	VerifyMFAOTP(ctx context.Context, authRequestID, userID, resourceOwner, code, userAgentID, instanceID string, info *domain.BrowserInfo) error
	BeginMFAU2FLogin(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID, instanceID string) (*domain.WebAuthNLogin, error)
	VerifyMFAU2F(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID, instanceID string, credentialData []byte, info *domain.BrowserInfo) error
	BeginPasswordlessSetup(ctx context.Context, instanceID, userID, resourceOwner string, preferredPlatformType domain.AuthenticatorAttachment) (login *domain.WebAuthNToken, err error)
	VerifyPasswordlessSetup(ctx context.Context, userID, resourceOwner, userAgentID, tokenName string, credentialData []byte) (err error)
	BeginPasswordlessInitCodeSetup(ctx context.Context, instanceID, userID, resourceOwner, codeID, verificationCode string, preferredPlatformType domain.AuthenticatorAttachment) (login *domain.WebAuthNToken, err error)
	VerifyPasswordlessInitCodeSetup(ctx context.Context, userID, resourceOwner, userAgentID, tokenName, codeID, verificationCode string, credentialData []byte) (err error)
	BeginPasswordlessLogin(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID, instanceID string) (*domain.WebAuthNLogin, error)
	VerifyPasswordless(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID, instanceID string, credentialData []byte, info *domain.BrowserInfo) error

	LinkExternalUsers(ctx context.Context, authReqID, userAgentID, instanceID string, info *domain.BrowserInfo) error
	AutoRegisterExternalUser(ctx context.Context, user *domain.Human, externalIDP *domain.UserIDPLink, orgMemberRoles []string, authReqID, userAgentID, resourceOwner, instanceID string, metadatas []*domain.Metadata, info *domain.BrowserInfo) error
	ResetLinkingUsers(ctx context.Context, authReqID, userAgentID, instanceID string) error
}
