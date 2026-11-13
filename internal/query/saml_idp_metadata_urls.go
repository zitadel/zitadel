package query

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//go:embed saml_idp_metadata_urls_list.sql
	samlIDPMetadataURLsListQuery string
)

// SAMLIDPMetadataURL describes a SAML IdP which fetches its metadata from a URL.
type SAMLIDPMetadataURL struct {
	ID            string
	InstanceID    string
	ResourceOwner string
	OwnerType     domain.IdentityProviderType
	MetadataURL   string
}

// ListSAMLIDPMetadataURLs returns a page of active SAML IdPs configured with a
// metadata URL. Results are ordered by instance and provider ID.
func (q *Queries) ListSAMLIDPMetadataURLs(ctx context.Context, afterInstanceID, afterID string, limit uint32) (result []SAMLIDPMetadataURL, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var idp SAMLIDPMetadataURL
			if err := rows.Scan(&idp.ID, &idp.InstanceID, &idp.ResourceOwner, &idp.OwnerType, &idp.MetadataURL); err != nil {
				return zerrors.ThrowInternal(err, "QUERY-u3o9Fa2l", "Errors.Internal")
			}
			result = append(result, idp)
		}
		return nil
	}, samlIDPMetadataURLsListQuery, domain.IDPTypeSAML, domain.IDPStateActive, afterInstanceID, afterID, limit)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-x7Bd1Qm4", "Errors.Internal")
	}
	return result, nil
}
