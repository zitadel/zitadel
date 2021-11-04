package view

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

// func (v *View) ApplicationByID(projectID, appID string) (*query.App, error) {
// 	return v.Query.AppByProjectAndAppID(context.TODO(), projectID, appID)
// }

func (v *View) ApplicationByOIDCClientID(clientID string) (*query.App, error) {
	return v.Query.AppByOIDCClientID(context.TODO(), clientID)
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
		Queries: make([]query.SearchQuery, 0, 2),
	}

	queries.Queries = append(queries.Queries, nameQuery, projectQuery)

	apps, err := v.Query.SearchApps(ctx, queries)
	if err != nil {
		return nil, err
	}
	if len(apps.Apps) != 0 {
		return nil, errors.ThrowNotFound(nil, "VIEW-svLQq", "app not found")
	}

	return apps.Apps[0], nil
}

func (v *View) SearchApplications(request *query.AppSearchQueries) (*query.Apps, error) {
	return v.Query.SearchApps(context.TODO(), request)
}
