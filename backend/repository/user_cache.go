package repository

import (
	"context"
	"log"

	"github.com/zitadel/zitadel/backend/storage/cache"
)

type UserCache struct {
	cache.Cache[UserIndex, string, *User]
}

type UserIndex uint8

var UserIndices = []UserIndex{
	UserByIDIndex,
	UserByUsernameIndex,
}

const (
	UserByIDIndex UserIndex = iota
	UserByUsernameIndex
)

var _ cache.Entry[UserIndex, string] = (*User)(nil)

// Keys implements [cache.Entry].
func (u *User) Keys(index UserIndex) (key []string) {
	switch index {
	case UserByIDIndex:
		return []string{u.ID}
	case UserByUsernameIndex:
		return []string{u.Username}
	}
	return nil
}

func NewUserCache(c cache.Cache[UserIndex, string, *User]) *UserCache {
	return &UserCache{c}
}

func (c *UserCache) ByID(ctx context.Context, id string) *User {
	log.Println("cached.user.byID")
	user, _ := c.Cache.Get(ctx, UserByIDIndex, id)
	return user
}

func (c *UserCache) Set(ctx context.Context, user *User) {
	log.Println("cached.user.set")
	c.Cache.Set(ctx, user)
}
