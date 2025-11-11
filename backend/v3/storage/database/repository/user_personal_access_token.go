package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userPersonalAccessToken struct{}

func UserPersonalAccessToken() domain.PersonalAccessTokenRepository {
	return new(userPersonalAccessToken)
}

// Create implements [domain.PersonalAccessTokenRepository].
func (u *userPersonalAccessToken) Create(ctx context.Context, client database.QueryExecutor, pat *domain.PersonalAccessToken) error {
	panic("unimplemented")
}

// Delete implements [domain.PersonalAccessTokenRepository].
func (u *userPersonalAccessToken) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	panic("unimplemented")
}

// Get implements [domain.PersonalAccessTokenRepository].
func (u *userPersonalAccessToken) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOpts) (*domain.PersonalAccessToken, error) {
	panic("unimplemented")
}

// InstanceIDCondition implements [domain.PersonalAccessTokenRepository].
func (u *userPersonalAccessToken) InstanceIDCondition(instanceID string) database.Condition {
	panic("unimplemented")
}

// List implements [domain.PersonalAccessTokenRepository].
func (u *userPersonalAccessToken) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOpts) ([]*domain.PersonalAccessToken, error) {
	panic("unimplemented")
}

// PrimaryKeyCondition implements [domain.PersonalAccessTokenRepository].
func (u *userPersonalAccessToken) PrimaryKeyCondition(instanceID string, tokenID string) database.Condition {
	panic("unimplemented")
}

// UserIDCondition implements [domain.PersonalAccessTokenRepository].
func (u *userPersonalAccessToken) UserIDCondition(userID string) database.Condition {
	panic("unimplemented")
}

func PersonalAccessTokenRepository() domain.PersonalAccessTokenRepository {
	return new(userPersonalAccessToken)
}
