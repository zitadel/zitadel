package http

import (
	"context"
	"fmt"
	"net"
	"net/url"
)

type DomainCtx struct {
	InstanceHost string
	PublicHost   string
	Protocol     string
}

func NewDomainCtx(instanceHostname, publicHostname, protocol string) *DomainCtx {
	return &DomainCtx{
		InstanceHost: instanceHostname,
		PublicHost:   publicHostname,
		Protocol:     protocol,
	}
}

func NewDomainCtxFromOrigin(origin *url.URL) *DomainCtx {
	return &DomainCtx{
		InstanceHost: origin.Host,
		PublicHost:   origin.Host,
		Protocol:     origin.Scheme,
	}
}

// InstanceDomain returns the hostname for which the request was handled.
func (r *DomainCtx) InstanceDomain() string {
	return hostnameFromHost(r.InstanceHost)
}

func hostnameFromHost(host string) string {
	hostname, _, err := net.SplitHostPort(host)
	if err != nil {
		return host
	}
	return hostname
}

// RequestedHost returns the host (hostname[:port]) for which the request was handled.
// The instance host is returned if not public host was set.
func (r *DomainCtx) RequestedHost() string {
	if r.PublicHost != "" {
		return r.PublicHost
	}
	return r.InstanceHost
}

// RequestedDomain returns the domain (hostname) for which the request was handled.
// The instance domain is returned if no public host / domain was set.
func (r *DomainCtx) RequestedDomain() string {
	return hostnameFromHost(r.RequestedHost())
}

// Origin returns the origin (protocol://hostname[:port]) for which the request was handled.
// The instance host is used if no public host was set.
func (r *DomainCtx) Origin() string {
	host := r.PublicHost
	if host == "" {
		host = r.InstanceHost
	}
	return fmt.Sprintf("%s://%s", r.Protocol, host)
}

func DomainContext(ctx context.Context) *DomainCtx {
	o, ok := ctx.Value(domainCtx).(*DomainCtx)
	if !ok {
		return &DomainCtx{}
	}
	return o
}

func WithDomainContext(ctx context.Context, domainContext *DomainCtx) context.Context {
	return context.WithValue(ctx, domainCtx, domainContext)
}

func WithRequestedHost(ctx context.Context, host string) context.Context {
	i, ok := ctx.Value(domainCtx).(*DomainCtx)
	if !ok {
		i = new(DomainCtx)
	}

	i.PublicHost = host
	return context.WithValue(ctx, domainCtx, i)
}
