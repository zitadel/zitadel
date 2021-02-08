package repository

import (
	"context"
	"github.com/caos/zitadel/internal/v2/domain"
)

type AuthRequestCache interface {
	Health(ctx context.Context) error

	GetAuthRequestByID(ctx context.Context, id string) (*domain.AuthRequest, error)
	GetAuthRequestByCode(ctx context.Context, code string) (*domain.AuthRequest, error)
	SaveAuthRequest(ctx context.Context, request *domain.AuthRequest) error
	UpdateAuthRequest(ctx context.Context, request *domain.AuthRequest) error
	DeleteAuthRequest(ctx context.Context, id string) error
}
