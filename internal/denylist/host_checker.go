package denylist

import (
	"net"
	"net/url"
	"strings"

	internal_net "github.com/zitadel/zitadel/internal/net"
)

var _ AddressChecker = (*HostChecker)(nil)

type HostChecker struct {
	Net    *net.IPNet
	IP     net.IP
	Domain string
}

// NewHostChecker returns an [AddressChecker].
//
// If the input string is neither a network, nor an IP address, it is considered a domain and saved as such
func NewHostChecker(entry string) AddressChecker {
	if entry == "" {
		return nil
	}
	_, network, err := net.ParseCIDR(entry)
	if err == nil {
		return &HostChecker{Net: network}
	}
	if ip := net.ParseIP(entry); ip != nil {
		return &HostChecker{IP: ip}
	}
	return &HostChecker{Domain: entry}
}

// IsDenied checks if the address or one of the IPs is blocked by the [HostChecker].
// The address is checked using string comparison.
// If either is blocked, an [AddressDeniedError] will be returned.
// If either is not blocked, it will return nil.
func (c *HostChecker) IsDenied(ips []net.IP, address string) error {
	// if the address matches the domain, no additional checks as needed
	if c.Domain != "" && strings.EqualFold(c.Domain, address) {
		return NewAddressDeniedError(c.Domain)
	}
	// otherwise we need to check on ips (incl. the resolved ips of the host)
	for _, ip := range ips {
		if c.Net != nil && c.Net.Contains(ip) {
			return NewAddressDeniedError(c.Net.String())
		}
		if c.IP != nil && c.IP.Equal(ip) {
			return NewAddressDeniedError(c.IP.String())
		}
	}
	return nil
}

// IsHostBlocked checks address against denyList. If a match is found, an [AddressDeniedError] will be returned.
// If no match is found, it will return nil.
func IsHostBlocked(denyList []AddressChecker, host string, ips ...net.IP) error {
	for _, denied := range denyList {
		if err := denied.IsDenied(ips, host); err != nil {
			return err
		}
	}
	return nil
}

// IsHostnameBlocked checks a hostname against the denyList. If a match is found, an [AddressDeniedError] will be returned
// It takes an input lookupFunc to convert the hostname to a list of IP addresses.
// If lookupFunc is nil, it defaults to [net.LookupIP]
func IsHostnameBlocked(denyList []AddressChecker, hostname string, lookupFunc internal_net.IPLookupFunc) error {
	if lookupFunc == nil {
		lookupFunc = net.LookupIP
	}

	ips, err := internal_net.HostnameToIPList(hostname, lookupFunc)
	if err != nil {
		return err
	}

	return IsHostBlocked(denyList, hostname, ips...)
}

// IsURLBlocked checks a URL against denyList. If a match is found, an [AddressDeniedError] will be returned
// It takes an input lookupFunc to convert the URL's hostname to a list of IP addresses.
// If lookupFunc is nil, it defaults to [net.LookupIP]
func IsURLBlocked(denyList []AddressChecker, address *url.URL, lookupFunc internal_net.IPLookupFunc) error {
	return IsHostnameBlocked(denyList, address.Hostname(), lookupFunc)
}
