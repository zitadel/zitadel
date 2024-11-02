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
		return enrichCtx(ctx, origin)
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
	return enrichCtx(
		ctx,
		http_utils.BuildHTTP(domains.Domains[0].Domain, n.externalPort, n.externalSecure),
	)
}

func enrichCtx(ctx context.Context, origin string) (context.Context, error) {
	u, err := url.Parse(origin)
	if err != nil {
		return nil, err
	}
	ctx = http_utils.WithDomainContext(ctx, &http_utils.DomainCtx{
		InstanceHost: u.Host,
		PublicHost:   u.Host,
		Protocol:     u.Scheme,
	})
	return ctx, nil
}
