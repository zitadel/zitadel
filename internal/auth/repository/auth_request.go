package repository

import (
	"context"

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
	SelectUser(ctx context.Context, id, userID, userAgentID string) error
	SelectExternalIDP(ctx context.Context, authReqID, idpConfigID, userAgentID string) error
	VerifyPassword(ctx context.Context, id, userID, password, userAgentID string, info *model.BrowserInfo) error
	VerifyMfaOTP(ctx context.Context, agentID, authRequestID, code, userAgentID string, info *model.BrowserInfo) error
}
