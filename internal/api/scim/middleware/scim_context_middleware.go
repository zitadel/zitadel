package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zitadel/zitadel/internal/api/authz"
	zhttp "github.com/zitadel/zitadel/internal/api/http/middleware"
	smetadata "github.com/zitadel/zitadel/internal/api/scim/metadata"
	sresources "github.com/zitadel/zitadel/internal/api/scim/resources"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func ScimContextMiddleware(q *query.Queries) func(next zhttp.HandlerFuncWithError) zhttp.HandlerFuncWithError {
	return func(next zhttp.HandlerFuncWithError) zhttp.HandlerFuncWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx, err := initScimContext(r.Context(), q)
			if err != nil {
				return err
			}

			return next(w, r.WithContext(ctx))
		}
	}
}

func initScimContext(ctx context.Context, q *query.Queries) (context.Context, error) {
	data := smetadata.NewScimContextData()
	ctx = smetadata.SetScimContextData(ctx, data)

	userID := authz.GetCtxData(ctx).UserID

	// get the provisioningDomain and ignorePassword metadata keys associated with the service user
	metadataKeys := []smetadata.Key{
		smetadata.KeyProvisioningDomain,
		smetadata.KeyIgnorePasswordOnCreate,
	}
	queries := sresources.BuildMetadataQueries(ctx, metadataKeys)

	metadataList, err := q.SearchUserMetadata(ctx, false, userID, queries, nil)
	if err != nil {
		if zerrors.IsNotFound(err) {
			return ctx, nil
		}
		return ctx, err
	}

	if metadataList == nil || len(metadataList.Metadata) == 0 {
		return ctx, nil
	}

	for _, metadata := range metadataList.Metadata {
		switch metadata.Key {
		case string(smetadata.KeyProvisioningDomain):
			data.ProvisioningDomain = string(metadata.Value)
			if data.ProvisioningDomain != "" {
				data.ExternalIDScopedMetadataKey = smetadata.ScopeExternalIdKey(data.ProvisioningDomain)
			}
		case string(smetadata.KeyIgnorePasswordOnCreate):
			ignorePasswordOnCreate, err := strconv.ParseBool(string(metadata.Value))
			if err == nil {
				data.IgnorePasswordOnCreate = ignorePasswordOnCreate
			}
		}
	}
	return smetadata.SetScimContextData(ctx, data), nil
}
