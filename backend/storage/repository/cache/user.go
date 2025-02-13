package cache

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/cache"
	"github.com/zitadel/zitadel/backend/storage/repository"
)

type User struct {
	cache.Cache[string, *repository.User]

	next repository.UserRepository
}

// ByID implements repository.UserRepository.
func (u *User) ByID(ctx context.Context, id string) (*repository.User, error) {
	if user, ok := u.Get(id); ok {
		return user, nil
	}

	user, err := u.next.ByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.set(user)
	return user, nil
}

// Create implements repository.UserRepository.
func (u *User) Create(ctx context.Context, user *repository.User) error {
	err := u.next.Create(ctx, user)
	if err != nil {
		return err
	}
	u.set(user)
	return nil
}

var _ repository.UserRepository = (*User)(nil)

func (u *User) set(user *repository.User) {
	u.Cache.Set(user.ID, user)
}
