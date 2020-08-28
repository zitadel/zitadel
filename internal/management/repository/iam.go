package repository

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IamRepository interface {
	IAMByID(ctx context.Context, id string) (*iam_model.IAM, error)
}
