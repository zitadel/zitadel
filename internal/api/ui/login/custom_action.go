package login

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/actions/object"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/query"
)

func (l *Login) runPostExternalAuthenticationActions(
	user *domain.ExternalUser,
	tokens *oidc.Tokens[*oidc.IDTokenClaims],
	authRequest *domain.AuthRequest,
	httpRequest *http.Request,
	idpUser idp.User,
	authenticationError error,
) (_ *domain.ExternalUser, userChanged bool, err error) {
	ctx := httpRequest.Context()

	// use the request org (scopes or domain discovery) as default
	resourceOwner := authRequest.RequestedOrgID
	// if the user was already linked to an IDP and redirected to that, the requested org might be empty
	if resourceOwner == "" {
		resourceOwner = authRequest.UserOrgID
	}
	// if we will have no org (e.g. user clicked directly on the IDP on the login page)
	if resourceOwner == "" {
		// in this case the user might nevertheless already be linked to an IDP,
		// so let's do a workaround and resourceOwnerOfUserIDPLink if there would be a IDP link
		resourceOwner, err = l.resourceOwnerOfUserIDPLink(ctx, authRequest.SelectedIDPConfigID, user.ExternalUserID)
		logging.WithFields("authReq", authRequest.ID, "idpID", authRequest.SelectedIDPConfigID).OnError(err).
			Warn("could not determine resource owner for runPostExternalAuthenticationActions, fall back to default org id")
	}
	// fallback to default org id
	if resourceOwner == "" {
		resourceOwner = authz.GetInstance(ctx).DefaultOrganisationID()
	}
	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeExternalAuthentication, domain.TriggerTypePostAuthentication, resourceOwner)
	if err != nil {
		return nil, false, err
	}

	metadataList := object.MetadataListFromDomain(user.Metadatas)
	apiFields := actions.WithAPIFields(
		actions.SetFields("setFirstName", func(firstName string) {
			user.FirstName = firstName
			userChanged = true
		}),
		actions.SetFields("setLastName", func(lastName string) {
			user.LastName = lastName
			userChanged = true
		}),
		actions.SetFields("setNickName", func(nickName string) {
			user.NickName = nickName
			userChanged = true
		}),
		actions.SetFields("setDisplayName", func(displayName string) {
			user.DisplayName = displayName
			userChanged = true
		}),
		actions.SetFields("setPreferredLanguage", func(preferredLanguage string) {
			user.PreferredLanguage = language.Make(preferredLanguage)
			userChanged = true
		}),
		actions.SetFields("setPreferredUsername", func(username string) {
			user.PreferredUsername = username
			userChanged = true
		}),
		actions.SetFields("setEmail", func(email domain.EmailAddress) {
			user.Email = email
			userChanged = true
		}),
		actions.SetFields("setEmailVerified", func(verified bool) {
			user.IsEmailVerified = verified
			userChanged = true
		}),
		actions.SetFields("setPhone", func(phone domain.PhoneNumber) {
			user.Phone = phone
			userChanged = true
		}),
		actions.SetFields("setPhoneVerified", func(verified bool) {
			user.IsPhoneVerified = verified
			userChanged = true
		}),
		actions.SetFields("metadata", func(c *actions.FieldConfig) interface{} {
			return metadataList.MetadataListFromDomain(c.Runtime)
		}),
		actions.SetFields("v1",
			actions.SetFields("user",
				actions.SetFields("appendMetadata", metadataList.AppendMetadataFunc),
			),
		),
	)

	authErrStr := "none"
	if authenticationError != nil {
		authErrStr = authenticationError.Error()
	}

	for _, a := range triggerActions {
		actionCtx, cancel := context.WithTimeout(ctx, a.Timeout())

		ctxFieldOptions := append(tokenCtxFields(tokens),
			actions.SetFields("v1",
				actions.SetFields("externalUser", func(c *actions.FieldConfig) interface{} {
					return object.UserFromExternalUser(c, user)
				}),
				actions.SetFields("providerInfo", func(c *actions.FieldConfig) interface{} {
					return c.Runtime.ToValue(idpUser)
				}),
				actions.SetFields("authRequest", object.AuthRequestField(authRequest)),
				actions.SetFields("httpRequest", object.HTTPRequestField(httpRequest)),
				actions.SetFields("authError", authErrStr),
			),
		)

		ctxFields := actions.SetContextFields(ctxFieldOptions...)

		err = actions.Run(
			actionCtx,
			ctxFields,
			apiFields,
			a.Script,
			a.Name,
			append(actions.ActionToOptions(a), actions.WithHTTP(actionCtx), actions.WithUUID(actionCtx))...,
		)
		cancel()
		if err != nil {
			return nil, false, err
		}
	}
	user.Metadatas = object.MetadataListToDomain(metadataList)
	return user, userChanged, err
}

type authMethod string

const (
	authMethodPassword     authMethod = "password"
	authMethodOTP          authMethod = "OTP"
	authMethodOTPSMS       authMethod = "OTP SMS"
	authMethodOTPEmail     authMethod = "OTP Email"
	authMethodU2F          authMethod = "U2F"
	authMethodPasswordless authMethod = "passwordless"
)

func (l *Login) runPostInternalAuthenticationActions(
	authRequest *domain.AuthRequest,
	httpRequest *http.Request,
	authMethod authMethod,
	authenticationError error,
) ([]*domain.Metadata, error) {
	ctx := httpRequest.Context()

	resourceOwner := authRequest.RequestedOrgID
	if resourceOwner == "" {
		resourceOwner = authRequest.UserOrgID
	}

	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeInternalAuthentication, domain.TriggerTypePostAuthentication, resourceOwner)
	if err != nil {
		return nil, err
	}

	metadataList := object.MetadataListFromDomain(nil)
	apiFields := actions.WithAPIFields(
		actions.SetFields("metadata", func(c *actions.FieldConfig) interface{} {
			return metadataList.MetadataListFromDomain(c.Runtime)
		}),
		actions.SetFields("v1",
			actions.SetFields("user",
				actions.SetFields("appendMetadata", metadataList.AppendMetadataFunc),
			),
		),
	)
	for _, a := range triggerActions {
		actionCtx, cancel := context.WithTimeout(ctx, a.Timeout())

		authErrStr := "none"
		if authenticationError != nil {
			authErrStr = authenticationError.Error()
		}
		ctxFields := actions.SetContextFields(
			actions.SetFields("v1",
				actions.SetFields("authMethod", authMethod),
				actions.SetFields("authError", authErrStr),
				actions.SetFields("authRequest", object.AuthRequestField(authRequest)),
				actions.SetFields("httpRequest", object.HTTPRequestField(httpRequest)),
			),
		)

		err = actions.Run(
			actionCtx,
			ctxFields,
			apiFields,
			a.Script,
			a.Name,
			append(actions.ActionToOptions(a), actions.WithHTTP(actionCtx), actions.WithUUID(actionCtx))...,
		)
		cancel()
		if err != nil {
			return nil, err
		}
	}
	return object.MetadataListToDomain(metadataList), err
}

func (l *Login) runPreCreationActions(
	authRequest *domain.AuthRequest,
	httpRequest *http.Request,
	user *domain.Human,
	metadata []*domain.Metadata,
	resourceOwner string,
	flowType domain.FlowType,
) (*domain.Human, []*domain.Metadata, error) {
	ctx := httpRequest.Context()

	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, flowType, domain.TriggerTypePreCreation, resourceOwner)
	if err != nil {
		return nil, nil, err
	}

	metadataList := object.MetadataListFromDomain(metadata)
	apiFields := actions.WithAPIFields(
		actions.SetFields("setFirstName", func(firstName string) {
			user.FirstName = firstName
		}),
		actions.SetFields("setLastName", func(lastName string) {
			user.LastName = lastName
		}),
		actions.SetFields("setNickName", func(nickName string) {
			user.NickName = nickName
		}),
		actions.SetFields("setDisplayName", func(displayName string) {
			user.DisplayName = displayName
		}),
		actions.SetFields("setPreferredLanguage", func(preferredLanguage string) {
			user.PreferredLanguage = language.Make(preferredLanguage)
		}),
		actions.SetFields("setGender", func(gender domain.Gender) {
			user.Gender = gender
		}),
		actions.SetFields("setUsername", func(username string) {
			user.Username = username
		}),
		actions.SetFields("setEmail", func(email domain.EmailAddress) {
			if user.Email == nil {
				user.Email = &domain.Email{}
			}
			user.Email.EmailAddress = email
		}),
		actions.SetFields("setEmailVerified", func(verified bool) {
			if user.Email == nil {
				return
			}
			user.Email.IsEmailVerified = verified
		}),
		actions.SetFields("setPhone", func(phone domain.PhoneNumber) {
			if user.Phone == nil {
				user.Phone = &domain.Phone{}
			}
			user.Phone.PhoneNumber = phone
		}),
		actions.SetFields("setPhoneVerified", func(verified bool) {
			if user.Phone == nil {
				return
			}
			user.Phone.IsPhoneVerified = verified
		}),
		actions.SetFields("metadata", func(c *actions.FieldConfig) interface{} {
			return metadataList.MetadataListFromDomain(c.Runtime)
		}),
		actions.SetFields("v1",
			actions.SetFields("user",
				actions.SetFields("appendMetadata", metadataList.AppendMetadataFunc),
			),
		),
	)

	for _, a := range triggerActions {
		actionCtx, cancel := context.WithTimeout(ctx, a.Timeout())

		ctxOpts := actions.SetContextFields(
			actions.SetFields("v1",
				actions.SetFields("user", func(c *actions.FieldConfig) interface{} {
					return object.UserFromHuman(c, user)
				}),
				actions.SetFields("authRequest", object.AuthRequestField(authRequest)),
				actions.SetFields("httpRequest", object.HTTPRequestField(httpRequest)),
			),
		)

		err = actions.Run(
			actionCtx,
			ctxOpts,
			apiFields,
			a.Script,
			a.Name,
			append(actions.ActionToOptions(a), actions.WithHTTP(actionCtx), actions.WithUUID(actionCtx))...,
		)
		cancel()
		if err != nil {
			return nil, nil, err
		}
	}
	return user, object.MetadataListToDomain(metadataList), err
}

func (l *Login) runPostCreationActions(
	userID string,
	authRequest *domain.AuthRequest,
	httpRequest *http.Request,
	resourceOwner string,
	flowType domain.FlowType,
) ([]*domain.UserGrant, error) {
	ctx := httpRequest.Context()

	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, flowType, domain.TriggerTypePostCreation, resourceOwner)
	if err != nil {
		return nil, err
	}

	mutableUserGrants := &object.UserGrants{UserGrants: make([]object.UserGrant, 0)}

	apiFields := actions.WithAPIFields(
		actions.SetFields("userGrants", &mutableUserGrants.UserGrants),
		actions.SetFields("v1",
			actions.SetFields("appendUserGrant", object.AppendGrantFunc(mutableUserGrants)),
		),
	)

	for _, a := range triggerActions {
		actionCtx, cancel := context.WithTimeout(ctx, a.Timeout())

		ctxFields := actions.SetContextFields(
			actions.SetFields("v1",
				actions.SetFields("getUser", func(c *actions.FieldConfig) interface{} {
					return func(call goja.FunctionCall) goja.Value {
						user, err := l.query.GetUserByID(actionCtx, true, userID)
						if err != nil {
							panic(err)
						}
						return object.UserFromQuery(c, user)
					}
				}),
				actions.SetFields("authRequest", object.AuthRequestField(authRequest)),
				actions.SetFields("httpRequest", object.HTTPRequestField(httpRequest)),
			),
		)

		err = actions.Run(
			actionCtx,
			ctxFields,
			apiFields,
			a.Script,
			a.Name,
			append(actions.ActionToOptions(a), actions.WithHTTP(actionCtx), actions.WithUUID(actionCtx))...,
		)
		cancel()
		if err != nil {
			return nil, err
		}
	}
	return object.UserGrantsToDomain(userID, mutableUserGrants.UserGrants), err
}

func tokenCtxFields(tokens *oidc.Tokens[*oidc.IDTokenClaims]) []actions.FieldOption {
	var accessToken, idToken string
	getClaim := func(claim string) interface{} {
		return nil
	}
	claimsJSON := func() (string, error) {
		return "", nil
	}
	if tokens == nil {
		return []actions.FieldOption{
			actions.SetFields("accessToken", accessToken),
			actions.SetFields("idToken", idToken),
			actions.SetFields("getClaim", getClaim),
			actions.SetFields("claimsJSON", claimsJSON),
		}
	}
	accessToken = tokens.AccessToken
	idToken = tokens.IDToken
	if tokens.IDTokenClaims != nil {
		getClaim = func(claim string) interface{} {
			return tokens.IDTokenClaims.Claims[claim]
		}
		claimsJSON = func() (string, error) {
			c, err := json.Marshal(tokens.IDTokenClaims)
			if err != nil {
				return "", err
			}
			return string(c), nil
		}
	}
	return []actions.FieldOption{
		actions.SetFields("accessToken", accessToken),
		actions.SetFields("idToken", idToken),
		actions.SetFields("getClaim", getClaim),
		actions.SetFields("claimsJSON", claimsJSON),
	}
}

func (l *Login) resourceOwnerOfUserIDPLink(ctx context.Context, idpConfigID string, externalUserID string) (string, error) {
	idQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(idpConfigID)
	if err != nil {
		return "", err
	}
	externalIDQuery, err := query.NewIDPUserLinksExternalIDSearchQuery(externalUserID)
	if err != nil {
		return "", err
	}
	queries := []query.SearchQuery{
		idQuery, externalIDQuery,
	}
	links, err := l.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{Queries: queries}, false)
	if err != nil {
		return "", err
	}
	if len(links.Links) != 1 {
		return "", nil
	}
	return links.Links[0].ResourceOwner, nil
}
