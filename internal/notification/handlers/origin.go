package handlers

import (
	"context"
	"fmt"

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

func (n *NotificationQueries) Origin(ctx context.Context, e eventstore.Event) (string, error) {
	originEvent, ok := e.(OriginEvent)
	if !ok {
		return "", errors.ThrowInternal(fmt.Errorf("event of type %T doesn't implement OriginEvent", e), "NOTIF-3m9fs", "Errors.Internal")
	}
	origin := originEvent.TriggerOrigin()
	if origin != "" {
		return origin, nil
	}
	primary, err := query.NewInstanceDomainPrimarySearchQuery(true)
	if err != nil {
		return "", err
	}
	domains, err := n.SearchInstanceDomains(ctx, &query.InstanceDomainSearchQueries{
		Queries: []query.SearchQuery{primary},
	})
	if err != nil {
		return "", err
	}
	if len(domains.Domains) < 1 {
		return "", errors.ThrowInternal(nil, "NOTIF-Ef3r1", "Errors.Notification.NoDomain")
	}
	ctx = authz.WithRequestedDomain(ctx, domains.Domains[0].Domain)
	return http_utils.BuildHTTP(domains.Domains[0].Domain, n.externalPort, n.externalSecure), nil
}
