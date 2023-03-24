package ldap

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
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
			s.Password, s.Provider.timeout)
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

func tryBind(
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
	conn, err := getConnection(server, startTLS, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := conn.Bind(bindDN, bindPassword); err != nil {
		return nil, err
	}

	return trySearchAndUserBind(
		conn,
		baseDN,
		attributes,
		objectClasses,
		userFilters,
		username,
		password,
		timeout,
	)
}

func getConnection(
	server string,
	startTLS bool,
	timeout time.Duration,
) (*ldap.Conn, error) {
	if timeout == 0 {
		timeout = ldap.DefaultTimeout
	}

	conn, err := ldap.DialURL(server, ldap.DialWithDialer(&net.Dialer{Timeout: timeout}))
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "ldaps" && startTLS {
		err = conn.StartTLS(&tls.Config{ServerName: u.Host})
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func trySearchAndUserBind(
	conn *ldap.Conn,
	baseDN string,
	attributes []string,
	objectClasses []string,
	userFilters []string,
	username string,
	password string,
	timeout time.Duration,
) (*ldap.Entry, error) {
	searchQuery := queriesAndToSearchQuery(
		objectClassesToSearchQuery(objectClasses),
		queriesOrToSearchQuery(
			userFiltersToSearchQuery(userFilters, username),
		),
	)

	// Search for user with the unique attribute for the userDN
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, int(timeout.Seconds()), false,
		searchQuery,
		attributes,
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	if len(sr.Entries) != 1 {
		return nil, ErrNoSingleUser
	}

	user := sr.Entries[0]
	// Bind as the user to verify their password
	if err = conn.Bind(user.DN, password); err != nil {
		return nil, ErrFailedLogin
	}
	return user, nil
}

func queriesAndToSearchQuery(queries ...string) string {
	if len(queries) == 0 {
		return ""
	}
	if len(queries) == 1 {
		return queries[0]
	}
	joinQueries := "(&"
	for _, s := range queries {
		joinQueries += s
	}
	return joinQueries + ")"
}

func queriesOrToSearchQuery(queries ...string) string {
	if len(queries) == 0 {
		return ""
	}
	if len(queries) == 1 {
		return queries[0]
	}
	joinQueries := "(|"
	for _, s := range queries {
		joinQueries += s
	}
	return joinQueries + ")"
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

func mapLDAPEntryToUser(
	user *ldap.Entry,
	idAttribute,
	firstNameAttribute,
	lastNameAttribute,
	displayNameAttribute,
	nickNameAttribute,
	preferredUsernameAttribute,
	emailAttribute,
	emailVerifiedAttribute,
	phoneAttribute,
	phoneVerifiedAttribute,
	preferredLanguageAttribute,
	avatarURLAttribute,
	profileAttribute string,
) (_ *User, err error) {
	var emailVerified bool
	if v := user.GetAttributeValue(emailVerifiedAttribute); v != "" {
		emailVerified, err = strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
	}
	var phoneVerified bool
	if v := user.GetAttributeValue(phoneVerifiedAttribute); v != "" {
		phoneVerified, err = strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
	}

	return NewUser(
		user.GetAttributeValue(idAttribute),
		user.GetAttributeValue(firstNameAttribute),
		user.GetAttributeValue(lastNameAttribute),
		user.GetAttributeValue(displayNameAttribute),
		user.GetAttributeValue(nickNameAttribute),
		user.GetAttributeValue(preferredUsernameAttribute),
		domain.EmailAddress(user.GetAttributeValue(emailAttribute)),
		emailVerified,
		domain.PhoneNumber(user.GetAttributeValue(phoneAttribute)),
		phoneVerified,
		language.Make(user.GetAttributeValue(preferredLanguageAttribute)),
		user.GetAttributeValue(avatarURLAttribute),
		user.GetAttributeValue(profileAttribute),
	), nil
}
