package smtp

import (
	"errors"
	"fmt"
	"maps"
	"net"
	"net/textproto"
	"slices"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/notification/messages"
)

// RuleConfig is an operator defined rule, which is applied to the emails sent through a matching SMTP provider.
// Rules are part of the runtime configuration and can therefore not be changed through the API.
type RuleConfig struct {
	Match RuleMatch
	// RestrictCustomHTML removes HTML from custom message texts.
	// Simple formatting (line breaks, bold, italic, underline, paragraphs) is kept,
	// all other elements like links, images or styles are removed, their text is kept.
	RestrictCustomHTML bool
	// Headers are added to every email sent through the matching provider.
	// The names are canonicalized (e.g. x-instance-id becomes X-Instance-Id),
	// the values can contain text and the placeholders {{.InstanceID}} and {{.OrgID}}.
	Headers map[string]string
}

// RuleMatch defines the criteria a SMTP provider must fulfill for the rule to be applied.
// Every defined criterion must match. If none is defined, the rule matches all SMTP providers.
type RuleMatch struct {
	// Hosts of the SMTP provider. An entry without port matches any port.
	Hosts []string
	// Users are the usernames of the SMTP authentication, e.g. the token of the provider.
	// They identify the account of the provider, as opposed to the sender address, which can be changed in the provider config.
	Users []string
	// SenderDomains are the domains of the sender address.
	SenderDomains []string
}

// RuleData provides the values of the placeholders of the header values.
// It intentionally does not contain any personal data.
type RuleData struct {
	InstanceID string
	OrgID      string
}

// headerPlaceholders are the only placeholders supported in header values.
// Anything else is rejected on startup, because whether it fails could depend on the data
// and would only be detected when an email is sent.
var headerPlaceholders = map[string]func(RuleData) string{
	"{{.InstanceID}}": func(data RuleData) string { return data.InstanceID },
	"{{.OrgID}}":      func(data RuleData) string { return data.OrgID },
}

// Rule is the result of the first matching [RuleConfig].
// The zero value is returned if no rule matches.
type Rule struct {
	RestrictCustomHTML bool
	Headers            map[string]string
}

// Rules are the compiled [RuleConfig]s in the order they were defined.
type Rules []*compiledRule

type compiledRule struct {
	hosts              []hostPort
	users              []string
	senderDomains      []string
	restrictCustomHTML bool
	// headers with canonical names
	headers map[string]string
}

type hostPort struct {
	host string
	// port is empty if any port matches
	port string
}

// CompileRules validates the rules.
// It is intended to be called on startup to fail fast on an invalid configuration.
func CompileRules(configs []RuleConfig) (Rules, error) {
	rules := make(Rules, len(configs))
	for i, config := range configs {
		headers, err := compileHeaders(config.Headers)
		if err != nil {
			return nil, fmt.Errorf("smtp rule %d: %w", i, err)
		}
		rules[i] = &compiledRule{
			hosts:              normalizeHosts(config.Match.Hosts),
			users:              normalizeUsers(config.Match.Users),
			senderDomains:      normalizeDomains(config.Match.SenderDomains),
			restrictCustomHTML: config.RestrictCustomHTML,
			headers:            headers,
		}
	}
	return rules, nil
}

var errUnsupportedHeaderValue = errors.New("value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}")

// compileHeaders validates the headers and canonicalizes their names.
// Names are matched case-insensitively, so they must not only differ in case.
func compileHeaders(configs map[string]string) (map[string]string, error) {
	headers := make(map[string]string, len(configs))
	// iterate in a defined order to report the same error on every start
	for _, name := range slices.Sorted(maps.Keys(configs)) {
		value := configs[name]
		if !messages.IsValidEmailHeaderName(name) {
			return nil, fmt.Errorf("header %q: invalid name", name)
		}
		if messages.IsReservedEmailHeader(name) {
			return nil, fmt.Errorf("header %q: reserved name", name)
		}
		if !messages.IsValidEmailHeaderValue(value) {
			return nil, fmt.Errorf("header %q: value must not contain line breaks", name)
		}
		if !isValidHeaderValue(value) {
			return nil, fmt.Errorf("header %q: %w", name, errUnsupportedHeaderValue)
		}
		canonical := textproto.CanonicalMIMEHeaderKey(name)
		if _, ok := headers[canonical]; ok {
			return nil, fmt.Errorf("header %q: duplicate name", name)
		}
		headers[canonical] = value
	}
	return headers, nil
}

// isValidHeaderValue reports whether the value consists of text and the supported placeholders only.
func isValidHeaderValue(value string) bool {
	for placeholder := range headerPlaceholders {
		value = strings.ReplaceAll(value, placeholder, "")
	}
	return !strings.Contains(value, "{{") && !strings.Contains(value, "}}")
}

// Match returns the [Rule] of the first rule matching the SMTP provider.
// If no rule matches, the zero value is returned.
func (r Rules) Match(config *Config, data RuleData) Rule {
	if config == nil {
		return Rule{}
	}
	for _, rule := range r {
		if !rule.matches(config) {
			continue
		}
		return Rule{
			RestrictCustomHTML: rule.restrictCustomHTML,
			Headers:            rule.renderHeaders(data),
		}
	}
	return Rule{}
}

func (r *compiledRule) matches(config *Config) bool {
	return matchesHost(r.hosts, config.SMTP.Host) &&
		matchesUser(r.users, config.SMTP) &&
		matchesSenderDomain(r.senderDomains, config.From)
}

func (r *compiledRule) renderHeaders(data RuleData) map[string]string {
	if len(r.headers) == 0 {
		return nil
	}
	headers := make(map[string]string, len(r.headers))
	for name, value := range r.headers {
		for placeholder, render := range headerPlaceholders {
			value = strings.ReplaceAll(value, placeholder, render(data))
		}
		headers[name] = value
	}
	return headers
}

func matchesHost(hosts []hostPort, hostAndPort string) bool {
	if len(hosts) == 0 {
		return true
	}
	configured := parseHostPort(hostAndPort)
	return slices.ContainsFunc(hosts, func(h hostPort) bool {
		return h.host == configured.host && (h.port == "" || h.port == configured.port)
	})
}

func matchesUser(users []string, config SMTP) bool {
	if len(users) == 0 {
		return true
	}
	var user string
	switch {
	case config.PlainAuth != nil:
		user = config.PlainAuth.User
	case config.XOAuth2Auth != nil:
		user = config.XOAuth2Auth.User
	}
	return user != "" && slices.Contains(users, strings.TrimSpace(user))
}

func matchesSenderDomain(domains []string, sender string) bool {
	if len(domains) == 0 {
		return true
	}
	senderDomain := domain.EmailAddress(sender).Domain()
	return senderDomain != "" && slices.Contains(domains, strings.ToLower(senderDomain))
}

// parseHostPort splits a host with an optional port. Brackets of IPv6 addresses are removed.
func parseHostPort(value string) hostPort {
	value = strings.ToLower(strings.TrimSpace(value))
	if host, port, err := net.SplitHostPort(value); err == nil {
		return hostPort{host: host, port: port}
	}
	return hostPort{host: strings.TrimSuffix(strings.TrimPrefix(value, "["), "]")}
}

func normalizeHosts(values []string) []hostPort {
	hosts := make([]hostPort, 0, len(values))
	for _, value := range values {
		if host := parseHostPort(value); host.host != "" {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

func normalizeUsers(values []string) []string {
	users := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			users = append(users, value)
		}
	}
	return users
}

func normalizeDomains(values []string) []string {
	domains := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.ToLower(strings.TrimSpace(value)); value != "" {
			domains = append(domains, value)
		}
	}
	return domains
}
