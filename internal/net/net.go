package net

import (
	builtin_net "net"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type IPLookupFunc func(string) ([]builtin_net.IP, error)

// HostnameToIPList converts the input URL to a list [net.IP].
//
// Returns an internal error if the lookup fails.
// Returns an invalid argument if the lookup function is nil
func HostnameToIPList(hostname string, lookupFunc IPLookupFunc) ([]builtin_net.IP, error) {
	ip := builtin_net.ParseIP(hostname)
	ips := []builtin_net.IP{ip}
	if ip != nil {
		return ips, nil
	}
	// if the hostname is a domain, we need to check resolve the ip(s), since it might be denied
	if lookupFunc == nil {
		return nil, zerrors.ThrowInvalidArgument(nil, "NET-naSn77", "lookup function must not be nil")
	}
	ips, err := lookupFunc(hostname)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "NET-4m9s2", "lookup failed")
	}
	return ips, nil
}
