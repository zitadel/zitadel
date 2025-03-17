package repository

import "github.com/zitadel/zitadel/backend/storage/cache"

type User struct {
	ID       string
	Username string
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
