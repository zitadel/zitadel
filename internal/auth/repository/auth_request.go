package repository

import (
	"context"

	"github.com/caos/zitadel/internal/auth_request/model"
)

type AuthRequestRepository interface {
	CreateAuthRequest(ctx context.Context, authRequest *model.AuthRequest) (*model.AuthRequest, error)
	AuthRequestByID(ctx context.Context, id string) (*model.AuthRequest, error)
}
