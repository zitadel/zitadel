package repository

import "context"

type Repository interface {
	Health(ctx context.Context) error
	AdministratorRepository
}
