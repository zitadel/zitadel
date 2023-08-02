package idp

import (
	"context"
)

// Session is the minimal implementation for a session of a 3rd party authentication [Provider]
type Session interface {
	GetAuthURL() string
	FetchUser(ctx context.Context) (User, error)
}

type SessionSupportsMigration interface {
	RetrieveOldID() (oldID string, err error)
}
