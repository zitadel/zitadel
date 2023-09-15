package idp

import (
	"context"
)

// Session is the minimal implementation for a session of a 3rd party authentication [Provider]
type Session interface {
	GetAuth(ctx context.Context) (content string, redirect bool)
	FetchUser(ctx context.Context) (User, error)
}

// SessionSupportsMigration is an optional extension to the Session interface.
// It can be implemented to support migrating users, were the initial external id has changed because of a migration of the Provider type.
// E.g. when a user was linked on a generic OIDC provider and this provider has now been migrated to an AzureAD provider.
// In this case OIDC used the `sub` claim and Azure now uses the id of the user endpoint, which differ.
// The RetrievePreviousID will return the `sub` claim again, so that the user can be matched and safely migrated to the new id.
type SessionSupportsMigration interface {
	RetrievePreviousID() (previousID string, err error)
}

func Redirect(redirectURL string) (string, bool) {
	return redirectURL, true
}

func Form(html string) (string, bool) {
	return html, false
}
