package idp

import (
	"context"
	"time"
)

// Session is the minimal implementation for a session of a 3rd party authentication [Provider]
type Session interface {
	GetAuth(ctx context.Context) (Auth, error)
	PersistentParameters() map[string]any
	FetchUser(ctx context.Context) (User, error)
	ExpiresAt() time.Time
}

type Auth interface {
	auth()
}

type RedirectAuth struct {
	RedirectURL string
}

func (r *RedirectAuth) auth() {}

type FormAuth struct {
	URL    string
	Fields map[string]string
}

func (f *FormAuth) auth() {}

// SessionSupportsMigration is an optional extension to the Session interface.
// It can be implemented to support migrating users, were the initial external id has changed because of a migration of the Provider type.
// E.g. when a user was linked on a generic OIDC provider and this provider has now been migrated to an AzureAD provider.
// In this case OIDC used the `sub` claim and Azure now uses the id of the user endpoint, which differ.
// The RetrievePreviousID will return the `sub` claim again, so that the user can be matched and safely migrated to the new id.
type SessionSupportsMigration interface {
	RetrievePreviousID() (previousID string, err error)
}

func Redirect(redirectURL string) (*RedirectAuth, error) {
	return &RedirectAuth{RedirectURL: redirectURL}, nil
}

func Form(url string, fields map[string]string) (*FormAuth, error) {
	return &FormAuth{
		URL:    url,
		Fields: fields,
	}, nil
}
