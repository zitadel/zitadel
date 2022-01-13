package repository

import (
	"context"

	"golang.org/x/text/language"
)

type IAMRepository interface {
	Languages(ctx context.Context) ([]language.Tag, error)

	GetIAMMemberRoles() []string
}
