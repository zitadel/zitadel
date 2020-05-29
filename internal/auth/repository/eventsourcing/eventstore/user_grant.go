package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
)

type UserGrantRepo struct {
	View *view.View
}

func (repo *UserGrantRepo) SearchUserGrants(ctx context.Context, request *grant_model.UserGrantSearchRequest) ([]*grant_model.UserGrantSearchResponse, error) {
	return nil, nil
}
