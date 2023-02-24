package ldap

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp"
)

var ErrNoSingleUser = errors.New("user does not exist or too many entries returned")

var _ idp.Session = (*Session)(nil)

type Session struct {
	Provider *Provider
	loginUrl string
	User     string
	Password string
}

func (s *Session) GetAuthURL() string {
	return s.loginUrl
}
func (s *Session) FetchUser(_ context.Context) (idp.User, error) {
	l, err := ldap.DialURL("ldap://" + s.Provider.host + ":" + s.Provider.port)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	if s.Provider.tls {
		err = l.StartTLS(&tls.Config{ServerName: s.Provider.host})
		if err != nil {
			return nil, err
		}
	}

	// Bind as the admin to search for user
	err = l.Bind("cn="+s.Provider.admin+","+s.Provider.baseDN, s.Provider.password)
	if err != nil {
		return nil, err
	}

	// Search for user with the unique attribute for the userDN
	searchRequest := ldap.NewSearchRequest(
		s.Provider.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass="+s.Provider.userObjectClass+")("+s.Provider.userUniqueAttribute+"=%s))", ldap.EscapeFilter(s.User)),
		[]string{"dn",
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
		},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	if len(sr.Entries) != 1 {
		return nil, ErrNoSingleUser
	}

	user := sr.Entries[0]
	// Bind as the user to verify their password
	err = l.Bind(user.DN, s.Password)
	if err != nil {
		return nil, err
	}
	var emailVerified bool
	if v := user.GetAttributeValue(s.Provider.emailVerifiedAttribute); v != "" {
		emailVerified, err = strconv.ParseBool(user.GetAttributeValue(s.Provider.emailVerifiedAttribute))
		if err != nil {
			return nil, err
		}
	}
	var phoneVerified bool
	if v := user.GetAttributeValue(s.Provider.phoneVerifiedAttribute); v != "" {
		phoneVerified, err = strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
	}

	return NewUser(
		user.GetAttributeValue(s.Provider.idAttribute),
		user.GetAttributeValue(s.Provider.firstNameAttribute),
		user.GetAttributeValue(s.Provider.lastNameAttribute),
		user.GetAttributeValue(s.Provider.displayNameAttribute),
		user.GetAttributeValue(s.Provider.nickNameAttribute),
		user.GetAttributeValue(s.Provider.preferredUsernameAttribute),
		user.GetAttributeValue(s.Provider.emailAttribute),
		emailVerified,
		user.GetAttributeValue(s.Provider.phoneAttribute),
		phoneVerified,
		language.Make(user.GetAttributeValue(s.Provider.preferredLanguageAttribute)),
		user.GetAttributeValue(s.Provider.avatarURLAttribute),
		user.GetAttributeValue(s.Provider.profileAttribute),
	), nil
}
