package denylist

import (
	"net"
	"net/url"

	internal_net "github.com/zitadel/zitadel/internal/net"
)

var _ AddressChecker = (*HostChecker)(nil)

type HostChecker struct {
	Net    *net.IPNet
	IP     net.IP
	Domain string
}

func NewHostChecker(entry string) (AddressChecker, error) {
	if entry == "" {
		return nil, nil
	}
	_, network, err := net.ParseCIDR(entry)
	if err == nil {
		return &HostChecker{Net: network}, nil
	}
	if ip := net.ParseIP(entry); ip != nil {
		return &HostChecker{IP: ip}, nil
	}
	return &HostChecker{Domain: entry}, nil
}

func (c *HostChecker) IsDenied(ips []net.IP, address string) error {
	// if the address matches the domain, no additional checks as needed
	if c.Domain == address {
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

// IsHostBlocked checks address against denyList. If a match is found, an [AddressDeniedError] will be returned
// Takes an input lookupFunc to convert address to a list of IP addresses.
// If lookupFunc is nil, it defaults to [net.LookupIP]
func IsHostBlocked(denyList []AddressChecker, address *url.URL, lookupFunc internal_net.IPLookupFunc) error {
	if lookupFunc == nil {
		lookupFunc = net.LookupIP
	}

	host := address.Hostname()
	ips, err := internal_net.HostnameToIPList(host, lookupFunc)
	if err != nil {
		return err
	}

	for _, denied := range denyList {
		if err := denied.IsDenied(ips, host); err != nil {
			return err
		}
	}
	return nil
}
