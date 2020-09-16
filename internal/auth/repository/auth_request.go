package repository

import (
	"context"
	org_model "github.com/caos/zitadel/internal/org/model"
	user_model "github.com/caos/zitadel/internal/user/model"

	"github.com/caos/zitadel/internal/auth_request/model"
)

type AuthRequestRepository interface {
	CreateAuthRequest(ctx context.Context, request *model.AuthRequest) (*model.AuthRequest, error)
	AuthRequestByID(ctx context.Context, id, userAgentID string) (*model.AuthRequest, error)
	AuthRequestByIDCheckLoggedIn(ctx context.Context, id, userAgentID string) (*model.AuthRequest, error)
	AuthRequestByCode(ctx context.Context, code string) (*model.AuthRequest, error)
	SaveAuthCode(ctx context.Context, id, code, userAgentID string) error
	DeleteAuthRequest(ctx context.Context, id string) error
	CheckLoginName(ctx context.Context, id, loginName, userAgentID string) error
	CheckExternalUserLogin(ctx context.Context, authReqID, userAgentID string, user *model.ExternalUser) error
	SelectUser(ctx context.Context, id, userID, userAgentID string) error
	SelectExternalIDP(ctx context.Context, authReqID, idpConfigID, userAgentID string) error
	VerifyPassword(ctx context.Context, id, userID, password, userAgentID string, info *model.BrowserInfo) error
	VerifyMfaOTP(ctx context.Context, agentID, authRequestID, code, userAgentID string, info *model.BrowserInfo) error
	AddUserExternalIDPs(ctx context.Context, userID string, linkingUsers []*model.ExternalUser) error
	LinkExternalUsers(ctx context.Context, authReqID, userAgentID string) error
	AutoRegisterExternalUser(ctx context.Context, user *user_model.User, externalIDP *user_model.ExternalIDP, member *org_model.OrgMember, authReqID, userAgentID, resourceOwner string) error
}
