package http

import (
	"context"
	"fmt"
	"strings"
)

type DomainCtx struct {
	InstanceHost string
	PublicHost   string
	Protocol     string
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
// The instance domain is returned if not public host / domain was set.
func (r *DomainCtx) RequestedDomain() string {
	return strings.Split(r.RequestedHost(), ":")[0]
}

// Origin returns the origin (protocol://hostname[:port]) for which the request was handled.
// The instance host is used if not public host was set.
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
