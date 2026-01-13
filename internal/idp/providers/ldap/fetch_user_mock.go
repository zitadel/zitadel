//go:build integration

package ldap

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
)

// FetchUser implements the [idp.Session] interface.
//
// This is a mock implementation of [idp.Session.FetchUser] that is used
// during integration tests and substitutes the real [idp.Session.FetchUser]
// implementation.
//
// This has been done to allow writing integration tests against a
// fake (external) LDAP server.
//
// The substitution/overloading is done through build tags.
func (s *Session) FetchUser(_ context.Context) (_ idp.User, err error) {
	usr := NewUser(
		gofakeit.UUID(),
		gofakeit.FirstName(),
		gofakeit.LastName(),
		gofakeit.Name(),
		gofakeit.Name(),
		gofakeit.Name(),
		domain.EmailAddress(gofakeit.Email()),
		true,
		domain.PhoneNumber(gofakeit.Phone()),
		true,
		language.English,
		gofakeit.URL(),
		gofakeit.UUID(),
	)

	s.Entry = &ldap.Entry{
		DN: s.Provider.baseDN,
		Attributes: []*ldap.EntryAttribute{
			{Name: "uid", Values: []string{gofakeit.UUID()}},
			{Name: "givenName", Values: []string{gofakeit.FirstName()}},
			{Name: "sn", Values: []string{gofakeit.LastName()}},
			{Name: "mail", Values: []string{gofakeit.Email()}},
			{Name: "telephoneNumber", Values: []string{gofakeit.Phone()}},
			{Name: "displayName", Values: []string{gofakeit.Name()}},
		},
	}
	return usr, nil
}
