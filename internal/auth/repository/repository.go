package repository

import (
	"context"
)

type Repository interface {
	Health(context.Context) error
	UserRepository
	AuthRequestRepository
	TokenRepository
	ApplicationRepository
	UserSessionRepository
	UserGrantRepository
	OrgRepository
	RefreshTokenRepository
}
