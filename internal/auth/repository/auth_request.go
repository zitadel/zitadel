package repository

import (
	"context"

	"github.com/caos/zitadel/internal/auth_request/model"
)

type AuthRequestRepository interface {
	CreateAuthRequest(ctx context.Context, request *model.AuthRequest) (*model.AuthRequest, error)
	AuthRequestByID(ctx context.Context, id string) (*model.AuthRequest, error)
	CheckUsername(ctx context.Context, id, username string) (*model.AuthRequest, error)
	VerifyPassword(ctx context.Context, id, userID, password string, info *model.BrowserInfo) (*model.AuthRequest, error)
	RequestPasswordReset(ctx context.Context, id, userID string, info *model.BrowserInfo) (*model.AuthRequest, error) //?
	SkipMfaInit(ctx context.Context, id, userID string) (*model.AuthRequest, error)
	AddMfa(ctx context.Context, agentID, authRequestID string, mfa interface{}, info *model.BrowserInfo) (*model.AuthRequest, error)
	VerifyMfa(ctx context.Context, agentID, authRequestID string, mfa interface{}, info *model.BrowserInfo) (*model.AuthRequest, error)
}
