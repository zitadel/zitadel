package handlers

import (
	"context"
	"net/url"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
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
		originURL, err := url.Parse(origin)
		if err != nil {
			return ctx, err
		}
		return enrichCtx(ctx, originURL.Hostname(), origin), nil
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
		return ctx, errors.ThrowInternal(nil, "NOTIF-Ef3r1", "Errors.Notification.NoDomain")
	}
	return enrichCtx(
		ctx,
		domains.Domains[0].Domain,
		http_utils.BuildHTTP(domains.Domains[0].Domain, n.externalPort, n.externalSecure),
	), nil
}

func enrichCtx(ctx context.Context, host, origin string) context.Context {
	ctx = authz.WithRequestedDomain(ctx, host)
	ctx = http_utils.WithComposedOrigin(ctx, origin)
	return ctx
}
