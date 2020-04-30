package repository

import (
	"context"

	auth_request "github.com/caos/zitadel/internal/auth_request/repository"
)

type Repository interface {
	Health(context.Context) error
	UserRepository
	auth_request.Repository
}
