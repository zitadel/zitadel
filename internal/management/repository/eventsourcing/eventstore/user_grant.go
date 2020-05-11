package eventstore

import (
	"context"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	grant_event "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing"
)

type UserGrantRepo struct {
	UserGrantEvents *grant_event.UserGrantEventStore
}

func (repo *UserGrantRepo) UserGrantByID(ctx context.Context, grantID string) (*grant_model.UserGrant, error) {
	return repo.UserGrantEvents.UserGrantByID(ctx, grantID)
}

func (repo *UserGrantRepo) AddUserGrant(ctx context.Context, grant *grant_model.UserGrant) (*grant_model.UserGrant, error) {
	return repo.UserGrantEvents.AddUserGrant(ctx, grant)
}

func (repo *UserGrantRepo) ChangeUserGrant(ctx context.Context, grant *grant_model.UserGrant) (*grant_model.UserGrant, error) {
	return repo.UserGrantEvents.ChangeUserGrant(ctx, grant)
}

func (repo *UserGrantRepo) DeactivateUserGrant(ctx context.Context, grantID string) (*grant_model.UserGrant, error) {
	return repo.UserGrantEvents.DeactivateUserGrant(ctx, grantID)
}

func (repo *UserGrantRepo) ReactivateUserGrant(ctx context.Context, grantID string) (*grant_model.UserGrant, error) {
	return repo.UserGrantEvents.ReactivateUserGrant(ctx, grantID)
}

func (repo *UserGrantRepo) RemoveUserGrant(ctx context.Context, grantID string) error {
	return repo.UserGrantEvents.RemoveUserGrant(ctx, grantID)
}
