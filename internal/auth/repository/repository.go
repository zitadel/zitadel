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
	ProjectRepository
	KeyRepository
	UserSessionRepository
	UserGrantRepository
	OrgRepository
	IAMRepository
	FeaturesRepository
}
