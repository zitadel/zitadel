package repository

import (
	"context"

	features_model "github.com/caos/zitadel/internal/features/model"
)

type FeaturesRepository interface {
	GetDefaultFeatures(ctx context.Context) (*features_model.FeaturesView, error)
	GetOrgFeatures(ctx context.Context, id string) (*features_model.FeaturesView, error)
}
