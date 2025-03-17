package cached

import (
	"context"
	"log"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/cache"
)

type User struct {
	cache.Cache[repository.UserIndex, string, *repository.User]
}

func NewUser(c cache.Cache[repository.UserIndex, string, *repository.User]) *User {
	return &User{c}
}

func (i *User) ByID(ctx context.Context, id string) *repository.User {
	log.Println("cached.user.byid")
	user, _ := i.Cache.Get(ctx, repository.UserByIDIndex, id)
	return user
}

func (i *User) Set(ctx context.Context, user *repository.User) {
	log.Println("cached.user.set")
	i.Cache.Set(ctx, user)
}
