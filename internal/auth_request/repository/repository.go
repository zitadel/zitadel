package repository

import (
	"context"

	"github.com/caos/zitadel/internal/auth_request/model"
)

type AuthRequestCache interface {
	Health(ctx context.Context) error

	GetAuthRequestByID(ctx context.Context, id string) (*model.AuthRequest, error)
	GetAuthRequestByCode(ctx context.Context, code string) (*model.AuthRequest, error)
	SaveAuthRequest(ctx context.Context, request *model.AuthRequest) error
	UpdateAuthRequest(ctx context.Context, request *model.AuthRequest) error
	DeleteAuthRequest(ctx context.Context, id string) error
}
