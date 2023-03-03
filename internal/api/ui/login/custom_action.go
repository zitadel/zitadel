package login

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dop251/goja"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/actions/object"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
)

func (l *Login) runPostExternalAuthenticationActions(
	user *domain.ExternalUser,
	tokens *oidc.Tokens,
	authRequest *domain.AuthRequest,
	httpRequest *http.Request,
	idpUser idp.User,
	authenticationError error,
) (*domain.ExternalUser, error) {
	ctx := httpRequest.Context()

	resourceOwner := authRequest.RequestedOrgID
	if resourceOwner == "" {
		resourceOwner = authz.GetInstance(ctx).DefaultOrganisationID()
	}
	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeExternalAuthentication, domain.TriggerTypePostAuthentication, resourceOwner, false)
	if err != nil {
		return nil, err
	}

	metadataList := object.MetadataListFromDomain(user.Metadatas)
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
		actions.SetFields("setPreferredUsername", func(username string) {
			user.PreferredUsername = username
		}),
		actions.SetFields("setEmail", func(email string) {
			user.Email = email
		}),
		actions.SetFields("setEmailVerified", func(verified bool) {
			user.IsEmailVerified = verified
		}),
		actions.SetFields("setPhone", func(phone string) {
			user.Phone = phone
		}),
		actions.SetFields("setPhoneVerified", func(verified bool) {
			user.IsPhoneVerified = verified
		}),
		actions.SetFields("metadata", &metadataList.Metadata),
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
			append(actions.ActionToOptions(a), actions.WithHTTP(actionCtx))...,
		)
		cancel()
		if err != nil {
			return nil, err
		}
	}
	user.Metadatas = object.MetadataListToDomain(metadataList)
	return user, err
}

type authMethod string

const (
	authMethodPassword     authMethod = "password"
	authMethodOTP          authMethod = "OTP"
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

	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeInternalAuthentication, domain.TriggerTypePostAuthentication, resourceOwner, false)
	if err != nil {
		return nil, err
	}

	metadataList := object.MetadataListFromDomain(nil)
	apiFields := actions.WithAPIFields(
		actions.SetFields("metadata", &metadataList.Metadata),
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
			append(actions.ActionToOptions(a), actions.WithHTTP(actionCtx))...,
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

	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, flowType, domain.TriggerTypePreCreation, resourceOwner, false)
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
		actions.SetFields("setEmail", func(email string) {
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
		actions.SetFields("setPhone", func(email string) {
			if user.Phone == nil {
				user.Phone = &domain.Phone{}
			}
			user.Phone.PhoneNumber = email
		}),
		actions.SetFields("setPhoneVerified", func(verified bool) {
			if user.Phone == nil {
				return
			}
			user.Phone.IsPhoneVerified = verified
		}),
		actions.SetFields("metadata", &metadataList.Metadata),
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
			append(actions.ActionToOptions(a), actions.WithHTTP(actionCtx))...,
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

	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, flowType, domain.TriggerTypePostCreation, resourceOwner, false)
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
						user, err := l.query.GetUserByID(actionCtx, true, userID, false)
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
			append(actions.ActionToOptions(a), actions.WithHTTP(actionCtx))...,
		)
		cancel()
		if err != nil {
			return nil, err
		}
	}
	return object.UserGrantsToDomain(userID, mutableUserGrants.UserGrants), err
}

func tokenCtxFields(tokens *oidc.Tokens) []actions.FieldOption {
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
			return tokens.IDTokenClaims.GetClaim(claim)
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
