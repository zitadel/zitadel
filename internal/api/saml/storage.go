package saml

import (
	"context"
	"github.com/zitadel/saml/pkg/provider"
	"github.com/zitadel/saml/pkg/provider/key"
	"github.com/zitadel/saml/pkg/provider/models"
	"github.com/zitadel/saml/pkg/provider/serviceprovider"
	"github.com/zitadel/saml/pkg/provider/xml"
	"github.com/zitadel/saml/pkg/provider/xml/samlp"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"time"
)

var _ provider.EntityStorage = &Storage{}
var _ provider.IdentityProviderStorage = &Storage{}
var _ provider.AuthStorage = &Storage{}
var _ provider.UserStorage = &Storage{}

type StorageConfig struct {
}

type Storage struct {
	certChan                   <-chan interface{}
	defaultCertificateLifetime time.Duration

	currentCACertificate       query.Certificate
	currentMetadataCertificate query.Certificate
	currentResponseCertificate query.Certificate

	locker               crdb.Locker
	certificateAlgorithm string
	encAlg               crypto.EncryptionAlgorithm

	eventstore *eventstore.Eventstore
	repo       repository.Repository
	command    *command.Commands
	query      *query.Queries

	defaultLoginURL string
}

func (p *Storage) GetEntityByID(ctx context.Context, entityID string) (*serviceprovider.ServiceProvider, error) {
	app, err := p.query.AppBySAMLEntityID(ctx, entityID)
	if err != nil {
		return nil, err
	}
	metadata := app.SAMLConfig.Metadata

	return serviceprovider.NewServiceProvider(
		app.ID,
		&serviceprovider.ServiceProviderConfig{
			Metadata: metadata,
			URL:      app.SAMLConfig.MetadataURL,
		},
		p.defaultLoginURL,
	)
}

func (p *Storage) GetEntityIDByAppID(ctx context.Context, appID string) (string, error) {
	app, err := p.query.AppByID(ctx, appID)
	if err != nil {
		return "", err
	}
	metadata, err := xml.ParseMetadataXmlIntoStruct([]byte(app.SAMLConfig.Metadata))
	if err != nil {
		return "", err
	}
	return string(metadata.EntityID), nil
}

func (p *Storage) Health(context.Context) error {
	return nil
}

func (p *Storage) GetCA(ctx context.Context) (*key.CertificateAndKey, error) {
	return p.GetCertificateAndKey(ctx, domain.KeyUsageSAMLCA)
}

func (p *Storage) GetMetadataSigningKey(ctx context.Context) (*key.CertificateAndKey, error) {
	return p.GetCertificateAndKey(ctx, domain.KeyUsageSAMLMetadataSigning)
}

func (p *Storage) GetResponseSigningKey(ctx context.Context) (*key.CertificateAndKey, error) {
	return p.GetCertificateAndKey(ctx, domain.KeyUsageSAMLResponseSinging)
}

func (p *Storage) CreateAuthRequest(ctx context.Context, req *samlp.AuthnRequestType, acsUrl, protocolBinding, relayState, applicationID string) (_ models.AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sd436", "no user agent id")
	}

	authRequest := CreateAuthRequestToBusiness(req, acsUrl, protocolBinding, applicationID, relayState, userAgentID)

	resp, err := p.repo.CreateAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}

	return AuthRequestFromBusiness(resp)
}

func (p *Storage) AuthRequestByID(ctx context.Context, id string) (_ models.AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-D3g21", "no user agent id")
	}
	resp, err := p.repo.AuthRequestByIDCheckLoggedIn(ctx, id, userAgentID)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (p *Storage) AuthRequestByCode(ctx context.Context, code string) (_ models.AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	resp, err := p.repo.AuthRequestByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (p *Storage) SetUserinfoWithUserID(ctx context.Context, userinfo models.AttributeSetter, userID string, attributes []int) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	user, err := p.query.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	for _, attribute := range attributes {
		switch attribute {
		case provider.AttributeEmail:
			userinfo.SetEmail(user.Human.Email)
		case provider.AttributeSurname:
			userinfo.SetSurname(user.Human.LastName)
		case provider.AttributeFullName:
			userinfo.SetFullName(user.Human.DisplayName)
		case provider.AttributeGivenName:
			userinfo.SetGivenName(user.Human.FirstName)
		case provider.AttributeUsername:
			userinfo.SetUsername(user.PreferredLoginName)
		case provider.AttributeUserID:
			userinfo.SetUserID(userID)
		}
	}
	if attributes == nil || len(attributes) == 0 {
		userinfo.SetEmail(user.Human.Email)
		userinfo.SetSurname(user.Human.LastName)
		userinfo.SetGivenName(user.Human.FirstName)
		userinfo.SetFullName(user.Human.DisplayName)
		userinfo.SetUsername(user.PreferredLoginName)
		userinfo.SetUserID(userID)
	}
	return nil
}

func (p *Storage) SetUserinfoWithLoginName(ctx context.Context, userinfo models.AttributeSetter, loginName string, attributes []int) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	loginNameSQ, err := query.NewUserLoginNamesSearchQuery(loginName)
	if err != nil {
		return err
	}
	user, err := p.query.GetUser(ctx, loginNameSQ)
	if err != nil {
		return err
	}

	for _, attribute := range attributes {
		switch attribute {
		case provider.AttributeEmail:
			userinfo.SetEmail(user.Human.Email)
		case provider.AttributeSurname:
			userinfo.SetSurname(user.Human.LastName)
		case provider.AttributeFullName:
			userinfo.SetFullName(user.Human.DisplayName)
		case provider.AttributeGivenName:
			userinfo.SetGivenName(user.Human.FirstName)
		case provider.AttributeUsername:
			userinfo.SetUsername(user.PreferredLoginName)
		case provider.AttributeUserID:
			userinfo.SetUserID(user.ID)
		}
	}
	if attributes == nil || len(attributes) == 0 {
		userinfo.SetEmail(user.Human.Email)
		userinfo.SetSurname(user.Human.LastName)
		userinfo.SetGivenName(user.Human.FirstName)
		userinfo.SetFullName(user.Human.DisplayName)
		userinfo.SetUsername(user.PreferredLoginName)
		userinfo.SetUserID(user.ID)
	}
	return nil
}
