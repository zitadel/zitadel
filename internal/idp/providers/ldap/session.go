package ldap

import (
	"context"
	"crypto/tls"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp"
)

var ErrNoSingleUser = errors.New("user does not exist or too many entries returned")
var ErrFailedLogin = errors.New("user failed to login")

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
	var user *ldap.Entry
	for _, server := range s.Provider.servers {
		userT, err := tryLogin(
			server,
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
		)
		if err != nil {
			return nil, err
		}
		user = userT
		break
	}

	var err error
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

func tryLogin(
	server string,
	startTLS bool,
	bindDN string,
	bindPassword string,
	baseDN string,
	attributes []string,
	objectClasses []string,
	userFilters []string,
	username string,
	password string,
	timeout time.Duration,
) (*ldap.Entry, error) {
	l, err := ldap.DialURL(server)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	if !strings.HasPrefix(server, "ldaps") && startTLS {
		parts := strings.Split(server, "://")
		parts = strings.Split(parts[1], ":")

		err = l.StartTLS(&tls.Config{ServerName: parts[0]})
		if err != nil {
			return nil, err
		}
	}
	// Bind as the admin to search for user
	err = l.Bind(bindDN, bindPassword)
	if err != nil {
		return nil, err
	}

	searchQuery := queriesAndToSearchQuery(objectClassesToSearchQuery(objectClasses), queriesOrToSearchQuery(userFiltersToSearchQuery(userFilters, username)))

	// Search for user with the unique attribute for the userDN
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, int(timeout.Seconds()), false,
		searchQuery,
		attributes,
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

	if err = l.Bind(user.DN, password); err != nil {
		return nil, ErrFailedLogin
	}
	return user, nil
}

func queriesAndToSearchQuery(queries ...string) string {
	joinQueries := []string{"(&"}
	joinQueries = append(joinQueries, queries...)
	joinQueries = append(joinQueries, ")")
	return strings.Join(joinQueries, "")
}

func queriesOrToSearchQuery(queries ...string) string {
	joinQueries := []string{"(|"}
	joinQueries = append(joinQueries, queries...)
	joinQueries = append(joinQueries, ")")
	return strings.Join(joinQueries, "")
}

func objectClassesToSearchQuery(classes []string) string {
	searchQuery := ""
	for _, class := range classes {
		searchQuery += "(objectClass=" + class + ")"
	}
	return searchQuery
}

func userFiltersToSearchQuery(filters []string, username string) string {
	searchQuery := ""
	for _, filter := range filters {
		searchQuery += "(" + filter + "=" + ldap.EscapeFilter(username) + ")"
	}
	return searchQuery
}
