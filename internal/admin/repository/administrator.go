package repository

import (
	"context"
)

type AdministratorRepository interface {
	ClearView(ctx context.Context, db, viewName string) error
}
