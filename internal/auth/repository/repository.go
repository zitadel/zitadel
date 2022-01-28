package repository

import (
	"context"
)

type Repository interface {
	Health(context.Context) error
	UserRepository
	AuthRequestRepository
	TokenRepository
	UserSessionRepository
	OrgRepository
	RefreshTokenRepository
}
