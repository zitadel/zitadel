package idp

import (
	"context"
)

type Session interface {
	GetAuthURL() string
	FetchUser(ctx context.Context) (User, error)
}
