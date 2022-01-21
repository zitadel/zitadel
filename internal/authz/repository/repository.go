package repository

import (
	"context"
)

type Repository interface {
	Health(context.Context) error
	UserGrantRepository
}
