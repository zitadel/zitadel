package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userIdentityProviderLink struct{}

// CreatedAtColumn implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) CreatedAtColumn() database.Column {
	panic("unimplemented")
}

// IdentityProviderIDColumn implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) IdentityProviderIDColumn() database.Column {
	panic("unimplemented")
}

// IdentityProviderIDCondition implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) IdentityProviderIDCondition(op database.TextOperation, id string) database.Condition {
	panic("unimplemented")
}

// InstanceIDColumn implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) InstanceIDColumn() database.Column {
	panic("unimplemented")
}

// InstanceIDCondition implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) InstanceIDCondition(op database.TextOperation, id string) database.Condition {
	panic("unimplemented")
}

// PrimaryKeyCondition implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) PrimaryKeyCondition(instanceID string, idpID string, providedID string) database.Condition {
	panic("unimplemented")
}

// ProvidedIDColumn implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) ProvidedIDColumn() database.Column {
	panic("unimplemented")
}

// ProvidedIDCondition implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) ProvidedIDCondition(op database.TextOperation, id string) database.Condition {
	panic("unimplemented")
}

// ProvidedUsernameColumn implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) ProvidedUsernameColumn() database.Column {
	panic("unimplemented")
}

// ProvidedUsernameCondition implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) ProvidedUsernameCondition(op database.TextOperation, username string) database.Condition {
	panic("unimplemented")
}

// SetProvidedUserID implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) SetProvidedUserID(providedID string) database.Change {
	panic("unimplemented")
}

// SetProvidedUsername implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) SetProvidedUsername(username string) database.Change {
	panic("unimplemented")
}

// SetUpdatedAt implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) SetUpdatedAt(updatedAt time.Time) database.Change {
	panic("unimplemented")
}

// UpdatedAtColumn implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) UpdatedAtColumn() database.Column {
	panic("unimplemented")
}

// UserIDColumn implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) UserIDColumn() database.Column {
	panic("unimplemented")
}

// UserIDCondition implements domain.UserIdentityProviderLinkRepository.
func (u *userIdentityProviderLink) UserIDCondition(op database.TextOperation, id string) database.Condition {
	panic("unimplemented")
}

func UserIdentityProviderLinkRepository() domain.UserIdentityProviderLinkRepository {
	return &userIdentityProviderLink{}
}

// Create implements [domain.UserIdentityProviderLinkRepository].
func (u *userIdentityProviderLink) Create(ctx context.Context, client database.QueryExecutor, link *domain.UserIdentityProviderLink) error {
	panic("unimplemented")
}

// Delete implements [domain.UserIdentityProviderLinkRepository].
func (u *userIdentityProviderLink) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	panic("unimplemented")
}

// Get implements [domain.UserIdentityProviderLinkRepository].
func (u *userIdentityProviderLink) Get(ctx context.Context, client database.QueryExecutor, condition database.Condition, opts ...database.QueryOption) (*domain.UserIdentityProviderLink, error) {
	panic("unimplemented")
}

// List implements [domain.UserIdentityProviderLinkRepository].
func (u *userIdentityProviderLink) List(ctx context.Context, client database.QueryExecutor, condition database.Condition, opts ...database.QueryOption) ([]*domain.UserIdentityProviderLink, error) {
	panic("unimplemented")
}

// PrimaryKeyColumns implements [domain.UserIdentityProviderLinkRepository].
func (u *userIdentityProviderLink) PrimaryKeyColumns() []database.Column {
	panic("unimplemented")
}

// Update implements [domain.UserIdentityProviderLinkRepository].
func (u *userIdentityProviderLink) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	panic("unimplemented")
}
