package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/zitadel/logging"

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

	// get the provisioningDomain and ignorePasswordOnCreate metadata keys associated with the service user
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
			ignorePasswordOnCreate, parseErr := strconv.ParseBool(strings.TrimSpace(string(metadata.Value)))
			if parseErr != nil {
				return ctx,
					zerrors.ThrowInvalidArgumentf(nil, "SMCM-yvw2rt", "Invalid value for metadata key %s: %s", smetadata.KeyIgnorePasswordOnCreate, metadata.Value)
			}
			data.IgnorePasswordOnCreate = ignorePasswordOnCreate
		default:
			logging.WithFields("user_metadata_key", metadata.Key).Warn("unexpected metadata key")
		}
	}
	return smetadata.SetScimContextData(ctx, data), nil
}
