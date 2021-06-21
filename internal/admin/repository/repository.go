package repository

import "context"

type Repository interface {
	Health(ctx context.Context) error
	OrgRepository
	IAMRepository
	AdministratorRepository
	FeaturesRepository
	UserRepository
}
