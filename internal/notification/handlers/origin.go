package handlers

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/api/authz"

	"github.com/zitadel/zitadel/internal/eventstore"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
)

type BaseURLEvent interface {
	eventstore.Event
	GetBaseURL() string
}

func (n *NotificationQueries) Origin(ctx context.Context, e eventstore.Event) (string, error) {
	baseURLEvent, ok := e.(BaseURLEvent)
	if !ok {
		return "", errors.ThrowInternal(fmt.Errorf("event of type %T doesn't implement BaseURLEvent", e), "NOTIF-3m9fs", "Errors.Internal")
	}
	baseURL := baseURLEvent.GetBaseURL()
	if baseURL != "" {
		return baseURL, nil
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
