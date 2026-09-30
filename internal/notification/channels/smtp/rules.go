package smtp

import (
	"fmt"
	"net"
	"slices"
	"strings"
	"text/template"

	"github.com/zitadel/zitadel/internal/notification/messages"
)

// RuleConfig is an operator defined rule, which is applied to the emails sent through a matching SMTP provider.
// Rules are part of the runtime configuration and can therefore not be changed through the API.
type RuleConfig struct {
	Match RuleMatch
	// DisableCustomHTML removes HTML from custom message texts.
	// Simple formatting (line breaks, bold, italic, underline, paragraphs) is kept,
	// all other elements like links, images or styles are removed, their text is kept.
	DisableCustomHTML bool
	// Headers are added to every email sent through the matching provider.
	Headers []RuleHeader
}

// RuleMatch defines the criteria a SMTP provider must fulfill for the rule to be applied.
// Every defined criterion must match. If none is defined, the rule matches all SMTP providers.
type RuleMatch struct {
	// Hosts of the SMTP provider. An entry without port matches any port.
	Hosts []string
	// SenderDomains are the domains of the sender address.
	SenderDomains []string
}

// RuleHeader is a list entry instead of a map, because keys of maps
// are lower-cased when read from the runtime configuration.
type RuleHeader struct {
	Name string
	// Value is a text/template, [RuleData] is passed as data.
	Value string
}

// RuleData is passed to the templates of the header values.
// It intentionally does not contain any personal data.
type RuleData struct {
	InstanceID string
	OrgID      string
}

// Rule is the result of the first matching [RuleConfig].
// The zero value is returned if no rule matches.
type Rule struct {
	DisableCustomHTML bool
	Headers           map[string]string
}

// Rules are the compiled [RuleConfig]s in the order they were defined.
type Rules []*compiledRule

type compiledRule struct {
	hosts             []string
	senderDomains     []string
	disableCustomHTML bool
	headers           []*compiledHeader
}

type compiledHeader struct {
	name  string
	value *template.Template
}

// CompileRules validates the rules and parses the templates of the header values.
// It is intended to be called on startup to fail fast on an invalid configuration.
func CompileRules(configs []RuleConfig) (Rules, error) {
	rules := make(Rules, len(configs))
	for i, config := range configs {
		headers, err := compileHeaders(config.Headers)
		if err != nil {
			return nil, fmt.Errorf("smtp rule %d: %w", i, err)
		}
		rules[i] = &compiledRule{
			hosts:             normalize(config.Match.Hosts),
			senderDomains:     normalize(config.Match.SenderDomains),
			disableCustomHTML: config.DisableCustomHTML,
			headers:           headers,
		}
	}
	return rules, nil
}

func compileHeaders(configs []RuleHeader) ([]*compiledHeader, error) {
	headers := make([]*compiledHeader, len(configs))
	for i, config := range configs {
		if !messages.IsValidEmailHeaderName(config.Name) {
			return nil, fmt.Errorf("header %q: invalid name", config.Name)
		}
		if messages.IsReservedEmailHeader(config.Name) {
			return nil, fmt.Errorf("header %q: reserved name", config.Name)
		}
		if !messages.IsValidEmailHeaderValue(config.Value) {
			return nil, fmt.Errorf("header %q: value must not contain line breaks", config.Name)
		}
		value, err := template.New(config.Name).Parse(config.Value)
		if err != nil {
			return nil, fmt.Errorf("header %q: %w", config.Name, err)
		}
		// fields unknown to RuleData are only detected on execution
		if err = value.Execute(new(strings.Builder), RuleData{}); err != nil {
			return nil, fmt.Errorf("header %q: %w", config.Name, err)
		}
		headers[i] = &compiledHeader{name: config.Name, value: value}
	}
	return headers, nil
}

// Match returns the [Rule] of the first rule matching the SMTP provider.
// If no rule matches, the zero value is returned.
func (r Rules) Match(config *Config, data RuleData) (Rule, error) {
	if config == nil {
		return Rule{}, nil
	}
	for _, rule := range r {
		if !rule.matches(config) {
			continue
		}
		headers, err := rule.renderHeaders(data)
		if err != nil {
			return Rule{}, err
		}
		return Rule{
			DisableCustomHTML: rule.disableCustomHTML,
			Headers:           headers,
		}, nil
	}
	return Rule{}, nil
}

func (r *compiledRule) matches(config *Config) bool {
	return matchesHost(r.hosts, config.SMTP.Host) &&
		matchesSenderDomain(r.senderDomains, config.From)
}

func (r *compiledRule) renderHeaders(data RuleData) (map[string]string, error) {
	if len(r.headers) == 0 {
		return nil, nil
	}
	headers := make(map[string]string, len(r.headers))
	for _, header := range r.headers {
		var value strings.Builder
		if err := header.value.Execute(&value, data); err != nil {
			return nil, fmt.Errorf("smtp rule header %q: %w", header.name, err)
		}
		headers[header.name] = value.String()
	}
	return headers, nil
}

func matchesHost(hosts []string, hostAndPort string) bool {
	if len(hosts) == 0 {
		return true
	}
	hostAndPort = strings.ToLower(strings.TrimSpace(hostAndPort))
	host, _, err := net.SplitHostPort(hostAndPort)
	if err != nil {
		// no port defined
		host = hostAndPort
	}
	return slices.Contains(hosts, hostAndPort) || slices.Contains(hosts, host)
}

func matchesSenderDomain(domains []string, sender string) bool {
	if len(domains) == 0 {
		return true
	}
	index := strings.LastIndex(sender, "@")
	if index < 0 {
		return false
	}
	return slices.Contains(domains, strings.ToLower(strings.TrimSpace(sender[index+1:])))
}

func normalize(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.ToLower(strings.TrimSpace(value)); value != "" {
			normalized = append(normalized, value)
		}
	}
	return normalized
}
