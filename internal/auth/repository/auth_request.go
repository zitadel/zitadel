package repository

import (
	"context"
	"github.com/caos/zitadel/internal/auth_request/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	user_model "github.com/caos/zitadel/internal/user/model"
)

type AuthRequestRepository interface {
	CreateAuthRequest(ctx context.Context, request *model.AuthRequest) (*model.AuthRequest, error)
	AuthRequestByID(ctx context.Context, id, userAgentID string) (*model.AuthRequest, error)
	AuthRequestByIDCheckLoggedIn(ctx context.Context, id, userAgentID string) (*model.AuthRequest, error)
	AuthRequestByCode(ctx context.Context, code string) (*model.AuthRequest, error)
	SaveAuthCode(ctx context.Context, id, code, userAgentID string) error
	DeleteAuthRequest(ctx context.Context, id string) error

	CheckLoginName(ctx context.Context, id, loginName, userAgentID string) error
	CheckExternalUserLogin(ctx context.Context, authReqID, userAgentID string, user *model.ExternalUser, info *model.BrowserInfo) error
	SelectUser(ctx context.Context, id, userID, userAgentID string) error
	SelectExternalIDP(ctx context.Context, authReqID, idpConfigID, userAgentID string) error
	VerifyPassword(ctx context.Context, id, userID, password, userAgentID string, info *model.BrowserInfo) error

	VerifyMFAOTP(ctx context.Context, agentID, authRequestID, code, userAgentID string, info *model.BrowserInfo) error
	BeginMFAU2FLogin(ctx context.Context, userID, authRequestID, userAgentID string) (*user_model.WebAuthNLogin, error)
	VerifyMFAU2F(ctx context.Context, userID, authRequestID, userAgentID string, credentialData []byte, info *model.BrowserInfo) error
	BeginPasswordlessLogin(ctx context.Context, userID, authRequestID, userAgentID string) (*user_model.WebAuthNLogin, error)
	VerifyPasswordless(ctx context.Context, userID, authRequestID, userAgentID string, credentialData []byte, info *model.BrowserInfo) error

	LinkExternalUsers(ctx context.Context, authReqID, userAgentID string, info *model.BrowserInfo) error
	AutoRegisterExternalUser(ctx context.Context, user *user_model.User, externalIDP *user_model.ExternalIDP, member *org_model.OrgMember, authReqID, userAgentID, resourceOwner string, info *model.BrowserInfo) error
	ResetLinkingUsers(ctx context.Context, authReqID, userAgentID string) error

	GetOrgByPrimaryDomain(primaryDomain string) (*org_model.OrgView, error)
}
