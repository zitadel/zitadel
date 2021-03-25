package repository

import (
	"context"

	features_model "github.com/caos/zitadel/internal/features/model"
)

type FeaturesRepository interface {
	GetOrgFeatures(ctx context.Context, id string) (*features_model.FeaturesView, error)
}
