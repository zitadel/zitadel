package ldap

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"net"
	"net/url"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/go-ldap/ldap/v3"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
)

var ErrNoSingleUser = errors.New("user does not exist or too many entries returned")
var ErrFailedLogin = errors.New("user failed to login")
var ErrUnableToAppendRootCA = errors.New("unable to append rootCA")

var _ idp.Session = (*Session)(nil)

type Session struct {
	Provider *Provider
	loginUrl string
	User     string
	Password string
	Entry    *ldap.Entry
}

func NewSession(provider *Provider, username, password string) *Session {
	return &Session{Provider: provider, User: username, Password: password}
}

// GetAuth implements the [idp.Session] interface.
func (s *Session) GetAuth(ctx context.Context) (idp.Auth, error) {
	return idp.Redirect(s.loginUrl)
}

// PersistentParameters implements the [idp.Session] interface.
func (s *Session) PersistentParameters() map[string]any {
	return nil
}

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

func (s *Session) ExpiresAt() time.Time {
	return time.Time{} // falls back to the default expiration time
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
	rootCA []byte,
) (*ldap.Entry, error) {
	conn, err := getConnection(server, startTLS, timeout, rootCA)
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
	rootCA []byte,
) (*ldap.Conn, error) {
	if timeout == 0 {
		timeout = ldap.DefaultTimeout
	}

	dialer := make([]ldap.DialOpt, 1, 2)
	dialer[0] = ldap.DialWithDialer(&net.Dialer{Timeout: timeout})

	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "ldaps" && len(rootCA) > 0 {
		rootCAs := x509.NewCertPool()
		if ok := rootCAs.AppendCertsFromPEM(rootCA); !ok {
			return nil, ErrUnableToAppendRootCA
		}

		dialer = append(dialer, ldap.DialWithTLSConfig(&tls.Config{
			RootCAs: rootCAs,
		}))
	}

	conn, err := ldap.DialURL(server, dialer...)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "ldap" && startTLS {
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
			userFiltersToSearchQuery(userFilters, username)...,
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
		logging.WithFields("entries", len(sr.Entries)).Info("ldap: no single user found")
		return nil, ErrNoSingleUser
	}

	user := sr.Entries[0]
	// Bind as the user to verify their password
	_, err = ldap.ParseDN(user.DN)
	if err != nil {
		logging.WithFields("userDN", user.DN).WithError(err).Info("ldap user parse DN failed")
		return nil, err
	}
	if err = conn.Bind(user.DN, password); err != nil {
		logging.WithFields("userDN", user.DN).WithError(err).Info("ldap user bind failed")
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

func userFiltersToSearchQuery(filters []string, username string) []string {
	searchQueries := make([]string, len(filters))
	for i, filter := range filters {
		searchQueries[i] = "(" + filter + "=" + username + ")"
	}
	return searchQueries
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
		getAttributeValue(user, idAttribute),
		getAttributeValue(user, firstNameAttribute),
		getAttributeValue(user, lastNameAttribute),
		getAttributeValue(user, displayNameAttribute),
		getAttributeValue(user, nickNameAttribute),
		getAttributeValue(user, preferredUsernameAttribute),
		domain.EmailAddress(user.GetAttributeValue(emailAttribute)),
		emailVerified,
		domain.PhoneNumber(user.GetAttributeValue(phoneAttribute)),
		phoneVerified,
		language.Make(user.GetAttributeValue(preferredLanguageAttribute)),
		user.GetAttributeValue(avatarURLAttribute),
		user.GetAttributeValue(profileAttribute),
	), nil
}

func getAttributeValue(user *ldap.Entry, attribute string) string {
	// return an empty string if no attribute is needed
	if attribute == "" {
		return ""
	}
	value := user.GetAttributeValue(attribute)
	if utf8.ValidString(value) {
		return value
	}
	return base64.StdEncoding.EncodeToString(user.GetRawAttributeValue(attribute))
}
