package handlers

import (
	"context"
	"net/url"

	"github.com/zitadel/logging"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type OriginEvent interface {
	eventstore.Event
	TriggerOrigin() string
}

func (n *NotificationQueries) Origin(ctx context.Context, e eventstore.Event) (context.Context, error) {
	var origin string
	originEvent, ok := e.(OriginEvent)
	if !ok {
		logging.Errorf("event of type %T doesn't implement OriginEvent", e)
	} else {
		origin = originEvent.TriggerOrigin()
	}
	if origin != "" {
		originCtx, err := enrichCtx(ctx, origin)
		if err == nil {
			return originCtx, nil
		}
		// Origins of requests are validated before an event is created, but older events can contain an untrusted host
		// (e.g. "https://trusted.example.com:443@evil.example"). Fall back to the primary instance domain instead of failing the event.
		logging.WithFields("origin", origin, "event_type", e.Type()).WithError(err).Warning("invalid trigger origin, using the primary instance domain")
	}
	primary, err := query.NewInstanceDomainPrimarySearchQuery(true)
	if err != nil {
		return ctx, err
	}
	domains, err := n.SearchInstanceDomains(ctx, &query.InstanceDomainSearchQueries{
		Queries: []query.SearchQuery{primary},
	})
	if err != nil {
		return ctx, err
	}
	if len(domains.Domains) < 1 {
		return ctx, zerrors.ThrowInternal(nil, "NOTIF-Ef3r1", "Errors.Notification.NoDomain")
	}
	// the primary instance domain is known and does not need to be validated as an origin stored on an event
	u, err := url.Parse(http_utils.BuildHTTP(domains.Domains[0].Domain, n.externalPort, n.externalSecure))
	if err != nil {
		return ctx, err
	}
	return http_utils.WithDomainContext(ctx, http_utils.NewDomainCtxFromOrigin(u)), nil
}

// enrichCtx sets the domain context from an origin stored on an event or job.
// An empty origin (e.g. a back-channel logout not triggered by a request) results in an empty domain context.
func enrichCtx(ctx context.Context, origin string) (context.Context, error) {
	if origin == "" {
		return http_utils.WithDomainContext(ctx, http_utils.NewDomainCtxFromOrigin(&url.URL{})), nil
	}
	u, err := parseOrigin(origin)
	if err != nil {
		return nil, err
	}
	return http_utils.WithDomainContext(ctx, http_utils.NewDomainCtxFromOrigin(u)), nil
}

// parseOrigin parses an origin stored on an event and returns a plain http(s)://host[:port] URL,
// with the host validated the same way as the host headers of a request.
// Only the scheme and the validated host are kept; any path, query or fragment is dropped so a link
// cannot be built with a different host than the one checked, e.g. url.Parse treats
// "https://trusted.example.com:443@evil.example" as host "evil.example".
func parseOrigin(origin string) (*url.URL, error) {
	u, err := url.Parse(origin)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "NOTIF-Ws6dQ", "invalid origin")
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.User != nil {
		return nil, zerrors.ThrowInvalidArgument(nil, "NOTIF-Hq3zW", "invalid origin")
	}
	host, err := http_utils.ValidateHost(u.Host)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "NOTIF-Vb8kE", "invalid origin host")
	}
	return &url.URL{Scheme: u.Scheme, Host: host}, nil
}
