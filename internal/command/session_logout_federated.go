package command

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/url"

	"github.com/crewjam/saml"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	saml_provider "github.com/zitadel/zitadel/internal/idp/providers/saml"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/sessionlogout"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// FederatedLogoutDataFetcher abstracts the data fetching required for federated logout
// allowing the command layer to be decoupled from the full query layer
type FederatedLogoutDataFetcher interface {
	IDPUserLinks(ctx context.Context, queries *query.IDPUserLinksSearchQuery, permissionCheck domain.PermissionCheck) (*query.IDPUserLinks, error)
	IDPTemplateByID(ctx context.Context, shouldTriggerBulk bool, id string, withOwnerRemoved bool, permissionCheck domain.PermissionCheck, queries ...query.SearchQuery) (*query.IDPTemplate, error)
}

// FederatedLogoutEventstore abstracts the eventstore operations required for federated logout
type FederatedLogoutEventstore interface {
	Push(ctx context.Context, cmds ...eventstore.Command) ([]eventstore.Event, error)
	FilterToQueryReducer(ctx context.Context, reducer eventstore.QueryReducer) error
	Filter(ctx context.Context, searchQuery *eventstore.SearchQueryBuilder) ([]eventstore.Event, error)
}

type FederatedLogoutRequest struct {
	LogoutID              string
	SessionID             string
	PostLogoutRedirectURI string

	// For SAML
	SAMLRequestID   string
	SAMLBindingType string // "redirect" or "post"
	SAMLRedirectURL string
	SAMLPostURL     string
	SAMLRequest     string
	SAMLRelayState  string
}

// StartFederatedLogout initiates a federated logout for a V2 session
// It checks if the session requires federated logout and creates the appropriate logout intent
func (c *Commands) StartFederatedLogout(
	ctx context.Context,
	fetcher FederatedLogoutDataFetcher,
	es FederatedLogoutEventstore,
	sessionID string,
	postLogoutRedirectURI string,
) (_ *FederatedLogoutRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// 1. Get the V2 session to retrieve user and IdP information
	instanceID := authz.GetInstance(ctx).InstanceID()
	session := NewSessionWriteModel(sessionID, instanceID)
	err = es.FilterToQueryReducer(ctx, session)
	if err != nil {
		return nil, err
	}

	if session.UserID == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sf3g2", "Errors.Session.NotFound")
	}

	// 2. Check if session was authenticated via IdP
	// For V2 sessions, we need to check if there's an IdP link
	idpID, err := c.findSessionIDPID(ctx, fetcher, session.UserID, instanceID)
	if err != nil {
		logging.WithFields("sessionID", sessionID, "userID", session.UserID).
			WithError(err).Info("error finding session IDP ID for federated logout")
		return nil, nil
	}
	if idpID == "" {
		// No IdP authentication, no federated logout needed
		logging.WithFields("sessionID", sessionID, "userID", session.UserID).
			Info("no IDP link found for user - skipping federated logout")
		return nil, nil
	}

	logging.WithFields("sessionID", sessionID, "userID", session.UserID, "idpID", idpID).
		Info("found IDP link for federated logout")

	// 3. Get IdP configuration
	idp, err := fetcher.IDPTemplateByID(ctx, false, idpID, false, nil)
	if err != nil {
		logging.WithFields("sessionID", sessionID, "idpID", idpID).
			WithError(err).Error("failed to get IDP template for federated logout")
		return nil, err
	}
	if idp == nil {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-IDP404", "Errors.IDP.NotFound")
	}

	logging.WithFields("sessionID", sessionID, "idpID", idpID,
		"isSAML", idp.SAMLIDPTemplate != nil,
		"federatedLogoutEnabled", idp.FederatedLogoutEnabled).
		Info("checked IDP configuration for federated logout")

	// 4. Check if IdP has federated logout enabled
	if idp.SAMLIDPTemplate == nil {
		logging.WithFields("sessionID", sessionID, "idpID", idpID).
			Info("IDP is not a SAML IDP - skipping federated logout")
		return nil, nil
	}
	if !idp.FederatedLogoutEnabled {
		logging.WithFields("sessionID", sessionID, "idpID", idpID).
			Warn("federated logout is NOT enabled for this IDP - skipping federated logout")
		return nil, nil
	}

	logging.WithFields("sessionID", sessionID, "idpID", idpID).
		Info("federated logout is enabled - proceeding with SAML logout")

	// 5. Create logout intent
	logoutID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	logoutAggregate := sessionlogout.NewAggregate(logoutID, instanceID)

	// 6. Create started event
	startedEvent := sessionlogout.NewStartedEvent(
		ctx,
		&logoutAggregate.Aggregate,
		sessionID,
		idpID,
		session.UserID,
		postLogoutRedirectURI,
	)

	// 7. Get nameID for SAML logout
	nameID, err := c.findIDPUserNameID(ctx, fetcher, session.UserID, idpID)
	if err != nil {
		return nil, err
	}

	// 8. Generate SAML logout request using the crewjam/saml library (with signature)
	samlRequest, err := c.generateSAMLLogoutRequest(ctx, es, idp, session.UserID, nameID, logoutID, instanceID)
	if err != nil {
		logging.WithFields("sessionID", sessionID, "idpID", idpID).
			WithError(err).Error("failed to generate SAML logout request")
		return nil, err
	}

	logging.WithFields(
		"sessionID", sessionID,
		"logoutID", logoutID,
		"samlRequestID", samlRequest.RequestID,
		"bindingType", samlRequest.BindingType,
		"redirectURL", samlRequest.RedirectURL,
		"postURL", samlRequest.PostURL,
	).Info("generated SAML logout request successfully")

	// 9. Create SAML request event
	samlEvent := sessionlogout.NewSAMLRequestCreatedEvent(
		ctx,
		&logoutAggregate.Aggregate,
		samlRequest.RequestID,
		samlRequest.BindingType,
		samlRequest.RedirectURL,
		samlRequest.PostURL,
		samlRequest.SAMLRequest,
		logoutID, // Use logoutID as RelayState to correlate the response
	)

	// 10. Push events
	_, err = es.Push(ctx, startedEvent, samlEvent)
	if err != nil {
		return nil, err
	}

	return &FederatedLogoutRequest{
		LogoutID:              logoutID,
		SessionID:             sessionID,
		PostLogoutRedirectURI: postLogoutRedirectURI,
		SAMLRequestID:         samlRequest.RequestID,
		SAMLBindingType:       samlRequest.BindingType,
		SAMLRedirectURL:       samlRequest.RedirectURL,
		SAMLPostURL:           samlRequest.PostURL,
		SAMLRequest:           samlRequest.SAMLRequest,
		SAMLRelayState:        logoutID,
	}, nil
}

// GetFederatedLogoutRequest retrieves an existing federated logout request
func (c *Commands) GetFederatedLogoutRequest(
	ctx context.Context,
	logoutID string,
) (_ *FederatedLogoutRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	logoutModel := NewFederatedLogoutWriteModel(logoutID, instanceID)

	err = c.eventstore.FilterToQueryReducer(ctx, logoutModel)
	if err != nil {
		return nil, err
	}

	if !logoutModel.IsActive() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Mf4g2", "Errors.SessionLogout.NotFound")
	}

	return &FederatedLogoutRequest{
		LogoutID:              logoutID,
		SessionID:             logoutModel.SessionID,
		PostLogoutRedirectURI: logoutModel.PostLogoutRedirectURI,
		SAMLRequestID:         logoutModel.SAMLRequestID,
		SAMLBindingType:       logoutModel.SAMLBindingType,
		SAMLRedirectURL:       logoutModel.SAMLRedirectURL,
		SAMLPostURL:           logoutModel.SAMLPostURL,
		SAMLRequest:           logoutModel.SAMLRequest,
		SAMLRelayState:        logoutModel.SAMLRelayState,
	}, nil
}

// CompleteFederatedLogout marks a federated logout as completed
// This is called when the IdP responds to the logout request
func (c *Commands) CompleteFederatedLogout(
	ctx context.Context,
	logoutID string,
) (postLogoutRedirectURI string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	logoutModel := NewFederatedLogoutWriteModel(logoutID, instanceID)

	err = c.eventstore.FilterToQueryReducer(ctx, logoutModel)
	if err != nil {
		return "", err
	}

	if !logoutModel.IsActive() {
		return "", zerrors.ThrowNotFound(nil, "COMMAND-Nf5g3", "Errors.SessionLogout.NotFound")
	}

	logoutAggregate := sessionlogout.NewAggregate(logoutID, instanceID)

	responseEvent := sessionlogout.NewSAMLResponseReceivedEvent(
		ctx,
		&logoutAggregate.Aggregate,
		logoutModel.SAMLRequestID,
	)

	completedEvent := sessionlogout.NewCompletedEvent(
		ctx,
		&logoutAggregate.Aggregate,
	)

	_, err = c.eventstore.Push(ctx, responseEvent, completedEvent)
	if err != nil {
		return "", err
	}

	return logoutModel.PostLogoutRedirectURI, nil
}

// Helper methods

func (c *Commands) findSessionIDPID(ctx context.Context, fetcher FederatedLogoutDataFetcher, userID, instanceID string) (string, error) {
	// Query the user's IdP links to find the most recently used IdP
	userIDQuery, err := query.NewIDPUserLinksUserIDSearchQuery(userID)
	if err != nil {
		return "", err
	}

	links, err := fetcher.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{
		Queries: []query.SearchQuery{userIDQuery},
	}, nil) // Check what permission check should be used. Passing nil for now as command might bypass or handle internally
	if err != nil {
		return "", err
	}

	if len(links.Links) == 0 {
		return "", nil
	}

	// In a real scenario with multiple linked IdPs, we might want to store the
	// IdP used for the specific session in the session events.
	// For now, we take the most recent one or the first one if only one exists.
	// Since IDPUserLinks returns all links, we'll pick the first one.
	return links.Links[0].IDPID, nil
}

func (c *Commands) findIDPUserNameID(ctx context.Context, fetcher FederatedLogoutDataFetcher, userID, idpID string) (string, error) {
	// Query the user's external IdP link to get the nameID
	userIDQuery, err := query.NewIDPUserLinksUserIDSearchQuery(userID)
	if err != nil {
		return "", err
	}
	idpIDQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(idpID)
	if err != nil {
		return "", err
	}

	links, err := fetcher.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{
		Queries: []query.SearchQuery{userIDQuery, idpIDQuery},
	}, nil)
	if err != nil || len(links.Links) != 1 {
		return "", zerrors.ThrowPreconditionFailed(err, "COMMAND-Sf4g3", "Errors.User.ExternalIDP.NotFound")
	}

	return links.Links[0].ProvidedUserID, nil
}

func (c *Commands) findIDPSessionIndex(ctx context.Context, es FederatedLogoutEventstore, userID, idpID, instanceID string) (string, error) {
	// Try to find the most recent SAML IDPIntent for this user and IDP
	// The SessionIndex is stored in the encrypted Assertion in the IDPIntent

	logging.WithFields("userID", userID, "idpID", idpID, "instanceID", instanceID).
		Info("findIDPSessionIndex: starting search for SessionIndex")

	// Query for IDPIntents that belong to this user
	// We need to search by UserID in the SAMLSucceededEvent
	// IMPORTANT: Use correct aggregate and event type names from idpintent package
	// ORDER DESC to get the most recent events first!
	intentQuery := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(instanceID).
		AddQuery().
		AggregateTypes(idpintent.AggregateType).      // Use constant from package: "idpintent"
		EventTypes(idpintent.SAMLSucceededEventType). // Use constant from package: "idpintent.saml.succeeded"
		Builder().
		OrderDesc(). // CRITICAL: Get newest events first!
		Limit(10)    // Get last 10 to find the most recent one

	events, err := es.Filter(ctx, intentQuery)
	if err != nil {
		logging.WithFields("userID", userID, "idpID", idpID).
			WithError(err).Warn("findIDPSessionIndex: failed to query IDPIntent events for SessionIndex")
		return "", nil // Return empty string, not an error - SessionIndex is optional
	}

	logging.WithFields("userID", userID, "idpID", idpID, "eventCount", len(events)).
		Info("findIDPSessionIndex: found IDPIntent events (ordered DESC - newest first)")

	// Find the most recent SAMLSucceededEvent for this user and IDP
	var mostRecentAssertion *crypto.CryptoValue
	var foundEventUserID string
	var foundEventIntentID string

	// CRITICAL: Iterate from 0 to len (newest to oldest due to OrderDesc)
	for i := 0; i < len(events); i++ {
		event := events[i]
		logging.WithFields(
			"eventIndex", i,
			"eventType", event.Type(),
			"aggregateID", event.Aggregate().ID,
			"aggregateType", event.Aggregate().Type,
			"eventSequence", event.Sequence(),
		).Info("findIDPSessionIndex: examining event (index 0 = most recent, OrderDesc)")

		// Try type assertion with detailed logging
		samlEvent, ok := event.(*idpintent.SAMLSucceededEvent)
		logging.WithFields(
			"eventIndex", i,
			"typeAssertionOK", ok,
			"eventGoType", fmt.Sprintf("%T", event),
		).Info("findIDPSessionIndex: type assertion result")

		if ok {
			logging.WithFields(
				"eventUserID", samlEvent.UserID,
				"eventIDPUserID", samlEvent.IDPUserID,
				"targetUserID", userID,
				"hasAssertion", samlEvent.Assertion != nil,
				"intentID", event.Aggregate().ID,
			).Info("findIDPSessionIndex: found SAMLSucceededEvent")

			// IMPORTANT: During the FIRST login, UserID is empty because the user is created AFTER
			// the IDPIntent succeeds. So we can't filter by UserID.
			// Instead, we take the most recent event with an Assertion.
			// In a production system, you might want to add a link between Session and IDPIntent
			// to make this lookup more precise.
			if samlEvent.Assertion != nil {
				// We found a SAML succeeded event with an assertion
				mostRecentAssertion = samlEvent.Assertion
				foundEventUserID = samlEvent.UserID
				if foundEventUserID == "" {
					foundEventUserID = samlEvent.IDPUserID // Use IDP UserID as fallback
				}
				foundEventIntentID = event.Aggregate().ID

				logging.WithFields(
					"targetUserID", userID,
					"eventUserID", samlEvent.UserID,
					"eventIDPUserID", samlEvent.IDPUserID,
					"idpID", idpID,
					"intentID", foundEventIntentID,
				).Info("findIDPSessionIndex: found matching SAML assertion for SessionIndex extraction (taking most recent)")
				break
			}
		} else {
			logging.WithFields("eventIndex", i, "actualType", fmt.Sprintf("%T", event)).
				Warn("findIDPSessionIndex: event is not SAMLSucceededEvent")
		}
	}

	if mostRecentAssertion == nil {
		logging.WithFields("userID", userID, "idpID", idpID, "eventsChecked", len(events)).
			Warn("findIDPSessionIndex: no SAML assertion found in IDPIntent events")
		return "", nil
	}

	logging.WithFields("userID", foundEventUserID, "idpID", idpID, "intentID", foundEventIntentID).
		Info("findIDPSessionIndex: decrypting SAML assertion")

	// Decrypt the assertion
	assertionXML, err := crypto.Decrypt(mostRecentAssertion, c.idpConfigEncryption)
	if err != nil {
		logging.WithFields("userID", userID, "idpID", idpID).
			WithError(err).Warn("findIDPSessionIndex: failed to decrypt SAML assertion for SessionIndex")
		return "", nil
	}

	logging.WithFields("userID", userID, "idpID", idpID, "assertionLength", len(assertionXML)).
		Info("findIDPSessionIndex: assertion decrypted, parsing XML")

	// Parse the assertion XML to extract SessionIndex
	var assertion saml.Assertion
	err = xml.Unmarshal(assertionXML, &assertion)
	if err != nil {
		logging.WithFields("userID", userID, "idpID", idpID).
			WithError(err).Warn("findIDPSessionIndex: failed to unmarshal SAML assertion for SessionIndex")
		return "", nil
	}

	logging.WithFields("userID", userID, "idpID", idpID, "authnStatementCount", len(assertion.AuthnStatements)).
		Info("findIDPSessionIndex: assertion parsed")

	// Extract SessionIndex from AuthnStatement
	if len(assertion.AuthnStatements) > 0 {
		sessionIndex := assertion.AuthnStatements[0].SessionIndex
		logging.WithFields("userID", userID, "idpID", idpID, "sessionIndex", sessionIndex, "isEmpty", sessionIndex == "").
			Info("findIDPSessionIndex: extracted SessionIndex from AuthnStatement")

		if sessionIndex != "" {
			logging.WithFields("userID", userID, "idpID", idpID, "sessionIndex", sessionIndex).
				Info("findIDPSessionIndex: successfully extracted SessionIndex from SAML assertion")
			return sessionIndex, nil
		}
	}

	logging.WithFields("userID", userID, "idpID", idpID).
		Warn("findIDPSessionIndex: no SessionIndex found in SAML assertion AuthnStatement")
	return "", nil
}

type SAMLLogoutRequestData struct {
	RequestID   string
	BindingType string
	RedirectURL string
	PostURL     string
	SAMLRequest string
}

func (c *Commands) generateSAMLLogoutRequest(
	ctx context.Context,
	es FederatedLogoutEventstore,
	idp *query.IDPTemplate,
	userID string,
	nameID string,
	relayState string,
	instanceID string,
) (*SAMLLogoutRequestData, error) {
	if idp.SAMLIDPTemplate == nil {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sg43g", "Errors.IDP.SAML.NotConfigured")
	}

	samlTemplate := idp.SAMLIDPTemplate

	// Decrypt the private key
	key, err := crypto.Decrypt(samlTemplate.Key, c.idpConfigEncryption)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "COMMAND-SAMLKey", "Errors.IDP.SAML.KeyDecryptFailed")
	}

	// Build the SP Entity ID and root URL using Origin() like Login V1 does
	// Origin() returns protocol://host where host can be PublicHost or InstanceHost
	origin := http_util.DomainContext(ctx).Origin()

	// The EntityID must match exactly what was registered in Keycloak
	// This is the same format used in Login V1: Origin + /idps/ + ID + /saml/metadata
	entityID := origin + "/idps/" + idp.ID + "/saml/metadata"

	// The rootURL is used by the SAML library to construct callback URLs
	// In Login V1 they use baseURL + EndpointExternalLogin + "/"
	// But for the IDP handler, we use /idps/{id}/
	rootURL := origin + "/idps/" + idp.ID + "/"

	// Try to retrieve SessionIndex from the most recent SAML IDPIntent
	sessionIndex, err := c.findIDPSessionIndex(ctx, es, userID, idp.ID, instanceID)
	if err != nil {
		// Log error but continue - SessionIndex is optional
		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
			"userID", userID,
			"idpID", idp.ID,
			"error", err.Error(),
		).Warn("SAML request LOGOUT! - error retrieving SessionIndex, continuing without it")
		sessionIndex = ""
	}

	if sessionIndex != "" {
		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
			"userID", userID,
			"idpID", idp.ID,
			"sessionIndex", sessionIndex,
		).Info("SAML request LOGOUT! - SessionIndex retrieved from IDPIntent Assertion")
	} else {
		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
			"userID", userID,
			"idpID", idp.ID,
		).Warn("SAML request LOGOUT! - no SessionIndex found in IDPIntent Assertion")
	}

	// Build provider options
	opts := []saml_provider.ProviderOpts{
		saml_provider.WithEntityID(entityID),
	}
	if samlTemplate.WithSignedRequest {
		opts = append(opts, saml_provider.WithSignedRequest())
	}
	if samlTemplate.Binding != "" {
		opts = append(opts, saml_provider.WithBinding(samlTemplate.Binding))
	}
	if samlTemplate.WithSignedRequest && samlTemplate.SignatureAlgorithm != "" {
		opts = append(opts, saml_provider.WithSignatureAlgorithm(samlTemplate.SignatureAlgorithm))
	}
	if samlTemplate.NameIDFormat.Valid {
		opts = append(opts, saml_provider.WithNameIDFormat(samlTemplate.NameIDFormat.V))
	}

	// Create the SAML provider using the internal library
	provider, err := saml_provider.New(
		idp.Name,
		rootURL,
		samlTemplate.Metadata,
		samlTemplate.Certificate,
		key,
		opts...,
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "COMMAND-SAMLProv", "Errors.IDP.SAML.ProviderCreationFailed")
	}

	// Get the ServiceProvider middleware
	sp, err := provider.GetSP()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "COMMAND-SAMLSP", "Errors.IDP.SAML.ServiceProviderFailed")
	}

	// Find SingleLogoutService endpoint
	// TEMPORARY: Prefer HTTP-Redirect first until Login V2 implements POST form rendering
	// POST binding requires direct HTML rendering which needs access to http.ResponseWriter
	// that is not available in the current flow (TerminateSessionFromRequest returns only a string)
	sloRedirect := sp.ServiceProvider.GetSLOBindingLocation(saml.HTTPRedirectBinding)
	sloPost := sp.ServiceProvider.GetSLOBindingLocation(saml.HTTPPostBinding)

	// Choose binding: TEMPORARY prefer Redirect until POST rendering is implemented in Login V2

	// Choose binding: TEMPORARY prefer Redirect until POST rendering is implemented in Login V2
	var sloLocation string
	var bindingType string
	if sloRedirect != "" {
		sloLocation = sloRedirect
		bindingType = "redirect"
		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
			"chosenBinding", "HTTP-Redirect",
			"location", sloLocation,
		).Info("SAML request LOGOUT! - using HTTP-Redirect binding (preferred for now)")
	} else if sloPost != "" {
		sloLocation = sloPost
		bindingType = "post"
		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
			"chosenBinding", "HTTP-POST",
			"location", sloLocation,
		).Warn("SAML request LOGOUT! - using HTTP-POST binding (may require Login V2 implementation)")
	} else {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SAMLNoSLO", "Errors.IDP.SAML.NoSingleLogoutService")
	}

	// Create the LogoutRequest using the library (this handles IssueInstant, ID generation, etc.)
	logoutRequest, err := sp.ServiceProvider.MakeLogoutRequest(sloLocation, nameID)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "COMMAND-SAMLLRReq", "Errors.IDP.SAML.LogoutRequestFailed")
	}

	// IMPORTANT: Remove NameQualifier and SPNameQualifier from NameID
	// The library MakeLogoutRequest() adds them automatically, but some IdPs don't want them
	// or expect different values, so we remove them to send a cleaner LogoutRequest
	if logoutRequest.NameID != nil {
		logoutRequest.NameID.NameQualifier = ""
		logoutRequest.NameID.SPNameQualifier = ""

		// Ensure NameID Format is set correctly if not already set
		if logoutRequest.NameID.Format == "" {
			logoutRequest.NameID.Format = string(saml.PersistentNameIDFormat)
		}
	}

	// Add SessionIndex to LogoutRequest if available
	// The SessionIndex is retrieved from the IDPIntent Assertion stored during login
	if sessionIndex != "" {
		logoutRequest.SessionIndex = &saml.SessionIndex{Value: sessionIndex}
		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
			"sessionIndex", sessionIndex,
		).Info("SAML request LOGOUT! - SessionIndex added to LogoutRequest")
	} else {
		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
		).Info("SAML request LOGOUT! - no SessionIndex available, proceeding without it (some IdPs work with just NameID)")
	}

	// Log the raw LogoutRequest details
	var sessionIndexValue string
	if logoutRequest.SessionIndex != nil {
		sessionIndexValue = logoutRequest.SessionIndex.Value
	}
	logging.WithFields(
		"SAML_TYPE", "SAML request LOGOUT!",
		"logoutRequest.ID", logoutRequest.ID,
		"logoutRequest.NameID.Value", logoutRequest.NameID.Value,
		"logoutRequest.SessionIndex", sessionIndexValue,
	).Info("SAML request LOGOUT! - LogoutRequest created with SessionIndex")

	var redirectURL, postURL, encodedRequest string

	if bindingType == "redirect" {
		// For LogoutRequest, Redirect() only accepts relayState
		// Unlike AuthnRequest, it does NOT automatically sign for HTTP-Redirect binding
		redirectURLParsed := logoutRequest.Redirect(relayState)

		// IMPORTANT: Manually add signature for HTTP-Redirect binding
		// This replicates what AuthnRequest.Redirect() does when passed a ServiceProvider
		if sp.ServiceProvider.SignatureMethod != "" {
			// Get the SAMLRequest from the generated URL
			samlRequest := redirectURLParsed.Query().Get("SAMLRequest")

			// Construct the query string manually to ensure the correct order:
			// SAMLRequest=...&RelayState=...&SigAlg=...
			// This is REQUIRED by the SAML specification (Section 3.4.4.1)
			// Go's url.Values.Encode() sorts keys alphabetically, which puts RelayState first

			var queryBuf bytes.Buffer
			queryBuf.WriteString("SAMLRequest=")
			queryBuf.WriteString(url.QueryEscape(samlRequest))

			if relayState != "" {
				queryBuf.WriteString("&RelayState=")
				queryBuf.WriteString(url.QueryEscape(relayState))
			}

			queryBuf.WriteString("&SigAlg=")
			queryBuf.WriteString(url.QueryEscape(sp.ServiceProvider.SignatureMethod))

			query := queryBuf.String()

			// Get signing context
			signingContext, err := saml.GetSigningContext(&sp.ServiceProvider)
			if err != nil {
				logging.WithFields(
					"SAML_TYPE", "SAML request LOGOUT!",
					"error", err.Error(),
				).Warn("SAML request LOGOUT! - failed to get signing context, sending unsigned")
			} else {
				// Sign the query string
				sig, err := signingContext.SignString(query)
				if err != nil {
					logging.WithFields(
						"SAML_TYPE", "SAML request LOGOUT!",
						"error", err.Error(),
					).Warn("SAML request LOGOUT! - failed to sign query string, sending unsigned")
				} else {
					// Add Signature parameter
					query += "&Signature=" + url.QueryEscape(base64.StdEncoding.EncodeToString(sig))
					logging.WithFields(
						"SAML_TYPE", "SAML request LOGOUT!",
						"signatureMethod", sp.ServiceProvider.SignatureMethod,
					).Info("SAML request LOGOUT! - LogoutRequest signed successfully with strict parameter ordering")
				}
			}

			// Update the URL with the signed query
			redirectURLParsed.RawQuery = query
		} else {
			logging.WithFields(
				"SAML_TYPE", "SAML request LOGOUT!",
			).Warn("SAML request LOGOUT! - SignatureMethod not configured, sending unsigned LogoutRequest")
		}

		redirectURL = redirectURLParsed.String()

		// Extract the SAMLRequest from the URL for storage
		parsedURL, err := url.Parse(redirectURL)
		if err == nil {
			encodedRequest = parsedURL.Query().Get("SAMLRequest")
		}

		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
			"requestID", logoutRequest.ID,
			"redirectURL", redirectURL,
		).Info("SAML request LOGOUT! - generated LogoutRequest (HTTP-Redirect binding)")
	} else {
		// For POST binding, use the library's Post method
		postURL = sloLocation
		postForm := logoutRequest.Post(relayState)
		// The Post method returns []byte HTML form, we store the encoded request as string
		encodedRequest = string(postForm)

		// DECODE AND LOG THE SAML REQUEST FOR DEBUGGING (POST binding)
		// Extract SAMLRequest from the HTML form
		decodedXML := decodeSAMLRequestFromPost(encodedRequest)
		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
			"requestID", logoutRequest.ID,
			"postURL", postURL,
			"decodedXML", decodedXML,
			"postFormLength", len(encodedRequest),
		).Info("SAML request LOGOUT! - LogoutRequest DECODED (HTTP-POST binding)")

		logging.WithFields(
			"SAML_TYPE", "SAML request LOGOUT!",
			"requestID", logoutRequest.ID,
			"postURL", postURL,
			"postForm", encodedRequest,
		).Info("SAML request LOGOUT! - generated signed LogoutRequest (HTTP-POST binding)")
	}

	return &SAMLLogoutRequestData{
		RequestID:   logoutRequest.ID,
		BindingType: bindingType,
		RedirectURL: redirectURL,
		PostURL:     postURL,
		SAMLRequest: encodedRequest,
	}, nil
}

// decodeSAMLRequest decodes a SAML request from base64+deflate encoding for logging purposes
// decodeSAMLRequestFromPost extracts and decodes a SAML request from HTML POST form
func decodeSAMLRequestFromPost(htmlForm string) string {
	if htmlForm == "" {
		return "<empty>"
	}

	// Extract the SAMLRequest value from the HTML form
	// The form typically contains: <input type="hidden" name="SAMLRequest" value="base64-encoded-xml"/>
	// We'll parse it simply by looking for the value attribute
	start := bytes.Index([]byte(htmlForm), []byte(`name="SAMLRequest" value="`))
	if start == -1 {
		return "<SAMLRequest field not found in POST form>"
	}
	start += len(`name="SAMLRequest" value="`)

	end := bytes.Index([]byte(htmlForm[start:]), []byte(`"`))
	if end == -1 {
		return "<SAMLRequest value end quote not found>"
	}

	encodedRequest := htmlForm[start : start+end]

	// For POST binding, the request is only base64 encoded (NOT deflated)
	decoded, err := base64.StdEncoding.DecodeString(encodedRequest)
	if err != nil {
		return "<base64 decode error: " + err.Error() + ">"
	}

	return string(decoded)
}
