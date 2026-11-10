//go:build !integration

package ldap

import (
	"context"
	"errors"

	"github.com/go-ldap/ldap/v3"

	"github.com/zitadel/zitadel/internal/idp"
)

// FetchUser implements the [idp.Session] interface.
func (s *Session) FetchUser(_ context.Context) (_ idp.User, err error) {
	var user *ldap.Entry
	for _, server := range s.Provider.servers {
		user, err = tryBind(server,
			s.Provider.startTLS,
			s.Provider.bindDN,
			s.Provider.bindPassword,
			s.Provider.baseDN,
			s.Provider.getNecessaryAttributes(),
			s.Provider.userObjectClasses,
			s.Provider.userFilters,
			s.User,
			s.Password,
			s.Provider.timeout,
			s.Provider.rootCA)
		// If there were invalid credentials or multiple users with the credentials cancel process
		if err != nil && (errors.Is(err, ErrFailedLogin) || errors.Is(err, ErrNoSingleUser)) {
			return nil, err
		}
		// If a user bind was successful and user is filled continue with login, otherwise try next server
		if err == nil && user != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	s.Entry = user

	return mapLDAPEntryToUser(
		user,
		s.Provider.idAttribute,
		s.Provider.firstNameAttribute,
		s.Provider.lastNameAttribute,
		s.Provider.displayNameAttribute,
		s.Provider.nickNameAttribute,
		s.Provider.preferredUsernameAttribute,
		s.Provider.emailAttribute,
		s.Provider.emailVerifiedAttribute,
		s.Provider.phoneAttribute,
		s.Provider.phoneVerifiedAttribute,
		s.Provider.preferredLanguageAttribute,
		s.Provider.avatarURLAttribute,
		s.Provider.profileAttribute,
	)
}
