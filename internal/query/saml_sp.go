package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"net/url"

	"github.com/zitadel/zitadel/internal/api/authz"
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
	LoginBaseURI         *url.URL            `json:"login_base_uri,omitempty"`
}

//go:embed saml_sp_by_id.sql
var samlSPQuery string

func (q *Queries) ActiveSAMLServiceProviderByID(ctx context.Context, entityID string) (sp *SAMLServiceProvider, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		sp, err = scanSAMLServiceProviderByID(row)
		return err
	}, samlSPQuery,
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
		sp.LoginBaseURI = loginV2.BaseURI
	}
	return sp, err
}

func scanSAMLServiceProviderByID(row *sql.Row) (*SAMLServiceProvider, error) {
	var instanceID, appID, entityID, metadataURL, projectID sql.NullString
	var projectRoleAssertion sql.NullBool
	var metadata []byte
	var state, loginVersion sql.NullInt16
	var loginBaseURI sql.NullString

	err := row.Scan(
		&instanceID,
		&appID,
		&state,
		&entityID,
		&metadata,
		&metadataURL,
		&projectID,
		&projectRoleAssertion,
		&loginVersion,
		&loginBaseURI,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowNotFound(err, "QUERY-8cjj8ao6yY", "Errors.App.NotFound")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-1xzFD209Bp", "Errors.Internal")
	}
	sp := &SAMLServiceProvider{
		InstanceID:           instanceID.String,
		AppID:                appID.String,
		State:                domain.AppState(state.Int16),
		EntityID:             entityID.String,
		Metadata:             metadata,
		MetadataURL:          metadataURL.String,
		ProjectID:            projectID.String,
		ProjectRoleAssertion: projectRoleAssertion.Bool,
	}
	if loginVersion.Valid {
		sp.LoginVersion = domain.LoginVersion(loginVersion.Int16)
	}
	if loginBaseURI.Valid && loginBaseURI.String != "" {
		url, err := url.Parse(loginBaseURI.String)
		if err != nil {
			return nil, err
		}
		sp.LoginBaseURI = url
	}
	return sp, nil
}
