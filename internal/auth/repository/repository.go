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
	KeyRepository
	UserSessionRepository
	UserGrantRepository
	OrgRepository
	IAMRepository
	FeaturesRepository
	RefreshTokenRepository
}
