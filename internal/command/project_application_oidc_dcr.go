package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// DCRProjectName is the name of the dedicated project that holds clients created through
// OAuth 2.0 Dynamic Client Registration (RFC 7591). It is auto-provisioned per
// organization on the first registration.
const DCRProjectName = "ZITADEL DCR"

// EnsureDCRProject returns the dedicated project that holds dynamically registered OIDC
// clients in the given organization, creating it on first use. Like AddDynamicOIDCClient it
// does not perform a permission check, because the provisioning is authorized at the
// registration endpoint through the feature flag and the configured registration mode.
//
// Both the existence check and the recovery from a concurrent creation read from the
// eventstore, so that several clients self-registering into the same organization at the
// same time converge on a single project: the per-organization uniqueness of the project
// name lets exactly one push win, and the racing callers resolve the winner's id strongly
// consistently. Resolving from the projection instead would be eventually consistent and
// could still report no project in that race.
func (c *Commands) EnsureDCRProject(ctx context.Context, resourceOwner string) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if resourceOwner == "" {
		return "", zerrors.ThrowInvalidArgument(nil, "COMMAND-Oht9a", "Errors.ResourceOwnerMissing")
	}

	projectID, err := c.dcrProjectID(ctx, resourceOwner)
	if err != nil {
		return "", err
	}
	if projectID != "" {
		return projectID, nil
	}

	projectID, err = c.idGenerator.Next()
	if err != nil {
		return "", err
	}
	wm, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return "", err
	}
	events := []eventstore.Command{
		project.NewProjectAddedEvent(
			ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &wm.WriteModel),
			DCRProjectName,
			false,
			false,
			false,
			domain.PrivateLabelingSettingUnspecified,
		),
	}
	postCommit, err := c.projectCreatedMilestone(ctx, &events)
	if err != nil {
		return "", err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		if zerrors.IsErrorAlreadyExists(err) {
			// A concurrent registration created the project first; resolve the winner.
			return c.dcrProjectID(ctx, resourceOwner)
		}
		return "", err
	}
	postCommit(ctx)
	if err = AppendAndReduce(wm, pushedEvents...); err != nil {
		return "", err
	}
	return projectID, nil
}

// dcrProjectID resolves the id of an organization's dedicated DCR project from the
// eventstore (strongly consistent). It returns an empty string when the project does not
// exist yet.
func (c *Commands) dcrProjectID(ctx context.Context, resourceOwner string) (string, error) {
	wm := newDCRProjectWriteModel(resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
		return "", err
	}
	return wm.projectID, nil
}

// dcrProjectWriteModel resolves the dedicated DCR project of an organization by its name.
type dcrProjectWriteModel struct {
	eventstore.WriteModel
	projectID string
}

func newDCRProjectWriteModel(resourceOwner string) *dcrProjectWriteModel {
	return &dcrProjectWriteModel{
		WriteModel: eventstore.WriteModel{ResourceOwner: resourceOwner},
	}
}

func (wm *dcrProjectWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.ProjectAddedEvent:
			if e.Name == DCRProjectName {
				wm.projectID = e.Aggregate().ID
			}
		case *project.ProjectRemovedEvent:
			if e.Aggregate().ID == wm.projectID {
				wm.projectID = ""
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *dcrProjectWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		EventTypes(project.ProjectAddedType, project.ProjectRemovedType).
		Builder()
}

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
// (see the registration endpoint's ensureDCRProject helper) and the owning organization.
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

	// Persist the registration access token (RFC 7592 §3) atomically with the application,
	// so a registered client can always be managed and never ends up without a token.
	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &addedApplication.WriteModel)
	registrationTokenEvent, registrationToken, err := c.newOIDCRegistrationTokenEvent(ctx, projectAgg, appID)
	if err != nil {
		return nil, err
	}

	registered, err := c.pushOIDCApplication(ctx, addedApplication, oidcApp, appID, registrationTokenEvent)
	if err != nil {
		return nil, err
	}
	registered.RegistrationAccessToken = registrationToken
	return registered, nil
}

// UpdateDynamicOIDCClient applies an RFC 7592 update to a dynamically registered client and
// rotates its registration access token in the same push. Unlike UpdateOIDCApplication it does
// not perform an app.write permission check: the caller is authorized through the registration
// access token (see VerifyDynamicClientRegistrationToken). The returned application carries the
// new registration access token, which the client must use from then on. Updating metadata to
// the values it already has is not an error: the token is rotated and the current state is
// returned, as RFC 7592 expects a successful read-back from an update.
func (c *Commands) UpdateDynamicOIDCClient(ctx context.Context, oidcApp *domain.OIDCApp, resourceOwner string) (_ *domain.OIDCApp, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if oidcApp == nil || !oidcApp.IsValid() || oidcApp.AppID == "" || oidcApp.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Aef8a", "Errors.Project.App.OIDCConfigInvalid")
	}

	existingOIDC, err := c.getOIDCAppWriteModel(ctx, oidcApp.AggregateID, oidcApp.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOIDC.State == domain.AppStateUnspecified || existingOIDC.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Voo8i", "Errors.Project.App.NotExisting")
	}
	if !existingOIDC.IsOIDC() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Geph9", "Errors.Project.App.IsNotOIDC")
	}
	if err = c.eventstore.FilterToQueryReducer(ctx, existingOIDC); err != nil {
		return nil, err
	}

	changedEvent, hasChanged, err := c.oidcApplicationChangeEvent(ctx, existingOIDC, oidcApp)
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &existingOIDC.WriteModel)
	rotationEvent, registrationToken, err := c.newOIDCRegistrationTokenEvent(ctx, projectAgg, oidcApp.AppID)
	if err != nil {
		return nil, err
	}

	events := make([]eventstore.Command, 0, 2)
	if hasChanged {
		events = append(events, changedEvent)
	}
	events = append(events, rotationEvent)

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	if err = AppendAndReduce(existingOIDC, pushedEvents...); err != nil {
		return nil, err
	}

	result := oidcWriteModelToOIDCConfig(existingOIDC)
	result.FillCompliance()
	result.RegistrationAccessToken = registrationToken
	return result, nil
}

// RemoveDynamicOIDCClient deletes a dynamically registered client (RFC 7592 §2.3). Unlike
// RemoveApplication it does not perform an app.delete permission check: the caller is
// authorized through the registration access token. Removing the application also invalidates
// its registration access token.
func (c *Commands) RemoveDynamicOIDCClient(ctx context.Context, projectID, appID, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if projectID == "" || appID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Eiv0u", "Errors.IDMissing")
	}
	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Quu7n", "Errors.Project.App.NotExisting")
	}

	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &existingApp.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, project.NewApplicationRemovedEvent(ctx, projectAgg, appID, existingApp.Name, ""))
	if err != nil {
		return nil, err
	}
	if err = AppendAndReduce(existingApp, pushedEvents...); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingApp.WriteModel), nil
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

// newOIDCRegistrationTokenEvent generates a fresh registration access token secret for the
// application, returning the event that persists its hash and the plain secret to hand back
// to the client. Like a client secret, only the hash is ever stored; the plain secret leaves
// the server exactly once. The caller is responsible for pushing the event.
func (c *Commands) newOIDCRegistrationTokenEvent(ctx context.Context, projectAgg *eventstore.Aggregate, appID string) (eventstore.Command, string, error) {
	encodedHash, plain, err := c.newHashedSecret(ctx, c.eventstore.Filter) //nolint:staticcheck
	if err != nil {
		return nil, "", err
	}
	return project.NewOIDCConfigRegistrationTokenChangedEvent(ctx, projectAgg, appID, encodedHash), plain, nil
}

// VerifyDynamicClientRegistrationToken checks a presented registration access token secret
// (RFC 7592 §3) against the stored hash of the application's current token. It returns an
// unauthenticated error when no token is set or the secret does not match, so the management
// endpoints can answer with 401. The hash is read strongly consistently from the eventstore;
// the management endpoints serve the common case from the projection and only fall back here,
// most importantly for a token that was just rotated and is not projected yet.
func (c *Commands) VerifyDynamicClientRegistrationToken(ctx context.Context, projectID, appID, resourceOwner, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if projectID == "" || appID == "" || secret == "" {
		return zerrors.ThrowUnauthenticated(nil, "COMMAND-Ohj6e", "Errors.Token.Invalid")
	}
	wm := newOIDCRegistrationTokenWriteModel(projectID, appID, resourceOwner)
	if err = c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
		return err
	}
	if wm.hashedToken == "" {
		return zerrors.ThrowUnauthenticated(nil, "COMMAND-Eereu", "Errors.Token.Invalid")
	}
	// The registration access token is short lived (it is rotated on every update), so an
	// outdated hash is not persisted back; the updated hash from Verify is intentionally
	// ignored.
	if _, err = c.secretHasher.Verify(wm.hashedToken, secret); err != nil {
		return zerrors.ThrowUnauthenticated(err, "COMMAND-Ush2a", "Errors.Token.Invalid")
	}
	return nil
}

// oidcRegistrationTokenWriteModel resolves the current registration access token hash of an
// application from the eventstore. It tracks the latest token-changed event and clears the
// hash when the application is removed.
type oidcRegistrationTokenWriteModel struct {
	eventstore.WriteModel
	appID       string
	hashedToken string
}

func newOIDCRegistrationTokenWriteModel(projectID, appID, resourceOwner string) *oidcRegistrationTokenWriteModel {
	return &oidcRegistrationTokenWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		appID: appID,
	}
}

func (wm *oidcRegistrationTokenWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.OIDCConfigRegistrationTokenChangedEvent:
			if e.AppID == wm.appID {
				wm.hashedToken = e.HashedToken
			}
		case *project.ApplicationRemovedEvent:
			if e.AppID == wm.appID {
				wm.hashedToken = ""
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *oidcRegistrationTokenWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			project.OIDCConfigRegistrationTokenChangedType,
			project.ApplicationRemovedType,
		).
		Builder()
}
