package repository

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IamRepository interface {
	IamByID(ctx context.Context, id string) (*iam_model.Iam, error)
}
