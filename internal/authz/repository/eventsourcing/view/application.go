package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (v *View) ApplicationByOIDCClientID(ctx context.Context, clientID string) (*query.App, error) {
	return v.Query.AppByOIDCClientID(ctx, clientID)
}

func (v *View) ApplicationByProjecIDAndAppName(ctx context.Context, projectID, appName string) (_ *query.App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	nameQuery, err := query.NewAppNameSearchQuery(query.TextEquals, appName)
	if err != nil {
		return nil, err
	}
	projectQuery, err := query.NewAppProjectIDSearchQuery(projectID)
	if err != nil {
		return nil, err
	}

	queries := &query.AppSearchQueries{
		Queries: []query.SearchQuery{
			nameQuery,
			projectQuery,
		},
	}

	apps, err := v.Query.SearchApps(ctx, queries)
	if err != nil {
		return nil, err
	}
	if len(apps.Apps) != 1 {
		return nil, errors.ThrowNotFound(nil, "VIEW-svLQq", "app not found")
	}

	return apps.Apps[0], nil
}

func (v *View) SearchApplications(request *query.AppSearchQueries) (*query.Apps, error) {
	return v.Query.SearchApps(context.TODO(), request)
}
