package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// AddDynamicOIDCClient registers a new OIDC application through OAuth 2.0 Dynamic Client
// Registration (RFC 7591). It reuses the regular OIDC application persistence
// (pushOIDCApplication, i.e. the same application-added and oidc-config-added events) so a
// dynamically registered client is an ordinary OIDC app and the whole token, authorization
// and introspection flow keeps working unchanged.
//
// Unlike AddOIDCApplication it does NOT perform an app.write permission check: dynamic
// client registration is authorized at the registration endpoint through the
// oidc_dynamic_client_registration feature flag and the configured registration mode (open
// or initial access token). The caller is responsible for providing the target project
// (see the registration endpoint's EnsureDCRProject helper) and the owning organization.
func (c *Commands) AddDynamicOIDCClient(ctx context.Context, projectID, resourceOwner string, oidcApp *domain.OIDCApp) (_ *domain.OIDCApp, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if oidcApp == nil || projectID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Eip3p", "Errors.Project.App.Invalid")
	}
	oidcApp.AggregateID = projectID

	if !oidcApp.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Joh2a", "Errors.Project.App.Invalid")
	}

	if _, err = c.checkProjectExists(ctx, projectID, resourceOwner); err != nil {
		return nil, err
	}

	appID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	// The application name carries a unique constraint per project. As the client name
	// from the registration request is optional and not guaranteed to be unique (several
	// MCP clients self-register into the same project), derive a unique name from the
	// generated application ID.
	oidcApp.AppName = dynamicOIDCClientName(oidcApp.AppName, appID)

	addedApplication := NewOIDCApplicationWriteModel(projectID, resourceOwner)
	if err = c.eventstore.FilterToQueryReducer(ctx, addedApplication); err != nil {
		return nil, err
	}

	return c.pushOIDCApplication(ctx, addedApplication, oidcApp, appID)
}

// dynamicOIDCClientName builds a per-project unique application name for a dynamically
// registered client. The optional, non-unique client name from the request is kept for
// readability and disambiguated with the unique application ID.
func dynamicOIDCClientName(requestedName, appID string) string {
	requestedName = strings.TrimSpace(requestedName)
	if requestedName == "" {
		return "DCR Client " + appID
	}
	return requestedName + " (" + appID + ")"
}
