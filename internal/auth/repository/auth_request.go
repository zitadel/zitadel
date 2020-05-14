package repository

import (
	"context"

	"github.com/caos/zitadel/internal/auth_request/model"
)

type AuthRequestRepository interface {
	CreateAuthRequest(ctx context.Context, request *model.AuthRequest) (*model.AuthRequest, error)
	AuthRequestByID(ctx context.Context, id string) (*model.AuthRequest, error)
	CheckUsername(ctx context.Context, id, username string) error
	VerifyPassword(ctx context.Context, id, userID, password string, info *model.BrowserInfo) error
	VerifyMfaOTP(ctx context.Context, agentID, authRequestID string, code string, info *model.BrowserInfo) error
}
