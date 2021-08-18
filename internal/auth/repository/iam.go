package repository

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository interface {
	Languages(ctx context.Context) ([]language.Tag, error)
	GetIAM(ctx context.Context) (*model.IAM, error)
}
