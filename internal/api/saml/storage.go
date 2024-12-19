package saml

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"
	"github.com/zitadel/saml/pkg/provider"
	"github.com/zitadel/saml/pkg/provider/key"
	"github.com/zitadel/saml/pkg/provider/models"
	"github.com/zitadel/saml/pkg/provider/serviceprovider"
	"github.com/zitadel/saml/pkg/provider/xml/samlp"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/actions/object"
	"github.com/zitadel/zitadel/internal/activity"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ provider.EntityStorage = &Storage{}
var _ provider.IdentityProviderStorage = &Storage{}
var _ provider.AuthStorage = &Storage{}
var _ provider.UserStorage = &Storage{}

const (
	LoginClientHeader = "x-zitadel-login-client"
)

type Storage struct {
	certChan                   <-chan interface{}
	defaultCertificateLifetime time.Duration

	currentCACertificate       query.Certificate
	currentMetadataCertificate query.Certificate
	currentResponseCertificate query.Certificate

	locker               crdb.Locker
	certificateAlgorithm string
	encAlg               crypto.EncryptionAlgorithm
	certEncAlg           crypto.EncryptionAlgorithm

	eventstore *eventstore.Eventstore
	repo       repository.Repository
	command    *command.Commands
	query      *query.Queries

	defaultLoginURL   string
	defaultLoginURLv2 string
}

func (p *Storage) GetEntityByID(ctx context.Context, entityID string) (*serviceprovider.ServiceProvider, error) {
	app, err := p.query.ActiveAppBySAMLEntityID(ctx, entityID)
	if err != nil {
		return nil, err
	}
	return serviceprovider.NewServiceProvider(
		app.ID,
		&serviceprovider.Config{
			Metadata: app.SAMLConfig.Metadata,
		},
		func(id string) string {
			if strings.HasPrefix(id, command.IDPrefixV2) {
				return p.defaultLoginURLv2 + id
			}
			return p.defaultLoginURL + id
		},
	)
}

func (p *Storage) GetEntityIDByAppID(ctx context.Context, appID string) (string, error) {
	app, err := p.query.AppByID(ctx, appID, true)
	if err != nil {
		return "", err
	}
	return app.SAMLConfig.EntityID, nil
}

func (p *Storage) Health(context.Context) error {
	return nil
}

func (p *Storage) GetCA(ctx context.Context) (*key.CertificateAndKey, error) {
	return p.GetCertificateAndKey(ctx, crypto.KeyUsageSAMLCA)
}

func (p *Storage) GetMetadataSigningKey(ctx context.Context) (*key.CertificateAndKey, error) {
	return p.GetCertificateAndKey(ctx, crypto.KeyUsageSAMLMetadataSigning)
}

func (p *Storage) GetResponseSigningKey(ctx context.Context) (*key.CertificateAndKey, error) {
	return p.GetCertificateAndKey(ctx, crypto.KeyUsageSAMLResponseSinging)
}

func (p *Storage) CreateAuthRequest(ctx context.Context, req *samlp.AuthnRequestType, acsUrl, protocolBinding, relayState, applicationID string) (_ models.AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	headers, _ := http_utils.HeadersFromCtx(ctx)
	if loginClient := headers.Get(LoginClientHeader); loginClient != "" {
		return p.createAuthRequestLoginClient(ctx, req, acsUrl, protocolBinding, relayState, applicationID, loginClient)
	}
	return p.createAuthRequest(ctx, req, acsUrl, protocolBinding, relayState, applicationID)
}

func (p *Storage) createAuthRequestLoginClient(ctx context.Context, req *samlp.AuthnRequestType, acsUrl, protocolBinding, relayState, applicationID, loginClient string) (_ models.AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	samlRequest := &command.SAMLRequest{
		ApplicationID: applicationID,
		ACSURL:        acsUrl,
		RelayState:    relayState,
		RequestID:     req.Id,
		Binding:       protocolBinding,
		Issuer:        req.Issuer.Text,
		Destination:   req.Destination,
		LoginClient:   loginClient,
	}

	aar, err := p.command.AddSAMLRequest(ctx, samlRequest)
	if err != nil {
		return nil, err
	}
	return &AuthRequestV2{aar}, nil
}

func (p *Storage) createAuthRequest(ctx context.Context, req *samlp.AuthnRequestType, acsUrl, protocolBinding, relayState, applicationID string) (_ models.AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, zerrors.ThrowPreconditionFailed(nil, "SAML-sd436", "no user agent id")
	}

	authRequest := CreateAuthRequestToBusiness(ctx, req, acsUrl, protocolBinding, applicationID, relayState, userAgentID)

	resp, err := p.repo.CreateAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}

	return AuthRequestFromBusiness(resp)
}

func (p *Storage) AuthRequestByID(ctx context.Context, id string) (_ models.AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if strings.HasPrefix(id, command.IDPrefixV2) {
		req, err := p.command.GetCurrentSAMLRequest(ctx, id)
		if err != nil {
			return nil, err
		}
		return &AuthRequestV2{req}, nil
	}

	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, zerrors.ThrowPreconditionFailed(nil, "SAML-D3g21", "no user agent id")
	}
	resp, err := p.repo.AuthRequestByIDCheckLoggedIn(ctx, id, userAgentID)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (p *Storage) SetUserinfoWithUserID(ctx context.Context, applicationID string, userinfo models.AttributeSetter, userID string, attributes []int) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	user, err := p.query.GetUserByID(ctx, true, userID)
	if err != nil {
		return err
	}
	if user.State != domain.UserStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "SAML-S3gFd", "Errors.User.NotActive")
	}

	userGrants, err := p.getGrants(ctx, userID, applicationID)
	if err != nil {
		return err
	}

	customAttributes, err := p.getCustomAttributes(ctx, user, userGrants)
	if err != nil {
		return err
	}

	setUserinfo(user, userinfo, attributes, customAttributes)

	// trigger activity log for authentication for user
	activity.Trigger(ctx, user.ResourceOwner, user.ID, activity.SAMLResponse, p.eventstore.FilterToQueryReducer)
	return nil
}

func (p *Storage) SetUserinfoWithLoginName(ctx context.Context, userinfo models.AttributeSetter, loginName string, attributes []int) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	user, err := p.query.GetUserByLoginName(ctx, true, loginName)
	if err != nil {
		return err
	}
	if user.State != domain.UserStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "SAML-FJ262", "Errors.User.NotActive")
	}

	setUserinfo(user, userinfo, attributes, map[string]*customAttribute{})
	return nil
}

func setUserinfo(user *query.User, userinfo models.AttributeSetter, attributes []int, customAttributes map[string]*customAttribute) {
	for name, attr := range customAttributes {
		userinfo.SetCustomAttribute(name, "", attr.nameFormat, attr.attributeValue)
	}
	if len(attributes) == 0 {
		userinfo.SetUsername(user.PreferredLoginName)
		userinfo.SetUserID(user.ID)
		if user.Human == nil {
			return
		}
		userinfo.SetEmail(string(user.Human.Email))
		userinfo.SetSurname(user.Human.LastName)
		userinfo.SetGivenName(user.Human.FirstName)
		userinfo.SetFullName(user.Human.DisplayName)
		return
	}
	for _, attribute := range attributes {
		switch attribute {
		case provider.AttributeEmail:
			if user.Human != nil {
				userinfo.SetEmail(string(user.Human.Email))
			}
		case provider.AttributeSurname:
			if user.Human != nil {
				userinfo.SetSurname(user.Human.LastName)
			}
		case provider.AttributeFullName:
			if user.Human != nil {
				userinfo.SetFullName(user.Human.DisplayName)
			}
		case provider.AttributeGivenName:
			if user.Human != nil {
				userinfo.SetGivenName(user.Human.FirstName)
			}
		case provider.AttributeUsername:
			userinfo.SetUsername(user.PreferredLoginName)
		case provider.AttributeUserID:
			userinfo.SetUserID(user.ID)
		}
	}
}

func (p *Storage) getCustomAttributes(ctx context.Context, user *query.User, userGrants *query.UserGrants) (map[string]*customAttribute, error) {
	customAttributes := make(map[string]*customAttribute, 0)
	queriedActions, err := p.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeCustomizeSAMLResponse, domain.TriggerTypePreSAMLResponseCreation, user.ResourceOwner)
	if err != nil {
		return nil, err
	}
	ctxFields := actions.SetContextFields(
		actions.SetFields("v1",
			actions.SetFields("getUser", func(c *actions.FieldConfig) interface{} {
				return func(call goja.FunctionCall) goja.Value {
					return object.UserFromQuery(c, user)
				}
			}),
			actions.SetFields("user",
				actions.SetFields("getMetadata", func(c *actions.FieldConfig) interface{} {
					return func(goja.FunctionCall) goja.Value {
						resourceOwnerQuery, err := query.NewUserMetadataResourceOwnerSearchQuery(user.ResourceOwner)
						if err != nil {
							logging.WithError(err).Debug("unable to create search query")
							panic(err)
						}
						metadata, err := p.query.SearchUserMetadata(
							ctx,
							true,
							user.ID,
							&query.UserMetadataSearchQueries{Queries: []query.SearchQuery{resourceOwnerQuery}},
							false,
						)
						if err != nil {
							logging.WithError(err).Info("unable to get md in action")
							panic(err)
						}
						return object.UserMetadataListFromQuery(c, metadata)
					}
				}),
				actions.SetFields("grants", func(c *actions.FieldConfig) interface{} {
					return object.UserGrantsFromQuery(ctx, p.query, c, userGrants)
				}),
			),
			actions.SetFields("org",
				actions.SetFields("getMetadata", func(c *actions.FieldConfig) interface{} {
					return func(goja.FunctionCall) goja.Value {
						return object.GetOrganizationMetadata(ctx, p.query, c, user.ResourceOwner)
					}
				}),
			),
		),
	)

	for _, action := range queriedActions {
		actionCtx, cancel := context.WithTimeout(ctx, action.Timeout())

		apiFields := actions.WithAPIFields(
			actions.SetFields("v1",
				actions.SetFields("attributes",
					actions.SetFields("setCustomAttribute", func(name string, nameFormat string, attributeValue ...string) {
						if _, ok := customAttributes[name]; !ok {
							customAttributes = appendCustomAttribute(customAttributes, name, nameFormat, attributeValue)
							return
						}
					}),
				),
				actions.SetFields("user",
					actions.SetFields("setMetadata", func(call goja.FunctionCall) goja.Value {
						if len(call.Arguments) != 2 {
							panic("exactly 2 (key, value) arguments expected")
						}
						key := call.Arguments[0].Export().(string)
						val := call.Arguments[1].Export()

						value, err := json.Marshal(val)
						if err != nil {
							logging.WithError(err).Debug("unable to marshal")
							panic(err)
						}

						metadata := &domain.Metadata{
							Key:   key,
							Value: value,
						}
						if _, err = p.command.SetUserMetadata(ctx, metadata, user.ID, user.ResourceOwner); err != nil {
							logging.WithError(err).Info("unable to set md in action")
							panic(err)
						}
						return nil
					}),
				),
			),
		)

		err = actions.Run(
			actionCtx,
			ctxFields,
			apiFields,
			action.Script,
			action.Name,
			append(actions.ActionToOptions(action), actions.WithHTTP(actionCtx))...,
		)
		cancel()
		if err != nil {
			return nil, err
		}
	}
	return customAttributes, nil
}

func (p *Storage) getGrants(ctx context.Context, userID, applicationID string) (*query.UserGrants, error) {
	projectID, err := p.query.ProjectIDFromClientID(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(projectID)
	if err != nil {
		return nil, err
	}
	userIDQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, err
	}
	activeQuery, err := query.NewUserGrantStateQuery(domain.UserGrantStateActive)
	if err != nil {
		return nil, err
	}
	return p.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{
			projectQuery,
			userIDQuery,
			activeQuery,
		},
	}, true)
}

type customAttribute struct {
	nameFormat     string
	attributeValue []string
}

func appendCustomAttribute(customAttributes map[string]*customAttribute, name string, nameFormat string, attributeValue []string) map[string]*customAttribute {
	if customAttributes == nil {
		customAttributes = make(map[string]*customAttribute)
	}
	customAttributes[name] = &customAttribute{
		nameFormat:     nameFormat,
		attributeValue: attributeValue,
	}
	return customAttributes
}
