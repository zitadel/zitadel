package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
)

func (n *NotificationQueries) Origin(ctx context.Context) (context.Context, string, error) {
	primary, err := query.NewInstanceDomainPrimarySearchQuery(true)
	if err != nil {
		return ctx, "", err
	}
	domains, err := n.SearchInstanceDomains(ctx, &query.InstanceDomainSearchQueries{
		Queries: []query.SearchQuery{primary},
	})
	if err != nil {
		return ctx, "", err
	}
	if len(domains.Domains) < 1 {
		return ctx, "", errors.ThrowInternal(nil, "NOTIF-Ef3r1", "Errors.Notification.NoDomain")
	}
	ctx = authz.WithRequestedDomain(ctx, domains.Domains[0].Domain)
	return ctx, http_utils.BuildHTTP(domains.Domains[0].Domain, n.externalPort, n.externalSecure), nil
}
