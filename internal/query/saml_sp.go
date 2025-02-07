package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SAMLServiceProvider struct {
	InstanceID           string              `json:"instance_id,omitempty"`
	AppID                string              `json:"app_id,omitempty"`
	State                domain.AppState     `json:"state,omitempty"`
	EntityID             string              `json:"entity_id,omitempty"`
	Metadata             []byte              `json:"metadata,omitempty"`
	MetadataURL          string              `json:"metadata_url,omitempty"`
	ProjectID            string              `json:"project_id,omitempty"`
	ProjectRoleAssertion bool                `json:"project_role_assertion,omitempty"`
	LoginVersion         domain.LoginVersion `json:"login_version,omitempty"`
	LoginBaseURI         *URL                `json:"login_base_uri,omitempty"`
	ProjectRoleKeys      []string            `json:"project_role_keys,omitempty"`
}

//go:embed saml_sp_by_id.sql
var samlSPQuery string

func (q *Queries) ActiveSAMLServiceProviderByID(ctx context.Context, entityID string) (sp *SAMLServiceProvider, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	sp, err = database.QueryJSONObject[SAMLServiceProvider](ctx, q.client,
		samlSPQuery,
		authz.GetInstance(ctx).InstanceID(),
		entityID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, zerrors.ThrowNotFound(err, "QUERY-HeOcis2511", "Errors.App.NotFound")
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-OyJx1Rp30z", "Errors.Internal")
	}
	instance := authz.GetInstance(ctx)
	loginV2 := instance.Features().LoginV2
	if loginV2.Required {
		sp.LoginVersion = domain.LoginVersion2
		sp.LoginBaseURI = (*URL)(loginV2.BaseURI)
	}
	return sp, err
}
