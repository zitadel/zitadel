package repository

import (
	"context"

	"github.com/caos/zitadel/internal/auth_request/model"
)

type Repository interface {
	Health(ctx context.Context) error

	GetAuthRequestByID(ctx context.Context, id string) (*model.AuthRequest, error)
	SaveAuthRequest(ctx context.Context, id string) (*model.AuthRequest, error)
}
