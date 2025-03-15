package cache

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/cache"
	"github.com/zitadel/zitadel/backend/storage/cache/gomap"
)

type User struct {
	cache.Cache[string, *repository.User]
}

func NewUser() *User {
	return &User{
		Cache: gomap.New[string, *repository.User](),
	}
}

// ByID implements repository.UserRepository.
func (u *User) ByID(ctx context.Context, id string) (*repository.User, error) {
	user, _ := u.Get(id)
	return user, nil

}

func (u *User) Set(ctx context.Context, user *repository.User) (*repository.User, error) {
	u.set(user)
	return user, nil
}

func (u *User) set(user *repository.User) {
	u.Cache.Set(user.ID, user)
}
